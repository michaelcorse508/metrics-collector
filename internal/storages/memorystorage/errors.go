package memorystorage

type StorageError string

func (sE StorageError) Error() string {
	return string(sE)
}

const (
	ErrGaugeMetricNotFound   StorageError = "gauge metric not found in storage"
	ErrCounterMetricNotFound StorageError = "counter metric not found in storage"
)
