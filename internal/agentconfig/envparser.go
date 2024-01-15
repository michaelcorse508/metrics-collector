package agentconfig

import (
	"github.com/bazookajoe1/metrics-collector/internal/auxiliary"
	"github.com/caarlos0/env/v9"
)

const InvalidDuration uint64 = 123456789

type EnvParams struct {
	Endpoint       []string `env:"ADDRESS"         envSeparator:":"`
	Address        string
	Port           string
	ReportInterval uint64 `env:"REPORT_INTERVAL"`
	PollInterval   uint64 `env:"POLL_INTERVAL"`
	SecretKey      string `env:"KEY"`
	WorkersLimit   uint64 `env:"RATE_LIMIT"`
}

func (a *EnvParams) GetAddress() string {
	return a.Address
}

func (a *EnvParams) GetPort() string {
	return a.Port
}

func (a *EnvParams) GetPollInterval() uint64 {
	return a.PollInterval
}

func (a *EnvParams) GetReportInterval() uint64 {
	return a.ReportInterval
}

func (a *EnvParams) GetSecretKey() string {
	return a.SecretKey
}

func (a *EnvParams) GetWorkersLimit() uint64 {
	return a.WorkersLimit
}

func EnvParse() *EnvParams {
	var ep = &EnvParams{
		PollInterval:   InvalidDuration,
		ReportInterval: InvalidDuration,
	}

	err := env.Parse(ep)
	if err != nil {
		return nil
	}
	ep.splitAddressFromPort()

	return ep
}

func (a *EnvParams) splitAddressFromPort() {
	if len(a.Endpoint) >= 2 {
		err := auxiliary.ValidateIP(a.Endpoint[0])
		if err != nil {
			return
		}
		err = auxiliary.ValidatePort(a.Endpoint[1])
		if err != nil {
			return
		}
		a.Address = a.Endpoint[0]
		a.Port = a.Endpoint[1]
	}
}
