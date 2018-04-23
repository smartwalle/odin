package odin

import (
	"github.com/smartwalle/dbr"
	"github.com/smartwalle/dbs"
)

type Service struct {
	m *manager
	r *redisManager
}

func NewService(db dbs.DB, redis *dbr.Pool) *Service {
	var s = &Service{}
	var m = &manager{}
	m.db = db
	m.groupTable = "odin_group"
	m.permissionTable = "odin_permission"
	m.roleTable = "odin_role"
	m.rolePermissionTable = "odin_role_permission"
	m.roleGrantTable = "odin_grant"
	s.m = m

	if redis != nil {
		var r = &redisManager{}
		r.r = redis
		s.r = r
	}
	return s
}

// GetPermissionTree 获取权限组列表，会返回该组包含的权限列表
// 如果 roleId 大于 0，则会返回各权限是否有授权给该角色
func (this *Service) GetPermissionTree(roleId int64, status int, name string) (result []*Group, err error) {
	return this.m.getPermissionTree(roleId, status, name)
}

// GetRoleTree 获取角色组列表，会返回该组包含的角色列表
// 如果 objectId 不为空字符串，则会返回各角色是否有授权给该对象
func (this *Service) GetRoleTree(objectId string, status int, name string) (result []*Group, err error) {
	return this.m.getRoleTree(objectId, status, name)
}

// --------------------------------------------------------------------------------
// GetPermissionGroupList 获取权限组列表，组信息不包含权限列表
func (this *Service) GetPermissionGroupList(status int, name string) (result []*Group, err error) {
	return this.m.getGroupListWithType(K_GROUP_TYPE_PERMISSION, status, name)
}

// GetRoleGroupList 获取角色组列表，组信息不包含角色列表
func (this *Service) GetRoleGroupList(status int, name string) (result []*Group, err error) {
	return this.m.getGroupListWithType(K_GROUP_TYPE_ROLE, status, name)
}

// GetPermissionGroupWithId 获取权限组详情，包含权限列表或者角色列表
func (this *Service) GetPermissionGroupWithId(id int64) (result *Group, err error) {
	return this.m.getGroupWithId(id, K_GROUP_TYPE_PERMISSION)
}

// GetRoleGroupWithId 获取角色组详情，包含权限列表或者角色列表
func (this *Service) GetRoleGroupWithId(id int64) (result *Group, err error) {
	return this.m.getGroupWithId(id, K_GROUP_TYPE_ROLE)
}

// GetPermissionGroupWithName 根据组名称查询权限组信息（精确匹配），返回数据不包含该组的权限列表
func (this *Service) GetPermissionGroupWithName(name string) (result *Group, err error) {
	return this.m.getGroupWithName(name, K_GROUP_TYPE_PERMISSION)
}

// GetRoleGroupWithName 根据组名称查询角色组信息（精确匹配），返回数据不包含该组的角色列表
func (this *Service) GetRoleGroupWithName(name string) (result *Group, err error) {
	return this.m.getGroupWithName(name, K_GROUP_TYPE_ROLE)
}

// AddPermissionGroup 添加权限组
func (this *Service) AddPermissionGroup(name string, status int) (result *Group, err error) {
	return this.m.addGroup(K_GROUP_TYPE_PERMISSION, name, status)
}

// AddRoleGroup 添加角色组
func (this *Service) AddRoleGroup(name string, status int) (result *Group, err error) {
	return this.m.addGroup(K_GROUP_TYPE_ROLE, name, status)
}

// UpdateGroup 更新组的基本信息
func (this *Service) UpdateGroup(id int64, name string, status int) (err error) {
	return this.m.updateGroup(id, name, status)
}

// UpdateGroupStatus 更新组的状态信息
func (this *Service) UpdateGroupStatus(id int64, status int) (err error) {
	return this.m.updateGroupStatus(id, status)
}

