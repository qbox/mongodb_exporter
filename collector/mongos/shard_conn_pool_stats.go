package mongos

import (
	"context"
	"github.com/percona/mongodb_exporter/collector/util"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ShardConnPoolStats struct {
	TotalInUse     float64
	TotalAvailable float64
	TotalCreated   float64
	GlobalPool     *util.PoolConnStats
}

type shardConnPoolStats struct {
	TotalInUse     float64                `bson:"totalInUse"`
	TotalAvailable float64                `bson:"totalAvailable"`
	TotalCreated   float64                `bson:"totalCreated"`
	Pools          map[string]interface{} `bson:"pools,omitempty"`

	Ok float64 `bson:"ok"`
}

var (
	shardTotalInUseDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "sc_stats", "in_use"),
		"Corresponds to the total number of client connections to shard mongo currently in use.",
		nil,
		nil,
	)

	shardTotalAvailableDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "sc_stats", "available"),
		"Corresponds to the total number of client connections to shard mongo that are currently available.",
		nil,
		nil,
	)

	shardTotalCreatedDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "sc_stats", "created_total"),
		"Corresponds to the total number of client connections to shard mongo created since instance start",
		nil,
		nil,
	)

	shardTotalInUseInPoolDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "sc_stats", "in_use_in_pool"),
		"Corresponds to the total number of client connections to shard mongo currently in use in a pool.",
		[]string{"pool"},
		nil,
	)

	shardTotalAvailableInPoolDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "sc_stats", "available_in_pool"),
		"Corresponds to the total number of client connections to shard mongo that are currently available in a pool.",
		[]string{"pool"},
		nil,
	)

	shardTotalCreatedInPoolDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "sc_stats", "created_in_pool_total"),
		"Corresponds to the total number of client connections to shard mongo created since instance start in a pool",
		[]string{"pool"},
		nil,
	)

	shardInUsePreShardDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "sc_stats", "in_use_preshard"),
		"Corresponds to the pre shard number of client connections to shard mongo currently in use.",
		[]string{"set", "pool"},
		nil,
	)

	shardAvailablePreShardDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "sc_stats", "available_preshard"),
		"Corresponds to the pre shard number of client connections to shard mongo that are currently available.",
		[]string{"set", "pool"},
		nil,
	)

	shardCreatedPreShardDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "sc_stats", "created_preshard_total"),
		"Corresponds to the pre shard number of client connections to shard mongo created since instance start",
		[]string{"set", "pool"},
		nil,
	)
)

func (stats *ShardConnPoolStats) Export(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(shardTotalInUseDesc, prometheus.GaugeValue, stats.TotalInUse)
	ch <- prometheus.MustNewConstMetric(shardTotalCreatedDesc, prometheus.CounterValue, stats.TotalCreated)
	ch <- prometheus.MustNewConstMetric(shardTotalAvailableDesc, prometheus.GaugeValue, stats.TotalInUse)

	if stats.GlobalPool != nil {
		ch <- prometheus.MustNewConstMetric(shardTotalAvailableInPoolDesc, prometheus.GaugeValue, stats.GlobalPool.PoolAvailable, "global")
		ch <- prometheus.MustNewConstMetric(shardTotalCreatedInPoolDesc, prometheus.CounterValue, stats.GlobalPool.PoolCreated, "global")
		ch <- prometheus.MustNewConstMetric(shardTotalInUseInPoolDesc, prometheus.GaugeValue, stats.GlobalPool.PoolInUse, "global")

		for k, v := range stats.GlobalPool.ShardStatus {
			if len(k) == 0 || v == nil {
				continue
			}
			ch <- prometheus.MustNewConstMetric(shardAvailablePreShardDesc, prometheus.GaugeValue, v.Available, k, "global")
			ch <- prometheus.MustNewConstMetric(shardCreatedPreShardDesc, prometheus.CounterValue, v.Created, k, "global")
			ch <- prometheus.MustNewConstMetric(shardInUsePreShardDesc, prometheus.GaugeValue, v.InUse, k, "global")
		}
	}
}

func (stats *ShardConnPoolStats) Describe(ch chan<- *prometheus.Desc) {
	ch <- shardTotalAvailableDesc
	ch <- shardTotalCreatedDesc
	ch <- shardTotalInUseDesc

	ch <- shardCreatedPreShardDesc
	ch <- shardAvailablePreShardDesc
	ch <- shardInUsePreShardDesc

	ch <- shardTotalCreatedInPoolDesc
	ch <- shardTotalInUseInPoolDesc
	ch <- shardTotalAvailableInPoolDesc
}

func GetShardConnPoolStats(client *mongo.Client, hostToShardname map[string]string) *ShardConnPoolStats {
	result := &shardConnPoolStats{}
	err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"shardConnPoolStats", 1}}).Decode(result)
	if err != nil || int(result.Ok) != 1 {
		log.Errorf("Failed to get shardConnPoolStats: %s.", err)
		return nil
	}

	//解析shardConnPoolStats to ShardConnPoolStats
	tmp := &ShardConnPoolStats{
		TotalAvailable: result.TotalAvailable,
		TotalCreated:   result.TotalCreated,
		TotalInUse:     result.TotalInUse,
	}

	if hostToShardname != nil && len(hostToShardname) != 0 {
		poolsStat := util.PoolsAnaliyze(result.Pools, hostToShardname)
		if poolsStat != nil {
			if global, ok := poolsStat["global"]; ok {
				tmp.GlobalPool = global
			}
		}
	}
	return tmp
}
