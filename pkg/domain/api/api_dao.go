package api

import (
	"context"
	"github.com/3scale-labs/kamwiel/pkg/datasources/kuadrant"
	kctlrv1beta1 "github.com/kuadrant/kuadrant-controller/apis/networking/v1beta1"
	"k8s.io/client-go/rest"
)

var kuadrantClient *rest.RESTClient

const (
	namespace = "kamwiel" // get from cluster config namespace
	resource = "apis"
)

func init() {
	kuadrantClient = kuadrant.Client
}

func (a *API) Get() error {
	api := &kctlrv1beta1.API{}
	if err := kuadrantClient.
		Get().
		Namespace(namespace).
		Resource(resource).
		Name(a.Name).
		Do(context.Background()).
		Into(api);  err != nil { return err }

	a.Spec = api.Spec

	return nil
}
