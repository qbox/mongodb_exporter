package common

import (
	"context"
	"github.com/percona/mongodb_exporter/collector/util"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// server connections -- all of these!
var (
	syncClientConnectionsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "connpoolstats", "connection_sync"),
		"Corresponds to the total number of client connections to mongo.",
		nil,
		nil,
	)

	numAScopedConnectionsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "connpoolstats", "connections_scoped_sync"),
		"Corresponds to the number of active and stored outgoing scoped synchronous connections from the current instance to other members of the sharded cluster or replica set.",
		nil,
		nil,
	)

	totalInUseDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "connpoolstats", "connections_in_use"),
		"Corresponds to the total number of client connections to mongo currently in use.",
		nil,
		nil,
	)

	totalAvailableDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "connpoolstats", "connections_available"),
		"Corresponds to the total number of client connections to mongo that are currently available.",
		nil,
		nil,
	)

	totalCreatedDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "connpoolstats", "connections_created_total"),
		"Corresponds to the total number of client connections to mongo created since instance start",
		nil,
		nil,
	)

	totalRefreshingDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "connpoolstats", "connections_refreshing"),
		"Corresponds to the total number of client connections to mongo that currency refreshing",
		nil,
		nil,
	)

	totalReqQueueLimitDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "connpoolstats", "c_req_queue_limit"),
		"Corresponds to the total number of req in the queue that wait a connection",
		nil,
		nil,
	)

	totalInUseInPoolDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "connpoolstats", "c_in_use_in_pool"),
		"Corresponds to the total number of client connections to mongo currently in use in a pool.",
		[]string{"pool"},
		nil,
	)

	totalRefreshingInPoolDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "connpoolstats", "c_refreshing_in_pool"),
		"Corresponds to the total number of client connections to mongo that are currently refreshing in a pool.",
		[]string{"pool"},
		nil,
	)

	totalReqQueueLimitInPoolDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "connpoolstats", "c_req_queue_limit_in_pool"),
		"Corresponds to the total number of req in the queue that wait a connection in a pool",
		[]string{"pool"},
		nil,
	)

	totalAvailableInPoolDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "connpoolstats", "c_available_in_pool"),
		"Corresponds to the total number of client connections to mongo that are currently available in a pool.",
		[]string{"pool"},
		nil,
	)

	totalCreatedInPoolDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "connpoolstats", "c_created_in_pool_total"),
		"Corresponds to the total number of client connections to mongo created since instance start in a pool",
		[]string{"pool"},
		nil,
	)

	inUsePreShardDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "connpoolstats", "c_in_use_preshard"),
		"Corresponds to the pre shard number of client connections to mongo currently in use.",
		[]string{"set", "pool"},
		nil,
	)

	availablePreShardDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "connpoolstats", "c_available_preshard"),
		"Corresponds to the pre shard number of client connections to mongo that are currently available.",
		[]string{"set", "pool"},
		nil,
	)

	refreshingPreShardDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "connpoolstats", "c_refreshing_preshard"),
		"Corresponds to the pre shard number of client connections to mongo that are currently refreshing.",
		[]string{"set", "pool"},
		nil,
	)
	reqQueueLimitPreShardDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "connpoolstats", "c_req_queue_limit_preshard"),
		"Corresponds to the total number of req in the queue that wait a connection in a pool pre shard",
		[]string{"set", "pool"},
		nil,
	)

	createdPreShardDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "connpoolstats", "c_created_preshard_total"),
		"Corresponds to the pre shard number of client connections to mongo created since instance start",
		[]string{"set", "pool"},
		nil,
	)

	getShardHostLimitDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "connpoolstats", "c_refresh_limit"),
		"get shard host limit for a shard",
		[]string{"name"},
		nil,
	)
)

// ConnPoolStats keeps the data returned by the connPoolStats command.
type ConnPoolStats struct {
	SyncClientConnections float64
	ASScopedConnections   float64
	TotalInUse            float64
	TotalAvailable        float64
	TotalCreated          float64
	TotalRefreshing       float64
	TotalReqQueueLimit    float64
	Pools                 map[string]*util.PoolConnStats
	ReplicaSets           map[string]util.ReplicaSetStats
}

type connPoolStats struct {
	SyncClientConnections float64                         `bson:"numClientConnections"`
	ASScopedConnections   float64                         `bson:"numAScopedConnections"`
	TotalInUse            float64                         `bson:"totalInUse"`
	TotalAvailable        float64                         `bson:"totalAvailable"`
	TotalCreated          float64                         `bson:"totalCreated"`
	TotalRefreshing       float64                         `bson:"totalRefreshing"`
	TotalReqQueueLimit    float64                         `bson:"totalReqQueueLimit,omitempty"`
	Pools                 map[string]interface{}          `bson:"pools,omitempty"`
	ReplicaSets           map[string]util.ReplicaSetStats `bson:"replicaSets,omitempty"`

	Ok float64 `bson:"ok"`
}

