package odin

import (
	"github.com/smartwalle/dbr"
	"github.com/smartwalle/dbs"
	"strings"
)

type Service struct {
	m *manager
	r *redisManager
}

func NewService(db dbs.DB, redis *dbr.Pool, tablePrefix string) *Service {
	var s = &Service{}
	var m = &manager{}

	tablePrefix = strings.TrimSpace(tablePrefix)
	if tablePrefix == "" {
		tablePrefix = "odin"
	}

	m.db = db
	m.groupTable = tablePrefix + "_group"
	m.permissionTable = tablePrefix + "_permission"
	m.roleTable = tablePrefix + "_role"
	m.rolePermissionTable = tablePrefix + "_role_permission"
	m.roleGrantTable = tablePrefix + "_grant"
	s.m = m

	if redis != nil {
		var r = &redisManager{}
		r.r = redis
		r.tPrefix = tablePrefix
		s.r = r
	}
	return s
}

// GetPermissionTree 获取权限组列表，会返回该组包含的权限列表
// 如果 roleId 大于 0，则会返回各权限是否有授权给该角色
func (this *Service) GetPermissionTree(ctxId, roleId int64, status int, name string) (result []*Group, err error) {
	return this.m.getPermissionTree(ctxId, roleId, status, name)
}

// GetRoleTree 获取角色组列表，会返回该组包含的角色列表
// 如果 objectId 不为空字符串，则会返回各角色是否有授权给该对象
func (this *Service) GetRoleTree(ctxId int64, objectId string, status int, name string) (result []*Group, err error) {
	return this.m.getRoleTree(ctxId, objectId, status, name)
}

// --------------------------------------------------------------------------------
// GetPermissionGroupList 获取权限组列表，组信息不包含权限列表
func (this *Service) GetPermissionGroupList(ctxId int64, status int, name string) (result []*Group, err error) {
	return this.m.getGroupListWithType(ctxId, K_GROUP_TYPE_PERMISSION, status, name)
}

// GetRoleGroupList 获取角色组列表，组信息不包含角色列表
func (this *Service) GetRoleGroupList(ctxId int64, status int, name string) (result []*Group, err error) {
	return this.m.getGroupListWithType(ctxId, K_GROUP_TYPE_ROLE, status, name)
}

// GetPermissionGroupWithId 获取权限组详情，包含权限列表或者角色列表
func (this *Service) GetPermissionGroupWithId(ctxId int64, id int64) (result *Group, err error) {
	return this.m.getGroupWithId(ctxId, id, K_GROUP_TYPE_PERMISSION)
}

// GetRoleGroupWithId 获取角色组详情，包含权限列表或者角色列表
func (this *Service) GetRoleGroupWithId(ctxId, id int64) (result *Group, err error) {
	return this.m.getGroupWithId(ctxId, id, K_GROUP_TYPE_ROLE)
}

// GetPermissionGroupWithName 根据组名称查询权限组信息（精确匹配），返回数据不包含该组的权限列表
func (this *Service) GetPermissionGroupWithName(ctxId int64, name string) (result *Group, err error) {
	return this.m.getGroupWithName(ctxId, name, K_GROUP_TYPE_PERMISSION)
}

// GetRoleGroupWithName 根据组名称查询角色组信息（精确匹配），返回数据不包含该组的角色列表
func (this *Service) GetRoleGroupWithName(ctxId int64, name string) (result *Group, err error) {
	return this.m.getGroupWithName(ctxId, name, K_GROUP_TYPE_ROLE)
}

// AddPermissionGroup 添加权限组
func (this *Service) AddPermissionGroup(ctxId int64, name string, status int) (result *Group, err error) {
	if result, err = this.m.getGroupWithName(ctxId, name, K_GROUP_TYPE_PERMISSION); err != nil {
		return nil, err
	}
	if result != nil {
		return nil, ErrGroupExists
	}
	return this.m.addGroup(ctxId, K_GROUP_TYPE_PERMISSION, name, status)
}

// AddRoleGroup 添加角色组
func (this *Service) AddRoleGroup(ctxId int64, name string, status int) (result *Group, err error) {
	if result, err = this.m.getGroupWithName(ctxId, name, K_GROUP_TYPE_ROLE); err != nil {
		return nil, err
	}
	if result != nil {
		return nil, ErrGroupExists
	}
	return this.m.addGroup(ctxId, K_GROUP_TYPE_ROLE, name, status)
}

// UpdatePermissionGroup 更新权限组的基本信息
func (this *Service) UpdatePermissionGroup(ctxId int64, id int64, name string, status int) (err error) {
	result, err := this.m.getGroupWithName(ctxId, name, K_GROUP_TYPE_PERMISSION)
	if err != nil {
		return err
	}
	if result != nil && result.Id != id {
		return ErrGroupExists
	}
	return this.m.updateGroup(ctxId, id, name, status)
}

// UpdateRoleGroup 更新权限组的基本信息
func (this *Service) UpdateRoleGroup(ctxId int64, id int64, name string, status int) (err error) {
	result, err := this.m.getGroupWithName(ctxId, name, K_GROUP_TYPE_ROLE)
	if err != nil {
		return err
	}
	if result != nil && result.Id != id {
		return ErrGroupExists
	}
	return this.m.updateGroup(ctxId, id, name, status)
}

