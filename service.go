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
func (this *Service) GetPermissionTree(ctx, roleId int64, status int, name string) (result []*Group, err error) {
	return this.m.getPermissionTree(ctx, roleId, status, name)
}

// GetRoleTree 获取角色组列表，会返回该组包含的角色列表
// 如果 objectId 不为空字符串，则会返回各角色是否有授权给该对象
func (this *Service) GetRoleTree(ctx int64, objectId string, status int, name string) (result []*Group, err error) {
	return this.m.getRoleTree(ctx, objectId, status, name)
}

// --------------------------------------------------------------------------------
// GetPermissionGroupList 获取权限组列表，组信息不包含权限列表
func (this *Service) GetPermissionGroupList(ctx int64, status int, name string) (result []*Group, err error) {
	return this.m.getGroupListWithType(ctx, K_GROUP_TYPE_PERMISSION, status, name)
}

// GetRoleGroupList 获取角色组列表，组信息不包含角色列表
func (this *Service) GetRoleGroupList(ctx int64, status int, name string) (result []*Group, err error) {
	return this.m.getGroupListWithType(ctx, K_GROUP_TYPE_ROLE, status, name)
}

// GetPermissionGroupWithId 获取权限组详情，包含权限列表或者角色列表
func (this *Service) GetPermissionGroupWithId(ctx int64, id int64) (result *Group, err error) {
	return this.m.getGroupWithId(ctx, id, K_GROUP_TYPE_PERMISSION)
}

// GetRoleGroupWithId 获取角色组详情，包含权限列表或者角色列表
func (this *Service) GetRoleGroupWithId(ctx, id int64) (result *Group, err error) {
	return this.m.getGroupWithId(ctx, id, K_GROUP_TYPE_ROLE)
}

// GetPermissionGroupWithName 根据组名称查询权限组信息（精确匹配），返回数据不包含该组的权限列表
func (this *Service) GetPermissionGroupWithName(ctx int64, name string) (result *Group, err error) {
	return this.m.getGroupWithName(ctx, name, K_GROUP_TYPE_PERMISSION)
}

// GetRoleGroupWithName 根据组名称查询角色组信息（精确匹配），返回数据不包含该组的角色列表
func (this *Service) GetRoleGroupWithName(ctx int64, name string) (result *Group, err error) {
	return this.m.getGroupWithName(ctx, name, K_GROUP_TYPE_ROLE)
}

// AddPermissionGroup 添加权限组
func (this *Service) AddPermissionGroup(ctx int64, name string, status int) (result *Group, err error) {
	if result, err = this.m.getGroupWithName(ctx, name, K_GROUP_TYPE_PERMISSION); err != nil {
		return nil, err
	}
	if result != nil {
		return nil, ErrGroupExists
	}
	return this.m.addGroup(ctx, K_GROUP_TYPE_PERMISSION, name, status)
}

// AddRoleGroup 添加角色组
func (this *Service) AddRoleGroup(ctx int64, name string, status int) (result *Group, err error) {
	if result, err = this.m.getGroupWithName(ctx, name, K_GROUP_TYPE_ROLE); err != nil {
		return nil, err
	}
	if result != nil {
		return nil, ErrGroupExists
	}
	return this.m.addGroup(ctx, K_GROUP_TYPE_ROLE, name, status)
}

// UpdatePermissionGroup 更新权限组的基本信息
func (this *Service) UpdatePermissionGroup(ctx int64, id int64, name string, status int) (err error) {
	result, err := this.m.getGroupWithName(ctx, name, K_GROUP_TYPE_PERMISSION)
	if err != nil {
		return err
	}
	if result != nil && result.Id != id {
		return ErrGroupExists
	}
	return this.m.updateGroup(ctx, id, name, status)
}

