package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigConnectionError(t *testing.T) {
	// a failed GetConfig connection must return configConnectionError
	c := redisConn{c: fakeConn{sendErr: errors.New("this is a fake")}}
	_, err := c.GetConfig()
	assert.Error(t, err)
	assert.IsType(t, configConnectionError{}, err)
}

func TestConfigReceiveError(t *testing.T) {
	// a successful GetConfig connection that fails in the data deserialization
	// MUST NOT return configConnectionError
	c := redisConn{c: fakeConn{}}
	_, err := c.GetConfig()
	assert.Error(t, err)
	_, ok := err.(configConnectionError)
	assert.False(t, ok, "the returned error should not be configConnectionError")
}

type fakeConn struct {
	sendErr error
}

func (f fakeConn) Send(_ string, _ ...interface{}) error {
	return f.sendErr
}

func (f fakeConn) Close() error                                       { return nil }
func (f fakeConn) Err() error                                         { return nil }
func (f fakeConn) Do(_ string, _ ...interface{}) (interface{}, error) { return nil, nil }
func (f fakeConn) Flush() error                                       { return nil }
func (f fakeConn) Receive() (interface{}, error)                      { return nil, nil }
