package main

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/newrelic/infra-integrations-sdk/log"
)

const defaultTimeout = time.Second * 5

type redisConn struct {
	c redis.Conn

	// renamedCommands is the renamed-version of Redis commands used throughout nri-redis
	// This is used to allow usage of 'renamed-command' in Redis server.
	// Example Redis server config:
	//     rename-command CONFIG "SUPER-SECRET-CONFIG-COMMAND"
	// We will have renamedCommands["CONFIG"] = "SUPER-SECRET-CONFIG-COMMAND"
	// Ref: https://redis.io/topics/security
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

func newSocketRedisCon(unixSocket string, options ...redis.DialOption) (*redisConn, error) {
	c, err := redis.Dial("unix", unixSocket, options...)
	if err != nil {
		return nil, fmt.Errorf("connecting through Unix Socket: %w\", err", err)
	}
	log.Debug("Connected to Redis through Unix Socket %s", unixSocket)
	return &redisConn{c, nil}, nil
}

func newNetworkRedisCon(redisURL string, options ...redis.DialOption) (*redisConn, error) {
	c, err := redis.Dial("tcp", redisURL, options...)
	if err != nil {
		return nil, fmt.Errorf("connecting through TCP: %w", err)
	}
	log.Debug("Connected to Redis through TCP %s", redisURL)

	return &redisConn{c, nil}, nil
}

func standardDialOptions(username string, password string) []redis.DialOption {
	return []redis.DialOption{
		redis.DialConnectTimeout(defaultTimeout),
		redis.DialReadTimeout(defaultTimeout),
		redis.DialWriteTimeout(defaultTimeout),
		redis.DialUsername(username),
		redis.DialPassword(password),
	}
}

func tlsDialOptions(useTLS bool, tlsInsecureSkipVerify bool) []redis.DialOption {
	return []redis.DialOption{
		redis.DialTLSHandshakeTimeout(defaultTimeout),
		redis.DialUseTLS(useTLS),
		redis.DialTLSSkipVerify(tlsInsecureSkipVerify),
	}
}

func (r redisConn) GetInfo() (string, error) {
	if err := r.c.Send(r.command("INFO")); err != nil {
		return "", fmt.Errorf("can't write INFO Redis command: %v", err.Error())
	}
	if err := r.c.Flush(); err != nil {
		return "", fmt.Errorf("can't send INFO Redis command: %v", err.Error())
	}
	return redis.String(r.c.Receive())
}

func (r redisConn) GetConfig() (map[string]string, error) {
	if err := r.c.Send(r.command("CONFIG"), "GET", "*"); err != nil {
		return nil, configConnectionError{cause: err}
	}
	if err := r.c.Flush(); err != nil {
		return nil, configConnectionError{cause: err}
	}
	return redis.StringMap(r.c.Receive())
}

// RenameCommands will populate internal renamedCommands mapping
func (r *redisConn) RenameCommands(renamedCommands map[string]string) {
	r.renamedCommands = renamedCommands
}

func (r redisConn) Close() error {
	return r.c.Close()
}

func (r redisConn) setKeysType(db string, keys []string, info map[string]keyInfo) error {
	_, err := r.c.Do(r.command("SELECT"), db)
	if err != nil {
		return fmt.Errorf("Cannot connect to db: %s, information for keys: %v will not be reported, got error: %v ", db, keys, err)
	}

	for _, key := range keys {
		if err = r.c.Send(r.command("TYPE"), key); err != nil {
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
	_, err := r.c.Do(r.command("SELECT"), db)
	if err != nil {
		return fmt.Errorf("Cannot connect to db: %s, information for keys: %v will not be reported, got error: %v ", db, keys, err)
	}

	for _, key := range keys {
		switch info[key].keyType {
		case "list":
			if err = r.c.Send(r.command("LLEN"), key); err != nil {
				log.Warn("Cannot retrieve a length of the key: %s from db: %s, got error: %v", key, db, err)
			}
		case "set":
			if err = r.c.Send(r.command("SCARD"), key); err != nil {
				log.Warn("Cannot retrieve a length of the key: %s from db: %s, got error: %v", key, db, err)
			}
		case "zset":
			if err = r.c.Send(r.command("ZCOUNT"), key, "-inf", "+inf"); err != nil {
				log.Warn("Cannot retrieve a length of the key: %s from db: %s, got error: %v", key, db, err)
			}
		case "hash":
			if err = r.c.Send(r.command("HLEN"), key); err != nil {
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

// command returns Redis command that should be used by nri-redis
// Supports:
//   - Renamed version of 'command' if exists
func (r redisConn) command(command string) string {
	if renamedCommand, ok := r.renamedCommands[command]; ok {
		return renamedCommand
	}
	return command
}