// UpdateRoleGroup 更新权限组的基本信息
func (this *Service) UpdateRoleGroup(ctx int64, id int64, name string, status int) (err error) {
	result, err := this.m.getGroupWithName(ctx, name, K_GROUP_TYPE_ROLE)
	if err != nil {
		return err
	}
	if result != nil && result.Id != id {
		return ErrGroupExists
	}
	return this.m.updateGroup(ctx, id, name, status)
}

// UpdateGroupStatus 更新组的状态信息
func (this *Service) UpdateGroupStatus(ctx, id int64, status int) (err error) {
	return this.m.updateGroupStatus(ctx, id, status)
}

// RemoveGroup 删除组信息
func (this *Service) RemoveGroup(ctx, id int64) (err error) {
	group, err := this.m.getGroupWithId(ctx, id, 0)
	if err != nil {
		return err
	}
	if group == nil {
		return nil
	}

	if group.Ctx != ctx {
		return ErrRemoveGroupNotAllowed
	}

	// 如果 group 下还有内容，则不能删除
	if group.Type == K_GROUP_TYPE_PERMISSION {
		pList, err := this.m.getPermissionList(ctx, []int64{id}, 0, "")
		if err != nil {
			return err
		}
		if len(pList) > 0 {
			return ErrRemoveGroupNotAllowed
		}
	} else if group.Type == K_GROUP_TYPE_ROLE {
		rList, err := this.m.getRoleList(ctx, id, 0, "")
		if err != nil {
			return err
		}
		if len(rList) > 0 {
			return ErrRemoveGroupNotAllowed
		}
	}
	return this.m.removeGroup(ctx, id)
}

// --------------------------------------------------------------------------------
// GetPermissionList 获取指定组的权限列表
func (this *Service) GetPermissionList(ctx, groupId int64, status int, keyword string) (result []*Permission, err error) {
	var groupIdList []int64
	if groupId > 0 {
		groupIdList = append(groupIdList, groupId)
	}
	return this.m.getPermissionList(ctx, groupIdList, status, keyword)
}

// GetPermissionWithId 获取权限详情
func (this *Service) GetPermissionWithId(ctx, id int64) (result *Permission, err error) {
	return this.m.getPermissionWithId(ctx, id)
}

// GetPermissionWithName 根据权限名称获取权限信息（精确匹配）
func (this *Service) GetPermissionWithName(ctx int64, name string) (result *Permission, err error) {
	return this.m.getPermissionWithName(ctx, name)
}

// GetPermissionWithIdentifier 权限权限标识符获取权限信息（精确匹配）
func (this *Service) GetPermissionWithIdentifier(ctx int64, identifier string) (result *Permission, err error) {
	return this.m.getPermissionWithIdentifier(ctx, identifier)
}

// AddPermission 添加权限
func (this *Service) AddPermission(ctx, groupId int64, name, identifier string, status int) (result *Permission, err error) {
	if this.CheckPermissionIsExists(ctx, identifier) == true {
		return nil, ErrPermissionIdentifierExists
	}
	if this.CheckPermissionNameIsExists(ctx, name) == true {
		return nil, ErrPermissionNameExists
	}

	if groupId <= 0 {
		return nil, ErrGroupNotExist
	}

	group, err := this.m.getGroupWithId(ctx, groupId, K_GROUP_TYPE_PERMISSION)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, ErrGroupNotExist
	}
	return this.m.addPermission(ctx, groupId, name, identifier, status)
}

// UpdatePermission 更新权限信息
func (this *Service) UpdatePermission(ctx, id, groupId int64, name, identifier string, status int) (err error) {
	p, err := this.m.getPermissionWithIdentifier(ctx, identifier)
	if err != nil {
		return err
	}
	if p != nil && p.Id != id {
		return ErrPermissionIdentifierExists
	}

	p, err = this.m.getPermissionWithName(ctx, name)
	if err != nil {
		return err
	}
	if p != nil && p.Id != id {
		return ErrPermissionNameExists
	}

	if groupId <= 0 {
		return ErrGroupNotExist
	}

	group, err := this.m.getGroupWithId(ctx, groupId, K_GROUP_TYPE_PERMISSION)
	if err != nil {
		return err
	}
	if group == nil {
		return ErrGroupNotExist
	}
	return this.m.updatePermission(ctx, id, groupId, name, identifier, status)
}

