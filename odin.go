package odin

import "time"

type Status int

const (
	Enable  Status = 1
	Disable Status = 2
)

type GroupType int

const (
	GroupPermission GroupType = 1
	GroupRole       GroupType = 2
)

type Group struct {
	Id        int64      `json:"id"              sql:"id"`
	Ctx       int64      `json:"ctx"             sql:"ctx"`
	Type      GroupType  `json:"type"            sql:"type"`
	Name      string     `json:"name"            sql:"name"`
	AliasName string     `json:"alias_name"      sql:"alias_name"`
	Status    Status     `json:"status"          sql:"status"`
	CreatedOn *time.Time `json:"created_on"      sql:"created_on"`
	UpdatedOn *time.Time `json:"updated_on"      sql:"updated_on"`
}

type Permission struct {
	Id          int64      `json:"id"              sql:"id"`
	GroupId     int64      `json:"group_id"        sql:"group_id"`
	Ctx         int64      `json:"ctx"             sql:"ctx"`
	Name        string     `json:"name"            sql:"name"`
	AliasName   string     `json:"alias_name"      sql:"alias_name"`
	Status      Status     `json:"status"          sql:"status"`
	Description string     `json:"description"     sql:"description"`
	Granted     bool       `json:"granted"         sql:"granted"`
	CreatedOn   *time.Time `json:"created_on"      sql:"created_on"`
	UpdatedOn   *time.Time `json:"updated_on"      sql:"updated_on"`
}

type Role struct {
	Id             int64         `json:"id"                              sql:"id"`
	Ctx            int64         `json:"ctx"                             sql:"ctx"`
	Name           string        `json:"name"                            sql:"name"`
	AliasName      string        `json:"alias_name"                      sql:"alias_name"`
	Status         Status        `json:"status"                          sql:"status"`
	Description    string        `json:"description"                     sql:"description"`
	Granted        bool          `json:"granted"                         sql:"granted"`
	ParentId       int64         `json:"parent_id"                       sql:"parent_id"`
	CreatedOn      *time.Time    `json:"created_on"                      sql:"created_on"`
	UpdatedOn      *time.Time    `json:"updated_on"                      sql:"updated_on"`
	PermissionList []*Permission `json:"permission_list,omitempty"       sql:"-"`
}

type RolePermission struct {
	Ctx          int64      `json:"ctx"             sql:"ctx"`
	RoleId       int64      `json:"role_id"         sql:"role_id"`
	PermissionId int64      `json:"permission_id"   sql:"permission_id"`
	CreatedOn    *time.Time `json:"created_on"      sql:"created_on"`
}

type Grant struct {
	Ctx       int64      `json:"ctx"             sql:"ctx"`
	RoleId    int64      `json:"role_id"         sql:"role_id"`
	TargetId  string     `json:"target_id"       sql:"target_id"`
	CreatedOn *time.Time `json:"created_on"      sql:"created_on"`
}

