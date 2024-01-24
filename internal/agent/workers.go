package agent

import (
	"github.com/bazookajoe1/metrics-collector/internal/pcstats"
)

type WCConfigurator interface {
	GetWorkersLimit() uint64
}

type WorkerController struct {
	workersLimit uint64
	sendQueue    chan pcstats.Metrics
}

func NewWorkerController(config WCConfigurator) *WorkerController {
	return &WorkerController{
		workersLimit: config.GetWorkersLimit(),
		sendQueue:    make(chan pcstats.Metrics, 10),
	}
}

func (wc *WorkerController) PutMetricsInQueue(metrics pcstats.Metrics) {
	wc.sendQueue <- metrics
}

func (wc *WorkerController) CloseQueue() {
	close(wc.sendQueue)
}
