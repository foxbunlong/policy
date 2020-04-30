package mysql

import (
	"context"
	"time"

	"github.com/oeoen/policy/helper/errorp"
	"github.com/oeoen/policy/pkg/police"
	"github.com/oeoen/policy/pkg/storage/sqls/mysql/queries"
)

func (m *MYSQLManager) UpsertRole(ctx context.Context, rbac *police.RBAC) error {
	sp := startSpan(ctx, "upsert-role", queries.InsertRole)
	defer finishSpan(sp)
	stmt, err := m.Prepare(queries.InsertRole)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(rbac.Subject, rbac.Tenant, rbac.Policy)
	if err != nil {
		return errorp.NewPolicyError(500, "ERR_INSERT_POLICY", "error_insert_policy", err.Error())
	}
	rbac.Created = time.Now()
	return nil
}

func (m *MYSQLManager) GetRoles(ctx context.Context, tenant string) ([]string, error) {
	sp := startSpan(ctx, "get-roles", queries.GetRoles)
	defer finishSpan(sp)
	result := []string{}
	r, err := m.DBService().Query(queries.GetRoles, tenant)
	if err != nil {
		return nil, errorp.NewPolicyError(500, "ERR_FETCH_ROLES", "error_fetch_query_roles", err.Error())
	}
	defer r.Close()
	for r.Next() {
		re := ""
		err = r.Scan(&re)
		if err != nil {
			return nil, errorp.NewPolicyError(500, "ERR_FETCH_ROLES", "error_fetch_scan_roles", err.Error())
		}
		result = append(result, re)
	}
	if len(result) == 0 {
		return nil, errorp.NewPolicyError(404, "ERR_FETCH_ROLES", "error_not_found_roles", "No data")
	}
	return result, nil
}

func (m *MYSQLManager) GetRoleSubjects(ctx context.Context, tenant, policy string) ([]string, error) {
	sp := startSpan(ctx, "get-role-subjects", queries.GetRoleSubjects)
	defer finishSpan(sp)
	result := []string{}
	r, err := m.DBService().Query(queries.GetRoleSubjects, tenant, policy)

	if err != nil {
		return nil, errorp.NewPolicyError(500, "ERR_FETCH_ROLE", "error_fetch_query_role", err.Error())
	}
	defer r.Close()
	for r.Next() {
		re := ""
		err = r.Scan(&re)
		if err != nil {
			return nil, errorp.NewPolicyError(500, "ERR_FETCH_ROLE", "error_fetch_scan_role", err.Error())
		}
		result = append(result, re)
	}
	if len(result) == 0 {
		return nil, errorp.NewPolicyError(404, "ERR_FETCH_ROLE", "error_not_found_role", "No data")
	}
	return result, nil
}
func (m *MYSQLManager) GetSubjectRoles(ctx context.Context, tenant, subject string) ([]string, error) {
	sp := startSpan(ctx, "get-subject-roles", queries.GetSubjectRoles)
	defer finishSpan(sp)
	result := []string{}
	r, err := m.DBService().Query(queries.GetSubjectRoles, tenant, subject)
	if err != nil {
		return nil, errorp.NewPolicyError(500, "ERR_FETCH_SUBJECT", "error_fetch_query_subjects", err.Error())
	}
	defer r.Close()
	for r.Next() {
		re := ""
		err = r.Scan(&re)
		if err != nil {
			return nil, errorp.NewPolicyError(500, "ERR_FETCH_SUBJECT", "error_fetch_scan_subjects", err.Error())
		}
		result = append(result, re)
	}
	if len(result) == 0 {
		return nil, errorp.NewPolicyError(404, "ERR_FETCH_SUBJECT", "error_not_found_subjects", "No data")
	}
	return result, nil
}
func (m *MYSQLManager) DeleteRole(ctx context.Context, tenant, subject, policy string) error {
	sp := startSpan(ctx, "delete-role", queries.DeleteRole)
	defer finishSpan(sp)
	stmt, err := m.Prepare(queries.DeleteRole)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(tenant, subject, policy)
	if err != nil {
		return errorp.NewPolicyError(500, "ERR_DELETE_ROLE", "error_delete_role", err.Error())
	}
	return nil
}
