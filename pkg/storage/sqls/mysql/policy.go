package mysql

import (
	"context"
	"math"
	"strconv"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/oeoen/policy/helper/errorp"
	"github.com/oeoen/policy/pkg/police"
	"github.com/oeoen/policy/pkg/storage/sqls/mysql/queries"
)

func (m *MYSQLManager) UpsertPolicy(ctx context.Context, acl *police.ACL) error {
	sp := startSpan(ctx, "upsert-policy", queries.InsertPolicy)
	defer finishSpan(sp)
	stmt, err := m.Prepare(queries.InsertPolicy)
	if err != nil {
		return err
	}
	u2, err := uuid.NewV4()
	if err != nil {
		return errorp.NewPolicyError(500, "ERR_GENERATE_UUID", "error_generate_uuid", err.Error())
	}
	acl.ID = u2.String()
	_, err = stmt.Exec(acl.ID, acl.Subject, acl.Tentant, acl.Resource, acl.Action, acl.Effect, acl.Active, acl.Expired)
	if err != nil {
		return errorp.NewPolicyError(500, "ERR_INSERT_POLICY", "error_insert_policy", err.Error())
	}
	return nil
}
func (m *MYSQLManager) GetPolicy(ctx context.Context, policeID string) (*police.ACL, error) {
	acls, err := m.FetchPolicy(ctx, [3]string{"uuid", "=", policeID})
	if err != nil {
		return nil, err
	}
	return acls[0], nil
}
func (m *MYSQLManager) FetchPolicy(ctx context.Context, filter ...[3]string) ([]*police.ACL, error) {
	query, k := constructWhereQuery(queries.GetPolicy, filter)
	acls := []*police.ACL{}
	sp := startSpan(ctx, "fetch-policy", query)
	defer finishSpan(sp)

	r, err := m.DBService().Query(query, k...)
	if err != nil {
		return nil, errorp.NewPolicyError(500, "ERR_FETCH_POLICY", "error_fetch_query_policy", err.Error())
	}
	defer r.Close()
	for r.Next() {
		acl := police.ACL{}
		err = r.Scan(
			&acl.ID,
			&acl.Subject,
			&acl.Tentant,
			&acl.Resource,
			&acl.Action,
			&acl.Effect,
			&acl.Active,
			&acl.Expired,
			&acl.Created,
			&acl.Updated,
		)
		if acl.Expired != nil {
			if acl.Expired.IsZero() {
				acl.Expired = nil
			}
		}
		acls = append(acls, &acl)
		if err != nil {
			return nil, errorp.NewPolicyError(500, "ERR_FETCH_POLICY", "error_fetch_scan_policy", err.Error())
		}
	}
	if len(acls) == 0 {
		return nil, errorp.NewPolicyError(404, "ERR_FETCH_POLICY", "error_not_found_policy", "No data")
	}
	return acls, nil
}
func (m *MYSQLManager) DeletePolicy(ctx context.Context, policeID string) error {
	sp := startSpan(ctx, "delete-policy", queries.DeletePolicy)
	defer finishSpan(sp)
	stmt, err := m.Prepare(queries.DeletePolicy)
	if err != nil {
		return err
	}
	result, err := stmt.Exec(policeID)
	if err != nil {
		return errorp.NewPolicyError(500, "ERR_DELETE_POLICY", "error_delete_policy", err.Error())
	}
	if affected, err := result.RowsAffected(); affected == 0 || err != nil {
		return errorp.NewPolicyError(404, "ERR_DELETE_POLICY", "error_delete_policy", "No data to Delete")
	}
	return nil
}

func (m *MYSQLManager) GetResources(ctx context.Context) ([]string, error) {
	sp := startSpan(ctx, "fetch-resources", queries.GetResources)
	defer finishSpan(sp)
	result := []string{}
	r, err := m.DBService().Query(queries.GetResources)
	if err != nil {
		return nil, errorp.NewPolicyError(500, "ERR_FETCH_RESOURCES", "error_fetch_query_resources", err.Error())
	}
	defer r.Close()
	for r.Next() {
		resource := ""
		err = r.Scan(&resource)
		if err != nil {
			return nil, errorp.NewPolicyError(500, "ERR_FETCH_RESOURCES", "error_fetch_scan_resources", err.Error())
		}
		result = append(result, resource)
	}
	if len(result) == 0 {
		return nil, errorp.NewPolicyError(404, "ERR_FETCH_POLICY", "error_not_found_policy", "No data")
	}
	return result, nil
}

