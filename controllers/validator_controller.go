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
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	terrav1alpha1 "github.com/terra-rebels/terra-operator/api/v1alpha1"
)

type ValidatorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *ValidatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&terrav1alpha1.Validator{}).
		Owns(&terrav1alpha1.TerradNode{}).
		Owns(&terrav1alpha1.OracleNode{}).
		Owns(&terrav1alpha1.IndexerNode{}).
		Owns(&corev1.Service{}).
		Complete(r)
}

// +kubebuilder:rbac:groups=terra.terra-rebels.org,resources=validators,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=terra.terra-rebels.org,resources=validators/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=terra.terra-rebels.org,resources=validators/finalizers,verbs=update
func (r *ValidatorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("Reconciling Validator object")

	validator := &terrav1alpha1.Validator{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, validator)

	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, err
	}

	terradNode := newTerradNodeForValidator(validator)

	if err := controllerutil.SetControllerReference(validator, terradNode, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	foundTerrad := &terrav1alpha1.TerradNode{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: terradNode.Name, Namespace: terradNode.Namespace}, foundTerrad)

	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating a new TerradNode", "TerradNode.Namespace", terradNode.Namespace, "TerradNode.Name", terradNode.Name)

		err = r.Client.Create(context.TODO(), terradNode)

		if err != nil {
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	}

	oracleNode := newOracleNodeForValidator(validator)

	if err := controllerutil.SetControllerReference(validator, oracleNode, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	foundOracleNode := &terrav1alpha1.OracleNode{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: oracleNode.Name, Namespace: oracleNode.Namespace}, foundOracleNode)

	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating a new OracleNode", "OracleNode.Namespace", oracleNode.Namespace, "OracleNode.Name", oracleNode.Name)

		err = r.Client.Create(context.TODO(), oracleNode)

		if err != nil {
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	}

	indexerNode := newIndexerNodeForValidator(validator)

	if err := controllerutil.SetControllerReference(validator, indexerNode, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	foundIndexerNode := &terrav1alpha1.IndexerNode{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: indexerNode.Name, Namespace: indexerNode.Namespace}, foundIndexerNode)

	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating a new IndexerNode", "IndexerNode.Namespace", indexerNode.Namespace, "IndexerNode.Name", indexerNode.Name)

		err = r.Client.Create(context.TODO(), indexerNode)

		if err != nil {
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	}

	if validator.Spec.IsPublic {
		service := newServiceForValidator(validator)

		if err := controllerutil.SetControllerReference(validator, service, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}

		foundService := &corev1.Service{}
		err = r.Client.Get(context.TODO(), types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, foundService)

		if err != nil && errors.IsNotFound(err) {
			logger.Info("Creating a new Service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)

			err = r.Client.Create(context.TODO(), service)

			if err != nil {
				return ctrl.Result{}, err
			}

			return ctrl.Result{}, nil
		} else if err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func newIndexerNodeForValidator(cr *terrav1alpha1.Validator) *terrav1alpha1.IndexerNode {
	labels := map[string]string{
		"app": cr.Name,
	}

	envVars := []corev1.EnvVar{
		{
			Name:  "CHAIN_ID",
			Value: cr.Spec.ChainId,
		},
		{
			//TODO: Make this configurable from spec
			Name:  "INDEXER_INITIAL_HEIGHT",
			Value: "1",
		},
		{
			//TODO: Make this configurable from spec
			Name:  "INDEXER_LCD_URI",
			Value: "https://lcd.terra.dev",
		},
		{
			//TODO: Make this configurable from spec
			Name:  "INDEXER_FCD_URI",
			Value: "https://fcd.terra.dev",
		},
		{
			//TODO: Make this configurable from spec
			Name:  "INDEXER_RPC_URI",
			Value: "https://localhost:26657",
		},
		{
			//TODO: Make this configurable from spec
			Name:  "INDEXER_SERVER_PORT",
			Value: "3060",
		},
		{
			//TODO: Make this configurable from spec
			Name:  "INDEXER_TOKEN_NETWORK",
			Value: "mainnet",
		},
	}

	indexerNode := &terrav1alpha1.IndexerNode{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-indexernode",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Env: envVars,
		Spec: terrav1alpha1.IndexerNodeSpec{
			NodeImages: cr.Spec.IndexerNodeImages,
		},
	}

	return indexerNode
}

func newOracleNodeForValidator(cr *terrav1alpha1.Validator) *terrav1alpha1.OracleNode {
	labels := map[string]string{
		"app": cr.Name,
	}

	envVars := []corev1.EnvVar{
		{
			Name:  "CHAIN_ID",
			Value: cr.Spec.ChainId,
		},
		{
			Name:  "ORACLE_FEEDER_LCD_URIS",
			Value: "https://lcd.terra.dev",
		},
		{
			//TODO: Make this configurable from spec
			Name:  "ORACLE_FEEDER_PRICE_SERVER_URI",
			Value: "http://localhost:8532/latest",
		},
		{
			//TODO: Make this configurable from spec
			Name:  "ORACLE_FEEDER_VALIDATOR_ADDRESSES",
			Value: "terravaloper1xx",
		},
		{
			Name:  "ORACLE_FEEDER_PASSPHRASE",
			Value: cr.Spec.Passphrase,
		},
		{
			Name:  "ORACLE_FEEDER_MNEMONIC",
			Value: cr.Spec.Mnenomic,
		},
	}

	oracleNode := &terrav1alpha1.OracleNode{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-oraclenode",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Env: envVars,
		Spec: terrav1alpha1.OracleNodeSpec{
			NodeImages: cr.Spec.OracleNodeImages,
		},
	}

	return oracleNode
}

func newTerradNodeForValidator(cr *terrav1alpha1.Validator) *terrav1alpha1.TerradNode {
	labels := map[string]string{
		"app": cr.Name,
	}

	envVars := []corev1.EnvVar{
		{
			Name:  "CHAINID",
			Value: cr.Spec.ChainId,
		},
		{
			Name:  "VALIDATOR_KEYNAME",
			Value: cr.Name,
		},
		{
			Name:  "VALIDATOR_PASSPHRASE",
			Value: cr.Spec.Passphrase,
		},
		{
			Name:  "VALIDATOR_MNENOMIC",
			Value: cr.Spec.Mnenomic,
		},
		{
			Name:  "VALIDATOR_AMOUNT",
			Value: cr.Spec.Amount,
		},
		{
			Name:  "VALIDATOR_COMMISSION_RATE",
			Value: cr.Spec.CommissionRate,
		},
		{
			Name:  "VALIDATOR_COMMISSION_RATE_MAX",
			Value: cr.Spec.CommissionRateMax,
		},
		{
			Name:  "VALIDATOR_COMMISSION_RATE_MAX_CHANGE",
			Value: cr.Spec.CommissionRateMaxChange,
		},
		{
			Name:  "VALIDATOR_MIN_SELF_DELEGATION",
			Value: cr.Spec.MinimumSelfDelegation,
		},
	}

	if cr.Spec.AutoConfig {
		envVars = append(envVars, corev1.EnvVar{Name: "VALIDATOR_AUTO_CONFIG", Value: "1"})
	}

	terrad := &terrav1alpha1.TerradNode{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-terrad",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Env: envVars,
		Spec: terrav1alpha1.TerradNodeSpec{
			Container: terrav1alpha1.ContainerSpec{
				Image:           cr.Spec.TerradNodeImage,
				ImagePullPolicy: string(corev1.PullAlways),
			},
			IsFullNode: true,
			DataVolume: cr.Spec.DataVolume,
		},
	}

	return terrad
}

func newServiceForValidator(cr *terrav1alpha1.Validator) *corev1.Service {
	labels := map[string]string{
		"app": cr.Name,
	}

	selector := map[string]string{
		"app": cr.Name + "-terrad",
	}

	ports := []corev1.ServicePort{
		{
			Name:       "p2p",
			Port:       26656,
			TargetPort: intstr.FromString("p2p"),
		},
		{
			Name:       "rpc",
			Port:       26657,
			TargetPort: intstr.FromString("rpc"),
		},
		{
			Name:       "lcd",
			Port:       1317,
			TargetPort: intstr.FromString("lcd"),
		},
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-service",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Ports:    ports,
			Selector: selector,
		},
	}
}
