package main

import (
	"reflect"
	"testing"

	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/stretchr/testify/assert"
)

func TestGetDbAndKeys(t *testing.T) {
	keysFlag := sdkArgs.NewJSON(map[string]interface{}{"0": []interface{}{"key1"}, "2": []interface{}{"key2", "key3"}})
	expectedValue := map[string][]string{
		"0": {"key1"},
		"2": {"key2", "key3"},
	}

	databaseKey := getDBAndKeys(*keysFlag)
	if !reflect.DeepEqual(databaseKey, expectedValue) {
		t.Error()
	}
}

func TestGetDbAndKeysDb0(t *testing.T) {
	keysFlag := sdkArgs.NewJSON([]interface{}{"key1", "key2"})
	expectedValue := map[string][]string{
		"0": {"key1", "key2"},
	}
	databaseKey := getDBAndKeys(*keysFlag)
	if !reflect.DeepEqual(databaseKey, expectedValue) {
		t.Error()
	}
}

func TestGetDbAndKeysEmpty(t *testing.T) {
	keysFlag := sdkArgs.NewJSON(nil)

	databaseKey := getDBAndKeys(*keysFlag)
	if len(databaseKey) != 0 {
		t.Error()
	}
}

func TestRemoveDuplicates(t *testing.T) {
	elements := []string{"k1", "k2", "k1", "k2", "k2", "k3"}
	expectedElements := []string{"k1", "k2", "k3"}

	result := removeDuplicates(elements)
	if !reflect.DeepEqual(result, expectedElements) {
		t.Error()
	}
}

func TestRemoveDuplicatesNoDuplicates(t *testing.T) {
	elements := []string{"k1", "k2", "k3"}
	expectedElements := []string{"k1", "k2", "k3"}

	result := removeDuplicates(elements)
	if !reflect.DeepEqual(result, expectedElements) {
		t.Error()
	}
}

func TestValidateKeysFlag(t *testing.T) {
	databaseKeys := map[string][]string{
		"0": {"key1", "key2", "key3"},
	}
	keysLimit := 3
	expectedKeysNumber := 3

	keysNumber, err := validateKeysFlag(databaseKeys, keysLimit)
	assert.NoError(t, err)

	if keysNumber != expectedKeysNumber {
		t.Error()
	}
}

func TestValidateKeysFlagExceedLimit(t *testing.T) {
	databaseKeys := map[string][]string{
		"0": {"key1", "key2", "key3"},
	}
	keysLimit := 2
	expectedKeysNumber := 3

	keysNumber, err := validateKeysFlag(databaseKeys, keysLimit)

	if reflect.DeepEqual(err.Error(), "Keys limit was exeeded; keys limit: 2, keys number: %3") {
		t.Error()
	}
	if keysNumber != expectedKeysNumber {
		t.Error()
	}
}

func Test_getRenamedCommands(t *testing.T) {
	type args struct {
		renamedCommandsFlag sdkArgs.JSON
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			"An empty renamed command should be fine",
			args{*sdkArgs.NewJSON(map[string]interface{}{})},
			map[string]string{},
			false,
		},
		{
			"Rename one command",
			args{*sdkArgs.NewJSON(map[string]interface{}{"CONFIG": "ZmtlbmZ3ZWZl-CONFIG"})},
			map[string]string{"CONFIG": "ZmtlbmZ3ZWZl-CONFIG"},
			false,
		},
		{
			"Rename multiple commands",
			args{
				*sdkArgs.NewJSON(map[string]interface{}{
					"NON-RENAMED-COMMAND":    "NON-RENAMED-COMMAND",
					"RENAMED-CONFIG":         "NEW-RENAMED-CONFIG",
					"ANOTHER-RENAMED-CONFIG": "ANOTHER-RENAMED-CONFIG",
				},
				),
			},
			map[string]string{
				"NON-RENAMED-COMMAND":    "NON-RENAMED-COMMAND",
				"RENAMED-CONFIG":         "NEW-RENAMED-CONFIG",
				"ANOTHER-RENAMED-CONFIG": "ANOTHER-RENAMED-CONFIG",
			},
			false,
		},
		{
			"An invalid renamed command value should be skipped",
			args{*sdkArgs.NewJSON(map[string]interface{}{"RENAMED-COMMAND": 1})},
			map[string]string{},
			false,
		},
		{
			"An invalid renamed command key should return error",
			args{*sdkArgs.NewJSON(map[int]interface{}{1: "RENAMED-COMMAND"})},
			map[string]string{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getRenamedCommands(tt.args.renamedCommandsFlag)
			if (err != nil) != tt.wantErr {
				t.Errorf("getRenamedCommands() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getRenamedCommands() = %v, want %v", got, tt.want)
			}
		})
	}
}
