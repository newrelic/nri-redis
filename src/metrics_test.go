package main

import (
	"testing"

	"io/ioutil"
	"os"

	"reflect"

	"fmt"

	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/metric"
	"github.com/newrelic/infra-integrations-sdk/sdk"
)

var expectedRawInfoFromSample = map[string]interface{}{
	"active_defrag_hits":             0,
	"active_defrag_key_hits":         0,
	"active_defrag_key_misses":       0,
	"active_defrag_misses":           0,
	"active_defrag_running":          0,
	"aof_current_rewrite_time_sec":   -1,
	"aof_enabled":                    0,
	"aof_last_bgrewrite_status":      "ok",
	"aof_last_cow_size":              0,
	"aof_last_rewrite_time_sec":      -1,
	"aof_last_write_status":          "ok",
	"aof_rewrite_in_progress":        0,
	"aof_rewrite_scheduled":          0,
	"arch_bits":                      64,
	"atomicvar_api":                  "atomic-builtin",
	"blocked_clients":                0,
	"client_biggest_input_buf":       0,
	"client_longest_output_list":     0,
	"cluster_enabled":                0,
	"config_file":                    "",
	"connected_clients":              1,
	"connected_slaves":               0,
	"evicted_keys":                   0,
	"executable":                     "/Users/newrelic/redis-server",
	"expired_keys":                   0,
	"gcc_version":                    "4.2.1",
	"hz":                             10,
	"instantaneous_input_kbps":       0.00,
	"instantaneous_ops_per_sec":      0,
	"instantaneous_output_kbps":      0.00,
	"keyspace_hits":                  0,
	"keyspace_misses":                0,
	"latest_fork_usec":               0,
	"lazyfree_pending_objects":       0,
	"loading":                        0,
	"lru_clock":                      622980,
	"master_repl_offset":             0,
	"master_replid":                  "f01929bda7bae06c4aaf8eb319ed04ec64e97965",
	"master_replid2":                 "0000000000000000000000000000000000000000",
	"maxmemory":                      0,
	"maxmemory_human":                "0B",
	"maxmemory_policy":               "noeviction",
	"mem_allocator":                  "libc",
	"mem_fragmentation_ratio":        2.23,
	"migrate_cached_sockets":         0,
	"multiplexing_api":               "kqueue",
	"os":                             "Darwin 17.2.0 x86_64",
	"process_id":                     11432,
	"pubsub_channels":                0,
	"pubsub_patterns":                0,
	"rdb_bgsave_in_progress":         0,
	"rdb_changes_since_last_save":    0,
	"rdb_current_bgsave_time_sec":    -1,
	"rdb_last_bgsave_status":         "ok",
	"rdb_last_bgsave_time_sec":       -1,
	"rdb_last_save_time":             1510570985,
	"redis_build_id":                 "993aa70a2300c21e",
	"redis_git_dirty":                0,
	"redis_git_sha1":                 "00000000",
	"redis_mode":                     "standalone",
	"redis_version":                  "4.0.2",
	"rejected_connections":           0,
	"repl_backlog_active":            0,
	"repl_backlog_first_byte_offset": 0,
	"repl_backlog_histlen":           0,
	"repl_backlog_size":              1048576,
	"role":                           "master",
	"run_id":                         "18d4d1e817d8ce8388cfecc70dc4ec7fcd4767b1",
	"second_repl_offset":             -1,
	"slave_expires_tracked_keys":     0,
	"sync_full":                      0,
	"sync_partial_err":               0,
	"sync_partial_ok":                0,
	"tcp_port":                       6379,
	"total_commands_processed":       19,
	"total_connections_received":     10,
	"total_net_input_bytes":          394,
	"total_net_output_bytes":         67347,
	"total_system_memory":            17179869184,
	"total_system_memory_human":      "16.00G",
	"uptime_in_days":                 0,
	"uptime_in_seconds":              1435,
	"used_cpu_sys":                   0.58,
	"used_cpu_sys_children":          0.00,
	"used_cpu_user":                  0.30,
	"used_cpu_user_children":         0.00,
	"used_memory":                    1014816,
	"used_memory_dataset":            1986,
	"used_memory_dataset_perc":       "3.85%",
	"used_memory_human":              "991.03K",
	"used_memory_lua":                37888,
	"used_memory_lua_human":          "37.00K",
	"used_memory_overhead":           1012830,
	"used_memory_peak":               1032128,
	"used_memory_peak_human":         "1007.94K",
	"used_memory_peak_perc":          "98.32%",
	"used_memory_rss":                2260992,
	"used_memory_rss_human":          "2.16M",
	"used_memory_startup":            963200,
	"rdb_last_cow_size":              0,
	"server.name":                    "test_instance",
}

