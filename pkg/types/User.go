package types

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: подходящей записи не найдено")

type User struct {
	Login    string    `json:"email"`
	Password string    `json:"password"`
	Created  time.Time `json:"-"`
}
