package server

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/bazookajoe1/metrics-collector/internal/logging"
	"github.com/bazookajoe1/metrics-collector/internal/serverconfig"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
)

type HTTPServer struct {
	server  *echo.Echo
	config  ServConfigurator
	storage Storage
	logger  Logger
}

func HTTPServerNew(c ServConfigurator, s Storage, l Logger) *HTTPServer {
	return &HTTPServer{
		server:  echo.New(),
		config:  c,
		storage: s,
		logger:  l,
	}
}

func (s *HTTPServer) InitRoutes() {
	s.server.Use(middleware.Decompress())
	s.server.Use(middleware.Logger())
	s.server.Use(middleware.Gzip())
	s.server.Use(s.HMACChecker)
	s.server.Use(s.HMACSigner)
	s.server.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Header.Get("Content-Type") != echo.MIMEApplicationJSON {
				return next(c)
			}

			data, err := io.ReadAll(c.Request().Body)
			if err != nil {
				return err
			}

			c.Logger().Error(string(data))
			c.Request().Body = io.NopCloser(bytes.NewReader(data))

			return next(c)
		}
	})

	s.server.RouteNotFound(
		"/*",
		func(c echo.Context) error { return c.NoContent(http.StatusNotFound) },
	)

	s.server.GET("/", s.SendAllMetricsHTML)

	s.server.GET("/ping", s.Ping)

	s.server.POST("/updates/", s.ReceiveBatchOfMetricsJSON)

	gUpdate := s.server.Group("/update")
	gUpdate.POST("/:type/:id/:value", s.ReceiveMetricFromURLParams)
	gUpdate.POST("/", s.ReceiveMetricFromBodyJSON)

	gValue := s.server.Group("/value")
	gValue.GET("/:type/:id", s.SendMetricText)
	gValue.POST("/", s.SendMetricJSON)
}

func (s *HTTPServer) Run() {
	aP := fmt.Sprintf("%s:%s", s.config.GetAddress(), s.config.GetPort())
	s.server.HideBanner = true
	if err := s.server.Start(aP); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Fatal("shutting down the server")
	}
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func RunServer() error {
	wg := sync.WaitGroup{}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM|syscall.SIGINT|syscall.SIGQUIT,
	)
	defer stop()

	logger := logging.NewZapLogger()

	config := serverconfig.NewConfig()

	storage := ChooseStorage(ctx, config, logger)

	storageCtx, storageCancel := context.WithCancel(ctx)
	defer storageCancel()
	wg.Add(1)
	go func() {
		storage.RunStorage(storageCtx)
		wg.Done()
	}()

	server := HTTPServerNew(config, storage, logger)
	server.InitRoutes()
	go server.Run()

	wg.Wait()
	logger.Info("all contexts done; starting to stop server")

	serverCtx, serverCancel := context.WithTimeout(ctx, 5*time.Second)
	defer serverCancel()

	if err := server.Stop(serverCtx); err != nil {
		return err
	}

	return nil
}
