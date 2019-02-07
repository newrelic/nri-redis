package main

import (
	"flag"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/newrelic/nri-redis/tests/integration/helpers"
	"github.com/newrelic/nri-redis/tests/integration/jsonschema"
)

var (
	iName = "redis"

	//default values
	defaultContainer = "integration_nri-redis_1"
	defaultBinPath   = "/nr-redis"
	defaultHost = "redis"
	defaultPort = 6379

	// cli flags
	iVersion = "1.2.0"
	container = flag.String("container", defaultContainer, "container where the integration is installed")
	update    = flag.Bool("test.update", false, "update json-schema file")
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

// changeRedisPassword sends request to Redis server through CONFIG SET command
// to change Redis server password. It returns standard output of the CONFIG
// SET command and error (if any)
func changeRedisPassword(currentPassword string, newPassword string) ([]byte, error) {
	return exec.Command("redis-cli", "-a", currentPassword, "config", "set", "requirepass", newPassword).Output()
}

func setup() error {
	flag.Parse()

	return nil
}

func teardown() error {
	return nil
}

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		tErr := teardown()
		if tErr != nil {
			fmt.Printf("Error during the teardown of the tests: %s\n", tErr)
		}
		os.Exit(1)
	}
	result := m.Run()
	err = teardown()
	if err != nil {
		fmt.Printf("Error during the teardown of the tests: %s\n", err)
	}
	os.Exit(result)
}

func TestRedisIntegration(t *testing.T) {
	testName := helpers.GetTestName(t)
	stdout, stderr := runIntegration(t, fmt.Sprintf("NRIA_CACHE_PATH=/tmp/%v.json", testName))

	schemaPath := filepath.Join("json-schema-files", "redis-schema.json")
	if *update {
		schema, err := jsonschema.Generate(stdout)
		if err != nil {
			t.Fatal(err)
		}

		schemaJSON, err := simplejson.NewJson(schema)
		if err != nil {
			t.Fatalf("Cannot unmarshal JSON schema, got error: %v", err)
		}
		err = helpers.ModifyJSONSchemaGlobal(schemaJSON, iName, 1, iVersion)
		if err != nil {
			t.Fatal(err)
		}
		err = helpers.ModifyJSONSchemaInventoryPresent(schemaJSON)
		if err != nil {
			t.Fatal(err)
		}
		err = jsonschema.AddNewElements(
			schemaJSON.GetPath("properties", "inventory", "properties", "requirepass", "properties"),
			map[string]jsonschema.ValidationField{
				"value": {"pattern", "^\\(omitted value\\)$"},
			})
		if err != nil {
			t.Fatal(err)
		}

		err = helpers.ModifyJSONSchemaMetricsPresent(schemaJSON, "RedisSample")
		if err != nil {
			t.Fatal(err)
		}
		schema, err = schemaJSON.MarshalJSON()
		if err != nil {
			t.Fatalf("Cannot marshal JSON schema, got error: %v", err)
		}
		err = ioutil.WriteFile(schemaPath, schema, 0644)
		if err != nil {
			t.Fatal(err)
		}
	}
	err := jsonschema.Validate(schemaPath, stdout)
	if err != nil {
		t.Fatalf("The output of Redis integration doesn't have expected format. Err: %s", err)
	}

	if stderr != "" {
		t.Fatalf("unexpected stderr %s", stderr)
	}
}

func TestRedisIntegration_OnlyMetrics(t *testing.T) {
	testName := helpers.GetTestName(t)

	stdout, stderr := runIntegration(t, fmt.Sprintf("NRIA_CACHE_PATH=/tmp/%v.json", testName))

	schemaPath := filepath.Join("json-schema-files", "redis-schema-metrics.json")
	if *update {
		schema, err := jsonschema.Generate(stdout)
		if err != nil {
			t.Fatal(err)
		}

		schemaJSON, err := simplejson.NewJson(schema)
		if err != nil {
			t.Fatalf("Cannot unmarshal JSON schema, got error: %v", err)
		}
		err = helpers.ModifyJSONSchemaGlobal(schemaJSON, iName, 1, iVersion)
		if err != nil {
			t.Fatal(err)
		}
		err = helpers.ModifyJSONSchemaMetricsPresent(schemaJSON, "RedisSample")
		if err != nil {
			t.Fatal(err)
		}
		err = helpers.ModifyJSONSchemaNoInventory(schemaJSON)
		if err != nil {
			t.Fatal(err)
		}
		schema, err = schemaJSON.MarshalJSON()
		if err != nil {
			t.Fatalf("Cannot marshal JSON schema, got error: %v", err)
		}
		err = ioutil.WriteFile(schemaPath, schema, 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	err := jsonschema.Validate(schemaPath, stdout)
	if err != nil {
		t.Fatalf("The output of Redis integration doesn't have expected format. Err: %s", err)
	}

	if stderr != "" {
		t.Fatalf("unexpected stderr %s", stderr)
	}
}

func TestRedisIntegration_OnlyInventory(t *testing.T) {
	testName := helpers.GetTestName(t)

	stdout, stderr := runIntegration(t, fmt.Sprintf("NRIA_CACHE_PATH=/tmp/%v.json", testName))

	schemaPath := filepath.Join("json-schema-files", "redis-schema-inventory.json")
	if *update {
		schema, err := jsonschema.Generate(stdout)
		if err != nil {
			t.Fatal(err)
		}
		schemaJSON, err := simplejson.NewJson(schema)
		if err != nil {
			t.Fatalf("Cannot unmarshal JSON schema, got error: %v", err)
		}

		err = helpers.ModifyJSONSchemaGlobal(schemaJSON, iName, 1, iVersion)
		if err != nil {
			t.Fatal(err)
		}
		err = helpers.ModifyJSONSchemaNoMetrics(schemaJSON)
		if err != nil {
			t.Fatal(err)
		}
		err = helpers.ModifyJSONSchemaInventoryPresent(schemaJSON)
		if err != nil {
			t.Fatal(err)
		}
		err = jsonschema.AddNewElements(
			schemaJSON.GetPath("properties", "inventory", "properties", "requirepass", "properties"),
			map[string]jsonschema.ValidationField{
				"value": {"pattern", "^\\(omitted value\\)$"},
			})
		if err != nil {
			t.Fatal(err)
		}
		schema, err = schemaJSON.MarshalJSON()
		if err != nil {
			t.Fatalf("Cannot marshal JSON schema, got error: %v", err)
		}
		if err = ioutil.WriteFile(schemaPath, schema, 0644); err != nil {
			t.Fatal(err)
		}
	}

	err := jsonschema.Validate(schemaPath, stdout)
	if err != nil {
		t.Fatalf("The output of Redis integration doesn't have expected format. Err: %s", err)
	}

	if stderr != "" {
		t.Fatalf("unexpected stderr %s", stderr)
	}
}

