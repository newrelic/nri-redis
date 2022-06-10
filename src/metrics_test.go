package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/newrelic/infra-integrations-sdk/data/attribute"
	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/persist"
	"github.com/stretchr/testify/assert"
)

var expectedRawInfoFromSample = map[string]interface{}{
	"active_defrag_hits":                0,
	"active_defrag_key_hits":            0,
	"active_defrag_key_misses":          0,
	"active_defrag_misses":              0,
	"active_defrag_running":             0,
	"aof_current_rewrite_time_sec":      -1,
	"aof_enabled":                       0,
	"aof_last_bgrewrite_status":         "ok",
	"aof_last_cow_size":                 0,
	"aof_last_rewrite_time_sec":         -1,
	"aof_last_write_status":             "ok",
	"aof_rewrite_in_progress":           0,
	"aof_rewrite_scheduled":             0,
	"arch_bits":                         64,
	"atomicvar_api":                     "atomic-builtin",
	"blocked_clients":                   0,
	"client_biggest_input_buf":          0,
	"client_longest_output_list":        0,
	"cluster_enabled":                   0,
	"config_file":                       "",
	"connected_clients":                 1,
	"connected_slaves":                  0,
	"evicted_keys":                      0,
	"executable":                        "/Users/newrelic/redis-server",
	"expired_keys":                      0,
	"gcc_version":                       "4.2.1",
	"hz":                                10,
	"instantaneous_input_kbps":          0.00,
	"instantaneous_ops_per_sec":         0,
	"instantaneous_output_kbps":         0.00,
	"keyspace_hits":                     0,
	"keyspace_misses":                   0,
	"latest_fork_usec":                  0,
	"lazyfree_pending_objects":          0,
	"loading":                           0,
	"lru_clock":                         622980,
	"master_repl_offset":                0,
	"master_replid":                     "f01929bda7bae06c4aaf8eb319ed04ec64e97965",
	"master_replid2":                    "0000000000000000000000000000000000000000",
	"maxmemory":                         123,
	"maxmemory_human":                   "123B",
	"maxmemory_policy":                  "noeviction",
	"mem_allocator":                     "libc",
	"mem_fragmentation_ratio":           2.23,
	"migrate_cached_sockets":            0,
	"multiplexing_api":                  "kqueue",
	"os":                                "Darwin 17.2.0 x86_64",
	"process_id":                        11432,
	"pubsub_channels":                   0,
	"pubsub_patterns":                   0,
	"rdb_bgsave_in_progress":            0,
	"rdb_changes_since_last_save":       0,
	"rdb_current_bgsave_time_sec":       -1,
	"rdb_last_bgsave_status":            "ok",
	"rdb_last_bgsave_time_sec":          -1,
	"rdb_last_save_time":                1510570985,
	"redis_build_id":                    "993aa70a2300c21e",
	"redis_git_dirty":                   0,
	"redis_git_sha1":                    "00000000",
	"redis_mode":                        "standalone",
	"redis_version":                     "4.0.2",
	"rejected_connections":              0,
	"repl_backlog_active":               0,
	"repl_backlog_first_byte_offset":    0,
	"repl_backlog_histlen":              0,
	"repl_backlog_size":                 1048576,
	"role":                              "master",
	"run_id":                            "18d4d1e817d8ce8388cfecc70dc4ec7fcd4767b1",
	"second_repl_offset":                -1,
	"slave_expires_tracked_keys":        0,
	"sync_full":                         0,
	"sync_partial_err":                  0,
	"sync_partial_ok":                   0,
	"tcp_port":                          6379,
	"total_commands_processed":          19,
	"total_connections_received":        10,
	"total_net_input_bytes":             394,
	"total_net_output_bytes":            67347,
	"total_system_memory":               17179869184,
	"total_system_memory_human":         "16.00G",
	"uptime_in_days":                    0,
	"uptime_in_seconds":                 1435,
	"used_cpu_sys":                      0.58,
	"used_cpu_sys_children":             0.00,
	"used_cpu_user":                     0.30,
	"used_cpu_user_children":            0.00,
	"used_memory":                       1014816,
	"used_memory_dataset":               1986,
	"used_memory_dataset_perc":          "3.85%",
	"used_memory_human":                 "991.03K",
	"used_memory_lua":                   37888,
	"used_memory_lua_human":             "37.00K",
	"used_memory_overhead":              1012830,
	"used_memory_peak":                  1032128,
	"used_memory_peak_human":            "1007.94K",
	"used_memory_peak_perc":             "98.32%",
	"used_memory_rss":                   2260992,
	"used_memory_rss_human":             "2.16M",
	"used_memory_startup":               963200,
	"rdb_last_cow_size":                 0,
	"monotonic_clock":                   "POSIX clock_gettime",
	"process_supervised":                "no",
	"server_time_usec":                  1654871464152599,
	"configured_hz":                     0,
	"io_threads_active":                 0,
	"cluster_connections":               0,
	"maxclients":                        0,
	"client_recent_max_input_buffer":    0,
	"client_recent_max_output_buffer":   0,
	"tracking_clients":                  0,
	"clients_in_timeout_table":          0,
	"allocator_allocated":               0,
	"allocator_active":                  0,
	"allocator_resident":                0,
	"used_memory_vm_eval":               0,
	"used_memory_scripts_eval":          0,
	"number_of_cached_scripts":          0,
	"number_of_functions":               0,
	"number_of_libraries":               0,
	"used_memory_vm_functions":          0,
	"used_memory_vm_total":              0,
	"used_memory_vm_total_human":        "1.0K",
	"used_memory_functions":             0,
	"used_memory_scripts":               0,
	"used_memory_scripts_human":         "216B",
	"allocator_frag_ratio":              0,
	"allocator_frag_bytes":              0,
	"allocator_rss_ratio":               0,
	"allocator_rss_bytes":               0,
	"rss_overhead_ratio":                1.02,
	"rss_overhead_bytes":                0,
	"mem_fragmentation_bytes":           0,
	"mem_not_counted_for_evict":         0,
	"mem_replication_backlog":           0,
	"mem_total_replication_buffers":     0,
	"mem_clients_slaves":                0,
	"mem_clients_normal":                0,
	"mem_cluster_links":                 0,
	"mem_aof_buffer":                    0,
	"lazyfreed_objects":                 0,
	"async_loading":                     0,
	"current_cow_peak":                  0,
	"current_cow_size":                  0,
	"current_cow_size_age":              0,
	"current_fork_perc":                 0,
	"current_save_keys_processed":       0,
	"current_save_keys_total":           0,
	"rdb_saves":                         0,
	"rdb_last_load_keys_expired":        0,
	"rdb_last_load_keys_loaded":         0,
	"aof_rewrites":                      0,
	"aof_rewrites_consecutive_failures": 0,
	"module_fork_in_progress":           0,
	"module_fork_last_cow_size":         0,
	"expired_stale_perc":                0,
	"expired_time_cap_reached_count":    0,
	"expire_cycle_cpu_milliseconds":     0,
	"evicted_clients":                   0,
	"total_eviction_exceeded_time":      0,
	"current_eviction_exceeded_time":    0,
	"total_forks":                       0,
	"total_active_defrag_time":          0,
	"current_active_defrag_time":        0,
	"tracking_total_keys":               0,
	"tracking_total_items":              0,
	"tracking_total_prefixes":           0,
	"unexpected_error_replies":          0,
	"total_error_replies":               0,
	"dump_payload_sanitizations":        0,
	"total_reads_processed":             0,
	"total_writes_processed":            0,
	"io_threaded_reads_processed":       0,
	"io_threaded_writes_processed":      0,
	"reply_buffer_shrinks":              0,
	"reply_buffer_expands":              0,
	"master_failover_state":             "no-failover",
}

