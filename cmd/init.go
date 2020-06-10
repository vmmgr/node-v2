package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize",
	Long: `initialize command. For example:

database init: init database
`,
}
var initDBCmd = &cobra.Command{
	Use:   "db",
	Short: "db init",
	Long:  "db init cmd",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
var initNodeCmd = &cobra.Command{
	Use:   "client",
	Short: "client init",
	Long:  "client init cmd",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("Not implemented")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.AddCommand(initDBCmd)
	initCmd.AddCommand(initNodeCmd)
}