// RemoveGroup 删除组信息
func (this *Service) RemoveGroup(id int64) (err error) {
	group, err := this.m.getGroupWithId(id, 0)
	if err != nil {
		return err
	}
	if group == nil {
		return nil
	}

	// 如果 group 下还有内容，则不能删除
	if group.Type == K_GROUP_TYPE_PERMISSION {
		pList, err := this.m.getPermissionList([]int64{id}, 0, "")
		if err != nil {
			return err
		}
		if len(pList) > 0 {
			return ErrRemoveGroupNotAllowed
		}
	} else if group.Type == K_GROUP_TYPE_ROLE {
		rList, err := this.m.getRoleList(id, 0, "")
		if err != nil {
			return err
		}
		if len(rList) > 0 {
			return ErrRemoveGroupNotAllowed
		}
	}
	return this.m.removeGroup(id)
}

// --------------------------------------------------------------------------------
// GetPermissionList 获取指定组的权限列表
func (this *Service) GetPermissionList(groupId int64, status int, keyword string) (result []*Permission, err error) {
	var groupIdList []int64
	if groupId > 0 {
		groupIdList = append(groupIdList, groupId)
	}
	return this.m.getPermissionList(groupIdList, status, keyword)
}

// GetPermissionWithId 获取权限详情
func (this *Service) GetPermissionWithId(id int64) (result *Permission, err error) {
	return this.m.getPermissionWithId(id)
}

// GetPermissionWithName 根据权限名称获取权限信息（精确匹配）
func (this *Service) GetPermissionWithName(name string) (result *Permission, err error) {
	return this.m.getPermissionWithName(name)
}

// GetPermissionWithIdentifier 权限权限标识符获取权限信息（精确匹配）
func (this *Service) GetPermissionWithIdentifier(identifier string) (result *Permission, err error) {
	return this.m.getPermissionWithIdentifier(identifier)
}

// AddPermission 添加权限
func (this *Service) AddPermission(groupId int64, name, identifier string, status int) (result *Permission, err error) {
	if this.CheckPermissionIsExists(identifier) == true {
		return nil, ErrPermissionIdentifierExists
	}
	group, err := this.m.getGroupWithId(groupId, K_GROUP_TYPE_PERMISSION)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, ErrGroupNotExists
	}
	return this.m.addPermission(groupId, name, identifier, status)
}

// UpdatePermission 更新权限信息
func (this *Service) UpdatePermission(id, groupId int64, name, identifier string, status int) (err error) {
	p, err := this.m.getPermissionWithIdentifier(identifier)
	if err != nil {
		return err
	}
	if p != nil && p.Id != id {
		return ErrPermissionIdentifierExists
	}
	group, err := this.m.getGroupWithId(groupId, K_GROUP_TYPE_PERMISSION)
	if err != nil {
		return err
	}
	if group == nil {
		return ErrGroupNotExists
	}
	return this.m.updatePermission(id, groupId, name, identifier, status)
}

func (this *Service) CheckPermissionIsExists(identifier string) (result bool) {
	p, err := this.m.getPermissionWithIdentifier(identifier)
	if p != nil || err != nil {
		return true
	}
	return false
}

// UpdatePermissionStatus 更新权限的状态信息
func (this *Service) UpdatePermissionStatus(id int64, status int) (err error) {
	return this.m.updatePermissionStatus(id, status)
}

// --------------------------------------------------------------------------------
// GetRoleList 获取指定组的角色组列表
func (this *Service) GetRoleList(groupId int64, status int, keyword string) (result []*Role, err error) {
	return this.m.getRoleList(groupId, status, keyword)
}

// GetPermissionListWithRole 获取指定角色的权限列表
func (this *Service) GetPermissionListWithRole(roleId int64) (result []*Permission, err error) {
	return this.m.getPermissionListWithRoleId(roleId)
}

// GetRoleWithId 获取角色详情，会返回该角色拥有的权限列表
func (this *Service) GetRoleWithId(id int64) (result *Role, err error) {
	return this.m.getRoleWithId(id, true)
}

// GetRoleWithName 根据角色名称获取角色信息（精确匹配），会返回该角色拥有的权限列表
func (this *Service) GetRoleWithName(name string) (result *Role, err error) {
	return this.m.getRoleWithName(name, true)
}

