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
	Id             int64         `json:"id"                              sql:"id"`
	Ctx            int64         `json:"ctx"                             sql:"ctx"`
	Type           GroupType     `json:"type"                            sql:"type"`
	Name           string        `json:"name"                            sql:"name"`
	AliasName      string        `json:"alias_name"                      sql:"alias_name"`
	Status         Status        `json:"status"                          sql:"status"`
	CreatedOn      *time.Time    `json:"created_on"                      sql:"created_on"`
	UpdatedOn      *time.Time    `json:"updated_on"                      sql:"updated_on"`
	PermissionList []*Permission `json:"permission_list,omitempty"       sql:"-"`
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
	Children       []*Role       `json:"children,omitempty"`
}

type RolePermission struct {
	Ctx          int64      `json:"ctx"             sql:"ctx"`
	RoleId       int64      `json:"role_id"         sql:"role_id"`
	PermissionId int64      `json:"permission_id"   sql:"permission_id"`
	CreatedOn    *time.Time `json:"created_on"      sql:"created_on"`
}

type Grant struct {
	Ctx            int64  `json:"ctx"                sql:"ctx"`
	TargetId       string `json:"target_id"          sql:"target_id"`
	RoleId         int64  `json:"role_id"            sql:"role_id"`
	RoleName       string `json:"role_name"          sql:"role_name"`
	PermissionId   int64  `json:"permission_id"      sql:"permission_id"`
	PermissionName string `json:"permission_name"    sql:"permission_name"`
}

