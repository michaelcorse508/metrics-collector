package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/bazookajoe1/metrics-collector/internal/pcstats"
	"github.com/bazookajoe1/metrics-collector/internal/storages/memorystorage/jsonfilesaver"
	"github.com/pkg/errors"
)

// MemoryStorage describes storage that stores all metrics in RAM.
type MemoryStorage struct {
	gauge         map[string]float64
	counter       map[string]int64
	logger        Logger
	fileSaver     FileSaver
	isToRestore   bool
	storeInterval time.Duration
	mu            sync.RWMutex
}

// NewMemoryStorage creates the instance of MemoryStorage and configures it with params from Configurator.
func NewMemoryStorage(config Configurator, logger Logger) *MemoryStorage {
	storage := &MemoryStorage{
		gauge:         make(map[string]float64),
		counter:       make(map[string]int64),
		logger:        logger,
		isToRestore:   config.GetRestoreFlag(),
		storeInterval: config.GetStoreInterval(),
	}

	fileSaver := jsonfilesaver.NewRawFileSaver(config.GetStoreFilePath())
	storage.fileSaver = fileSaver

	storage.restoreStorageFromFile()

	return storage
}

// CheckMetricAndSave performs validating of pcstats.Metric and saving to storage.
func (s *MemoryStorage) CheckMetricAndSave(ctx context.Context, metric pcstats.Metric) error {
	defer s.saveStorageToFileIfSynchroSavingEnabled(ctx)

	err := pcstats.CheckMetricParamsIsValid(metric)
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}

	if metric.GetType() == pcstats.Gauge {
		return s.saveGaugeMetric(metric)
	}

	return s.saveCounterMetric(metric)
}

// CheckMetricBatchAndSave performs checking of each pcstats.Metric in batch and then saves them to storage.
func (s *MemoryStorage) CheckMetricBatchAndSave(ctx context.Context, batch pcstats.Metrics) error {
	err := CheckMetricBatchIsValid(batch)
	if err != nil {
		return errors.Wrap(err, "cannot save batch")
	}

	s.saveBatch(ctx, batch)
	return nil
}

// GetMetric tries to get value of metric with input type and ID from storage. Then returns new metric with got params if no
// error occurred.
func (s *MemoryStorage) GetMetric(
	ctx context.Context,
	metricID string,
	metricType pcstats.MetricType,
) (*pcstats.Metric, error) {
	switch metricType {
	case pcstats.Gauge:
		return s.getGaugeMetric(metricID)
	case pcstats.Counter:
		return s.getCounterMetric(metricID)
	default:
		return nil, errors.Wrap(pcstats.ErrBadType, metricType.String())
	}
}

// GetAllMetrics returns all metrics from storage.
func (s *MemoryStorage) GetAllMetrics(ctx context.Context) (pcstats.Metrics, error) {
	outMetrics := make(pcstats.Metrics, 0, len(s.gauge)+len(s.counter))

	gaugeMetrics := s.getAllGauges()
	counterMetrics := s.getAllCounters()

	outMetrics = append(outMetrics, gaugeMetrics...)
	outMetrics = append(outMetrics, counterMetrics...)

	return outMetrics, nil

}

// RunStorage starts interval saving of storage into file if store interval is not 0.
func (s *MemoryStorage) RunStorage(ctx context.Context) {
	if s.storeInterval != 0 {
		s.startIntervalSaving(ctx)
	}
}

func (s *MemoryStorage) Ping(ctx context.Context) error {
	return nil
}

func (s *MemoryStorage) startIntervalSaving(ctx context.Context) {
	saveTicker := time.NewTicker(s.storeInterval)
	for {
		select {
		case <-saveTicker.C:
			metrics, _ := s.GetAllMetrics(ctx)
			err := s.saveStorageToFile(metrics)
			if err != nil {
				s.logger.Error(err.Error())
			}
		case <-ctx.Done():
			saveTicker.Stop()
			s.logger.Info("interval saving context stopped now")
			return
		}
	}
}

func (s *MemoryStorage) saveGaugeMetric(metric pcstats.Metric) error {
	value, err := metric.GetGaugeValue()
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}

	s.mu.Lock()
	s.gauge[metric.GetID()] = value
	s.mu.Unlock()

	return nil
}