func (m *MYSQLManager) GetPolicySubjects(ctx context.Context) ([]string, error) {
	sp := startSpan(ctx, "fetch-policy-subjects", queries.GetPolicySubjects)
	defer finishSpan(sp)
	result := []string{}
	r, err := m.DBService().Query(queries.GetPolicySubjects)
	if err != nil {
		return nil, errorp.NewPolicyError(500, "ERR_FETCH_SUBJECTS", "error_fetch_query_subjects", err.Error())
	}
	defer r.Close()
	for r.Next() {
		subject := ""
		err = r.Scan(&subject)
		if err != nil {
			return nil, errorp.NewPolicyError(500, "ERR_FETCH_SUBJECTS", "error_fetch_scan_subjects", err.Error())
		}
		result = append(result, subject)
	}
	if len(result) == 0 {
		return nil, errorp.NewPolicyError(404, "ERR_FETCH_POLICY", "error_not_found_policy", "No data")
	}
	return result, nil
}

func (m *MYSQLManager) Enforce(ctx context.Context, tenant, subject, action, resource string) (*police.ACL, error) {
	wSub, vSub := hStringWhereQuery("subject", subject)
	wSub = wSub + " OR subject in ( SELECT policy from roles WHERE subject = ? )"
	vSub = append(vSub, subject)
	wAct, vAct := hStringWhereQuery("action", action)
	wRes, vRes := hStringWhereQuery("resource", resource)
	values := []interface{}{tenant}
	values = append(values, vSub...)
	values = append(values, vAct...)
	values = append(values, vRes...)
	values = append(values, []interface{}{"0", "1000"}...)
	where := strings.Join([]string{"( " + wSub + " )", "( " + wAct + " )", "( " + wRes + " )"}, " AND ")
	q := strings.ReplaceAll(queries.GetPolicy, ":where", "WHERE tenant = ? AND "+where)
	sp := startSpan(ctx, "fetch-enforcement", q)
	defer finishSpan(sp)

	r, err := m.DBService().Query(q, values...)
	if err != nil {
		return nil, errorp.NewPolicyError(500, "ERR_ENFORCE_POLICY", "error_enforce_policy", err.Error())
	}
	defer r.Close()
	r.Next()
	acl := police.ACL{}
	err = r.Scan(
		&acl.ID,
		&acl.Subject,
		&acl.Tentant,
		&acl.Resource,
		&acl.Action,
		&acl.Effect,
		&acl.Active,
		&acl.Expired,
		&acl.Created,
		&acl.Updated,
	)

	if err != nil {
		return nil, errorp.NewPolicyError(401, "ERR_FETCH_POLICY", "error_not_found_policy", "No Policy")
	}
	return &acl, nil
}

func (m *MYSQLManager) GetRolePolicy(ctx context.Context, tenant, subject string) (*police.ACL, error) {
	sp := startSpan(ctx, "fetch-role-policy", tenant+":"+subject)
	defer finishSpan(sp)
	acls, err := m.FetchPolicy(ctx, [][3]string{{"tenant", "=", tenant}, {"subject", "=", subject}}...)
	if err != nil {
		return nil, err
	}
	return acls[0], nil
}

func (m *MYSQLManager) UpdatePolicy(ctx context.Context, policeID string, acl *police.ACL) error {
	sp := startSpan(ctx, "updatePolicy", queries.UpdatePolicy)
	defer finishSpan(sp)
	stmt, err := m.Prepare(queries.UpdatePolicy)
	if err != nil {
		return err
	}
	acl.ID = policeID
	r, err := stmt.Exec(acl.Subject, acl.Tentant, acl.Resource, acl.Action, acl.Effect, acl.Active, acl.Expired, acl.Updated, acl.ID)
	if acl.Expired.IsZero() {
		acl.Expired = nil
	}
	if err != nil {
		return errorp.NewPolicyError(500, "ERR_UPDATE_POLICY", "error_update_policy", err.Error())
	}
	if affected, err := r.RowsAffected(); affected == 0 || err != nil {
		return errorp.NewPolicyError(404, "ERR_UPDATE_POLICY", "error_update_policy", "No data to Delete")
	}
	acl, _ = m.GetPolicy(ctx, policeID)
	return nil
}

func hStringWhereQuery(field, str string) (query string, v []interface{}) {

	res := []string{}
	hstring := police.ConverHString(str)
	pow := math.Pow(2, float64(len(hstring)))

	for i := 0; i < int(pow); i++ {

		bin := strLvar(len(hstring), strconv.FormatInt(int64(i), 2))
		c := strings.Split(bin, "")
		for j := 0; j < len(hstring); j++ {
			if c[j] == "0" {
				c[j] = "*"
			} else {
				c[j] = hstring[j]
			}
		}
		lastIndex := len(c)
		for j := len(c) - 1; j >= 0; j-- {
			if c[j] == "*" {
				lastIndex = j + 1
			} else {
				j = -1
			}
		}
		c = c[0:lastIndex]
		v = append(v, strings.Join(c, police.DELIMITER_HSTRING))
		res = append(res, field+" = ? ")
	}
	query = strings.Join(res, " OR ")
	return
}

func strLvar(l int, str string) string {
	for i := len(str); i < l; i++ {
		str = "0" + str
	}
	return str
}
