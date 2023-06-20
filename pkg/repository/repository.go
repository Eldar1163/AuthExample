package repository

import (
	"TestTask/model"
	"github.com/jmoiron/sqlx"
	"time"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GetUser(username, passwordHash string) (model.User, error)
	WriteUserToken(userId int, token string, expireDate time.Time) error
}

type Repository struct {
	Authorization
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
