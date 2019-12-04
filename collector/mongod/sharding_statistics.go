package mongod

import (
	"github.com/prometheus/client_golang/prometheus"
)

const shardingStatPrefix = "sharding_stat_"
const shardingStatCatalogCachePrefix = shardingStatPrefix + "catalog_cache_"

var (
	countStaleConfigErrors = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      shardingStatPrefix + "stale_config_errors_total",
		Help:      "The total number of times that threads hit stale config exception",
	}, []string{})

	countDonorMoveChunkStarted = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      shardingStatPrefix + "donor_move_chunk_start_total",
		Help:      "The total number of times that the moveChunk command has started on the shard",
	}, []string{})
	totalDonorChunkCloneTimeMillis = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      shardingStatPrefix + "donor_chunk_clone_time_ms_total",
		Help:      "The cumulative time, in milliseconds, taken by the clone phase of the chunk migrations from this shard, of which this node is a member.",
	}, []string{})
	totalCriticalSectionCommitTimeMillis = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      shardingStatPrefix + "critical_section_commit_time_ms_total",
		Help:      "The cumulative time, in milliseconds, taken by the update metadata phase of the chunk migrations from this shard, of which this node is a member.",
	}, []string{})
	totalCriticalSectionTimeMillis = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      shardingStatPrefix + "critical_section_time_ms_total",
		Help:      "The cumulative time, in milliseconds, taken by the catch-up phase and the update metadata phase of the chunk migrations from this shard, of which this node is a member.",
	}, []string{})
	numDatabaseEntries = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      shardingStatCatalogCachePrefix + "num_database_entries",
		Help:      "The total number of database entries that are currently in the catalog cache.",
	}, []string{})
	numCollectionEntries = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      shardingStatCatalogCachePrefix + "num_collection_entries",
		Help:      "The total number of database entries that are currently in the catalog cache.",
	}, []string{})
	catalogCacheCountStaleConfigErrors = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      shardingStatCatalogCachePrefix + "count_stale_config_errors",
		Help:      "The total number of database entries that are currently in the catalog cache.",
	}, []string{})
	totalRefreshWaitTimeMicros = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      shardingStatCatalogCachePrefix + "total_refresh_wait_time_micros",
		Help:      "The cumulative time, in microseconds, that threads had to wait for a refresh of the metadata.",
	}, []string{})
	numActiveIncrementalRefreshes = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      shardingStatCatalogCachePrefix + "num_active_incremental_refreshes",
		Help:      "The number of incremental catalog cache refreshes that are currently waiting to complete.",
	}, []string{})
	countIncrementalRefreshesStarted = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      shardingStatCatalogCachePrefix + "count_incremental_refreshes_started",
		Help:      "The cumulative number of incremental refreshes that have started.",
	}, []string{})
	numActiveFullRefreshes = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      shardingStatCatalogCachePrefix + "num_active_full_refreshes",
		Help:      "The number of full catalog cache refreshes that are currently waiting to complete.",
	}, []string{})
	countFullRefreshesStarted = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      shardingStatCatalogCachePrefix + "count_full_refreshes_started",
		Help:      "The cumulative number of full refreshes that have started",
	}, []string{})
	countFailedRefreshes = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      shardingStatCatalogCachePrefix + "count_failed_refreshes",
		Help:      "The cumulative number of full or incremental refreshes that have failed.",
	}, []string{})
)

// ShardingStatistics https://docs.mongodb.com/manual/reference/command/serverStatus/#shardingstatistics
type ShardingStatistics struct {
	CountStaleConfigErrors               float64      `bson:"countStaleConfigErrors"`
	CountDonorMoveChunkStarted           float64      `bson:"countDonorMoveChunkStarted"`
	TotalDonorChunkCloneTimeMillis       float64      `bson:"totalDonorChunkCloneTimeMillis"`
	TotalCriticalSectionCommitTimeMillis float64      `bson:"totalCriticalSectionCommitTimeMillis"`
	TotalCriticalSectionTimeMillis       float64      `bson:"totalCriticalSectionTimeMillis"`
	CatalogCache                         catalogCache `bson:"catalogCache"`
}

