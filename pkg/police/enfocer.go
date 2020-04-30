package police

import (
	"time"

	"github.com/oeoen/policy/helper/errorp"
)

func (p *ACL) Check() error {
	if !p.Active {
		return errorp.NewPolicyError(401, "ERROR_POLICY_NOT_ACTIVE", "Policy has been blocked", "Blocked access")
	}
	if p.Effect != "allow" {
		return errorp.NewPolicyError(401, "ERROR_POLICY_NOT_ALLOWED", "Policy not allowed", "Blocked access")
	}
	if p.Expired == nil {
		return nil
	}
	if p.Expired.IsZero() {
		return nil
	}
	now := time.Now().UnixNano()
	if p.Expired.UnixNano() < now {
		return errorp.NewPolicyError(401, "ERROR_POLICY_EXPIRED", "Policy expired", "Blocked access")
	}
	return nil
}
