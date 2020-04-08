package odin

import "time"

// Status 状态信息。
type Status int

const (
	Enable  Status = 1 // 启用
	Disable Status = 2 // 禁用
)

// GroupType 组的类型，目前分为权限组和角色组，组没有实质的意义，主要是对权限数据或者角色数据进行分类， 目前只实现了权限组的管理。
type GroupType int

const (
	GroupPermission GroupType = 1 // 权限组
	GroupRole       GroupType = 2 // 角色组
)

// Group 组数据结构，用于描述组信息。
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

// Permission 权限数据结构，用于描述权限信息。
type Permission struct {
	Id                int64            `json:"id"                               sql:"id"`
	GroupId           int64            `json:"group_id"                         sql:"group_id"`
	Ctx               int64            `json:"ctx"                              sql:"ctx"`
	Name              string           `json:"name"                             sql:"name"`
	AliasName         string           `json:"alias_name"                       sql:"alias_name"`
	Status            Status           `json:"status"                           sql:"status"`
	Description       string           `json:"description"                      sql:"description"`
	Granted           bool             `json:"granted"                          sql:"granted"` // 权限是否授予给指定角色
	CreatedOn         *time.Time       `json:"created_on"                       sql:"created_on"`
	UpdatedOn         *time.Time       `json:"updated_on"                       sql:"updated_on"`
	PrePermissionList []*PrePermission `json:"pre_permission_list,omitempty"    sql:"-"`
}

// Role 角色数据结构，用于描述角色信息。
type Role struct {
	Id             int64         `json:"id"                              sql:"id"`
	Ctx            int64         `json:"ctx"                             sql:"ctx"`
	Name           string        `json:"name"                            sql:"name"`
	AliasName      string        `json:"alias_name"                      sql:"alias_name"`
	Status         Status        `json:"status"                          sql:"status"`
	Description    string        `json:"description"                     sql:"description"`
	Granted        bool          `json:"granted"                         sql:"granted"`    // 角色是否授予给指定 target
	Accessible     bool          `json:"accessible"                      sql:"can_access"` // 角色是否能够被指定 target 操作访问
	ParentId       int64         `json:"parent_id"                       sql:"parent_id"`
	LeftValue      int64         `json:"left_value"                      sql:"left_value"`
	RightValue     int64         `json:"right_value"                     sql:"right_value"`
	Depth          int           `json:"depth"                           sql:"depth"`
	CreatedOn      *time.Time    `json:"created_on"                      sql:"created_on"`
	UpdatedOn      *time.Time    `json:"updated_on"                      sql:"updated_on"`
	PermissionList []*Permission `json:"permission_list,omitempty"       sql:"-"`
	MutexRoleList  []*RoleMutex  `json:"mutex_role_list,omitempty"       sql:"-"`
	PreRoleList    []*PreRole    `json:"pre_role_list,omitempty"         sql:"-"`
}

// RolePermission 角色-权限数据结构，用于描述角色与权限的关联关系。
type RolePermission struct {
	Ctx          int64      `json:"ctx"             sql:"ctx"`
	RoleId       int64      `json:"role_id"         sql:"role_id"`
	PermissionId int64      `json:"permission_id"   sql:"permission_id"`
	CreatedOn    *time.Time `json:"created_on"      sql:"created_on"`
}

// RoleMutex 角色互斥数据结构，用于描述角色与角色之间的互斥关系。
// 主要应用于更新（授予）某一 target 的角色时，判断该 target 的角色是否有互斥的情况。
// 比如：角色A与角色B互斥，则角色A和角色B不能同时授予同一 target。
type RoleMutex struct {
	Ctx                int64      `json:"ctx"                       sql:"ctx"`
	RoleId             int64      `json:"role_id"                   sql:"role_id"`
	RoleName           string     `json:"role_name"                 sql:"role_name"`
	RoleAliasName      string     `json:"role_alias_name"           sql:"role_alias_name"`
	MutexRoleId        int64      `json:"mutex_role_id"             sql:"mutex_role_id"`
	MutexRoleName      string     `json:"mutex_role_name"           sql:"mutex_role_name"`
	MutexRoleAliasName string     `json:"mutex_role_alias_name"     sql:"mutex_role_alias_name"`
	CreatedOn          *time.Time `json:"created_on"                sql:"created_on"`
}

