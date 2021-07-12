package main

import (
	"errors"
	"fmt"
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

func Test_redisConn_command(t *testing.T) {
	renamedCommands := make(map[string]string)
	renamedCommands["NON-RENAMED-COMMAND"] = "NON-RENAMED-COMMAND"
	renamedCommands["RENAMED-CONFIG"] = "NEW-RENAMED-CONFIG"
	renamedCommands["DISABLED-COMMAND"] = ""

	c := redisConn{c: fakeConn{}, renamedCommands: renamedCommands}

	type fields struct {
		c redisConn
	}
	type args struct {
		command string
	}
	type testcase struct {
		name   string
		fields fields
		args   args
		want   string
	}
	// Test cases
	tests := []testcase{}
	for cmd, alias := range renamedCommands {
		tests = append(tests, testcase{
			fmt.Sprintf("rename-command %v %v", cmd, alias),
			fields{c},
			args{cmd},
			alias,
		})
	}
	// Extra test case where the command is not specified the renamedCommands
	tests = append(tests, testcase{
		"rename-command UNSPECIFIED UNSPECIFIED",
		fields{c},
		args{"UNSPECIFIED"},
		"UNSPECIFIED",
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.c.command(tt.args.command); got != tt.want {
				t.Errorf("redisConn.command() = %v, want %v", got, tt.want)
			}
		})
	}
}
