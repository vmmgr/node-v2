package cmd

import (
	"fmt"
	"github.com/yoneyan/vm_mgr/node/data"
	"github.com/yoneyan/vm_mgr/node/vm"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start client server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		vm.StartupProcess()
		data.Server()
		fmt.Println("end")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
