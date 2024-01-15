package postgresdb

import (
	"context"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"reflect"
	"testing"
)

func TestPostgresDB_retriableExec(t *testing.T) {
	type fields struct {
		connectionString string
		db               *pgxpool.Pool
		logger           Logger
	}
	type args struct {
		ctx       context.Context
		sql       string
		arguments []any
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantCommandTag pgconn.CommandTag
		wantErr        bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PostgresDB{
				connectionString: tt.fields.connectionString,
				db:               tt.fields.db,
				logger:           tt.fields.logger,
			}
			gotCommandTag, err := s.retriableExec(tt.args.ctx, tt.args.sql, tt.args.arguments...)
			if (err != nil) != tt.wantErr {
				t.Errorf("retriableExec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCommandTag, tt.wantCommandTag) {
				t.Errorf("retriableExec() gotCommandTag = %v, want %v", gotCommandTag, tt.wantCommandTag)
			}
		})
	}
}