// UpdateGroupStatus 更新组的状态信息
func (this *Service) UpdateGroupStatus(ctxId, id int64, status int) (err error) {
	return this.m.updateGroupStatus(ctxId, id, status)
}

// RemoveGroup 删除组信息
func (this *Service) RemoveGroup(ctxId, id int64) (err error) {
	group, err := this.m.getGroupWithId(ctxId, id, 0)
	if err != nil {
		return err
	}
	if group == nil {
		return nil
	}

	if group.CtxId != ctxId {
		return ErrRemoveGroupNotAllowed
	}

	// 如果 group 下还有内容，则不能删除
	if group.Type == K_GROUP_TYPE_PERMISSION {
		pList, err := this.m.getPermissionList(ctxId, []int64{id}, 0, "")
		if err != nil {
			return err
		}
		if len(pList) > 0 {
			return ErrRemoveGroupNotAllowed
		}
	} else if group.Type == K_GROUP_TYPE_ROLE {
		rList, err := this.m.getRoleList(ctxId, id, 0, "")
		if err != nil {
			return err
		}
		if len(rList) > 0 {
			return ErrRemoveGroupNotAllowed
		}
	}
	return this.m.removeGroup(ctxId, id)
}

// --------------------------------------------------------------------------------
// GetPermissionList 获取指定组的权限列表
func (this *Service) GetPermissionList(ctxId, groupId int64, status int, keyword string) (result []*Permission, err error) {
	var groupIdList []int64
	if groupId > 0 {
		groupIdList = append(groupIdList, groupId)
	}
	return this.m.getPermissionList(ctxId, groupIdList, status, keyword)
}

// GetPermissionWithId 获取权限详情
func (this *Service) GetPermissionWithId(ctxId, id int64) (result *Permission, err error) {
	return this.m.getPermissionWithId(ctxId, id)
}

// GetPermissionWithName 根据权限名称获取权限信息（精确匹配）
func (this *Service) GetPermissionWithName(ctxId int64, name string) (result *Permission, err error) {
	return this.m.getPermissionWithName(ctxId, name)
}

// GetPermissionWithIdentifier 权限权限标识符获取权限信息（精确匹配）
func (this *Service) GetPermissionWithIdentifier(ctxId int64, identifier string) (result *Permission, err error) {
	return this.m.getPermissionWithIdentifier(ctxId, identifier)
}

// AddPermission 添加权限
func (this *Service) AddPermission(ctxId, groupId int64, name, identifier string, status int) (result *Permission, err error) {
	if this.CheckPermissionIsExists(ctxId, identifier) == true {
		return nil, ErrPermissionIdentifierExists
	}
	if this.CheckPermissionNameIsExists(ctxId, name) == true {
		return nil, ErrPermissionNameExists
	}

	group, err := this.m.getGroupWithId(ctxId, groupId, K_GROUP_TYPE_PERMISSION)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, ErrGroupNotExist
	}
	return this.m.addPermission(ctxId, groupId, name, identifier, status)
}

// UpdatePermission 更新权限信息
func (this *Service) UpdatePermission(ctxId, id, groupId int64, name, identifier string, status int) (err error) {
	p, err := this.m.getPermissionWithIdentifier(ctxId, identifier)
	if err != nil {
		return err
	}
	if p != nil && p.Id != id {
		return ErrPermissionIdentifierExists
	}

	p, err = this.m.getPermissionWithName(ctxId, name)
	if err != nil {
		return err
	}
	if p != nil && p.Id != id {
		return ErrPermissionNameExists
	}

	group, err := this.m.getGroupWithId(ctxId, groupId, K_GROUP_TYPE_PERMISSION)
	if err != nil {
		return err
	}
	if group == nil {
		return ErrGroupNotExist
	}
	return this.m.updatePermission(ctxId, id, groupId, name, identifier, status)
}

// CheckPermissionIsExists 验证权限标识已经是否已经存在
func (this *Service) CheckPermissionIsExists(ctxId int64, identifier string) (result bool) {
	p, err := this.m.getPermissionWithIdentifier(ctxId, identifier)
	if p != nil || err != nil {
		return true
	}
	return false
}

// CheckPermissionNameIsExists 验证权限名称是否已经存在
func (this *Service) CheckPermissionNameIsExists(ctxId int64, name string) (result bool) {
	p, err := this.m.getPermissionWithName(ctxId, name)
	if p != nil || err != nil {
		return true
	}
	return false
}

// UpdatePermissionStatus 更新权限的状态信息
func (this *Service) UpdatePermissionStatus(ctxId, id int64, status int) (err error) {
	return this.m.updatePermissionStatus(ctxId, id, status)
}

// --------------------------------------------------------------------------------
// GetRoleList 获取指定组的角色组列表
func (this *Service) GetRoleList(ctxId, groupId int64, status int, keyword string) (result []*Role, err error) {
	return this.m.getRoleList(ctxId, groupId, status, keyword)
}

