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
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
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

// TerradNetReconciler reconciles a TerradNet object
type TerradNetReconciler struct {
	client.Client
	Scheme  *runtime.Scheme
	Replica int32
}

//+kubebuilder:rbac:groups=terra.terra-rebels.org,resources=TerradNets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=terra.terra-rebels.org,resources=TerradNets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=terra.terra-rebels.org,resources=TerradNets/finalizers,verbs=update
//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch,resources=jobs/status,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the TerradNet object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *TerradNetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("Reconciling TerradNet object")

	terradNet := &terrav1alpha1.TerradNet{}
	err := r.Client.Get(ctx, req.NamespacedName, terradNet)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, err
	}

	job := newJobForTerradNet(terradNet, r.Replica)

	if err := controllerutil.SetControllerReference(terradNet, job, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	foundJob := &batchv1.Job{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: job.Name, Namespace: job.Namespace}, foundJob)

	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating a new Job", "job.Namespace", job.Namespace, "job.Name", job.Name)

		err = r.Client.Create(context.TODO(), job)

		if err != nil {
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	}

	statefulSet := newStatefulSetForTerradNet(terradNet, r.Replica)

	if err := controllerutil.SetControllerReference(terradNet, statefulSet, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	foundSet := &appsv1.StatefulSet{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: statefulSet.Name, Namespace: statefulSet.Namespace}, foundSet)

	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating a new StatefulSet", "StatefulSet.Namespace", statefulSet.Namespace, "StatefulSet.Name", statefulSet.Name)

		err = r.Client.Create(context.TODO(), statefulSet)

		if err != nil {
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TerradNetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&terrav1alpha1.TerradNet{}).
		Owns(&batchv1.Job{}).
		Owns(&appsv1.StatefulSet{}).
		Complete(r)
}

func newJobForTerradNet(cr *terrav1alpha1.TerradNet, replica int32) *batchv1.Job {
	containers := []corev1.Container{
		{
			Name:            "job",
			Image:           cr.Spec.Container.Image,
			ImagePullPolicy: corev1.PullPolicy(cr.Spec.Container.ImagePullPolicy),
			Command:         []string{"/bin/sh", "/setup-localnet.sh"},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      cr.Spec.DataSource.Name,
					MountPath: "/terra",
				},
			},
			Env: []corev1.EnvVar{
				{
					Name:  "CHAINID",
					Value: cr.Spec.ChainId,
				},
				{
					Name:  "REPLICA",
					Value: strconv.FormatInt(int64(replica), 10),
				},
			},
		},
	}

	backOffLimit := int32(3)

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("localnet-%s-setup", cr.Name),
			Namespace: "default",
		},
		Spec: batchv1.JobSpec{
			BackoffLimit: &backOffLimit,
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers:    containers,
					RestartPolicy: corev1.RestartPolicyNever,
					Volumes:       []corev1.Volume{cr.Spec.DataSource},
				},
			},
		},
	}

	return job
}

func newStatefulSetForTerradNet(cr *terrav1alpha1.TerradNet, replica int32) *appsv1.StatefulSet {
	nodeDataVolumeName := "node-data"

	labels := map[string]string{
		"net": cr.Name,
	}

	node_labels := map[string]string{
		"node": cr.Spec.ChainId,
	}

	// Replica defined in CLI will overwrite the one in yaml
	if replica == 0 {
		replica = cr.Spec.Replica
	}

	// if replica is still 0, set it to default 1
	if replica == 0 {
		replica = 1
	}

	initContainer := []corev1.Container{{
		Name:    "init",
		Image:   "busybox",
		Command: []string{"sh", "-c", "cp -R /data/mytestnet /target && chown -R 1000:1000 /target"},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      cr.Spec.DataSource.Name,
				MountPath: "/data",
			}, {
				Name:      nodeDataVolumeName,
				MountPath: "/target",
			},
		},
	}}

	container := []corev1.Container{
		{
			Name:            "node",
			Image:           cr.Spec.Container.Image,
			ImagePullPolicy: corev1.PullPolicy(cr.Spec.Container.ImagePullPolicy),
			Ports: []corev1.ContainerPort{
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
			},
			Env: []corev1.EnvVar{
				{
					Name:  "IS_LOCALNET",
					Value: "true",
				},
				{
					Name:  "SERVICE_NAME",
					Value: cr.Spec.ServiceName,
				},
				{
					Name:  "CHAINID",
					Value: cr.Spec.ChainId,
				},
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      nodeDataVolumeName,
					MountPath: "/terra",
				},
			},
		},
	}

	persistentVolumeClaim := []corev1.PersistentVolumeClaim{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: nodeDataVolumeName,
			},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceStorage: resource.MustParse("50Mi"),
					},
				},
			},
		},
	}

	statefulSet := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.StatefulSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: node_labels,
			},
			ServiceName: cr.Spec.ServiceName,
			Replicas:    &replica,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: node_labels,
				},
				Spec: corev1.PodSpec{
					InitContainers: initContainer,
					Containers:     container,
					Volumes:        []corev1.Volume{cr.Spec.DataSource},
				},
			},
			VolumeClaimTemplates: persistentVolumeClaim,
		},
	}

	return statefulSet
}
