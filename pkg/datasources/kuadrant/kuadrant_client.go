package kuadrant

import (
	"fmt"
	kctlrv1beta1 "github.com/kuadrant/kuadrant-controller/apis/networking/v1beta1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

var Client client.Client

func init() {
	if err := kctlrv1beta1.AddToScheme(scheme.Scheme); err != nil { panic(err) }

	configuration, err := config.GetConfig()
	if err != nil {
		panic(fmt.Sprintf("It was impossible to get cluster config: %s", err))
	}

	Client, err = client.New(configuration, client.Options{Scheme: scheme.Scheme})
	if err != nil {
		panic(fmt.Sprintf("It was impossible to setup KuadrantClient: %s", err))
	}
}
