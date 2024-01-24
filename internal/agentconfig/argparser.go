package agentconfig

import (
	"flag"
	"fmt"
	"github.com/bazookajoe1/metrics-collector/internal/auxiliary"
	"strings"
)

type CLArgsParams struct {
	Address        string
	Port           string
	PollInterval   uint64
	ReportInterval uint64
	SecretKey      string
	WorkersLimit   uint64
}

func (a *CLArgsParams) String() string {
	return a.Address + ":" + a.Port
}

func (a *CLArgsParams) Set(s string) error {
	hp := strings.Split(s, ":")
	if len(hp) != 2 {
		return fmt.Errorf("need address in a form Host:Port")
	}

	err := auxiliary.ValidateIP(hp[0])
	if err != nil {
		return err
	}

	err = auxiliary.ValidatePort(hp[1])
	if err != nil {
		return err
	}

	a.Address = hp[0]
	a.Port = hp[1]

	return nil
}

func (a *CLArgsParams) GetAddress() string {
	return a.Address
}

func (a *CLArgsParams) GetPort() string {
	return a.Port
}

func (a *CLArgsParams) GetReportInterval() uint64 {
	return a.ReportInterval
}

func (a *CLArgsParams) GetPollInterval() uint64 {
	return a.PollInterval
}

func (a *CLArgsParams) GetSecretKey() string {
	return a.SecretKey
}

func (a *CLArgsParams) GetWorkersLimit() uint64 {
	return a.WorkersLimit
}

func ArgParse() *CLArgsParams {
	params := &CLArgsParams{
		Address:        "localhost",
		Port:           "8080",
		ReportInterval: 10,
		PollInterval:   2,
	}

	flag.Var(params, "a", "Server listen point in format: `Host:Port`")
	flag.Uint64Var(&params.ReportInterval, "r", 10, "Report interval in seconds")
	flag.Uint64Var(&params.PollInterval, "p", 2, "Collect interval in seconds")
	flag.StringVar(&params.SecretKey, "k", "", "Sign key")
	flag.Uint64Var(&params.WorkersLimit, "l", 10, "Rate limit of request senders")
	flag.Parse()

	return params
}
