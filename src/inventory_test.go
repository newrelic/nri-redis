package main

import (
	"github.com/newrelic/infra-integrations-sdk/data/inventory"
	"testing"
)

func TestGetRawInventory(t *testing.T) {
	config := map[string]string{
		"dbfilename":  "dump.rdb",
		"requirepass": "",
		"logfile":     "/var/log/redis/redis.log",
	}

	metrics := map[string]interface{}{
		"redis_version": "3.2.3",
		"executable":    "/usr/bin/redis-server",
		"config_file":   "/etc/redis.conf",
		"mem_allocator": "jemalloc-3.6.0",
	}
	expectedLength := 7
	expectedDbfilename := "dump.rdb"
	expectedRequirepass := ""
	expectedMemAllocator := "jemalloc-3.6.0"
	rawInventory := getRawInventory(config, metrics)

	if len(rawInventory) != expectedLength {
		t.Errorf("Not all values processed, got %d length of the rawInventory, expected: %d", len(rawInventory), expectedLength)
	}
	if rawInventory["dbfilename"] != expectedDbfilename {
		t.Errorf("For key \"dbfilename\" wrong value returned, got: %s, expected: %s", rawInventory["dbfilename"], expectedDbfilename)
	}
	if rawInventory["requirepass"] != expectedRequirepass {
		t.Errorf("For key \"requirepass\" wrong value returned, got: %s, expected: %s", rawInventory["requirepass"], expectedRequirepass)
	}
	if rawInventory["mem-allocator"] != expectedMemAllocator {
		t.Errorf("For key \"mem_allocator\" wrong value returned, got: %s, expected: %s", rawInventory["mem-allocator"], expectedMemAllocator)
	}
}

func TestGetRawInventoryEmpty(t *testing.T) {
	config := map[string]string{}
	metrics := map[string]interface{}{}

	rawInventory := getRawInventory(config, metrics)
	expectedLength := 0

	if len(rawInventory) != expectedLength {
		t.Errorf("rawInventory not empty, got %d length of the rawInventory, expected: %d", len(rawInventory), expectedLength)
	}
}

func TestPopulateInventory(t *testing.T) {
	inventory := inventory.New()
	rawInventory := map[string]interface{}{
		"redis_version": "3.2.3",
		"requirepass":   "",
		"save":          "900 1 300 10 60 10000",
		"client-output-buffer-limit": "normal 0 0 0 slave 268435456 67108864 60 pubsub 33554432 8388608 60",
	}
	expectedRedisVersion := "3.2.3"
	expectedRequirePass := "(omitted value)"
	expectedLength := 4
	expectedSaveItemsLength := 4
	expectedBufferItemsLength := 10

	populateInventory(inventory, rawInventory)
	if len(inventory.Items()) != expectedLength {
		t.Errorf("Not all values processed, got %d length of the Inventory, expected: %d", len(inventory.Items()), expectedLength)
	}
	if len(inventory.Items()["save"]) != expectedSaveItemsLength {
		t.Errorf("Not all values processed for Inventory Save, got %d, expected: %d", len(inventory.Items()["save"]), expectedSaveItemsLength)
	}
	if len(inventory.Items()["client-output-buffer-limit"]) != expectedBufferItemsLength {
		t.Errorf("Not all values processed for Inventory buffer, got %d, expected: %d", len(inventory.Items()["client-output-buffer-limit"]), expectedBufferItemsLength)
	}
	if inventory.Items()["redis_version"]["value"] != expectedRedisVersion {
		t.Errorf("For key \"redis_version\" wrong value returned, got: %s, expected: %s", inventory.Items()["redis_version"]["value"], expectedRedisVersion)
	}
	if inventory.Items()["requirepass"]["value"] != expectedRequirePass {
		t.Errorf("For key \"requirepass\" wrong value returned, got: %s, expected: %s", inventory.Items()["requirepass"]["value"], expectedRequirePass)
	}
}

func TestPopulateInventorySaveEmptyAndBuffer(t *testing.T) {
	inventory := inventory.New()
	rawInventory := map[string]interface{}{
		"save": "",
	}
	expectedSave := ""
	expectedLength := 1

	populateInventory(inventory, rawInventory)
	if len(inventory.Items()) != expectedLength {
		t.Error()
	}
	if len(inventory.Items()["save"]) != 1 {
		t.Error()
	}
	if inventory.Items()["save"]["value"] != expectedSave {
		t.Errorf("For key \"save\" wrong value returned, got: %s, expected: %s", inventory.Items()["save"]["value"], expectedSave)
	}
}