// PreRole 角色先决条件数据结构。
// 主要应用于更新（授予）某一 target 的角色时，判断该 target 是否已经拥有某一角色。
//
// 比如：角色A是角色B的先决条件，向 target 授予角色B时，需要 target 已经拥有角色A。
type PreRole struct {
	Ctx              int64      `json:"ctx"                       sql:"ctx"`
	RoleId           int64      `json:"role_id"                   sql:"role_id"`
	RoleName         string     `json:"role_name"                 sql:"role_name"`
	RoleAliasName    string     `json:"role_alias_name"           sql:"role_alias_name"`
	PreRoleId        int64      `json:"pre_role_id"               sql:"pre_role_id"`
	PreRoleName      string     `json:"pre_role_name"             sql:"pre_role_name"`
	PreRoleAliasName string     `json:"pre_role_alias_name"       sql:"pre_role_alias_name"`
	CreatedOn        *time.Time `json:"created_on"                sql:"created_on"`
}

// PrePermission 权限先决条件数据结构。
// 主要应用于更新（授予）某一角色的权限时，判断该角色是否已经拥有某一权限。
//
// 比如：权限A是权限B的先决条件，向角色A授予权限B时，需要角色A已经拥有权限A。
type PrePermission struct {
	Ctx                    int64      `json:"ctx"                           sql:"ctx"`
	PermissionId           int64      `json:"permission_id"                 sql:"permission_id"`
	PermissionName         string     `json:"permission_name"               sql:"permission_name"`
	PermissionAliasName    string     `json:"permission_alias_name"         sql:"permission_alias_name"`
	PrePermissionId        int64      `json:"pre_permission_id"             sql:"pre_permission_id"`
	PrePermissionName      string     `json:"pre_permission_name"           sql:"pre_permission_name"`
	PrePermissionAliasName string     `json:"pre_permission_alias_name"     sql:"pre_permission_alias_name"`
	AutoGrant              bool       `json:"auto_grant"                    sql:"auto_grant"`
	CreatedOn              *time.Time `json:"created_on"                    sql:"created_on"`
}

// Grant 用于描述 target、角色、权限之间的关系。
type Grant struct {
	Ctx            int64  `json:"ctx"                sql:"ctx"`
	Target         string `json:"target"             sql:"target"`
	RoleId         int64  `json:"role_id"            sql:"role_id"`
	RoleName       string `json:"role_name"          sql:"role_name"`
	PermissionId   int64  `json:"permission_id"      sql:"permission_id"`
	PermissionName string `json:"permission_name"    sql:"permission_name"`
}