type catalogCache struct {
	NumDatabaseEntries               float64 `bson:"numDatabaseEntries"`
	NumCollectionEntries             float64 `bson:"numCollectionEntries"`
	CountStaleConfigErrors           float64 `bson:"countStaleConfigErrors"`
	TotalRefreshWaitTimeMicros       float64 `bson:"totalRefreshWaitTimeMicros"`
	NumActiveIncrementalRefreshes    float64 `bson:"numActiveIncrementalRefreshes"`
	CountIncrementalRefreshesStarted float64 `bson:"countIncrementalRefreshesStarted"`
	NumActiveFullRefreshes           float64 `bson:"numActiveFullRefreshes"`
	CountFullRefreshesStarted        float64 `bson:"countFullRefreshesStarted"`
	CountFailedRefreshes             float64 `bson:"countFailedRefreshes"`
}

func (s *ShardingStatistics) update() {
	countStaleConfigErrors.WithLabelValues().Set(s.CountStaleConfigErrors)
	countDonorMoveChunkStarted.WithLabelValues().Set(s.CountDonorMoveChunkStarted)
	totalDonorChunkCloneTimeMillis.WithLabelValues().Set(s.TotalDonorChunkCloneTimeMillis)
	totalCriticalSectionCommitTimeMillis.WithLabelValues().Set(s.TotalCriticalSectionCommitTimeMillis)
	totalCriticalSectionTimeMillis.WithLabelValues().Set(s.TotalCriticalSectionTimeMillis)

	s.CatalogCache.update()
}

func (c *catalogCache) update() {
	numDatabaseEntries.WithLabelValues().Set(c.NumDatabaseEntries)
	numCollectionEntries.WithLabelValues().Set(c.NumCollectionEntries)
	catalogCacheCountStaleConfigErrors.WithLabelValues().Set(c.CountStaleConfigErrors)
	totalRefreshWaitTimeMicros.WithLabelValues().Set(c.TotalRefreshWaitTimeMicros)
	numActiveIncrementalRefreshes.WithLabelValues().Set(c.NumActiveIncrementalRefreshes)
	countIncrementalRefreshesStarted.WithLabelValues().Set(c.CountIncrementalRefreshesStarted)
	numActiveFullRefreshes.WithLabelValues().Set(c.NumActiveFullRefreshes)
	countFullRefreshesStarted.WithLabelValues().Set(c.CountFullRefreshesStarted)
	countFailedRefreshes.WithLabelValues().Set(c.CountFailedRefreshes)
}

// Export exports the data to prometheus.
func (s *ShardingStatistics) Export(ch chan<- prometheus.Metric) {
	s.update()
	countStaleConfigErrors.Collect(ch)
	countDonorMoveChunkStarted.Collect(ch)
	totalDonorChunkCloneTimeMillis.Collect(ch)
	totalCriticalSectionCommitTimeMillis.Collect(ch)
	totalCriticalSectionTimeMillis.Collect(ch)

	s.CatalogCache.Export(ch)
}

// Export exports the data to prometheus.
func (c *catalogCache) Export(ch chan<- prometheus.Metric) {
	numDatabaseEntries.Collect(ch)
	numCollectionEntries.Collect(ch)
	catalogCacheCountStaleConfigErrors.Collect(ch)
	totalRefreshWaitTimeMicros.Collect(ch)
	numActiveIncrementalRefreshes.Collect(ch)
	countIncrementalRefreshesStarted.Collect(ch)
	numActiveFullRefreshes.Collect(ch)
	countFullRefreshesStarted.Collect(ch)
	countFailedRefreshes.Collect(ch)
}

// Describe describes the metrics for prometheus
func (s *ShardingStatistics) Describe(ch chan<- *prometheus.Desc) {
	countStaleConfigErrors.Describe(ch)
	countDonorMoveChunkStarted.Describe(ch)
	totalDonorChunkCloneTimeMillis.Describe(ch)
	totalCriticalSectionCommitTimeMillis.Describe(ch)
	totalCriticalSectionTimeMillis.Describe(ch)

	s.CatalogCache.Describe(ch)
}

// Describe describes the metrics for prometheus
func (c *catalogCache) Describe(ch chan<- *prometheus.Desc) {
	numDatabaseEntries.Describe(ch)
	numCollectionEntries.Describe(ch)
	catalogCacheCountStaleConfigErrors.Describe(ch)
	totalRefreshWaitTimeMicros.Describe(ch)
	numActiveIncrementalRefreshes.Describe(ch)
	countIncrementalRefreshesStarted.Describe(ch)
	numActiveFullRefreshes.Describe(ch)
	countFailedRefreshes.Describe(ch)
}
