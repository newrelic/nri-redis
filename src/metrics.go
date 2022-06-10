package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/log"
)

var metricsDefinition = map[string][]interface{}{
	// Server
	"software.uptimeMilliseconds":     {secondsToMilliseconds("uptime_in_seconds"), metric.GAUGE},
	"software.version":                {"redis_version", metric.ATTRIBUTE},
	"software.monotonicClock":         {"monotonic_clock", metric.ATTRIBUTE},
	"software.processSupervised":      {"process_supervised", metric.ATTRIBUTE},
	"software.serverTimeMilliseconds": {microsecondsToMilliseconds("server_time_usec"), metric.GAUGE},
	"software.configuredHz":           {"configured_hz", metric.GAUGE},
	"software.ioThreadsActive":        {"io_threads_active", metric.GAUGE},
	// Clients
	"net.connectedClients":            {"connected_clients", metric.GAUGE},
	"net.clientLongestOutputList":     {"client_longest_output_list", metric.GAUGE},
	"net.clientBiggestInputBufBytes":  {"client_biggest_input_buf", metric.GAUGE},
	"net.blockedClients":              {"blocked_clients", metric.GAUGE},
	"net.clusterConnections":          {"cluster_connections", metric.GAUGE},
	"net.maxClients":                  {"maxclients", metric.GAUGE},
	"net.clientRecentMaxInputBuffer":  {"client_recent_max_input_buffer", metric.GAUGE},
	"net.clientRecentMaxOutputBuffer": {"client_recent_max_output_buffer", metric.GAUGE},
	"net.trackingClients":             {"tracking_clients", metric.GAUGE},
	"net.clientsInTimeoutTable":       {"clients_in_timeout_table", metric.GAUGE},
	// Memory
	"system.usedMemoryBytes":            {"used_memory", metric.GAUGE},
	"system.usedMemoryRssBytes":         {"used_memory_rss", metric.GAUGE},
	"system.usedMemoryPeakBytes":        {"used_memory_peak", metric.GAUGE},
	"system.usedMemoryLuaBytes":         {"used_memory_lua", metric.GAUGE},
	"system.usedMemoryHuman":            {"used_memory_human", metric.ATTRIBUTE},
	"system.usedMemoryOverhead":         {"used_memory_overhead", metric.GAUGE},
	"system.totalSystemMemoryBytes":     {"total_system_memory", metric.GAUGE},
	"system.maxmemoryBytes":             {"maxmemory", metric.GAUGE},
	"system.memFragmentationRatio":      {"mem_fragmentation_ratio", metric.GAUGE},
	"system.allocatorAllocated":         {"allocator_allocated", metric.GAUGE},
	"system.allocatorActive":            {"allocator_active", metric.GAUGE},
	"system.allocatorResident":          {"allocator_resident", metric.GAUGE},
	"system.usedMemoryVmEval":           {"used_memory_vm_eval", metric.GAUGE},
	"system.memAllocator":               {"mem_allocator", metric.ATTRIBUTE},
	"system.usedMemoryScriptsEval":      {"used_memory_scripts_eval", metric.GAUGE},
	"system.numberOfCachedScripts":      {"number_of_cached_scripts", metric.GAUGE},
	"system.numberOfFunctions":          {"number_of_functions", metric.GAUGE},
	"system.numberOfLibraries":          {"number_of_libraries", metric.GAUGE},
	"system.usedMemoryVmFunctions":      {"used_memory_vm_functions", metric.GAUGE},
	"system.usedMemoryVmTotal":          {"used_memory_vm_total", metric.GAUGE},
	"system.usedMemoryVmTotalHuman":     {"used_memory_vm_total_human", metric.ATTRIBUTE},
	"system.usedMemoryFunctions":        {"used_memory_functions", metric.GAUGE},
	"system.usedMemoryScripts":          {"used_memory_scripts", metric.GAUGE},
	"system.usedMemoryScriptsHuman":     {"used_memory_scripts_human", metric.ATTRIBUTE},
	"system.usedMemoryPeakHuman":        {"used_memory_peak_human", metric.ATTRIBUTE},
	"system.allocatorFragRatio":         {"allocator_frag_ratio", metric.GAUGE},
	"system.allocatorFragBytes":         {"allocator_frag_bytes", metric.GAUGE},
	"system.allocatorRssRatio":          {"allocator_rss_ratio", metric.GAUGE},
	"system.allocatorRssBytes":          {"allocator_rss_bytes", metric.GAUGE},
	"system.rssOverheadRatio":           {"rss_overhead_ratio", metric.GAUGE},
	"system.memFragmentationBytes":      {"mem_fragmentation_bytes", metric.GAUGE},
	"system.memNotCountedForEvict":      {"mem_not_counted_for_evict", metric.GAUGE},
	"system.memReplicationBacklog":      {"mem_replication_backlog", metric.GAUGE},
	"system.memTotalReplicationBuffers": {"mem_total_replication_buffers", metric.GAUGE},
	"system.memClientsSlaves":           {"mem_clients_slaves", metric.GAUGE},
	"system.memClientsNormal":           {"mem_clients_normal", metric.GAUGE},
	"system.memClusterLinks":            {"mem_cluster_links", metric.GAUGE},
	"system.memAofBuffer":               {"mem_aof_buffer", metric.GAUGE},
	"system.lazyfreedObjects":           {"lazyfreed_objects", metric.GAUGE},
	// Persistence
	"db.rdbChangesSinceLastSave":           {"rdb_changes_since_last_save", metric.GAUGE},
	"db.rdbBgsaveInProgress":               {"rdb_bgsave_in_progress", metric.GAUGE},
	"db.rdbLastSaveTime":                   {"rdb_last_save_time", metric.GAUGE},
	"db.rdbLastBgsaveStatus":               {"rdb_last_bgsave_status", metric.ATTRIBUTE},
	"db.rdbLastBgsaveTimeMilliseconds":     {secondsToMilliseconds("rdb_last_bgsave_time_sec"), metric.GAUGE},
	"db.aofLastRewriteTimeMiliseconds":     {secondsToMilliseconds("aof_last_rewrite_time_sec"), metric.GAUGE},
	"db.aofCurrentRewriteTimeMilliseconds": {secondsToMilliseconds("aof_current_rewrite_time_sec"), metric.GAUGE},
	"db.aofLastBgrewriteStatus":            {"aof_last_bgrewrite_status", metric.ATTRIBUTE},
	"db.aofLastWriteStatus":                {"aof_last_write_status", metric.ATTRIBUTE},
	"db.asyncLoading":                      {"async_loading", metric.GAUGE},
	"db.currentCowPeak":                    {"current_cow_peak", metric.GAUGE},
	"db.currentCowSize":                    {"current_cow_size", metric.GAUGE},
	"db.currentCowSizeAge":                 {"current_cow_size_age", metric.GAUGE},
	"db.currentForkPerc":                   {"current_fork_perc", metric.GAUGE},
	"db.currentSaveKeysProcessed":          {"current_save_keys_processed", metric.GAUGE},
	"db.currentSaveKeysTotal":              {"current_save_keys_total", metric.GAUGE},
	"db.rdbSaves":                          {"rdb_saves", metric.GAUGE},
	"db.rdbLastLoadKeysExpired":            {"rdb_last_load_keys_expired", metric.GAUGE},
	"db.rdbLastLoadKeysLoaded":             {"rdb_last_load_keys_loaded", metric.GAUGE},
	"db.aofRewrites":                       {"aof_rewrites", metric.GAUGE},
	"db.aofRewritesConsecutiveFailures":    {"aof_rewrites_consecutive_failures", metric.GAUGE},
	"db.moduleForkInProgress":              {"module_fork_in_progress", metric.GAUGE},
	"db.moduleForkLastCowSize":             {"module_fork_last_cow_size", metric.GAUGE},
	// Stats
	"net.connectionsReceivedPerSecond": {"total_connections_received", metric.RATE},
	"net.commandsProcessedPerSecond":   {"total_commands_processed", metric.RATE},
	"net.inputBytesPerSecond":          {"total_net_input_bytes", metric.RATE},
	"net.outputBytesPerSecond":         {"total_net_output_bytes", metric.RATE},
	"net.rejectedConnectionsPerSecond": {"rejected_connections", metric.RATE},
	"db.syncFull":                      {"sync_full", metric.GAUGE},
	"db.syncPartialOk":                 {"sync_partial_ok", metric.GAUGE},
	"db.syncPartialErr":                {"sync_partial_err", metric.GAUGE},
	"db.expiredKeysPerSecond":          {"expired_keys", metric.RATE},
	"db.evictedKeysPerSecond":          {"evicted_keys", metric.RATE},
	"db.keyspaceHitsPerSecond":         {"keyspace_hits", metric.RATE},
	"db.keyspaceMissesPerSecond":       {"keyspace_misses", metric.RATE},
	"net.pubsubChannels":               {"pubsub_channels", metric.GAUGE},
	"net.pubsubPatterns":               {"pubsub_patterns", metric.GAUGE},
	"db.latestForkMilliseconds":        {microsecondsToMilliseconds("latest_fork_usec"), metric.GAUGE},
	"db.expiredStalePercent":           {"expired_stale_perc", metric.GAUGE},
	"db.expiredTimecapReachedCount":    {"expired_time_cap_reached_count", metric.GAUGE},
	"db.expireCycleCpuMilliseconds":    {secondsToMilliseconds("expire_cycle_cpu_milliseconds"), metric.GAUGE},
	"db.evictedClients":                {"evicted_clients", metric.GAUGE},
	"db.totalEvictionExceededTime":     {"total_eviction_exceeded_time", metric.GAUGE},
	"db.currentEvictionExceededTime":   {"current_eviction_exceeded_time", metric.GAUGE},
	"db.totalForks":                    {"total_forks", metric.GAUGE},
	"db.totalActiveDefragTime":         {"total_active_defrag_time", metric.GAUGE},
	"db.currentActiveDefragTime":       {"current_active_defrag_time", metric.GAUGE},
	"db.trackingTotalKeys":             {"tracking_total_keys", metric.GAUGE},
	"db.trackingTotalItems":            {"tracking_total_items", metric.GAUGE},
	"db.trackingTotalPrefixes":         {"tracking_total_prefixes", metric.GAUGE},
	"db.unexpectedErrorReplies":        {"unexpected_error_replies", metric.GAUGE},
	"db.totalErrorReplies":             {"total_error_replies", metric.GAUGE},
	"db.dumpPayloadSanitizations":      {"dump_payload_sanitizations", metric.GAUGE},
	"db.totalReadsProcessed":           {"total_reads_processed", metric.GAUGE},
	"db.totalWritesProcessed":          {"total_writes_processed", metric.GAUGE},
	"db.ioThreadedReadsProcessed":      {"io_threaded_reads_processed", metric.GAUGE},
	"db.ioThreadedWritesProcessed":     {"io_threaded_writes_processed", metric.GAUGE},
	"db.replyBufferShrinks":            {"reply_buffer_shrinks", metric.GAUGE},
	"db.replyBufferExpands":            {"reply_buffer_expands", metric.GAUGE},

	// Replication
	"cluster.role":                {"role", metric.ATTRIBUTE},
	"cluster.connectedSlaves":     {"connected_slaves", metric.GAUGE},
	"cluster.masterFailoverState": {"master_failover_state", metric.ATTRIBUTE},
	// CPU
	"system.usedCpuSysMilliseconds":          {secondsToMilliseconds("used_cpu_sys"), metric.GAUGE},
	"system.usedCpuUserMilliseconds":         {secondsToMilliseconds("used_cpu_user"), metric.GAUGE},
	"system.usedCpuSysChildrenMilliseconds":  {secondsToMilliseconds("used_cpu_sys_children"), metric.GAUGE},
	"system.usedCpuUserChildrenMilliseconds": {secondsToMilliseconds("used_cpu_user_children"), metric.GAUGE},
}

