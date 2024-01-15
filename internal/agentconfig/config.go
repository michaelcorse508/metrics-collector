package agentconfig

import (
	"time"

	"github.com/bazookajoe1/metrics-collector/internal/auxiliary"
)

type Config struct {
	address        string
	port           string
	pollInterval   time.Duration
	reportInterval time.Duration
	secretKey      []byte
	workersLimit   uint64
}

func NewConfig() *Config {
	c := &Config{
		pollInterval:   2 * time.Second,
		reportInterval: 10 * time.Second,
	}

	clArgsParams := ArgParse()
	envArgsParams := EnvParse()

	c.UpdateConfig(clArgsParams, envArgsParams)

	return c
}

func (c *Config) UpdateConfig(p ...Parameters) {
	for _, paramInstance := range p {
		address := paramInstance.GetAddress()
		c.updateAddress(address)

		port := paramInstance.GetPort()
		c.updatePort(port)

		pollInterval := paramInstance.GetPollInterval()
		c.updatePollInterval(pollInterval)

		reportInterval := paramInstance.GetReportInterval()
		c.updateReportInterval(reportInterval)

		secretKey := paramInstance.GetSecretKey()
		c.updateSecretKey(secretKey)

		workersLimit := paramInstance.GetWorkersLimit()
		c.updateWorkersLimit(workersLimit)
	}
}

func (c *Config) GetAddress() string {
	return c.address
}

func (c *Config) GetPort() string {
	return c.port
}

func (c *Config) GetPollInterval() time.Duration {
	return c.pollInterval
}

func (c *Config) GetReportInterval() time.Duration {
	return c.reportInterval
}

func (c *Config) GetSecretKey() []byte {
	return c.secretKey
}

func (c *Config) GetWorkersLimit() uint64 {
	return c.workersLimit
}

func (c *Config) updateAddress(address string) {
	err := auxiliary.ValidateIP(address)
	if err == nil {
		c.address = address
	}
}

func (c *Config) updatePort(port string) {
	err := auxiliary.ValidatePort(port)
	if err == nil {
		c.port = port
	}
}

func (c *Config) updatePollInterval(pollInterval uint64) {
	if pollInterval != InvalidDuration {
		c.pollInterval = time.Duration(pollInterval) * time.Second
	}
}

func (c *Config) updateReportInterval(reportInterval uint64) {
	if reportInterval != InvalidDuration {
		c.reportInterval = time.Duration(reportInterval) * time.Second
	}
}

func (c *Config) updateSecretKey(secretKey string) {
	if secretKey == "" {
		return
	}
	c.secretKey = []byte(secretKey)
}

func (c *Config) updateWorkersLimit(workersLimit uint64) {
	if workersLimit == 0 {
		return
	}
	c.workersLimit = workersLimit
}
