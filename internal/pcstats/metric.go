package pcstats

import (
	"github.com/bazookajoe1/metrics-collector/internal/auxiliary"
	"github.com/pkg/errors"
	"strconv"
)

type MetricType string

func (mT MetricType) String() string {
	return string(mT)
}

const (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
)

type Metric struct {
	ID    string     `json:"id"`
	MType MetricType `json:"type"`
	Delta *int64     `json:"delta,omitempty"`
	Value *float64   `json:"value,omitempty"`
}

// NewMetric creates the new instance of Metric. It validates all given params. Use as counterValue and gaugeValue only
// heap allocated vars to prevent unexpected behavior.
func NewMetric(metricID string, metricType MetricType, counterValue *int64, gaugeValue *float64) (*Metric, error) {
	if err := CheckIDIsNotEmpty(metricID); err != nil {
		return nil, errors.Wrap(err, metricID)
	}

	switch metricType {
	case Gauge:
		return createGaugeMetric(metricID, gaugeValue)
	case Counter:
		return createCounterMetric(metricID, counterValue)
	default:
		return nil, errors.Wrap(ErrBadType, metricType.String())
	}
}

// NewMetricFromString creates the new instance of Metric from given string params. It validates all given params.
func NewMetricFromString(metricID string, metricType string, metricValue string) (*Metric, error) {
	var (
		delta *int64
		value *float64
		mType MetricType
	)

	mType = MetricType(metricType)

	delta, _ = auxiliary.ConvertStringToInt64(metricValue)
	value, _ = auxiliary.ConvertStringToFloat64(metricValue)
	// Above we don't care about errors. Underlying function will check arguments on its own.

	return NewMetric(metricID, mType, delta, value)
}

// GetID returns ID of metric.
func (m *Metric) GetID() string {
	return m.ID
}

// GetType returns type of metric.
func (m *Metric) GetType() MetricType {
	return m.MType
}

// GetGaugeValue returns float64 value of Metric if type is Gauge. If not, returns error.
func (m *Metric) GetGaugeValue() (float64, error) {
	if m.MType == Gauge {
		err := CheckGaugeIsNotNil(m.Value)
		if err != nil {
			return 0, err
		}
		return *m.Value, nil
	}

	return 0, errors.Wrap(ErrInvalidValueSelector, m.MType.String())
}

// GetCounterValue returns int64 value of Metric if type is Counter. If not, returns error.
func (m *Metric) GetCounterValue() (int64, error) {
	if m.MType == Counter {
		err := CheckCounterIsNotNil(m.Delta)
		if err != nil {
			return 0, err
		}
		return *m.Delta, nil
	}

	return 0, errors.Wrap(ErrInvalidValueSelector, m.MType.String())
}

// GetStringValue returns string representation of metric value.
func (m *Metric) GetStringValue() (string, error) {
	if m.GetType() == Counter {
		delta, err := m.GetCounterValue()
		if err != nil {
			return "", err
		}
		return strconv.FormatInt(delta, 10), nil
	}

	value, err := m.GetGaugeValue()
	if err != nil {
		return "", err
	}
	return strconv.FormatFloat(value, 'f', -1, 64), nil
}

func createGaugeMetric(metricName string, gaugeValue *float64) (*Metric, error) {
	if err := CheckGaugeIsNotNil(gaugeValue); err != nil {
		return nil, err
	}

	metric := &Metric{
		ID:    metricName,
		MType: Gauge,
		Value: gaugeValue,
		Delta: nil,
	}

	return metric, nil
}

func createCounterMetric(metricName string, counterValue *int64) (*Metric, error) {
	if err := CheckCounterIsNotNil(counterValue); err != nil {
		return nil, err
	}

	metric := &Metric{
		ID:    metricName,
		MType: Counter,
		Value: nil,
		Delta: counterValue,
	}

	return metric, nil
}