var keyspaceMetricsDefinition = map[string][]interface{}{
	"db.keys":               {"keys", metric.GAUGE},
	"db.expires":            {"expires", metric.GAUGE},
	"db.avgTtlMilliseconds": {"avg_ttl", metric.GAUGE},
	"db.keyspace":           {"keyspace", metric.ATTRIBUTE},
}

type manipulator func(map[string]interface{}) (interface{}, bool)

func populateMetrics(sample *metric.Set, metrics map[string]interface{}, definition map[string][]interface{}) error {
	notFoundMetric := make([]string, 0)

	for metricName, metricInfo := range definition {
		rawSource := metricInfo[0]
		metricType := metricInfo[1].(metric.SourceType)

		var rawMetric interface{}
		var ok bool

		switch source := rawSource.(type) {
		case string:
			rawMetric, ok = metrics[source]
		case manipulator:
			rawMetric, ok = source(metrics)
		default:
			log.Warn("Invalid raw source metric for %s", metricName)
			continue
		}

		if !ok {
			notFoundMetric = append(notFoundMetric, metricName)
			continue
		}
		err := sample.SetMetric(metricName, rawMetric, metricType)
		if err != nil {
			log.Warn("Error setting value: %s", err)
			continue
		}
	}
	if len(notFoundMetric) > 0 {
		log.Warn("Can't find raw metrics in results for %v", notFoundMetric)
	}
	return nil
}

