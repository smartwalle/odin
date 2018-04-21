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
	s.m = m
	return s
}

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
