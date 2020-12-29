package http

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"

	"github.com/arnaz06/users"
)

type myCustomClaims struct {
    Email string `json:"email"`
    jwt.StandardClaims
}

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

// AuthenticationMiddleware is a function to check a user based on key authentication.
func AuthenticationMiddleware(secretKey string) middleware.KeyAuthValidator {
	return func(key string, c echo.Context) (bool, error) {
		tokenString := c.Request().Header.Get("Authorization")

		splitedString := strings.Split(tokenString, " ")

		if len(splitedString) < 2 {
			return false, users.UnauthorizedErrorf("invalid token format")
		}

		_, err := jwt.ParseWithClaims(splitedString[1], &myCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
		if err != nil {
			return false, users.UnauthorizedErrorf("invalid token: %+v", err)
		}

		return true, nil
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

			if _, ok := err.(users.ConstraintError); ok {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			if _, ok := err.(users.UnauthorizedError); ok {
				lg.Errorln(err.Error())
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}

			switch err {
			case context.DeadlineExceeded, context.Canceled:
				lg.Errorln(err.Error())
				return echo.NewHTTPError(http.StatusRequestTimeout, err.Error())

			case users.ErrNotFound:
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			}

			lg.Errorln(err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
}
