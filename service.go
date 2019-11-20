package odin

import (
	"github.com/smartwalle/dbs"
)

type Repository interface {
	BeginTx() (dbs.TX, Repository)

	WithTx(tx dbs.TX) Repository

	GetGroupList(ctx int64, gType GroupType, status Status, name string) (result []*Group, err error)

	GetGroupWithId(ctx, id int64, gType GroupType) (result *Group, err error)

	GetGroupWithName(ctx int64, name string, gType GroupType) (result *Group, err error)

	AddGroup(ctx int64, gType GroupType, name string, status Status) (result int64, err error)

	UpdateGroup(ctx, id int64, name string, status Status) (err error)

	UpdateGroupStatus(ctx, id int64, status Status) (err error)

	RemoveGroup(ctx, id int64) (err error)

	GetPermissionList(ctx int64, groupIdList []int64, status Status, keyword string, roleId int64) (result []*Permission, err error)

	GetPermissionListWithIds(ctx int64, idList []int64) (result []*Permission, err error)

	GetPermissionWithId(ctx, id int64) (result *Permission, err error)

	GetPermissionWithName(ctx int64, name string) (result *Permission, err error)

	GetPermissionWithIdentifier(ctx int64, identifier string) (result *Permission, err error)

	AddPermission(ctx int64, groupId int64, name, identifier string, status Status) (result int64, err error)

	UpdatePermission(ctx, id, groupId int64, name, identifier string, status Status) (err error)

	UpdatePermissionStatus(ctx, id int64, status Status) (err error)

	GetRolePermissionList(ctx, roleId int64) (result []*Permission, err error)

	GetGrantedPermissionList(ctx int64, target string) (result []*Permission, err error)

	GetRoleList(ctx int64, target string, groupIdList []int64, status Status, keyword string) (result []*Role, err error)

	GetRoleWithId(ctx, id int64, withPermissionList bool) (result *Role, err error)

	GetRoleWithName(ctx int64, name string, withPermissionList bool) (result *Role, err error)

	AddRole(ctx, groupId int64, name string, status Status) (result int64, err error)

	UpdateRole(ctx, id, groupId int64, name string, status Status) (err error)

	UpdateRoleStatus(ctx, id int64, status Status) (err error)

	GetRoleWithIdList(ctx int64, idList []int64) (result []*Role, err error)

	GrantPermission(ctx, roleId int64, permissionIdList []int64) (err error)

	RevokePermission(ctx, roleId int64, permissionIdList []int64) (err error)

	ReGrantPermission(ctx, roleId int64, permissionIdList []int64) (err error)

	GrantRole(ctx int64, target string, roleIdList []int64) (err error)

	RevokeRole(ctx int64, target string, roleIdList []int64) (err error)

	ReGrantRole(ctx int64, target string, roleIdList []int64) (err error)

	Check(ctx int64, target, identifier string) (result bool)

	CheckList(ctx int64, target string, identifiers ...string) (result map[string]bool)

	GetGrantedRoleList(ctx int64, target string) (result []*Role, err error)

	CleanCache(ctx int64, target string)
}

type odinService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	var s = &odinService{}
	s.repo = repo
	return s
}

// GetPermissionTree 获取权限组列表，会返回该组包含的权限列表
// 如果 roleId 大于 0，则会返回各权限是否有授权给该角色
func (this *odinService) GetPermissionTree(ctx, roleId int64, status Status, name string) (result []*Group, err error) {
	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if result, err = nRepo.GetGroupList(ctx, GroupTypeOfPermission, status, name); err != nil {
		return nil, err
	}

	var gMap = make(map[int64]*Group)
	var gIdList = make([]int64, 0, len(result))
	for _, group := range result {
		gMap[group.Id] = group
		gIdList = append(gIdList, group.Id)
	}

	pList, err := nRepo.GetPermissionList(ctx, gIdList, status, "", roleId)
	if err != nil {
		return nil, err
	}

	for _, p := range pList {
		var group = gMap[p.GroupId]
		if group != nil {
			group.PermissionList = append(group.PermissionList, p)
		}
	}

	tx.Commit()
	return result, nil
}

