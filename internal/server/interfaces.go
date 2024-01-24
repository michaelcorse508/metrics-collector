package server

import (
	"context"
	"github.com/bazookajoe1/metrics-collector/internal/pcstats"
	"time"
)

type Storage interface {
	CheckMetricAndSave(ctx context.Context, metric pcstats.Metric) error
	CheckMetricBatchAndSave(ctx context.Context, batch pcstats.Metrics) (err error)
	GetMetric(ctx context.Context, metricID string, metricType pcstats.MetricType) (*pcstats.Metric, error)
	GetAllMetrics(ctx context.Context) (pcstats.Metrics, error)
	Ping(ctx context.Context) error
	RunStorage(ctx context.Context)
}

type Logger interface {
	Info(string)
	Debug(string)
	Error(string)
	Fatal(string)
}

type ServConfigurator interface {
	GetAddress() string
	GetPort() string
	GetSecretKey() []byte
}

type StorageConfigurator interface {
	GetDatabaseDSN() string
	GetRestoreFlag() bool
	GetStoreFilePath() string
	GetStoreInterval() time.Duration
}
