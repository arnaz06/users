package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/arnaz06/users"
	"github.com/arnaz06/users/mocks"
	"github.com/arnaz06/users/testdata"
	"github.com/arnaz06/users/user"
)

func TestLoginUserService(t *testing.T) {
	var mockUser users.User
	testdata.GoldenJSONUnmarshal(t, "user", &mockUser)
	mockUser.Password = "$2a$04$s06QpvFZ6QeQouJEozOLjeqtUhnCAY307dTgq.aThUH6d8/5VO89u"
	tests := []struct {
		testName      string
		email         string
		password      string
		repo          testdata.FuncCall
		expectedError error
	}{
		{
			testName: "success",
			email:    mockUser.Email,
			password: "secret-123",
			repo: testdata.FuncCall{
				Called: true,
				Input:  []interface{}{mock.Anything, mockUser.Email},
				Output: []interface{}{mockUser, nil},
			},
		},
		{
			testName: "invalid password",
			email:    mockUser.Email,
			password: "invalid-password",
			repo: testdata.FuncCall{
				Called: true,
				Input:  []interface{}{mock.Anything, mockUser.Email},
				Output: []interface{}{mockUser, nil},
			},
			expectedError: users.UnauthorizedErrorf("invalid password"),
		},
		{
			testName: "unexpected error from service",
			email:    mockUser.Email,
			password: "secret-123",
			repo: testdata.FuncCall{
				Called: true,
				Input:  []interface{}{mock.Anything, mockUser.Email},
				Output: []interface{}{users.User{}, errors.New("unexpected error")},
			},
			expectedError: errors.New("unexpected error"),
		},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			mockRepo := new(mocks.UserRepository)

			if test.repo.Called {
				mockRepo.On("GetByEmail", test.repo.Input...).
					Return(test.repo.Output...).Once()
			}

			service := user.NewUserService(mockRepo)
			err := service.Login(context.Background(), test.email, test.password)
			mockRepo.AssertExpectations(t)

			if test.expectedError != nil {
				require.EqualError(t, err, test.expectedError.Error())
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestCreateUserService(t *testing.T) {
	var mockUser users.User
	testdata.GoldenJSONUnmarshal(t, "user", &mockUser)

	tests := []struct {
		testName       string
		input          users.User
		repo           testdata.FuncCall
		expectedResult users.User
		expectedError  error
	}{
		{
			testName: "success",
			input:    mockUser,
			repo: testdata.FuncCall{
				Called: true,
				Input:  []interface{}{mock.Anything, mockUser},
				Output: []interface{}{mockUser, nil},
			},
			expectedResult: mockUser,
		},
		{
			testName: "error from service",
			input:    mockUser,
			repo: testdata.FuncCall{
				Called: true,
				Input:  []interface{}{mock.Anything, mockUser},
				Output: []interface{}{users.User{}, errors.New("unexpected error")},
			},
			expectedError: errors.New("unexpected error"),
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			mockRepo := new(mocks.UserRepository)
			if test.repo.Called {
				mockRepo.On("Create", test.repo.Input...).
					Return(test.repo.Output...).Once()
			}

			service := user.NewUserService(mockRepo)
			res, err := service.Create(context.Background(), test.input)
			mockRepo.AssertExpectations(t)

			if test.expectedError != nil {
				require.EqualError(t, err, test.expectedError.Error())
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.expectedResult, res)
		})
	}
}

func TestUpdateUserService(t *testing.T) {
	var mockUser users.User
	testdata.GoldenJSONUnmarshal(t, "user", &mockUser)

	tests := []struct {
		testName      string
		input         users.User
		repo          testdata.FuncCall
		expectedError error
	}{
		{
			testName: "success",
			input:    mockUser,
			repo: testdata.FuncCall{
				Called: true,
				Input:  []interface{}{mock.Anything, mockUser},
				Output: []interface{}{nil},
			},
		},
		{
			testName: "error from service",
			input:    mockUser,
			repo: testdata.FuncCall{
				Called: true,
				Input:  []interface{}{mock.Anything, mockUser},
				Output: []interface{}{errors.New("unexpected error")},
			},
			expectedError: errors.New("unexpected error"),
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			mockRepo := new(mocks.UserRepository)
			if test.repo.Called {
				mockRepo.On("Update", test.repo.Input...).
					Return(test.repo.Output...).Once()
			}

			service := user.NewUserService(mockRepo)
			err := service.Update(context.Background(), test.input)
			mockRepo.AssertExpectations(t)

			if test.expectedError != nil {
				require.EqualError(t, err, test.expectedError.Error())
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestGetUserService(t *testing.T) {
	var mockUser users.User
	testdata.GoldenJSONUnmarshal(t, "user", &mockUser)

	tests := []struct {
		testName       string
		input          string
		repo           testdata.FuncCall
		expectedResult users.User
		expectedError  error
	}{
		{
			testName: "success",
			input:    mockUser.ID,
			repo: testdata.FuncCall{
				Called: true,
				Input:  []interface{}{mock.Anything, mockUser.ID},
				Output: []interface{}{mockUser, nil},
			},
			expectedResult: mockUser,
		},
		{
			testName: "error from service",
			input:    mockUser.ID,
			repo: testdata.FuncCall{
				Called: true,
				Input:  []interface{}{mock.Anything, mockUser.ID},
				Output: []interface{}{users.User{}, errors.New("unexpected error")},
			},
			expectedError: errors.New("unexpected error"),
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			mockRepo := new(mocks.UserRepository)
			if test.repo.Called {
				mockRepo.On("Get", test.repo.Input...).
					Return(test.repo.Output...).Once()
			}

			service := user.NewUserService(mockRepo)
			res, err := service.Get(context.Background(), test.input)
			mockRepo.AssertExpectations(t)

			if test.expectedError != nil {
				require.EqualError(t, err, test.expectedError.Error())
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.expectedResult, res)
		})
	}
}

func TestDeleteUserService(t *testing.T) {
	var mockUser users.User
	testdata.GoldenJSONUnmarshal(t, "user", &mockUser)

	tests := []struct {
		testName      string
		input         string
		repo          testdata.FuncCall
		expectedError error
	}{
		{
			testName: "success",
			input:    mockUser.ID,
			repo: testdata.FuncCall{
				Called: true,
				Input:  []interface{}{mock.Anything, mockUser.ID},
				Output: []interface{}{nil},
			},
		},
		{
			testName: "error from service",
			input:    mockUser.ID,
			repo: testdata.FuncCall{
				Called: true,
				Input:  []interface{}{mock.Anything, mockUser.ID},
				Output: []interface{}{errors.New("unexpected error")},
			},
			expectedError: errors.New("unexpected error"),
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			mockRepo := new(mocks.UserRepository)
			if test.repo.Called {
				mockRepo.On("Delete", test.repo.Input...).
					Return(test.repo.Output...).Once()
			}

			service := user.NewUserService(mockRepo)
			err := service.Delete(context.Background(), test.input)
			mockRepo.AssertExpectations(t)

			if test.expectedError != nil {
				require.EqualError(t, err, test.expectedError.Error())
				return
			}

			require.NoError(t, err)
		})
	}
}
