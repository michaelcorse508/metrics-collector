package serverconfig

import (
	"github.com/bazookajoe1/metrics-collector/internal/auxiliary"
	"sync"
	"time"
)

type Config struct {
	address       string
	port          string
	databaseDSN   string
	storeInterval time.Duration
	restoreFlag   bool
	storeFilePath string
	secretKey     []byte
	mu            sync.Mutex
}

// NewConfig creates Config instance with parameters collected from
// environment variables and command line arguments.
func NewConfig() *Config {
	c := &Config{}

	clArgParams := ArgParse()
	envParams := EnvParse()

	c.UpdateConfig(clArgParams, envParams)

	return c
}

// UpdateConfig sequentially fill Config with data provided from
// different sources (cl args and env vars). Last provided data (if exists) overrides
// previous provided data.
func (c *Config) UpdateConfig(p ...Parameters) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, paramInstance := range p {
		address := paramInstance.GetAddress()
		c.updateAddress(address)

		port := paramInstance.GetPort()
		c.updatePort(port)

		databaseDSN := paramInstance.GetDatabaseDSN()
		c.updateDatabaseDSN(databaseDSN)

		storeFilePath := paramInstance.GetFilePath()
		c.updateStoreFilePath(storeFilePath)

		storeInterval := paramInstance.GetStoreInterval()
		c.updateStoreInterval(storeInterval)

		restoreFlag := paramInstance.GetRestoreFlag()
		c.updateRestoreFlag(restoreFlag)

		secretKey := paramInstance.GetSecretKey()
		c.updateSecretKey(secretKey)
	}
}

func (c *Config) GetAddress() string {
	return c.address
}

func (c *Config) GetPort() string {
	return c.port
}

func (c *Config) GetDatabaseDSN() string {
	return c.databaseDSN
}

func (c *Config) GetRestoreFlag() bool {
	return c.restoreFlag
}

func (c *Config) GetStoreFilePath() string {
	return c.storeFilePath
}

func (c *Config) GetStoreInterval() time.Duration {
	return c.storeInterval
}

func (c *Config) GetSecretKey() []byte {
	return c.secretKey
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

func (c *Config) updateDatabaseDSN(databaseDSN string) {
	if databaseDSN != "" {
		c.databaseDSN = databaseDSN
	}
}

func (c *Config) updateStoreFilePath(storeFilePath string) {
	if storeFilePath != "" {
		c.storeFilePath = storeFilePath
	}
}

func (c *Config) updateStoreInterval(storeInterval int) {
	if storeInterval != InvalidDuration {
		c.storeInterval = time.Duration(storeInterval) * time.Second
	}
}

func (c *Config) updateRestoreFlag(restoreFlag *bool) {
	if restoreFlag != nil {
		c.restoreFlag = *restoreFlag
	}
}

func (c *Config) updateSecretKey(secretKey string) {
	if secretKey == "" {
		return
	}
	c.secretKey = []byte(secretKey)
}
