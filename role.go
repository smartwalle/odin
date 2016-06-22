package odin

type Role struct {
	Name       string
	Identifier string
	RoleGroup  *RoleGroup
}

type RoleGroup struct {
	Name       string
	Identifier string
}
