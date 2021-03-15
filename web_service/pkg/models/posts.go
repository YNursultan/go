package models

import (
	"errors"
)

var ErrNoRecord = errors.New("models: no matching record found")

type Post struct {
	ID          int64
	Title       string
	Description string
	Category    string
	UserId      int64
}
