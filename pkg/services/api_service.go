package services

import (
	"github.com/3scale-labs/kamwiel/pkg/domain/api"
)

var APIService apiServiceInterface = &apiService{}

type apiServiceInterface interface {
	GetAPI(string) (*api.API, error)
}

type apiService struct {}

func (s *apiService) GetAPI(name string) (*api.API, error) {
	dao := &api.API{Name: name}
	if err := dao.Get(); err != nil { return nil, err }
	return dao, nil
}
