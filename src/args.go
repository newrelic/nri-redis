package main

import (
	"fmt"

	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/log"
)

func getDBAndKeys(keysFlag sdkArgs.JSON) map[string][]string {
	databaseKeys := make(map[string][]string)

	convertF := func(listInterface []interface{}) []string {
		listString := make([]string, 0)
		for _, elem := range listInterface {
			if _, ok := elem.(string); ok {
				listString = append(listString, elem.(string))
			} else {
				log.Warn("Not expected type for key: %v, string type required", elem)
			}
		}
		return listString
	}

	switch source := keysFlag.Get().(type) {
	case []interface{}:
		databaseKeys["0"] = removeDuplicates(convertF(source))
	case map[string]interface{}:
		for db, keys := range source {
			databaseKeys[db] = removeDuplicates(convertF(keys.([]interface{})))
		}
	default:
		log.Warn("Invalid format, value of keys flag: %v ", keysFlag)
	}
	return databaseKeys
}

// getRenamedCommands returns a map containing command and renamed command pairs
// Example flag value: '{"CONFIG": "ZmtlbmZ3ZWZl-CONFIG", "ANOTHER-COMMAND": ""}'
func getRenamedCommands(renamedCommandsFlag sdkArgs.JSON) (map[string]string, error) {
	renamedCommands := make(map[string]string)

	convertF := func(stringInterface interface{}) (string, error) {
		if str, ok := stringInterface.(string); ok {
			return str, nil
		}
		return "", fmt.Errorf("Unexpected type %T, value must be a string", stringInterface)
	}

	switch source := renamedCommandsFlag.Get().(type) {
	case map[string]interface{}:
		for cmd, alias := range source {
			if renamedCmd, err := convertF(alias); err == nil {
				renamedCommands[cmd] = renamedCmd
			} else {
				log.Warn("Invalid format, value of a renamed command: %v=%v", cmd, alias)
			}
		}
	default:
		return renamedCommands, fmt.Errorf("Invalid format, value of renamed commands flag: %v", renamedCommandsFlag)
	}

	return renamedCommands, nil
}

func removeDuplicates(elements []string) []string {
	found := map[string]struct{}{}
	result := []string{}

	for v := range elements {
		if _, ok := found[elements[v]]; !ok {
			found[elements[v]] = struct{}{}
			result = append(result, elements[v])
		}
	}
	return result
}

func validateKeysFlag(databaseKeys map[string][]string, keysLimit int) (int, error) {
	keysNumber := 0
	for _, keys := range databaseKeys {
		keysNumber += len(keys)
	}
	if keysNumber > keysLimit {
		return keysNumber, fmt.Errorf("Keys limit was exceeded; keys limit: %d, keys number: %d", keysNumber, keysLimit)
	}
	return keysNumber, nil
}
