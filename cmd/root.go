package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"

	terrav1alpha1 "github.com/terra-rebels/terra-operator/api/v1alpha1"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")

	rootCmd = &cobra.Command{
		Use:   "kubectl-terra_operator",
		Short: "A Kubernetes operator for Terra",
	}
)

const (
	RelayerService = "enable-relayer"
	LeaderElection = "leader-elect"
	MetricsAddr    = "metrics-bind-address"
	ProbeAddr      = "health-probe-bind-address"
)

// Execute executes the root command.
func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.AddCommand(GetNetworkCmd())
	rootCmd.AddCommand(GetNodeCmd())

	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(terrav1alpha1.AddToScheme(scheme))
}

func AddCommonFlagsToCmd(cmd *cobra.Command) {
	cmd.Flags().Bool(RelayerService, false, "enable relayer service")
	cmd.Flags().String(MetricsAddr, ":8080", "The address the metric endpoint binds to.")
	cmd.Flags().String(ProbeAddr, ":8081", "The address the probe endpoint binds to.")
	cmd.Flags().Bool(LeaderElection, false, "Enable leader election for controller manager. "+
		"Enabling this will ensure there is only one active controller manager.")
}
