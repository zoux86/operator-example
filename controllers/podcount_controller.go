/*
Copyright 2022 zoux86.

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
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	zouxappv1 "github.com/zoux86/operator-example/api/v1"
)

// PodCountReconciler reconciles a PodCount object
type PodCountReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=zouxapp.github.com,resources=podcounts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=zouxapp.github.com,resources=podcounts/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=zouxapp.github.com,resources=podcounts/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PodCount object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile
func (r *PodCountReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	rlog := log.FromContext(ctx)
	rlog.Info("start to reconciling podCount %s", req.Name)
	podCount := &zouxappv1.PodCount{}
	err := r.Client.Get(ctx, req.NamespacedName, podCount)
	if err != nil {
		rlog.Error(err, fmt.Sprintf("get podcount %s/%s err during reconcile.", req.Namespace, req.Name))
		return ctrl.Result{}, nil
	}
	podCountCopy := podCount.DeepCopy()
	if podCount.Spec.Count <= 0 {
		podCountCopy.Status.Count = 0
	} else {
		podCountCopy.Status.Count = podCount.Spec.Count
	}

	err = r.Client.Status().Update(ctx, podCountCopy)
	if err != nil {
		rlog.Error(err, fmt.Sprintf("update crd podcount status error %s/%s  during reconcile.", req.Namespace, req.Name))
	}
	//r.Status().Update(ctx, podCountCopy, metav1.UpdateOptions{})
	// TODO(user): your logic here

	return ctrl.Result{}, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *PodCountReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&zouxappv1.PodCount{}).
		Complete(r)
}