func TestSetInventorySaveValue(t *testing.T) {
	inventory := inventory.New()
	inventory.SetItem("save", "value", "900 1 300 10 60 10000")

	expectedLength := 1
	expectedItemsLength := 4
	expectedRawValue := "900 1 300 10 60 10000"
	expectedAfter900Seconds := "at-least-1-key-changes"
	expectedAfter300Seconds := "at-least-10-key-changes"
	expectedAfter60Seconds := "at-least-10000-key-changes"

	setInventorySaveValue(inventory)
	if len(inventory.Items()) != expectedLength {
		t.Error()
	}
	if len(inventory.Items()["save"]) != expectedItemsLength {
		t.Error()
	}
	if inventory.Items()["save"]["raw-value"] != expectedRawValue {
		t.Errorf("For key [\"save\"][\"raw-value\"] wrong value returned, got: %s, expected: %s", inventory.Items()["save"]["raw-value"], expectedRawValue)
	}
	if inventory.Items()["save"]["after-900-seconds"] != expectedAfter900Seconds {
		t.Errorf("For key [\"save\"][\"after-900-seconds\"] wrong value returned, got: %s, expected: %s", inventory.Items()["save"]["after-900-seconds"], expectedAfter900Seconds)
	}
	if inventory.Items()["save"]["after-300-seconds"] != expectedAfter300Seconds {
		t.Errorf("For key [\"save\"][\"after-300-seconds\"] wrong value returned, got: %s, expected: %s", inventory.Items()["save"]["after-300-seconds"], expectedAfter300Seconds)
	}
	if inventory.Items()["save"]["after-60-seconds"] != expectedAfter60Seconds {
		t.Errorf("For key [\"save\"][\"after-60-seconds\"] wrong value returned, got: %s, expected: %s", inventory.Items()["save"]["after-60-seconds"], expectedAfter60Seconds)
	}
}

func TestSetInventorySaveEmptyValue(t *testing.T) {
	inventory := inventory.New()
	inventory.SetItem("save", "value", "")

	expectedLength := 1
	expectedItemsLength := 1
	expectedSave := ""

	setInventorySaveValue(inventory)
	if len(inventory.Items()) != expectedLength {
		t.Error()
	}
	if len(inventory.Items()["save"]) != expectedItemsLength {
		t.Error()
	}
	if inventory.Items()["save"]["value"] != expectedSave {
		t.Errorf("For key \"save\" wrong value returned, got: %s, expected: %s", inventory.Items()["save"]["value"], expectedSave)
	}
}

func TestSetInventoryClientBufferValue(t *testing.T) {
	inventory := inventory.New()
	inventory.SetItem("client-output-buffer-limit", "value", "normal 0 0 0 slave 268435456 67108864 60 pubsub 33554432 8388608 60")

	expectedLength := 1
	expectedItemsLength := 10
	expectedRawValue := "normal 0 0 0 slave 268435456 67108864 60 pubsub 33554432 8388608 60"
	expectedNormalHardLimit := "0"
	expectedNormalSoftLimit := "0"
	expectedNormalSoftSeconds := "0"

	setInventoryClientBufferValue(inventory)
	if len(inventory.Items()) != expectedLength {
		t.Error()
	}
	if len(inventory.Items()["client-output-buffer-limit"]) != expectedItemsLength {
		t.Error()
	}
	if _, ok := inventory.Items()["client-output-buffer-limit"]["value"]; ok {
		t.Error()
	}
	if inventory.Items()["client-output-buffer-limit"]["raw-value"] != expectedRawValue {
		t.Errorf("For key [\"client-output-buffer-limit\"][\"raw-value\"] wrong value returned, got: %s, expected: %s", inventory.Items()["client-output-buffer-limit"]["raw-value"], expectedRawValue)
	}
	if inventory.Items()["client-output-buffer-limit"]["normal-hard-limit"] != expectedNormalHardLimit {
		t.Errorf("For key [\"client-output-buffer-limit\"][\"normal-hard-limit\"] wrong value returned, got: %s, expected: %s", inventory.Items()["client-output-buffer-limit"]["normal-hard-limit"], expectedNormalHardLimit)
	}
	if inventory.Items()["client-output-buffer-limit"]["normal-soft-limit"] != expectedNormalSoftLimit {
		t.Errorf("For key [\"client-output-buffer-limit\"][\"normal-soft-limit\"] wrong value returned, got: %s, expected: %s", inventory.Items()["client-output-buffer-limit"]["normal-soft-limit"], expectedNormalSoftLimit)
	}
	if inventory.Items()["client-output-buffer-limit"]["normal-soft-seconds"] != expectedNormalSoftSeconds {
		t.Errorf("For key [\"client-output-buffer-limit\"][\"normal-soft-seconds\"] wrong value returned, got: %s, expected: %s", inventory.Items()["client-output-buffer-limit"]["normal-soft-seconds"], expectedNormalSoftSeconds)
	}
}

func TestSetInventoryClientBufferNotPresent(t *testing.T) {
	inventory := inventory.New()
	inventory.SetItem("other-key", "value", "normal 0 0 0 slave 268435456 67108864 60 pubsub 33554432 8388608 60")
	expectedItemsLength := 0

	setInventoryClientBufferValue(inventory)
	if len(inventory.Items()["client-output-buffer-limit"]) != expectedItemsLength {
		t.Error()
	}
}