func (s *MemoryStorage) saveCounterMetric(metric pcstats.Metric) error {
	delta, err := metric.GetCounterValue()
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}

	s.mu.Lock()
	s.counter[metric.GetID()] += delta
	s.mu.Unlock()

	return nil
}

func (s *MemoryStorage) getGaugeMetric(metricID string) (*pcstats.Metric, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if value, ok := s.gauge[metricID]; ok {
		inputGaugeValue := new(float64)
		*inputGaugeValue = value
		metric, err := pcstats.NewMetric(metricID, pcstats.Gauge, nil, inputGaugeValue)
		if err != nil {
			s.logger.Error(err.Error())
			return nil, errors.Wrap(err, "cannot create metric")
		}

		return metric, nil
	}

	return nil, ErrGaugeMetricNotFound
}

func (s *MemoryStorage) getCounterMetric(metricID string) (*pcstats.Metric, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if delta, ok := s.counter[metricID]; ok {
		inputCounterValue := new(int64)
		*inputCounterValue = delta
		metric, err := pcstats.NewMetric(metricID, pcstats.Counter, inputCounterValue, nil)
		if err != nil {
			s.logger.Error(err.Error())
			return nil, errors.Wrap(err, "cannot create metric")
		}

		return metric, nil
	}

	return nil, ErrCounterMetricNotFound
}

func (s *MemoryStorage) getAllGauges() pcstats.Metrics {
	s.mu.RLock()
	defer s.mu.RUnlock()

	metrics := make(pcstats.Metrics, 0, len(s.gauge))

	for ID, value := range s.gauge {
		inputGaugeValue := new(float64)
		*inputGaugeValue = value
		metric, err := pcstats.NewMetric(ID, pcstats.Gauge, nil, inputGaugeValue)
		if err != nil {
			s.logger.Error(err.Error())
			continue
		}
		metrics = append(metrics, *metric)
	}

	return metrics
}

func (s *MemoryStorage) getAllCounters() pcstats.Metrics {
	s.mu.RLock()
	defer s.mu.RUnlock()

	metrics := make(pcstats.Metrics, 0, len(s.counter))

	for ID, delta := range s.counter {
		inputCounterValue := new(int64)
		*inputCounterValue = delta
		metric, err := pcstats.NewMetric(ID, pcstats.Counter, inputCounterValue, nil)
		if err != nil {
			s.logger.Error(err.Error())
			continue
		}
		metrics = append(metrics, *metric)
	}

	return metrics
}

func (s *MemoryStorage) saveStorageToFile(metrics pcstats.Metrics) (err error) {
	defer func() {
		if err != nil {
			s.logger.Error(errors.Wrap(err, "storage content has not been saved").Error())
			return
		}
		s.logger.Info("storage successfully saved to file")
	}()

	jsonData, err := metrics.MarshalJSON()
	if err != nil {
		return err
	}
	return s.fileSaver.Save(jsonData)
}

func (s *MemoryStorage) saveStorageToFileIfSynchroSavingEnabled(ctx context.Context) {
	if s.storeInterval == 0 {
		metrics, _ := s.GetAllMetrics(ctx)
		err := s.saveStorageToFile(metrics)

		if err != nil {
			s.logger.Error(err.Error())
		}
	}
}

func (s *MemoryStorage) saveBatch(ctx context.Context, batch pcstats.Metrics) {
	for _, metric := range batch {
		err := s.CheckMetricAndSave(ctx, metric)
		if err != nil {
			s.logger.Error(err.Error() + "; metric hasn't been saved")
		}
	}
}

func (s *MemoryStorage) restoreStorageFromFile() {
	if s.isToRestore {
		data, err := s.fileSaver.Load()
		if err != nil {
			s.logger.Error(errors.Wrap(err, "cannot restore storage from file").Error())
			return
		}

		metrics, err := JSONDataToMetrics(data)
		if err != nil {
			s.logger.Error(errors.Wrap(err, "cannot restore storage from file").Error())
			return
		}

		err = s.CheckMetricBatchAndSave(context.Background(), metrics)
		if err != nil {
			s.logger.Error(errors.Wrap(err, "cannot restore storage from file").Error())
			return
		}

		s.logger.Info("storage successfully restored from file")
		return
	}
	s.logger.Info("restore flag has not been set, storage is not to be restored")
}
