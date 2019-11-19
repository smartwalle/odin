package service

import (
	"github.com/smartwalle/odin"
)

type OdinRepository interface {
	GetGroupListWithType(ctx int64, gType odin.GroupType, status odin.Status, name string) (result []*odin.Group, err error)

	GetGroupWithId(ctx, id int64, gType odin.GroupType) (result *odin.Group, err error)

	GetGroupWithName(ctx int64, name string, gType odin.GroupType) (result *odin.Group, err error)

	AddGroup(ctx int64, gType odin.GroupType, name string, status odin.Status) (result *odin.Group, err error)

	UpdateGroup(ctx, id int64, name string, status odin.Status) (err error)

	UpdateGroupStatus(ctx, id int64, status odin.Status) (err error)

	RemoveGroup(ctx, id int64) (err error)

	GetPermissionTree(ctx, roleId int64, status odin.Status, name string) (result []*odin.Group, err error)

	GetPermissionList(ctx int64, groupIdList []int64, status odin.Status, keyword string) (result []*odin.Permission, err error)

	GetPermissionWithIdList(ctx int64, idList []int64) (result []*odin.Permission, err error)

	GetPermissionWithId(ctx, id int64) (result *odin.Permission, err error)

	GetPermissionWithName(ctx int64, name string) (result *odin.Permission, err error)

	GetPermissionWithIdentifier(ctx int64, identifier string) (result *odin.Permission, err error)

	AddPermission(ctx int64, groupId int64, name, identifier string, status odin.Status) (result *odin.Permission, err error)

	UpdatePermission(ctx, id, groupId int64, name, identifier string, status odin.Status) (err error)

	UpdatePermissionStatus(ctx, id int64, status odin.Status) (err error)

	GetPermissionListWithRoleId(ctx, roleId int64) (result []*odin.Permission, err error)

	GetGrantedPermissionList(ctx int64, target string) (result []*odin.Permission, err error)

	GetRoleTree(ctx int64, target string, status odin.Status, name string) (result []*odin.Group, err error)

	GetRoleList(ctx, groupId int64, status odin.Status, keyword string) (result []*odin.Role, err error)

	GetRoleWithId(ctx, id int64, withPermissionList bool) (result *odin.Role, err error)

	GetRoleWithName(ctx int64, name string, withPermissionList bool) (result *odin.Role, err error)

	AddRole(ctx, groupId int64, name string, status odin.Status) (result *odin.Role, err error)

	UpdateRole(ctx, id, groupId int64, name string, status odin.Status) (err error)

	UpdateRoleStatus(ctx, id int64, status odin.Status) (err error)

	GetRoleWithIdList(ctx int64, idList []int64) (result []*odin.Role, err error)

	GrantPermission(ctx, roleId int64, permissionIdList []int64) (err error)

	RevokePermission(ctx, roleId int64, permissionIdList []int64) (err error)

	ReGrantPermission(ctx, roleId int64, permissionIdList []int64) (err error)

	GrantRole(ctx int64, target string, roleIdList []int64) (err error)

	RevokeRole(ctx int64, target string, roleIdList []int64) (err error)

	ReGrantRole(ctx int64, target string, roleIdList []int64) (err error)

	Check(ctx int64, target, identifier string) (result bool)

	CheckList(ctx int64, target string, identifiers ...string) (result map[string]bool)

	GetGrantedRoleList(ctx int64, target string) (result []*odin.Role, err error)

	ClearCache(ctx int64, target string)
}

type odinService struct {
	repo OdinRepository
}

func NewOdinService(repo OdinRepository) odin.OdinService {
	var s = &odinService{}
	s.repo = repo
	return s
}

// GetPermissionTree 获取权限组列表，会返回该组包含的权限列表
// 如果 roleId 大于 0，则会返回各权限是否有授权给该角色
func (this *odinService) GetPermissionTree(ctx, roleId int64, status odin.Status, name string) (result []*odin.Group, err error) {
	return this.repo.GetPermissionTree(ctx, roleId, status, name)
}

