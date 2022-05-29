package terradnode

import (
	"context"

	terrav1alpha1 "github.com/terra-rebels/terra-operator/pkg/apis/terra/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_terradnode")

// Add creates a new TerradNode Controller and adds it to the Manage
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileTerradNode{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	c, err := controller.New("terradnode-controller", mgr, controller.Options{Reconciler: r})

	if err != nil {
		return err
	}

	// Watch for changes to primary resource TerradNode
	err = c.Watch(&source.Kind{Type: &terrav1alpha1.TerradNode{}}, &handler.EnqueueRequestForObject{})

	if err != nil {
		return err
	}

	// Watch for changes to secondary resources and requeue the owner TerradNode
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &terrav1alpha1.TerradNode{},
	})

	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &terrav1alpha1.TerradNode{},
	})

	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileTerradNode implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileTerradNode{}

// ReconcileTerradNode reconciles a TerradNode object
type ReconcileTerradNode struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a TerradNode object and makes changes based on the state read
// and what is in the TerradNode.Spec
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileTerradNode) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling TerradNode")

	instance := &terrav1alpha1.TerradNode{}

	err := r.client.Get(context.TODO(), request.NamespacedName, instance)

	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}

		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Define a new Pod object
	pod := newPodForCR(instance)

	// Set TerradNode instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Pod already exists
	foundPod := &corev1.Pod{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, foundPod)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)

		//Create pod
		err = r.client.Create(context.TODO(), pod)

		if err != nil {
			// Pod creation failed - requeue
			return reconcile.Result{}, err
		}

		// Pod created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Define a new Service object
	service := newServiceForCR(instance)

	// TODO: Figure out if we should have 1 service per pod or one for all pods.
	service.Spec.Selector = pod.ObjectMeta.Labels

	// Set TerradNode instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, service, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Service already exists
	foundService := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, foundService)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
		err = r.client.Create(context.TODO(), service)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Pod Service successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Everything is fine, dont requeue
	return reconcile.Result{}, nil
}

func newPodForCR(cr *terrav1alpha1.TerradNode) *corev1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}

	// Terrad default ports @ https://docs.terra.money/docs/full-node/run-a-full-terra-node/system-config.html
	ports := []corev1.ContainerPort{
		{
			Name:          "LCD",
			ContainerPort: 1317,
		},
		{
			Name:          "P2P",
			ContainerPort: 26656,
		},
		{
			Name:          "RPC",
			ContainerPort: 26657,
		},
		{
			Name:          "Prometheus",
			ContainerPort: 26660,
		},
	}

	// 4 CPUs, 32GB memory & 2TB of storage as minimum requirement @ https://docs.terra.money/docs/full-node/run-a-full-terra-node/system-config.html
	minimumRequestLimits := corev1.ResourceList{
		corev1.ResourceCPU:     resource.MustParse("4000m"),
		corev1.ResourceMemory:  resource.MustParse("32GiB"),
		corev1.ResourceStorage: resource.MustParse("2TiB"),
	}

	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "terrad",
					Image:   "terramoney/core-node:v0.5.11-oracle",
					EnvFrom: cr.EnvFrom,
					Ports:   ports,
					Resources: corev1.ResourceRequirements{
						Requests: minimumRequestLimits,
					},
				},
			},
		},
	}
}

func newServiceForCR(cr *terrav1alpha1.TerradNode) *corev1.Service {
	labels := map[string]string{
		"app": cr.Name,
	}

	// TODO: Figure out how to dynamically generate service ports (managerCount + 1 style logic)
	const (
		p2pPort = 99900 + iota
		rpcPort = 99900 + iota
		lcdPort = 99900 + iota
	)

	ports := []corev1.ServicePort{
		{
			Port:       p2pPort,
			TargetPort: intstr.FromString("P2P"),
		},
		{
			Port:       rpcPort,
			TargetPort: intstr.FromString("RPC"),
		},
		{
			Port:       lcdPort,
			TargetPort: intstr.FromString("LCD"),
		},
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-service",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: ports,
		},
	}
}
