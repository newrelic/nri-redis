package main

import (
	"flag"
	"fmt"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/newrelic/nri-redis/tests/integration/helpers"
	"github.com/newrelic/nri-redis/tests/integration/jsonschema"
)

var (
	//default values
	defaultContainer = "integration_nri-redis_1"
	defaultBinPath   = "/nr-redis"
	defaultHost      = "redis"
	defaultPort      = 6379

	// cli flags
	container = flag.String("container", defaultContainer, "container where the integration is installed")
	binPath   = flag.String("bin", defaultBinPath, "Integration binary path")
	host      = flag.String("host", defaultHost, "Redis host ip address")
	port      = flag.Int("port", defaultPort, "Redis port")
)

func runIntegration(t *testing.T, envVars ...string) (stdout string, stderr string) {
	t.Helper()

	command := make([]string, 0)
	command = append(command, *binPath)
	if host != nil {
		command = append(command, "--hostname", *host)
	}
	if port != nil {
		command = append(command, "--port", strconv.Itoa(*port))
	}
	stdout, stderr, err := helpers.ExecInContainer(*container, command, envVars...)

	if stderr != "" {
		log.Debug("Integration command Standard Error: ", stderr)
	}
	require.NoError(t, err)

	return stdout, stderr
}

func TestMain(m *testing.M) {
	flag.Parse()

	result := m.Run()

	os.Exit(result)
}

func TestRedisIntegration(t *testing.T) {
	testName := helpers.GetTestName(t)
	stdout, stderr := runIntegration(t, fmt.Sprintf("NRIA_CACHE_PATH=/tmp/%v.json", testName))

	schemaPath := filepath.Join("json-schema-files", "redis-schema.json")

	err := jsonschema.Validate(schemaPath, stdout)

	assert.NoError(t, err, "The output of Redis integration doesn't have expected format")
	assert.NotNil(t, stderr, "unexpected stderr")
}

func TestRedisIntegration_OnlyMetrics(t *testing.T) {
	testName := helpers.GetTestName(t)

	stdout, stderr := runIntegration(t, fmt.Sprintf("NRIA_CACHE_PATH=/tmp/%v.json", testName))

	schemaPath := filepath.Join("json-schema-files", "redis-schema-metrics.json")

	err := jsonschema.Validate(schemaPath, stdout)

	assert.NoError(t, err, "The output of Redis integration doesn't have expected format")
	assert.NotNil(t, stderr, "unexpected stderr")
}

func TestRedisIntegration_OnlyInventory(t *testing.T) {
	testName := helpers.GetTestName(t)

	stdout, stderr := runIntegration(t, fmt.Sprintf("NRIA_CACHE_PATH=/tmp/%v.json", testName))

	schemaPath := filepath.Join("json-schema-files", "redis-schema-inventory.json")

	err := jsonschema.Validate(schemaPath, stdout)

	assert.NoError(t, err, "The output of Redis integration doesn't have expected format")
	assert.NotNil(t, stderr, "unexpected stderr")
}
