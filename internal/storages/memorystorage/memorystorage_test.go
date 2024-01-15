package memorystorage

import (
	"context"
	"reflect"
	"testing"

	"github.com/bazookajoe1/metrics-collector/internal/pcstats"
)

func TestMemoryStorage_CheckMetricAndSave(t *testing.T) {
	storage := MemoryStorage{
		gauge:         make(map[string]float64),
		counter:       make(map[string]int64),
		logger:        MockLogger{},
		isToRestore:   false,
		storeInterval: 300,
	}
	type args struct {
		metric pcstats.Metric
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		Err     error
	}{
		{
			name: "Good gauge save",
			args: args{
				pcstats.Metric{
					ID:    "test",
					MType: pcstats.Gauge,
					Delta: nil,
					Value: new(float64),
				},
			},
			wantErr: false,
		},
		{
			name: "Good counter save",
			args: args{
				pcstats.Metric{
					ID:    "test",
					MType: pcstats.Counter,
					Delta: new(int64),
					Value: nil,
				},
			},
			wantErr: false,
		},
		{
			name: "Bad gauge save",
			args: args{
				pcstats.Metric{
					ID:    "test",
					MType: pcstats.Gauge,
					Delta: new(int64),
					Value: nil,
				},
			},
			wantErr: true,
		},
		{
			name: "Bad counter save",
			args: args{
				pcstats.Metric{
					ID:    "test",
					MType: pcstats.Counter,
					Delta: nil,
					Value: new(float64),
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid type save",
			args: args{
				pcstats.Metric{
					ID:    "test",
					MType: "invalid",
					Delta: new(int64),
					Value: new(float64),
				},
			},
			wantErr: true,
		},
		{
			name: "Empty name save 1",
			args: args{
				pcstats.Metric{
					ID:    "",
					MType: pcstats.Gauge,
					Delta: nil,
					Value: new(float64),
				},
			},
			wantErr: true,
		},
		{
			name: "Empty name save 2",
			args: args{
				pcstats.Metric{
					ID:    "",
					MType: pcstats.Counter,
					Delta: new(int64),
					Value: nil,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &storage
			if err := s.CheckMetricAndSave(context.Background(), tt.args.metric); (err != nil) != tt.wantErr {
				t.Errorf("CheckMetricAndSave() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type MockLogger struct{}

func (l MockLogger) Error(string) {}
func (l MockLogger) Info(string)  {}
func (l MockLogger) Debug(string) {}

type MockFileSaver struct{}

func (s *MockFileSaver) Load() ([]byte, error) {
	return nil, nil
}
func (s *MockFileSaver) Save(data []byte) error {
	return nil
}

func TestMemoryStorage_GetMetric(t *testing.T) {
	storage := MemoryStorage{
		gauge:         make(map[string]float64),
		counter:       make(map[string]int64),
		logger:        MockLogger{},
		isToRestore:   false,
		storeInterval: 300,
	}
	type args struct {
		metricID   string
		metricType pcstats.MetricType
	}
	tests := []struct {
		name    string
		args    args
		want    *pcstats.Metric
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &storage
			got, err := s.GetMetric(context.Background(), tt.args.metricID, tt.args.metricType)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMetric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMetric() got = %v, want %v", got, tt.want)
			}
		})
	}
}
