package service

import (
	"TestTask/model"
	"TestTask/pkg/repository"
)

type ApiService struct {
	repo *repository.Repository
}

func NewApiService(repo *repository.Repository) *ApiService {
	return &ApiService{repo: repo}
}

func (s *ApiService) CheckToken(token string) (int, error) {
	return s.repo.CheckToken(token)
}

func (s *ApiService) ClearAudit(userId int) error {
	return s.repo.ClearAudit(userId)
}

func (s *ApiService) GetUserEvents(token string) ([]model.Audit, error) {
	return s.repo.GetUserEvents(token)
}
