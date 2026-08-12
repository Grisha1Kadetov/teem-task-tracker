package role

import "fmt"

type Role string

const (
	Owner  Role = "owner"
	Admin  Role = "admin"
	Member Role = "member"
)

func (r Role) IsValid() bool {
	switch r {
	case Owner, Admin, Member:
		return true
	default:
		return false
	}
}

func Parse(value string) (Role, error) {
	r := Role(value)
	if !r.IsValid() {
		return "", fmt.Errorf("unknown team role %q", value)
	}

	return r, nil
}
