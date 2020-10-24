package cmd

import (
	"github.com/vmmgr/node/data"
	"github.com/vmmgr/node/vm"
	"log"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start client server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err := vm.StartUPVM(); err != nil {
			log.Println("Error: vm auto start failed...")
		}
		data.Server()
		log.Println("end")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
