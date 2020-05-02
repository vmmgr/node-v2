package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/yoneyan/vm_mgr/node/db"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize",
	Long: `initialize command. For example:

database init: init database
`,
}
var initdbCmd = &cobra.Command{
	Use:   "db",
	Short: "db init",
	Long:  "db init cmd",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(db.Createdb())
		return nil
	},
}
var initNodeCmd = &cobra.Command{
	Use:   "client",
	Short: "client init",
	Long:  "client init cmd",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Not implemented")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.AddCommand(initdbCmd)
	initCmd.AddCommand(initNodeCmd)
}
