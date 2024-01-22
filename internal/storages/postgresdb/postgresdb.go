package postgresdb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bazookajoe1/metrics-collector/internal/pcstats"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	errorsWithWrap "github.com/pkg/errors"
)

// "user=server password=lD093hf09L20 host=192.168.163.130 port=5432 dbname=metrics sslmode=disable"

// PostgresDB describes PostgreSQL DB storage.
type PostgresDB struct {
	connectionString string
	db               *pgxpool.Pool
	logger           Logger
}

// NewPostgresDB creates the instance of PostgresDB struct. It opens pgxpool.Pool, pings it and creates needed tables.
func NewPostgresDB(
	ctx context.Context,
	connectionString string,
	logger Logger,
) (*PostgresDB, error) {
	pdb := &PostgresDB{
		connectionString: connectionString,
		logger:           logger,
	}

	db, err := createDBWithTablesAndCheckConnection(ctx, pdb.connectionString)
	if err != nil {
		return nil, err
	}

	pdb.db = db
	return pdb, nil
}

// CheckMetricAndSave performs validity checking of pcstats.Metric and saves it to storage.
func (s *PostgresDB) CheckMetricAndSave(ctx context.Context, metric pcstats.Metric) error {
	err := pcstats.CheckMetricParamsIsValid(metric)
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}

	err = s.saveMetricToDB(ctx, metric)
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}

	return nil
}

func (s *PostgresDB) CheckMetricBatchAndSave(
	ctx context.Context,
	batch pcstats.Metrics,
) (err error) {
	txCtx, cancelTx := context.WithTimeout(ctx, 10*time.Second)
	defer cancelTx()

	tx, err := s.db.Begin(txCtx)
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}

	defer func() {
		if err != nil {
			errRb := tx.Rollback(txCtx)
			if errRb != nil {
				err = fmt.Errorf("%s: cannot rollback", errRb.Error())
				return
			}
			return
		}
		err = tx.Commit(txCtx)
	}()

	return func() error { // here we in cycle construct transaction, if error in any query, transaction will be rolled back
		for _, metric := range batch {
			name := metric.GetID()
			mType := metric.GetType()
			value, _ := metric.GetGaugeValue() // here I don't care about error
			delta, _ := metric.GetCounterValue()

			result, err := tx.Exec(txCtx, InsertUpdateRowQuery, name, mType, value, delta)
			if err != nil {
				return fmt.Errorf("%s: upsert failure", err.Error())
			}
			rowsAffected := result.RowsAffected()
			if rowsAffected == 0 {
				s.logger.Debug(fmt.Sprintf("UPSERT rows affected on metric: %v is 0", metric))
			}
		}
		return nil
	}()
}

// GetMetric selects pcstats.Metric from storage by given params.
func (s *PostgresDB) GetMetric(
	ctx context.Context,
	metricID string,
	metricType pcstats.MetricType,
) (*pcstats.Metric, error) {
	metricFields, err := s.selectRowByParams(ctx, metricID, metricType)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	fmt.Println("===========>", metricFields.gotID)

	metric, err := pcstats.NewMetric(
		metricFields.gotID,
		pcstats.MetricType(metricFields.gotType),
		metricFields.gotDelta,
		metricFields.gotValue,
	)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return metric, nil
}

// GetAllMetrics selects all metrics from storage and returns them as pcstats.Metrics.
func (s *PostgresDB) GetAllMetrics(ctx context.Context) (pcstats.Metrics, error) {
	allFields, err := s.selectAllRows(ctx)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, errorsWithWrap.Wrap(err, "cannot get all metrics from db")
	}

	metrics := SelectedMetricFieldsToMetrics(allFields)
	return metrics, nil
}

// RunStorage spies for context cancelling and aborts DB connection if context cancelled.
func (s *PostgresDB) RunStorage(ctx context.Context) {
	<-ctx.Done()
	s.stop()
	s.logger.Info("got context cancelling; database connection aborted")
}

func (s *PostgresDB) Ping(ctx context.Context) error {
	return checkConnectionWithDB(ctx, s.db)
}

func (s *PostgresDB) stop() {
	s.db.Close()
}

func (s *PostgresDB) saveMetricToDB(ctx context.Context, metric pcstats.Metric) error {
	name := metric.GetID()
	mType := metric.GetType()
	value, _ := metric.GetGaugeValue() // here I don't care about error
	delta, _ := metric.GetCounterValue()

	_, err := s.retriableExec(ctx, InsertUpdateRowQuery, name, mType, value, delta)
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}

	return nil
}

func (s *PostgresDB) selectRowByParams(
	ctx context.Context,
	metricID string,
	metricType pcstats.MetricType,
) (SelectedMetricFields, error) {

	execCtx, cancelExec := context.WithTimeout(ctx, 5*time.Second)
	defer cancelExec()

	selectedMetricValues := s.retriableQueryRow(
		execCtx,
		SelectRowQuery,
		metricID,
		metricType.String(),
	)

	return selectedMetricValues, nil
}

func (s *PostgresDB) selectAllRows(ctx context.Context) ([]SelectedMetricFields, error) {
	queryCtx, cancelQuery := context.WithTimeout(ctx, 5*time.Second)
	defer cancelQuery()

	rows, err := s.retriableQuery(queryCtx, SelectAllRowsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	outFields := UnmarshallRows(rows)

	return outFields, nil
}

// createDBWithTablesAndCheckConnection creates new pgxpool.Pool, pings it and creates tables if DB is available.
func createDBWithTablesAndCheckConnection(
	ctx context.Context,
	connectionString string,
) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		return nil, err
	}

	if err = checkConnectionWithDB(ctx, db); err != nil {
		return nil, err
	}

	err = createTables(ctx, db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// createTables sends to DB query that creates needed table in DB.
func createTables(ctx context.Context, db *pgxpool.Pool) error {
	var pgErr *pgconn.PgError

	crtCtx, cancelCrt := context.WithTimeout(ctx, 5*time.Second)
	defer cancelCrt()

	_, err := db.Exec(crtCtx, CreateTableQuery)
	if err != nil && errors.As(err, &pgErr) {
		if pgErr.Code != pgerrcode.DuplicateTable { // if this error - do nothing (relation exists)
			return err
		}
	}

	return nil
}

func checkConnectionWithDB(ctx context.Context, db *pgxpool.Pool) error {
	pingCtx, cancelPing := context.WithTimeout(ctx, 5*time.Second)
	defer cancelPing()

	return db.Ping(pingCtx)
}