var expectedMetricSetFromSample = metric.MetricSet{
	"cluster.connectedSlaves":                0,
	"cluster.role":                           "master",
	"db.aofLastBgrewriteStatus":              "ok",
	"db.aofLastRewriteTimeMiliseconds":       -1,
	"db.aofLastWriteStatus":                  "ok",
	"db.evictedKeysPerSecond":                0.0,
	"db.expiredKeysPerSecond":                0.0,
	"db.keyspaceHitsPerSecond":               0.0,
	"db.keyspaceMissesPerSecond":             0.0,
	"db.latestForkMilliseconds":              0,
	"db.rdbBgsaveInProgress":                 0,
	"db.rdbChangesSinceLastSave":             0,
	"db.rdbLastBgsaveStatus":                 "ok",
	"db.rdbLastBgsaveTimeMilliseconds":       -1,
	"db.rdbLastSaveTime":                     1510570985,
	"db.syncFull":                            0,
	"db.syncPartialErr":                      0,
	"db.syncPartialOk":                       0,
	"event_type":                             "metricsTestSample",
	"net.blockedClients":                     0,
	"net.clientBiggestInputBufBytes":         0,
	"net.clientLongestOutputList":            0,
	"net.commandsProcessedPerSecond":         0.0,
	"net.connectedClients":                   1,
	"net.connectionsReceivedPerSecond":       0.0,
	"net.inputBytesPerSecond":                0.0,
	"net.outputBytesPerSecond":               0.0,
	"net.pubsubChannels":                     0,
	"net.pubsubPatterns":                     0,
	"net.rejectedConnectionsPerSecond":       0.0,
	"software.uptimeMilliseconds":            1435000,
	"software.version":                       "4.0.2",
	"system.totalSystemMemoryBytes":          17179869184,
	"system.usedCpuSysMilliseconds":          580,
	"system.usedCpuSysChildrenMilliseconds":  0,
	"system.usedCpuUserMilliseconds":         300,
	"system.usedCpuUserChildrenMilliseconds": 0,
	"system.usedMemoryBytes":                 1014816,
	"system.usedMemoryLuaBytes":              37888,
	"system.usedMemoryPeakBytes":             1032128,
	"system.usedMemoryRssBytes":              2260992,
	"system.memFragmentationRatio":           2.23,
	"server.name":                            "testInstance",
}

var expectedRawKeyspaceInfoFromSample = map[string]map[string]interface{}{
	"db0": {
		"keys":     1,
		"expires":  1,
		"avg_ttl":  7853,
		"keyspace": "db0",
	},
}

var expectedKeyspaceMetricSetFromSample = metric.MetricSet{
	"db.keys":               1,
	"db.expires":            1,
	"db.avgTtlMilliseconds": 7853,
	"db.keyspace":           "db0",
	"event_type":            "keyspaceTestSample_db0",
}

func readInfoSample() string {
	f, err := os.Open("data/info_sample.txt")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	sample, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	return string(sample)
}

func TestGetRawMetricsFromSample(t *testing.T) {
	metrics, keyspace, err := getRawMetrics(readInfoSample())
	if err != nil {
		t.Error("error getting raw metrics")
	}

	if len(expectedRawInfoFromSample) != len(metrics) {
		t.Error("got info with different size")
	}

	if !reflect.DeepEqual(expectedRawInfoFromSample, metrics) {
		t.Error("metrics are different than expected")
	}

	if len(expectedRawKeyspaceInfoFromSample) != len(keyspace) {
		t.Error("got keyspace with different size")
	}

	if !reflect.DeepEqual(expectedRawKeyspaceInfoFromSample, keyspace) {
		t.Error("keyspace is different than expected")
	}
}

func TestAsValue(t *testing.T) {
	if asValue("1") != 1 ||
		asValue("1.49") != 1.49 ||
		asValue("00000000") != "00000000" ||
		asValue("0.0") != 0.0 {
		t.Error("wrong conversion")
	}
}

