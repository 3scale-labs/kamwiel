package kuadrant

import (
	"context"
	"github.com/3scale-labs/kamwiel/pkg/domain/api"
	kctlrv1beta1 "github.com/kuadrant/kuadrant-controller/apis/networking/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type kuadrantRepository struct{}

const namespace = "kamwiel" // TODO: get from cluster config namespace

var kuadrantClient client.Client

func init() {
	kuadrantClient = Client
}

func NewKuadrantRepository() api.Repository {
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

	return &api.API{Name: name, Spec: *kAPI.Spec.Mappings.OAS}, nil
}

func (r *kuadrantRepository) ListAPI() (*api.APIs, error) {
	kAPIs := &kctlrv1beta1.APIList{}
	if err := kuadrantClient.List(context.Background(), kAPIs, &client.ListOptions{
		Namespace: namespace,
	}); err != nil {
		return nil, err
	}

	apis := make(api.APIs, len(kAPIs.Items))
	for _, kAPI := range kAPIs.Items {
		apis = append(apis, api.API{Name: kAPI.Name, Spec: *kAPI.Spec.Mappings.OAS})
	}

	return &apis, nil
}