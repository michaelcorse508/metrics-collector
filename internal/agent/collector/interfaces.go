package collector

import "time"

type Configurator interface {
	GetPollInterval() time.Duration
}

type Logger interface {
	Info(string)
	Debug(string)
	Error(string)
	Fatal(string)
}
