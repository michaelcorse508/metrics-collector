package memorystorage

import "testing"

func TestStorageError_Error(t *testing.T) {
	tests := []struct {
		name string
		sE   StorageError
		want string
	}{
		{
			name: "Err to string 1",
			sE:   ErrGaugeMetricNotFound,
			want: "gauge metric not found in storage",
		},
		{
			name: "Err to string 2",
			sE:   ErrCounterMetricNotFound,
			want: "counter metric not found in storage",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sE.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
