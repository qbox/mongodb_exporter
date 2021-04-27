package util

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestHostConnStats_Add(t *testing.T) {
	tmp := HostConnStats{
		InUse: 1,
		Available: 2,
		Created: 3,
		Refreshing: 4,
		ReqQueueLimit: 5,
	}

	tmp2 := HostConnStats{
		InUse: 5,
		Available: 4,
		Created: 3,
		Refreshing: 2,
		ReqQueueLimit: 1,
	}

	tmp.Add(&tmp2)

	assert.True(t, true, reflect.DeepEqual(tmp, HostConnStats{
		InUse: 6,
		Available: 6,
		Created: 6,
		Refreshing: 6,
		ReqQueueLimit: 6,
	}))
}

func TestGetPoolConnStats(t *testing.T) {
	data := "{\n\t\t\t\"poolInUse\" : 1,\n\t\t\t\"poolAvailable\" : 3,\n\t\t\t\"poolCreated\" : 614,\n\t\t\t\"poolRefreshing\" : 0,\n\t\t\t\"poolReqQueueLimit\" : 0,\n\t\t\t\"10.34.62.46:15350\" : {\n\t\t\t\t\"inUse\" : 0,\n\t\t\t\t\"available\" : 1,\n\t\t\t\t\"created\" : 66,\n\t\t\t\t\"refreshing\" : 0,\n\t\t\t\t\"reqQueueLimit\" : 0\n\t\t\t},\n\t\t\t\"10.34.60.45:15350\" : {\n\t\t\t\t\"inUse\" : 0,\n\t\t\t\t\"available\" : 1,\n\t\t\t\t\"created\" : 56,\n\t\t\t\t\"refreshing\" : 0,\n\t\t\t\t\"reqQueueLimit\" : 0\n\t\t\t},\n\t\t\t\"10.34.60.43:15350\" : {\n\t\t\t\t\"inUse\" : 1,\n\t\t\t\t\"available\" : 1,\n\t\t\t\t\"created\" : 492,\n\t\t\t\t\"refreshing\" : 0,\n\t\t\t\t\"reqQueueLimit\" : 0\n\t\t\t}\n\t\t}"

	res := make(map[string]interface{})
	err := json.Unmarshal([]byte(data), &res)
	assert.Equal(t, nil, err, "json must be correct")

	pool := GetPoolConnStats(res)
	assert.NotEqual(t, nil, pool, "pools must be not nil")

	assert.Equal(t, 1, int(pool.PoolInUse))
	assert.Equal(t, 3, int(pool.PoolAvailable))
	assert.Equal(t, 614, int(pool.PoolCreated))
	assert.Equal(t, 0, int(pool.PoolRefreshing))
	assert.Equal(t, 0, int(pool.PoolReqQueueLimit))
	assert.Equal(t, 3, len(pool.ShardStatus))

	if _, ok := pool.ShardStatus["10.34.60.45:15350"]; !ok {
		assert.Equal(t, true, false, "host is error")
	}
	if _, ok := pool.ShardStatus["10.34.60.43:15350"]; !ok {
		assert.Equal(t, true, false, "host is error")
	}
	if _, ok := pool.ShardStatus["10.34.62.46:15350"]; !ok {
		assert.Equal(t, true, false, "host is error")
	}
}

