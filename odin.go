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
	Id             int64         `json:"id,string"                       sql:"id"`
	Ctx            int64         `json:"ctx,string"                      sql:"ctx"`
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
	Id                int64            `json:"id,string"                        sql:"id"`
	GroupId           int64            `json:"group_id,string"                  sql:"group_id"`
	Ctx               int64            `json:"ctx,string"                       sql:"ctx"`
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
	Id             int64         `json:"id,string"                       sql:"id"`
	Ctx            int64         `json:"ctx,string"                      sql:"ctx"`
	Name           string        `json:"name"                            sql:"name"`
	AliasName      string        `json:"alias_name"                      sql:"alias_name"`
	Status         Status        `json:"status"                          sql:"status"`
	Description    string        `json:"description"                     sql:"description"`
	Granted        bool          `json:"granted"                         sql:"granted"`    // 角色是否授予给指定 target
	Accessible     bool          `json:"accessible"                      sql:"can_access"` // 角色是否能够被指定 target 操作访问
	ParentId       int64         `json:"parent_id,string"                sql:"parent_id"`
	LeftValue      int64         `json:"left_value,string"               sql:"left_value"`
	RightValue     int64         `json:"right_value,string"              sql:"right_value"`
	Depth          int           `json:"depth"                           sql:"depth"`
	CreatedOn      *time.Time    `json:"created_on"                      sql:"created_on"`
	UpdatedOn      *time.Time    `json:"updated_on"                      sql:"updated_on"`
	PermissionList []*Permission `json:"permission_list,omitempty"       sql:"-"`
	MutexRoleList  []*RoleMutex  `json:"mutex_role_list,omitempty"       sql:"-"`
	PreRoleList    []*PreRole    `json:"pre_role_list,omitempty"         sql:"-"`
}

// RolePermission 角色-权限数据结构，用于描述角色与权限的关联关系。
type RolePermission struct {
	Ctx          int64      `json:"ctx,string"             sql:"ctx"`
	RoleId       int64      `json:"role_id,string"         sql:"role_id"`
	PermissionId int64      `json:"permission_id,string"   sql:"permission_id"`
	CreatedOn    *time.Time `json:"created_on"             sql:"created_on"`
}

// RoleMutex 角色互斥数据结构，用于描述角色与角色之间的互斥关系。
// 主要应用于更新（授予）某一 target 的角色时，判断该 target 的角色是否有互斥的情况。
// 比如：角色A与角色B互斥，则角色A和角色B不能同时授予同一 target。
type RoleMutex struct {
	Ctx                int64      `json:"ctx,string"                sql:"ctx"`
	RoleId             int64      `json:"role_id,string"            sql:"role_id"`
	RoleName           string     `json:"role_name"                 sql:"role_name"`
	RoleAliasName      string     `json:"role_alias_name"           sql:"role_alias_name"`
	MutexRoleId        int64      `json:"mutex_role_id,string"      sql:"mutex_role_id"`
	MutexRoleName      string     `json:"mutex_role_name"           sql:"mutex_role_name"`
	MutexRoleAliasName string     `json:"mutex_role_alias_name"     sql:"mutex_role_alias_name"`
	CreatedOn          *time.Time `json:"created_on"                sql:"created_on"`
}

// PreRole 角色先决条件数据结构。
// 主要应用于更新（授予）某一 target 的角色时，判断该 target 是否已经拥有某一角色。
//
// 比如：角色A是角色B的先决条件，向 target 授予角色B时，需要 target 已经拥有角色A。
type PreRole struct {
	Ctx              int64      `json:"ctx,string"                sql:"ctx"`
	RoleId           int64      `json:"role_id,string"            sql:"role_id"`
	RoleName         string     `json:"role_name"                 sql:"role_name"`
	RoleAliasName    string     `json:"role_alias_name"           sql:"role_alias_name"`
	PreRoleId        int64      `json:"pre_role_id,string"        sql:"pre_role_id"`
	PreRoleName      string     `json:"pre_role_name"             sql:"pre_role_name"`
	PreRoleAliasName string     `json:"pre_role_alias_name"       sql:"pre_role_alias_name"`
	CreatedOn        *time.Time `json:"created_on"                sql:"created_on"`
}

// PrePermission 权限先决条件数据结构。
// 主要应用于更新（授予）某一角色的权限时，判断该角色是否已经拥有某一权限。
//
// 比如：权限A是权限B的先决条件，向角色A授予权限B时，需要角色A已经拥有权限A。
type PrePermission struct {
	Ctx                    int64      `json:"ctx,string"                    sql:"ctx"`
	PermissionId           int64      `json:"permission_id,string"          sql:"permission_id"`
	PermissionName         string     `json:"permission_name"               sql:"permission_name"`
	PermissionAliasName    string     `json:"permission_alias_name"         sql:"permission_alias_name"`
	PrePermissionId        int64      `json:"pre_permission_id,string"      sql:"pre_permission_id"`
	PrePermissionName      string     `json:"pre_permission_name"           sql:"pre_permission_name"`
	PrePermissionAliasName string     `json:"pre_permission_alias_name"     sql:"pre_permission_alias_name"`
	AutoGrant              bool       `json:"auto_grant"                    sql:"auto_grant"`
	CreatedOn              *time.Time `json:"created_on"                    sql:"created_on"`
}

// Grant 用于描述 target、角色、权限之间的关系。
type Grant struct {
	Ctx            int64  `json:"ctx,string"                sql:"ctx"`
	Target         string `json:"target"                    sql:"target"`
	RoleId         int64  `json:"role_id,string"            sql:"role_id"`
	RoleName       string `json:"role_name"                 sql:"role_name"`
	PermissionId   int64  `json:"permission_id,string"      sql:"permission_id"`
	PermissionName string `json:"permission_name"           sql:"permission_name"`
}
