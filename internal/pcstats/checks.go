package pcstats

func CheckIDIsNotEmpty(name string) error {
	if name == "" {
		return ErrBadID
	}
	return nil
}

func CheckTypeIsValid(mType MetricType) error {
	if mType != Gauge && mType != Counter {
		return ErrBadType
	}
	return nil
}

func CheckGaugeIsNotNil(value *float64) error {
	if value == nil {
		return ErrBadGauge
	}
	return nil
}

func CheckCounterIsNotNil(delta *int64) error {
	if delta == nil {
		return ErrBadCounter
	}
	return nil
}

func CheckMetricParamsIsValid(metric Metric) error {
	err := CheckIDIsNotEmpty(metric.GetID())
	if err != nil {
		return err
	}

	switch metric.GetType() {
	case Gauge:
		err = CheckGaugeIsNotNil(metric.Value)
		if err != nil {
			return err
		}
		return nil
	case Counter:
		err = CheckCounterIsNotNil(metric.Delta)
		if err != nil {
			return err
		}
		return nil
	default:
		return ErrBadType
	}
}
