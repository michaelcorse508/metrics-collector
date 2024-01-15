package memorystorage

import "time"

type FileSaver interface {
	Load() ([]byte, error)
	Save(data []byte) error
}

// Logger is the interfaces that allows to work with different loggers.
type Logger interface {
	Info(string)
	Debug(string)
	Error(string)
}

type Configurator interface {
	GetRestoreFlag() bool
	GetStoreInterval() time.Duration
	GetStoreFilePath() string
}
