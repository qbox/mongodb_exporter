package mongos

import (
	"context"
	"github.com/percona/mongodb_exporter/collector/util"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/percona/mongodb_exporter/testutils"
)

func TestGetShardConnPoolStatsDecodesFine(t *testing.T) {
	// setup
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t.Run("mongod", func(t *testing.T) {
		// setup
		t.Parallel()
		defaultClient := testutils.MustGetConnectedMongodClient(ctx, t)
		defer defaultClient.Disconnect(ctx)

		// run
		statusDefault := GetShardConnPoolStats(defaultClient, nil)

		// test
		assert.NotNil(t, statusDefault)
		assert.Equal(t, int(statusDefault.TotalInUse), 0, "mongod totalInUse must be 0")
		assert.Equal(t, int(statusDefault.TotalCreated), 0, "mongod totalCreated must be 0")
		assert.Equal(t, int(statusDefault.TotalAvailable), 0, "mongod totalAvailable must be 0")
		assert.Equal(t, statusDefault.GlobalPool, (*util.PoolConnStats)(nil), "mongod globalPool must be nil")
	})

	t.Run("replset", func(t *testing.T) {
		// setup
		t.Parallel()
		replSetClient := testutils.MustGetConnectedReplSetClient(ctx, t)
		defer replSetClient.Disconnect(ctx)

		// run
		statusReplSet := GetShardConnPoolStats(replSetClient, nil)

		// test
		assert.NotNil(t, statusReplSet)
		// test
		assert.NotNil(t, statusReplSet)
		assert.Equal(t, int(statusReplSet.TotalInUse), 0, "mongod totalInUse must be 0")
		assert.Equal(t, int(statusReplSet.TotalCreated), 0, "mongod totalCreated must be 0")
		assert.Equal(t, int(statusReplSet.TotalAvailable), 0, "mongod totalAvailable must be 0")
		assert.Equal(t, statusReplSet.GlobalPool, (*util.PoolConnStats)(nil), "mongod globalPool must be nil")
	})


	t.Run("mongos", func(t *testing.T) {
		// setup
		t.Parallel()
		mongosClient := testutils.MustGetConnectedMongosClient(ctx, t)
		defer mongosClient.Disconnect(ctx)
		hostToShardName, err := GetHostToShardNameMap(mongosClient)

		assert.Equal(t, nil, err, "err must be nil")
		assert.NotEqual(t, len(hostToShardName), 0, "hostToShardName must be great than 0")
		// run
		statusMongos := GetShardConnPoolStats(mongosClient, hostToShardName)

		// test
		assert.NotNil(t, statusMongos)
		assert.NotEqual(t, statusMongos.GlobalPool, (*util.PoolConnStats)(nil), "mongod globalPool must be not nil")
	})

}
