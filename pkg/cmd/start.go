package cmd

import (
	"github.com/spf13/cobra"
	gateway "github.com/vmmgr/node/pkg/api/core/gateway/v0"
	"log"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start client server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		gateway.NodeAPI()

		log.Println("end")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
