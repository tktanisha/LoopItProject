package enums

import (
	"encoding/json"
	"fmt"
)

type Role int

const (
	RoleUser Role = iota
	RoleLender
	RoleAdmin
)

var roleNames = [...]string{
	"user",
	"lender",
	"admin",
}

func (r Role) String() string {
	if r < 0 || int(r) >= len(roleNames) {
		return "unknown"
	}
	return roleNames[r]
}

func ParseRole(roleStr string) (Role, error) {
	for i, v := range roleNames {
		if v == roleStr {
			return Role(i), nil
		}
	}
	return RoleUser, fmt.Errorf("invalid role: %s", roleStr)
}

func (r Role) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

func (r *Role) UnmarshalJSON(data []byte) error {
	var roleStr string
	if err := json.Unmarshal(data, &roleStr); err != nil {
		return err
	}
	parsed, err := ParseRole(roleStr)
	if err != nil {
		return err
	}
	*r = parsed
	return nil
}

// Role.string() => Enum

// Role.string() => "user", "lender", "admin"
// Role.ParseRole("user") => Role(0) => RoleUser
// Role.ParseRole("lender") => Role(1) => RoleLender
// Role.ParseRole("admin") => Role(2) => RoleAdmin

// enum.Role => json ????  -> role.string()
// "user" => parserole("user") => RoleUser Role(0)
