package http_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/arnaz06/users"
	handler "github.com/arnaz06/users/internal/http"
)

func TestErrorMiddleware(t *testing.T) {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)

	mw := handler.ErrorMiddleware()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	t.Run("with not found object", func(t *testing.T) {
		h := func(c echo.Context) error {
			return users.ErrNotFound
		}
		err := mw(h)(c).(*echo.HTTPError)

		require.Error(t, err)
		require.Equal(t, http.StatusNotFound, err.Code)
	})

	t.Run("with constraint error", func(t *testing.T) {
		h := func(c echo.Context) error {
			return users.ConstraintErrorf("this is a constraint error")
		}

		buf := new(bytes.Buffer)
		log.SetOutput(buf)

		err := mw(h)(c).(*echo.HTTPError)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.Code)
	})

	t.Run("with unauthorized user", func(t *testing.T) {
		h := func(c echo.Context) error {
			return users.UnauthorizedErrorf("invalid user")
		}

		buf := new(bytes.Buffer)
		log.SetOutput(buf)

		err := mw(h)(c).(*echo.HTTPError)
		require.Error(t, err)
		require.Equal(t, http.StatusUnauthorized, err.Code)
	})

	t.Run("with unexpected error", func(t *testing.T) {
		h := func(c echo.Context) error {
			return errors.New("unexpected error")
		}

		buf := new(bytes.Buffer)
		log.SetOutput(buf)

		err := mw(h)(c).(*echo.HTTPError)
		require.Error(t, err)
		require.Equal(t, http.StatusInternalServerError, err.Code)
		require.Contains(t, buf.String(), "unexpected error")
	})

}