// GetRoleTree 获取角色组列表，会返回该组包含的角色列表
// 如果 target 不为空字符串，则会返回各角色是否有授权给该对象
func (this *odinService) GetRoleTree(ctx int64, target string, status odin.Status, name string) (result []*odin.Group, err error) {
	return this.repo.GetRoleTree(ctx, target, status, name)
}

// --------------------------------------------------------------------------------
// GetPermissionGroupList 获取权限组列表，组信息不包含权限列表
func (this *odinService) GetPermissionGroupList(ctx int64, status odin.Status, name string) (result []*odin.Group, err error) {
	return this.repo.GetGroupListWithType(ctx, odin.GroupTypeOfPermission, status, name)
}

// GetRoleGroupList 获取角色组列表，组信息不包含角色列表
func (this *odinService) GetRoleGroupList(ctx int64, status odin.Status, name string) (result []*odin.Group, err error) {
	return this.repo.GetGroupListWithType(ctx, odin.GroupTypeOfRole, status, name)
}

// GetPermissionGroupWithId 获取权限组详情，包含权限列表或者角色列表
func (this *odinService) GetPermissionGroupWithId(ctx int64, id int64) (result *odin.Group, err error) {
	return this.repo.GetGroupWithId(ctx, id, odin.GroupTypeOfPermission)
}

// GetRoleGroupWithId 获取角色组详情，包含权限列表或者角色列表
func (this *odinService) GetRoleGroupWithId(ctx, id int64) (result *odin.Group, err error) {
	return this.repo.GetGroupWithId(ctx, id, odin.GroupTypeOfRole)
}

// GetPermissionGroupWithName 根据组名称查询权限组信息（精确匹配），返回数据不包含该组的权限列表
func (this *odinService) GetPermissionGroupWithName(ctx int64, name string) (result *odin.Group, err error) {
	return this.repo.GetGroupWithName(ctx, name, odin.GroupTypeOfPermission)
}

// GetRoleGroupWithName 根据组名称查询角色组信息（精确匹配），返回数据不包含该组的角色列表
func (this *odinService) GetRoleGroupWithName(ctx int64, name string) (result *odin.Group, err error) {
	return this.repo.GetGroupWithName(ctx, name, odin.GroupTypeOfRole)
}

// AddPermissionGroup 添加权限组
func (this *odinService) AddPermissionGroup(ctx int64, name string, status odin.Status) (result *odin.Group, err error) {
	if result, err = this.repo.GetGroupWithName(ctx, name, odin.GroupTypeOfPermission); err != nil {
		return nil, err
	}
	if result != nil {
		return nil, odin.ErrGroupExists
	}
	return this.repo.AddGroup(ctx, odin.GroupTypeOfPermission, name, status)
}

// AddRoleGroup 添加角色组
func (this *odinService) AddRoleGroup(ctx int64, name string, status odin.Status) (result *odin.Group, err error) {
	if result, err = this.repo.GetGroupWithName(ctx, name, odin.GroupTypeOfRole); err != nil {
		return nil, err
	}
	if result != nil {
		return nil, odin.ErrGroupExists
	}
	return this.repo.AddGroup(ctx, odin.GroupTypeOfRole, name, status)
}

// UpdatePermissionGroup 更新权限组的基本信息
func (this *odinService) UpdatePermissionGroup(ctx int64, id int64, name string, status odin.Status) (err error) {
	result, err := this.repo.GetGroupWithName(ctx, name, odin.GroupTypeOfPermission)
	if err != nil {
		return err
	}
	if result != nil && result.Id != id {
		return odin.ErrGroupExists
	}
	return this.repo.UpdateGroup(ctx, id, name, status)
}

// UpdateRoleGroup 更新权限组的基本信息
func (this *odinService) UpdateRoleGroup(ctx int64, id int64, name string, status odin.Status) (err error) {
	result, err := this.repo.GetGroupWithName(ctx, name, odin.GroupTypeOfRole)
	if err != nil {
		return err
	}
	if result != nil && result.Id != id {
		return odin.ErrGroupExists
	}
	return this.repo.UpdateGroup(ctx, id, name, status)
}

