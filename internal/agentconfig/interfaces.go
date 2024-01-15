package agentconfig

type Parameters interface {
	GetAddress() string
	GetPort() string
	GetReportInterval() uint64
	GetPollInterval() uint64
	GetSecretKey() string
	GetWorkersLimit() uint64
}
