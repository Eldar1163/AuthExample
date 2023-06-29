package service

import (
	"TestTask/model"
	"TestTask/pkg/repository"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GenerateToken(userId int) (string, error)
	CheckUserCredentials(username, password string) (model.User, error)
	GetBadAuthAttemptsCnt(userId int) (int, error)
	BlockUser(username string) error
}

type Api interface {
	CheckToken(token string) (int, error)
	ClearAudit(userId int) error
	GetUserEvents(token string) ([]model.Audit, error)
}

type Service struct {
	Authorization
	Api
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos),
		Api:           NewApiService(repos),
	}
}
