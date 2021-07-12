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

	databaseKey := getDbAndKeys(*keysFlag)
	if !reflect.DeepEqual(databaseKey, expectedValue) {
		t.Error()
	}
}

func TestGetDbAndKeysDb0(t *testing.T) {
	keysFlag := sdkArgs.NewJSON([]interface{}{"key1", "key2"})
	expectedValue := map[string][]string{
		"0": {"key1", "key2"},
	}
	databaseKey := getDbAndKeys(*keysFlag)
	if !reflect.DeepEqual(databaseKey, expectedValue) {
		t.Error()
	}
}

func TestGetDbAndKeysEmpty(t *testing.T) {
	keysFlag := sdkArgs.NewJSON(nil)

	databaseKey := getDbAndKeys(*keysFlag)
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
