//go:generate goversioninfo

package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"

	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/data/attribute"
	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/infra-integrations-sdk/persist"
)

type argumentList struct {
	sdkArgs.DefaultArgumentList
	Hostname              string       `default:"localhost" help:"Hostname or IP where Redis server is running."`
	Port                  int          `default:"6379" help:"Port on which Redis server is listening."`
	UnixSocketPath        string       `default:"" help:"Unix socket path on which Redis server is listening."`
	Keys                  sdkArgs.JSON `default:"" help:"List of the keys for retrieving their lengths"`
	KeysLimit             int          `default:"30" help:"Max number of the keys to retrieve their lengths"`
	Username              string       `help:"Username to use when connecting to the Redis server."`
	Password              string       `help:"Password to use when connecting to the Redis server."`
	UseUnixSocket         bool         `default:"false" help:"Adds the UnixSocketPath value to the entity. If you are monitoring more than one Redis instance on the same host using Unix sockets, then you should set it to true."`
	RemoteMonitoring      bool         `default:"false" help:"Allows to monitor multiple instances as 'remote' entity. Set to 'FALSE' value for backwards compatibility otherwise set to 'TRUE'"`
	RenamedCommands       sdkArgs.JSON `default:"" help:"Map of default redis commands to their renamed form, if rename-command config has been used in the redis server."`
	ConfigInventory       bool         `default:"true" help:"Provides CONFIG inventory information. Set it to 'false' in environments where the Redis CONFIG command is prohibited (e.g. AWS ElastiCache)"`
	ShowVersion           bool         `default:"false" help:"Print build information and exit"`
	UseTLS                bool         `default:"false" help:"Use TLS when communicating with the Redis server."`
	TLSInsecureSkipVerify bool         `default:"false" help:"Disable server name verification when connecting over TLS"`
}

const (
	integrationName  = "com.newrelic.redis"
	entityRemoteType = "instance"
)

var (
	args               argumentList
	integrationVersion = "0.0.0"
	gitCommit          = ""
	buildDate          = ""

	errorArgs = errors.New("no connection method available, UnixSocketPath or Hostname and Port should be populated")
)

func main() {
	i, err := createIntegration()
	fatalIfErr(err)

	if args.ShowVersion {
		printVersion()
		os.Exit(0)
	}

	dialOptions := standardDialOptions(args.Username, args.Password)

	var c *redisConn
	switch {
	// Notice that we are not checking UseUnixSocket since it is not used to define how to connect, but merely the entity name.
	// There are users having use_unix_socket=true and then connecting with hostname and port,
	// or use_unix_socket=false and then connecting with the unix socket.
	case args.UnixSocketPath != "":
		c, err = newSocketRedisCon(args.UnixSocketPath, dialOptions...)
		fatalIfErr(err)
	case args.Hostname != "" && args.Port > 0:
		if args.UseTLS {
			tlsDialOptions := tlsDialOptions(args.UseTLS, args.TLSInsecureSkipVerify)
			dialOptions = append(dialOptions, tlsDialOptions...)
		}
		redisURL := net.JoinHostPort(args.Hostname, strconv.Itoa(args.Port))
		c, err = newNetworkRedisCon(redisURL, dialOptions...)
		fatalIfErr(err)
	default:
		log.Fatal(errorArgs)
	}

	defer c.Close()

	// Support using renamed form of redis commands, if 'renamed-command' config is used in Redis server
	if args.RenamedCommands.Get() != nil {
		renamedCommands, err := getRenamedCommands(args.RenamedCommands)
		fatalIfErr(err)

		c.RenameCommands(renamedCommands)
	}

	info, err := c.GetInfo()
	fatalIfErr(err)

	rawMetrics, rawKeyspaceMetrics, metricsErr := getRawMetrics(info)

	e, err := entity(i, &args)
	fatalIfErr(err)

	if args.HasInventory() {
		var config map[string]string
		if args.ConfigInventory {
			config, err = c.GetConfig()
			if err != nil {
				fmtStr := "%v. Configuration inventory won't be reported"
				if _, ok := err.(configConnectionError); ok {
					fmtStr += ". This may be expected if you are monitoring a managed " +
						"Redis instance with limited permissions. " +
						"Set the 'config_inventory' argument to 'false' to remove this message"
				}
				log.Warn(fmtStr, err)
			}
		}
		rawInventory := getRawInventory(config, rawMetrics)
		populateInventory(e.Inventory, rawInventory)
	}

	if args.HasMetrics() {
		fatalIfErr(metricsErr)

		ms := metricSet(e, "RedisSample", args.Hostname, args.Port, args.RemoteMonitoring)
		fatalIfErr(populateMetrics(ms, rawMetrics, metricsDefinition))

		var rawCustomKeysMetric map[string]map[string]keyInfo
		keysFlagPresent := args.Keys.Get() != nil

		if keysFlagPresent {
			databaseKeys := getDBAndKeys(args.Keys)
			_, keysFlagErr := validateKeysFlag(databaseKeys, args.KeysLimit)

			if keysFlagErr != nil {
				log.Warn("Error processing keys flag: %v", keysFlagErr)
			} else {
				rawCustomKeysMetric, err = c.GetRawCustomKeys(databaseKeys)
				if err != nil {
					log.Warn("Got error: %v", err)
				}
			}
		}

		ms = metricSet(e, "RedisKeyspaceSample", args.Hostname, args.Port, args.RemoteMonitoring)
		for db, keyspaceMetrics := range rawKeyspaceMetrics {
			fatalIfErr(populateMetrics(ms, keyspaceMetrics, keyspaceMetricsDefinition))
			if _, ok := rawCustomKeysMetric[db]; ok && keysFlagPresent {
				populateCustomKeysMetric(ms, rawCustomKeysMetric[db])
			}
		}
	}

	fatalIfErr(i.Publish())
}

func printVersion() {
	fmt.Printf(
		"New Relic %s integration Version: %s, Platform: %s, GoVersion: %s, GitCommit: %s, BuildDate: %s\n",
		strings.Title(strings.Replace(integrationName, "com.newrelic.", "", 1)),
		integrationVersion,
		fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		runtime.Version(),
		gitCommit,
		buildDate)
}

func metricSet(e *integration.Entity, eventType, hostname string, port int, remote bool) *metric.Set {
	strPort := strconv.Itoa(port)
	if remote {
		return e.NewMetricSet(
			eventType,
			attribute.Attr("hostname", hostname),
			attribute.Attr("port", strPort),
		)
	}

	return e.NewMetricSet(
		eventType,
		attribute.Attr("port", strPort),
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
		var n string
		if args.UseUnixSocket && args.UnixSocketPath != "" {
			n = fmt.Sprintf("%s:%s", args.Hostname, args.UnixSocketPath)
		} else {
			n = net.JoinHostPort(args.Hostname, strconv.Itoa(args.Port))
		}
		return i.Entity(n, entityRemoteType)
	}

	return i.LocalEntity(), nil
}

func fatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