var expectedMetricSetFromSample = map[string]interface{}{
	"event_type":                      "metricsTestSample",
	"software.uptimeMilliseconds":     1.435e+06,
	"software.version":                "4.0.2",               //
	"software.monotonicClock":         "POSIX clock_gettime", //
	"software.processSupervised":      "no",
	"software.serverTimeMilliseconds": 1.654871464152e+12,
	"software.configuredHz":           0.0,
	"software.ioThreadsActive":        0.0,
	// Clients
	"net.connectedClients":            1.0,
	"net.clientLongestOutputList":     0.0,
	"net.clientBiggestInputBufBytes":  0.0,
	"net.blockedClients":              0.0,
	"net.clusterConnections":          0.0,
	"net.maxClients":                  0.0,
	"net.clientRecentMaxInputBuffer":  0.0,
	"net.clientRecentMaxOutputBuffer": 0.0,
	"net.trackingClients":             0.0,
	"net.clientsInTimeoutTable":       0.0,
	// Memory
	"system.usedMemoryBytes":            1.014816e+06,
	"system.usedMemoryRssBytes":         2.260992e+06,
	"system.usedMemoryPeakBytes":        1.032128e+06,
	"system.usedMemoryLuaBytes":         37888.0,
	"system.usedMemoryHuman":            "991.03K", //
	"system.usedMemoryOverhead":         1.01283e+06,
	"system.totalSystemMemoryBytes":     1.7179869184e+10,
	"system.maxmemoryBytes":             123.0,
	"system.memFragmentationRatio":      2.23,
	"system.allocatorAllocated":         0.0,
	"system.allocatorActive":            0.0,
	"system.allocatorResident":          0.0,
	"system.usedMemoryVmEval":           0.0,
	"system.memAllocator":               "libc", //
	"system.usedMemoryScriptsEval":      0.0,
	"system.numberOfCachedScripts":      0.0,
	"system.numberOfFunctions":          0.0,
	"system.numberOfLibraries":          0.0,
	"system.usedMemoryVmFunctions":      0.0,
	"system.usedMemoryVmTotal":          0.0,
	"system.usedMemoryVmTotalHuman":     "1.0K", //
	"system.usedMemoryFunctions":        0.0,
	"system.usedMemoryScripts":          0.0,
	"system.usedMemoryScriptsHuman":     "216B",     //
	"system.usedMemoryPeakHuman":        "1007.94K", //
	"system.allocatorFragRatio":         0.0,
	"system.allocatorFragBytes":         0.0,
	"system.allocatorRssRatio":          0.0,
	"system.allocatorRssBytes":          0.0,
	"system.rssOverheadRatio":           1.02,
	"system.memFragmentationBytes":      0.0,
	"system.memNotCountedForEvict":      0.0,
	"system.memReplicationBacklog":      0.0,
	"system.memTotalReplicationBuffers": 0.0,
	"system.memClientsSlaves":           0.0,
	"system.memClientsNormal":           0.0,
	"system.memClusterLinks":            0.0,
	"system.memAofBuffer":               0.0,
	"system.lazyfreedObjects":           0.0,
	// Persistence
	"db.rdbChangesSinceLastSave":           0.0,
	"db.rdbBgsaveInProgress":               0.0,
	"db.rdbLastSaveTime":                   1.510570985e+09,
	"db.rdbLastBgsaveStatus":               "ok", //
	"db.rdbLastBgsaveTimeMilliseconds":     -1.0,
	"db.aofLastRewriteTimeMiliseconds":     -1.0,
	"db.aofLastBgrewriteStatus":            "ok", //
	"db.aofLastWriteStatus":                "ok", //
	"db.aofCurrentRewriteTimeMilliseconds": -1.0,
	"db.asyncLoading":                      0.0,
	"db.currentCowPeak":                    0.0,
	"db.currentCowSize":                    0.0,
	"db.currentCowSizeAge":                 0.0,
	"db.currentForkPerc":                   0.0,
	"db.currentSaveKeysProcessed":          0.0,
	"db.currentSaveKeysTotal":              0.0,
	"db.rdbSaves":                          0.0,
	"db.rdbLastLoadKeysExpired":            0.0,
	"db.rdbLastLoadKeysLoaded":             0.0,
	"db.aofRewrites":                       0.0,
	"db.aofRewritesConsecutiveFailures":    0.0,
	"db.moduleForkInProgress":              0.0,
	"db.moduleForkLastCowSize":             0.0,
	// Stats
	"net.connectionsReceivedPerSecond": 0.0,
	"net.commandsProcessedPerSecond":   0.0,
	"net.inputBytesPerSecond":          0.0,
	"net.outputBytesPerSecond":         0.0,
	"net.rejectedConnectionsPerSecond": 0.0,
	"db.syncFull":                      0.0,
	"db.syncPartialOk":                 0.0,
	"db.syncPartialErr":                0.0,
	"db.expiredKeysPerSecond":          0.0,
	"db.evictedKeysPerSecond":          0.0,
	"db.keyspaceHitsPerSecond":         0.0,
	"db.keyspaceMissesPerSecond":       0.0,
	"net.pubsubChannels":               0.0,
	"net.pubsubPatterns":               0.0,
	"db.latestForkMilliseconds":        0.0,
	"db.expiredStalePercent":           0.0,
	"db.expiredTimecapReachedCount":    0.0,
	"db.expireCycleCpuMilliseconds":    0.0,
	"db.evictedClients":                0.0,
	"db.totalEvictionExceededTime":     0.0,
	"db.currentEvictionExceededTime":   0.0,
	"db.totalForks":                    0.0,
	"db.totalActiveDefragTime":         0.0,
	"db.currentActiveDefragTime":       0.0,
	"db.trackingTotalKeys":             0.0,
	"db.trackingTotalItems":            0.0,
	"db.trackingTotalPrefixes":         0.0,
	"db.unexpectedErrorReplies":        0.0,
	"db.totalErrorReplies":             0.0,
	"db.dumpPayloadSanitizations":      0.0,
	"db.totalReadsProcessed":           0.0,
	"db.totalWritesProcessed":          0.0,
	"db.ioThreadedReadsProcessed":      0.0,
	"db.ioThreadedWritesProcessed":     0.0,
	"db.replyBufferShrinks":            0.0,
	"db.replyBufferExpands":            0.0,

	// Replication
	"cluster.role":                "master",
	"cluster.connectedSlaves":     0.0,
	"cluster.masterFailoverState": "no-failover", //
	// CPU
	"system.usedCpuSysMilliseconds":          580.0,
	"system.usedCpuUserMilliseconds":         300.0,
	"system.usedCpuSysChildrenMilliseconds":  0.0,
	"system.usedCpuUserChildrenMilliseconds": 0.0,
}

