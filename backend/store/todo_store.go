package store

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/thimc/go-backend/types"

	_ "github.com/lib/pq"
)

type DatabaseStorer interface {
	GetTodos(context.Context) ([]*types.Todo, error)
	GetTodoByID(context.Context, int64) (*types.Todo, error)
	InsertTodo(context.Context, *types.Todo) (*types.Todo, error)
	UpdateTodoByID(context.Context, types.UpdateTodoParams, int64) error
	DeleteTodoByID(context.Context, int64) error
	PatchTodoByID(context.Context, int64, types.UpdateTodoParams) error

	Close() error
}

type PostgreStore struct {
	db *sql.DB
}

func NewPostgreStore(connectionStr string) (*PostgreStore, error) {
	db, err := sql.Open("postgres", connectionStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	store := &PostgreStore{
		db: db,
	}
	if err := store.init(); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *PostgreStore) init() error {
	query := `CREATE TABLE IF NOT EXISTS todo (
		id SERIAL PRIMARY KEY,
		title VARCHAR(100),
		content VARCHAR(1000),
		created TIMESTAMP,
		updated TIMESTAMP,
		created_by SERIAL,
		updated_by SERIAL,
		done BOOLEAN
	);`
	_, err := s.db.Exec(query)

	return err
}

func (s *PostgreStore) Close() error {
	return s.db.Close()
}

func (s *PostgreStore) GetTodos(ctx context.Context) ([]*types.Todo, error) {
	todos := []*types.Todo{}

	rows, err := s.db.QueryContext(ctx, `SELECT * FROM todo`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		todo, err := scanTodo(rows)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

func (s *PostgreStore) GetTodoByID(ctx context.Context, id int64) (*types.Todo, error) {
	var todo *types.Todo

	rows, err := s.db.QueryContext(ctx, `SELECT * FROM todo
								WHERE id = $1 LIMIT 1`, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		todo, err = scanTodo(rows)
		if err != nil {
			return nil, err
		}
	}

	return todo, nil
}

func (s *PostgreStore) InsertTodo(ctx context.Context, t *types.Todo) (*types.Todo, error) {
	query := `INSERT INTO todo(title, content, created, created_by, done)
				VALUES($1, $2, NOW(), $3, $4) RETURNING id;`
	rows, err := s.db.QueryContext(ctx, query, t.Title, t.Content, t.CreatedBy, t.Done)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err := rows.Scan(&t.ID)
		if err != nil {
			return nil, err
		}
	}

	return t, nil
}

func (s *PostgreStore) UpdateTodoByID(ctx context.Context, t types.UpdateTodoParams, id int64) error {
	query := `UPDATE todo
				SET title = $1, content = $2, created = $3, updated = $4, created_by = $5, updated_by = $6, done = $7
				WHERE id = $8`
	_, err := s.db.ExecContext(ctx, query, t.Title, t.Content, t.UpdatedBy, id)

	return err
}

func (s *PostgreStore) DeleteTodoByID(ctx context.Context, id int64) error {
	todo, err := s.GetTodoByID(ctx, id)
	if err != nil {
		return err
	}
	if todo == nil {
		return fmt.Errorf("unknown id: %d", id)
	}
	_, err = s.db.ExecContext(ctx, `DELETE FROM todo
						WHERE id = $1`, id)

	return err
}

func (s *PostgreStore) PatchTodoByID(ctx context.Context, id int64, t types.UpdateTodoParams) error {
	var (
		sb  strings.Builder
		ref = reflect.Indirect(reflect.ValueOf(t))
	)
	for i := 0; i < ref.NumField(); i++ {
		field := ref.Type().Field(i)
		tag := field.Tag.Get("sql")
		if tag == "" {
			continue
		}

		val := ref.Field(i)
		if val.IsNil() {
			continue
		}
		fieldValue := val.Elem().Interface()

		switch val.Elem().Kind() {
		case reflect.Bool:
			sb.WriteString(fmt.Sprintf("%s = %v", tag, fieldValue))
		case reflect.Struct:
			date := fieldValue.(time.Time)
			sb.WriteString(fmt.Sprintf("%s = '%s'", tag, date.Format("2006-01-02 15:04:05.99-07")))
		case reflect.String:
			sb.WriteString(fmt.Sprintf("%s = '%v'", tag, fieldValue))
		case reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64,
			reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64,
			reflect.Float32, reflect.Float64:
			sb.WriteString(fmt.Sprintf("%s = %v", tag, fieldValue))
		}
		sb.WriteString(", ")
	}
	query := fmt.Sprintf("UPDATE todo SET %s WHERE id = %d", sb.String()[:len(sb.String())-2], id)
	_, err := s.db.ExecContext(ctx, query)

	return err
}

func scanTodo(rows *sql.Rows) (*types.Todo, error) {
	todo := new(types.Todo)
	// FIXME: if we return the error here any NULL value in the
	// database will cause the internal SQL scanner to throw an error.
	_ = rows.Scan(
		&todo.ID,
		&todo.Title,
		&todo.Content,
		&todo.Created,
		&todo.Updated,
		&todo.CreatedBy,
		&todo.UpdatedBy,
		&todo.Done)

	return todo, nil
}
