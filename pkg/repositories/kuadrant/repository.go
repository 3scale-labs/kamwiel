package kuadrant

import (
	"context"
	"github.com/3scale-labs/kamwiel/pkg/domain/api"
	kctlrv1beta1 "github.com/kuadrant/kuadrant-controller/apis/networking/v1beta1"
	v1 "k8s.io/api/core/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sort"
)

type Repository interface {
	GetAPI(context.Context, string) (*api.API, error)
	ListAPI(context.Context) (*api.APIs, error)
	GetAPIListState(context.Context) (string, error)
	UpdateAPIListState(context.Context, string) error
}

type kuadrantRepository struct {
	client client.Client
}

const kuadrantNamespace = "kamwiel" // TODO: get from cluster config namespace
const apiListStatusConfigMap = "kamwiel-api-list-status"

func NewRepository(client client.Client) Repository {
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

func (r *kuadrantRepository) GetAPIListState(ctx context.Context) (string, error) {
	apiListStatus, err := r.getApiListStatus(ctx)
	if err != nil {
		return "", err
	}
	return apiListStatus.Data["fresh"], nil
}

func (r *kuadrantRepository) UpdateAPIListState(ctx context.Context, newHash string) error {
	apiListStatus, err := r.getApiListStatus(ctx)
	if err != nil {
		return err
	}

	if apiListStatus.Data["hash"] == newHash && apiListStatus.Data["fresh"] != "false" {
		apiListStatus.Data["fresh"] = "false"
		if updateErr := r.client.Update(ctx, apiListStatus); updateErr != nil {
			return updateErr
		}
	}
	return nil
}

func (r *kuadrantRepository) getApiListStatus(ctx context.Context) (*v1.ConfigMap, error) {
	apiListStatus := &v1.ConfigMap{}
	namespacedName := k8stypes.NamespacedName{Namespace: kuadrantNamespace, Name: apiListStatusConfigMap}
	if err := r.client.Get(ctx, namespacedName, apiListStatus); err != nil {
		return nil, err
	}
	return apiListStatus, nil
}
