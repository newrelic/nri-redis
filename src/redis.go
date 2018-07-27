package main

import (
	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/infra-integrations-sdk/sdk"
)

type argumentList struct {
	sdkArgs.DefaultArgumentList
	Hostname       string       `default:"localhost" help:"Hostname or IP where Redis server is running."`
	Port           int          `default:"6379" help:"Port on which Redis server is listening."`
	UnixSocketPath string       `default:"" help:"Unix socket path on which Redis server is listening."`
	Keys           sdkArgs.JSON `default:"" help:"List of the keys for retrieving their lengths"`
	KeysLimit      int          `default:"30" help:"Max number of the keys to retrieve their lengths"`
	Password       string       `help:"Password to use when connecting to the Redis server."`
}

const (
	integrationName    = "com.newrelic.redis"
	integrationVersion = "1.0.0"
)

var args argumentList

func main() {
	integration, err := sdk.NewIntegration(integrationName, integrationVersion, &args)
	fatalIfErr(err)

	conn, err := newRedisCon(args.Hostname, args.Port, args.UnixSocketPath, args.Password)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	info, config, err := conn.GetData()
	fatalIfErr(err)

	rawMetrics, rawKeyspaceMetrics, metricsErr := getRawMetrics(info)

	if args.All || args.Inventory {
		rawInventory := getRawInventory(config, rawMetrics)
		populateInventory(integration.Inventory, rawInventory)
	}

	if args.All || args.Metrics {
		fatalIfErr(metricsErr)
		ms := integration.NewMetricSet("RedisSample")
		fatalIfErr(populateMetrics(ms, rawMetrics, metricsDefinition))

		var rawCustomKeysMetric map[string]map[string]keyInfo
		keysFlagPresent := args.Keys.Get() != nil

		if keysFlagPresent {
			databaseKeys := getDbAndKeys(args.Keys)
			_, keysFlagErr := validateKeysFlag(databaseKeys, args.KeysLimit)

			if keysFlagErr != nil {
				log.Warn("Error processing keys flag: %v", keysFlagErr)
			} else {
				rawCustomKeysMetric, err = conn.GetRawCustomKeys(databaseKeys)
				if err != nil {
					log.Warn("Got error: %v", err)
				}
			}
		}

		for db, keyspaceMetrics := range rawKeyspaceMetrics {
			ms = integration.NewMetricSet("RedisKeyspaceSample")
			fatalIfErr(populateMetrics(ms, keyspaceMetrics, keyspaceMetricsDefinition))

			if _, ok := rawCustomKeysMetric[db]; ok && keysFlagPresent {
				populateCustomKeysMetric(ms, rawCustomKeysMetric[db])
			}
		}
	}

	fatalIfErr(integration.Publish())
}

func fatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