// UpdateGroupStatus 更新组的状态信息
func (this *odinService) UpdateGroupStatus(ctx, id int64, status odin.Status) (err error) {
	return this.repo.UpdateGroupStatus(ctx, id, status)
}

// RemoveGroup 删除组信息
func (this *odinService) RemoveGroup(ctx, id int64) (err error) {
	group, err := this.repo.GetGroupWithId(ctx, id, 0)
	if err != nil {
		return err
	}
	if group == nil {
		return nil
	}

	if group.Ctx != ctx {
		return odin.ErrRemoveGroupNotAllowed
	}

	// 如果 group 下还有内容，则不能删除
	if group.Type == odin.GroupTypeOfPermission {
		pList, err := this.repo.GetPermissionList(ctx, []int64{id}, 0, "")
		if err != nil {
			return err
		}
		if len(pList) > 0 {
			return odin.ErrRemoveGroupNotAllowed
		}
	} else if group.Type == odin.GroupTypeOfRole {
		rList, err := this.repo.GetRoleList(ctx, id, 0, "")
		if err != nil {
			return err
		}
		if len(rList) > 0 {
			return odin.ErrRemoveGroupNotAllowed
		}
	}
	return this.repo.RemoveGroup(ctx, id)
}

// --------------------------------------------------------------------------------
// GetPermissionList 获取指定组的权限列表
func (this *odinService) GetPermissionList(ctx, groupId int64, status odin.Status, keyword string) (result []*odin.Permission, err error) {
	var groupIdList []int64
	if groupId > 0 {
		groupIdList = append(groupIdList, groupId)
	}
	return this.repo.GetPermissionList(ctx, groupIdList, status, keyword)
}

// GetPermissionWithId 获取权限详情
func (this *odinService) GetPermissionWithId(ctx, id int64) (result *odin.Permission, err error) {
	return this.repo.GetPermissionWithId(ctx, id)
}

// GetPermissionWithName 根据权限名称获取权限信息（精确匹配）
func (this *odinService) GetPermissionWithName(ctx int64, name string) (result *odin.Permission, err error) {
	return this.repo.GetPermissionWithName(ctx, name)
}

// GetPermissionWithIdentifier 权限权限标识符获取权限信息（精确匹配）
func (this *odinService) GetPermissionWithIdentifier(ctx int64, identifier string) (result *odin.Permission, err error) {
	return this.repo.GetPermissionWithIdentifier(ctx, identifier)
}

// AddPermission 添加权限
func (this *odinService) AddPermission(ctx, groupId int64, name, identifier string, status odin.Status) (result *odin.Permission, err error) {
	if this.CheckPermissionIsExists(ctx, identifier) == true {
		return nil, odin.ErrPermissionIdentifierExists
	}
	if this.CheckPermissionNameIsExists(ctx, name) == true {
		return nil, odin.ErrPermissionNameExists
	}

	if groupId <= 0 {
		return nil, odin.ErrGroupNotExist
	}

	group, err := this.repo.GetGroupWithId(ctx, groupId, odin.GroupTypeOfPermission)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, odin.ErrGroupNotExist
	}
	return this.repo.AddPermission(ctx, groupId, name, identifier, status)
}

// UpdatePermission 更新权限信息
func (this *odinService) UpdatePermission(ctx, id, groupId int64, name, identifier string, status odin.Status) (err error) {
	p, err := this.repo.GetPermissionWithIdentifier(ctx, identifier)
	if err != nil {
		return err
	}
	if p != nil && p.Id != id {
		return odin.ErrPermissionIdentifierExists
	}

	p, err = this.repo.GetPermissionWithName(ctx, name)
	if err != nil {
		return err
	}
	if p != nil && p.Id != id {
		return odin.ErrPermissionNameExists
	}

	if groupId <= 0 {
		return odin.ErrGroupNotExist
	}

	group, err := this.repo.GetGroupWithId(ctx, groupId, odin.GroupTypeOfPermission)
	if err != nil {
		return err
	}
	if group == nil {
		return odin.ErrGroupNotExist
	}
	return this.repo.UpdatePermission(ctx, id, groupId, name, identifier, status)
}

