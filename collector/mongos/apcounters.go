package mongos

import "github.com/prometheus/client_golang/prometheus"

var (
	apCountersTotalDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "", "ap_counters_total"),
		"The apcounters data structure",
		[]string{"type"},
		nil,
	)
)

// apcounter for new add monitor
type ApCounters struct {
	ReadAp              float64 `bson:"readAp"`
	ReadTp              float64 `bson:"readTp"`
	ErrorApExecutorPool float64 `bson:"error_apexecutor_pool"`

	ReadSlowLog    float64 `bson:"read_slowlog"`
	ReadApSlowLog  float64 `bson:"read_ap_slowlog"`
	ReadDSlowLog   float64 `bson:"read_d_slowlog"`
	ReadApDSlowLog float64 `bson:"read_ap_d_slowlog"`
	ReadUnSlowLog  float64 `bson:"read_un_slowlog"`
	WriteSlowLog   float64 `bson:"write_slowlog"`
	FamSlowLog     float64 `bson:"fam_slowlog"`
	CmdSlowLog     float64 `bson:"cmd_slowlog"`

	LimitForLegacy   float64 `bson:"limitForLegacy"`
	LimitForAsioReqQ float64 `bson:"limitForAsioReqQ"`
	LimitForRefresh  float64 `bson:"limitForRefresh"`
}

func (this *ApCounters) Export(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(apCountersTotalDesc, prometheus.CounterValue, this.ReadAp, "readap")
	ch <- prometheus.MustNewConstMetric(apCountersTotalDesc, prometheus.CounterValue, this.ReadTp, "readtp")
	ch <- prometheus.MustNewConstMetric(apCountersTotalDesc, prometheus.CounterValue, this.ErrorApExecutorPool, "error_ap_exec")
	ch <- prometheus.MustNewConstMetric(apCountersTotalDesc, prometheus.CounterValue, this.ReadSlowLog, "r_s_l")
	ch <- prometheus.MustNewConstMetric(apCountersTotalDesc, prometheus.CounterValue, this.ReadApSlowLog, "r_a_s_l")
	ch <- prometheus.MustNewConstMetric(apCountersTotalDesc, prometheus.CounterValue, this.ReadDSlowLog, "r_d_s_l")
	ch <- prometheus.MustNewConstMetric(apCountersTotalDesc, prometheus.CounterValue, this.ReadApDSlowLog, "r_a_d_s_l")
	ch <- prometheus.MustNewConstMetric(apCountersTotalDesc, prometheus.CounterValue, this.ReadUnSlowLog, "r_un_s_l")
	ch <- prometheus.MustNewConstMetric(apCountersTotalDesc, prometheus.CounterValue, this.WriteSlowLog, "w_s_l")
	ch <- prometheus.MustNewConstMetric(apCountersTotalDesc, prometheus.CounterValue, this.CmdSlowLog, "c_s_l")
	ch <- prometheus.MustNewConstMetric(apCountersTotalDesc, prometheus.CounterValue, this.FamSlowLog, "f_s_l")
	ch <- prometheus.MustNewConstMetric(apCountersTotalDesc, prometheus.CounterValue, this.LimitForLegacy, "limitForLegacy")
	ch <- prometheus.MustNewConstMetric(apCountersTotalDesc, prometheus.CounterValue, this.LimitForAsioReqQ, "limitForAsioReqQ")
	ch <- prometheus.MustNewConstMetric(apCountersTotalDesc, prometheus.CounterValue, this.LimitForRefresh, "limitForRefresh")
}

func (this *ApCounters) Describe(ch chan<- *prometheus.Desc) {
	ch <- apCountersTotalDesc
}
