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

	// GetPermissionGroup 根据 groupName 获取组信息
	GetPermissionGroup(ctx int64, gType GroupType, groupName string) (result *Group, err error)

	// AddPermissionGroup 添加组信息
	AddPermissionGroup(ctx int64, groupName, aliasName string, status Status) (result int64, err error)

	// UpdatePermissionGroupWithId 根据 groupId 更新组信息
	UpdatePermissionGroupWithId(ctx, groupId int64, aliasName string, status Status) (err error)

	// UpdatePermissionGroup 根据 groupName 更新组信息
	UpdatePermissionGroup(ctx int64, groupName string, aliasName string, status Status) (err error)

	// UpdatePermissionGroupStatusWithId 根据 groupId 更新组状态
	UpdatePermissionGroupStatusWithId(ctx int64, groupId int64, status Status) (err error)

	// UpdatePermissionGroupStatus 根据 groupName 更新组状态
	UpdatePermissionGroupStatus(ctx int64, groupName string, status Status) (err error)

	// permission

	// GetPermissionList 获取权限列表
	GetPermissionList(ctx, groupId int64, status Status, keywords string) (result []*Permission, err error)

	// GetPermissionWithId 根据 permissionId 获取权限信息
	GetPermissionWithId(ctx, permissionId int64) (result *Permission, err error)

	// GetPermission 根据 permissionName 获取权限信息
	GetPermission(ctx int64, permissionName string) (result *Permission, err error)

	// AddPermissionWithGroupId 添加权限
	AddPermissionWithGroupId(ctx, groupId int64, permissionName, aliasName, description string, status Status) (result int64, err error)

	// AddPermission 添加权限
	AddPermission(ctx int64, groupName, permissionName, aliasName, description string, status Status) (result int64, err error)

	// UpdatePermissionWithId 根据 permissionId 更新权限信息
	UpdatePermissionWithId(ctx, permissionId, groupId int64, aliasName, description string, status Status) (err error)

	// UpdatePermission 根据 permissionName 更新权限信息
	UpdatePermission(ctx int64, permissionName, groupName, aliasName, description string, status Status) (err error)

	// UpdatePermissionStatusWithId 根据 permissionId 更新权限状态
	UpdatePermissionStatusWithId(ctx, permissionId int64, status Status) (err error)

	// UpdatePermissionStatus 根据 permissionName 更新权限状态
	UpdatePermissionStatus(ctx int64, permissionName string, status Status) (err error)

	// GrantPermissionWithId 授予权限给角色
	GrantPermissionWithId(ctx int64, roleId int64, permissionIds ...int64) (err error)

	// GrantPermissionWithId 授予权限给角色
	GrantPermission(ctx int64, roleName string, permissionNames ...string) (err error)

	// ReGrantPermissionWithId 授予权限给角色，会将原有的权限先取消掉
	ReGrantPermissionWithId(ctx int64, roleId int64, permissionIds ...int64) (err error)

	// ReGrantPermission 授予权限给角色，会将原有的权限先取消掉
	ReGrantPermission(ctx int64, roleName string, permissionNames ...string) (err error)

	// RevokePermissionWithId 取消对角色的指定权限授权
	RevokePermissionWithId(ctx int64, roleId int64, permissionIds ...int64) (err error)

	// RevokePermission 取消对角色的指定权限授权
	RevokePermission(ctx int64, roleName string, permissionNames ...string) (err error)

	// RevokeAllPermissionWithId 取消对角色的所有权限授权
	RevokeAllPermissionWithId(ctx, roleId int64) (err error)

	// RevokeAllPermission 取消对角色的所有权限授权
	RevokeAllPermission(ctx int64, roleName string) (err error)

	// role

	// GetRoleList 获取角色列表，如果有传递 targetId 参数，则返回的角色数据中将附带该角色是否已授权给该 targetId
	GetRoleList(ctx int64, targetId string, status Status, keywords string) (result []*Role, err error)

	// GetRoleWithId 根据 roleId 获取角色信息
	GetRoleWithId(ctx, roleId int64) (result *Role, err error)

	// GetRole 根据 roleName 获取角色信息
	GetRole(ctx int64, roleName string) (result *Role, err error)

	// AddRole 添加角色，如果 parentName 参数的值大于 0，则会验证 parent 是否存在
	AddRoleWithParentId(ctx, parentId int64, roleName, aliasName, description string, status Status) (result int64, err error)

	// AddRole 添加角色，如果 parentName 不为空字符串，则会验证 parent 是否存在
	AddRole(ctx int64, parentName, roleName, aliasName, description string, status Status) (result int64, err error)

	// UpdateRoleWithId 根据 roleId 更新角色信息
	UpdateRoleWithId(ctx, roleId int64, aliasName, description string, status Status) (err error)

	// UpdateRole 根据 roleName 更新角色信息
	UpdateRole(ctx int64, roleName, aliasName, description string, status Status) (err error)

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
