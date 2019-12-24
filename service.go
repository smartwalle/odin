package odin

import (
	"fmt"
	"github.com/smartwalle/dbs"
)

type Repository interface {
	BeginTx() (dbs.TX, Repository)

	WithTx(tx dbs.TX) Repository

	// InitTable 初始化数据库表
	InitTable() error

	// GetGroups 获取组列表
	GetGroups(ctx int64, gType GroupType, status Status, keywords string) (result []*Group, err error)

	// GetGroupWithId 获取组信息
	GetGroupWithId(ctx int64, gType GroupType, groupId int64) (result *Group, err error)

	// GetGroupWithName 获取组信息
	GetGroupWithName(ctx int64, gType GroupType, name string) (result *Group, err error)

	// AddGroup 添加组信息
	AddGroup(ctx int64, gType GroupType, name, aliasName string, status Status) (result int64, err error)

	// UpdateGroup 更新组信息
	UpdateGroup(ctx int64, gType GroupType, groupId int64, aliasName string, status Status) (err error)

	// UpdateGroupStatus 更新组状态
	UpdateGroupStatus(ctx int64, gType GroupType, groupId int64, status Status) (err error)

	// GetPermissions 获取角色列表
	// 如果参数 limitedInRole 的值大于 0，则返回的权限数据将限定在已授权给 limitedInRole 的权限范围之内
	// 如果参数 isGrantedToRole 的值大于 0，则返回的权限数据中将附带该权限是否已授权给该 isGrantedToRole
	GetPermissions(ctx int64, status Status, keywords string, groupIds []int64, limitedInRole, isGrantedToRole int64) (result []*Permission, err error)

	// GetPermissionsWithIds 根据权限 id 列表获取权限信息
	GetPermissionsWithIds(ctx int64, permissionIds ...int64) (result []*Permission, err error)

	// GetPermissionsWithNames 根据权限名称列表获取权限信息
	GetPermissionsWithNames(ctx int64, names ...string) (result []*Permission, err error)

	// GetPermissionsWithRoleId 获取已授予给指定角色的权限列表
	GetPermissionsWithRoleId(ctx int64, roleId int64) (result []*Permission, err error)

	// GetPermissionWithId 获取权限信息
	GetPermissionWithId(ctx, permissionId int64) (result *Permission, err error)

	// GetPermissionWithName 获取权限信息
	GetPermissionWithName(ctx int64, name string) (result *Permission, err error)

	// AddPermission 添加权限信息
	AddPermission(ctx, groupId int64, name, aliasName, description string, status Status) (result int64, err error)

	// UpdatePermission 更新权限信息
	UpdatePermission(ctx, permissionId, groupId int64, aliasName, description string, status Status) (err error)

	// UpdatePermissionStatus 更新权限状态
	UpdatePermissionStatus(ctx, permissionId int64, status Status) (err error)

	// GrantPermissionWithIds 授予权限给角色
	GrantPermissionWithIds(ctx, roleId int64, permissionIds []int64) (err error)

	// RevokePermissionWithIds 取消对角色的权限授权
	RevokePermissionWithIds(ctx, roleId int64, permissionIds []int64) (err error)

	// RevokeAllPermission 权限对角色的所有权限授权
	RevokeAllPermission(ctx, roleId int64) (err error)

	// GetGrantedPermissions 获取指定 target 拥有的权限信息
	GetGrantedPermissions(ctx int64, target string) (result []*Permission, err error)

	// AddPrePermission 添加权限先决条件
	AddPrePermission(ctx, permissionId int64, prePermissionIds []int64) (err error)

	// RemovePrePermission 移除权限先决条件
	RemovePrePermission(ctx, permissionId int64, prePermissionIds []int64) (err error)

	// CleanPrePermission 清除权限先决条件，即移除该权限的所有先决条件
	CleanPrePermission(ctx, permissionId int64) (err error)

	// GetPrePermissions 获取指定权限的所有先决条件
	GetPrePermissions(ctx, permissionId int64) (result []*PrePermission, err error)

	// GetPrePermissionsWithIds 获取指定权限列表的所有先决条件
	GetPrePermissionsWithIds(ctx int64, permissionIds []int64) (result []*PrePermission, err error)

	// GetRoles 获取角色列表
	// 如果参数 parentId 的值大于等于 0，则表示查询 parentId 的子角色列表
	// 如果参数 isGrantedToTarget 的值不为空字符串，则返回的角色数据中将包含该角色（通过 Granted 判断）是否已授权给 isGrantedToTarget
	GetRoles(ctx int64, parentId int64, status Status, keywords, isGrantedToTarget string) (result []*Role, err error)

	// GetRolesInTarget 获取角色列表
	// 如果参数 limitedInTarget 的值不为空字符串， 则返回的角色数据将限定在 limitedInTarget 已拥有的角色及其子角色范围内
	// 如果参数 isGrantedToTarget 的值不为空字符串，则返回的角色数据中将包含该角色（通过 Granted 判断）是否已授权给 isGrantedToTarget
	GetRolesInTarget(ctx int64, limitedInTarget string, status Status, keywords, isGrantedToTarget string) (result []*Role, err error)

	// GetRolesWithIds 根据角色 id 列表获取角色列表信息
	GetRolesWithIds(ctx int64, roleIds ...int64) (result []*Role, err error)

	// GetRolesWithNames 根据角色名称列表获取角色列表信息
	GetRolesWithNames(ctx int64, names ...string) (result []*Role, err error)

	// GetRoleWithId 获取角色信息
	GetRoleWithId(ctx, roleId int64) (result *Role, err error)

	// GetRoleWithName 获取角色信息
	GetRoleWithName(ctx int64, name string) (result *Role, err error)

	// AddRole 添加角色
	AddRole(ctx int64, parent *Role, name, aliasName, description string, status Status) (result int64, err error)

	// UpdateRole 更新角色
	UpdateRole(ctx, roleId int64, aliasName, description string, status Status) (err error)

	// UpdateRoleStatus 更新角色状态
	UpdateRoleStatus(ctx, roleId int64, status Status) (err error)

	// AddRoleMutex 添加角色互斥关系
	AddRoleMutex(ctx, roleId int64, mutexRoleIds []int64) (err error)

	// RemoveRoleMutex 移除角色互斥关系
	RemoveRoleMutex(ctx, roleId int64, mutexRoleIds []int64) (err error)

	// CleanRoleMutex 清除角色互斥关系，即移除该角色的所有互斥角色
	CleanRoleMutex(ctx, roleId int64) (err error)

	// GetMutexRoles 获取指定角色的所有互斥角色
	GetMutexRoles(ctx, roleId int64) (result []*RoleMutex, err error)

	// GetMutexRolesWithIds 获取指定角色列表的所有互斥角色
	GetMutexRolesWithIds(ctx int64, roleIds []int64) (result []*RoleMutex, err error)

	// CheckRoleMutex 验证角色是否是互斥关系
	CheckRoleMutex(ctx, roleId, mutexRoleId int64) bool

	// AddPreRole 添加角色先决条件
	AddPreRole(ctx, roleId int64, preRoleIds []int64) (err error)

	// RemovePreRole 移除角色先决条件
	RemovePreRole(ctx, roleId int64, preRoleIds []int64) (err error)

	// CleanPreRole 清除角色先决条件，即移除该角色的所有先决条件
	CleanPreRole(ctx, roleId int64) (err error)

	// GetPreRoles 获取角色的先决条件
	GetPreRoles(ctx, roleId int64) (result []*PreRole, err error)

	// GetPreRolesWithIds 获取指定角色列表的所有先决条件
	GetPreRolesWithIds(ctx int64, roleIds []int64) (result []*PreRole, err error)

	// GetGrantedRoles 获取已授权给 target 的角色列表
	// 如果参数 withChildren 的值为 true，则返回的角色数据中将包含该角色的子角色列表（子角色列表不一定授权给 target）
	GetGrantedRoles(ctx int64, target string, withChildren bool) (result []*Role, err error)

	// GrantRoleWithIds 授予角色给 target
	GrantRoleWithIds(ctx int64, target string, roleIds ...int64) (err error)

	// RevokeRoleWithIds 取消对 target 的角色授权
	RevokeRoleWithIds(ctx int64, target string, roleIds ...int64) (err error)

	// RevokeAllRole 取消对 target 的所有角色授权
	RevokeAllRole(ctx int64, target string) (err error)

	// CheckPermission 验证 target 是否拥有指定权限
	CheckPermission(ctx int64, target string, permissionName string) bool

	// CheckPermissionWithId 验证 target 是否拥有指定权限
	CheckPermissionWithId(ctx int64, target string, permissionId int64) bool

	// CheckRole 验证 target 是否拥有指定角色
	CheckRole(ctx int64, target string, roleName string) bool

	// CheckRoleWithId 验证 target 是否拥有指定角色
	CheckRoleWithId(ctx int64, target string, roleId int64) bool

	// CheckRoleAccessible 验证 target 是否拥有指定角色的操作权限
	CheckRoleAccessible(ctx int64, target string, roleName string) bool

	// CheckRoleAccessibleWithId 验证 target 是否拥有指定角色的操作权限
	CheckRoleAccessibleWithId(ctx int64, target string, roleId int64) bool

	// CheckRolePermission 验证角色是否拥有指定权限
	CheckRolePermission(ctx int64, roleName, permissionName string) bool

	// CheckRolePermissionWithId 验证角色是否拥有指定权限
	CheckRolePermissionWithId(ctx, roleId, permissionId int64) bool

	// CleanCache 清除缓存
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

func (this *odinService) Init() error {
	var tx, nRepo = this.repo.BeginTx()
	if err := nRepo.InitTable(); err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (this *odinService) GetPermissionGroups(ctx int64, status Status, keywords string) (result []*Group, err error) {
	return this.repo.GetGroups(ctx, GroupPermission, status, keywords)
}

func (this *odinService) GetPermissionGroup(ctx int64, groupName string) (result *Group, err error) {
	return this.repo.GetGroupWithName(ctx, GroupPermission, groupName)
}

func (this *odinService) GetPermissionGroupWithId(ctx, groupId int64) (result *Group, err error) {
	return this.repo.GetGroupWithId(ctx, GroupPermission, groupId)
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

func (this *odinService) UpdatePermissionGroup(ctx int64, groupName string, aliasName string, status Status) (err error) {
	return this.updateGroup(ctx, GroupPermission, groupName, aliasName, status)
}

func (this *odinService) UpdatePermissionGroupWithId(ctx, groupId int64, aliasName string, status Status) (err error) {
	return this.updateGroupWithId(ctx, GroupPermission, groupId, aliasName, status)
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

func (this *odinService) UpdatePermissionGroupStatusWithId(ctx int64, groupId int64, status Status) (err error) {
	return this.updateGroupStatusWithId(ctx, GroupPermission, groupId, status)
}

// 权限

func (this *odinService) GetPermissions(ctx int64, status Status, keywords string, groupIds []int64) (result []*Permission, err error) {
	return this.repo.GetPermissions(ctx, status, keywords, groupIds, 0, 0)
}

func (this *odinService) GetPermission(ctx int64, permissionName string) (result *Permission, err error) {
	result, err = this.repo.GetPermissionWithName(ctx, permissionName)
	if err != nil {
		return nil, err
	}
	if result != nil {
		if result.PrePermissionList, err = this.repo.GetPrePermissions(ctx, result.Id); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (this *odinService) GetPermissionWithId(ctx, permissionId int64) (result *Permission, err error) {
	result, err = this.repo.GetPermissionWithId(ctx, permissionId)
	if err != nil {
		return nil, err
	}
	if result != nil {
		if result.PrePermissionList, err = this.repo.GetPrePermissions(ctx, result.Id); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (this *odinService) AddPermissionWithGroup(ctx int64, groupName, permissionName, aliasName, description string, status Status) (result int64, err error) {
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

	// 获取当前角色的父角色
	if role.ParentId > 0 {
		parent, err := nRepo.GetRoleWithId(ctx, role.ParentId)
		if err != nil {
			return err
		}
		if parent == nil || parent.Status != Enable {
			return ErrInvalidParentRole
		}

		// 验证将要授权的权限是否超出父角色的权限
		parentPermissions, err := nRepo.GetPermissionsWithRoleId(ctx, parent.Id)
		if err != nil {
			return err
		}

		var permissionMap = make(map[string]struct{})
		for _, p := range parentPermissions {
			permissionMap[p.Name] = struct{}{}
		}

		for _, pName := range permissionNames {
			if _, ok := permissionMap[pName]; ok == false {
				return ErrPermissionOutOfParent
			}
		}
	}

	permissionList, err := nRepo.GetPermissionsWithNames(ctx, permissionNames...)
	if err != nil {
		return err
	}

	var gIds = make([]int64, 0, len(permissionList)) // 本次新授权权限 id 列表加上原来已授权权限 id 列表
	var nIds = make([]int64, 0, len(permissionList)) // 本次新授权权限 id 列表
	var gIdm = make(map[int64]struct{})              // 本次新授权权限 id 加上原来已授权权限 id 组成的 map
	for _, permission := range permissionList {
		gIds = append(gIds, permission.Id)
		nIds = append(nIds, permission.Id)
		gIdm[permission.Id] = struct{}{}
	}
	if len(nIds) == 0 {
		return ErrGrantFailed
	}

	// 查询出已授予给该角色的权限
	grantedPermissionList, err := nRepo.GetPermissionsWithRoleId(ctx, role.Id)
	if err != nil {
		return err
	}
	for _, permission := range grantedPermissionList {
		gIds = append(gIds, permission.Id)
		gIdm[permission.Id] = struct{}{}
	}

	// 获取并验证所有权限所需要的权限先决条件
	prePermissionList, err := nRepo.GetPrePermissionsWithIds(ctx, gIds)
	if err != nil {
		return err
	}
	for _, pPermission := range prePermissionList {
		if _, ok := gIdm[pPermission.PrePermissionId]; ok == false {
			return fmt.Errorf("授予权限 %s 时需要先授予权限 %s", pPermission.PermissionAliasName, pPermission.PrePermissionAliasName)
		}
	}

	if err = nRepo.GrantPermissionWithIds(ctx, role.Id, nIds); err != nil {
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

	// 获取当前角色的父角色
	if role.ParentId > 0 {
		parent, err := nRepo.GetRoleWithId(ctx, role.ParentId)
		if err != nil {
			return err
		}
		if parent == nil || parent.Status != Enable {
			return ErrInvalidParentRole
		}

		// 验证将要授权的权限是否超出父角色的权限
		parentPermissions, err := nRepo.GetPermissionsWithRoleId(ctx, parent.Id)
		if err != nil {
			return err
		}

		var permissionMap = make(map[int64]struct{})
		for _, p := range parentPermissions {
			permissionMap[p.Id] = struct{}{}
		}

		for _, pId := range permissionIds {
			if _, ok := permissionMap[pId]; ok == false {
				return ErrPermissionOutOfParent
			}
		}
	}

	permissionList, err := nRepo.GetPermissionsWithIds(ctx, permissionIds...)
	if err != nil {
		return err
	}

	var gIds = make([]int64, 0, len(permissionList))
	var nIds = make([]int64, 0, len(permissionList))
	var gIdm = make(map[int64]struct{})
	for _, permission := range permissionList {
		gIds = append(gIds, permission.Id)
		nIds = append(nIds, permission.Id)
		gIdm[permission.Id] = struct{}{}
	}
	if len(nIds) == 0 {
		return ErrGrantFailed
	}

	// 查询出已授予给该角色的权限
	grantedPermissionList, err := nRepo.GetPermissionsWithRoleId(ctx, role.Id)
	if err != nil {
		return err
	}
	for _, permission := range grantedPermissionList {
		gIds = append(gIds, permission.Id)
		gIdm[permission.Id] = struct{}{}
	}

	// 获取并验证所有权限所需要的权限先决条件
	prePermissionList, err := nRepo.GetPrePermissionsWithIds(ctx, gIds)
	if err != nil {
		return err
	}
	for _, pPermission := range prePermissionList {
		if _, ok := gIdm[pPermission.PrePermissionId]; ok == false {
			return fmt.Errorf("授予权限 %s 时需要先授予权限 %s", pPermission.PermissionAliasName, pPermission.PrePermissionAliasName)
		}
	}

	if err = nRepo.GrantPermissionWithIds(ctx, roleId, nIds); err != nil {
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

	// 获取当前角色的父角色
	if role.ParentId > 0 {
		parent, err := nRepo.GetRoleWithId(ctx, role.ParentId)
		if err != nil {
			return err
		}
		if parent == nil || parent.Status != Enable {
			return ErrInvalidParentRole
		}

		// 验证将要授权的权限是否超出父角色的权限
		parentPermissions, err := nRepo.GetPermissionsWithRoleId(ctx, parent.Id)
		if err != nil {
			return err
		}

		var permissionMap = make(map[string]struct{})
		for _, p := range parentPermissions {
			permissionMap[p.Name] = struct{}{}
		}

		for _, pName := range permissionNames {
			if _, ok := permissionMap[pName]; ok == false {
				return ErrPermissionOutOfParent
			}
		}
	}

	permissionList, err := nRepo.GetPermissionsWithNames(ctx, permissionNames...)
	if err != nil {
		return err
	}

	var nIds = make([]int64, 0, len(permissionList))
	var gIdm = make(map[int64]struct{})
	for _, permission := range permissionList {
		nIds = append(nIds, permission.Id)
		gIdm[permission.Id] = struct{}{}
	}
	if len(nIds) == 0 {
		return ErrGrantFailed
	}

	// 查出原有的权限
	rolePermissions, err := nRepo.GetPermissionsWithRoleId(ctx, role.Id)
	if err != nil {
		return err
	}

	// 查出需要取消掉权限
	var revokeIds = make([]int64, 0, len(rolePermissions))
	for _, p := range rolePermissions {
		if _, ok := gIdm[p.Id]; ok == false {
			revokeIds = append(revokeIds, p.Id)
		}
	}

	// 获取并验证所有权限所需要的权限先决条件
	prePermissionList, err := nRepo.GetPrePermissionsWithIds(ctx, nIds)
	if err != nil {
		return err
	}
	for _, pPermission := range prePermissionList {
		if _, ok := gIdm[pPermission.PrePermissionId]; ok == false {
			return fmt.Errorf("授予权限 %s 时需要先授予权限 %s", pPermission.PermissionAliasName, pPermission.PrePermissionAliasName)
		}
	}

	if len(revokeIds) > 0 {
		if err = nRepo.RevokePermissionWithIds(ctx, role.Id, revokeIds); err != nil {
			return err
		}
	}

	if err = nRepo.GrantPermissionWithIds(ctx, role.Id, nIds); err != nil {
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

	// 获取当前角色的父角色
	if role.ParentId > 0 {
		parent, err := nRepo.GetRoleWithId(ctx, role.ParentId)
		if err != nil {
			return err
		}
		if parent == nil || parent.Status != Enable {
			return ErrInvalidParentRole
		}

		// 验证将要授权的权限是否超出父角色的权限
		parentPermissions, err := nRepo.GetPermissionsWithRoleId(ctx, parent.Id)
		if err != nil {
			return err
		}

		var permissionMap = make(map[int64]struct{})
		for _, p := range parentPermissions {
			permissionMap[p.Id] = struct{}{}
		}

		for _, pId := range permissionIds {
			if _, ok := permissionMap[pId]; ok == false {
				return ErrPermissionOutOfParent
			}
		}
	}

	permissions, err := nRepo.GetPermissionsWithIds(ctx, permissionIds...)
	if err != nil {
		return err
	}

	var nIds = make([]int64, 0, len(permissions))
	var gIdm = make(map[int64]struct{})
	for _, permission := range permissions {
		nIds = append(nIds, permission.Id)
		gIdm[permission.Id] = struct{}{}
	}
	if len(nIds) == 0 {
		return ErrGrantFailed
	}

	// 查出原有的权限
	rolePermissions, err := nRepo.GetPermissionsWithRoleId(ctx, roleId)
	if err != nil {
		return err
	}

	// 查出需要取消掉权限
	var revokeIds = make([]int64, 0, len(rolePermissions))
	for _, p := range rolePermissions {
		if _, ok := gIdm[p.Id]; ok == false {
			revokeIds = append(revokeIds, p.Id)
		}
	}

	// 获取并验证所有权限所需要的权限先决条件
	prePermissionList, err := nRepo.GetPrePermissionsWithIds(ctx, nIds)
	if err != nil {
		return err
	}
	for _, pPermission := range prePermissionList {
		if _, ok := gIdm[pPermission.PrePermissionId]; ok == false {
			return fmt.Errorf("授予权限 %s 时需要先授予权限 %s", pPermission.PermissionAliasName, pPermission.PrePermissionAliasName)
		}
	}

	if len(revokeIds) > 0 {
		if err = nRepo.RevokePermissionWithIds(ctx, role.Id, revokeIds); err != nil {
			return err
		}
	}

	if err = nRepo.GrantPermissionWithIds(ctx, role.Id, nIds); err != nil {
		return err
	}

	tx.Commit()
	return nil
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

	rPermissionList, err := nRepo.GetPermissionsWithNames(ctx, permissionNames...)
	if err != nil {
		return err
	}

	var rIds = make([]int64, 0, len(rPermissionList))
	for _, permission := range rPermissionList {
		rIds = append(rIds, permission.Id)
	}
	if len(rIds) == 0 {
		return ErrRevokeFailed
	}

	if err = nRepo.RevokePermissionWithIds(ctx, role.Id, rIds); err != nil {
		return err
	}

	// 查询出已授予给角色的权限
	grantedPermissionList, err := nRepo.GetPermissionsWithRoleId(ctx, role.Id)
	if err != nil {
		return err
	}
	var gIds = make([]int64, 0, len(grantedPermissionList))
	var gIdm = make(map[int64]struct{})

	for _, permission := range grantedPermissionList {
		gIds = append(gIds, permission.Id)
		gIdm[permission.Id] = struct{}{}
	}

	// 获取并验证所有权限所需要的权限先决条件
	if len(gIds) > 0 {
		prePermissionList, err := nRepo.GetPrePermissionsWithIds(ctx, gIds)
		if err != nil {
			return err
		}
		for _, pPermission := range prePermissionList {
			if _, ok := gIdm[pPermission.PrePermissionId]; ok == false {
				return fmt.Errorf("权限 %s 依赖于权限 %s", pPermission.PermissionAliasName, pPermission.PrePermissionAliasName)
			}
		}
	}

	tx.Commit()
	return nil
}

func (this *odinService) RevokePermissionWithId(ctx int64, roleId int64, permissionIds ...int64) (err error) {
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

	rPermissionList, err := nRepo.GetPermissionsWithIds(ctx, permissionIds...)
	if err != nil {
		return err
	}

	var rIds = make([]int64, 0, len(rPermissionList))
	for _, permission := range rPermissionList {
		rIds = append(rIds, permission.Id)
	}
	if len(rIds) == 0 {
		return ErrRevokeFailed
	}

	if err = nRepo.RevokePermissionWithIds(ctx, role.Id, rIds); err != nil {
		return err
	}

	// 查询出已授予给角色的权限
	grantedPermissionList, err := nRepo.GetPermissionsWithRoleId(ctx, role.Id)
	if err != nil {
		return err
	}
	var gIds = make([]int64, 0, len(grantedPermissionList))
	var gIdm = make(map[int64]struct{})

	for _, permission := range grantedPermissionList {
		gIds = append(gIds, permission.Id)
		gIdm[permission.Id] = struct{}{}
	}

	// 获取并验证所有权限所需要的权限先决条件
	if len(gIds) > 0 {
		prePermissionList, err := nRepo.GetPrePermissionsWithIds(ctx, gIds)
		if err != nil {
			return err
		}
		for _, pPermission := range prePermissionList {
			if _, ok := gIdm[pPermission.PrePermissionId]; ok == false {
				return fmt.Errorf("权限 %s 依赖于权限 %s", pPermission.PermissionAliasName, pPermission.PrePermissionAliasName)
			}
		}
	}

	tx.Commit()
	return nil
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

func (this *odinService) RevokeAllPermissionWithId(ctx, roleId int64) (err error) {
	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err = nRepo.RevokeAllPermission(ctx, roleId); err != nil {
		return err
	}
	tx.Commit()
	return nil
}

// 权限先决条件

func (this *odinService) AddPrePermission(ctx int64, permissionName string, prePermissionNames ...string) (err error) {
	if len(prePermissionNames) == 0 {
		return ErrPrePermissionNotExist
	}

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

	prePermissionList, err := nRepo.GetPermissionsWithNames(ctx, prePermissionNames...)
	if err != nil {
		return err
	}

	var preIds = make([]int64, 0, len(prePermissionList))
	for _, permission := range prePermissionList {
		preIds = append(preIds, permission.Id)
	}
	if len(preIds) == 0 {
		return ErrPreRoleNotExist
	}

	if err = nRepo.AddPrePermission(ctx, permission.Id, preIds); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) AddPrePermissionWithId(ctx, permissionId int64, prePermissionIds ...int64) (err error) {
	if len(prePermissionIds) == 0 {
		return ErrPrePermissionNotExist
	}

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

	prePermissionList, err := nRepo.GetPermissionsWithIds(ctx, prePermissionIds...)
	if err != nil {
		return err
	}

	var preIds = make([]int64, 0, len(prePermissionList))
	for _, permission := range prePermissionList {
		preIds = append(preIds, permission.Id)
	}
	if len(preIds) == 0 {
		return ErrPreRoleNotExist
	}

	if err = nRepo.AddPrePermission(ctx, permission.Id, preIds); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) RemovePrePermission(ctx int64, permissionName string, prePermissionNames ...string) (err error) {
	if len(prePermissionNames) == 0 {
		return ErrPrePermissionNotExist
	}

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

	prePermissionList, err := nRepo.GetPermissionsWithNames(ctx, prePermissionNames...)
	if err != nil {
		return err
	}

	var preIds = make([]int64, 0, len(prePermissionList))
	for _, permission := range prePermissionList {
		preIds = append(preIds, permission.Id)
	}
	if len(preIds) == 0 {
		return ErrPrePermissionNotExist
	}

	if err = nRepo.RemovePrePermission(ctx, permission.Id, preIds); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) RemovePrePermissionWithId(ctx, permissionId int64, prePermissionIds ...int64) (err error) {
	if len(prePermissionIds) == 0 {
		return ErrPrePermissionNotExist
	}

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

	prePermissionList, err := nRepo.GetPermissionsWithIds(ctx, prePermissionIds...)
	if err != nil {
		return err
	}

	var preIds = make([]int64, 0, len(prePermissionList))
	for _, permission := range prePermissionList {
		preIds = append(preIds, permission.Id)
	}
	if len(preIds) == 0 {
		return ErrPrePermissionNotExist
	}

	if err = nRepo.RemovePrePermission(ctx, permission.Id, preIds); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) RemoveAllPrePermission(ctx int64, permissionName string) (err error) {
	// 验证权限是否存在
	permission, err := this.repo.GetPermissionWithName(ctx, permissionName)
	if err != nil {
		return err
	}
	if permission == nil {
		return ErrPermissionNotExist
	}
	return this.repo.CleanPrePermission(ctx, permission.Id)
}

func (this *odinService) RemoveAllPrePermissionWithId(ctx, permissionId int64) (err error) {
	// 验证权限是否存在
	permission, err := this.repo.GetPermissionWithId(ctx, permissionId)
	if err != nil {
		return err
	}
	if permission == nil {
		return ErrPermissionNotExist
	}
	return this.repo.CleanPrePermission(ctx, permission.Id)
}

func (this *odinService) GetPrePermissions(ctx int64, permissionName string) (result []*PrePermission, err error) {
	// 验证权限是否存在
	permission, err := this.repo.GetPermissionWithName(ctx, permissionName)
	if err != nil {
		return nil, err
	}
	if permission == nil {
		return nil, ErrRoleNotExist
	}
	return this.repo.GetPrePermissions(ctx, permission.Id)
}

func (this *odinService) GetPrePermissionsWithId(ctx int64, permissionId int64) (result []*PrePermission, err error) {
	// 验证权限是否存在
	permission, err := this.repo.GetPermissionWithId(ctx, permissionId)
	if err != nil {
		return nil, err
	}
	if permission == nil {
		return nil, ErrRoleNotExist
	}
	return this.repo.GetPrePermissions(ctx, permission.Id)
}

// 角色

func (this *odinService) GetRoles(ctx int64, status Status, keywords, isGrantedToTarget, limitedInTarget string) (result []*Role, err error) {
	if limitedInTarget == "" {
		return this.repo.GetRoles(ctx, -1, status, keywords, isGrantedToTarget)
	}
	return this.repo.GetRolesInTarget(ctx, limitedInTarget, status, keywords, isGrantedToTarget)
}

func (this *odinService) GetRolesWithParent(ctx int64, parentRoleName string, status Status, keywords, isGrantedToTarget string) (result []*Role, err error) {
	var parentRoleId int64 = 0
	if parentRoleName != "" {
		// 验证 parentRoleName 是否存在
		role, err := this.repo.GetRoleWithName(ctx, parentRoleName)
		if err != nil {
			return nil, err
		}
		if role == nil {
			return nil, ErrRoleNotExist
		}
		parentRoleId = role.Id
	}
	return this.repo.GetRoles(ctx, parentRoleId, status, keywords, isGrantedToTarget)
}

func (this *odinService) GetRolesWithParentId(ctx, parentRoleId int64, status Status, keywords, isGrantedToTarget string) (result []*Role, err error) {
	if parentRoleId < 0 {
		parentRoleId = 0
	}
	if parentRoleId > 0 {
		// 验证 parentRoleId 是否存在
		role, err := this.repo.GetRoleWithId(ctx, parentRoleId)
		if err != nil {
			return nil, err
		}
		if role == nil {
			return nil, ErrRoleNotExist
		}
		parentRoleId = role.Id
	}
	return this.repo.GetRoles(ctx, parentRoleId, status, keywords, isGrantedToTarget)
}

func (this *odinService) GetRole(ctx int64, name string) (result *Role, err error) {
	result, err = this.repo.GetRoleWithName(ctx, name)
	if err != nil {
		return nil, err
	}
	if result != nil {
		if result.PermissionList, err = this.repo.GetPermissionsWithRoleId(ctx, result.Id); err != nil {
			return nil, err
		}
		if result.MutexRoleList, err = this.repo.GetMutexRoles(ctx, result.Id); err != nil {
			return nil, err
		}
		if result.PreRoleList, err = this.repo.GetPreRoles(ctx, result.Id); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (this *odinService) GetRoleWithId(ctx, roleId int64) (result *Role, err error) {
	result, err = this.repo.GetRoleWithId(ctx, roleId)
	if err != nil {
		return nil, err
	}
	if result != nil {
		if result.PermissionList, err = this.repo.GetPermissionsWithRoleId(ctx, result.Id); err != nil {
			return nil, err
		}
		if result.MutexRoleList, err = this.repo.GetMutexRoles(ctx, result.Id); err != nil {
			return nil, err
		}
		if result.PreRoleList, err = this.repo.GetPreRoles(ctx, result.Id); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (this *odinService) AddRole(ctx int64, roleName, aliasName, description string, status Status) (result int64, err error) {
	return this.AddRoleWithParentId(ctx, 0, roleName, aliasName, description, status)
}

func (this *odinService) AddRoleWithParent(ctx int64, parentRoleName, roleName, aliasName, description string, status Status) (result int64, err error) {
	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var parentRole *Role
	if parentRoleName != "" {
		// 验证 parentRoleName 是否存在
		parentRole, err = nRepo.GetRoleWithName(ctx, parentRoleName)
		if err != nil {
			return 0, err
		}
		if parentRole == nil {
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

	if result, err = nRepo.AddRole(ctx, parentRole, roleName, aliasName, description, status); err != nil {
		return 0, err
	}

	tx.Commit()
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
		return 0, ErrParentRoleNotExist
	}

	var parentRole *Role
	if parentRoleId > 0 {
		// 验证 parentRoleName 是否存在
		parentRole, err = nRepo.GetRoleWithId(ctx, parentRoleId)
		if err != nil {
			return 0, err
		}
		if parentRole == nil {
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

	if result, err = nRepo.AddRole(ctx, parentRole, roleName, aliasName, description, status); err != nil {
		return 0, err
	}

	tx.Commit()
	return result, nil
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

func (this *odinService) GrantRole(ctx int64, target string, roleNames ...string) (err error) {
	if len(roleNames) == 0 {
		return ErrRoleNotExist
	}

	if target == "" {
		return ErrTargetNotAllowed
	}

	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	roleList, err := nRepo.GetRolesWithNames(ctx, roleNames...)
	if err != nil {
		return err
	}

	var gIds = make([]int64, 0, len(roleList))
	var nIds = make([]int64, 0, len(roleList))
	var gIdm = make(map[int64]struct{})
	for _, role := range roleList {
		gIds = append(gIds, role.Id)
		nIds = append(nIds, role.Id)
		gIdm[role.Id] = struct{}{}
	}
	if len(nIds) == 0 {
		return ErrGrantFailed
	}

	// 查询出已授予给 target 的角色
	grantedRoleList, err := nRepo.GetGrantedRoles(ctx, target, false)
	if err != nil {
		return err
	}
	for _, role := range grantedRoleList {
		gIds = append(gIds, role.Id)
		gIdm[role.Id] = struct{}{}
	}

	// 获取并验证互斥关系
	mutexRoleList, err := nRepo.GetMutexRolesWithIds(ctx, gIds)
	if err != nil {
		return err
	}
	for _, role := range mutexRoleList {
		return fmt.Errorf("角色 %s 与角色 %s 互斥", role.RoleAliasName, role.MutexRoleAliasName)
	}

	// 获取并验证所有角色所需要的角色先决条件
	preRoleList, err := nRepo.GetPreRolesWithIds(ctx, gIds)
	if err != nil {
		return err
	}
	for _, pRole := range preRoleList {
		if _, ok := gIdm[pRole.PreRoleId]; ok == false {
			return fmt.Errorf("授予角色 %s 时需要先授予角色 %s", pRole.RoleAliasName, pRole.PreRoleAliasName)
		}
	}

	if err = nRepo.GrantRoleWithIds(ctx, target, nIds...); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) GrantRoleWithId(ctx int64, target string, roleIds ...int64) (err error) {
	if len(roleIds) == 0 {
		return ErrRoleNotExist
	}

	if target == "" {
		return ErrTargetNotAllowed
	}

	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	roleList, err := nRepo.GetRolesWithIds(ctx, roleIds...)
	if err != nil {
		return err
	}

	var gIds = make([]int64, 0, len(roleList)) // 本次新授权角色 id 列表加上原来已授权角色 id 列表
	var nIds = make([]int64, 0, len(roleList)) // 本次新授权角色 id 列表
	var gIdm = make(map[int64]struct{})        // 本次新授权角色 id 加上原来已授权角色 id 组成的 map
	for _, role := range roleList {
		gIds = append(gIds, role.Id)
		nIds = append(nIds, role.Id)
		gIdm[role.Id] = struct{}{}
	}
	if len(nIds) == 0 {
		return ErrGrantFailed
	}

	// 查询出已授予给 target 的角色
	grantedRoleList, err := nRepo.GetGrantedRoles(ctx, target, false)
	if err != nil {
		return err
	}
	for _, role := range grantedRoleList {
		gIds = append(gIds, role.Id)
		gIdm[role.Id] = struct{}{}
	}

	// 获取并验证互斥关系
	mutexRoleList, err := nRepo.GetMutexRolesWithIds(ctx, gIds)
	if err != nil {
		return err
	}
	for _, role := range mutexRoleList {
		return fmt.Errorf("角色 %s 与角色 %s 互斥", role.RoleAliasName, role.MutexRoleAliasName)
	}

	// 获取并验证所有角色所需要的角色先决条件
	preRoleList, err := nRepo.GetPreRolesWithIds(ctx, gIds)
	if err != nil {
		return err
	}
	for _, pRole := range preRoleList {
		if _, ok := gIdm[pRole.PreRoleId]; ok == false {
			return fmt.Errorf("授予角色 %s 时需要先授予角色 %s", pRole.RoleAliasName, pRole.PreRoleAliasName)
		}
	}

	if err = nRepo.GrantRoleWithIds(ctx, target, nIds...); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) ReGrantRole(ctx int64, target string, roleNames ...string) (err error) {
	if len(roleNames) == 0 {
		return ErrRoleNotExist
	}

	if target == "" {
		return ErrTargetNotAllowed
	}

	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	roleList, err := nRepo.GetRolesWithNames(ctx, roleNames...)
	if err != nil {
		return err
	}

	var nIds = make([]int64, 0, len(roleList))
	var gIdm = make(map[int64]struct{})
	for _, role := range roleList {
		nIds = append(nIds, role.Id)
		gIdm[role.Id] = struct{}{}
	}
	if len(nIds) == 0 {
		return ErrGrantFailed
	}

	// 获取并验证互斥关系
	mutexRoleList, err := nRepo.GetMutexRolesWithIds(ctx, nIds)
	if err != nil {
		return err
	}
	for _, role := range mutexRoleList {
		return fmt.Errorf("角色 %s 与角色 %s 互斥", role.RoleAliasName, role.MutexRoleAliasName)
	}

	// 获取并验证所有角色所需要的角色先决条件
	preRoleList, err := nRepo.GetPreRolesWithIds(ctx, nIds)
	if err != nil {
		return err
	}
	for _, pRole := range preRoleList {
		if _, ok := gIdm[pRole.PreRoleId]; ok == false {
			return fmt.Errorf("授予角色 %s 时需要先授予角色 %s", pRole.RoleAliasName, pRole.PreRoleAliasName)
		}
	}

	if err = nRepo.RevokeAllRole(ctx, target); err != nil {
		return err
	}

	if err = nRepo.GrantRoleWithIds(ctx, target, nIds...); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) ReGrantRoleWithId(ctx int64, target string, roleIds ...int64) (err error) {
	if len(roleIds) == 0 {
		return ErrRoleNotExist
	}

	if target == "" {
		return ErrTargetNotAllowed
	}

	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	roleList, err := nRepo.GetRolesWithIds(ctx, roleIds...)
	if err != nil {
		return err
	}

	var nIds = make([]int64, 0, len(roleList))
	var gIdm = make(map[int64]struct{})
	for _, role := range roleList {
		nIds = append(nIds, role.Id)
		gIdm[role.Id] = struct{}{}
	}
	if len(nIds) == 0 {
		return ErrGrantFailed
	}

	// 获取并验证互斥关系
	mutexRoleList, err := nRepo.GetMutexRolesWithIds(ctx, nIds)
	if err != nil {
		return err
	}
	for _, role := range mutexRoleList {
		return fmt.Errorf("角色 %s 与角色 %s 互斥", role.RoleAliasName, role.MutexRoleAliasName)
	}

	// 获取并验证所有角色所需要的角色先决条件
	preRoleList, err := nRepo.GetPreRolesWithIds(ctx, nIds)
	if err != nil {
		return err
	}
	for _, pRole := range preRoleList {
		if _, ok := gIdm[pRole.PreRoleId]; ok == false {
			return fmt.Errorf("授予角色 %s 时需要先授予角色 %s", pRole.RoleAliasName, pRole.PreRoleAliasName)
		}
	}

	if err = nRepo.RevokeAllRole(ctx, target); err != nil {
		return err
	}

	if err = nRepo.GrantRoleWithIds(ctx, target, nIds...); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) RevokeRole(ctx int64, target string, roleNames ...string) (err error) {
	if len(roleNames) == 0 {
		return ErrRoleNotExist
	}

	if target == "" {
		return ErrTargetNotAllowed
	}

	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	rRoleList, err := nRepo.GetRolesWithNames(ctx, roleNames...)
	if err != nil {
		return err
	}

	var rIds = make([]int64, 0, len(rRoleList))
	for _, role := range rRoleList {
		rIds = append(rIds, role.Id)
	}
	if len(rIds) == 0 {
		return ErrGrantFailed
	}

	if err = nRepo.RevokeRoleWithIds(ctx, target, rIds...); err != nil {
		return err
	}

	var gIds = make([]int64, 0, len(rRoleList))
	var gIdm = make(map[int64]struct{})

	// 查询出已授予给 target 的角色
	grantedRoleList, err := nRepo.GetGrantedRoles(ctx, target, false)
	if err != nil {
		return err
	}
	for _, role := range grantedRoleList {
		gIds = append(gIds, role.Id)
		gIdm[role.Id] = struct{}{}
	}

	// 获取并验证所有角色所需要的角色先决条件
	if len(gIds) > 0 {
		preRoleList, err := nRepo.GetPreRolesWithIds(ctx, gIds)
		if err != nil {
			return err
		}
		for _, pRole := range preRoleList {
			if _, ok := gIdm[pRole.PreRoleId]; ok == false {
				return fmt.Errorf("角色 %s 依赖于角色 %s", pRole.RoleAliasName, pRole.PreRoleAliasName)
			}
		}
	}

	tx.Commit()
	return nil
}

func (this *odinService) RevokeRoleWithId(ctx int64, target string, roleIds ...int64) (err error) {
	if len(roleIds) == 0 {
		return ErrRoleNotExist
	}

	if target == "" {
		return ErrTargetNotAllowed
	}

	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	rRoleList, err := nRepo.GetRolesWithIds(ctx, roleIds...)
	if err != nil {
		return err
	}

	var rIds = make([]int64, 0, len(rRoleList))
	for _, role := range rRoleList {
		rIds = append(rIds, role.Id)
	}
	if len(rIds) == 0 {
		return ErrGrantFailed
	}

	if err = nRepo.RevokeRoleWithIds(ctx, target, rIds...); err != nil {
		return err
	}

	var gIds = make([]int64, 0, len(rRoleList))
	var gIdm = make(map[int64]struct{})

	// 查询出已授予给 target 的角色
	grantedRoleList, err := nRepo.GetGrantedRoles(ctx, target, false)
	if err != nil {
		return err
	}
	for _, role := range grantedRoleList {
		gIds = append(gIds, role.Id)
		gIdm[role.Id] = struct{}{}
	}

	// 获取并验证所有角色所需要的角色先决条件
	if len(gIds) > 0 {
		preRoleList, err := nRepo.GetPreRolesWithIds(ctx, gIds)
		if err != nil {
			return err
		}
		for _, pRole := range preRoleList {
			if _, ok := gIdm[pRole.PreRoleId]; ok == false {
				return fmt.Errorf("角色 %s 依赖于角色 %s", pRole.RoleAliasName, pRole.PreRoleAliasName)
			}
		}
	}

	tx.Commit()
	return nil
}

func (this *odinService) RevokeAllRole(ctx int64, target string) (err error) {
	return this.repo.RevokeAllRole(ctx, target)
}

// 角色互斥

func (this *odinService) AddRoleMutex(ctx int64, roleName string, mutexRoleNames ...string) (err error) {
	if len(mutexRoleNames) == 0 {
		return ErrMutexRoleNotExist
	}

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

	mutexRoleList, err := nRepo.GetRolesWithNames(ctx, mutexRoleNames...)
	if err != nil {
		return err
	}

	var mutexIds = make([]int64, 0, len(mutexRoleList))
	for _, role := range mutexRoleList {
		mutexIds = append(mutexIds, role.Id)
	}
	if len(mutexIds) == 0 {
		return ErrMutexRoleNotExist
	}

	if err = nRepo.AddRoleMutex(ctx, role.Id, mutexIds); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) AddRoleMutexWithId(ctx, roleId int64, mutexRoleIds ...int64) (err error) {
	if len(mutexRoleIds) == 0 {
		return ErrMutexRoleNotExist
	}

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

	mutexRoleList, err := nRepo.GetRolesWithIds(ctx, mutexRoleIds...)
	if err != nil {
		return err
	}

	var mutexIds = make([]int64, 0, len(mutexRoleList))
	for _, role := range mutexRoleList {
		mutexIds = append(mutexIds, role.Id)
	}
	if len(mutexIds) == 0 {
		return ErrMutexRoleNotExist
	}

	if err = nRepo.AddRoleMutex(ctx, role.Id, mutexIds); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) RemoveRoleMutex(ctx int64, roleName string, mutexRoleNames ...string) (err error) {
	if len(mutexRoleNames) == 0 {
		return ErrMutexRoleNotExist
	}

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

	mutexRoleList, err := nRepo.GetRolesWithNames(ctx, mutexRoleNames...)
	if err != nil {
		return err
	}

	var mutexIds = make([]int64, 0, len(mutexRoleList))
	for _, role := range mutexRoleList {
		mutexIds = append(mutexIds, role.Id)
	}
	if len(mutexIds) == 0 {
		return ErrMutexRoleNotExist
	}

	if err = nRepo.RemoveRoleMutex(ctx, role.Id, mutexIds); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) RemoveRoleMutexWithId(ctx, roleId int64, mutexRoleIds ...int64) (err error) {
	if len(mutexRoleIds) == 0 {
		return ErrMutexRoleNotExist
	}

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

	mutexRoleList, err := nRepo.GetRolesWithIds(ctx, mutexRoleIds...)
	if err != nil {
		return err
	}

	var mutexIds = make([]int64, 0, len(mutexRoleList))
	for _, role := range mutexRoleList {
		mutexIds = append(mutexIds, role.Id)
	}
	if len(mutexIds) == 0 {
		return ErrMutexRoleNotExist
	}

	if err = nRepo.RemoveRoleMutex(ctx, role.Id, mutexIds); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) RemoveAllRoleMutex(ctx int64, roleName string) (err error) {
	// 验证角色是否存在
	role, err := this.repo.GetRoleWithName(ctx, roleName)
	if err != nil {
		return err
	}
	if role == nil {
		return ErrRoleNotExist
	}
	return this.repo.CleanRoleMutex(ctx, role.Id)
}

func (this *odinService) RemoveAllRoleMutexWithId(ctx, roleId int64) (err error) {
	// 验证角色是否存在
	role, err := this.repo.GetRoleWithId(ctx, roleId)
	if err != nil {
		return err
	}
	if role == nil {
		return ErrRoleNotExist
	}
	return this.repo.CleanRoleMutex(ctx, role.Id)
}

func (this *odinService) GetMutexRoles(ctx int64, roleName string) (result []*RoleMutex, err error) {
	// 验证角色是否存在
	role, err := this.repo.GetRoleWithName(ctx, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, ErrRoleNotExist
	}
	return this.repo.GetMutexRoles(ctx, role.Id)
}

func (this *odinService) GetMutexRolesWithId(ctx, roleId int64) (result []*RoleMutex, err error) {
	// 验证角色是否存在
	role, err := this.repo.GetRoleWithId(ctx, roleId)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, ErrRoleNotExist
	}
	return this.repo.GetMutexRoles(ctx, role.Id)
}

func (this *odinService) CheckRoleMutex(ctx int64, roleName, mutexRoleName string) bool {
	var tx, nRepo = this.repo.BeginTx()
	var err error
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 验证角色是否存在
	role, err := nRepo.GetRoleWithName(ctx, roleName)
	if err != nil {
		return true
	}
	if role == nil {
		return false
	}

	// 验证角色是否存在
	mutexRole, err := nRepo.GetRoleWithName(ctx, mutexRoleName)
	if err != nil {
		return true
	}
	if mutexRole == nil {
		return false
	}

	var ok = nRepo.CheckRoleMutex(ctx, role.Id, mutexRole.Id)
	tx.Commit()
	return ok
}

func (this *odinService) CheckRoleMutexWithId(ctx, roleId, mutexRoleId int64) bool {
	return this.repo.CheckRoleMutex(ctx, roleId, mutexRoleId)
}

// 角色先决条件

func (this *odinService) AddPreRole(ctx int64, roleName string, preRoleNames ...string) (err error) {
	if len(preRoleNames) == 0 {
		return ErrPreRoleNotExist
	}

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

	preRoleList, err := nRepo.GetRolesWithNames(ctx, preRoleNames...)
	if err != nil {
		return err
	}

	var preIds = make([]int64, 0, len(preRoleList))
	for _, role := range preRoleList {
		preIds = append(preIds, role.Id)
	}
	if len(preIds) == 0 {
		return ErrPreRoleNotExist
	}

	if err = nRepo.AddPreRole(ctx, role.Id, preIds); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) AddPreRoleWithId(ctx, roleId int64, preRoleIds ...int64) (err error) {
	if len(preRoleIds) == 0 {
		return ErrPreRoleNotExist
	}

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

	preRoleList, err := nRepo.GetRolesWithIds(ctx, preRoleIds...)
	if err != nil {
		return err
	}

	var preIds = make([]int64, 0, len(preRoleList))
	for _, role := range preRoleList {
		preIds = append(preIds, role.Id)
	}
	if len(preIds) == 0 {
		return ErrPreRoleNotExist
	}

	if err = nRepo.AddPreRole(ctx, role.Id, preIds); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) RemovePreRole(ctx int64, roleName string, preRoleNames ...string) (err error) {
	if len(preRoleNames) == 0 {
		return ErrPreRoleNotExist
	}

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

	preRoleList, err := nRepo.GetRolesWithNames(ctx, preRoleNames...)
	if err != nil {
		return err
	}

	var preIds = make([]int64, 0, len(preRoleList))
	for _, role := range preRoleList {
		preIds = append(preIds, role.Id)
	}
	if len(preIds) == 0 {
		return ErrPreRoleNotExist
	}

	if err = nRepo.RemovePreRole(ctx, role.Id, preIds); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) RemovePreRoleWithId(ctx, roleId int64, preRoleIds ...int64) (err error) {
	if len(preRoleIds) == 0 {
		return ErrPreRoleNotExist
	}

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

	preRoleList, err := nRepo.GetRolesWithIds(ctx, preRoleIds...)
	if err != nil {
		return err
	}

	var preIds = make([]int64, 0, len(preRoleList))
	for _, role := range preRoleList {
		preIds = append(preIds, role.Id)
	}
	if len(preIds) == 0 {
		return ErrPreRoleNotExist
	}

	if err = nRepo.RemovePreRole(ctx, role.Id, preIds); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) RemoveAllPreRole(ctx int64, roleName string) (err error) {
	// 验证角色是否存在
	role, err := this.repo.GetRoleWithName(ctx, roleName)
	if err != nil {
		return err
	}
	if role == nil {
		return ErrRoleNotExist
	}
	return this.repo.CleanPreRole(ctx, role.Id)
}

func (this *odinService) RemoveAllPreRoleWithId(ctx, roleId int64) (err error) {
	// 验证角色是否存在
	role, err := this.repo.GetRoleWithId(ctx, roleId)
	if err != nil {
		return err
	}
	if role == nil {
		return ErrRoleNotExist
	}
	return this.repo.CleanPreRole(ctx, role.Id)
}

func (this *odinService) GetPreRoles(ctx int64, roleName string) (result []*PreRole, err error) {
	// 验证角色是否存在
	role, err := this.repo.GetRoleWithName(ctx, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, ErrRoleNotExist
	}
	return this.repo.GetPreRoles(ctx, role.Id)
}

func (this *odinService) GetPreRolesWithId(ctx, roleId int64) (result []*PreRole, err error) {
	// 验证角色是否存在
	role, err := this.repo.GetRoleWithId(ctx, roleId)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, ErrRoleNotExist
	}
	return this.repo.GetPreRoles(ctx, roleId)
}

// 其它

func (this *odinService) GetGrantedRoles(ctx int64, target string) (result []*Role, err error) {
	return this.repo.GetGrantedRoles(ctx, target, false)
}

func (this *odinService) GetRolesWithTarget(ctx int64, target string) (result []*Role, err error) {
	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	result, err = nRepo.GetGrantedRoles(ctx, target, true)
	if err != nil {
		return nil, err
	}

	tx.Commit()
	return result, nil
}

func (this *odinService) CheckRole(ctx int64, target string, roleName string) bool {
	return this.repo.CheckRole(ctx, target, roleName)
}

func (this *odinService) CheckRoleWithId(ctx int64, target string, roleId int64) bool {
	return this.repo.CheckRoleWithId(ctx, target, roleId)
}

func (this *odinService) CheckRoleAccessible(ctx int64, target string, roleName string) bool {
	return this.repo.CheckRoleAccessible(ctx, target, roleName)
}

func (this *odinService) CheckRoleAccessibleWithId(ctx int64, target string, roleId int64) bool {
	return this.repo.CheckRoleAccessibleWithId(ctx, target, roleId)
}

func (this *odinService) GetPermissionsWithRole(ctx int64, roleName string) (result []*Permission, err error) {
	role, err := this.repo.GetRoleWithName(ctx, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, ErrRoleNotExist
	}
	return this.repo.GetPermissionsWithRoleId(ctx, role.Id)
}

func (this *odinService) GetPermissionsWithRoleId(ctx int64, roleId int64) (result []*Permission, err error) {
	role, err := this.repo.GetRoleWithId(ctx, roleId)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, ErrRoleNotExist
	}
	return this.repo.GetPermissionsWithRoleId(ctx, role.Id)
}

func (this *odinService) GetGrantedPermissions(ctx int64, target string) (result []*Permission, err error) {
	return this.repo.GetGrantedPermissions(ctx, target)
}

func (this *odinService) GetPermissionsTreeWithRole(ctx int64, roleName string, status Status, limitedInParentRole bool) (result []*Group, err error) {
	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 验证角色信息
	var roleId int64
	var parentRoleId int64 = 0
	if roleName != "" {
		role, err := nRepo.GetRoleWithName(ctx, roleName)
		if err != nil {
			return nil, err
		}
		if role == nil {
			return nil, ErrRoleNotExist
		}
		roleId = role.Id
		if limitedInParentRole {
			parentRoleId = role.ParentId
		}
	}

	groupList, err := nRepo.GetGroups(ctx, GroupPermission, status, "")
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

	pList, err := nRepo.GetPermissions(ctx, status, "", groupIds, parentRoleId, roleId)
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

func (this *odinService) GetPermissionsTreeWithRoleId(ctx, roleId int64, status Status, limitedInParentRole bool) (result []*Group, err error) {
	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 验证角色信息
	var parentRoleId int64 = 0
	if roleId > 0 {
		role, err := nRepo.GetRoleWithId(ctx, roleId)
		if err != nil {
			return nil, err
		}
		if role == nil {
			return nil, ErrRoleNotExist
		}
		roleId = role.Id
		if limitedInParentRole {
			parentRoleId = role.ParentId
		}
	}

	groupList, err := nRepo.GetGroups(ctx, GroupPermission, status, "")
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

	pList, err := nRepo.GetPermissions(ctx, status, "", groupIds, parentRoleId, roleId)
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

func (this *odinService) CheckPermission(ctx int64, target string, permissionName string) bool {
	return this.repo.CheckPermission(ctx, target, permissionName)
}

func (this *odinService) CheckPermissionWithId(ctx int64, target string, permissionId int64) bool {
	return this.repo.CheckPermissionWithId(ctx, target, permissionId)
}

func (this *odinService) CheckRolePermission(ctx int64, roleName, permissionName string) bool {
	return this.repo.CheckRolePermission(ctx, roleName, permissionName)
}

func (this *odinService) CheckRolePermissionWithId(ctx, roleId, permissionId int64) bool {
	return this.repo.CheckRolePermissionWithId(ctx, roleId, permissionId)
}

func (this *odinService) CleanCache(ctx int64, target string) {
	this.repo.CleanCache(ctx, target)
}