// GetRoleTree 获取角色组列表，会返回该组包含的角色列表
// 如果 target 不为空字符串，则会返回各角色是否有授权给该对象
func (this *odinService) GetRoleTree(ctx int64, target string, status Status, name string) (result []*Group, err error) {
	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	result, err = nRepo.GetGroupList(ctx, GroupTypeOfRole, status, name)
	if err != nil {
		return nil, err
	}

	var gMap = make(map[int64]*Group)
	var gIdList []int64
	for _, group := range result {
		gMap[group.Id] = group
		gIdList = append(gIdList, group.Id)
	}

	rList, err := nRepo.GetRoleList(ctx, target, gIdList, status, "")
	if err != nil {
		return nil, err
	}

	for _, r := range rList {
		var group = gMap[r.GroupId]
		if group != nil {
			group.RoleList = append(group.RoleList, r)
		}
	}

	tx.Commit()
	return result, nil
}

// --------------------------------------------------------------------------------
// GetPermissionGroupList 获取权限组列表，组信息不包含权限列表
func (this *odinService) GetPermissionGroupList(ctx int64, status Status, name string) (result []*Group, err error) {
	return this.repo.GetGroupList(ctx, GroupTypeOfPermission, status, name)
}

// GetRoleGroupList 获取角色组列表，组信息不包含角色列表
func (this *odinService) GetRoleGroupList(ctx int64, status Status, name string) (result []*Group, err error) {
	return this.repo.GetGroupList(ctx, GroupTypeOfRole, status, name)
}

// GetPermissionGroupWithId 获取权限组详情，包含权限列表或者角色列表
func (this *odinService) GetPermissionGroupWithId(ctx int64, id int64) (result *Group, err error) {
	return this.repo.GetGroupWithId(ctx, id, GroupTypeOfPermission)
}

// GetRoleGroupWithId 获取角色组详情，包含权限列表或者角色列表
func (this *odinService) GetRoleGroupWithId(ctx, id int64) (result *Group, err error) {
	return this.repo.GetGroupWithId(ctx, id, GroupTypeOfRole)
}

// GetPermissionGroupWithName 根据组名称查询权限组信息（精确匹配），返回数据不包含该组的权限列表
func (this *odinService) GetPermissionGroupWithName(ctx int64, name string) (result *Group, err error) {
	return this.repo.GetGroupWithName(ctx, name, GroupTypeOfPermission)
}

// GetRoleGroupWithName 根据组名称查询角色组信息（精确匹配），返回数据不包含该组的角色列表
func (this *odinService) GetRoleGroupWithName(ctx int64, name string) (result *Group, err error) {
	return this.repo.GetGroupWithName(ctx, name, GroupTypeOfRole)
}

// AddPermissionGroup 添加权限组
func (this *odinService) AddPermissionGroup(ctx int64, name string, status Status) (result *Group, err error) {
	return this.addGroup(ctx, GroupTypeOfPermission, name, status)
}

// AddRoleGroup 添加角色组
func (this *odinService) AddRoleGroup(ctx int64, name string, status Status) (result *Group, err error) {
	return this.addGroup(ctx, GroupTypeOfRole, name, status)
}

func (this *odinService) addGroup(ctx int64, gType GroupType, name string, status Status) (result *Group, err error) {
	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 判断组名称是否已经存在
	if result, err = nRepo.GetGroupWithName(ctx, name, gType); err != nil {
		return nil, err
	}
	if result != nil {
		return nil, ErrGroupExists
	}

	// 添加组信息
	nId, err := nRepo.AddGroup(ctx, gType, name, status)
	if err != nil {
		return nil, err
	}

	// 获取新添加的组
	if result, err = nRepo.GetGroupWithId(ctx, nId, gType); err != nil {
		return nil, err
	}

	tx.Commit()
	return result, nil
}

// UpdatePermissionGroup 更新权限组的基本信息
func (this *odinService) UpdatePermissionGroup(ctx int64, id int64, name string, status Status) (err error) {
	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	result, err := nRepo.GetGroupWithName(ctx, name, GroupTypeOfPermission)
	if err != nil {
		return err
	}
	if result != nil && result.Id != id {
		return ErrGroupExists
	}
	if err = nRepo.UpdateGroup(ctx, id, name, status); err != nil {
		return err
	}
	tx.Commit()
	return nil
}