func asValue(value string) interface{} {
	value = strings.TrimRight(value, "\n")

	if i, err := strconv.Atoi(value); err == nil {
		if i == 0 && len(value) > 1 {
			// "It is a valid string with zeros like 000000000"
			return value
		}
		return i
	}

	if f, err := strconv.ParseFloat(value, 64); err == nil {
		return f
	}

	if b, err := strconv.ParseBool(value); err == nil {
		return b
	}

	return value
}

func getRawMetrics(info string) (map[string]interface{}, map[string]map[string]interface{}, error) {
	metrics := make(map[string]interface{})
	keyspaceMetrics := make(map[string]map[string]interface{})

	reader := bufio.NewReader(strings.NewReader(info))
	_, err := reader.Peek(1)
	if err != nil {
		return nil, nil, fmt.Errorf("Empty result")
	}
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
		}

		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}
		if strings.HasPrefix(parts[0], "db") {
			keyspaceMetrics[parts[0]], err = parseKeyspaceMetrics(parts[0], parts[1])
			if err != nil {
				return nil, nil, err
			}
		} else {
			value := strings.TrimSuffix(parts[1], "\r\n")
			metrics[parts[0]] = asValue(value)
		}
	}

	if len(metrics) == 0 {
		log.Debug("Empty result for non-keyspace metrics")
	}
	if len(keyspaceMetrics) == 0 {
		log.Debug("Empty result for keyspace metrics")
	}

	return metrics, keyspaceMetrics, nil
}

