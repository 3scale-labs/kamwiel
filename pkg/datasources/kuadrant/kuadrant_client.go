package kuadrant

import (
	"fmt"
	kctlrv1beta1 "github.com/kuadrant/kuadrant-controller/apis/networking/v1beta1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

var Client *rest.RESTClient

func setupRestClient() *rest.RESTClient {
	configuration, err := rest.InClusterConfig()
	if err != nil {
		panic(fmt.Sprintf("It was impossible to get cluster config: %s", err))
	}

	crdConfig := *configuration
	crdConfig.ContentConfig.GroupVersion = &schema.GroupVersion{Group: kctlrv1beta1.GroupVersion.Group, Version: kctlrv1beta1.GroupVersion.Version}
	crdConfig.APIPath = "/apis"
	crdConfig.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	crdConfig.UserAgent = rest.DefaultKubernetesUserAgent()

	restClient, err := rest.RESTClientFor(&crdConfig)
	if err != nil {
		panic(fmt.Sprintf("It was impossible to setup KuadrantClient: %s", err))
	}

	return restClient
}

func init() {
	if err := kctlrv1beta1.AddToScheme(scheme.Scheme); err != nil { panic(err) }
	Client = setupRestClient()
}
