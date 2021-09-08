/*
Copyright 2021 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	apiservice "github.com/3scale-labs/kamwiel/pkg/services/api"
	kctlrv1beta1 "github.com/kuadrant/kuadrant-controller/apis/networking/v1beta1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	k8stypes "k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// APIReconciler reconciles a API object
type APIReconciler struct {
	client.Client
	apiservice.Service
	*runtime.Scheme
}

const apiListStatusConfigMap = "kamwiel-api-list-status"

//+kubebuilder:rbac:groups=networking.kuadrant.io,resources=apis,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=networking.kuadrant.io,resources=apis/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=networking.kuadrant.io,resources=apis/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *APIReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	apiListStatus := r.getApiListStatus(ctx, req)

	apiList, listErr := r.ListAPI(ctx)
	if listErr != nil {
		return ctrl.Result{}, listErr
	}

	apiListMarshalled, marshallErr := json.Marshal(apiList)
	if marshallErr != nil {
		return ctrl.Result{}, marshallErr
	}

	newApiListHash := fmt.Sprintf("%x", md5.Sum(apiListMarshalled))

	if newApiListHash != apiListStatus.Data["hash"] {
		logger.Info("Reconcile", "New API List hash:", newApiListHash)
		newValues := map[string]string{
			"hash":    newApiListHash,
			"payload": string(apiListMarshalled),
			"fresh":   "true",
		}
		if patchErr := r.patchApiListStatus(ctx, apiListStatus, newValues); patchErr != nil {
			logger.Error(patchErr, "Error patching the APIList ConfigMap")
			return ctrl.Result{}, patchErr
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *APIReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kctlrv1beta1.API{}).
		Complete(r)
}

func (r *APIReconciler) getApiListStatus(ctx context.Context, req ctrl.Request) *v1.ConfigMap {
	apiListStatus := &v1.ConfigMap{}
	namespacedName := k8stypes.NamespacedName{Namespace: req.Namespace, Name: apiListStatusConfigMap}
	if err := r.Client.Get(ctx, namespacedName, apiListStatus); err != nil && errors.IsNotFound(err) {
		// Crete ApiListStatus ConfigMap
	}
	return apiListStatus
}

func (r *APIReconciler) patchApiListStatus(ctx context.Context, apiListStatus *v1.ConfigMap, newValues map[string]string) error {
	apiListStatusPatch := client.MergeFrom(apiListStatus.DeepCopy())
	for k, v := range newValues {
		apiListStatus.Data[k] = v
	}
	if patchErr := r.Client.Patch(ctx, apiListStatus, apiListStatusPatch); patchErr != nil {
		return patchErr
	}
	return nil
}
