package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/maxvasylets/kube-pod-rescheduler/controller"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run Kubernetes controller.",
	Long:  "Run Kubernetes controller.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("the controller is starting ...")
		controller.Run()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
