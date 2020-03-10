package cli

import (
	"github.com/oeoen/policy/driver"
	"github.com/oeoen/policy/pkg/storage/sqls"
	"github.com/ory/x/logrusx"
	"github.com/spf13/cobra"
)

func MigrateSQLHandler(cmd *cobra.Command, args []string) {
	var d driver.Driver = driver.NewDefaultDriver(logrusx.New(), false)
	sqlsDB := d.Registry().PoliceManager().(*sqls.SQLs)
	sqlsDB.DoMigration(d.Configuration(), args...)

}
