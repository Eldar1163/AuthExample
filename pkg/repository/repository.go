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
	WrongPasswordEnterCnt(userId int) (int, error)
	WriteEvent(username string, event int) error
	GetUserByUsername(username string) (model.User, error)
}

type Api interface {
	CheckToken(token string) (int, error)
	ClearAudit(userId int) error
	GetUserEvents(token string) ([]model.Audit, error)
}

type Repository struct {
	Authorization
	Api
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Api:           NewApiPostgres(db),
	}
}
