package odin

const (
	k_ODIN_GROUP_TYPE_PERMISSION = "1" // 权限组
	k_ODIN_GROUP_TYPE_ROLE       = "2" // 角色组
)

// odin_gps_[type] 用于存储 Group 的 Id 列表

// Group 用于存储分组信息，不参与权限的控制
// 特殊限定：创建之后，只能修改 name，不能修改 type
// Key：odin_gp_[id]
type Group struct {
	Id   string `json:"id"     redis:"id"`   // id
	Type string `json:"type"   redis:"type"` // 组类型：1、权限组；2、角色组
	Name string `json:"name"   redis:"name"` // 组的名字
}

// Permission 用于存储权限信息
// Key：odin_pn_[id]
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
	DestinationId string   `json:"destination_id"   redis:"destination_id"`
	RoleIdList    []string `json:"role_id_list"     redis:"-"`
}
