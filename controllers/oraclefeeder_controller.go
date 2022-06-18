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

type OracleFeederReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *OracleFeederReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&terrav1alpha1.OracleFeeder{}).
		Owns(&corev1.Pod{}).
		Complete(r)
}

//+kubebuilder:rbac:groups=terra.terra-rebels.org,resources=oraclefeeders,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=terra.terra-rebels.org,resources=oraclefeeders/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=terra.terra-rebels.org,resources=oraclefeeders/finalizers,verbs=update
func (r *OracleFeederReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("Reconciling OracleFeeder object")

	oracleFeeder := &terrav1alpha1.OracleFeeder{}
	err := r.Client.Get(ctx, req.NamespacedName, oracleFeeder)

	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, err
	}

	pod := newPodForOracleFeeder(oracleFeeder)

	if err := controllerutil.SetControllerReference(oracleFeeder, pod, r.Scheme); err != nil {
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

func newPodForOracleFeeder(cr *terrav1alpha1.OracleFeeder) *corev1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "oraclefeeder",
					Image: cr.Spec.NodeImage,
					Env:   cr.Env,
				},
			},
		},
	}

	return pod
}
