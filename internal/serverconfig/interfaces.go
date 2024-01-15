package serverconfig

type Parameters interface {
	GetAddress() string
	GetPort() string
	GetStoreInterval() int
	GetFilePath() string
	GetRestoreFlag() *bool
	GetDatabaseDSN() string
	GetSecretKey() string
}