type Service interface {
	// Init 执行初始化操作，目前主要功能为初始化数据库表。
	// 虽然此方法可以被重复调用，但是外部应该尽量控制此方法只在需要的时候调用。
	Init() error

	// GetPermissionGroups 获取权限组列表
	GetPermissionGroups(ctx int64, status Status, keywords string) (result []*Group, err error)

	// GetPermissionGroup 根据 groupName 获取权限组信息
	GetPermissionGroup(ctx int64, groupName string) (result *Group, err error)

	// GetPermissionGroupWithId 根据组 groupId 获取权限组信息
	GetPermissionGroupWithId(ctx, groupId int64) (result *Group, err error)

	// AddPermissionGroup 添加权限组信息
	AddPermissionGroup(ctx int64, groupName, aliasName string, status Status) (result int64, err error)

	// UpdatePermissionGroup 根据 groupName 更新权限组信息
	UpdatePermissionGroup(ctx int64, groupName string, aliasName string, status Status) (err error)

	// UpdatePermissionGroupWithId 根据 groupId 更新权限组信息
	UpdatePermissionGroupWithId(ctx, groupId int64, aliasName string, status Status) (err error)

	// UpdatePermissionGroupStatus 根据 groupName 更新权限组状态
	UpdatePermissionGroupStatus(ctx int64, groupName string, status Status) (err error)

	// UpdatePermissionGroupStatusWithId 根据 groupId 更新权限组状态
	UpdatePermissionGroupStatusWithId(ctx int64, groupId int64, status Status) (err error)

	// GetPermissions 获取权限列表
	GetPermissions(ctx int64, status Status, keywords string, groupIds []int64) (result []*Permission, err error)

	// GetPermission 根据 permissionName 获取权限信息
	GetPermission(ctx int64, permissionName string) (result *Permission, err error)

	// GetPermissionWithId 根据 permissionId 获取权限信息
	GetPermissionWithId(ctx, permissionId int64) (result *Permission, err error)

	// AddPermissionWithGroup 添加权限
	AddPermissionWithGroup(ctx int64, groupName, permissionName, aliasName, description string, status Status) (result int64, err error)

	// AddPermissionWithGroupId 添加权限
	AddPermissionWithGroupId(ctx, groupId int64, permissionName, aliasName, description string, status Status) (result int64, err error)

	// UpdatePermission 根据 permissionName 更新权限信息
	UpdatePermission(ctx int64, permissionName, groupName, aliasName, description string, status Status) (err error)

	// UpdatePermissionWithId 根据 permissionId 更新权限信息
	UpdatePermissionWithId(ctx, permissionId, groupId int64, aliasName, description string, status Status) (err error)

	// UpdatePermissionStatus 根据 permissionName 更新权限状态
	UpdatePermissionStatus(ctx int64, permissionName string, status Status) (err error)

	// UpdatePermissionStatusWithId 根据 permissionId 更新权限状态
	UpdatePermissionStatusWithId(ctx, permissionId int64, status Status) (err error)

	// GrantPermission 授予权限给角色
	GrantPermission(ctx int64, roleName string, permissionNames ...string) (err error)

	// GrantPermissionWithId 授予权限给角色
	GrantPermissionWithId(ctx int64, roleId int64, permissionIds ...int64) (err error)

	// ReGrantPermission 授予权限给角色，会将原有的权限先取消掉
	ReGrantPermission(ctx int64, roleName string, permissionNames ...string) (err error)

	// ReGrantPermissionWithId 授予权限给角色，会将原有的权限先取消掉
	ReGrantPermissionWithId(ctx int64, roleId int64, permissionIds ...int64) (err error)

	// RevokePermission 取消对角色的指定权限授权
	RevokePermission(ctx int64, roleName string, permissionNames ...string) (err error)

	// RevokePermissionWithId 取消对角色的指定权限授权
	RevokePermissionWithId(ctx int64, roleId int64, permissionIds ...int64) (err error)

	// RevokeAllPermission 取消对角色的所有权限授权
	RevokeAllPermission(ctx int64, roleName string) (err error)

	// RevokeAllPermissionWithId 取消对角色的所有权限授权
	RevokeAllPermissionWithId(ctx, roleId int64) (err error)

	// AddPrePermission 添加授予该权限时需要的先决条件
	AddPrePermission(ctx int64, permissionName string, prePermissionNames ...string) (err error)

	// AddPrePermissionWithId 添加授予该权限时需要的先决条件
	AddPrePermissionWithId(ctx, permissionId int64, prePermissionIds ...int64) (err error)

	// RemovePrePermission 删除授予该权限时需要的先决条件
	RemovePrePermission(ctx int64, permissionName string, prePermissionNames ...string) (err error)

	// RemovePrePermissionWithId 删除授予该权限时需要的先决条件
	RemovePrePermissionWithId(ctx, permissionId int64, prePermissionIds ...int64) (err error)

	// RemoveAllPrePermission 删除授予该权限时需要的所有先决条件
	RemoveAllPrePermission(ctx int64, permissionName string) (err error)

	// RemoveAllPrePermissionWithId 删除授予该权限时需要的所有先决条件
	RemoveAllPrePermissionWithId(ctx, permissionId int64) (err error)

	// GetPrePermissions 获取授予该权限时需要的所有先决条件
	GetPrePermissions(ctx int64, permissionName string) (result []*PrePermission, err error)

	// GetPrePermissionsWithId 获取授予该权限时需要的所有先决条件
	GetPrePermissionsWithId(ctx int64, permissionId int64) (result []*PrePermission, err error)

	// GetRoles 获取角色列表
	// 如果参数 isGrantedToTarget 的值不为空字符串，则返回的角色数据中将包含该角色（通过 Granted 判断）是否已授权给 isGrantedToTarget
	// 如果参数 limitedInTarget 的值不为空字符串， 则返回的角色数据将限定在 limitedInTarget 已拥有的角色及其子角色范围内
	// 返回的角色数据的 Granted 字段参照的是 isGrantedToTarget
	// 返回的角色数据的 Accessible 字段参照的是 limitedInTarget
	GetRoles(ctx int64, status Status, keywords, isGrantedToTarget, limitedInTarget string) (result []*Role, err error)

	// GetRolesWithParent 获取角色列表
	// 如果参数 isGrantedToTarget 的值不为空字符串，则返回的角色数据中将包含该角色（通过 Granted 判断）是否已授权给 isGrantedToTarget
	GetRolesWithParent(ctx int64, parentRoleName string, status Status, keywords, isGrantedToTarget string) (result []*Role, err error)

	// GetRolesWithParentId 获取角色列表
	// 如果参数 isGrantedToTarget 的值不为空字符串，则返回的角色数据中将包含该角色（通过 Granted 判断）是否已授权给 isGrantedToTarget
	GetRolesWithParentId(ctx, parentRoleId int64, status Status, keywords, isGrantedToTarget string) (result []*Role, err error)

	// GetRole 根据 roleName 获取角色信息
	GetRole(ctx int64, roleName string) (result *Role, err error)

	// GetRoleWithId 根据 roleId 获取角色信息
	GetRoleWithId(ctx, roleId int64) (result *Role, err error)

	// AddRole 添加角色
	AddRole(ctx int64, roleName, aliasName, description string, status Status) (result int64, err error)

	// AddRoleWithParent 添加角色，新添加的角色将作为 parentRoleName 的子角色
	// 调用时应该确认操作者是否有访问 parentRoleName 的权限，即 parentRoleName 是否为当前操作者拥有的角色及其子角色
	AddRoleWithParent(ctx int64, parentRoleName, roleName, aliasName, description string, status Status) (result int64, err error)

	// AddRoleWithParentId 添加角色，新添加的角色将作为 parentRoleId 的子角色
	// 调用时应该确认操作者是否有访问 parentRoleId 的权限，即 parentRoleId 是否为当前操作者拥有的角色及其子角色
	AddRoleWithParentId(ctx, parentRoleId int64, roleName, aliasName, description string, status Status) (result int64, err error)

	// UpdateRole 根据 roleName 更新角色信息
	UpdateRole(ctx int64, roleName, aliasName, description string, status Status) (err error)

	// UpdateRoleWithId 根据 roleId 更新角色信息
	UpdateRoleWithId(ctx, roleId int64, aliasName, description string, status Status) (err error)

	// UpdateRoleStatus 根据 roleName 更新角色的状态
	UpdateRoleStatus(ctx int64, roleName string, status Status) (err error)

	// UpdateRoleStatusWithId 根据 roleId 更新角色的状态
	UpdateRoleStatusWithId(ctx, roleId int64, status Status) (err error)

	// GrantRole 授权角色给 target
	GrantRole(ctx int64, target string, roleNames ...string) (err error)

	// GrantRoleWithId 授权角色给 target
	GrantRoleWithId(ctx int64, target string, roleIds ...int64) (err error)

	// ReGrantRole 授权角色给 target，会将原有的角色授权先取消掉
	ReGrantRole(ctx int64, target string, roleNames ...string) (err error)

	// ReGrantRoleWithId 授权角色给 target，会将原有的角色授权先取消掉
	ReGrantRoleWithId(ctx int64, target string, roleIds ...int64) (err error)

	// RevokeRole 取消对 target 的角色授权
	RevokeRole(ctx int64, target string, roleNames ...string) (err error)

	// RevokeRoleWithId 取消对 target 的角色授权
	RevokeRoleWithId(ctx int64, target string, roleIds ...int64) (err error)

	// RevokeAllRole 取消对 target 的所有角色授权
	RevokeAllRole(ctx int64, target string) (err error)

	// AddRoleMutex 添加角色互斥关系
	AddRoleMutex(ctx int64, roleName string, mutexRoleNames ...string) (err error)

	// AddRoleMutexWithId 添加角色互斥关系
	AddRoleMutexWithId(ctx, roleId int64, mutexRoleIds ...int64) (err error)

	// RemoveRoleMutex 删除角色互斥关系
	RemoveRoleMutex(ctx int64, roleName string, mutexRoleNames ...string) (err error)

	// RemoveRoleMutexWithId 删除角色互斥关系
	RemoveRoleMutexWithId(ctx, roleId int64, mutexRoleIds ...int64) (err error)

	// RemoveAllRoleMutex 删除该角色所有的互斥关系
	RemoveAllRoleMutex(ctx int64, roleName string) (err error)

	// RemoveAllRoleMutexWithId 删除该角色所有的互斥关系
	RemoveAllRoleMutexWithId(ctx, roleId int64) (err error)

	// GetMutexRoles 获取与该角色互斥的角色列表
	GetMutexRoles(ctx int64, roleName string) (result []*RoleMutex, err error)

	// GetMutexRolesWithId 获取与该角色互斥的角色列表
	GetMutexRolesWithId(ctx, roleId int64) (result []*RoleMutex, err error)

	// CheckRoleMutex 验证两个角色是否互斥
	CheckRoleMutex(ctx int64, roleName, mutexRoleName string) bool

	// CheckRoleMutexWithId 验证两个角色是否互斥
	CheckRoleMutexWithId(ctx, roleId, mutexRoleId int64) bool

	// AddPreRole 添加授予该角色时需要的先决条件
	AddPreRole(ctx int64, roleName string, preRoleNames ...string) (err error)

	// AddPreRoleWithId 添加授予该角色时需要的先决条件
	AddPreRoleWithId(ctx, roleId int64, preRoleIds ...int64) (err error)

	// RemovePreRole 删除授予该角色时需要的先决条件
	RemovePreRole(ctx int64, roleName string, preRoleNames ...string) (err error)

	// RemovePreRoleWithId 删除授予该角色时需要的先决条件
	RemovePreRoleWithId(ctx, roleId int64, preRoleIds ...int64) (err error)

	// RemoveAllPreRole 删除授予该角色时需要的所有先决条件
	RemoveAllPreRole(ctx int64, roleName string) (err error)

	// RemoveAllPreRoleWithId 删除授予该角色时需要的所有先决条件
	RemoveAllPreRoleWithId(ctx, roleId int64) (err error)

	// GetPreRoles 获取授予该角色时需要的所有先决条件
	GetPreRoles(ctx int64, roleName string) (result []*PreRole, err error)

	// GetPreRolesWithId 获取授予该角色时需要的所有先决条件
	GetPreRolesWithId(ctx, roleId int64) (result []*PreRole, err error)

	// GetGrantedRoles 获取已授权给 target 的角色列表
	GetGrantedRoles(ctx int64, target string) (result []*Role, err error)

	// GetRolesWithTarget 获取已授权给 target 的角色，及其角色的子角色
	GetRolesWithTarget(ctx int64, target string) (result []*Role, err error)

	// CheckRole 验证 target 是否拥有指定角色
	CheckRole(ctx int64, target string, roleName string) bool

	// CheckRoleWithId 验证 target 是否拥有指定角色
	CheckRoleWithId(ctx int64, target string, roleId int64) bool

	// CheckRoleAccessible 验证 target 是否拥有操作访问 roleName 的权限
	CheckRoleAccessible(ctx int64, target string, roleName string) bool

	// CheckRoleAccessibleWithId 验证 target 是否拥有操作访问 roleId 的权限
	CheckRoleAccessibleWithId(ctx int64, target string, roleId int64) bool

	// GetPermissionsWithRole 获取已授权给 roleName 的权限列表
	GetPermissionsWithRole(ctx int64, roleName string) (result []*Permission, err error)

	// GetPermissionsWithRoleId 获取已授权给 roleId 的权限列表
	GetPermissionsWithRoleId(ctx int64, roleId int64) (result []*Permission, err error)

	// GetGrantedPermissions 获取已授权给 target 的权限列表
	GetGrantedPermissions(ctx int64, target string) (result []*Permission, err error)

	// GetPermissionsTreeWithRole 获取权限组列表，组中包含该组所有的权限信息
	// 如果参数 roleName 的值不为空字符串，则返回的权限数据中将附带该权限是否已授权给该 roleName
	// 如果参数 limitedInParentRole 的值为 true，并且 roleName 的值不为空字符串，则返回的权限数据将限定在 roleName 的父角色拥有的权限范围内
	GetPermissionsTreeWithRole(ctx int64, roleName string, status Status, limitedInParentRole bool) (result []*Group, err error)

	// GetPermissionsTreeWithRoleId 获取权限组列表，组中包含该组所有的权限信息
	// 如果参数 roleId 的值大于 0，则返回的权限数据中将附带该权限是否已授权给该 roleId
	// 如果参数 limitedInParentRole 的值为 true，并且 roleId 的值大于 0，则返回的权限数据将限定在 roleId 的父角色拥有的权限范围内
	GetPermissionsTreeWithRoleId(ctx, roleId int64, status Status, limitedInParentRole bool) (result []*Group, err error)

	// CheckPermission 验证 target 是否拥有指定权限
	CheckPermission(ctx int64, target string, permissionName string) bool

	// CheckPermissionWithId 验证 target 是否拥有指定权限
	CheckPermissionWithId(ctx int64, target string, permissionId int64) bool

	// CheckRolePermission 验证角色是否拥有指定权限
	CheckRolePermission(ctx int64, roleName, permissionName string) bool

	// CheckRolePermissionWithId 验证角色是否拥有指定权限
	CheckRolePermissionWithId(ctx, roleId, permissionId int64) bool

	// CleanCache 清除缓存，如果 target 为空字符串或者 target 的值为星号(*)，则会清空所有缓存
	CleanCache(ctx int64, target string)
}
