package pcstats

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMetric(t *testing.T) {
	type args struct {
		metricName   string
		metricType   MetricType
		counterValue *int64
		gaugeValue   *float64
	}
	tests := []struct {
		name    string
		args    args
		want    *Metric
		wantErr bool
	}{
		{
			name: "good gauge",
			args: args{
				metricName:   "test",
				metricType:   Gauge,
				counterValue: nil,
				gaugeValue:   new(float64),
			},
			want: &Metric{
				ID:    "test",
				MType: Gauge,
				Delta: nil,
				Value: new(float64),
			},
			wantErr: false,
		},
		{
			name: "good counter",
			args: args{
				metricName:   "test",
				metricType:   Counter,
				counterValue: new(int64),
				gaugeValue:   nil,
			},
			want: &Metric{
				ID:    "test",
				MType: Counter,
				Delta: new(int64),
				Value: nil,
			},
			wantErr: false,
		},
		{
			name: "empty name metric",
			args: args{
				metricName:   "",
				metricType:   Gauge,
				counterValue: nil,
				gaugeValue:   new(float64),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "bad gauge",
			args: args{
				metricName:   "test",
				metricType:   Gauge,
				counterValue: new(int64),
				gaugeValue:   nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "bad counter",
			args: args{
				metricName:   "test",
				metricType:   Counter,
				counterValue: nil,
				gaugeValue:   new(float64),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "bad type",
			args: args{
				metricName:   "test",
				metricType:   "random",
				counterValue: new(int64),
				gaugeValue:   new(float64),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMetric(tt.args.metricName, tt.args.metricType, tt.args.counterValue, tt.args.gaugeValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMetric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			if !assert.Equal(t, got.ID, tt.want.ID) {
				t.Errorf("NewMetric() ID: got = %v, want %v", got, tt.want)
			}

			if !assert.Equal(t, got.MType, tt.want.MType) {
				t.Errorf("NewMetric() MType: got = %v, want %v", got, tt.want)
			}

			if !assert.Equal(t, got.Value, tt.want.Value) {
				t.Errorf("NewMetric() Value: got = %v, want %v", got, tt.want)
			}

			if !assert.Equal(t, got.Delta, tt.want.Delta) {
				t.Errorf("NewMetric() Delta: got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMetricFromString(t *testing.T) {
	type args struct {
		metricName  string
		metricType  string
		metricValue string
	}
	type Test struct {
		name    string
		args    args
		want    *Metric
		wantErr bool
	}

	tests := func() []Test {
		tests := make([]Test, 0)

		// good_gauge
		test := Test{
			name: "Good gauge",
			args: args{
				metricName:  "test",
				metricType:  "gauge",
				metricValue: "1.011",
			},
			want: &Metric{
				ID:    "test",
				MType: Gauge,
				Delta: nil,
				Value: new(float64),
			},
			wantErr: false,
		}
		*test.want.Value = 1.011
		tests = append(tests, test)

		// good_counter
		test = Test{
			name: "Good counter",
			args: args{
				metricName:  "test",
				metricType:  "counter",
				metricValue: "10",
			},
			want: &Metric{
				ID:    "test",
				MType: Counter,
				Delta: new(int64),
				Value: nil,
			},
			wantErr: false,
		}
		*test.want.Delta = 10
		tests = append(tests, test)

		// empty_name
		test = Test{
			name: "Empty name",
			args: args{
				metricName:  "",
				metricType:  "counter",
				metricValue: "10",
			},
			want:    nil,
			wantErr: true,
		}
		tests = append(tests, test)

		// bad_gauge
		test = Test{
			name: "Bad gauge",
			args: args{
				metricName:  "test",
				metricType:  "gauge",
				metricValue: "random",
			},
			want:    nil,
			wantErr: true,
		}
		tests = append(tests, test)

		// bad_counter
		test = Test{
			name: "Bad counter",
			args: args{
				metricName:  "test",
				metricType:  "counter",
				metricValue: "10.111",
			},
			want:    nil,
			wantErr: true,
		}
		tests = append(tests, test)

		// bad_type
		test = Test{
			name: "Bad type",
			args: args{
				metricName:  "test",
				metricType:  "random_type",
				metricValue: "10",
			},
			want:    nil,
			wantErr: true,
		}
		tests = append(tests, test)
		return tests
	}()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMetricFromString(tt.args.metricName, tt.args.metricType, tt.args.metricValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMetric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			if !assert.Equal(t, got.ID, tt.want.ID) {
				t.Errorf("NewMetric() ID: got = %v, want %v", got, tt.want)
			}

			if !assert.Equal(t, got.MType, tt.want.MType) {
				t.Errorf("NewMetric() MType: got = %v, want %v", got, tt.want)
			}

			if !assert.Equal(t, got.Value, tt.want.Value) {
				t.Errorf("NewMetric() Value: got = %v, want %v", got, tt.want)
			}

			if !assert.Equal(t, got.Delta, tt.want.Delta) {
				t.Errorf("NewMetric() Delta: got = %v, want %v", got, tt.want)
			}
		})
	}
}
