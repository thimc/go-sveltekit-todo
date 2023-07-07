package store

import (
	"context"
	"database/sql"

	"github.com/thimc/go-backend/types"
)

type UserStorer interface {
	GetUserByID(context.Context, int64) (*types.User, error)
	GetUserByEmail(context.Context, string) (*types.User, error)
	DeleteUserByID(context.Context, int64) error
	CreateUser(context.Context, *types.User) (*types.User, error)
	UpdateUser(context.Context, int64, *types.User) error

	Close() error
	init() error
}

type PostgreUserStore struct {
	db *sql.DB
}

func NewPostgreUserStore(s *PostgreStore) (*PostgreUserStore, error) {
	store := &PostgreUserStore{
		db: s.db,
	}
	err := store.init()
	return store, err
}

func (s *PostgreUserStore) init() error {
	query := `
		CREATE TABLE IF NOT EXISTS todo_user (
		id SERIAL PRIMARY KEY,
		email VARCHAR(100),
		password VARCHAR(100)
	);`
	_, err := s.db.Query(query)

	return err
}

func (s *PostgreUserStore) Close() error {
	return s.db.Close()
}

func (s *PostgreUserStore) GetUserByID(ctx context.Context, id int64) (*types.User, error) {
	var user *types.User

	rows, err := s.db.QueryContext(ctx, `SELECT * FROM todo_user
								WHERE id = $1 LIMIT 1`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		user, err = scanUser(rows)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (s *PostgreUserStore) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	var user *types.User

	rows, err := s.db.QueryContext(ctx, `SELECT * FROM todo_user
								WHERE email = $1 LIMIT 1`, email)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		user, err = scanUser(rows)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (s *PostgreUserStore) CreateUser(ctx context.Context, u *types.User) (*types.User, error) {
	query := `INSERT INTO todo_user(email, encrypted_password)
				VALUES($1, $2)`
	result, err := s.db.ExecContext(ctx, query, u.Email, u.EncryptedPassword)
	if err != nil {
		return nil, err
	}

	insertedId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	u.ID = int(insertedId)

	return u, nil
}

func (s *PostgreUserStore) UpdateUser(ctx context.Context, id int64, u *types.User) error {
	return nil
}

func (s *PostgreUserStore) DeleteUserByID(ctx context.Context, id int64) error {
	_, err := s.db.QueryContext(ctx, `DELETE FROM todo_user WHERE id = $1`, id)

	return err
}

func scanUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)
	_ = rows.Scan(&user.ID,
		&user.Email,
		&user.EncryptedPassword)
	return user, nil
}
