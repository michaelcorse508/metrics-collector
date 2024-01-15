package postgresdb

import (
	"database/sql"
	"github.com/bazookajoe1/metrics-collector/internal/pcstats"
	"github.com/jackc/pgx/v5"
)

type SelectedMetricFields struct {
	gotID    string
	gotType  string
	gotValue *float64
	gotDelta *int64
}

func (qp *SelectedMetricFields) UnmarshallRow(row pgx.Row) error {
	var (
		gotID    string
		gotType  string
		gotValue sql.NullFloat64
		gotDelta sql.NullInt64
	)

	err := row.Scan(&gotID, &gotType, &gotDelta, &gotValue)
	if err != nil {
		return err
	}

	qp.gotID = gotID
	qp.gotType = gotType
	qp.checkAndSetValue(gotValue)
	qp.checkAndSetDelta(gotDelta)

	return nil
}

func (qp *SelectedMetricFields) checkAndSetValue(value sql.NullFloat64) {
	if value.Valid {
		writtenValue := new(float64)
		*writtenValue = value.Float64
		qp.gotValue = writtenValue
	}
}

func (qp *SelectedMetricFields) checkAndSetDelta(value sql.NullInt64) {
	if value.Valid {
		writtenValue := new(int64)
		*writtenValue = value.Int64
		qp.gotDelta = writtenValue
	}
}

func UnmarshallRows(rows pgx.Rows) []SelectedMetricFields {
	outFields := make([]SelectedMetricFields, 0)

	for rows.Next() {
		selectedMetricFields := new(SelectedMetricFields)
		err := selectedMetricFields.UnmarshallRow(rows)
		if err != nil {
			continue
		}
		outFields = append(outFields, *selectedMetricFields)
	}

	return outFields
}

func SelectedMetricFieldsToMetrics(allFields []SelectedMetricFields) pcstats.Metrics {
	metrics := make(pcstats.Metrics, len(allFields))

	for _, fields := range allFields {
		metric, err := pcstats.NewMetric(fields.gotID, pcstats.MetricType(fields.gotType), fields.gotDelta, fields.gotValue)
		if err != nil {
			continue
		}

		metrics = append(metrics, *metric)
	}

	return metrics
}
