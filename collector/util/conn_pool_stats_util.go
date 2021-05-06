package util

import (
	"github.com/prometheus/common/log"
	"reflect"
	"strings"
)

type ReplicaSetStats struct {
	RefreshLimiter float64 `bson:"refreshLimiter,omitempty"`
}

type HostConnStats struct {
	InUse         float64
	Available     float64
	Created       float64
	Refreshing    float64
	ReqQueueLimit float64
}

func (this *HostConnStats) Add(lhs *HostConnStats) {
	if lhs == nil {
		return
	}
	this.InUse += lhs.InUse
	this.Available += lhs.Available
	this.Created += lhs.Created
	this.Refreshing += lhs.Refreshing
	this.ReqQueueLimit += lhs.ReqQueueLimit
}

func getHostConnStatsFromInterface(obj map[string]interface{}) *HostConnStats {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("getHostConnStatsFromInterface panic, err:%v", r)
		}
	}()

	if obj == nil {
		return nil
	}
	tmp := HostConnStats{}

	for key, value := range obj {
		if len(key) == 0 {
			continue
		}

		if key == "inUse" {
			tmp.InUse = valueType(value)
		} else if key == "available" {
			tmp.Available = valueType(value)
		} else if key == "created" {
			tmp.Created = valueType(value)
		} else if key == "refreshing" {
			tmp.Refreshing = valueType(value)
		} else if key == "reqQueueLimit" {
			tmp.ReqQueueLimit = valueType(value)
		} else {
			log.Errorf("invalid hostStat key:%s, value:%+v", key, value)
		}
	}
	return &tmp
}

type PoolConnStats struct {
	PoolInUse         float64                   `json:"poolInUse"`
	PoolAvailable     float64                   `json:"poolAvailable"`
	PoolCreated       float64                   `json:"poolCreated"`
	PoolRefreshing    float64                   `json:"poolRefreshing"`
	PoolReqQueueLimit float64                   `json:"poolReqQueueLimit,omitempty"`
	ShardStatus       map[string]*HostConnStats `json:"shardStatus,omitempty"`
}

func GetPoolConnStats(obj map[string]interface{}) *PoolConnStats {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("GetPoolConnStats panic, err:%v", r)
		}
	}()

	if obj == nil || len(obj) == 0 {
		return nil
	}
	tmp := PoolConnStats{ShardStatus: make(map[string]*HostConnStats)}

	for key, value := range obj {
		if len(key) == 0 {
			continue
		}

		if key == "poolReqQueueLimit" {
			tmp.PoolReqQueueLimit = valueType(value)
		} else if key == "poolRefreshing" {
			tmp.PoolRefreshing = valueType(value)
		} else if key == "poolCreated" {
			tmp.PoolCreated = valueType(value)
		} else if key == "poolAvailable" {
			tmp.PoolAvailable = valueType(value)
		} else if key == "poolInUse" {
			tmp.PoolInUse = valueType(value)
		} else if strings.Index(key, ":") != -1 {
			host := getHostConnStatsFromInterface(value.(map[string]interface{}))
			if host == nil {
				continue
			}
			tmp.ShardStatus[key] = host
		} else {
			log.Errorf("invalid key:%s, value:%+v", key, value)
		}
	}
	return &tmp
}

func valueType(value interface{}) float64 {
	if value == nil {
		return 0
	}

	switch value.(type) {
	case float64:
		return value.(float64)
	case float32:
		return float64(value.(float32))
	case int64:
		return float64(value.(int64))
	case int32:
		return float64(value.(int32))
	case int:
		return float64(value.(int))
	case uint64:
		return float64(value.(uint64))
	case uint32:
		return float64(value.(uint32))
	case uint:
		return float64(value.(uint))
	default:
		floatType := reflect.TypeOf(float64(0))
		log.Errorf("%+v", reflect.TypeOf(value))
		v := reflect.ValueOf(value)
		v = reflect.Indirect(v)
		if v.Type().ConvertibleTo(floatType) {
			fv := v.Convert(floatType)
			return fv.Float()
		} else {
			log.Errorf("interface to float is error, interface:%+v, type:%v", value, reflect.TypeOf(value))
			return 0
		}
	}
}

func PoolsAnaliyze(pools map[string]interface{}, hostToShardName map[string]string) map[string]*PoolConnStats {
	if len(pools) == 0 {
		return nil
	}

	poolsStat := make(map[string]*PoolConnStats)
	for poolName, poolStat := range pools {
		if len(poolName) == 0 || poolStat == nil {
			continue
		}
		switch poolStat.(type) {
		case map[string]interface{}:
			{
				tmpPoolStats := GetPoolConnStats(poolStat.(map[string]interface{}))
				shardToHostStats := make(map[string]*HostConnStats)
				if tmpPoolStats != nil {
					if len(tmpPoolStats.ShardStatus) == 0 {
						break
					}

					for hostKey, hostStats := range tmpPoolStats.ShardStatus {
						if hostStats == nil || len(hostKey) == 0 {
							continue
						}
						if name, ok := hostToShardName[hostKey]; ok {
							if len(name) == 0 {
								continue
							}
							if tmpH, isOk := shardToHostStats[name]; isOk {
								tmpH.Add(hostStats)
							} else {
								shardToHostStats[name] = hostStats
							}
						} else {
							shardToHostStats["other"] = hostStats
						}
					}
					tmpPoolStats.ShardStatus = shardToHostStats
				}
				poolsStat[poolName] = tmpPoolStats
			}
		default:
			log.Errorf("%s pool type is not map[string]interface{}, type:%+v", poolName, reflect.TypeOf(poolStat))
		}
	}
	return poolsStat
}
