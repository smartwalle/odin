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
	return this.m.getGroupWithName(name)
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