// CheckPermissionIsExists 验证权限标识已经是否已经存在
func (this *Service) CheckPermissionIsExists(ctx int64, identifier string) (result bool) {
	p, err := this.m.getPermissionWithIdentifier(ctx, identifier)
	if p != nil || err != nil {
		return true
	}
	return false
}

// CheckPermissionNameIsExists 验证权限名称是否已经存在
func (this *Service) CheckPermissionNameIsExists(ctx int64, name string) (result bool) {
	p, err := this.m.getPermissionWithName(ctx, name)
	if p != nil || err != nil {
		return true
	}
	return false
}

// UpdatePermissionStatus 更新权限的状态信息
func (this *Service) UpdatePermissionStatus(ctx, id int64, status int) (err error) {
	return this.m.updatePermissionStatus(ctx, id, status)
}

// --------------------------------------------------------------------------------
// GetRoleList 获取指定组的角色组列表
func (this *Service) GetRoleList(ctx, groupId int64, status int, keyword string) (result []*Role, err error) {
	return this.m.getRoleList(ctx, groupId, status, keyword)
}

// GetPermissionListWithRole 获取指定角色的权限列表
func (this *Service) GetPermissionListWithRole(ctx, roleId int64) (result []*Permission, err error) {
	return this.m.getPermissionListWithRoleId(ctx, roleId)
}

// GetRoleWithId 获取角色详情，会返回该角色拥有的权限列表
func (this *Service) GetRoleWithId(ctx, id int64) (result *Role, err error) {
	return this.m.getRoleWithId(ctx, id, true)
}

// GetRoleWithName 根据角色名称获取角色信息（精确匹配），会返回该角色拥有的权限列表
func (this *Service) GetRoleWithName(ctx int64, name string) (result *Role, err error) {
	return this.m.getRoleWithName(ctx, name, true)
}

// CheckRoleNameIsExists 检测角色名是否已经存在
func (this *Service) CheckRoleNameIsExists(ctx int64, name string) (result bool) {
	role, err := this.m.getRoleWithName(ctx, name, false)
	if role != nil || err != nil {
		return true
	}
	return false
}

// AddRole 添加角色
func (this *Service) AddRole(ctx, groupId int64, name string, status int) (result *Role, err error) {
	if this.CheckRoleNameIsExists(ctx, name) == true {
		return nil, ErrRoleNameExists
	}

	if groupId <= 0 {
		return nil, ErrGroupNotExist
	}

	group, err := this.m.getGroupWithId(ctx, groupId, K_GROUP_TYPE_ROLE)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, ErrGroupNotExist
	}
	return this.m.addRole(ctx, groupId, name, status)
}

// UpdateRole 更新角色信息
func (this *Service) UpdateRole(ctx, id, groupId int64, name string, status int) (err error) {
	role, err := this.m.getRoleWithName(ctx, name, false)
	if err != nil {
		return err
	}
	if role != nil && role.Id != id {
		return ErrRoleNameExists
	}

	if groupId <= 0 {
		return ErrGroupNotExist
	}

	group, err := this.m.getGroupWithId(ctx, groupId, K_GROUP_TYPE_ROLE)
	if err != nil {
		return err
	}
	if group == nil {
		return ErrGroupNotExist
	}
	return this.m.updateRole(ctx, id, groupId, name, status)
}

// UpdateRoleStatus 更新角色状态信息
func (this *Service) UpdateRoleStatus(ctx, id int64, status int) (err error) {
	return this.m.updateRoleStatus(ctx, id, status)
}

