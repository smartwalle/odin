package odin

import (
	"github.com/smartwalle/dbs"
)

type Repository interface {
	BeginTx() (dbs.TX, Repository)

	WithTx(tx dbs.TX) Repository

	// group
	GetGroupList(ctx int64, gType GroupType, status Status, keywords string) (result []*Group, err error)

	GetGroupWithId(ctx int64, gType GroupType, groupId int64) (result *Group, err error)

	GetGroupWithName(ctx int64, gType GroupType, name string) (result *Group, err error)

	AddGroup(ctx int64, gType GroupType, name, aliasName string, status Status) (result int64, err error)

	UpdateGroup(ctx int64, gType GroupType, groupId int64, aliasName string, status Status) (err error)

	UpdateGroupStatus(ctx int64, gType GroupType, groupId int64, status Status) (err error)

	// permission

	// GetPermissionList 获取角色列表，如果有传递 roleId 参数，则返回的权限数据中将附带该权限是否已授权给该 roleId
	GetPermissionList(ctx, roleId int64, status Status, keywords string, groupIds ...int64) (result []*Permission, err error)

	GetPermissionListWithIds(ctx int64, permissionIds ...int64) (result []*Permission, err error)

	GetPermissionListWithNames(ctx int64, names ...string) (result []*Permission, err error)

	GetPermissionListWithRoleId(ctx int64, roleId int64) (result []*Permission, err error)

	GetPermissionWithId(ctx, permissionId int64) (result *Permission, err error)

	GetPermissionWithName(ctx int64, name string) (result *Permission, err error)

	AddPermission(ctx, groupId int64, name, aliasName, description string, status Status) (result int64, err error)

	UpdatePermission(ctx, permissionId, groupId int64, aliasName, description string, status Status) (err error)

	UpdatePermissionStatus(ctx, permissionId int64, status Status) (err error)

	GrantPermissionWithIds(ctx, roleId int64, permissionIds ...int64) (err error)

	RevokePermissionWithIds(ctx, roleId int64, permissionIds ...int64) (err error)

	RevokeAllPermission(ctx, roleId int64) (err error)

	GetGrantedPermissionList(ctx int64, targetId string) (result []*Permission, err error)

	// role
	// GetRoleList 获取角色列表，如果有传递 targetId 参数，则返回的角色数据中将附带该角色是否已授权给该 targetId
	GetRoleList(ctx int64, targetId string, status Status, keywords string) (result []*Role, err error)

	GetRoleListWithIds(ctx int64, roleIds ...int64) (result []*Role, err error)

	GetRoleListWithNames(ctx int64, names ...string) (result []*Role, err error)

	GetRoleWithId(ctx, roleId int64) (result *Role, err error)

	GetRoleWithName(ctx int64, name string) (result *Role, err error)

	AddRole(ctx, parentId int64, name, aliasName, description string, status Status) (result int64, err error)

	UpdateRole(ctx, roleId int64, aliasName, description string, status Status) (err error)

	UpdateRoleStatus(ctx, roleId int64, status Status) (err error)

	GetGrantedRoleList(ctx int64, targetId string) (result []*Role, err error)

	GrantRoleWithIds(ctx int64, targetId string, roleIds ...int64) (err error)

	RevokeRoleWithIds(ctx int64, targetId string, roleIds ...int64) (err error)

	RevokeAllRole(ctx int64, targetId string) (err error)

	CleanCache(ctx int64, targetId string)
}

type odinService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	var s = &odinService{}
	s.repo = repo
	return s
}

// group

func (this *odinService) GetPermissionGroupList(ctx int64, status Status, keywords string) (result []*Group, err error) {
	return this.repo.GetGroupList(ctx, GroupPermission, status, keywords)
}

func (this *odinService) GetPermissionGroupWithId(ctx, groupId int64) (result *Group, err error) {
	return this.repo.GetGroupWithId(ctx, GroupPermission, groupId)
}

func (this *odinService) GetPermissionGroup(ctx int64, gType GroupType, groupName string) (result *Group, err error) {
	return this.repo.GetGroupWithName(ctx, GroupPermission, groupName)
}

