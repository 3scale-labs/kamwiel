package api

import (
	"fmt"
	"strings"
)

type Repository interface {
	GetAPI(string) (*API, error)
	ListAPI() (*APIs, error)
}

type Service interface {
	GetAPI(string) (*API, error)
	ListAPI() (*APIs, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetAPI(name string) (*API, error) {
	apiName := strings.TrimSpace(name)
	if len(apiName) == 0 {
		return nil, fmt.Errorf("invalid API name")
	}
	api, err := s.repo.GetAPI(apiName)
	if err != nil {
		return nil, err
	}
	return api, nil
}

func (s *service) ListAPI() (*APIs, error) {
	apis, err := s.repo.ListAPI()
	if err != nil {
		return nil, err
	}
	return apis, nil
}
