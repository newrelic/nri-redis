package main

import (
	"regexp"
	"strings"

	"github.com/newrelic/infra-integrations-sdk/data/inventory"

	"github.com/newrelic/infra-integrations-sdk/log"
)

func getRawInventory(config map[string]string, metrics map[string]interface{}) map[string]interface{} {
	rawInventory := make(map[string]interface{})
	for key, value := range config {
		rawInventory[key] = value
	}

	if _, ok := metrics["redis_version"]; ok {
		rawInventory["redis_version"] = metrics["redis_version"]
	}
	if _, ok := metrics["executable"]; ok {
		rawInventory["executable"] = metrics["executable"]
	}
	if _, ok := metrics["config_file"]; ok {
		rawInventory["config-file"] = metrics["config_file"]
	}
	if _, ok := metrics["mem_allocator"]; ok {
		rawInventory["mem-allocator"] = metrics["mem_allocator"]
	}

	if len(rawInventory) == 0 {
		log.Debug("Empty result for inventory")
	}
	return rawInventory
}

func populateInventory(inventory *inventory.Inventory, rawInventory map[string]interface{}) {
	re, _ := regexp.Compile("(?i)(requirepass|masterauth)")

	for key, value := range rawInventory {
		if re.MatchString(key) {
			value = "(omitted value)"
		}
		if err := inventory.SetItem(key, "value", value); err != nil {
			log.Error("fail to set inventory item: %s", err)
		}
	}

	setInventoryClientBufferValue(inventory)
	setInventorySaveValue(inventory)
}

func setInventorySaveValue(inventory *inventory.Inventory) {
	if save, ok := inventory.Items()["save"]["value"]; ok {
		if save != "" {
			delete(inventory.Items()["save"], "value")
			if err := inventory.SetItem("save", "raw-value", save); err != nil {
				log.Error("fail to set inventory item: %s", err)
			}
			saveSplited := strings.Split(save.(string), " ")
			for i := 0; i <= len(saveSplited)-1; i += 2 {
				if err := inventory.SetItem("save", "after-"+saveSplited[i]+"-seconds", "at-least-"+saveSplited[i+1]+"-key-changes"); err != nil {
					log.Error("fail to set inventory item: %s", err)
				}
			}
		}
	} else {
		log.Debug("Key \"save\" is not present")
	}
}

func setInventoryClientBufferValue(inventory *inventory.Inventory) {
	if clientBuffer, ok := inventory.Items()["client-output-buffer-limit"]["value"]; ok {
		delete(inventory.Items()["client-output-buffer-limit"], "value")
		if err := inventory.SetItem("client-output-buffer-limit", "raw-value", clientBuffer); err != nil {
			log.Error("fail to set inventory item: %s", err)
		}

		clientBufferSplited := strings.Split(clientBuffer.(string), " ")
		for i := 0; i <= len(clientBufferSplited)-3; i += 4 {
			if err := inventory.SetItem("client-output-buffer-limit", clientBufferSplited[i]+"-hard-limit", clientBufferSplited[i+1]); err != nil {
				log.Error("fail to set inventory item: %s", err)
			}
			if err := inventory.SetItem("client-output-buffer-limit", clientBufferSplited[i]+"-soft-limit", clientBufferSplited[i+2]); err != nil {
				log.Error("fail to set inventory item: %s", err)
			}
			if err := inventory.SetItem("client-output-buffer-limit", clientBufferSplited[i]+"-soft-seconds", clientBufferSplited[i+3]); err != nil {
				log.Error("fail to set inventory item: %s", err)
			}
		}
	} else {
		log.Debug("Key \"client-output-buffer-limit\" is not present")
	}
}