func (this *odinService) addGroup(ctx int64, gType GroupType, groupName, aliasName string, status Status) (result int64, err error) {
	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 验证 groupName 是否已经存在
	group, err := nRepo.GetGroupWithName(ctx, gType, groupName)
	if err != nil {
		return 0, err
	}
	if group != nil {
		return 0, ErrGroupNameExists
	}

	if result, err = nRepo.AddGroup(ctx, gType, groupName, aliasName, status); err != nil {
		return 0, err
	}

	tx.Commit()
	return result, nil
}

func (this *odinService) AddPermissionGroup(ctx int64, groupName, aliasName string, status Status) (result int64, err error) {
	return this.addGroup(ctx, GroupPermission, groupName, aliasName, status)
}

func (this *odinService) updateGroupWithId(ctx int64, gType GroupType, groupId int64, aliasName string, status Status) (err error) {
	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 验证组是否存在
	group, err := nRepo.GetGroupWithId(ctx, gType, groupId)
	if err != nil {
		return err
	}
	if group == nil {
		return ErrGroupNotExist
	}

	if err = nRepo.UpdateGroup(ctx, gType, groupId, aliasName, status); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) UpdatePermissionGroupWithId(ctx, groupId int64, aliasName string, status Status) (err error) {
	return this.updateGroupWithId(ctx, GroupPermission, groupId, aliasName, status)
}

func (this *odinService) updateGroup(ctx int64, gType GroupType, groupName, aliasName string, status Status) (err error) {
	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 验证组是否存在
	group, err := nRepo.GetGroupWithName(ctx, gType, groupName)
	if err != nil {
		return err
	}
	if group == nil {
		return ErrGroupNotExist
	}

	if err = nRepo.UpdateGroup(ctx, gType, group.Id, aliasName, status); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) UpdatePermissionGroup(ctx int64, groupName string, aliasName string, status Status) (err error) {
	return this.updateGroup(ctx, GroupPermission, groupName, aliasName, status)
}

func (this *odinService) updateGroupStatusWithId(ctx int64, gType GroupType, groupId int64, status Status) (err error) {
	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 验证组是否存在
	group, err := nRepo.GetGroupWithId(ctx, gType, groupId)
	if err != nil {
		return err
	}
	if group == nil {
		return ErrGroupNotExist
	}

	if err = nRepo.UpdateGroupStatus(ctx, gType, groupId, status); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) UpdatePermissionGroupStatusWithId(ctx int64, groupId int64, status Status) (err error) {
	return this.updateGroupStatusWithId(ctx, GroupPermission, groupId, status)
}

func (this *odinService) updateGroupStatus(ctx int64, gType GroupType, groupName string, status Status) (err error) {
	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 验证组是否存在
	group, err := nRepo.GetGroupWithName(ctx, gType, groupName)
	if err != nil {
		return err
	}
	if group == nil {
		return ErrGroupNotExist
	}

	if err = nRepo.UpdateGroupStatus(ctx, gType, group.Id, status); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) UpdatePermissionGroupStatus(ctx int64, groupName string, status Status) (err error) {
	return this.updateGroupStatus(ctx, GroupPermission, groupName, status)
}

// permission

func (this *odinService) GetPermissionList(ctx int64, status Status, keywords string, groupIds ...int64) (result []*Permission, err error) {
	return this.repo.GetPermissionList(ctx, 0, status, keywords, groupIds...)
}