func TestPopulateMetrics(t *testing.T) {
	rawMetrics, rawKeyspace, err := getRawMetrics(readInfoSample())
	if err != nil {
		t.Error(err)
	}

	integration := sdk.Integration{}
	ms := integration.NewMetricSet("metricsTestSample")

	err = populateMetrics(ms, "test instance", rawMetrics, metricsDefinition)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expectedMetricSetFromSample, *ms) {
		t.Error("unexpected metric set")
	}

	for db, keyspaceMetrics := range rawKeyspace {
		ms = integration.NewMetricSet(fmt.Sprintf("keyspaceTestSample_%s", db))
		err = populateMetrics(ms, "testInstance", keyspaceMetrics, keyspaceMetricsDefinition)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(expectedKeyspaceMetricSetFromSample, *ms) {
			t.Errorf("unexpected keyspace metric set for %s", db)
		}
	}
}

func TestGetRawMetricsEmptyInput(t *testing.T) {
	info := ""
	metrics, keyspaceMetrics, err := getRawMetrics(info)
	expectedMetricsLength := 0
	expectedKeyspaceMetricsLength := 0

	if !reflect.DeepEqual(err.Error(), "Empty result") {
		t.Error()
	}
	if len(metrics) != expectedMetricsLength {
		t.Error()
	}
	if len(keyspaceMetrics) != expectedKeyspaceMetricsLength {
		t.Error()
	}
}

func TestGetRawMetricsNotValidInput(t *testing.T) {
	info := "."
	expectedMetricsLength := 0
	expectedKeyspaceMetricsLength := 0
	expectedMetrics := make(map[string]interface{})
	expectedKeyspaceMetrics := make(map[string]map[string]interface{})
	metrics, keyspaceMetrics, err := getRawMetrics(info)

	if err != nil {
		t.Error()
	}
	if len(metrics) != expectedMetricsLength {
		t.Error()
	}
	if len(keyspaceMetrics) != expectedKeyspaceMetricsLength {
		t.Error()
	}
	if !reflect.DeepEqual(metrics, expectedMetrics) {
		t.Error()
	}
	if !reflect.DeepEqual(keyspaceMetrics, expectedKeyspaceMetrics) {
		t.Error()
	}
}

func TestParseKeyspaceMetrics(t *testing.T) {
	dbName := "db0"
	keyspace := "keys=3,expires=1,avg_ttl=354949"
	metric, err := parseKeyspaceMetrics(dbName, keyspace)
	expectedLength := 4
	expectedMetric := map[string]interface{}{
		"keyspace": "db0",
		"keys":     3,
		"expires":  1,
		"avg_ttl":  354949,
	}

	if err != nil {
		t.Errorf("Error %v returned ", err)
	}
	if len(metric) != expectedLength {
		t.Errorf("Not all values processed, got %d length of the rawInventory, expected: %d", len(metric), expectedLength)
	}
	if !reflect.DeepEqual(metric, expectedMetric) {
		t.Error()
	}
}

func TestParseKeyspaceMetricsNoMatch(t *testing.T) {
	dbName := "db0"
	keyspace := "db0:keys=3,expires=1,"
	metrics, err := parseKeyspaceMetrics(dbName, keyspace)
	expectedLength := 0
	expectedMetrics := make(map[string]interface{})

	if err != nil {
		t.Errorf("Error %v returned ", err)
	}
	if len(metrics) != expectedLength {
		t.Errorf("Not all values processed, got %d length of the rawInventory, expected: %d", len(metrics), expectedLength)
	}
	if !reflect.DeepEqual(metrics, expectedMetrics) {
		t.Error()
	}
}

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
	if err != nil {
		t.Error()
	}
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

func TestPopulateCustomKeysMetric(t *testing.T) {
	rawCustomKeys := map[string]keyInfo{
		"myhash": {keyLength: 1, keyType: "hash"},
		"mylist": {keyLength: 5, keyType: "list"},
	}
	expectedHashKeyName := "db.keyLength.hash.myhash"
	expectedListKeyName := "db.keyLength.list.mylist"
	expectedHashKeyLength := int64(1)
	expectedListKeyLength := int64(5)

	var sample = metric.NewMetricSet("RedisKeyspaceSample")
	populateCustomKeysMetric(&sample, rawCustomKeys)

	if sample[expectedHashKeyName] != expectedHashKeyLength {
		t.Error()
	}
	if sample[expectedListKeyName] != expectedListKeyLength {
		t.Error()
	}
	if len(sample) != 3 {
		t.Error()
	}

}
