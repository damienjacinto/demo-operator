/*
Copyright 2023.

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
	"time"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	examplev1 "github.com/damienjacinto/demo-operator/api/v1"
)

// KillOldPodReconciler reconciles a KillOldPod object
type KillOldPodReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=example.com,resources=killoldpods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.com,resources=killoldpods/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.com,resources=killoldpods/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the KillOldPod object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *KillOldPodReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	mylog := log.FromContext(ctx)

	mylog.Info("Starting reconcile")
	// Fetch the KillOldPod custom resource
	killOldPod := &examplev1.KillOldPod{}
	err := r.Get(ctx, req.NamespacedName, killOldPod)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// List all Pods in the specified namespace
	podList := &corev1.PodList{}
	err = r.List(ctx, podList, client.InNamespace(killOldPod.Spec.Namespace))
	if err != nil {
		return ctrl.Result{}, err
	}
	mylog.Info("Number of pod in namespace", "nb", len(podList.Items))

	// Filter the Pods that haven't been updated for the specified number of days
	threshold := time.Now().AddDate(0, 0, -1*int(killOldPod.Spec.Minutes))
	mylog.Info("Threshold", "threshold", threshold.String())
	oldPods := []corev1.Pod{}
	for _, pod := range podList.Items {
		if pod.CreationTimestamp.Time.Before(threshold) {
			oldPods = append(oldPods, pod)
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KillOldPodReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&examplev1.KillOldPod{}).
		Complete(r)
}