func (this *odinService) GetPermissionWithId(ctx, permissionId int64) (result *Permission, err error) {
	result, err = this.repo.GetPermissionWithId(ctx, permissionId)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (this *odinService) GetPermission(ctx int64, permissionName string) (result *Permission, err error) {
	result, err = this.repo.GetPermissionWithName(ctx, permissionName)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (this *odinService) AddPermissionWithGroupId(ctx, groupId int64, permissionName, aliasName, description string, status Status) (result int64, err error) {
	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 判断组是否存在
	group, err := nRepo.GetGroupWithId(ctx, GroupPermission, groupId)
	if err != nil {
		return 0, err
	}
	if group == nil {
		return 0, ErrGroupNotExist
	}

	// 验证 permissionName 是否已经存在
	permission, err := nRepo.GetPermissionWithName(ctx, permissionName)
	if err != nil {
		return 0, err
	}
	if permission != nil {
		return 0, ErrPermissionNameExists
	}

	if result, err = nRepo.AddPermission(ctx, group.Id, permissionName, aliasName, description, status); err != nil {
		return 0, err
	}

	tx.Commit()
	return result, nil
}

func (this *odinService) AddPermission(ctx int64, groupName, permissionName, aliasName, description string, status Status) (result int64, err error) {
	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 判断组是否存在
	group, err := nRepo.GetGroupWithName(ctx, GroupPermission, groupName)
	if err != nil {
		return 0, err
	}
	if group == nil {
		return 0, ErrGroupNotExist
	}

	// 验证 permissionName 是否已经存在
	permission, err := nRepo.GetPermissionWithName(ctx, permissionName)
	if err != nil {
		return 0, err
	}
	if permission != nil {
		return 0, ErrPermissionNameExists
	}

	if result, err = nRepo.AddPermission(ctx, group.Id, permissionName, aliasName, description, status); err != nil {
		return 0, err
	}

	tx.Commit()
	return result, nil
}

func (this *odinService) UpdatePermissionWithId(ctx, permissionId, groupId int64, aliasName, description string, status Status) (err error) {
	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 验证权限是否存在
	permission, err := nRepo.GetPermissionWithId(ctx, permissionId)
	if err != nil {
		return err
	}
	if permission == nil {
		return ErrPermissionNotExist
	}

	// 判断组是否存在
	group, err := nRepo.GetGroupWithId(ctx, GroupPermission, groupId)
	if err != nil {
		return err
	}
	if group == nil {
		return ErrGroupNotExist
	}

	if err = nRepo.UpdatePermission(ctx, permission.Id, group.Id, aliasName, description, status); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) UpdatePermission(ctx int64, permissionName, groupName, aliasName, description string, status Status) (err error) {
	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 验证权限是否存在
	permission, err := nRepo.GetPermissionWithName(ctx, permissionName)
	if err != nil {
		return err
	}
	if permission == nil {
		return ErrPermissionNotExist
	}

	// 判断组是否存在
	group, err := nRepo.GetGroupWithName(ctx, GroupPermission, groupName)
	if err != nil {
		return err
	}
	if group == nil {
		return ErrGroupNotExist
	}

	if err = nRepo.UpdatePermission(ctx, permission.Id, group.Id, aliasName, description, status); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) UpdatePermissionStatusWithId(ctx, permissionId int64, status Status) (err error) {
	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 验证权限是否存在
	permission, err := nRepo.GetPermissionWithId(ctx, permissionId)
	if err != nil {
		return err
	}
	if permission == nil {
		return ErrPermissionNotExist
	}

	if err = nRepo.UpdatePermissionStatus(ctx, permission.Id, status); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) UpdatePermissionStatus(ctx int64, permissionName string, status Status) (err error) {
	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 验证权限是否存在
	permission, err := nRepo.GetPermissionWithName(ctx, permissionName)
	if err != nil {
		return err
	}
	if permission == nil {
		return ErrPermissionNotExist
	}

	if err = nRepo.UpdatePermissionStatus(ctx, permission.Id, status); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) GrantPermissionWithId(ctx int64, roleId int64, permissionIds ...int64) (err error) {
	if len(permissionIds) == 0 {
		return ErrPermissionNotExist
	}

	if roleId <= 0 {
		return ErrRoleNotExist
	}

	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	role, err := nRepo.GetRoleWithId(ctx, roleId)
	if err != nil {
		return err
	}
	if role == nil {
		return ErrRoleNotExist
	}

	permissionList, err := nRepo.GetPermissionListWithIds(ctx, permissionIds...)
	if err != nil {
		return err
	}

	var nIds = make([]int64, 0, len(permissionList))
	for _, role := range permissionList {
		nIds = append(nIds, role.Id)
	}
	if len(nIds) == 0 {
		return ErrGrantFailed
	}

	if err = nRepo.GrantPermissionWithIds(ctx, roleId, nIds...); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) GrantPermission(ctx int64, roleName string, permissionNames ...string) (err error) {
	if len(permissionNames) == 0 {
		return ErrPermissionNotExist
	}

	if roleName == "" {
		return ErrRoleNotExist
	}

	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	role, err := nRepo.GetRoleWithName(ctx, roleName)
	if err != nil {
		return err
	}
	if role == nil {
		return ErrRoleNotExist
	}

	permissionList, err := nRepo.GetPermissionListWithNames(ctx, permissionNames...)
	if err != nil {
		return err
	}

	var nIds = make([]int64, 0, len(permissionList))
	for _, permission := range permissionList {
		nIds = append(nIds, permission.Id)
	}
	if len(nIds) == 0 {
		return ErrGrantFailed
	}

	if err = nRepo.GrantPermissionWithIds(ctx, role.Id, nIds...); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) ReGrantPermissionWithId(ctx int64, roleId int64, permissionIds ...int64) (err error) {
	if len(permissionIds) == 0 {
		return ErrPermissionNotExist
	}

	if roleId <= 0 {
		return ErrRoleNotExist
	}

	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	role, err := nRepo.GetRoleWithId(ctx, roleId)
	if err != nil {
		return err
	}
	if role == nil {
		return ErrRoleNotExist
	}

	permissionList, err := nRepo.GetPermissionListWithIds(ctx, permissionIds...)
	if err != nil {
		return err
	}

	var nIds = make([]int64, 0, len(permissionList))
	for _, role := range permissionList {
		nIds = append(nIds, role.Id)
	}
	if len(nIds) == 0 {
		return ErrGrantFailed
	}

	if err = nRepo.RevokeAllPermission(ctx, role.Id); err != nil {
		return err
	}

	if err = nRepo.GrantPermissionWithIds(ctx, role.Id, nIds...); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) ReGrantPermission(ctx int64, roleName string, permissionNames ...string) (err error) {
	if len(permissionNames) == 0 {
		return ErrPermissionNotExist
	}

	if roleName == "" {
		return ErrRoleNotExist
	}

	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	role, err := nRepo.GetRoleWithName(ctx, roleName)
	if err != nil {
		return err
	}
	if role == nil {
		return ErrRoleNotExist
	}

	permissionList, err := nRepo.GetPermissionListWithNames(ctx, permissionNames...)
	if err != nil {
		return err
	}

	var nIds = make([]int64, 0, len(permissionList))
	for _, role := range permissionList {
		nIds = append(nIds, role.Id)
	}
	if len(nIds) == 0 {
		return ErrGrantFailed
	}

	if err = nRepo.RevokeAllPermission(ctx, role.Id); err != nil {
		return err
	}

	if err = nRepo.GrantPermissionWithIds(ctx, role.Id, nIds...); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) RevokePermissionWithId(ctx int64, roleId int64, permissionIds ...int64) (err error) {
	return this.repo.RevokePermissionWithIds(ctx, roleId, permissionIds...)
}

func (this *odinService) RevokePermission(ctx int64, roleName string, permissionNames ...string) (err error) {
	if len(permissionNames) == 0 {
		return ErrPermissionNotExist
	}

	if roleName == "" {
		return ErrRoleNotExist
	}

	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	role, err := nRepo.GetRoleWithName(ctx, roleName)
	if err != nil {
		return err
	}
	if role == nil {
		return ErrRoleNotExist
	}

	permissionList, err := nRepo.GetPermissionListWithNames(ctx, permissionNames...)
	if err != nil {
		return err
	}

	var nIds = make([]int64, 0, len(permissionList))
	for _, permission := range permissionList {
		nIds = append(nIds, permission.Id)
	}
	if len(nIds) == 0 {
		return ErrRevokeFailed
	}

	if err = nRepo.RevokePermissionWithIds(ctx, role.Id, nIds...); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) RevokeAllPermissionWithId(ctx, roleId int64) (err error) {
	return this.repo.RevokeAllPermission(ctx, roleId)
}

func (this *odinService) RevokeAllPermission(ctx int64, roleName string) (err error) {
	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	role, err := nRepo.GetRoleWithName(ctx, roleName)
	if err != nil {
		return err
	}
	if role == nil {
		return ErrRoleNotExist
	}

	if err = nRepo.RevokeAllPermission(ctx, role.Id); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

// role

func (this *odinService) GetRoleList(ctx int64, targetId string, status Status, keywords string) (result []*Role, err error) {
	return this.repo.GetRoleList(ctx, targetId, status, keywords)
}

func (this *odinService) GetRoleWithId(ctx, roleId int64) (result *Role, err error) {
	result, err = this.repo.GetRoleWithId(ctx, roleId)
	if err != nil {
		return nil, err
	}
	if result != nil {
		result.PermissionList, err = this.repo.GetPermissionListWithRoleId(ctx, result.Id)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (this *odinService) GetRole(ctx int64, name string) (result *Role, err error) {
	result, err = this.repo.GetRoleWithName(ctx, name)
	if err != nil {
		return nil, err
	}
	if result != nil {
		result.PermissionList, err = this.repo.GetPermissionListWithRoleId(ctx, result.Id)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (this *odinService) AddRoleWithParentId(ctx, parentRoleId int64, roleName, aliasName, description string, status Status) (result int64, err error) {
	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if parentRoleId < 0 {
		parentRoleId = 0
	}

	// 验证 parent 是否存在
	if parentRoleId > 0 {
		parent, err := nRepo.GetRoleWithId(ctx, parentRoleId)
		if err != nil {
			return 0, err
		}
		if parent == nil {
			return 0, ErrParentRoleNotExist
		}
	}

	// 验证 name 是否已经存在
	role, err := nRepo.GetRoleWithName(ctx, roleName)
	if err != nil {
		return 0, err
	}
	if role != nil {
		return 0, ErrRoleNameExists
	}

	if result, err = nRepo.AddRole(ctx, parentRoleId, roleName, aliasName, description, status); err != nil {
		return 0, err
	}

	tx.Commit()
	return result, nil
}

func (this *odinService) AddRole(ctx int64, parentRoleName, roleName, aliasName, description string, status Status) (result int64, err error) {
	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 验证 parent 是否存在
	var parentId int64 = 0
	if parentRoleName != "" {
		parent, err := nRepo.GetRoleWithName(ctx, parentRoleName)
		if err != nil {
			return 0, err
		}
		if parent == nil {
			return 0, ErrParentRoleNotExist
		}
	}

	// 验证 name 是否已经存在
	role, err := nRepo.GetRoleWithName(ctx, roleName)
	if err != nil {
		return 0, err
	}
	if role != nil {
		return 0, ErrRoleNameExists
	}

	if result, err = nRepo.AddRole(ctx, parentId, roleName, aliasName, description, status); err != nil {
		return 0, err
	}

	tx.Commit()
	return result, nil
}

func (this *odinService) UpdateRoleWithId(ctx, roleId int64, aliasName, description string, status Status) (err error) {
	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 验证角色是否存在
	role, err := nRepo.GetRoleWithId(ctx, roleId)
	if err != nil {
		return err
	}
	if role == nil {
		return ErrRoleNotExist
	}

	if err = nRepo.UpdateRole(ctx, roleId, aliasName, description, status); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) UpdateRole(ctx int64, roleName, aliasName, description string, status Status) (err error) {
	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 验证角色是否存在
	role, err := nRepo.GetRoleWithName(ctx, roleName)
	if err != nil {
		return err
	}
	if role == nil {
		return ErrRoleNotExist
	}

	if err = nRepo.UpdateRole(ctx, role.Id, aliasName, description, status); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) UpdateRoleStatusWithId(ctx, roleId int64, status Status) (err error) {
	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 验证角色是否存在
	role, err := nRepo.GetRoleWithId(ctx, roleId)
	if err != nil {
		return err
	}
	if role == nil {
		return ErrRoleNotExist
	}

	if err = nRepo.UpdateRoleStatus(ctx, role.Id, status); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) UpdateRoleStatus(ctx int64, roleName string, status Status) (err error) {
	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 验证角色是否存在
	role, err := nRepo.GetRoleWithName(ctx, roleName)
	if err != nil {
		return err
	}
	if role == nil {
		return ErrRoleNotExist
	}

	if err = nRepo.UpdateRoleStatus(ctx, role.Id, status); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) GrantRoleWithId(ctx int64, targetId string, roleIds ...int64) (err error) {
	if len(roleIds) == 0 {
		return ErrRoleNotExist
	}

	if targetId == "" {
		return ErrTargetNotAllowed
	}

	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	roleList, err := nRepo.GetRoleListWithIds(ctx, roleIds...)
	if err != nil {
		return err
	}

	var nIds = make([]int64, 0, len(roleList))
	for _, role := range roleList {
		nIds = append(nIds, role.Id)
	}
	if len(nIds) == 0 {
		return ErrGrantFailed
	}

	if err = nRepo.GrantRoleWithIds(ctx, targetId, nIds...); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) GrantRole(ctx int64, targetId string, roleNames ...string) (err error) {
	if len(roleNames) == 0 {
		return ErrRoleNotExist
	}

	if targetId == "" {
		return ErrTargetNotAllowed
	}

	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	roleList, err := nRepo.GetRoleListWithNames(ctx, roleNames...)
	if err != nil {
		return err
	}

	var nIds = make([]int64, 0, len(roleList))
	for _, role := range roleList {
		nIds = append(nIds, role.Id)
	}
	if len(nIds) == 0 {
		return ErrGrantFailed
	}

	if err = nRepo.GrantRoleWithIds(ctx, targetId, nIds...); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) ReGrantRoleWithId(ctx int64, targetId string, roleIds ...int64) (err error) {
	if len(roleIds) == 0 {
		return ErrRoleNotExist
	}

	if targetId == "" {
		return ErrTargetNotAllowed
	}

	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	roleList, err := nRepo.GetRoleListWithIds(ctx, roleIds...)
	if err != nil {
		return err
	}

	var nIds = make([]int64, 0, len(roleList))
	for _, role := range roleList {
		nIds = append(nIds, role.Id)
	}
	if len(nIds) == 0 {
		return ErrGrantFailed
	}

	if err = nRepo.RevokeAllRole(ctx, targetId); err != nil {
		return err
	}

	if err = nRepo.GrantRoleWithIds(ctx, targetId, nIds...); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) ReGrantRole(ctx int64, targetId string, roleNames ...string) (err error) {
	if len(roleNames) == 0 {
		return ErrRoleNotExist
	}

	if targetId == "" {
		return ErrTargetNotAllowed
	}

	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	roleList, err := nRepo.GetRoleListWithNames(ctx, roleNames...)
	if err != nil {
		return err
	}

	var nIds = make([]int64, 0, len(roleList))
	for _, role := range roleList {
		nIds = append(nIds, role.Id)
	}
	if len(nIds) == 0 {
		return ErrGrantFailed
	}

	if err = nRepo.RevokeAllRole(ctx, targetId); err != nil {
		return err
	}

	if err = nRepo.GrantRoleWithIds(ctx, targetId, nIds...); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) RevokeRoleWithId(ctx int64, targetId string, roleIds ...int64) (err error) {
	if len(roleIds) == 0 {
		return ErrRoleNotExist
	}
	if targetId == "" {
		return ErrTargetNotAllowed
	}
	return this.repo.RevokeRoleWithIds(ctx, targetId, roleIds...)
}

func (this *odinService) RevokeRole(ctx int64, targetId string, roleNames ...string) (err error) {
	if len(roleNames) == 0 {
		return ErrRoleNotExist
	}

	if targetId == "" {
		return ErrTargetNotAllowed
	}

	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	roleList, err := nRepo.GetRoleListWithNames(ctx, roleNames...)
	if err != nil {
		return err
	}

	var nIds = make([]int64, 0, len(roleList))
	for _, role := range roleList {
		nIds = append(nIds, role.Id)
	}
	if len(nIds) == 0 {
		return ErrGrantFailed
	}

	if err = nRepo.RevokeRoleWithIds(ctx, targetId, nIds...); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) RevokeAllRole(ctx int64, targetId string) (err error) {
	return this.repo.RevokeAllRole(ctx, targetId)
}

//

func (this *odinService) GetPermissionListWithRoleId(ctx int64, roleId int64) (result []*Permission, err error) {
	role, err := this.repo.GetRoleWithId(ctx, roleId)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, ErrRoleNotExist
	}
	return this.repo.GetPermissionListWithRoleId(ctx, role.Id)
}

func (this *odinService) GetPermissionListWithRole(ctx int64, roleName string) (result []*Permission, err error) {
	role, err := this.repo.GetRoleWithName(ctx, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, ErrRoleNotExist
	}
	return this.repo.GetPermissionListWithRoleId(ctx, role.Id)
}

func (this *odinService) GetGrantedRoleList(ctx int64, targetId string) (result []*Role, err error) {
	return this.repo.GetGrantedRoleList(ctx, targetId)
}

func (this *odinService) GetGrantedPermissionList(ctx int64, targetId string) (result []*Permission, err error) {
	return this.repo.GetGrantedPermissionList(ctx, targetId)
}

func (this *odinService) GetPermissionTreeWithRoleId(ctx, roleId int64, status Status) (result []*Group, err error) {
	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 获取权限信息
	if roleId > 0 {
		role, err := nRepo.GetRoleWithId(ctx, roleId)
		if err != nil {
			return nil, err
		}
		if role == nil {
			return nil, ErrRoleNotExist
		}
		roleId = role.Id
	}

	groupList, err := nRepo.GetGroupList(ctx, GroupPermission, status, "")
	if err != nil {
		return nil, err
	}
	if len(groupList) == 0 {
		tx.Commit()
		return nil, nil
	}

	var groupMap = make(map[int64]*Group)
	var groupIds = make([]int64, 0, len(groupList))
	for _, group := range groupList {
		groupMap[group.Id] = group
		groupIds = append(groupIds, group.Id)
	}

	pList, err := nRepo.GetPermissionList(ctx, roleId, status, "", groupIds...)
	if err != nil {
		return nil, err
	}

	for _, p := range pList {
		var group = groupMap[p.GroupId]
		if group != nil {
			group.PermissionList = append(group.PermissionList, p)
		}
	}

	tx.Commit()
	result = groupList
	return result, nil
}

func (this *odinService) GetPermissionTree(ctx int64, roleName string, status Status) (result []*Group, err error) {
	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 获取权限信息
	var roleId int64
	if roleName != "" {
		role, err := nRepo.GetRoleWithName(ctx, roleName)
		if err != nil {
			return nil, err
		}
		if role == nil {
			return nil, ErrRoleNotExist
		}
		roleId = role.Id
	}

	groupList, err := nRepo.GetGroupList(ctx, GroupPermission, status, "")
	if err != nil {
		return nil, err
	}
	if len(groupList) == 0 {
		tx.Commit()
		return nil, nil
	}

	var groupMap = make(map[int64]*Group)
	var groupIds = make([]int64, 0, len(groupList))
	for _, group := range groupList {
		groupMap[group.Id] = group
		groupIds = append(groupIds, group.Id)
	}

	pList, err := nRepo.GetPermissionList(ctx, roleId, status, "", groupIds...)
	if err != nil {
		return nil, err
	}

	for _, p := range pList {
		var group = groupMap[p.GroupId]
		if group != nil {
			group.PermissionList = append(group.PermissionList, p)
		}
	}

	tx.Commit()
	result = groupList
	return result, nil
}
