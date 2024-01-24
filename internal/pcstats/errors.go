package pcstats

type MetricError string

func (mE MetricError) Error() string {
	return string(mE)
}

const (
	ErrBadID                MetricError = "invalid metric ID"
	ErrBadType              MetricError = "invalid metric type"
	ErrBadGauge             MetricError = "gauge value is nil"
	ErrBadCounter           MetricError = "counter value is nil"
	ErrInvalidValueSelector MetricError = "inappropriate value getter"
)
