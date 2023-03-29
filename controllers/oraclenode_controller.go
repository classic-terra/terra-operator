/*
Copyright 2022.

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

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	terrav1alpha1 "github.com/terra-rebels/terra-operator/api/v1alpha1"
)

type OracleNodeReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *OracleNodeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&terrav1alpha1.OracleNode{}).
		Owns(&corev1.Pod{}).
		Complete(r)
}

// +kubebuilder:rbac:groups=terra.terra-rebels.org,resources=oraclenodes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=terra.terra-rebels.org,resources=oraclenodes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=terra.terra-rebels.org,resources=oraclenodes/finalizers,verbs=update
func (r *OracleNodeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("Reconciling OracleNode object")

	oracleNode := &terrav1alpha1.OracleNode{}
	err := r.Client.Get(ctx, req.NamespacedName, oracleNode)

	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, err
	}

	pod := newPodForOracleNode(oracleNode)

	if err := controllerutil.SetControllerReference(oracleNode, pod, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	foundPod := &corev1.Pod{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, foundPod)

	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)

		err = r.Client.Create(context.TODO(), pod)

		if err != nil {
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func newPodForOracleNode(cr *terrav1alpha1.OracleNode) *corev1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}

	containers := make([]corev1.Container, len(cr.Spec.NodeImages))

	for i := 0; i < len(cr.Spec.NodeImages); i++ {
		containers = append(containers, corev1.Container{
			Name:  fmt.Sprintf("oraclenode-%d", i),
			Image: cr.Spec.NodeImages[i],
			Env:   cr.Env,
		})
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: containers,
		},
	}

	return pod
}
