package http

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"github.com/arnaz06/soccer"
)

// TimeoutMiddleware is used to add timeout for context cancellation.
func TimeoutMiddleware(timeout time.Duration) echo.MiddlewareFunc {
	return func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctxWithTimeout, cancel := context.WithTimeout(c.Request().Context(), timeout)
			defer cancel()
			req := c.Request().WithContext(ctxWithTimeout)
			c.SetRequest(req)
			return handlerFunc(c)
		}
	}
}

// ErrorMiddleware is a function to generate http status code.
func ErrorMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err == nil {
				return nil
			}

			headers := []string{}
			for k, v := range c.Request().Header {
				headers = append(headers, fmt.Sprintf(`%s:%s`, k, strings.Join(v, ",")))
			}
			lg := log.WithFields(log.Fields{
				"headers": strings.Join(headers, " | "),
				"method":  c.Request().Method,
				"uri":     c.Request().RequestURI,
			})

			if e, ok := err.(*echo.HTTPError); ok {
				if e.Code >= http.StatusInternalServerError {
					lg.Errorln(e.Message)
				}
				return echo.NewHTTPError(e.Code, e.Message)
			}

			if _, ok := err.(soccer.ConstraintError); ok {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			switch err {
			case context.DeadlineExceeded, context.Canceled:
				lg.Errorln(err.Error())
				return echo.NewHTTPError(http.StatusRequestTimeout, err.Error())

			case soccer.ErrNotFound:
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			}

			lg.Errorln(err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
}
