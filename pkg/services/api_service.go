package services

import (
	"fmt"
	"github.com/3scale-labs/kamwiel/pkg/domain/api"
	"github.com/3scale-labs/kamwiel/pkg/repositories"
	"strings"
)

type APIService interface {
	GetAPI(string) (*api.API, error)
}

type apiService struct {
	kuadrantRepo repositories.KuadrantRepository
}

func NewAPIService(kuadrantRepo repositories.KuadrantRepository) APIService {
	return &apiService{
		kuadrantRepo: kuadrantRepo,
	}
}

func (s *apiService) GetAPI(name string) (*api.API, error) {
	apiName := strings.TrimSpace(name)
	if len(apiName) == 0 {
		return nil, fmt.Errorf("invalid API name")
	}
	api, err := s.kuadrantRepo.GetAPI(apiName)
	if err != nil {
		return nil, err
	}
	return api, nil
}
