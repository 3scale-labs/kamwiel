// +build unit

package controllers

import (
	"context"
	"fmt"
	apirepo "github.com/3scale-labs/kamwiel/pkg/repositories/kuadrant"
	apiservice "github.com/3scale-labs/kamwiel/pkg/services/api"
	kctlrv1beta1 "github.com/kuadrant/kuadrant-controller/apis/networking/v1beta1"
	"gotest.tools/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

type apiControllerTestData struct {
	fakeClient    client.WithWatch
	apiList       *kctlrv1beta1.APIList
	apiListStatus *v1.ConfigMap
}

func buildApiControllerTestData(apiListData, apiListStatusData map[string]string) *apiControllerTestData {
	scheme := runtime.NewScheme()
	_ = kctlrv1beta1.AddToScheme(scheme)
	_ = v1.AddToScheme(scheme)

	apiListStatus := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: "kamwiel-api-list-status", Namespace: "kamwiel"},
		Data:       apiListStatusData,
	}

	oasStringTemplate := "{\"openapi\":\"3.0.0\",\"info\":{\"title\":\"%s\"}}"

	var apiItems []kctlrv1beta1.API
	for k, v := range apiListData {
		spec := fmt.Sprintf(oasStringTemplate, v)
		apiItems = append(apiItems, kctlrv1beta1.API{
			ObjectMeta: metav1.ObjectMeta{
				Name:      k,
				Namespace: "kamwiel",
			},
			Spec: kctlrv1beta1.APISpec{
				Mappings: kctlrv1beta1.APIMappings{
					OAS: &spec,
				},
			},
		},
		)
	}
	apiList := &kctlrv1beta1.APIList{Items: apiItems}

	fakeClient := fake.NewClientBuilder().
		WithScheme(scheme).
		WithLists(apiList).
		WithObjects(apiListStatus).
		Build()

	return &apiControllerTestData{fakeClient, apiList, apiListStatus}
}

func TestConstants(t *testing.T) {
	assert.Check(t, "kamwiel-api-list-status" == apiListStatusConfigMap)
	assert.Check(t, "hash" == cMHash)
	assert.Check(t, "fresh" == cMFresh)
	assert.Check(t, "payload" == cMPayload)
}

func TestAPIControllerReconcileUpdateStatus(t *testing.T) {
	controllerTestData := buildApiControllerTestData(
		map[string]string{
			"dogs-api": "Best friend API",
		},
		map[string]string{
			"hash":    "someHashPreviouslyStored",
			"payload": "aJsonString",
			"fresh":   "false",
		},
	)

	cfMaps := make([]v1.ConfigMap, 1)
	cfMaps = append(cfMaps, *controllerTestData.apiListStatus)

	watchConfigMaps, _ := controllerTestData.fakeClient.Watch(
		context.Background(),
		&v1.ConfigMapList{Items: cfMaps},
		&client.ListOptions{Namespace: "kamwiel"})

	resultChan := watchConfigMaps.ResultChan()

	result, err := (&APIReconciler{
		controllerTestData.fakeClient,
		apiservice.NewService(
			apirepo.NewRepository(controllerTestData.fakeClient)),
		nil,
	}).Reconcile(context.Background(), controllerruntime.Request{
		NamespacedName: types.NamespacedName{
			Namespace: "kamwiel",
			Name:      "dogs-api",
		},
	})

	assert.Equal(t, result, controllerruntime.Result{})
	assert.NilError(t, err)

	events := <-resultChan
	apiListStatusDesiredMd5Hash := "df577f4565ce32c22b64f575de8dc5d8"
	apiListStatusDesiredPayload := "[{\"name\":\"dogs-api\",\"spec\":\"{\\\"openapi\\\":\\\"3.0.0\\\",\\\"info\\\":{\\\"title\\\":\\\"Best friend API\\\"}}\"}]"

	configMapData := events.Object.(*v1.ConfigMap).Data
	assert.Equal(t, events.Type, watch.EventType("MODIFIED"))
	assert.Equal(t, configMapData["hash"], apiListStatusDesiredMd5Hash)
	assert.Equal(t, configMapData["fresh"], "true")
	assert.Equal(t, configMapData["payload"], apiListStatusDesiredPayload)
}

func TestAPIControllerReconcileNoStatusUpdate(t *testing.T) {

	artifacts := buildApiControllerTestData(
		map[string]string{
			"dogs-api": "Best friend API",
		},
		map[string]string{
			"hash":    "df577f4565ce32c22b64f575de8dc5d8",
			"payload": "[{\"name\":\"dogs-api\",\"spec\":\"{\\\"openapi\\\":\\\"3.0.0\\\",\\\"info\\\":{\\\"title\\\":\\\"Best friend API\\\"}}\"}]",
			"fresh":   "true",
		},
	)

	cfMaps := make([]v1.ConfigMap, 1)
	cfMaps = append(cfMaps, *artifacts.apiListStatus)

	watchConfigMaps, _ := artifacts.fakeClient.Watch(context.Background(), &v1.ConfigMapList{Items: cfMaps}, &client.ListOptions{Namespace: "kamwiel"})
	resultChan := watchConfigMaps.ResultChan()

	result, err := (&APIReconciler{
		artifacts.fakeClient,
		apiservice.NewService(
			apirepo.NewRepository(artifacts.fakeClient)),
		nil,
	}).Reconcile(context.Background(), controllerruntime.Request{
		NamespacedName: types.NamespacedName{
			Namespace: "kamwiel",
			Name:      "dogs-api",
		},
	})

	assert.Equal(t, result, controllerruntime.Result{})
	assert.NilError(t, err)

	watchConfigMaps.Stop()
	events := <-resultChan

	assert.Equal(t, events.Type, watch.EventType("")) // Event is empty
	assert.Equal(t, events.Object, nil)               // No ConfigMap was changed
}
