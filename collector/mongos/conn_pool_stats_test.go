//// Copyright 2017 Percona LLC
////
//// Licensed under the Apache License, Version 2.0 (the "License");
//// you may not use this file except in compliance with the License.
//// You may obtain a copy of the License at
////
////   http://www.apache.org/licenses/LICENSE-2.0
////
//// Unless required by applicable law or agreed to in writing, software
//// distributed under the License is distributed on an "AS IS" BASIS,
//// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//// See the License for the specific language governing permissions and
//// limitations under the License.
//
//package mongos
//
//import (
//	"context"
//	"github.com/percona/mongodb_exporter/collector/common"
//	"testing"
//	"time"
//
//	"github.com/stretchr/testify/assert"
//
//	"github.com/percona/mongodb_exporter/testutils"
//)
//
//func TestGetConnPoolStatsDecodesFine(t *testing.T) {
//	// setup
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	t.Run("mongos", func(t *testing.T) {
//		// setup
//		t.Parallel()
//		mongosClient := testutils.MustGetConnectedMongosClient(ctx, t)
//		defer mongosClient.Disconnect(ctx)
//		hostToShardName, err := GetHostToShardNameMap(mongosClient)
//
//		assert.Equal(t, nil, err, "err must be nil")
//		assert.NotEqual(t, len(hostToShardName), 0, "hostToShardName must be great than 0")
//		// run
//		statusMongos := common.GetConnPoolStats(mongosClient, hostToShardName)
//
//		// test
//		assert.NotNil(t, statusMongos)
//		assert.NotEqual(t, len(statusMongos.Pools), 0, "mongos pools len must be not 0")
//	})
//}
