package main

import (
	"fmt"
	"os"
	"strconv"

	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/infra-integrations-sdk/persist"
)

type argumentList struct {
	sdkArgs.DefaultArgumentList
	Hostname         string       `default:"localhost" help:"Hostname or IP where Redis server is running."`
	Port             int          `default:"6379" help:"Port on which Redis server is listening."`
	UnixSocketPath   string       `default:"" help:"Unix socket path on which Redis server is listening."`
	Keys             sdkArgs.JSON `default:"" help:"List of the keys for retrieving their lengths"`
	KeysLimit        int          `default:"30" help:"Max number of the keys to retrieve their lengths"`
	Password         string       `help:"Password to use when connecting to the Redis server."`
	RemoteMonitoring bool         `default:"false" help:"Allows to monitor multiple instances as 'remote' entity. Set to 'FALSE' value for backwards compatibility otherwise set to 'TRUE'"`
}

const (
	integrationName    = "com.newrelic.redis"
	integrationVersion = "1.2.0"
	entityRemoteType   = "instance"
)

var args argumentList

func main() {
	i, err := createIntegration()
	fatalIfErr(err)

	conn, err := newRedisCon(args.Hostname, args.Port, args.UnixSocketPath, args.Password)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	info, config, err := conn.GetData()
	fatalIfErr(err)

	rawMetrics, rawKeyspaceMetrics, metricsErr := getRawMetrics(info)

	e, err := entity(i, &args)
	fatalIfErr(err)

	if args.HasInventory() {
		rawInventory := getRawInventory(config, rawMetrics)
		populateInventory(e.Inventory, rawInventory)
	}

	if args.HasMetrics() {
		fatalIfErr(metricsErr)

		ms := metricSet(e, "RedisSample", args.Hostname, args.Port)
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

		ms = metricSet(e, "RedisKeyspaceSample", args.Hostname, args.Port)
		for db, keyspaceMetrics := range rawKeyspaceMetrics {
			fatalIfErr(populateMetrics(ms, keyspaceMetrics, keyspaceMetricsDefinition))
			if _, ok := rawCustomKeysMetric[db]; ok && keysFlagPresent {
				populateCustomKeysMetric(ms, rawCustomKeysMetric[db])
			}
		}
	}

	fatalIfErr(i.Publish())
}

func metricSet(e *integration.Entity, eventType, hostname string, port int) *metric.Set {
	strPort := strconv.Itoa(port)
	if args.RemoteMonitoring {
		return e.NewMetricSet(
			eventType,
			metric.Attr("hostname", hostname),
			metric.Attr("port", strPort),
		)
	}

	return e.NewMetricSet(
		eventType,
		metric.Attr("port", strPort),
	)
}

func createIntegration() (*integration.Integration, error) {
	cachePath := os.Getenv("NRIA_CACHE_PATH")
	if cachePath == "" {
		return integration.New(integrationName, integrationVersion, integration.Args(&args))
	}

	l := log.NewStdErr(args.Verbose)
	s, err := persist.NewFileStore(cachePath, l, persist.DefaultTTL)
	if err != nil {
		return nil, err
	}

	return integration.New(integrationName, integrationVersion, integration.Args(&args), integration.Storer(s), integration.Logger(l))
}

func entity(i *integration.Integration, args *argumentList) (*integration.Entity, error) {
	if args.RemoteMonitoring {
		n := fmt.Sprintf("%s:%d", args.Hostname, args.Port)
		return i.Entity(n, entityRemoteType)
	}

	return i.LocalEntity(), nil
}

func fatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
