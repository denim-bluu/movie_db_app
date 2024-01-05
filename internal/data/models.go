package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Movies MovieModel
}

func NwewModels(db *sql.DB) *Models {
	return &Models{
		Movies: MovieModel{DB: db},
	}
}
