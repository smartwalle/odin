package odin

import (
	"time"
)

type GroupType int

const (
	GroupTypeOfPermission GroupType = 1
	GroupTypeOfRole       GroupType = 2
)

type Status int

const (
	StatusOfEnable  Status = 1
	StatusOfDisable Status = 2
)

type Group struct {
	Id             int64         `json:"id"                          sql:"id"`
	Ctx            int64         `json:"ctx"                         sql:"ctx"`
	Type           GroupType     `json:"type"                        sql:"type"`
	Name           string        `json:"name"                        sql:"name"`
	Status         Status        `json:"status"                      sql:"status"`
	CreatedOn      *time.Time    `json:"created_on"                  sql:"created_on"`
	UpdatedOn      *time.Time    `json:"updated_on"                  sql:"updated_on"`
	PermissionList []*Permission `json:"permission_list,omitempty"   sql:"-"`
	RoleList       []*Role       `json:"role_list,omitempty"         sql:"-"`
}

type Permission struct {
	Id         int64      `json:"id"            sql:"id"`
	Ctx        int64      `json:"ctx"           sql:"ctx"`
	GroupId    int64      `json:"group_id"      sql:"group_id"`
	Name       string     `json:"name"          sql:"name"`
	Identifier string     `json:"identifier"    sql:"identifier"`
	Status     Status     `json:"status"        sql:"status"`
	Granted    bool       `json:"granted"       sql:"granted"`
	CreatedOn  *time.Time `json:"created_on"    sql:"created_on"`
	UpdatedOn  *time.Time `json:"updated_on"    sql:"updated_on"`
}

type Role struct {
	Id             int64         `json:"id"                          sql:"id"`
	Ctx            int64         `json:"ctx"                         sql:"ctx"`
	GroupId        int64         `json:"group_id"                    sql:"group_id"`
	Name           string        `json:"name"                        sql:"name"`
	Status         Status        `json:"status"                      sql:"status"`
	CreatedOn      *time.Time    `json:"created_on"                  sql:"created_on"`
	UpdatedOn      *time.Time    `json:"updated_on"                  sql:"updated_on"`
	Granted        bool          `json:"granted"                     sql:"granted"`
	PermissionList []*Permission `json:"permission_list,omitempty"   sql:"-"`
}

type RolePermission struct {
	Ctx          int64      `json:"ctx"              sql:"ctx"`
	RoleId       int64      `json:"role_id"          sql:"role_id"`
	PermissionId int64      `json:"permission_id"    sql:"permission_id"`
	CreatedOn    *time.Time `json:"created_on"       sql:"created_on"`
}

type Grant struct {
	Ctx          int64  `json:"ctx"           sql:"ctx"`
	Target       string `json:"target"        sql:"target"`
	RoleId       int64  `json:"role_id"       sql:"role_id"`
	PermissionId int64  `json:"permission_id" sql:"permission_id"`
	Identifier   string `json:"identifier"    sql:"identifier"`
}

type Service interface {
	GetPermissionTree(ctx, roleId int64, status Status, name string) (result []*Group, err error)

	GetRoleTree(ctx int64, target string, status Status, name string) (result []*Group, err error)

	GetPermissionGroupList(ctx int64, status Status, name string) (result []*Group, err error)

	GetRoleGroupList(ctx int64, status Status, name string) (result []*Group, err error)

	GetPermissionGroupWithId(ctx int64, id int64) (result *Group, err error)

	GetRoleGroupWithId(ctx, id int64) (result *Group, err error)

	GetPermissionGroupWithName(ctx int64, name string) (result *Group, err error)

	GetRoleGroupWithName(ctx int64, name string) (result *Group, err error)

	AddPermissionGroup(ctx int64, name string, status Status) (result *Group, err error)

	AddRoleGroup(ctx int64, name string, status Status) (result *Group, err error)

	UpdatePermissionGroup(ctx int64, id int64, name string, status Status) (err error)

	UpdateRoleGroup(ctx int64, id int64, name string, status Status) (err error)

	UpdateGroupStatus(ctx, id int64, status Status) (err error)

	RemoveGroup(ctx, id int64) (err error)

	GetPermissionList(ctx, groupId int64, status Status, keyword string) (result []*Permission, err error)

	GetPermissionWithId(ctx, id int64) (result *Permission, err error)

	GetPermissionWithName(ctx int64, name string) (result *Permission, err error)

	GetPermissionWithIdentifier(ctx int64, identifier string) (result *Permission, err error)

	AddPermission(ctx, groupId int64, name, identifier string, status Status) (result *Permission, err error)

	UpdatePermission(ctx, id, groupId int64, name, identifier string, status Status) (err error)

	CheckPermissionIsExists(ctx int64, identifier string) (result bool)

	CheckPermissionNameIsExists(ctx int64, name string) (result bool)

	UpdatePermissionStatus(ctx, id int64, status Status) (err error)

	GetRoleList(ctx, groupId int64, status Status, keyword string) (result []*Role, err error)

	GetPermissionListWithRole(ctx, roleId int64) (result []*Permission, err error)

	GetRoleWithId(ctx, id int64) (result *Role, err error)

	GetRoleWithName(ctx int64, name string) (result *Role, err error)

	CheckRoleNameIsExists(ctx int64, name string) (result bool)

	AddRole(ctx, groupId int64, name string, status Status) (result *Role, err error)

	UpdateRole(ctx, id, groupId int64, name string, status Status) (err error)

	UpdateRoleStatus(ctx, id int64, status Status) (err error)

	GrantPermission(ctx, roleId int64, permissionIdList ...int64) (err error)

	RevokePermission(ctx, roleId int64, permissionIdList ...int64) (err error)

	ReGrantPermission(ctx, roleId int64, permissionIdList ...int64) (err error)

	GrantRole(ctx int64, target string, roleIdList ...int64) (err error)

	RevokeRole(ctx int64, target string, roleIdList ...int64) (err error)

	ReGrantRole(ctx int64, target string, roleIdList ...int64) (err error)

	Check(ctx int64, target, identifier string) (result bool)

	CheckList(ctx int64, target string, identifiers ...string) (result map[string]bool)

	GetGrantedRoleList(ctx int64, target string) (result []*Role, err error)

	GetGrantedPermissionList(ctx int64, target string) (result []*Permission, err error)

	ClearCache(ctx int64, target string)
}
