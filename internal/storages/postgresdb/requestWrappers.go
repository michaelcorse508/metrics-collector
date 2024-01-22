package postgresdb

import (
	"context"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// retriableExec is the wrapper that sets retry count and initial wait time between requests for pgxpool.Pool Exec.
func (s *PostgresDB) retriableExec(
	ctx context.Context,
	sql string,
	arguments ...any,
) (commandTag pgconn.CommandTag, err error) {
	retryCount := 3
	waitTime := 1 * time.Second

	return s.execWrapper(ctx, retryCount, waitTime, sql, arguments...)
}

func (s *PostgresDB) execWrapper(
	ctx context.Context,
	retryCount int,
	waitTime time.Duration,
	sql string,
	arguments ...any,
) (commandTag pgconn.CommandTag, err error) {
	timeIncrement := 2 * time.Second

	execCtx, cancelExec := context.WithTimeout(ctx, 5*time.Second)
	defer cancelExec()

	commandTag, err = s.db.Exec(execCtx, sql, arguments...)
	if err != nil && pgerrcode.IsConnectionException(err.Error()) && retryCount > 0 {
		retryCount--
		time.Sleep(waitTime)
		waitTime += timeIncrement

		commandTag, err = s.execWrapper(ctx, retryCount, waitTime, sql, arguments...)
	}
	return
}

// retriableQuery is the wrapper that sets retry count and initial wait time between requests for pgxpool.Pool Query.
func (s *PostgresDB) retriableQuery(
	ctx context.Context,
	sql string,
	args ...any,
) (pgx.Rows, error) {
	retryCount := 3
	waitTime := 1 * time.Second

	return s.queryWrapper(ctx, retryCount, waitTime, sql, args...)
}

func (s *PostgresDB) queryWrapper(
	ctx context.Context,
	retryCount int,
	waitTime time.Duration,
	sql string,
	args ...any,
) (rows pgx.Rows, err error) {
	timeIncrement := 2 * time.Second

	execCtx, cancelExec := context.WithTimeout(ctx, 5*time.Second)
	defer cancelExec()

	rows, err = s.db.Query(execCtx, sql, args...)
	if err != nil && pgerrcode.IsConnectionException(err.Error()) && retryCount > 0 {
		retryCount--
		time.Sleep(waitTime)
		waitTime += timeIncrement

		rows, err = s.queryWrapper(ctx, retryCount, waitTime, sql, args...)
	}
	return
}

// retriableQuery is the wrapper that sets retry count and initial wait time between requests for pgxpool.Pool Query.
func (s *PostgresDB) retriableQueryRow(
	ctx context.Context,
	sql string,
	args ...any,
) SelectedMetricFields {
	retryCount := 3
	waitTime := 1 * time.Second

	return s.queryRowWrapper(ctx, retryCount, waitTime, sql, args...)
}

func (s *PostgresDB) queryRowWrapper(
	ctx context.Context,
	retryCount int,
	waitTime time.Duration,
	sql string,
	args ...any,
) (metricsValues SelectedMetricFields) {

	timeIncrement := 2 * time.Second

	execCtx, cancelExec := context.WithTimeout(ctx, 5*time.Second)
	defer cancelExec()

	row := s.db.QueryRow(execCtx, sql, args...)

	err := metricsValues.UnmarshallRow(row)
	if err != nil && pgerrcode.IsConnectionException(err.Error()) && retryCount > 0 {
		retryCount--
		time.Sleep(waitTime)
		waitTime += timeIncrement

		metricsValues = s.queryRowWrapper(ctx, retryCount, waitTime, sql, args...)
	}
	return
}