func TestPoolsAnaliyze(t *testing.T) {
	data := "{\"NetworkInterfaceASIO-ShardRegistry\" : {\n\t\t\t\"poolInUse\" : 1,\n\t\t\t\"poolAvailable\" : 3,\n\t\t\t\"poolCreated\" : 614,\n\t\t\t\"poolRefreshing\" : 0,\n\t\t\t\"poolReqQueueLimit\" : 0,\n\t\t\t\"10.34.62.46:15350\" : {\n\t\t\t\t\"inUse\" : 0,\n\t\t\t\t\"available\" : 1,\n\t\t\t\t\"created\" : 66,\n\t\t\t\t\"refreshing\" : 0,\n\t\t\t\t\"reqQueueLimit\" : 0\n\t\t\t},\n\t\t\t\"10.34.60.45:15350\" : {\n\t\t\t\t\"inUse\" : 0,\n\t\t\t\t\"available\" : 1,\n\t\t\t\t\"created\" : 56,\n\t\t\t\t\"refreshing\" : 0,\n\t\t\t\t\"reqQueueLimit\" : 0\n\t\t\t},\n\t\t\t\"10.34.60.43:15350\" : {\n\t\t\t\t\"inUse\" : 1,\n\t\t\t\t\"available\" : 1,\n\t\t\t\t\"created\" : 492,\n\t\t\t\t\"refreshing\" : 0,\n\t\t\t\t\"reqQueueLimit\" : 0\n\t\t\t}\n\t\t},\n\t\t\"global\" : {\n\t\t\t\"poolInUse\" : 0,\n\t\t\t\"poolAvailable\" : 29,\n\t\t\t\"poolCreated\" : 29,\n\t\t\t\"poolRefreshing\" : 0,\n\t\t\t\"poolReqQueueLimit\" : 0,\n\t\t\t\"10.34.62.46:15450\" : {\n\t\t\t\t\"inUse\" : 0,\n\t\t\t\t\"available\" : 3,\n\t\t\t\t\"created\" : 3,\n\t\t\t\t\"refreshing\" : 0,\n\t\t\t\t\"reqQueueLimit\" : 0\n\t\t\t},\n\t\t\t\"10.34.62.46:15350\" : {\n\t\t\t\t\"inUse\" : 0,\n\t\t\t\t\"available\" : 1,\n\t\t\t\t\"created\" : 1,\n\t\t\t\t\"refreshing\" : 0,\n\t\t\t\t\"reqQueueLimit\" : 0\n\t\t\t},\n\t\t\t\"10.34.60.45:15550\" : {\n\t\t\t\t\"inUse\" : 0,\n\t\t\t\t\"available\" : 3,\n\t\t\t\t\"created\" : 3,\n\t\t\t\t\"refreshing\" : 0,\n\t\t\t\t\"reqQueueLimit\" : 0\n\t\t\t},\n\t\t\t\"10.34.60.45:15450\" : {\n\t\t\t\t\"inUse\" : 0,\n\t\t\t\t\"available\" : 1,\n\t\t\t\t\"created\" : 1,\n\t\t\t\t\"refreshing\" : 0,\n\t\t\t\t\"reqQueueLimit\" : 0\n\t\t\t},\n\t\t\t\"10.34.60.43:15350\" : {\n\t\t\t\t\"inUse\" : 0,\n\t\t\t\t\"available\" : 2,\n\t\t\t\t\"created\" : 2,\n\t\t\t\t\"refreshing\" : 0,\n\t\t\t\t\"reqQueueLimit\" : 0\n\t\t\t},\n\t\t\t\"10.34.60.45:15250\" : {\n\t\t\t\t\"inUse\" : 0,\n\t\t\t\t\"available\" : 3,\n\t\t\t\t\"created\" : 3,\n\t\t\t\t\"refreshing\" : 0,\n\t\t\t\t\"reqQueueLimit\" : 0\n\t\t\t},\n\t\t\t\"10.34.60.45:15350\" : {\n\t\t\t\t\"inUse\" : 0,\n\t\t\t\t\"available\" : 1,\n\t\t\t\t\"created\" : 1,\n\t\t\t\t\"refreshing\" : 0,\n\t\t\t\t\"reqQueueLimit\" : 0\n\t\t\t},\n\t\t\t\"10.34.62.46:15250\" : {\n\t\t\t\t\"inUse\" : 0,\n\t\t\t\t\"available\" : 3,\n\t\t\t\t\"created\" : 3,\n\t\t\t\t\"refreshing\" : 0,\n\t\t\t\t\"reqQueueLimit\" : 0\n\t\t\t},\n\t\t\t\"10.34.60.43:15250\" : {\n\t\t\t\t\"inUse\" : 0,\n\t\t\t\t\"available\" : 2,\n\t\t\t\t\"created\" : 2,\n\t\t\t\t\"refreshing\" : 0,\n\t\t\t\t\"reqQueueLimit\" : 0\n\t\t\t},\n\t\t\t\"10.34.60.43:15550\" : {\n\t\t\t\t\"inUse\" : 0,\n\t\t\t\t\"available\" : 2,\n\t\t\t\t\"created\" : 2,\n\t\t\t\t\"refreshing\" : 0,\n\t\t\t\t\"reqQueueLimit\" : 0\n\t\t\t},\n\t\t\t\"10.34.62.46:15550\" : {\n\t\t\t\t\"inUse\" : 0,\n\t\t\t\t\"available\" : 3,\n\t\t\t\t\"created\" : 3,\n\t\t\t\t\"refreshing\" : 0,\n\t\t\t\t\"reqQueueLimit\" : 0\n\t\t\t},\n\t\t\t\"10.34.60.43:15450\" : {\n\t\t\t\t\"inUse\" : 0,\n\t\t\t\t\"available\" : 5,\n\t\t\t\t\"created\" : 5,\n\t\t\t\t\"refreshing\" : 0,\n\t\t\t\t\"reqQueueLimit\" : 0\n\t\t\t}\n\t\t}}"

	res := make(map[string]interface{})
	err := json.Unmarshal([]byte(data), &res)
	assert.Equal(t, nil, err, "json must be correct")

	nameToShards := make(map[string]string)
	nameToShards["10.34.60.45:15350"] = "config"
	nameToShards["10.34.60.43:15350"] = "config"
	nameToShards["10.34.62.46:15350"] = "config"

	nameToShards["10.34.60.45:15550"] = "shard3"
	nameToShards["10.34.60.45:15450"] = "shard2"
	nameToShards["10.34.60.45:15250"] = "shard1"
	nameToShards["10.34.60.43:15550"] = "shard3"
	nameToShards["10.34.60.43:15450"] = "shard2"
	nameToShards["10.34.60.43:15250"] = "shard1"
	nameToShards["10.34.62.46:15550"] = "shard3"
	nameToShards["10.34.62.46:15450"] = "shard2"
	nameToShards["10.34.62.46:15250"] = "shard1"

	pools := PoolsAnaliyze(res, nameToShards)
	assert.Equal(t, 2, len(pools), "pool equal to 2")

	if pool,ok := pools["NetworkInterfaceASIO-ShardRegistry"]; ok {
		assert.NotEqual(t, nil, pool)
		assert.Equal(t, 1, int(pool.PoolInUse))
		assert.Equal(t, 3, int(pool.PoolAvailable))
		assert.Equal(t, 614, int(pool.PoolCreated))
		assert.Equal(t, 0, int(pool.PoolRefreshing))
		assert.Equal(t, 0, int(pool.PoolReqQueueLimit))
		assert.Equal(t, 1, len(pool.ShardStatus))
	} else {
		assert.Equal(t, true, false, "error")
	}

	if pool,ok := pools["global"]; ok {
		assert.NotEqual(t, nil, pool)
		assert.Equal(t, 0, int(pool.PoolInUse))
		assert.Equal(t, 29, int(pool.PoolAvailable))
		assert.Equal(t, 29, int(pool.PoolCreated))
		assert.Equal(t, 0, int(pool.PoolRefreshing))
		assert.Equal(t, 0, int(pool.PoolReqQueueLimit))
		assert.Equal(t, 4, len(pool.ShardStatus))
	} else {
		assert.Equal(t, true, false, "error")
	}


}
