package api

import (
	"context"
	"github.com/3scale-labs/kamwiel/pkg/datasources/kuadrant"
	kctlrv1beta1 "github.com/kuadrant/kuadrant-controller/apis/networking/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var kuadrantClient client.Client

const namespace = "kamwiel" // get from cluster config namespace

func init() {
	kuadrantClient = kuadrant.Client
}

func (a *API) Get() error {
	api := &kctlrv1beta1.API{}
	if err := kuadrantClient.Get(context.Background(), client.ObjectKey{
		Namespace: namespace,
		Name:      a.Name,
	}, api); err != nil {
		return err
	}
	a.Spec = api.Spec

	return nil
}