func parseKeyspaceMetrics(dbName string, keyspace string) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	re, err := regexp.Compile(`keys=(\d+),expires=(\d+),avg_ttl=(\d+\.*\d*)`)
	if err != nil {
		return nil, err
	}
	matches := re.FindStringSubmatch(keyspace)
	if matches != nil {
		m["keyspace"] = asValue(dbName)
		m["keys"] = asValue(matches[1])
		m["expires"] = asValue(matches[2])
		m["avg_ttl"] = asValue(matches[3])
	} else {
		log.Warn("Keyspace metrics cannot be parsed for %s", dbName)
	}

	return m, nil
}

func populateCustomKeysMetric(sample *metric.Set, rawCustomKeys map[string]keyInfo) {
	for key, value := range rawCustomKeys {
		err := sample.SetMetric(fmt.Sprintf("db.keyLength.%s.%s", value.keyType, key), value.keyLength, metric.GAUGE)
		if err != nil {
			log.Warn("Error setting value: %s", err)
		}
	}
}

func microsecondsToMilliseconds(source string) manipulator {
	return func(metrics map[string]interface{}) (interface{}, bool) {
		if metrics[source] == 0 || metrics[source] == -1 {
			return metrics[source], true
		}

		secs, ok := metrics[source].(int)
		if ok {
			return secs / 1000, true
		}

		return 0, false
	}
}

func secondsToMilliseconds(source string) manipulator {
	return func(metrics map[string]interface{}) (interface{}, bool) {
		if metrics[source] == 0 || metrics[source] == -1 {
			return metrics[source], true
		}

		switch metrics[source].(type) {
		case int:
			return metrics[source].(int) * 1000, true
		case float64:
			// NOTE: We return int because redis values are expressed as seconds 2 decimals precision,
			// so once converted to millis would not contain decimal values.
			return int(metrics[source].(float64) * 1000), true
		default:
			return 0, false
		}
	}
}