// CheckPermissionIsExists 验证权限标识已经是否已经存在
func (this *odinService) CheckPermissionIsExists(ctx int64, identifier string) (result bool) {
	p, err := this.repo.GetPermissionWithIdentifier(ctx, identifier)
	if p != nil || err != nil {
		return true
	}
	return false
}

// CheckPermissionNameIsExists 验证权限名称是否已经存在
func (this *odinService) CheckPermissionNameIsExists(ctx int64, name string) (result bool) {
	p, err := this.repo.GetPermissionWithName(ctx, name)
	if p != nil || err != nil {
		return true
	}
	return false
}

// UpdatePermissionStatus 更新权限的状态信息
func (this *odinService) UpdatePermissionStatus(ctx, id int64, status odin.Status) (err error) {
	return this.repo.UpdatePermissionStatus(ctx, id, status)
}

// --------------------------------------------------------------------------------
// GetRoleList 获取指定组的角色组列表
func (this *odinService) GetRoleList(ctx, groupId int64, status odin.Status, keyword string) (result []*odin.Role, err error) {
	return this.repo.GetRoleList(ctx, groupId, status, keyword)
}

// GetPermissionListWithRole 获取指定角色的权限列表
func (this *odinService) GetPermissionListWithRole(ctx, roleId int64) (result []*odin.Permission, err error) {
	return this.repo.GetPermissionListWithRoleId(ctx, roleId)
}

// GetRoleWithId 获取角色详情，会返回该角色拥有的权限列表
func (this *odinService) GetRoleWithId(ctx, id int64) (result *odin.Role, err error) {
	return this.repo.GetRoleWithId(ctx, id, true)
}

// GetRoleWithName 根据角色名称获取角色信息（精确匹配），会返回该角色拥有的权限列表
func (this *odinService) GetRoleWithName(ctx int64, name string) (result *odin.Role, err error) {
	return this.repo.GetRoleWithName(ctx, name, true)
}

// CheckRoleNameIsExists 检测角色名是否已经存在
func (this *odinService) CheckRoleNameIsExists(ctx int64, name string) (result bool) {
	role, err := this.repo.GetRoleWithName(ctx, name, false)
	if role != nil || err != nil {
		return true
	}
	return false
}

// AddRole 添加角色
func (this *odinService) AddRole(ctx, groupId int64, name string, status odin.Status) (result *odin.Role, err error) {
	if this.CheckRoleNameIsExists(ctx, name) == true {
		return nil, odin.ErrRoleNameExists
	}

	if groupId <= 0 {
		return nil, odin.ErrGroupNotExist
	}

	group, err := this.repo.GetGroupWithId(ctx, groupId, odin.GroupTypeOfRole)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, odin.ErrGroupNotExist
	}
	return this.repo.AddRole(ctx, groupId, name, status)
}

// UpdateRole 更新角色信息
func (this *odinService) UpdateRole(ctx, id, groupId int64, name string, status odin.Status) (err error) {
	role, err := this.repo.GetRoleWithName(ctx, name, false)
	if err != nil {
		return err
	}
	if role != nil && role.Id != id {
		return odin.ErrRoleNameExists
	}

	if groupId <= 0 {
		return odin.ErrGroupNotExist
	}

	group, err := this.repo.GetGroupWithId(ctx, groupId, odin.GroupTypeOfRole)
	if err != nil {
		return err
	}
	if group == nil {
		return odin.ErrGroupNotExist
	}
	return this.repo.UpdateRole(ctx, id, groupId, name, status)
}

// UpdateRoleStatus 更新角色状态信息
func (this *odinService) UpdateRoleStatus(ctx, id int64, status odin.Status) (err error) {
	return this.repo.UpdateRoleStatus(ctx, id, status)
}

