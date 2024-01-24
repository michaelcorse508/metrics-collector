package pcstats

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMetricError_Error(t *testing.T) {
	tests := []struct {
		name string
		mE   MetricError
		want string
	}{
		{
			name: "Error bad type",
			mE:   ErrBadType,
			want: "invalid metric type",
		},
		{
			name: "Error bad ID",
			mE:   ErrBadID,
			want: "invalid metric ID",
		},
		{
			name: "Error bad gauge",
			mE:   ErrBadGauge,
			want: "gauge value is nil",
		},
		{
			name: "Error bad counter",
			mE:   ErrBadCounter,
			want: "counter value is nil",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.mE.Error(), "Error()")
		})
	}
}
