package agent

import (
	"context"
	"github.com/bazookajoe1/metrics-collector/internal/pcstats"
	"time"
)

type Configurator interface {
	GetAddress() string
	GetPort() string
	GetReportInterval() time.Duration
	GetSecretKey() []byte
	GetWorkersLimit() uint64
}

type Logger interface {
	Info(string)
	Debug(string)
	Error(string)
	Fatal(string)
	Errorf(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Debugf(format string, v ...interface{})
}

type Collector interface {
	Run(ctx context.Context)
	GetMetrics() pcstats.Metrics
}
