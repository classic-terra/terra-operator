package cmd

import (
	"github.com/spf13/cobra"
	"github.com/terra-rebels/terra-operator/controllers"
)

func GetNodeCmd() *cobra.Command {
	nodeCmd := &cobra.Command{
		Use:   "node",
		Short: "start node deployment",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr, err := SetupManager(cmd)
			if err != nil {
				return err
			}

			if err = (&controllers.TerradNodeReconciler{
				Client: mgr.GetClient(),
				Scheme: mgr.GetScheme(),
			}).SetupWithManager(mgr); err != nil {
				setupLog.Error(err, "unable to create controller", "controller", "TerradNode")
				return err
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

	// default option
	AddCommonFlagsToCmd(nodeCmd)

	return nodeCmd
}
