package odin

import "github.com/smartwalle/dbs"

type Service struct {
	m *manager
}

func NewService(db dbs.DB) *Service {
	var s = &Service{}
	var m = &manager{}
	m.db = db
	m.groupTable = "odin_group"
	m.permissionTable = "odin_permission"
	m.roleTable = "odin_role"
	m.rolePermissionTable = "odin_role_permission"
	m.roleGrantTable = "odin_grant"
	s.m = m
	return s
}

// --------------------------------------------------------------------------------
func (this *Service) GetGroupList(gType, status int, name string) (result []*Group, err error) {
	return this.m.getGroupList(gType, status, name)
}

func (this *Service) GetGroupWithId(id int64) (result *Group, err error) {
	return this.m.getGroupWithId(id)
}

func (this *Service) GetGroupWithName(name string) (result *Group, err error) {
	return this.m.getGroupWithName(0, name)
}

func (this *Service) GetRoleGroupWithName(name string) (result *Group, err error) {
	return this.m.getGroupWithName(K_GROUP_TYPE_ROLE, name)
}

func (this *Service) GetPermissionGroupWithName(name string) (result *Group, err error) {
	return this.m.getGroupWithName(K_GROUP_TYPE_PERMISSION, name)
}

func (this *Service) AddRoleGroup(name string, status int) (result *Group, err error) {
	return this.m.addGroup(K_GROUP_TYPE_ROLE, name, status)
}

func (this *Service) AddPermissionGroup(name string, status int) (result *Group, err error) {
	return this.m.addGroup(K_GROUP_TYPE_PERMISSION, name, status)
}

func (this *Service) AddGroup(gType int, name string, status int) (result *Group, err error) {
	return this.m.addGroup(gType, name, status)
}

func (this *Service) UpdateGroup(id int64, name string, status int) (err error) {
	return this.m.updateGroup(id, name, status)
}

func (this *Service) UpdateGroupStatus(id int64, name string, status int) (err error) {
	return this.m.updateGroupStatus(id, status)
}

// --------------------------------------------------------------------------------
func (this *Service) GetPermissionList(groupId int64, status int, keyword string) (result []*Permission, err error) {
	return this.m.getPermissionList(groupId, status, keyword)
}

func (this *Service) GetPermissionWithId(id int64) (result *Permission, err error) {
	return this.m.getPermissionWithId(id)
}

func (this *Service) GetPermissionWithName(name string) (result *Permission, err error) {
	return this.m.getPermissionWithName(name)
}

func (this *Service) GetPermissionWithIdentifier(identifier string) (result *Permission, err error) {
	return this.m.getPermissionWithIdentifier(identifier)
}

func (this *Service) AddPermission(groupId int64, name, identifier string, status int) (result *Permission, err error) {
	return this.m.addPermission(groupId, name, identifier, status)
}

func (this *Service) UpdatePermission(id, groupId int64, name, identifier string, status int) (err error) {
	return this.m.updatePermission(id, groupId, name, identifier, status)
}

func (this *Service) UpdatePermissionStatus(id int64, status int) (err error) {
	return this.m.updatePermissionStatus(id, status)
}

// --------------------------------------------------------------------------------
func (this *Service) GetRoleList(groupId int64, status int, keyword string) (result []*Role, err error) {
	return this.m.getRoleList(groupId, status, keyword)
}

func (this *Service) GetPermissionListWithRole(roleId int64) (result []*Permission, err error) {
	return this.m.getPermissionListWithRoleId(roleId)
}

func (this *Service) GetRoleWithId(id int64) (result *Role, err error) {
	return this.m.getRoleWithId(id, true)
}

func (this *Service) GetRoleWithName(name string) (result *Role, err error) {
	return this.m.getRoleWithName(name, true)
}

func (this *Service) CheckRoleNameIsExists(name string) (result bool) {
	role, err := this.m.getRoleWithName(name, false)
	if role != nil || err != nil {
		return true
	}
	return false
}

func (this *Service) AddRole(groupId int64, name string, status int) (result *Role, err error) {
	return this.m.addRole(groupId, name, status)
}

func (this *Service) UpdateRole(id, groupId int64, name string, status int) (err error) {
	return this.m.updateRole(id, groupId, name, status)
}

func (this *Service) UpdateRoleStatus(id int64, status int) (err error) {
	return this.m.updateRoleStatus(id, status)
}

// --------------------------------------------------------------------------------
