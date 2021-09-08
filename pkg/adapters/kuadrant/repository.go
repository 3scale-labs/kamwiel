package kuadrant

import (
	"context"
	"github.com/3scale-labs/kamwiel/pkg/domain/api"
	kctlrv1beta1 "github.com/kuadrant/kuadrant-controller/apis/networking/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sort"
)

type Repository interface {
	GetAPI(context.Context, string) (*api.API, error)
	ListAPI(context.Context) (*api.APIs, error)
}

type kuadrantRepository struct {
	client client.Client
}

const kuadrantNamespace = "kamwiel" // TODO: get from cluster config namespace

func NewKuadrantRepository(client client.Client) Repository {
	return &kuadrantRepository{
		client,
	}
}

func (r *kuadrantRepository) GetAPI(ctx context.Context, name string) (*api.API, error) {
	kAPI := &kctlrv1beta1.API{}
	if err := r.client.Get(ctx, client.ObjectKey{
		Namespace: kuadrantNamespace,
		Name:      name,
	}, kAPI); err != nil {
		return nil, err
	}

	return &api.API{Name: name, Spec: *kAPI.Spec.Mappings.OAS}, nil
}

func (r *kuadrantRepository) ListAPI(ctx context.Context) (*api.APIs, error) {
	kAPIs := &kctlrv1beta1.APIList{}
	if err := r.client.List(ctx, kAPIs, &client.ListOptions{Namespace: kuadrantNamespace}); err != nil {
		return nil, err
	}
	apis := make(api.APIs, len(kAPIs.Items))
	for i, kAPI := range kAPIs.Items {
		apis[i] = api.API{Name: kAPI.Name, Spec: *kAPI.Spec.Mappings.OAS}
	}
	// Sorting by Name before returning
	sort.Slice(apis, func(i, j int) bool {
		return apis[i].Name < apis[j].Name
	})

	return &apis, nil
}
