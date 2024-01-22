package server

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/bazookajoe1/metrics-collector/internal/pcstats"
	"github.com/labstack/echo/v4"
)

// ReceiveMetricFromURLParams is the handler responsible for receiving metrics from request uri.
func (s *HTTPServer) ReceiveMetricFromURLParams(c echo.Context) error {
	metric, err := pcstats.NewMetricFromString(
		c.Param("id"),
		c.Param("type"),
		c.Param("value"),
	)

	if err != nil {
		s.logger.Error(err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	err = s.storage.CheckMetricAndSave(c.Request().Context(), *metric)
	if err != nil {
		s.logger.Error(err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

// SendMetricText is the handler responsible for sending metric value got from storage with parameters from uri.
func (s *HTTPServer) SendMetricText(c echo.Context) error {
	metric, err := s.storage.GetMetric(
		c.Request().Context(),
		c.Param("id"),
		pcstats.MetricType(c.Param("type")),
	)
	if err != nil {
		s.logger.Debug(err.Error())
		return c.NoContent(http.StatusNotFound)
	}

	stringValue, err := metric.GetStringValue()
	if err != nil {
		s.logger.Debug(err.Error())
		return c.NoContent(
			http.StatusInternalServerError,
		) // this case should not occur, but who knows
	}

	return c.String(http.StatusOK, stringValue)
}

// SendAllMetricsHTML is the handler responsible for sending back all metrics got from storage in html string format.
func (s *HTTPServer) SendAllMetricsHTML(c echo.Context) error {
	metrics, err := s.storage.GetAllMetrics(c.Request().Context())
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	responseString := MetricSliceToHTMLString(metrics)

	return c.HTMLBlob(http.StatusOK, []byte(responseString))
}

// ReceiveMetricFromBodyJSON is the handler responsible for receiving metrics in JSON format from body.
func (s *HTTPServer) ReceiveMetricFromBodyJSON(c echo.Context) error {
	metric := new(pcstats.Metric)

	buf := new(strings.Builder)
	_, _ = io.Copy(buf, c.Request().Body)
	// check errors
	fmt.Println(buf.String())

	if err := c.Bind(metric); err != nil {
		return err
	} // binding request body to pcstats.Metric

	err := s.storage.CheckMetricAndSave(c.Request().Context(), *metric)
	if err != nil {
		s.logger.Error(err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	responseMetric, err := s.storage.GetMetric(
		c.Request().Context(),
		metric.GetID(),
		metric.GetType(),
	)
	if err != nil {
		s.logger.Error(err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	sendData, err := responseMetric.MarshalJSON() // serializing into JSON
	if err != nil {
		s.logger.Error(err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSONBlob(http.StatusOK, sendData)
}

// SendMetricJSON is the handler responsible for getting metric from storage according to parameters got
// // from body in JSON format.
func (s *HTTPServer) SendMetricJSON(c echo.Context) error {
	metric := new(pcstats.Metric)

	if err := c.Bind(metric); err != nil {
		return err
	}

	responseMetric, err := s.storage.GetMetric(
		c.Request().Context(),
		metric.GetID(),
		metric.GetType(),
	)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	sendData, err := responseMetric.MarshalJSON()
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSONBlob(http.StatusOK, sendData)
}

// ReceiveBatchOfMetricsJSON is the handler responsible for saving of metric batch in json format.
func (s *HTTPServer) ReceiveBatchOfMetricsJSON(c echo.Context) error {
	metrics := new(pcstats.Metrics)
	if err := c.Bind(metrics); err != nil {
		return err
	}

	err := s.storage.CheckMetricBatchAndSave(c.Request().Context(), *metrics)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	return c.NoContent(http.StatusOK)
}

// Ping is the handler that checks connection with storage.
func (s *HTTPServer) Ping(c echo.Context) error {
	err := s.storage.Ping(c.Request().Context())
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}
