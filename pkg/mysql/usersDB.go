package mysql

import (
	"database/sql"
	"errors"

	"reg/pkg/types"

	_ "github.com/go-sql-driver/mysql"
)

// Users - Определяем тип который обертывает пул подключения sql.DB
type Users struct {
	db *sql.DB
}

func New(db *sql.DB) *Users {
	return &Users{
		db: db,
	}
}

// Insert - Метод для создания новой заметки в базе дынных.
func (m *Users) Insert(user *types.User) (int, error) {
	stmt := `INSERT INTO profiles (login, password, created) VALUES(?, ?, ?)`
	result, err := m.db.Exec(stmt, user.Login, user.Password, user.Created)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Get - Метод для возвращения данных заметки по её идентификатору ID.
func (m *Users) Get(id int) (*types.User, error) {

	stmt := `SELECT login, password, created FROM profiles
    WHERE id = ?`

	row := m.db.QueryRow(stmt, id)

	s := &types.User{}

	err := row.Scan(&s.Login, &s.Password, &s.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, types.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

// Latest - Метод возвращает 10 наиболее часто используемые заметки.
func (m *Users) Latest() ([]*types.User, error) {
	return nil, nil
}

func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
