package serverconfig

import (
	"github.com/bazookajoe1/metrics-collector/internal/auxiliary"
	"github.com/caarlos0/env/v9"
)

const InvalidDuration int = 123456789

type EnvParams struct {
	Address       []string `env:"ADDRESS" envSeparator:":"`
	Host          string
	Port          string
	StoreInterval int    `env:"STORE_INTERVAL"`
	FilePath      string `env:"FILE_STORAGE_PATH"`
	Restore       *bool  `env:"RESTORE"`
	DatabaseDSN   string `env:"DATABASE_DSN"`
	SecretKey     string `env:"KEY"`
}

func (a *EnvParams) GetAddress() string {
	return a.Host
}

func (a *EnvParams) GetPort() string {
	return a.Port
}

func (a *EnvParams) GetStoreInterval() int {
	return a.StoreInterval
}

func (a *EnvParams) GetFilePath() string {
	return a.FilePath
}

func (a *EnvParams) GetRestoreFlag() *bool {
	return a.Restore
}

func (a *EnvParams) GetDatabaseDSN() string {
	return a.DatabaseDSN
}

func (a *EnvParams) GetSecretKey() string {
	return a.SecretKey
}

func EnvParse() *EnvParams {
	var ep = &EnvParams{StoreInterval: InvalidDuration}

	err := env.Parse(ep)
	if err != nil {
		return nil
	}
	ep.splitHostFromPort()

	return ep
}

func (a *EnvParams) splitHostFromPort() {
	if len(a.Address) >= 2 {
		err := auxiliary.ValidateIP(a.Address[0])
		if err != nil {
			return
		}
		err = auxiliary.ValidatePort(a.Address[1])
		if err != nil {
			return
		}
		a.Host = a.Address[0]
		a.Port = a.Address[1]
	}
}
