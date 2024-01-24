package agent

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/bazookajoe1/metrics-collector/internal/agent/collector"
	"github.com/bazookajoe1/metrics-collector/internal/agentconfig"
	"github.com/bazookajoe1/metrics-collector/internal/auxiliary"
	"github.com/bazookajoe1/metrics-collector/internal/compressing"
	"github.com/bazookajoe1/metrics-collector/internal/cryptography"
	"github.com/bazookajoe1/metrics-collector/internal/logging"
	"github.com/bazookajoe1/metrics-collector/internal/pcstats"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

const ContentEncodingGZIP = "gzip"

var ReqHeaderParams = map[string]string{
	echo.HeaderContentType:     echo.MIMEApplicationJSONCharsetUTF8,
	echo.HeaderAcceptEncoding:  ContentEncodingGZIP,
	echo.HeaderContentEncoding: ContentEncodingGZIP,
}

type HTTPAgent struct {
	client           *resty.Client
	config           Configurator
	collector        Collector
	logger           Logger
	sendURL          string
	workerController *WorkerController
}

func HTTPAgentNew(config Configurator, collector Collector, logger Logger) (*HTTPAgent, error) {
	agent := &HTTPAgent{
		client:           resty.New(),
		config:           config,
		collector:        collector,
		logger:           logger,
		workerController: NewWorkerController(config),
	}

	agent.client = agent.client.SetRetryCount(3).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(5 * time.Second)

	err := agent.calculateSendURL()
	if err != nil {
		return nil, err
	}

	agent.client.SetLogger(logger)
	return agent, nil
}

func (a *HTTPAgent) SendMetrics(ctx context.Context, metrics pcstats.Metrics) {
	metricsJSON, err := metrics.MarshalJSON()
	if err != nil {
		a.logger.Error(errors.Wrap(err, "cannot send metrics").Error())
		return
	}

	gzMetricsJSON, err := compressing.GZIPCompress(metricsJSON)
	if err != nil {
		a.logger.Error(errors.Wrap(err, "cannot send metrics").Error())
		return
	}

	reqCtx, reqCancel := context.WithTimeout(ctx, 2*time.Second)
	defer reqCancel()

	request, err := a.SignedRequestConstructor(reqCtx, metricsJSON).
		SetBody(gzMetricsJSON).
		Post(a.sendURL)
	if err != nil {
		a.logger.Error(errors.Wrap(err, "cannot send metrics").Error())
	}
	a.logger.Debugf("request body: %s\nresponse code:%d", metricsJSON, request.StatusCode())
}

func (a *HTTPAgent) RunSender(ctx context.Context) {
	workerWaitGroup := a.StartWorkers(ctx)

	ticker := time.NewTicker(a.config.GetReportInterval())

	for {
		select {
		case <-ticker.C:
			metrics := a.collector.GetMetrics()
			a.workerController.PutMetricsInQueue(metrics) // sending metrics to workers
		case <-ctx.Done():
			a.workerController.CloseQueue()
			ticker.Stop()
			workerWaitGroup.Wait()
			a.logger.Info("sender context cancelled; return")
			return
		}
	}
}

func (a *HTTPAgent) SignedRequestConstructor(ctx context.Context, reqBody []byte) *resty.Request {

	sign, err := a.calculateHMACIfKeyIsSet(reqBody)
	if err != nil {
		return a.client.R().SetContext(ctx).SetHeaders(ReqHeaderParams)
	}
	signString := hex.EncodeToString(sign)
	return a.client.R().
		SetContext(ctx).
		SetHeaders(ReqHeaderParams).
		SetHeader(cryptography.HeaderHashSHA256, signString)
}

func (a *HTTPAgent) Worker(ctx context.Context, workerID uint64) {
	for metrics := range a.workerController.sendQueue {
		a.logger.Debugf("worker %d job\n", workerID)
		a.SendMetrics(ctx, metrics)
	}
	a.logger.Debugf("worker %d shutdown\n", workerID)
}

func (a *HTTPAgent) StartWorkers(ctx context.Context) *sync.WaitGroup {
	workerWaitGroup := sync.WaitGroup{}

	for workerID := uint64(0); workerID < a.workerController.workersLimit; workerID++ {
		workerID := workerID
		workerWaitGroup.Add(1)
		go func() {
			a.Worker(ctx, workerID)
			workerWaitGroup.Done()
		}()
	}

	return &workerWaitGroup
}

func (a *HTTPAgent) calculateSendURL() error {
	address := a.config.GetAddress()
	err := auxiliary.ValidateIP(address)
	if err != nil {
		return errors.Wrap(err, "cannot calculate send address")
	}

	port := a.config.GetPort()
	err = auxiliary.ValidatePort(port)
	if err != nil {
		return errors.Wrap(err, "cannot calculate send address")
	}

	a.sendURL = fmt.Sprintf("http://%s:%s/updates/", address, port)

	return nil
}

func (a *HTTPAgent) calculateHMACIfKeyIsSet(message []byte) ([]byte, error) {
	if len(a.config.GetSecretKey()) < 1 {
		return nil, fmt.Errorf("secret key has not been set; message will not be signed")
	}

	sign, err := cryptography.SignMessageHMAC(message, a.config.GetSecretKey())
	if err != nil {
		return nil, fmt.Errorf("cannot calculate sign on given message; message will not be signed")
	}

	return sign, nil
}

func RunAgent() error {
	wg := sync.WaitGroup{}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM|syscall.SIGINT|syscall.SIGQUIT,
	)
	defer stop()

	logger := logging.NewZapLogger()

	config := agentconfig.NewConfig()

	metricCollector := collector.NewCollector(collector.NeededMetrics, config, logger)

	agent, err := HTTPAgentNew(config, metricCollector, logger)
	fmt.Println("AGENT SECRET KEY =============>", agent.config.GetSecretKey())
	if err != nil {
		logger.Error(errors.Wrap(err, "cannot create agent").Error())
		return err
	}

	wg.Add(1)
	go func() {
		agent.RunSender(ctx)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		agent.collector.Run(ctx)
		wg.Done()
	}()

	wg.Wait()
	logger.Info("shutting down agent")
	return nil
}