// --------------------------------------------------------------------------------
// GrantPermission 为角色添加权限信息
func (this *odinService) GrantPermission(ctx, roleId int64, permissionIdList ...int64) (err error) {
	role, err := this.repo.GetRoleWithId(ctx, roleId, false)
	if err != nil {
		return err
	}
	if role == nil {
		return odin.ErrRoleNotExist
	}
	if role.Status != odin.StatusOfEnable {
		return odin.ErrRoleNotExist
	}

	pList, err := this.repo.GetPermissionWithIdList(ctx, permissionIdList)
	if err != nil {
		return err
	}
	var nIdList []int64
	for _, p := range pList {
		if p.Status == odin.StatusOfEnable {
			nIdList = append(nIdList, p.Id)
		}
	}
	if len(nIdList) == 0 {
		return odin.ErrGrantFailed
	}
	return this.repo.GrantPermission(ctx, roleId, nIdList)
}

func (this *odinService) RevokePermission(ctx, roleId int64, permissionIdList ...int64) (err error) {
	return this.repo.RevokePermission(ctx, roleId, permissionIdList)
}

// ReGrantPermission 移除之前已经授予的权限，添加新的权限
func (this *odinService) ReGrantPermission(ctx, roleId int64, permissionIdList ...int64) (err error) {
	role, err := this.repo.GetRoleWithId(ctx, roleId, false)
	if err != nil {
		return err
	}
	if role == nil {
		return odin.ErrRoleNotExist
	}
	if role.Status != odin.StatusOfEnable {
		return odin.ErrRoleNotExist
	}

	pList, err := this.repo.GetPermissionWithIdList(ctx, permissionIdList)
	if err != nil {
		return err
	}
	var nIdList []int64
	for _, p := range pList {
		if p.Status == odin.StatusOfEnable {
			nIdList = append(nIdList, p.Id)
		}
	}
	return this.repo.ReGrantPermission(ctx, roleId, nIdList)
}

// GrantRole 为目前对象添加角色信息
func (this *odinService) GrantRole(ctx int64, target string, roleIdList ...int64) (err error) {
	if len(roleIdList) == 0 {
		return odin.ErrRoleNotExist
	}
	if target == "" {
		return odin.ErrObjectNotAllowed
	}
	roleList, err := this.repo.GetRoleWithIdList(ctx, roleIdList)
	if err != nil {
		return err
	}

	var nIdList []int64
	for _, role := range roleList {
		if role.Status == odin.StatusOfEnable {
			nIdList = append(nIdList, role.Id)
		}
	}
	if len(nIdList) == 0 {
		return odin.ErrGrantFailed
	}

	err = this.repo.GrantRole(ctx, target, nIdList)
	return err
}

func (this *odinService) RevokeRole(ctx int64, target string, roleIdList ...int64) (err error) {
	return this.repo.RevokeRole(ctx, target, roleIdList)
}

// ReGrantRole 移除之前已经授予的角色，添加新的角色
func (this *odinService) ReGrantRole(ctx int64, target string, roleIdList ...int64) (err error) {
	if len(roleIdList) == 0 {
		return odin.ErrRoleNotExist
	}
	if target == "" {
		return odin.ErrObjectNotAllowed
	}
	roleList, err := this.repo.GetRoleWithIdList(ctx, roleIdList)
	if err != nil {
		return err
	}

	var nIdList []int64
	for _, role := range roleList {
		if role.Status == odin.StatusOfEnable {
			nIdList = append(nIdList, role.Id)
		}
	}

	err = this.repo.ReGrantRole(ctx, target, nIdList)
	return err
}

func (this *odinService) Check(ctx int64, target, identifier string) (result bool) {
	return this.repo.Check(ctx, target, identifier)
}

func (this *odinService) CheckList(ctx int64, target string, identifiers ...string) (result map[string]bool) {
	return this.repo.CheckList(ctx, target, identifiers...)
}

func (this *odinService) GetGrantedRoleList(ctx int64, target string) (result []*odin.Role, err error) {
	return this.repo.GetGrantedRoleList(ctx, target)
}

func (this *odinService) GetGrantedPermissionList(ctx int64, target string) (result []*odin.Permission, err error) {
	return this.repo.GetGrantedPermissionList(ctx, target)
}

func (this *odinService) ClearCache(ctx int64, target string) {
	this.repo.ClearCache(ctx, target)
}
