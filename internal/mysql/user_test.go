package mysql_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/arnaz06/users"
	"github.com/arnaz06/users/internal/mysql"
	"github.com/arnaz06/users/testdata"
)

type userSuite struct {
	mysqlSuite
}

func TestUserSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipped for short testing")
	}
	suite.Run(t, new(userSuite))
}

func (u *userSuite) SetupTest() {
	_, err := u.db.Exec("TRUNCATE users")
	require.NoError(u.T(), err)
}

func (u *userSuite) seedUser(user users.User) {
	query := `INSERT users SET id=?, email=?, password=?, address=?, updated_time=?, created_time=?`
	user.CreatedTime = time.Now()
	user.UpdatedTime = user.CreatedTime
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	_, err := u.db.Exec(query, user.ID, user.Email, user.Password, user.Address, user.UpdatedTime.Unix(), user.CreatedTime.Unix())
	require.NoError(u.T(), err)
}

func (u *userSuite) getUser(id string) users.User {
	query := `SELECT id, email, password, address, updated_time, created_time FROM users WHERE id=? AND deleted_time IS NULL`
	row := u.db.QueryRowContext(context.Background(), query, id)

	var res users.User
	updatedTime := int64(0)
	createdTime := int64(0)
	err := row.Scan(
		&res.ID,
		&res.Email,
		&res.Password,
		&res.Address,
		&updatedTime,
		&createdTime,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return users.User{}
		}
		require.NoError(u.T(), err)
	}

	res.UpdatedTime = time.Unix(updatedTime, 0)
	res.CreatedTime = time.Unix(createdTime, 0)
	return res
}

func (u *userSuite) TestCreateUser() {
	var mockUser users.User
	testdata.GoldenJSONUnmarshal(u.T(), "user", &mockUser)
	repo := mysql.NewUserRepository(u.db)

	u.T().Run("success", func(t *testing.T) {
		res, err := repo.Create(context.Background(), mockUser)
		require.NoError(t, err)

		tz := time.UTC
		if loc, err := time.LoadLocation("Asia/Jakarta"); err == nil {
			tz = loc
		}

		res.UpdatedTime = res.UpdatedTime.In(tz)
		res.CreatedTime = res.CreatedTime.In(tz)

		inserted := u.getUser(mockUser.ID)
		inserted.UpdatedTime = res.UpdatedTime.In(tz)
		inserted.CreatedTime = res.CreatedTime.In(tz)
		require.Equal(t, inserted, res)
	})
}

func (u *userSuite) TestGetUser() {
	var mockUser users.User
	testdata.GoldenJSONUnmarshal(u.T(), "user", &mockUser)
	u.seedUser(mockUser)

	tests := []struct {
		testName       string
		input          string
		expectedResult users.User
		expectedError  error
	}{
		{
			testName:       "success",
			input:          mockUser.ID,
			expectedResult: mockUser,
		},
		{
			testName:      "error not found",
			input:         "user-404",
			expectedError: users.ErrNotFound,
		},
	}

	for _, test := range tests {
		u.T().Run(test.testName, func(t *testing.T) {
			repo := mysql.NewUserRepository(u.db)
			res, err := repo.Get(context.Background(), test.input)

			if test.expectedError != nil {
				require.EqualError(t, err, test.expectedError.Error())
				return
			}

			test.expectedResult.UpdatedTime = time.Time{}
			test.expectedResult.CreatedTime = time.Time{}
			res.UpdatedTime = time.Time{}
			res.CreatedTime = time.Time{}
			require.Equal(t, test.expectedResult, res)
		})
	}
}

func (u *userSuite) TestGetUserByEmail() {
	var mockUser users.User
	testdata.GoldenJSONUnmarshal(u.T(), "user", &mockUser)
	u.seedUser(mockUser)

	tests := []struct {
		testName       string
		input          string
		expectedResult users.User
		expectedError  error
	}{
		{
			testName:       "success",
			input:          mockUser.Email,
			expectedResult: mockUser,
		},
		{
			testName:      "error not found",
			input:         "email-404",
			expectedError: users.ErrNotFound,
		},
	}

	for _, test := range tests {
		u.T().Run(test.testName, func(t *testing.T) {
			repo := mysql.NewUserRepository(u.db)
			res, err := repo.GetByEmail(context.Background(), test.input)

			if test.expectedError != nil {
				require.EqualError(t, err, test.expectedError.Error())
				return
			}
			require.NoError(t, err)

			test.expectedResult.UpdatedTime = time.Time{}
			test.expectedResult.CreatedTime = time.Time{}
			res.UpdatedTime = time.Time{}
			res.CreatedTime = time.Time{}
			require.Equal(t, test.expectedResult, res)
		})
	}
}

func (u *userSuite) TestUpdateUser() {
	var mockUser users.User
	testdata.GoldenJSONUnmarshal(u.T(), "user", &mockUser)
	u.seedUser(mockUser)

	updatedUser := users.User(mockUser)
	updatedUser.Address = "updated address"

	userNotFound := users.User(mockUser)
	userNotFound.ID = "404"
	tests := []struct {
		testName      string
		input         users.User
		expectedError error
	}{
		{
			testName: "success",
			input:    updatedUser,
		},
		{
			testName:      "success",
			input:         userNotFound,
			expectedError: users.ErrNotFound,
		},
	}

	for _, test := range tests {
		u.T().Run(test.testName, func(t *testing.T) {
			repo := mysql.NewUserRepository(u.db)

			err := repo.Update(context.Background(), test.input)

			if test.expectedError != nil {
				require.EqualError(t, err, test.expectedError.Error())
				return
			}
			require.NoError(t, err)
		})
	}
}

func (u *userSuite) TestDeleteUser() {
	var mockUser users.User
	testdata.GoldenJSONUnmarshal(u.T(), "user", &mockUser)
	u.seedUser(mockUser)
	tests := []struct {
		testName      string
		input         string
		expectedError error
	}{
		{
			testName: "success",
			input:    mockUser.ID,
		},
		{
			testName:      "user not found",
			input:         "404",
			expectedError: users.ErrNotFound,
		},
	}

	for _, test := range tests {
		u.T().Run(test.testName, func(t *testing.T) {
			repo := mysql.NewUserRepository(u.db)

			err := repo.Delete(context.Background(), test.input)

			if test.expectedError != nil {
				require.EqualError(t, err, test.expectedError.Error())
				return
			}
			require.NoError(t, err)
		})
	}
}
