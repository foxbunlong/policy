package police

import "time"

type ACL struct {
	ID       string     `json:"id"`
	Tentant  string     `json:"tenant"`
	Subject  string     `json:"subject"`
	Resource string     `json:"resource"`
	Action   string     `json:"action"`
	Effect   string     `json:"effect"`
	Active   bool       `json:"active"`
	Expired  *time.Time `json:"expired,omitempty"`
	Created  time.Time  `json:"created,omitempty"`
	Updated  time.Time  `json:"updated,omitempty"`
}

type RBAC struct {
	ID      int64     `json:"-"`
	Policy  string    `json:"policy"`
	Tenant  string    `json:"tenant"`
	Subject string    `json:"subject"`
	Created time.Time `json:"created,omitempty"`
}