// --------------------------------------------------------------------------------
// GrantPermission 为角色添加权限信息
func (this *Service) GrantPermission(ctx, roleId int64, permissionIdList ...int64) (err error) {
	role, err := this.m.getRoleWithId(ctx, roleId, false)
	if err != nil {
		return err
	}
	if role == nil {
		return ErrRoleNotExist
	}
	if role.Status != K_STATUS_ENABLE {
		return ErrRoleNotExist
	}

	pList, err := this.m.getPermissionWithIdList(ctx, permissionIdList)
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
	return this.m.grantPermission(ctx, roleId, nIdList)
}

func (this *Service) RevokePermission(ctx, roleId int64, permissionIdList ...int64) (err error) {
	return this.m.revokePermission(ctx, roleId, permissionIdList)
}

// ReGrantPermission 移除之前已经授予的权限，添加新的权限
func (this *Service) ReGrantPermission(ctx, roleId int64, permissionIdList ...int64) (err error) {
	role, err := this.m.getRoleWithId(ctx, roleId, false)
	if err != nil {
		return err
	}
	if role == nil {
		return ErrRoleNotExist
	}
	if role.Status != K_STATUS_ENABLE {
		return ErrRoleNotExist
	}

	pList, err := this.m.getPermissionWithIdList(ctx, permissionIdList)
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
	return this.m.reGrantPermission(ctx, roleId, nIdList)
}

// GrantRole 为目前对象添加角色信息
func (this *Service) GrantRole(ctx int64, objectId string, roleIdList ...int64) (err error) {
	if len(roleIdList) == 0 {
		return ErrRoleNotExist
	}
	if objectId == "" {
		return ErrObjectNotAllowed
	}
	roleList, err := this.m.getRoleWithIdList(ctx, roleIdList)
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

	err = this.m.grantRole(ctx, objectId, nIdList)
	return err
}

func (this *Service) RevokeRole(ctx int64, objectId string, roleIdList ...int64) (err error) {
	return this.m.revokeRole(ctx, objectId, roleIdList)
}

// ReGrantRole 移除之前已经授予的角色，添加新的角色
func (this *Service) ReGrantRole(ctx int64, objectId string, roleIdList ...int64) (err error) {
	if len(roleIdList) == 0 {
		return ErrRoleNotExist
	}
	if objectId == "" {
		return ErrObjectNotAllowed
	}
	roleList, err := this.m.getRoleWithIdList(ctx, roleIdList)
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

	err = this.m.reGrantRole(ctx, objectId, nIdList)
	return err
}

func (this *Service) Check(ctx int64, objectId, identifier string) (result bool) {
	if this.r != nil {
		result = this.r.check(ctx, objectId, identifier)
		if result == false {
			if this.r.exists(ctx, objectId) == false {
				pList, _ := this.m.getGrantedPermissionList(ctx, objectId)
				var identifierList []interface{}
				for _, p := range pList {
					if p.Identifier == identifier {
						result = true
					}
					identifierList = append(identifierList, p.Identifier)
				}
				this.r.grantPermissions(ctx, objectId, identifierList)
			}
		}
		return result
	}
	return this.m.check(ctx, objectId, identifier)
}

func (this *Service) CheckList(ctx int64, objectId string, identifiers ...string) (result map[string]bool) {
	if this.r != nil {
		return this.r.checkList(ctx, objectId, identifiers...)
	}
	return this.m.checkList(ctx, objectId, identifiers...)
}

func (this *Service) GetGrantedRoleList(ctx int64, objectId string) (result []*Role, err error) {
	return this.m.getGrantedRoleList(ctx, objectId)
}

func (this *Service) GetGrantedPermissionList(ctx int64, objectId string) (result []*Permission, err error) {
	return this.m.getGrantedPermissionList(ctx, objectId)
}

func (this *Service) ClearCache() {
	if this.r != nil {
		this.r.clear()
	}
}
