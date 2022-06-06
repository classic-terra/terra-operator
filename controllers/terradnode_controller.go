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
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	terrav1alpha1 "github.com/terra-rebels/terra-operator/api/v1alpha1"
)

type TerradNodeReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *TerradNodeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&terrav1alpha1.TerradNode{}).
		Owns(&corev1.Pod{}).
		Complete(r)
}

//+kubebuilder:rbac:groups=terra.terra-rebels.org,resources=terradnodes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=terra.terra-rebels.org,resources=terradnodes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=terra.terra-rebels.org,resources=terradnodes/finalizers,verbs=update
func (r *TerradNodeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("Reconciling TerradNode object")

	terradNode := &terrav1alpha1.TerradNode{}
	err := r.Client.Get(ctx, req.NamespacedName, terradNode)

	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, err
	}

	pod := newPodForTerradNode(terradNode)

	if err := controllerutil.SetControllerReference(terradNode, pod, r.Scheme); err != nil {
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

func newPodForTerradNode(cr *terrav1alpha1.TerradNode) *corev1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}

	// Terrad default ports @ https://docs.terra.money/docs/full-node/run-a-full-terra-node/system-config.html
	ports := []corev1.ContainerPort{
		{
			Name:          "lcd",
			ContainerPort: 1317,
		},
		{
			Name:          "p2p",
			ContainerPort: 26656,
		},
		{
			Name:          "rpc",
			ContainerPort: 26657,
		},
		{
			Name:          "prometheus",
			ContainerPort: 26660,
		},
	}

	// 4 CPUs & 32GB memory as minimum requirement @ https://docs.terra.money/docs/full-node/run-a-full-terra-node/system-config.html
	minimumRequestLimits := corev1.ResourceList{}

	if cr.Spec.IsFullNode {
		minimumRequestLimits = corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("4000m"),
			corev1.ResourceMemory: resource.MustParse("32Gi"),
		}
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
					Name:  "terrad",
					Image: cr.Spec.TerradNodeImage,
					Ports: ports,
					Resources: corev1.ResourceRequirements{
						Requests: minimumRequestLimits,
					},
				},
			},
		},
	}

	if (cr.Spec.DataVolume != corev1.Volume{}) {
		pod.Spec.Volumes = []corev1.Volume{cr.Spec.DataVolume}
		pod.Spec.Containers[0].VolumeMounts = []corev1.VolumeMount{
			{
				Name: cr.Spec.DataVolume.Name,
				//Data folder location for terrad: https://docs.terra.money/docs/full-node/run-a-full-terra-node/sync.html#quicksync
				MountPath: "/root/.terra/data/",
			},
		}
	}

	if len(cr.Spec.PostStartCommand) > 0 {
		pod.Spec.Containers[0].Lifecycle = &corev1.Lifecycle{
			PostStart: &corev1.LifecycleHandler{
				Exec: &corev1.ExecAction{
					Command: cr.Spec.PostStartCommand,
				},
			},
		}
	}

	return pod
}