// CheckRoleNameIsExists 检测角色名是否已经存在
func (this *Service) CheckRoleNameIsExists(name string) (result bool) {
	role, err := this.m.getRoleWithName(name, false)
	if role != nil || err != nil {
		return true
	}
	return false
}

// AddRole 添加角色
func (this *Service) AddRole(groupId int64, name string, status int) (result *Role, err error) {
	group, err := this.m.getGroupWithId(groupId, K_GROUP_TYPE_ROLE)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, ErrGroupNotExists
	}
	return this.m.addRole(groupId, name, status)
}

// UpdateRole 更新角色信息
func (this *Service) UpdateRole(id, groupId int64, name string, status int) (err error) {
	group, err := this.m.getGroupWithId(groupId, K_GROUP_TYPE_ROLE)
	if err != nil {
		return err
	}
	if group == nil {
		return ErrGroupNotExists
	}
	return this.m.updateRole(id, groupId, name, status)
}

// UpdateRoleStatus 更新角色状态信息
func (this *Service) UpdateRoleStatus(id int64, status int) (err error) {
	return this.m.updateRoleStatus(id, status)
}

// --------------------------------------------------------------------------------
// GrantPermission 为角色添加权限信息
func (this *Service) GrantPermission(roleId int64, permissionIdList ...int64) (err error) {
	role, err := this.m.getRoleWithId(roleId, false)
	if err != nil {
		return err
	}
	if role == nil {
		return ErrRoleNotExists
	}
	if role.Status != K_STATUS_ENABLE {
		return ErrRoleNotExists
	}

	pList, err := this.m.getPermissionWithIdList(permissionIdList)
	if err != nil {
		return err
	}
	var nIdList []int64
	for _, p := range pList {
		if p.Status == K_STATUS_ENABLE {
			nIdList = append(nIdList, p.Id)
		}
	}
	if len(nIdList) == 0 {
		return ErrGrantFailed
	}
	return this.m.grantPermission(roleId, nIdList)
}

func (this *Service) RevokePermission(roleId int64, permissionIdList ...int64) (err error) {
	return this.m.revokePermission(roleId, permissionIdList)
}

func (this *Service) GrantRole(objectId string, roleIdList ...int64) (err error) {
	if len(roleIdList) == 0 {
		return ErrRoleNotExists
	}
	if objectId == "" {
		return ErrObjectNotAllowed
	}
	roleList, err := this.m.getRoleWithIdList(roleIdList)
	if err != nil {
		return err
	}

	var nIdList []int64
	for _, role := range roleList {
		if role.Status == K_STATUS_ENABLE {
			nIdList = append(nIdList, role.Id)
		}
	}
	if len(nIdList) == 0 {
		return ErrGrantFailed
	}

	err = this.m.grantRole(objectId, nIdList)
	return err
}

func (this *Service) RevokeRole(roleId string, roleIdList ...int64) (err error) {
	return this.m.revokeRole(roleId, roleIdList)
}

func (this *Service) Check(objectId, identifier string) (result bool) {
	if this.r != nil {
		result = this.r.check(objectId, identifier)
		if result == false {
			if this.r.exists(objectId) == false {
				pList, _ := this.m.getGrantedPermissionList(objectId)
				var identifierList []interface{}
				for _, p := range pList {
					if p.Identifier == identifier {
						result = true
					}
					identifierList = append(identifierList, p.Identifier)
				}
				this.r.grantPermissions(objectId, identifierList)
			}
		}
		return result
	}
	return this.m.check(objectId, identifier)
}

func (this *Service) GetGrantedRoleList(objectId string) (result []*Role, err error) {
	return this.m.getGrantedRoleList(objectId)
}

func (this *Service) GetGrantedPermissionList(objectId string) (result []*Permission, err error) {
	return this.m.getGrantedPermissionList(objectId)
}

func (this *Service) ClearCache() {
	if this.r != nil {
		this.r.clear()
	}
}
