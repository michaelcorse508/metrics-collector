package collector

import (
	"context"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/bazookajoe1/metrics-collector/internal/pcstats"
)

type Collector struct {
	stats        map[string]pcstats.Metric
	pollInterval time.Duration
	logger       Logger
	mu           sync.RWMutex
}

// NewCollector create the instance of Collector. Pool of collected metrics is assigned in neededMetrics.
func NewCollector(neededMetrics pcstats.Metrics, config Configurator, logger Logger) *Collector {
	c := &Collector{
		stats:        make(map[string]pcstats.Metric),
		pollInterval: config.GetPollInterval(),
		logger:       logger,
	}

	c.FillCollectorWithMetrics(neededMetrics)

	return c
}

func (c *Collector) FillCollectorWithMetrics(neededMetrics pcstats.Metrics) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, metric := range neededMetrics {
		c.stats[metric.GetID()] = metric
	}
}

// Collect collects metric values from runtime.MemStats and updates appropriate metric in Collector.
// PollCount increments automatically. RandomValue also is updated.
func (c *Collector) Collect() {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	for ID := range c.stats {
		value, err := GetMemStatsFieldValueByName(&stats, ID)
		if err != nil {
			continue
		}

		c.update(ID, value)
	}
	c.updatePollCount()
	c.updateRandomValue()
}

// GetMetrics return all metrics from Collector.
func (c *Collector) GetMetrics() pcstats.Metrics {
	metrics := make(pcstats.Metrics, 0, len(c.stats))

	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, metric := range c.stats {
		metrics = append(metrics, metric)
	}

	return metrics
}

// Run starts Collector. Metrics are collected with interval = pollInterval.
func (c *Collector) Run(ctx context.Context) {
	ticker := time.NewTicker(c.pollInterval)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			c.logger.Info("collector context cancelled; return")
			return
		case <-ticker.C:
			c.Collect()
		}
	}
}

func (c *Collector) update(ID string, value float64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.stats[ID]; ok {
		*c.stats[ID].Value = value
	}
}

func (c *Collector) updatePollCount() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.stats["PollCount"]; ok {
		*c.stats["PollCount"].Delta += 1
	}
}

func (c *Collector) updateRandomValue() {
	for {
		randomValue := rand.NormFloat64()
		if randomValue != 0 { // we don't need zero random value
			c.mu.Lock()
			*c.stats["RandomValue"].Value = randomValue
			c.mu.Unlock()
			break
		}
	}
}
