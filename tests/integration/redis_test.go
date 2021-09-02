// +build integration

package integration

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/newrelic/nri-redis/tests/integration/helpers"
	"github.com/newrelic/nri-redis/tests/integration/jsonschema"
)

var (
	//default values
	defaultContainer = "integration_nri-redis_1"
	defaultBinPath   = "/nri-redis"
	defaultHost      = "redis"
	defaultPort      = 6379

	// cli flags
	container = flag.String("container", defaultContainer, "container where the integration is installed")
	binPath   = flag.String("bin", defaultBinPath, "Integration binary path")
)

func runIntegration(t *testing.T, host string, port int, enableTLS bool, username string, password string, envVars ...string) (string, string) {
	t.Helper()

	command := make([]string, 0)
	command = append(command, *binPath)
	command = append(command, "--hostname", host)
	command = append(command, "--port", strconv.Itoa(port))
	command = append(command, fmt.Sprintf("--use_tls=%s", strconv.FormatBool(enableTLS)))
	command = append(command, fmt.Sprintf("--tls_insecure_skip_verify=%s", strconv.FormatBool(enableTLS)))
	if password != "" {
		command = append(command, "--password", password)
		if username != "" {
			command = append(command, "--username", username)
		}
	}

	stdout, stderr, err := helpers.ExecInContainer(*container, command, envVars...)

	if stderr != "" {
		log.Debug("Integration command Standard Error: ", stderr)
	}
	require.NoError(t, err, stdout)

	return stdout, stderr
}

func TestMain(m *testing.M) {
	flag.Parse()
	log.SetupLogging(true)

	result := m.Run()

	os.Exit(result)
}

func TestRedisIntegration(t *testing.T) {
	testName := t.Name()

	envVars := []string{fmt.Sprintf("NRIA_CACHE_PATH=/tmp/%v.json", testName)}
	stdout, stderr := runIntegration(t, defaultHost, defaultPort, false, "", "", envVars...)

	schemaPath := filepath.Join("json-schema-files", "redis-schema.json")

	err := jsonschema.Validate(schemaPath, stdout)

	assert.NoError(t, err, "The output of Redis integration doesn't have expected format")
	assert.NotNil(t, stderr, "unexpected stderr")
}

func TestRedisIntegration_With_TLS(t *testing.T) {
	testName := t.Name()

	envVars := []string{fmt.Sprintf("NRIA_CACHE_PATH=/tmp/%v.json", testName)}
	stdout, stderr := runIntegration(t, defaultHost+"-tls", defaultPort, true, "", "", envVars...)

	schemaPath := filepath.Join("json-schema-files", "redis-schema.json")

	err := jsonschema.Validate(schemaPath, stdout)

	assert.NoError(t, err, "The output of Redis integration doesn't have expected format")
	assert.NotNil(t, stderr, "unexpected stderr")
}

func TestRedisIntegration_With_ACL(t *testing.T) {
	testName := t.Name()

	envVars := []string{fmt.Sprintf("NRIA_CACHE_PATH=/tmp/%v.json", testName)}
	stdout, stderr := runIntegration(t, defaultHost+"-acl", defaultPort, false, "testUsername", "testPassword", envVars...)

	schemaPath := filepath.Join("json-schema-files", "redis-schema.json")

	err := jsonschema.Validate(schemaPath, stdout)

	assert.NoError(t, err, "The output of Redis integration doesn't have expected format")
	assert.NotNil(t, stderr, "unexpected stderr")
}

func TestRedisIntegration_WithRemoteEntity(t *testing.T) {
	testName := t.Name()

	envVars := []string{fmt.Sprintf("NRIA_CACHE_PATH=/tmp/%v.json", testName), "REMOTE_MONITORING=true"}
	stdout, stderr := runIntegration(t, defaultHost, defaultPort, false, "", "", envVars...)

	schemaPath := filepath.Join("json-schema-files", "redis-schema-remote-entity.json")

	err := jsonschema.Validate(schemaPath, stdout)

	assert.NoError(t, err, "The output of Redis integration doesn't have expected format")
	assert.NotNil(t, stderr, "unexpected stderr")
}

func TestRedisIntegration_OnlyMetrics(t *testing.T) {
	testName := t.Name()

	envVars := []string{fmt.Sprintf("NRIA_CACHE_PATH=/tmp/%v.json", testName)}
	stdout, stderr := runIntegration(t, defaultHost, defaultPort, false, "", "", envVars...)

	schemaPath := filepath.Join("json-schema-files", "redis-schema-metrics.json")

	err := jsonschema.Validate(schemaPath, stdout)

	assert.NoError(t, err, "The output of Redis integration doesn't have expected format")
	assert.NotNil(t, stderr, "unexpected stderr")
}

func TestRedisIntegration_OnlyInventory(t *testing.T) {
	testName := t.Name()

	envVars := []string{fmt.Sprintf("NRIA_CACHE_PATH=/tmp/%v.json", testName)}
	stdout, stderr := runIntegration(t, defaultHost, defaultPort, false, "", "", envVars...)

	schemaPath := filepath.Join("json-schema-files", "redis-schema-inventory.json")

	err := jsonschema.Validate(schemaPath, stdout)

	assert.NoError(t, err, "The output of Redis integration doesn't have expected format")
	assert.NotNil(t, stderr, "unexpected stderr")

	// Verify that the CONFIG inventory is retrieved
	inventory := extractInventory(t, stdout)
	removeCommonInventory(inventory)
	require.NotZerof(t, len(inventory), "inventory should include CONFIG data. Was %#v", inventory)
}

func TestRedisIntegration_SkipConfig(t *testing.T) {
	testName := t.Name()

	envVars := []string{fmt.Sprintf("NRIA_CACHE_PATH=/tmp/%v.json", testName), "CONFIG_INVENTORY=false"}
	stdout, _ := runIntegration(t, defaultHost, defaultPort, false, "", "", envVars...)

	// Verify that the CONFIG inventory is NOT retrieved
	inventory := extractInventory(t, stdout)
	removeCommonInventory(inventory)
	require.Zerof(t, len(inventory), "inventory should NOT include CONFIG data. Was %#v", inventory)
}

func extractInventory(t *testing.T, jsonOutput string) map[string]interface{} {
	t.Helper()
	jsonMap := map[string]interface{}{}
	require.NoError(t, json.Unmarshal([]byte(jsonOutput), &jsonMap))
	require.Contains(t, jsonMap, "data")
	require.IsType(t, []interface{}{}, jsonMap["data"])
	require.Len(t, jsonMap["data"], 1)
	require.IsType(t, map[string]interface{}{}, jsonMap["data"].([]interface{})[0])
	require.Contains(t, jsonMap["data"].([]interface{})[0].(map[string]interface{}), "inventory")
	return jsonMap["data"].([]interface{})[0].(map[string]interface{})["inventory"].(map[string]interface{})
}

// removes inventory data that is not result of the CONFIG command
func removeCommonInventory(inventory map[string]interface{}) {
	delete(inventory, "redis_version")
	delete(inventory, "executable")
	delete(inventory, "config-file")
	delete(inventory, "mem-allocator")
}