// UpdateRoleGroup 更新权限组的基本信息
func (this *odinService) UpdateRoleGroup(ctx int64, id int64, name string, status Status) (err error) {
	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	result, err := nRepo.GetGroupWithName(ctx, name, GroupTypeOfRole)
	if err != nil {
		return err
	}
	if result != nil && result.Id != id {
		return ErrGroupExists
	}
	if err = nRepo.UpdateGroup(ctx, id, name, status); err != nil {
		return err
	}
	tx.Commit()
	return nil
}

// UpdateGroupStatus 更新组的状态信息
func (this *odinService) UpdateGroupStatus(ctx, id int64, status Status) (err error) {
	return this.repo.UpdateGroupStatus(ctx, id, status)
}

// RemoveGroup 删除组信息
func (this *odinService) RemoveGroup(ctx, id int64) (err error) {
	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	group, err := nRepo.GetGroupWithId(ctx, id, 0)
	if err != nil {
		return err
	}
	if group == nil {
		return ErrGroupNotExist
	}

	if group.Ctx != ctx {
		return ErrRemoveGroupNotAllowed
	}

	// 如果 group 下还有内容，则不能删除
	if group.Type == GroupTypeOfPermission {
		pList, err := nRepo.GetPermissionList(ctx, []int64{id}, 0, "", 0)
		if err != nil {
			return err
		}
		if len(pList) > 0 {
			return ErrRemoveGroupNotAllowed
		}
	} else if group.Type == GroupTypeOfRole {
		rList, err := nRepo.GetRoleList(ctx, "", []int64{id}, 0, "")
		if err != nil {
			return err
		}
		if len(rList) > 0 {
			return ErrRemoveGroupNotAllowed
		}
	}
	if err = nRepo.RemoveGroup(ctx, id); err != nil {
		return err
	}
	tx.Commit()

	return nil
}

// --------------------------------------------------------------------------------
// GetPermissionList 获取指定组的权限列表
func (this *odinService) GetPermissionList(ctx, groupId int64, status Status, keyword string) (result []*Permission, err error) {
	var groupIdList []int64
	if groupId > 0 {
		groupIdList = append(groupIdList, groupId)
	}
	return this.repo.GetPermissionList(ctx, groupIdList, status, keyword, 0)
}

// GetPermissionWithId 获取权限详情
func (this *odinService) GetPermissionWithId(ctx, id int64) (result *Permission, err error) {
	return this.repo.GetPermissionWithId(ctx, id)
}

// GetPermissionWithName 根据权限名称获取权限信息（精确匹配）
func (this *odinService) GetPermissionWithName(ctx int64, name string) (result *Permission, err error) {
	return this.repo.GetPermissionWithName(ctx, name)
}

// GetPermissionWithIdentifier 权限权限标识符获取权限信息（精确匹配）
func (this *odinService) GetPermissionWithIdentifier(ctx int64, identifier string) (result *Permission, err error) {
	return this.repo.GetPermissionWithIdentifier(ctx, identifier)
}

// AddPermission 添加权限
func (this *odinService) AddPermission(ctx, groupId int64, name, identifier string, status Status) (result *Permission, err error) {
	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	p, err := nRepo.GetPermissionWithIdentifier(ctx, identifier)
	if err != nil {
		return nil, err
	}
	if p != nil {
		return nil, ErrPermissionIdentifierExists
	}

	p, err = nRepo.GetPermissionWithName(ctx, name)
	if err != nil {
		return nil, err
	}
	if p != nil {
		return nil, ErrPermissionNameExists
	}

	if groupId <= 0 {
		return nil, ErrGroupNotExist
	}

	group, err := nRepo.GetGroupWithId(ctx, groupId, GroupTypeOfPermission)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, ErrGroupNotExist
	}
	nId, err := nRepo.AddPermission(ctx, groupId, name, identifier, status)
	if err != nil {
		return nil, err
	}

	result, err = nRepo.GetPermissionWithId(ctx, nId)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return result, nil
}

