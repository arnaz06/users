package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"

	"github.com/arnaz06/users"
)

type userRepo struct {
	db *sql.DB
}

// NewUserRepository is constructor for user repository.
func NewUserRepository(db *sql.DB) users.UserRepository {
	return userRepo{
		db: db,
	}
}

func (r userRepo) Create(ctx context.Context, user users.User) (users.User, error) {
	query := `INSERT users SET id=?, email=?, password=?, address=?, updated_time=?, created_time=?`
	now := time.Now()
	user.CreatedTime = now
	user.UpdatedTime = now
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	_, err := r.db.ExecContext(ctx, query, user.ID, user.Email, user.Password, user.Address, user.UpdatedTime.Unix(), user.CreatedTime.Unix())
	if err != nil {
		return users.User{}, err
	}
	return user, nil
}

func (r userRepo) Get(ctx context.Context, id string) (users.User, error) {
	query := `SELECT id, email, password, address, updated_time, created_time FROM users WHERE id=? AND deleted_time IS NULL`
	row := r.db.QueryRowContext(ctx, query, id)

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
			return users.User{}, users.ErrNotFound
		}
		return users.User{}, err
	}

	res.UpdatedTime = time.Unix(updatedTime, 0)
	res.CreatedTime = time.Unix(createdTime, 0)
	return res, nil
}

func (r userRepo) GetByEmail(ctx context.Context, email string) (users.User, error) {
	query := `SELECT id, email, password, address, updated_time, created_time FROM users WHERE email=? AND deleted_time IS NULL`
	row := r.db.QueryRowContext(ctx, query, email)

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
			return users.User{}, users.ErrNotFound
		}
		return users.User{}, err
	}

	res.UpdatedTime = time.Unix(updatedTime, 0)
	res.CreatedTime = time.Unix(createdTime, 0)
	return res, nil
}

func (r userRepo) Update(ctx context.Context, user users.User) error {
	query := `UPDATE users SET email=?, password=?, address=?, updated_time=? WHERE id=? AND deleted_time IS NULL`
	user.UpdatedTime = time.Now()

	res, err := r.db.ExecContext(ctx, query, user.Email, user.Password, user.Address, user.UpdatedTime.Unix(), user.ID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected != 1 {
		return users.ErrNotFound
	}

	return nil
}

func (r userRepo) Delete(ctx context.Context, id string) error {
	query := `UPDATE users SET deleted_time=? WHERE id=?`
	res, err := r.db.ExecContext(ctx, query, time.Now().Unix(), id)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected != 1 {
		return users.ErrNotFound
	}

	return nil
}