// GetPermissionListWithRole 获取指定角色的权限列表
func (this *Service) GetPermissionListWithRole(ctxId, roleId int64) (result []*Permission, err error) {
	return this.m.getPermissionListWithRoleId(ctxId, roleId)
}

// GetRoleWithId 获取角色详情，会返回该角色拥有的权限列表
func (this *Service) GetRoleWithId(ctxId, id int64) (result *Role, err error) {
	return this.m.getRoleWithId(ctxId, id, true)
}

// GetRoleWithName 根据角色名称获取角色信息（精确匹配），会返回该角色拥有的权限列表
func (this *Service) GetRoleWithName(ctxId int64, name string) (result *Role, err error) {
	return this.m.getRoleWithName(ctxId, name, true)
}

// CheckRoleNameIsExists 检测角色名是否已经存在
func (this *Service) CheckRoleNameIsExists(ctxId int64, name string) (result bool) {
	role, err := this.m.getRoleWithName(ctxId, name, false)
	if role != nil || err != nil {
		return true
	}
	return false
}

// AddRole 添加角色
func (this *Service) AddRole(ctxId, groupId int64, name string, status int) (result *Role, err error) {
	if this.CheckRoleNameIsExists(ctxId, name) == true {
		return nil, ErrRoleNameExists
	}

	group, err := this.m.getGroupWithId(ctxId, groupId, K_GROUP_TYPE_ROLE)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, ErrGroupNotExist
	}
	return this.m.addRole(ctxId, groupId, name, status)
}

// UpdateRole 更新角色信息
func (this *Service) UpdateRole(ctxId, id, groupId int64, name string, status int) (err error) {
	role, err := this.m.getRoleWithName(ctxId, name, false)
	if err != nil {
		return err
	}
	if role != nil && role.Id != id {
		return ErrRoleNameExists
	}

	group, err := this.m.getGroupWithId(ctxId, groupId, K_GROUP_TYPE_ROLE)
	if err != nil {
		return err
	}
	if group == nil {
		return ErrGroupNotExist
	}
	return this.m.updateRole(ctxId, id, groupId, name, status)
}

// UpdateRoleStatus 更新角色状态信息
func (this *Service) UpdateRoleStatus(ctxId, id int64, status int) (err error) {
	return this.m.updateRoleStatus(ctxId, id, status)
}

// --------------------------------------------------------------------------------
// GrantPermission 为角色添加权限信息
func (this *Service) GrantPermission(ctxId, roleId int64, permissionIdList ...int64) (err error) {
	role, err := this.m.getRoleWithId(ctxId, roleId, false)
	if err != nil {
		return err
	}
	if role == nil {
		return ErrRoleNotExist
	}
	if role.Status != K_STATUS_ENABLE {
		return ErrRoleNotExist
	}

	pList, err := this.m.getPermissionWithIdList(ctxId, permissionIdList)
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
	return this.m.grantPermission(ctxId, roleId, nIdList)
}

func (this *Service) RevokePermission(ctxId, roleId int64, permissionIdList ...int64) (err error) {
	return this.m.revokePermission(ctxId, roleId, permissionIdList)
}

func (this *Service) GrantRole(ctxId int64, objectId string, roleIdList ...int64) (err error) {
	if len(roleIdList) == 0 {
		return ErrRoleNotExist
	}
	if objectId == "" {
		return ErrObjectNotAllowed
	}
	roleList, err := this.m.getRoleWithIdList(ctxId, roleIdList)
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

	err = this.m.grantRole(ctxId, objectId, nIdList)
	return err
}

func (this *Service) RevokeRole(ctxId int64, objectId string, roleIdList ...int64) (err error) {
	return this.m.revokeRole(ctxId, objectId, roleIdList)
}

func (this *Service) Check(ctxId int64, objectId, identifier string) (result bool) {
	if this.r != nil {
		result = this.r.check(ctxId, objectId, identifier)
		if result == false {
			if this.r.exists(ctxId, objectId) == false {
				pList, _ := this.m.getGrantedPermissionList(ctxId, objectId)
				var identifierList []interface{}
				for _, p := range pList {
					if p.Identifier == identifier {
						result = true
					}
					identifierList = append(identifierList, p.Identifier)
				}
				this.r.grantPermissions(ctxId, objectId, identifierList)
			}
		}
		return result
	}
	return this.m.check(ctxId, objectId, identifier)
}

func (this *Service) CheckList(ctxId int64, objectId string, identifiers ...string) (result map[string]bool) {
	if this.r != nil {
		return this.r.checkList(ctxId, objectId, identifiers...)
	}
	return this.m.checkList(ctxId, objectId, identifiers...)
}

func (this *Service) GetGrantedRoleList(ctxId int64, objectId string) (result []*Role, err error) {
	return this.m.getGrantedRoleList(ctxId, objectId)
}

func (this *Service) GetGrantedPermissionList(ctxId int64, objectId string) (result []*Permission, err error) {
	return this.m.getGrantedPermissionList(ctxId, objectId)
}

func (this *Service) ClearCache() {
	if this.r != nil {
		this.r.clear()
	}
}
