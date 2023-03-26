package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/terra-rebels/terra-operator/controllers"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type Reconciler interface {
	Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error)
	SetupWithManager(mgr ctrl.Manager) error
}

func SetupManager(cmd *cobra.Command) (manager.Manager, error) {
	metricsAddr, err := cmd.Flags().GetString(MetricsAddr)
	if err != nil {
		return nil, err
	}

	probeAddr, err := cmd.Flags().GetString(ProbeAddr)
	if err != nil {
		return nil, err
	}

	enableLeaderElection, err := cmd.Flags().GetBool(LeaderElection)
	if err != nil {
		return nil, err
	}

	opts := zap.Options{
		Development: true,
	}

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "fea5d43e.terra-rebels.org",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		return nil, err
	}

	return mgr, nil
}

func AssembleReconciler(cmd *cobra.Command, mgr manager.Manager) error {

	reconcilers := []Reconciler{}

	if enable, _ := cmd.Flags().GetBool(RelayerService); enable {
		reconcilers = append(reconcilers, &controllers.RelayerReconciler{
			Client: mgr.GetClient(),
			Scheme: mgr.GetScheme(),
		})
	}

	for _, reconciler := range reconcilers {
		if err := reconciler.SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", reconciler)
			return err
		}
	}

	return nil
}

func StartManager(mgr manager.Manager) error {
	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		return err
	}

	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		return err
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		return err
	}

	return nil
}
