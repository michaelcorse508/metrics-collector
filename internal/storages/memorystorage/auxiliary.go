package memorystorage

import (
	"github.com/bazookajoe1/metrics-collector/internal/pcstats"
	"github.com/pkg/errors"
)

func CheckMetricBatchIsValid(metrics pcstats.Metrics) error {
	for _, metric := range metrics {
		err := pcstats.CheckMetricParamsIsValid(metric)
		if err != nil {
			return errors.Wrap(err, "batch invalid")
		}
	}
	return nil
}

func JSONDataToMetrics(data []byte) (pcstats.Metrics, error) {
	metrics := new(pcstats.Metrics)
	err := metrics.UnmarshalJSON(data)
	if err != nil {
		return nil, err
	}
	return *metrics, nil
}