var expectedRawKeyspaceInfoFromSample = map[string]map[string]interface{}{
	"db0": {
		"keys":     1,
		"expires":  1,
		"avg_ttl":  7853,
		"keyspace": "db0",
	},
}

var expectedKeyspaceMetricSetFromSample = map[string]interface{}{
	"db.keys":               1.0,
	"db.expires":            1.0,
	"db.avgTtlMilliseconds": 7853.0,
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
	assert.NoError(t, err)

	attr := attribute.Attr("metricsTestSample", "test")
	ms := metric.NewSet("metricsTestSample", persist.NewInMemoryStore(), attr)

	assert.NoError(t, populateMetrics(ms, rawMetrics, metricsDefinition))
	expectedMetricSetFromSample[attr.Key] = attr.Value
	assert.Equal(t, len(expectedMetricSetFromSample), len(ms.Metrics))
	if !reflect.DeepEqual(expectedMetricSetFromSample, ms.Metrics) {
		t.Error("unexpected metric set")
		for k, v := range ms.Metrics {
			if v != expectedMetricSetFromSample[k] {
				t.Errorf("key: %+v expected: %+v have: %+v", k, expectedMetricSetFromSample[k], ms.Metrics[k])
			}
		}
	}

	for db, ks := range rawKeyspace {
		ms = metric.NewSet(fmt.Sprintf("keyspaceTestSample_%s", db), persist.NewInMemoryStore())
		assert.NoError(t, populateMetrics(ms, ks, keyspaceMetricsDefinition))

		if !reflect.DeepEqual(expectedKeyspaceMetricSetFromSample, ms.Metrics) {
			t.Errorf("unexpected keyspace metric set for %+v", ms.Metrics)
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

	assert.NoError(t, err)

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
	m, err := parseKeyspaceMetrics(dbName, keyspace)
	expectedLength := 4
	expectedMetric := map[string]interface{}{
		"keyspace": "db0",
		"keys":     3,
		"expires":  1,
		"avg_ttl":  354949,
	}

	assert.NoError(t, err)

	if len(m) != expectedLength {
		t.Errorf("Not all values processed, got %d length of the rawInventory, expected: %d", len(m), expectedLength)
	}
	if !reflect.DeepEqual(m, expectedMetric) {
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

func TestPopulateCustomKeysMetric(t *testing.T) {
	rawCustomKeys := map[string]keyInfo{
		"myhash": {keyLength: 1, keyType: "hash"},
		"mylist": {keyLength: 5, keyType: "list"},
	}
	expectedHashKeyName := "db.keyLength.hash.myhash"
	expectedListKeyName := "db.keyLength.list.mylist"
	expectedHashKeyLength := float64(1)
	expectedListKeyLength := float64(5)

	sample := metric.NewSet("RedisKeyspaceSample", persist.NewInMemoryStore())

	populateCustomKeysMetric(sample, rawCustomKeys)

	if sample.Metrics[expectedHashKeyName] != expectedHashKeyLength {
		t.Error()
	}
	if sample.Metrics[expectedListKeyName] != expectedListKeyLength {
		t.Error()
	}
	if len(sample.Metrics) != 3 {
		t.Error()
	}
}
