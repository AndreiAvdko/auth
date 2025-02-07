package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID             int64
	Name           string
	Email          string
	Password       string
	PassworConfirm string
	IsAdmin        string
	CreatedAt      time.Time
	UpdatedAt      sql.NullTime
}
