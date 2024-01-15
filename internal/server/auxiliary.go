package server

import (
	"context"
	"fmt"
	"github.com/bazookajoe1/metrics-collector/internal/pcstats"
	"github.com/bazookajoe1/metrics-collector/internal/storages/memorystorage"
	"github.com/bazookajoe1/metrics-collector/internal/storages/postgresdb"
	"sort"
)

// MetricSliceToHTMLString performs conversion of Metric slice to formatted string,
// like in MetricSliceToString, but with <br> instead of "\n".
// Each string representation of Metric is separated from another by '\n'.
// Metrics in output string is sorted by type and name.
func MetricSliceToHTMLString(metrics pcstats.Metrics) string {
	outString := ""

	sort.Sort(metrics)

	for _, metric := range metrics {
		value, err := metric.GetStringValue()
		if err != nil {
			continue
		}
		outString += fmt.Sprintf("%s: %s", metric.GetID(), value)
		outString += "<br>"
	}

	return outString
}

func MetricSliceToJSON(metrics pcstats.Metrics) []byte {
	var outData = make([]byte, 0)

	sort.Sort(metrics)

	for _, metric := range metrics {
		data, err := metric.MarshalJSON()
		if err != nil {
			continue
		}
		outData = append(outData, data...)
		outData = append(outData, '\n')
	}

	return outData
}

func ChooseStorage(ctx context.Context, config StorageConfigurator, logger Logger) Storage {
	var (
		storage Storage
		err     error
	)

	storage, err = postgresdb.NewPostgresDB(ctx, config.GetDatabaseDSN(), logger)
	if err != nil {
		storage = memorystorage.NewMemoryStorage(config, logger)
		logger.Info("memory storage chosen")
	} else {
		logger.Info("database storage chosen")
	}

	return storage
}