type Service interface {
	// Init 执行初始化操作，目前主要功能为初始化数据库表。
	// 虽然此方法可以被重复调用，但是外部应该尽量控制此方法只在需要的时候调用。
	Init() error

	// group

	// GetPermissionGroups 获取权限组列表
	GetPermissionGroups(ctx int64, status Status, keywords string) (result []*Group, err error)

	// GetPermissionGroupWithId 根据组 id 获取组信息
	GetPermissionGroupWithId(ctx, groupId int64) (result *Group, err error)

	// GetPermissionGroup 根据 groupName 获取组信息
	GetPermissionGroup(ctx int64, groupName string) (result *Group, err error)

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

	// GetPermissions 获取权限列表
	GetPermissions(ctx int64, status Status, keywords string, groupIds []int64) (result []*Permission, err error)

	// GetPermissionWithId 根据 permissionId 获取权限信息
	GetPermissionWithId(ctx, permissionId int64) (result *Permission, err error)

	// GetPermission 根据 permissionName 获取权限信息
	GetPermission(ctx int64, permissionName string) (result *Permission, err error)

	// AddPermissionWithGroupId 添加权限
	AddPermissionWithGroupId(ctx, groupId int64, permissionName, aliasName, description string, status Status) (result int64, err error)

	// AddPermissionWithGroup 添加权限
	AddPermissionWithGroup(ctx int64, groupName, permissionName, aliasName, description string, status Status) (result int64, err error)

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

	// GrantPermission 授予权限给角色
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

	// GetRoles 获取角色列表
	// 如果参数 isGrantedToTarget 的值不为空字符串，则返回的角色数据中将包含该角色（通过 Granted 判断）是否已授权给 isGrantedToTarget
	GetRoles(ctx int64, status Status, keywords string, isGrantedToTarget string) (result []*Role, err error)

	// GetRolesTreeWithParentId 获取角色树
	GetRolesTreeWithParentId(ctx, parentRoleId int64, status Status, keywords string) (result []*Role, err error)

	// GetRolesTreeWithParent 获取角色树
	GetRolesTreeWithParent(ctx int64, parentRoleName string, status Status, keywords string) (result []*Role, err error)

	// GetRoleWithId 根据 roleId 获取角色信息
	GetRoleWithId(ctx, roleId int64) (result *Role, err error)

	// GetRole 根据 roleName 获取角色信息
	GetRole(ctx int64, roleName string) (result *Role, err error)

	// AddRole 添加角色
	AddRole(ctx int64, roleName, aliasName, description string, status Status) (result int64, err error)

	// AddRoleWithParentId 添加角色，新添加的角色将作为 parentRoleId 的子角色，调用时应该确认操作者是否有访问 parentRoleId 的权限
	AddRoleWithParentId(ctx, parentRoleId int64, roleName, aliasName, description string, status Status) (result int64, err error)

	// AddRoleWithParent 添加角色，新添加的角色将作为 parentRoleName 的子角色，调用时应该确认操作者是否有访问 parentRoleName 的权限
	AddRoleWithParent(ctx int64, parentRoleName, roleName, aliasName, description string, status Status) (result int64, err error)

	// UpdateRoleWithId 根据 roleId 更新角色信息
	UpdateRoleWithId(ctx, roleId int64, aliasName, description string, status Status) (err error)

	// UpdateRole 根据 roleName 更新角色信息
	UpdateRole(ctx int64, roleName, aliasName, description string, status Status) (err error)

	// UpdateRoleStatusWithId 根据 roleId 更新角色的状态
	UpdateRoleStatusWithId(ctx, roleId int64, status Status) (err error)

	// UpdateRoleStatus 根据 roleName 更新角色的状态
	UpdateRoleStatus(ctx int64, roleName string, status Status) (err error)

	// GrantRoleWithId 授权角色给 targetId
	GrantRoleWithId(ctx int64, targetId string, roleIds ...int64) (err error)

	// GrantRole 授权角色给 targetId
	GrantRole(ctx int64, targetId string, roleNames ...string) (err error)

	// ReGrantRoleWithId 授权角色给 targetId，会将原有的角色授权先取消掉
	ReGrantRoleWithId(ctx int64, targetId string, roleIds ...int64) (err error)

	// ReGrantRole 授权角色给 targetId，会将原有的角色授权先取消掉
	ReGrantRole(ctx int64, targetId string, roleNames ...string) (err error)

	// RevokeRoleWithId 取消对 targetId 的角色授权
	RevokeRoleWithId(ctx int64, targetId string, roleIds ...int64) (err error)

	// RevokeRole 取消对 targetId 的角色授权
	RevokeRole(ctx int64, targetId string, roleNames ...string) (err error)

	// RevokeAllRole 取消对 targetId 的所有角色授权
	RevokeAllRole(ctx int64, targetId string) (err error)

	// 其它

	// GetPermissionsWithRoleId 获取已授权给 roleId 的权限列表
	GetPermissionsWithRoleId(ctx int64, roleId int64) (result []*Permission, err error)

	// GetPermissionsWithRole 获取已授权给 roleName 的权限列表
	GetPermissionsWithRole(ctx int64, roleName string) (result []*Permission, err error)

	// GetGrantedRoles 获取已授权给 targetId 的角色列表
	GetGrantedRoles(ctx int64, targetId string) (result []*Role, err error)

	// GetRolesWithTargetId 获取已授权给 targetId 的角色列表，与方法 GetGrantedRoles 作用相同
	GetRolesWithTargetId(ctx int64, targetId string) (result []*Role, err error)

	// GetGrantedPermissions 获取已授权给 targetId 的权限列表
	GetGrantedPermissions(ctx int64, targetId string) (result []*Permission, err error)

	// GetPermissionsWithTargetId 获取已授权给 targetId 的权限列表，与方法 GetGrantedPermissions 作用相同
	GetPermissionsWithTargetId(ctx int64, targetId string) (result []*Permission, err error)

	// GetPermissionsTreeWithRoleId 获取权限组列表，组中包含该组所有的权限信息
	// 如果参数 roleId 的值大于 0，则返回的权限数据中将附带该权限是否已授权给该 roleId
	// 如果参数 limitedInParentRole 的值为 true，并且 roleId 的值大于 0，则返回的权限数据将限定在 roleId 的父角色拥有的权限范围内
	GetPermissionsTreeWithRoleId(ctx, roleId int64, status Status, limitedInParentRole bool) (result []*Group, err error)

	// GetPermissionsTreeWithRole 获取权限组列表，组中包含该组所有的权限信息
	// 如果参数 roleName 的值不为空字符串，则返回的权限数据中将附带该权限是否已授权给该 roleName
	// 如果参数 limitedInParentRole 的值为 true，并且 roleName 的值不为空字符串，则返回的权限数据将限定在 roleName 的父角色拥有的权限范围内
	GetPermissionsTreeWithRole(ctx int64, roleName string, status Status, limitedInParentRole bool) (result []*Group, err error)

	// CheckPermission 验证 targetId 是否拥有指定权限
	CheckPermission(ctx int64, targetId string, permissionName string) bool

	// CheckPermissionWithId 验证 targetId 是否拥有指定权限
	CheckPermissionWithId(ctx int64, targetId string, permissionId int64) bool

	// CheckRole 验证 targetId 是否拥有指定角色
	CheckRole(ctx int64, targetId string, roleName string) bool

	// CheckRoleWithId 验证 targetId 是否拥有指定角色
	CheckRoleWithId(ctx int64, targetId string, roleId int64) bool

	// CleanCache 清除缓存，如果 targetId 为空字符串或者 targetId 的值为星号(*)，则会清空所有缓存
	CleanCache(ctx int64, targetId string)
}