// UpdatePermission 更新权限信息
func (this *odinService) UpdatePermission(ctx, id, groupId int64, name, identifier string, status Status) (err error) {
	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	p, err := nRepo.GetPermissionWithIdentifier(ctx, identifier)
	if err != nil {
		return err
	}
	if p != nil && p.Id != id {
		return ErrPermissionIdentifierExists
	}

	p, err = nRepo.GetPermissionWithName(ctx, name)
	if err != nil {
		return err
	}
	if p != nil && p.Id != id {
		return ErrPermissionNameExists
	}

	if groupId <= 0 {
		return ErrGroupNotExist
	}

	group, err := nRepo.GetGroupWithId(ctx, groupId, GroupTypeOfPermission)
	if err != nil {
		return err
	}
	if group == nil {
		return ErrGroupNotExist
	}
	if err = nRepo.UpdatePermission(ctx, id, groupId, name, identifier, status); err != nil {
		return err
	}
	tx.Commit()
	return nil
}

// UpdatePermissionStatus 更新权限的状态信息
func (this *odinService) UpdatePermissionStatus(ctx, id int64, status Status) (err error) {
	return this.repo.UpdatePermissionStatus(ctx, id, status)
}

// CheckPermissionExists 验证权限标识已经是否已经存在
func (this *odinService) CheckPermissionExists(ctx int64, identifier string) (result bool) {
	p, _ := this.repo.GetPermissionWithIdentifier(ctx, identifier)
	if p != nil {
		return true
	}
	return false
}

// CheckPermissionNameExists 验证权限名称是否已经存在
func (this *odinService) CheckPermissionNameExists(ctx int64, name string) (result bool) {
	p, _ := this.repo.GetPermissionWithName(ctx, name)
	if p != nil {
		return true
	}
	return false
}

// --------------------------------------------------------------------------------
// GetRoleList 获取指定组的角色组列表
func (this *odinService) GetRoleList(ctx, groupId int64, status Status, keyword string) (result []*Role, err error) {
	return this.repo.GetRoleList(ctx, "", []int64{groupId}, status, keyword)
}

// GetRolePermissionList 获取指定角色的权限列表
func (this *odinService) GetRolePermissionList(ctx, roleId int64) (result []*Permission, err error) {
	return this.repo.GetRolePermissionList(ctx, roleId)
}

// GetRoleWithId 获取角色详情，会返回该角色拥有的权限列表
func (this *odinService) GetRoleWithId(ctx, id int64) (result *Role, err error) {
	return this.repo.GetRoleWithId(ctx, id, true)
}

// GetRoleWithName 根据角色名称获取角色信息（精确匹配），会返回该角色拥有的权限列表
func (this *odinService) GetRoleWithName(ctx int64, name string) (result *Role, err error) {
	return this.repo.GetRoleWithName(ctx, name, true)
}

// CheckRoleNameExists 检测角色名是否已经存在
func (this *odinService) CheckRoleNameExists(ctx int64, name string) (result bool) {
	role, err := this.repo.GetRoleWithName(ctx, name, false)
	if role != nil || err != nil {
		return true
	}
	return false
}

// AddRole 添加角色
func (this *odinService) AddRole(ctx, groupId int64, name string, status Status) (result *Role, err error) {
	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if groupId <= 0 {
		return nil, ErrGroupNotExist
	}

	group, err := nRepo.GetGroupWithId(ctx, groupId, GroupTypeOfRole)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, ErrGroupNotExist
	}

	role, err := nRepo.GetRoleWithName(ctx, name, false)
	if err != nil {
		return nil, err
	}
	if role != nil {
		return nil, ErrRoleNameExists
	}

	nId, err := nRepo.AddRole(ctx, groupId, name, status)
	if err != nil {
		return nil, err
	}

	result, err = nRepo.GetRoleWithId(ctx, nId, false)
	if err != nil {
		return nil, err
	}

	tx.Commit()
	return result, err
}

// UpdateRole 更新角色信息
func (this *odinService) UpdateRole(ctx, id, groupId int64, name string, status Status) (err error) {
	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	role, err := nRepo.GetRoleWithName(ctx, name, false)
	if err != nil {
		return err
	}
	if role != nil && role.Id != id {
		return ErrRoleNameExists
	}

	if groupId <= 0 {
		return ErrGroupNotExist
	}

	group, err := nRepo.GetGroupWithId(ctx, groupId, GroupTypeOfRole)
	if err != nil {
		return err
	}
	if group == nil {
		return ErrGroupNotExist
	}
	err = nRepo.UpdateRole(ctx, id, groupId, name, status)
	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}

// UpdateRoleStatus 更新角色状态信息
func (this *odinService) UpdateRoleStatus(ctx, id int64, status Status) (err error) {
	return this.repo.UpdateRoleStatus(ctx, id, status)
}