// Export exports the server status to be consumed by prometheus.
func (stats *ConnPoolStats) Export(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(syncClientConnectionsDesc, prometheus.GaugeValue, stats.SyncClientConnections)
	ch <- prometheus.MustNewConstMetric(numAScopedConnectionsDesc, prometheus.GaugeValue, stats.ASScopedConnections)

	ch <- prometheus.MustNewConstMetric(totalInUseDesc, prometheus.GaugeValue, stats.TotalInUse)
	ch <- prometheus.MustNewConstMetric(totalAvailableDesc, prometheus.GaugeValue, stats.TotalAvailable)
	ch <- prometheus.MustNewConstMetric(totalCreatedDesc, prometheus.CounterValue, stats.TotalCreated)
	ch <- prometheus.MustNewConstMetric(totalRefreshingDesc, prometheus.GaugeValue, stats.TotalRefreshing)
	ch <- prometheus.MustNewConstMetric(totalReqQueueLimitDesc, prometheus.GaugeValue, stats.TotalReqQueueLimit)

	for k, v := range stats.ReplicaSets {
		if len(k) == 0 {
			continue
		}
		ch <- prometheus.MustNewConstMetric(getShardHostLimitDesc, prometheus.GaugeValue, v.RefreshLimiter, k)
	}

	for k, v := range stats.Pools {
		if len(k) == 0 || v == nil {
			continue
		}

		ch <- prometheus.MustNewConstMetric(totalAvailableInPoolDesc, prometheus.GaugeValue, v.PoolAvailable, k)
		ch <- prometheus.MustNewConstMetric(totalCreatedInPoolDesc, prometheus.CounterValue, v.PoolCreated, k)
		ch <- prometheus.MustNewConstMetric(totalInUseInPoolDesc, prometheus.GaugeValue, v.PoolInUse, k)
		ch <- prometheus.MustNewConstMetric(totalRefreshingInPoolDesc, prometheus.GaugeValue, v.PoolRefreshing, k)
		ch <- prometheus.MustNewConstMetric(totalReqQueueLimitInPoolDesc, prometheus.GaugeValue, v.PoolReqQueueLimit, k)

		for shardName, shardStat := range v.ShardStatus {
			if len(shardName) == 0 || shardStat == nil {
				continue
			}

			ch <- prometheus.MustNewConstMetric(availablePreShardDesc, prometheus.GaugeValue, shardStat.Available, shardName, k)
			ch <- prometheus.MustNewConstMetric(createdPreShardDesc, prometheus.CounterValue, shardStat.Created, shardName, k)
			ch <- prometheus.MustNewConstMetric(inUsePreShardDesc, prometheus.GaugeValue, shardStat.InUse, shardName, k)
			ch <- prometheus.MustNewConstMetric(refreshingPreShardDesc, prometheus.GaugeValue, shardStat.Refreshing, shardName, k)
			ch <- prometheus.MustNewConstMetric(reqQueueLimitPreShardDesc, prometheus.GaugeValue, shardStat.ReqQueueLimit, shardName, k)
		}
	}
}

// Describe describes the server status for prometheus.
func (stats *ConnPoolStats) Describe(ch chan<- *prometheus.Desc) {
	ch <- syncClientConnectionsDesc
	ch <- numAScopedConnectionsDesc

	ch <- totalInUseDesc
	ch <- totalAvailableDesc
	ch <- totalCreatedDesc
	ch <- totalRefreshingDesc
	ch <- totalReqQueueLimitDesc

	ch <- totalReqQueueLimitInPoolDesc
	ch <- totalCreatedInPoolDesc
	ch <- totalAvailableInPoolDesc
	ch <- totalRefreshingInPoolDesc
	ch <- totalInUseInPoolDesc

	ch <- availablePreShardDesc
	ch <- createdPreShardDesc
	ch <- refreshingPreShardDesc
	ch <- refreshingPreShardDesc
	ch <- inUsePreShardDesc

	ch <- getShardHostLimitDesc
}

// GetConnPoolStats returns the server connPoolStats info.
func GetConnPoolStats(client *mongo.Client, hostToShardName map[string]string) *ConnPoolStats {
	result := &connPoolStats{}
	err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"connPoolStats", 1}, {"recordStats", 0}}).Decode(result)
	if err != nil {
		log.Errorf("Failed to get connPoolStats: %s.", err)
		return nil
	}

	if int(result.Ok) != 1 {
		log.Errorf("failed to get connPoolStats because ok != 1")
		return nil
	}
	tmp := ConnPoolStats{
		TotalAvailable:        result.TotalAvailable,
		TotalRefreshing:       result.TotalRefreshing,
		TotalCreated:          result.TotalCreated,
		TotalInUse:            result.TotalInUse,
		TotalReqQueueLimit:    result.TotalReqQueueLimit,
		SyncClientConnections: result.SyncClientConnections,
		ASScopedConnections:   result.ASScopedConnections,
		Pools:                 make(map[string]*util.PoolConnStats),
		ReplicaSets:           result.ReplicaSets,
	}

	if hostToShardName != nil && len(hostToShardName) != 0 {
		poolsStat := util.PoolsAnaliyze(result.Pools, hostToShardName)
		if poolsStat != nil {
			tmp.Pools = poolsStat
		}
	}
	return &tmp
}
