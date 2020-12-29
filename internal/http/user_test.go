package http_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/arnaz06/users"
	"github.com/arnaz06/users/internal"
	handler "github.com/arnaz06/users/internal/http"
	"github.com/arnaz06/users/mocks"
	"github.com/arnaz06/users/testdata"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func getEchoServer() *echo.Echo {
	e := echo.New()
	e.Validator = internal.NewValidator()
	e.Use(handler.ErrorMiddleware())
	return e
}

func TestCreateUserHandler(t *testing.T) {
	userJSON := testdata.GetGolden(t, "user")
	var mockUser users.User
	testdata.GoldenJSONUnmarshal(t, "user", &mockUser)

	missingEmail := users.User(mockUser)
	missingEmail.Email = ""
	missingEMailJSON, err := json.Marshal(missingEmail)
	require.NoError(t, err)

	tests := []struct {
		testName       string
		input          []byte
		service        testdata.FuncCall
		expectedStatus int
	}{
		{
			testName: "success",
			input:    userJSON,
			service: testdata.FuncCall{
				Called: true,
				Input:  []interface{}{mock.Anything, mock.AnythingOfType("users.User")},
				Output: []interface{}{mockUser, nil},
			},
			expectedStatus: http.StatusCreated,
		},
		{
			testName: "with invalid request body",
			input:    []byte(`invalid body`),
			service: testdata.FuncCall{
				Called: false,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			testName: "error validator",
			input:    missingEMailJSON,
			service: testdata.FuncCall{
				Called: false,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			testName: "with unexpected error from service",
			input:    userJSON,
			service: testdata.FuncCall{
				Called: true,
				Input:  []interface{}{mock.Anything, mock.AnythingOfType("users.User")},
				Output: []interface{}{users.User{}, errors.New("unexpected error")},
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	e := getEchoServer()
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			mockService := new(mocks.UserService)
			if test.service.Called {
				mockService.On("Create", test.service.Input...).
					Return(test.service.Output...).Once()
			}

			req := httptest.NewRequest(echo.POST, "/user", strings.NewReader(string(test.input)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			handler.AddUserHandler(e, mockService, "secret", time.Duration(3600))
			e.ServeHTTP(rec, req)

			mockService.AssertExpectations(t)

			require.Equal(t, test.expectedStatus, rec.Code)
		})
	}
}

func TestUpdateUserHandler(t *testing.T) {
	userJSON := testdata.GetGolden(t, "user")
	var mockUser users.User
	testdata.GoldenJSONUnmarshal(t, "user", &mockUser)

	missingEmail := users.User(mockUser)
	missingEmail.Email = ""
	missingEMailJSON, err := json.Marshal(missingEmail)
	require.NoError(t, err)

	tests := []struct {
		testName       string
		input          []byte
		service        testdata.FuncCall
		expectedStatus int
	}{
		{
			testName: "success",
			input:    userJSON,
			service: testdata.FuncCall{
				Called: true,
				Input:  []interface{}{mock.Anything, mock.AnythingOfType("users.User")},
				Output: []interface{}{nil},
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			testName: "with invalid request body",
			input:    []byte(`invalid body`),
			service: testdata.FuncCall{
				Called: false,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			testName: "error validator",
			input:    missingEMailJSON,
			service: testdata.FuncCall{
				Called: false,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			testName: "with unexpected error from service",
			input:    userJSON,
			service: testdata.FuncCall{
				Called: true,
				Input:  []interface{}{mock.Anything, mock.AnythingOfType("users.User")},
				Output: []interface{}{errors.New("unexpected error")},
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	e := getEchoServer()
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			mockService := new(mocks.UserService)
			if test.service.Called {
				mockService.On("Update", test.service.Input...).
					Return(test.service.Output...).Once()
			}

			req := httptest.NewRequest(echo.PUT, "/user/123", strings.NewReader(string(test.input)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			handler.AddUserHandler(e, mockService, "secret", time.Duration(3600))
			e.ServeHTTP(rec, req)

			mockService.AssertExpectations(t)

			require.Equal(t, test.expectedStatus, rec.Code)
		})
	}
}

func TestDeleteUserHandler(t *testing.T) {
	var mockUser users.User
	testdata.GoldenJSONUnmarshal(t, "user", &mockUser)

	tests := []struct {
		testName       string
		input          string
		service        testdata.FuncCall
		expectedStatus int
	}{
		{
			testName: "success",
			input:    mockUser.ID,
			service: testdata.FuncCall{
				Called: true,
				Input:  []interface{}{mock.Anything, mockUser.ID},
				Output: []interface{}{nil},
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			testName: "with unexpected error from service",
			input:    mockUser.ID,
			service: testdata.FuncCall{
				Called: true,
				Input:  []interface{}{mock.Anything, mockUser.ID},
				Output: []interface{}{errors.New("unexpected error")},
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	e := getEchoServer()
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			mockService := new(mocks.UserService)
			if test.service.Called {
				mockService.On("Delete", test.service.Input...).
					Return(test.service.Output...).Once()
			}

			req := httptest.NewRequest(echo.DELETE, "/user/123", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			handler.AddUserHandler(e, mockService, "secret", time.Duration(3600))
			e.ServeHTTP(rec, req)

			mockService.AssertExpectations(t)

			require.Equal(t, test.expectedStatus, rec.Code)
		})
	}
}
func TestGetUserHandler(t *testing.T) {
	var mockUser users.User
	testdata.GoldenJSONUnmarshal(t, "user", &mockUser)

	tests := []struct {
		testName       string
		input          string
		service        testdata.FuncCall
		expectedStatus int
	}{
		{
			testName: "success",
			input:    mockUser.ID,
			service: testdata.FuncCall{
				Called: true,
				Input:  []interface{}{mock.Anything, mockUser.ID},
				Output: []interface{}{mockUser, nil},
			},
			expectedStatus: http.StatusOK,
		},
		{
			testName: "with unexpected error from service",
			input:    mockUser.ID,
			service: testdata.FuncCall{
				Called: true,
				Input:  []interface{}{mock.Anything, mockUser.ID},
				Output: []interface{}{users.User{}, errors.New("unexpected error")},
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	e := getEchoServer()
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			mockService := new(mocks.UserService)
			if test.service.Called {
				mockService.On("Get", test.service.Input...).
					Return(test.service.Output...).Once()
			}

			req := httptest.NewRequest(echo.GET, "/user/123", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			handler.AddUserHandler(e, mockService, "secret", time.Duration(3600))
			e.ServeHTTP(rec, req)

			mockService.AssertExpectations(t)

			require.Equal(t, test.expectedStatus, rec.Code)
		})
	}
}

func TestLoginUserHandler(t *testing.T) {
	userJSON := testdata.GetGolden(t, "user")
	var mockUser users.User
	testdata.GoldenJSONUnmarshal(t, "user", &mockUser)

	tests := []struct {
		testName       string
		input          []byte
		service        testdata.FuncCall
		expectedStatus int
	}{
		{
			testName: "success",
			input:    userJSON,
			service: testdata.FuncCall{
				Called: true,
				Input:  []interface{}{mock.Anything, mockUser.Email, mockUser.Password},
				Output: []interface{}{nil},
			},
			expectedStatus: http.StatusOK,
		},
		{
			testName: "with invalid request body",
			input:    []byte(`invalid body`),
			service: testdata.FuncCall{
				Called: false,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			testName: "with unexpected error from service",
			input:    userJSON,
			service: testdata.FuncCall{
				Called: true,
				Input:  []interface{}{mock.Anything, mockUser.Email, mockUser.Password},
				Output: []interface{}{errors.New("unexpected error")},
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	e := getEchoServer()
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			mockService := new(mocks.UserService)
			if test.service.Called {
				mockService.On("Login", test.service.Input...).
					Return(test.service.Output...).Once()
			}

			req := httptest.NewRequest(echo.POST, "/user/login", strings.NewReader(string(test.input)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			handler.AddUserHandler(e, mockService, "secret", time.Duration(3600))
			e.ServeHTTP(rec, req)

			mockService.AssertExpectations(t)

			require.Equal(t, test.expectedStatus, rec.Code)
		})
	}
}