// --------------------------------------------------------------------------------
// GrantPermission 为角色添加权限信息
func (this *odinService) GrantPermission(ctx, roleId int64, permissionIdList ...int64) (err error) {
	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	role, err := nRepo.GetRoleWithId(ctx, roleId, false)
	if err != nil {
		return err
	}
	if role == nil {
		return ErrRoleNotExist
	}
	if role.Status != StatusOfEnable {
		return ErrRoleNotExist
	}

	pList, err := nRepo.GetPermissionListWithIds(ctx, permissionIdList)
	if err != nil {
		return err
	}
	var nIdList []int64
	for _, p := range pList {
		if p.Status == StatusOfEnable {
			nIdList = append(nIdList, p.Id)
		}
	}
	if len(nIdList) == 0 {
		return ErrGrantFailed
	}
	if err = nRepo.GrantPermission(ctx, roleId, nIdList); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) RevokePermission(ctx, roleId int64, permissionIdList ...int64) (err error) {
	return this.repo.RevokePermission(ctx, roleId, permissionIdList)
}

// ReGrantPermission 移除之前已经授予的权限，添加新的权限
func (this *odinService) ReGrantPermission(ctx, roleId int64, permissionIdList ...int64) (err error) {
	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	role, err := nRepo.GetRoleWithId(ctx, roleId, false)
	if err != nil {
		return err
	}
	if role == nil {
		return ErrRoleNotExist
	}
	if role.Status != StatusOfEnable {
		return ErrRoleNotExist
	}

	pList, err := nRepo.GetPermissionListWithIds(ctx, permissionIdList)
	if err != nil {
		return err
	}
	var nIdList []int64
	for _, p := range pList {
		if p.Status == StatusOfEnable {
			nIdList = append(nIdList, p.Id)
		}
	}
	if err = nRepo.ReGrantPermission(ctx, roleId, nIdList); err != nil {
		return err
	}
	tx.Commit()
	return nil
}

// GrantRole 为目前对象添加角色信息
func (this *odinService) GrantRole(ctx int64, target string, roleIdList ...int64) (err error) {
	if len(roleIdList) == 0 {
		return ErrRoleNotExist
	}
	if target == "" {
		return ErrObjectNotAllowed
	}

	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	roleList, err := nRepo.GetRoleWithIdList(ctx, roleIdList)
	if err != nil {
		return err
	}

	var nIdList []int64
	for _, role := range roleList {
		if role.Status == StatusOfEnable {
			nIdList = append(nIdList, role.Id)
		}
	}
	if len(nIdList) == 0 {
		return ErrGrantFailed
	}

	if err = nRepo.GrantRole(ctx, target, nIdList); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) RevokeRole(ctx int64, target string, roleIdList ...int64) (err error) {
	return this.repo.RevokeRole(ctx, target, roleIdList)
}

// ReGrantRole 移除之前已经授予的角色，添加新的角色
func (this *odinService) ReGrantRole(ctx int64, target string, roleIdList ...int64) (err error) {
	if len(roleIdList) == 0 {
		return ErrRoleNotExist
	}
	if target == "" {
		return ErrObjectNotAllowed
	}

	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	roleList, err := nRepo.GetRoleWithIdList(ctx, roleIdList)
	if err != nil {
		return err
	}

	var nIdList []int64
	for _, role := range roleList {
		if role.Status == StatusOfEnable {
			nIdList = append(nIdList, role.Id)
		}
	}

	if err = nRepo.ReGrantRole(ctx, target, nIdList); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) Check(ctx int64, target, identifier string) (result bool) {
	return this.repo.Check(ctx, target, identifier)
}

func (this *odinService) CheckList(ctx int64, target string, identifiers ...string) (result map[string]bool) {
	return this.repo.CheckList(ctx, target, identifiers...)
}

func (this *odinService) GetGrantedRoleList(ctx int64, target string) (result []*Role, err error) {
	return this.repo.GetGrantedRoleList(ctx, target)
}

func (this *odinService) GetGrantedPermissionList(ctx int64, target string) (result []*Permission, err error) {
	return this.repo.GetGrantedPermissionList(ctx, target)
}

func (this *odinService) CleanCache(ctx int64, target string) {
	this.repo.CleanCache(ctx, target)
}
