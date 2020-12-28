package users

import (
	"context"
	"time"
)

// User is the struct represent the user's data
type User struct {
	ID          string    `json:"id"`
	Email       string    `json:"email" validate:"required"`
	Address     string    `json:"address"`
	Password    string    `json:"password" validate:"required"`
	CreatedTime time.Time `json:"created_time"`
	UpdatedTime time.Time `json:"updated_time"`
}

// UserRepository is interface of user repository.
type UserRepository interface {
	Create(ctx context.Context, user User) (User, error)
	Get(ctx context.Context, id string) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, id string) error
}

// UserService is interface of user service.
type UserService interface {
	Create(ctx context.Context, user User) (User, error)
	Get(ctx context.Context, id string) (User, error)
	Login(ctx context.Context, email, password string)  error
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, id string) error
}
