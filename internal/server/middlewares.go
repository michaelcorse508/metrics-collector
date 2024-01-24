package server

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"

	"github.com/bazookajoe1/metrics-collector/internal/cryptography"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func (s *HTTPServer) HMACSigner(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		secretKey := s.config.GetSecretKey()
		if len(secretKey) < 1 {
			s.logger.Debug("secret key has not been set; message will not be signed")
			return next(c)
		}

		dsResponseWrite := cryptography.DsResponseWriter{
			ResponseWriter: c.Response().Writer,
			SecretKey:      secretKey,
		}

		c.Response().Writer = &dsResponseWrite
		return next(c)
	}
}

func (s *HTTPServer) HMACChecker(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		secretKey := s.config.GetSecretKey()
		if len(secretKey) < 1 {
			s.logger.Debug("secret key has not been set; message will not be checked")
			return next(c)
		}

		sign, err := GetHMAC(c)
		if err != nil {
			s.logger.Debug("client has not set digital sign; message will not be checked")
			return next(c)
		}

		message, err := ReadRequestBody(c)
		if err != nil {
			s.logger.Debug(errors.Wrap(err, "cannot check HMAC").Error())
			return c.NoContent(http.StatusInternalServerError)
		}

		switch cryptography.CheckMessageHMAC(message, sign, s.config.GetSecretKey()) {
		case cryptography.ErrInvalidSign:
			s.logger.Debug("client has given invalid digital sign")
			return c.NoContent(http.StatusBadRequest)
		case cryptography.ErrSignCalculatingError:
			s.logger.Debug("cannot calculate sigh on given message")
			return c.NoContent(http.StatusInternalServerError)
		default:
			s.logger.Debug("sign is valid; continue")
			return next(c)
		}

	}
}

func (s *HTTPServer) LogHeaders(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		headers := "Headers: "
		for key, value := range c.Request().Header {
			headers += fmt.Sprintf("{%s:%s}", key, value)
			headers += ", \n"
		}
		s.logger.Debug(headers)

		return next(c)
	}
}

func ReadRequestBody(c echo.Context) ([]byte, error) {
	requestBody, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read request body")
	}

	newBody := bytes.NewBuffer(requestBody)
	c.Request().Body = io.NopCloser(newBody)

	return requestBody, nil
}

func GetHMAC(c echo.Context) ([]byte, error) {
	value := c.Request().Header.Get(cryptography.HeaderHashSHA256)
	if value == "" {
		return nil, fmt.Errorf("no HashSHA256 header set")
	}

	hexValue, err := hex.DecodeString(value)
	if err != nil {
		return nil, fmt.Errorf("broken HashSHA256")
	}

	return hexValue, nil
}
