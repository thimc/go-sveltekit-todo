package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/thimc/go-svelte-todo/backend/types"
)

type UserStorer interface {
	GetUsers(context.Context) ([]*types.User, error)
	GetUserByID(context.Context, int64) (*types.User, error)
	GetUserByEmail(context.Context, string) (*types.User, error)
	DeleteUserByID(context.Context, int64) error
	CreateUser(context.Context, *types.User) (*types.User, error)
	UpdateUserPasswordByID(context.Context, string, int64) error

	init() error
}

type PostgreUserStore struct {
	db *sql.DB
}

func NewPostgreUserStore(s *PostgreTodoStore) (*PostgreUserStore, error) {
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
		email VARCHAR(100) UNIQUE,
		encrypted_password VARCHAR(100)
	);`
	_, err := s.db.Query(query)

	return err
}

func (s *PostgreUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	users := []*types.User{}

	rows, err := s.db.QueryContext(ctx, `SELECT * FROM todo_user`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (s *PostgreUserStore) GetUserByID(ctx context.Context, id int64) (*types.User, error) {
	var user *types.User

	rows, err := s.db.QueryContext(ctx, `SELECT * FROM todo_user WHERE id = $1 LIMIT 1`, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		user, err = scanUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if user == nil {
		return nil, fmt.Errorf("unknown ID: %d", id)
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
	user, err := s.GetUserByEmail(ctx, u.Email)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return nil, fmt.Errorf("user exists already")
	}

	query := `INSERT INTO todo_user(email, encrypted_password)
				VALUES($1, $2) RETURNING id;`
	rows, err := s.db.QueryContext(ctx, query, u.Email, u.EncryptedPassword)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err := rows.Scan(&u.ID)
		if err != nil {
			return nil, err
		}
	}

	return u, nil
}

func (s *PostgreUserStore) UpdateUserPasswordByID(ctx context.Context, password string, id int64) error {
	query := `UPDATE todo_user SET encrypted_password = $1
				WHERE id = $2`;
	res, err := s.db.ExecContext(ctx, query, password, id)
	if err !=  nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("unknonw ID: %d", id)
	}

	return nil
}


func (s *PostgreUserStore) DeleteUserByID(ctx context.Context, id int64) error {
	_, err := s.db.QueryContext(ctx, `DELETE FROM todo_user WHERE id = $1`, id)

	return err
}

func scanUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)
	_ = rows.Scan(&user.ID, &user.Email, &user.EncryptedPassword)

	return user, nil
}
