package api

import (
	"context"
	"fmt"
	"github.com/3scale-labs/kamwiel/pkg/domain/api"
	"github.com/3scale-labs/kamwiel/pkg/repositories/kuadrant"
	"strings"
)

type Service interface {
	GetAPI(context.Context, string) (*api.API, error)
	ListAPI(context.Context) (*api.APIs, error)
	GetAPIListState(context.Context) (string, error)
	UpdateAPIListState(context.Context, string) error
}

type service struct {
	kuadrantRepo kuadrant.Repository
}

func NewService(repo kuadrant.Repository) Service {
	return &service{
		kuadrantRepo: repo,
	}
}

func (s *service) GetAPI(ctx context.Context, name string) (*api.API, error) {
	apiName := strings.TrimSpace(name)
	if len(apiName) == 0 {
		return nil, fmt.Errorf("invalid API name")
	}
	api, err := s.kuadrantRepo.GetAPI(ctx, apiName)
	if err != nil {
		return nil, err
	}
	return api, nil
}

func (s *service) ListAPI(ctx context.Context) (*api.APIs, error) {
	apis, err := s.kuadrantRepo.ListAPI(ctx)
	if err != nil {
		return nil, err
	}
	return apis, nil
}

func (s *service) GetAPIListState(ctx context.Context) (string, error) {
	apiListState, err := s.kuadrantRepo.GetAPIListState(ctx)
	if err != nil {
		return "", err
	}
	return apiListState, nil
}

func (s *service) UpdateAPIListState(ctx context.Context, hash string) error {
	apisHash := strings.TrimSpace(hash)
	if len(apisHash) == 0 {
		return fmt.Errorf("invalid hash")
	}
	return s.kuadrantRepo.UpdateAPIListState(ctx, apisHash)
}
