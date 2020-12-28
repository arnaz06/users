package user

import (
	"context"

	"github.com/arnaz06/users"
)

type userService struct {
	repo users.UserRepository
}

// NewUserService creates a new user service
func NewUserService(repo users.UserRepository) users.UserService {
	return userService{
		repo: repo,
	}
}

func (s userService) Create(ctx context.Context, user users.User) (users.User, error) {
	return s.repo.Create(ctx, user)
}

func (s userService) Get(ctx context.Context, id string) (users.User, error) {
	return s.repo.Get(ctx, id)
}

func (s userService) Login(ctx context.Context, email, password string) error {
	savedUser, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return err
	}

	err = users.CompareHash(savedUser.Password, password)
	if err != nil {
		return users.UnauthorizedErrorf("invalid password")
	}

	return nil
}

func (s userService) Update(ctx context.Context, user users.User) error {
	return s.repo.Update(ctx, user)
}

func (s userService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
