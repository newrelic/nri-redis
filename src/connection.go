package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/newrelic/infra-integrations-sdk/log"
)

type conn interface {
	GetInfo() (string, error)
	GetConfig() (map[string]string, error)
	setKeysType(string, []string, map[string]keyInfo) error
	setKeysLength(string, []string, map[string]keyInfo) error
	GetRawCustomKeys(map[string][]string) (map[string]map[string]keyInfo, error)
	Close()
}

type redisConn struct {
	c redis.Conn

	// renamedCommands is the custom Redis command used throughout nri-redis
	// This is used to allow usage of renamed-command in Redis server.
	// Example Redis server config:
	//     rename-command CONFIG "SUPER-SECRET-CONFIG-COMMAND"
	// We will have renamedCommands["CONFIG"] = "SUPER-SECRET-CONFIG-COMMAND"
	renamedCommands map[string]string
}

type keyInfo struct {
	keyLength int64
	keyType   string
}

type configConnectionError struct {
	cause error
}

func (c configConnectionError) Error() string {
	return "can't execute redis 'CONFIG' command: " + c.cause.Error()
}

func newRedisCon(hostname string, port int, unixSocket string, password string, renamedCommands map[string]string) (conn, error) {
	connectTimeout := redis.DialConnectTimeout(time.Second * 5)
	readTimeout := redis.DialReadTimeout(time.Second * 5)
	writeTimeout := redis.DialWriteTimeout(time.Second * 5)
	redisPass := redis.DialPassword(password)

	var c redis.Conn
	var err error

	switch {
	case unixSocket != "":
		c, err = redis.Dial("unix", unixSocket, connectTimeout, readTimeout, writeTimeout, redisPass)
		if err != nil {
			return nil, fmt.Errorf("Redis connection through Unix Socket failed, got error: %v", err)
		}
		log.Debug("Connected to Redis through Unix Socket")
	case hostname != "" && port > 0:
		URL := hostname + ":" + strconv.Itoa(port)
		c, err = redis.Dial("tcp", URL, connectTimeout, readTimeout, writeTimeout, redisPass)
		if err != nil {
			return nil, fmt.Errorf("Redis connection through TCP failed, got error: %v", err)
		}
		log.Debug("Connected to Redis through TCP")
	default:
		return nil, fmt.Errorf("Redis connection failed, cannot connect either through TCP or Unix Socket")
	}
	return redisConn{c, renamedCommands}, nil
}

func (r redisConn) GetInfo() (string, error) {
	if err := r.c.Send(r.getCommand("INFO")); err != nil {
		return "", fmt.Errorf("can't write INFO Redis command: %v", err.Error())
	}
	if err := r.c.Flush(); err != nil {
		return "", fmt.Errorf("can't send INFO Redis command: %v", err.Error())
	}
	return redis.String(r.c.Receive())
}

func (r redisConn) GetConfig() (map[string]string, error) {
	if err := r.c.Send(r.getCommand("CONFIG"), "GET", "*"); err != nil {
		return nil, configConnectionError{cause: err}
	}
	if err := r.c.Flush(); err != nil {
		return nil, configConnectionError{cause: err}
	}
	return redis.StringMap(r.c.Receive())
}

func (r redisConn) Close() {
	r.c.Close()
}

func (r redisConn) setKeysType(db string, keys []string, info map[string]keyInfo) error {

	_, err := r.c.Do(r.getCommand("SELECT"), db)
	if err != nil {
		return fmt.Errorf("Cannot connect to db: %s, information for keys: %v will not be reported, got error: %v ", db, keys, err)
	}

	for _, key := range keys {
		if err = r.c.Send(r.getCommand("TYPE"), key); err != nil {
			log.Warn("Cannot get a type for key: %s, got error: %v", key, err)
		}
	}

	if err = r.c.Flush(); err != nil {
		return fmt.Errorf("Cannot get data for db: %s, got error: %v", db, err)
	}

	for _, key := range keys {
		keyType, err := r.c.Receive()
		if err != nil {
			log.Warn("For db: %s and key: %s, got error: %v", db, key, err)
			continue
		}

		tmp := info[key]
		tmp.keyType = keyType.(string)
		info[key] = tmp
	}

	return nil
}

func (r redisConn) setKeysLength(db string, keys []string, info map[string]keyInfo) error {

	_, err := r.c.Do(r.getCommand("SELECT"), db)
	if err != nil {
		return fmt.Errorf("Cannot connect to db: %s, information for keys: %v will not be reported, got error: %v ", db, keys, err)
	}

	for _, key := range keys {
		switch info[key].keyType {
		case "list":
			if err = r.c.Send(r.getCommand("LLEN"), key); err != nil {
				log.Warn("Cannot retrieve a length of the key: %s from db: %s, got error: %v", key, db, err)
			}
		case "set":
			if err = r.c.Send(r.getCommand("SCARD"), key); err != nil {
				log.Warn("Cannot retrieve a length of the key: %s from db: %s, got error: %v", key, db, err)
			}
		case "zset":
			if err = r.c.Send(r.getCommand("ZCOUNT"), key, "-inf", "+inf"); err != nil {
				log.Warn("Cannot retrieve a length of the key: %s from db: %s, got error: %v", key, db, err)
			}
		case "hash":
			if err = r.c.Send(r.getCommand("HLEN"), key); err != nil {
				log.Warn("Cannot retrieve a length of the key: %s from db: %s, got error: %v", key, db, err)
			}
		case "string":
			log.Warn("Key: %s from db: %s is a string type, cannot retrieve a length", key, db)
		default:
			log.Warn("Unknown type of the key: %s from db: %s, cannot retrieve a length", key, db)
		}
	}

	if err = r.c.Flush(); err != nil {
		return fmt.Errorf("Cannot get data for db: %s, got error: %v", db, err)
	}

	for _, key := range keys {
		if info[key].keyType != "string" && info[key].keyType != "none" {
			keyLength, err := r.c.Receive()
			if err != nil {
				log.Warn("For db: %s and key: %s, got error: %v", db, key, err)
				continue
			}
			tmp := info[key]
			tmp.keyLength = keyLength.(int64)
			info[key] = tmp
		} else {
			delete(info, key)
		}
	}
	return nil
}

func (r redisConn) GetRawCustomKeys(databaseKeys map[string][]string) (map[string]map[string]keyInfo, error) {
	customKeysMetric := make(map[string]map[string]keyInfo)

	for db, keys := range databaseKeys {
		info := make(map[string]keyInfo)

		err := r.setKeysType(db, keys, info)
		if err != nil {
			return nil, fmt.Errorf("Cannot get type for keys %s from db %s, got err: %v", keys, db, err)
		}
		err = r.setKeysLength(db, keys, info)
		if err != nil {
			return nil, fmt.Errorf("Cannot get length for keys %s from db %s, got err: %v", keys, db, err)
		}

		customKeysMetric["db"+db] = info
	}

	return customKeysMetric, nil
}

// getCommand returns alias of 'command' if exists
// Redis server security model supports disabling/renaming of specific commands
// Ref: https://redis.io/topics/security
func (r redisConn) getCommand(command string) string {
	if alias, ok := r.renamedCommands[command]; ok {
		return alias
	}
	return command
}
