package odin

type Permission struct {
	Id         string `json:"id"           redis:"id"`         // 将 identifier 进行 md5 编码
	Identifier string `json:"identifier"   redis:"identifier"` // 权限标识
	Name       string `json:"name"         redis:"name"`       // 权限名称
	Group      string `json:"group"        redis:"group"`      // 组名
}

type Role struct {
	Id               string   `json:"id"                   redis:"id"`
	Name             string   `json:"name"                 redis:"name"`
	Group            string   `json:"group"                redis:"group"` // 组名
	PermissionIdList []string `json:"permission_id_list"   redis:"-"`
}

type GrantInfo struct {
	DestinationId  string   `json:"destination_id"             redis:"destination_id"`
	RoleIdList     []string `json:"role_id_list,omitempty"     redis:"-"`
	PermissionList []string `json:"permission_list,omitempty"  redis:"-"`
}
