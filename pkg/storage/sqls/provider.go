package sqls

import (
	"context"
	"database/sql"
	"strings"
	"sync"

	"github.com/oeoen/policy/driver/config"
	"github.com/oeoen/policy/pkg/police"
	"github.com/oeoen/policy/pkg/storage/sqls/mysql"
)

type SQLSManager interface {
	DoMigration(c config.Provider, args ...string)
	DBService() *sql.DB
	Open() error
	Close() error
	police.Manager
	Stmts() sync.Map
	Prepare(query string) (*sql.Stmt, error)
}

type SQLs struct {
	dbManager SQLSManager
}

func NewSQLS(c config.Provider) *SQLs {
	p := &SQLs{}
	p.DBInit(c)
	return p
}

func (p *SQLs) DBInit(c config.Provider) error {
	databaseKind := getDSN(c.DSN())
	if databaseKind == "mysql" {
		p.dbManager = mysql.NewManager(c)
		err := p.Manager().Open()
		if err != nil {
			c.Logger().Fatal(err)
		}
	}
	return nil
}

func (p *SQLs) DBDefer() error {
	return p.Manager().Close()
}
func (p *SQLs) DoMigration(c config.Provider, args ...string) {
	p.dbManager.DoMigration(c, args...)
}

func getDSN(dsn string) string {
	if dsn == "" {
		return "No DSN"
	}
	if strings.Contains(dsn, "mysql") {
		return "mysql"
	}
	return ""
}
func (p *SQLs) Manager() SQLSManager {
	return p.dbManager
}
func (p *SQLs) UpsertPolicy(ctx context.Context, acl *police.ACL) error {
	return p.Manager().UpsertPolicy(ctx, acl)
}
func (p *SQLs) GetPolicy(ctx context.Context, policeID string) (*police.ACL, error) {
	return p.Manager().GetPolicy(ctx, policeID)
}
func (p *SQLs) FetchPolicy(ctx context.Context, filter ...[3]string) ([]*police.ACL, error) {
	return p.Manager().FetchPolicy(ctx, filter...)
}
func (p *SQLs) DeletePolicy(ctx context.Context, policeID string) error {
	return p.Manager().DeletePolicy(ctx, policeID)
}
func (p *SQLs) UpsertRole(ctx context.Context, rbac *police.RBAC) error {
	return p.Manager().UpsertRole(ctx, rbac)
}
func (p *SQLs) GetRoles(ctx context.Context, tenant string) ([]string, error) {
	return p.Manager().GetRoles(ctx, tenant)
}
func (p *SQLs) DeleteRole(ctx context.Context, tenant, subject, policy string) error {
	return p.Manager().DeleteRole(ctx, tenant, subject, policy)
}
func (p *SQLs) GetRoleSubjects(ctx context.Context, tenant, policy string) ([]string, error) {
	return p.Manager().GetRoleSubjects(ctx, tenant, policy)
}
func (p *SQLs) GetSubjectRoles(ctx context.Context, tenant, subject string) ([]string, error) {
	return p.Manager().GetRoleSubjects(ctx, tenant, subject)
}

func (p *SQLs) GetResources(ctx context.Context) ([]string, error) {
	return p.Manager().GetResources(ctx)

}
func (p *SQLs) GetPolicySubjects(ctx context.Context) ([]string, error) {
	return p.Manager().GetPolicySubjects(ctx)
}
func (p *SQLs) GetRolePolicy(ctx context.Context, tenant, subject string) (*police.ACL, error) {
	return p.Manager().GetRolePolicy(ctx, tenant, subject)
}
func (p *SQLs) UpdatePolicy(ctx context.Context, policeID string, acl *police.ACL) error {
	return p.Manager().UpdatePolicy(ctx, policeID, acl)
}

func (p *SQLs) Enforce(ctx context.Context, tenant, subject, action, resource string) (*police.ACL, error) {
	return p.Manager().Enforce(ctx, tenant, subject, action, resource)
}
