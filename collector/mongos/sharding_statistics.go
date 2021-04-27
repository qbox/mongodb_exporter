package mongos

import (
	"github.com/prometheus/client_golang/prometheus"
)

const mongosCatalogCachePrefix = "ccache_"

var (
	numDatabaseEntries = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      mongosCatalogCachePrefix + "num_database_entries",
		Help:      "The total number of database entries that are currently in the catalog cache.",
	}, []string{})
	numCollectionEntries = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      mongosCatalogCachePrefix + "num_collection_entries",
		Help:      "The total number of database entries that are currently in the catalog cache.",
	}, []string{})
	catalogCacheCountStaleConfigErrors = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      mongosCatalogCachePrefix + "count_stale_config_errors",
		Help:      "The total number of database entries that are currently in the catalog cache.",
	}, []string{})
	totalRefreshWaitTimeMicros = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      mongosCatalogCachePrefix + "total_refresh_wait_time_micros",
		Help:      "The cumulative time, in microseconds, that threads had to wait for a refresh of the metadata.",
	}, []string{})
	numActiveIncrementalRefreshes = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      mongosCatalogCachePrefix + "num_active_incremental_refreshes",
		Help:      "The number of incremental catalog cache refreshes that are currently waiting to complete.",
	}, []string{})
	countIncrementalRefreshesStarted = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      mongosCatalogCachePrefix + "count_incremental_refreshes_started",
		Help:      "The cumulative number of incremental refreshes that have started.",
	}, []string{})
	numActiveFullRefreshes = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      mongosCatalogCachePrefix + "num_active_full_refreshes",
		Help:      "The number of full catalog cache refreshes that are currently waiting to complete.",
	}, []string{})
	countFullRefreshesStarted = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      mongosCatalogCachePrefix + "count_full_refreshes_started",
		Help:      "The cumulative number of full refreshes that have started",
	}, []string{})
	countFailedRefreshes = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Name:      mongosCatalogCachePrefix + "count_failed_refreshes",
		Help:      "The cumulative number of full or incremental refreshes that have failed.",
	}, []string{})
)

// ShardingStatistics https://docs.mongodb.com/manual/reference/command/serverStatus/#shardingstatistics
type ShardingStatistics struct {
	CatalogCache catalogCache `bson:"catalogCache,omitempty"`
}

type catalogCache struct {
	NumDatabaseEntries               float64 `bson:"numDatabaseEntries,omitempty"`
	NumCollectionEntries             float64 `bson:"numCollectionEntries,omitempty"`
	CountStaleConfigErrors           float64 `bson:"countStaleConfigErrors,omitempty"`
	TotalRefreshWaitTimeMicros       float64 `bson:"totalRefreshWaitTimeMicros,omitempty"`
	NumActiveIncrementalRefreshes    float64 `bson:"numActiveIncrementalRefreshes,omitempty"`
	CountIncrementalRefreshesStarted float64 `bson:"countIncrementalRefreshesStarted,omitempty"`
	NumActiveFullRefreshes           float64 `bson:"numActiveFullRefreshes,omitempty"`
	CountFullRefreshesStarted        float64 `bson:"countFullRefreshesStarted,omitempty"`
	CountFailedRefreshes             float64 `bson:"countFailedRefreshes,omitempty"`
}

func (s *ShardingStatistics) update() {
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
