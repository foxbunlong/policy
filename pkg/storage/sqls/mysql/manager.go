package mysql

import (
	"database/sql"
	"flag"
	"strings"
	"sync"

	"github.com/oeoen/policy/driver/config"
	"github.com/oeoen/policy/helper/errorp"
	_ "github.com/oeoen/policy/pkg/storage/sqls/migrations"

	_ "github.com/go-sql-driver/mysql" // here our migrations will live  -- use your path
	"github.com/pressly/goose"
)

var (
	flags = flag.NewFlagSet("migrate", flag.ExitOnError)
	dir   = flags.String("dir", ".", "directory with migration files")
)

type MYSQLManager struct {
	db    *sql.DB
	c     config.Provider
	stmts sync.Map
}

func NewManager(c config.Provider) *MYSQLManager {

	db := &MYSQLManager{
		c:     c,
		stmts: sync.Map{},
	}

	return db
}

func (m *MYSQLManager) DBService() *sql.DB {
	return m.db
}

func (m *MYSQLManager) DoMigration(c config.Provider, args ...string) {
	connectionString := replaceDSN(c.DSN())
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		c.Logger().Fatalf("migrate: failed to open DB: %v\n", err)
	}
	arguments := []string{}
	if len(args) > 2 {
		arguments = append(arguments, args[1:]...)
	}
	command := args[0]
	goose.SetDialect("mysql")
	if err := goose.Run(command, db, *dir, arguments...); err != nil {
		c.Logger().Fatalf("migrate %v: %v", command, err)
	}
}

func (m *MYSQLManager) Open() error {
	db, err := sql.Open("mysql", replaceDSN(m.c.DSN()))
	if err != nil {
		return errorp.NewPolicyError(500, "ERR_CONN_DB", "error_connection_db_mysql", err.Error())
	}
	m.db = db
	m.c.Logger().Info("Connected to DB")
	return nil
}

func (m *MYSQLManager) Close() error {
	if m.db != nil {
		err := m.db.Close()
		if err != nil {
			return errorp.NewPolicyError(500, "ERR_CONN_DB", "error_close_db_mysql", err.Error())
		}
	}
	return nil
}

func (m *MYSQLManager) Prepare(query string) (*sql.Stmt, error) {
	var stmt *sql.Stmt
	var err error
	i, ok := m.stmts.Load(query)
	if ok {
		stmt, ok = i.(*sql.Stmt)
	}
	if !ok {
		stmt, err = m.db.Prepare(query)
		if err != nil {
			return nil, errorp.NewPolicyError(500, "ERR_STMT_DB", "error_prepare_stmt", err.Error())
		}
		m.stmts.Store(query, stmt)
	}
	return stmt, nil
}

func constructWhereQuery(query string, whereOperator [][3]string) (string, []interface{}) {
	whereStr := []string{}
	val, limit, offset := []interface{}{}, "100000", "0"
	if len(whereOperator) == 0 {
		return strings.ReplaceAll(query, ":where", ""), []interface{}{offset, limit}
	}
	for i := 0; i < len(whereOperator); i++ {
		if whereOperator[i][0] == "size" {
			limit = whereOperator[i][2]
		} else if whereOperator[i][0] == "from" {
			offset = whereOperator[i][2]
		} else {
			if whereOperator[i][1] == "LIKE" {
				whereOperator[i][2] = "%" + whereOperator[i][2] + "%"
			}

			whereStr = append(whereStr, " "+whereOperator[i][0]+" "+whereOperator[i][1]+" ?")
			val = append(val, whereOperator[i][2])
		}
	}
	val = append(val, offset, limit)
	where := strings.Join(whereStr, " AND ")
	if where != "" {
		where = "WHERE " + where
	}
	q := strings.ReplaceAll(query, ":where", where)
	return q, val
}
func replaceDSN(dsn string) string {
	connectionString := strings.ReplaceAll(dsn, "mysql://", "")
	return connectionString
}

func (m *MYSQLManager) Stmts() sync.Map {
	return m.stmts
}
