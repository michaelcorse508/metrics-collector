package collector

import "github.com/bazookajoe1/metrics-collector/internal/pcstats"

// NeededMetrics are prototypes of metrics we need to collect.
var NeededMetrics = pcstats.Metrics{
	{ID: "Alloc", MType: pcstats.Gauge, Delta: nil, Value: new(float64)},
	{ID: "BuckHashSys", MType: pcstats.Gauge, Delta: nil, Value: new(float64)},
	{ID: "Frees", MType: pcstats.Gauge, Delta: nil, Value: new(float64)},
	{ID: "GCCPUFraction", MType: pcstats.Gauge, Delta: nil, Value: new(float64)},
	{ID: "GCSys", MType: pcstats.Gauge, Delta: nil, Value: new(float64)},
	{ID: "HeapAlloc", MType: pcstats.Gauge, Delta: nil, Value: new(float64)},
	{ID: "HeapIdle", MType: pcstats.Gauge, Delta: nil, Value: new(float64)},
	{ID: "HeapInuse", MType: pcstats.Gauge, Delta: nil, Value: new(float64)},
	{ID: "HeapObjects", MType: pcstats.Gauge, Delta: nil, Value: new(float64)},
	{ID: "HeapReleased", MType: pcstats.Gauge, Delta: nil, Value: new(float64)},
	{ID: "HeapSys", MType: pcstats.Gauge, Delta: nil, Value: new(float64)},
	{ID: "LastGC", MType: pcstats.Gauge, Delta: nil, Value: new(float64)},
	{ID: "Lookups", MType: pcstats.Gauge, Delta: nil, Value: new(float64)},
	{ID: "MCacheInuse", MType: pcstats.Gauge, Delta: nil, Value: new(float64)},
	{ID: "MCacheSys", MType: pcstats.Gauge, Delta: nil, Value: new(float64)},
	{ID: "MSpanInuse", MType: pcstats.Gauge, Delta: nil, Value: new(float64)},
	{ID: "MSpanSys", MType: pcstats.Gauge, Delta: nil, Value: new(float64)},
	{ID: "Mallocs", MType: pcstats.Gauge, Delta: nil, Value: new(float64)},
	{ID: "NextGC", MType: pcstats.Gauge, Delta: nil, Value: new(float64)},
	{ID: "NumForcedGC", MType: pcstats.Gauge, Delta: nil, Value: new(float64)},
	{ID: "NumGC", MType: pcstats.Gauge, Delta: nil, Value: new(float64)},
	{ID: "OtherSys", MType: pcstats.Gauge, Delta: nil, Value: new(float64)},
	{ID: "PauseTotalNs", MType: pcstats.Gauge, Delta: nil, Value: new(float64)},
	{ID: "StackInuse", MType: pcstats.Gauge, Delta: nil, Value: new(float64)},
	{ID: "StackSys", MType: pcstats.Gauge, Delta: nil, Value: new(float64)},
	{ID: "Sys", MType: pcstats.Gauge, Delta: nil, Value: new(float64)},
	{ID: "TotalAlloc", MType: pcstats.Gauge, Delta: nil, Value: new(float64)},
	{ID: "RandomValue", MType: pcstats.Gauge, Delta: nil, Value: new(float64)},
	{ID: "PollCount", MType: pcstats.Counter, Delta: new(int64), Value: nil},
}
