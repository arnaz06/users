package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/arnaz06/users"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type userHandler struct {
	service     users.UserService
	secretKey   string
	expiresTime time.Duration
}

// AddUserHandler adds the user handler.
func AddUserHandler(e *echo.Echo, service users.UserService, secretKey string, expiresTime time.Duration) {
	if service == nil {
		panic("http: nil users service")
	}

	handler := &userHandler{
		service:     service,
		secretKey:   secretKey,
		expiresTime: expiresTime,
	}

	e.POST("/user", handler.create)
	e.GET("/user/:userId", handler.get)
	e.POST("/user/login", handler.login)
	e.PUT("/user/:userId", handler.update)
	e.DELETE("/user/:userId", handler.delete)
}

func (h userHandler) create(c echo.Context) error {
	var input users.User
	if err := c.Bind(&input); err != nil {
		return users.ConstraintErrorf("%s", err)
	}

	if err := c.Validate(input); err != nil {
		return users.ConstraintErrorf("error validating user: %+v", err)
	}

	hashedPassword, err := users.EncodeString(input.Password)
	if err != nil {
		return err
	}
	input.Password = hashedPassword

	res, err := h.service.Create(c.Request().Context(), input)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, res)
}

func (h userHandler) get(c echo.Context) error {
	res, err := h.service.Get(c.Request().Context(), c.Param("userId"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (h userHandler) login(c echo.Context) error {
	var input users.User
	if err := c.Bind(&input); err != nil {
		return users.ConstraintErrorf("%s", err)
	}

	err := h.service.Login(c.Request().Context(), input.Email, input.Password)
	if err != nil {
		return err
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), jwt.StandardClaims{
		ExpiresAt: time.Now().Add(h.expiresTime).Unix(),
		Audience:  c.Param("email"),
	})

	tokenString, err := token.SignedString([]byte(h.secretKey))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"token": fmt.Sprintf("Bearer %s", tokenString)})
}

func (h userHandler) update(c echo.Context) error {
	var input users.User
	if err := c.Bind(&input); err != nil {
		return users.ConstraintErrorf("%s", err)
	}

	if err := c.Validate(input); err != nil {
		return users.ConstraintErrorf("error validating user: %+v", err)
	}

	hashedPassword, err := users.EncodeString(input.Password)
	if err != nil {
		return err
	}

	input.Password = hashedPassword
	input.ID = c.Param("userId")

	err = h.service.Update(c.Request().Context(), input)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h userHandler) delete(c echo.Context) error {
	err := h.service.Delete(c.Request().Context(), c.Param("userId"))
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
