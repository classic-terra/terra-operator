package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/terra-rebels/terra-operator/controllers"
)

const (
	NetworkReplica = "replica"
)

func GetNetworkCmd() *cobra.Command {
	networkCmd := &cobra.Command{
		Use:   "network",
		Short: "start network deployment",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			replica, err := cmd.Flags().GetInt(NetworkReplica)
			if err != nil {
				return err
			}

			mgr, err := SetupManager(cmd)
			if err != nil {
				return err
			}

			if err = (&controllers.TerraNetDeploymentReconciler{
				Client:  mgr.GetClient(),
				Scheme:  mgr.GetScheme(),
				Replica: replica,
			}).SetupWithManager(mgr); err != nil {
				return fmt.Errorf("unable to create controller TerraNetDeployment: %v", err)
			}

			if err := AssembleReconciler(cmd, mgr); err != nil {
				return err
			}

			if err := StartManager(mgr); err != nil {
				return err
			}

			return nil
		},
	}

	networkCmd.Flags().Int(NetworkReplica, 1, "Number of network replica")

	// default option
	AddCommonFlagsToCmd(networkCmd)

	return networkCmd
}
