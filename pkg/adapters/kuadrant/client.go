package kuadrant

import (
	"context"
	"fmt"
	kctlrv1beta1 "github.com/kuadrant/kuadrant-controller/apis/networking/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	utilRuntime "k8s.io/apimachinery/pkg/util/runtime"
	clientGoScheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
	"sync"
)

var (
	clientOnce    sync.Once
	k8sClient     client.Client
	RuntimeScheme = runtime.NewScheme()
)

func newClientOrDie() client.Client {
	utilRuntime.Must(clientGoScheme.AddToScheme(RuntimeScheme))
	utilRuntime.Must(kctlrv1beta1.AddToScheme(RuntimeScheme))

	configuration := config.GetConfigOrDie()

	customCache, err := cache.New(configuration, cache.Options{Scheme: RuntimeScheme})
	if err != nil {
		panic(fmt.Sprintf("It was impossible to create a new cache: %s", err))
	}

	go func() {
		err := customCache.Start(context.Background())
		if err != nil {
			panic(fmt.Sprintf("An error occured while starting the cache: %s", err))
		}
	}()

	k8sClient, err = cluster.DefaultNewClient(customCache, configuration, client.Options{Scheme: RuntimeScheme})

	if err != nil {
		panic(fmt.Sprintf("It was impossible to setup KuadrantClient: %s", err))
	}
	return k8sClient
}

func GetClient() client.Client {
	clientOnce.Do(func() { // <-- atomic, does not allow repeating
		k8sClient = newClientOrDie() // <-- thread safe
	})
	return k8sClient
}
