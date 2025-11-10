package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	once sync.Once

	TotalScans = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "sentrinet_total_scans",
		Help: "Total number of scan operations executed",
	})

	OpenPorts = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "sentrinet_open_ports_total",
		Help: "Total number of open ports discovered",
	})

	ClosedPorts = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "sentrinet_closed_ports_total",
		Help: "Total number of closed ports discovered",
	})

	ScanDurationMs = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "sentrinet_scan_duration_ms",
		Help: "Histogram of individual port scan durations(ms)",
		Buckets: prometheus.ExponentialBuckets(1,2,12),
	})

	CleanupDeleted = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "sentrinet_cleanup_deleted_ports_total",
		Help: "Total number of closed ports deleted by cleanup jobs",
	})
)

func Register(reg prometheus.Registerer){
	once.Do(func() {
		reg.MustRegister(TotalScans)
		reg.MustRegister(OpenPorts)
		reg.MustRegister(ClosedPorts)
		reg.MustRegister(ScanDurationMs)
		reg.MustRegister(CleanupDeleted)
	})
}