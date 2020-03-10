package police

import "context"

type Manager interface {
	UpsertPolicy(ctx context.Context, acl *ACL) error
	GetPolicy(ctx context.Context, policeID string) (*ACL, error)
	FetchPolicy(ctx context.Context, filter ...[3]string) ([]*ACL, error)
	DeletePolicy(ctx context.Context, policeID string) error
	UpdatePolicy(ctx context.Context, policeID string, acl *ACL) error
	GetResources(ctx context.Context) ([]string, error)
	GetPolicySubjects(ctx context.Context) ([]string, error)

	UpsertRole(ctx context.Context, rbac *RBAC) error
	GetRoles(ctx context.Context, tenant string) ([]string, error)
	GetRoleSubjects(ctx context.Context, tenant, policy string) ([]string, error)
	GetSubjectRoles(ctx context.Context, tenant, subject string) ([]string, error)
	DeleteRole(ctx context.Context, tenant, subject, policy string) error

	GetRolePolicy(ctx context.Context, tenant, subject string) (*ACL, error)

	Enforce(ctx context.Context, tenant, subject, action, resource string) (*ACL, error)
}
