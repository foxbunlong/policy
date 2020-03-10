package cmd

import (
	"github.com/oeoen/policy/cmd/cli"
	"github.com/spf13/cobra"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Various migration helpers",
	Run:   cli.MigrateSQLHandler,
}

func init() {
	RootCmd.AddCommand(migrateCmd)
}