type Service interface {
	// group

	// GetPermissionGroupList 获取权限组列表
	GetPermissionGroupList(ctx int64, status Status, keywords string) (result []*Group, err error)

	// GetPermissionGroupWithId 根据组 id 获取组信息
	GetPermissionGroupWithId(ctx, groupId int64) (result *Group, err error)

	// GetPermissionGroupWithName 根据组名称获取组信息
	GetPermissionGroupWithName(ctx int64, gType GroupType, name string) (result *Group, err error)

	// AddPermissionGroup 添加组信息
	AddPermissionGroup(ctx int64, name, aliasName string, status Status) (result int64, err error)

	// UpdatePermissionGroup 更新组信息
	UpdatePermissionGroup(ctx, groupId int64, name, aliasName string, status Status) (err error)

	// UpdatePermissionGroupWithName 更新组信息
	UpdatePermissionGroupWithName(ctx int64, groupName string, name, aliasName string, status Status) (err error)

	// UpdatePermissionGroupStatus 更新组状态
	UpdatePermissionGroupStatus(ctx int64, groupId int64, status Status) (err error)

	// UpdatePermissionGroupStatusWithName 更新组状态
	UpdatePermissionGroupStatusWithName(ctx int64, groupName string, status Status) (err error)

	// permission

	// GetPermissionList 获取权限列表
	GetPermissionList(ctx, groupId int64, status Status, keywords string) (result []*Permission, err error)

	// GetPermissionWithId 根据权限 id 获取权限信息
	GetPermissionWithId(ctx, permissionId int64) (result *Permission, err error)

	// GetPermissionWithName 根据权限名称获取权限信息
	GetPermissionWithName(ctx int64, name string) (result *Permission, err error)

	// AddPermission 添加权限
	AddPermission(ctx int64, groupName, name, aliasName, description string, status Status) (result int64, err error)

	// UpdatePermission 更新权限
	UpdatePermission(ctx, permissionId int64, groupName string, name, aliasName, description string, status Status) (err error)

	// UpdatePermissionStatus 更新权限状态
	UpdatePermissionStatus(ctx, permissionId int64, status Status) (err error)

	// GrantPermission 授权权限给 roleId
	GrantPermission(ctx int64, roleId int64, names ...string) (err error)

	// GrantPermissionWithIds 授权权限给 roleId
	GrantPermissionWithIds(ctx int64, roleId int64, permissionIds ...int64) (err error)

	// ReGrantPermission 授权权限给 roleId，会将原有的权限授权先取消掉
	ReGrantPermission(ctx int64, roleId int64, names ...string) (err error)

	// ReGrantPermissionWithIds 授权权限给 roleId，会将原有的权限授权先取消掉
	ReGrantPermissionWithIds(ctx int64, roleId int64, permissionIds ...int64) (err error)

	// RevokePermission 取消对 roleId 的权限授权
	RevokePermission(ctx int64, roleId int64, names ...string) (err error)

	// RevokePermissionWithIds 取消对 roleId 的权限授权
	RevokePermissionWithIds(ctx int64, roleId int64, permissionIds ...int64) (err error)

	// RevokeAllPermission 取消对 roleId 的所有权限授权
	RevokeAllPermission(ctx, roleId int64) (err error)

	// role

	// GetRoleList 获取角色列表，如果有传递 targetId 参数，则返回的角色数据中将附带该角色是否已授权给该 targetId
	GetRoleList(ctx int64, targetId string, status Status, keywords string) (result []*Role, err error)

	// GetRoleWithId 根据角色 id 获取角色信息
	GetRoleWithId(ctx, roleId int64) (result *Role, err error)

	// GetRoleWithName 根据角色名称获取角色信息
	GetRoleWithName(ctx int64, name string) (result *Role, err error)

	// AddRole 添加角色，如果 parentId 参数的值大于 0，则会验证 parent 是否存在
	AddRole(ctx, parentId int64, name, aliasName, description string, status Status) (result int64, err error)

	// UpdateRole 更新角色信息
	UpdateRole(ctx, roleId int64, name, aliasName, description string, status Status) (err error)

	// UpdateRoleStatus 更新角色的状态
	UpdateRoleStatus(ctx, roleId int64, status Status) (err error)

	// GetGrantedRoleList 获取已授权给 targetId 的角色列表
	GetGrantedRoleList(ctx int64, targetId string) (result []*Role, err error)

	// GrantRole 授权角色给 targetId
	GrantRole(ctx int64, targetId string, names ...string) (err error)

	// GrantRoleWithIds 授权角色给 targetId
	GrantRoleWithIds(ctx int64, targetId string, roleIds ...int64) (err error)

	// ReGrantRole 授权角色给 targetId，会将原有的角色授权先取消掉
	ReGrantRole(ctx int64, targetId string, names ...string) (err error)

	// ReGrantRoleWithIds 授权角色给 targetId，会将原有的角色授权先取消掉
	ReGrantRoleWithIds(ctx int64, targetId string, roleIds ...int64) (err error)

	// RevokeRole 取消对 targetId 的角色授权
	RevokeRole(ctx int64, targetId string, names ...string) (err error)

	// RevokeRoleWithIds 取消对 targetId 的角色授权
	RevokeRoleWithIds(ctx int64, targetId string, roleIds ...int64) (err error)

	// RevokeAllRole 取消对 targetId 的所有角色授权
	RevokeAllRole(ctx int64, targetId string) (err error)
}
