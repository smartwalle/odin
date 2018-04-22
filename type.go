package odin

import "time"

const (
	K_STATUS_ENABLE  = 1000 // 启用
	K_STATUS_DISABLE = 2000 // 禁用
)

const (
	K_GROUP_TYPE_ROLE       = 1000 // role
	K_GROUP_TYPE_PERMISSION = 2000 // permission
)

type Group struct {
	Id             int64         `json:"id"                          sql:"id"`
	Type           int           `json:"type"                        sql:"type"`
	Name           string        `json:"name"                        sql:"name"`
	Status         int           `json:"status"                      sql:"status"`
	CreatedOn      *time.Time    `json:"created_on"                  sql:"created_on"`
	PermissionList []*Permission `json:"permission_list,omitempty"   sql:"-"`
	RoleList       []*Role       `json:"role_list,omitempty"         sql:"-"`
}

type Permission struct {
	Id         int64      `json:"id"            sql:"id"`
	GroupId    int64      `json:"group_id"      sql:"group_id"`
	Name       string     `json:"name"          sql:"name"`
	Identifier string     `json:"identifier"    sql:"identifier"`
	Status     int        `json:"status"        sql:"status"`
	CreatedOn  *time.Time `json:"created_on"    sql:"created_on"`
}

type Role struct {
	Id             int64         `json:"id"                          sql:"id"`
	GroupId        int64         `json:"group_id"                    sql:"group_id"`
	Name           string        `json:"name"                        sql:"name"`
	Status         int           `json:"status"                      sql:"status"`
	CreatedOn      *time.Time    `json:"created_on"                  sql:"created_on"`
	PermissionList []*Permission `json:"permission_list,omitempty"   sql:"-"`
}

type RolePermission struct {
	RoleId       int64      `json:"role_id"          sql:"role_id"`
	PermissionId int64      `json:"permission_id"    sql:"permission_id"`
	CreatedOn    *time.Time `json:"created_on"       sql:"created_on"`
}

type Grant struct {
	ObjectId  string     `json:"object_id"     sql:"object_id"`
	RoleId    int64      `json:"role_id"       sql:"role_id"`
	CreatedOn *time.Time `json:"created_on"    sql:"created_on"`
}
