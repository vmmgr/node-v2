package cmd

import (
	"github.com/spf13/cobra"
	gateway "github.com/vmmgr/node/pkg/api/core/gateway/v0"
	"github.com/vmmgr/node/pkg/api/core/tool/config"
	"log"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start client server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		confPath, err := cmd.Flags().GetString("config")
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}

		if config.GetConfig(confPath) != nil {
			log.Fatalf("error config process |%v", err)
		}

		gateway.NodeAPI()

		log.Println("end")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.PersistentFlags().StringP("config", "c", "", "config path")
}
