package repositories

import (
	"context"
	"github.com/3scale-labs/kamwiel/pkg/clients/kuadrant"
	"github.com/3scale-labs/kamwiel/pkg/domain/api"
	kctlrv1beta1 "github.com/kuadrant/kuadrant-controller/apis/networking/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type KuadrantRepository interface {
	GetAPI(string) (*api.API, error)
}

type kuadrantRepository struct{}

const namespace = "kamwiel" // get from cluster config namespace

var kuadrantClient client.Client

func init() {
	kuadrantClient = kuadrant.Client
}

func NewKuadrantRepository() KuadrantRepository {
	return &kuadrantRepository{}
}

func (r *kuadrantRepository) GetAPI(name string) (*api.API, error) {
	kAPI := &kctlrv1beta1.API{}
	if err := kuadrantClient.Get(context.Background(), client.ObjectKey{
		Namespace: namespace,
		Name:      name,
	}, kAPI); err != nil {
		return nil, err
	}

	return &api.API{Name: name, Spec: kAPI.Spec}, nil
}
