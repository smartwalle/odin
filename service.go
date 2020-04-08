package odin

import (
	"fmt"
	"github.com/smartwalle/dbs"
)

type Repository interface {
	BeginTx() (dbs.TX, Repository)

	WithTx(tx dbs.TX) Repository

	// UseIdGenerator 设置 id 生成器，默认使用 dbs 库提供的 id 生成器
	UseIdGenerator(g dbs.IdGenerator)

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

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	var s = &Service{}
	s.repo = repo
	return s
}

// Init 执行初始化操作，目前主要功能为初始化数据库表。
//
// 虽然此方法可以被重复调用，但是外部应该尽量控制此方法只在需要的时候调用。
func (this *Service) Init() error {
	var tx, nRepo = this.repo.BeginTx()
	if err := nRepo.InitTable(); err != nil {
		return err
	}
	tx.Commit()
	return nil
}

// GetPermissionGroups 获取权限组列表
func (this *Service) GetPermissionGroups(ctx int64, status Status, keywords string) (result []*Group, err error) {
	return this.repo.GetGroups(ctx, GroupPermission, status, keywords)
}

// GetPermissionGroup 根据 groupName 获取权限组信息
func (this *Service) GetPermissionGroup(ctx int64, groupName string) (result *Group, err error) {
	return this.repo.GetGroupWithName(ctx, GroupPermission, groupName)
}

// GetPermissionGroupWithId 根据组 groupId 获取权限组信息
func (this *Service) GetPermissionGroupWithId(ctx, groupId int64) (result *Group, err error) {
	return this.repo.GetGroupWithId(ctx, GroupPermission, groupId)
}

func (this *Service) addGroup(ctx int64, gType GroupType, groupName, aliasName string, status Status) (result int64, err error) {
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

// AddPermissionGroup 添加权限组信息
func (this *Service) AddPermissionGroup(ctx int64, groupName, aliasName string, status Status) (result int64, err error) {
	return this.addGroup(ctx, GroupPermission, groupName, aliasName, status)
}

func (this *Service) updateGroup(ctx int64, gType GroupType, groupName, aliasName string, status Status) (err error) {
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

func (this *Service) updateGroupWithId(ctx int64, gType GroupType, groupId int64, aliasName string, status Status) (err error) {
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

// UpdatePermissionGroup 根据 groupName 更新权限组信息
func (this *Service) UpdatePermissionGroup(ctx int64, groupName string, aliasName string, status Status) (err error) {
	return this.updateGroup(ctx, GroupPermission, groupName, aliasName, status)
}

// UpdatePermissionGroupWithId 根据 groupId 更新权限组信息
func (this *Service) UpdatePermissionGroupWithId(ctx, groupId int64, aliasName string, status Status) (err error) {
	return this.updateGroupWithId(ctx, GroupPermission, groupId, aliasName, status)
}

func (this *Service) updateGroupStatusWithId(ctx int64, gType GroupType, groupId int64, status Status) (err error) {
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

func (this *Service) updateGroupStatus(ctx int64, gType GroupType, groupName string, status Status) (err error) {
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

// UpdatePermissionGroupStatus 根据 groupName 更新权限组状态
func (this *Service) UpdatePermissionGroupStatus(ctx int64, groupName string, status Status) (err error) {
	return this.updateGroupStatus(ctx, GroupPermission, groupName, status)
}

// UpdatePermissionGroupStatusWithId 根据 groupId 更新权限组状态
func (this *Service) UpdatePermissionGroupStatusWithId(ctx int64, groupId int64, status Status) (err error) {
	return this.updateGroupStatusWithId(ctx, GroupPermission, groupId, status)
}

// GetPermissions 获取权限列表
func (this *Service) GetPermissions(ctx int64, status Status, keywords string, groupIds []int64) (result []*Permission, err error) {
	return this.repo.GetPermissions(ctx, status, keywords, groupIds, 0, 0)
}

// GetPermission 根据 permissionName 获取权限信息
func (this *Service) GetPermission(ctx int64, permissionName string) (result *Permission, err error) {
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

// GetPermissionWithId 根据 permissionId 获取权限信息
func (this *Service) GetPermissionWithId(ctx, permissionId int64) (result *Permission, err error) {
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

// AddPermissionWithGroup 添加权限
func (this *Service) AddPermissionWithGroup(ctx int64, groupName, permissionName, aliasName, description string, status Status) (result int64, err error) {
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

// AddPermissionWithGroupId 添加权限
func (this *Service) AddPermissionWithGroupId(ctx, groupId int64, permissionName, aliasName, description string, status Status) (result int64, err error) {
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

// UpdatePermission 根据 permissionName 更新权限信息
func (this *Service) UpdatePermission(ctx int64, permissionName, groupName, aliasName, description string, status Status) (err error) {
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

// UpdatePermissionWithId 根据 permissionId 更新权限信息
func (this *Service) UpdatePermissionWithId(ctx, permissionId, groupId int64, aliasName, description string, status Status) (err error) {
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

// UpdatePermissionStatus 根据 permissionName 更新权限状态
func (this *Service) UpdatePermissionStatus(ctx int64, permissionName string, status Status) (err error) {
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

// UpdatePermissionStatusWithId 根据 permissionId 更新权限状态
func (this *Service) UpdatePermissionStatusWithId(ctx, permissionId int64, status Status) (err error) {
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

// GrantPermission 授予权限给角色
func (this *Service) GrantPermission(ctx int64, roleName string, permissionNames ...string) (err error) {
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

// GrantPermissionWithId 授予权限给角色
func (this *Service) GrantPermissionWithId(ctx int64, roleId int64, permissionIds ...int64) (err error) {
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

// ReGrantPermission 授予权限给角色，会将原有的权限先取消掉
func (this *Service) ReGrantPermission(ctx int64, roleName string, permissionNames ...string) (err error) {
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

// ReGrantPermissionWithId 授予权限给角色，会将原有的权限先取消掉
func (this *Service) ReGrantPermissionWithId(ctx int64, roleId int64, permissionIds ...int64) (err error) {
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

// RevokePermission 取消对角色的指定权限授权
func (this *Service) RevokePermission(ctx int64, roleName string, permissionNames ...string) (err error) {
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

// RevokePermissionWithId 取消对角色的指定权限授权
func (this *Service) RevokePermissionWithId(ctx int64, roleId int64, permissionIds ...int64) (err error) {
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

// RevokeAllPermission 取消对角色的所有权限授权
func (this *Service) RevokeAllPermission(ctx int64, roleName string) (err error) {
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

// RevokeAllPermissionWithId 取消对角色的所有权限授权
func (this *Service) RevokeAllPermissionWithId(ctx, roleId int64) (err error) {
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

// AddPrePermission 添加授予该权限时需要的先决条件
func (this *Service) AddPrePermission(ctx int64, permissionName string, prePermissionNames ...string) (err error) {
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

// AddPrePermissionWithId 添加授予该权限时需要的先决条件
func (this *Service) AddPrePermissionWithId(ctx, permissionId int64, prePermissionIds ...int64) (err error) {
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

// RemovePrePermission 删除授予该权限时需要的先决条件
func (this *Service) RemovePrePermission(ctx int64, permissionName string, prePermissionNames ...string) (err error) {
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

// RemovePrePermissionWithId 删除授予该权限时需要的先决条件
func (this *Service) RemovePrePermissionWithId(ctx, permissionId int64, prePermissionIds ...int64) (err error) {
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

// RemoveAllPrePermission 删除授予该权限时需要的所有先决条件
func (this *Service) RemoveAllPrePermission(ctx int64, permissionName string) (err error) {
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

// RemoveAllPrePermissionWithId 删除授予该权限时需要的所有先决条件
func (this *Service) RemoveAllPrePermissionWithId(ctx, permissionId int64) (err error) {
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

// GetPrePermissions 获取授予该权限时需要的所有先决条件
func (this *Service) GetPrePermissions(ctx int64, permissionName string) (result []*PrePermission, err error) {
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

// GetPrePermissionsWithId 获取授予该权限时需要的所有先决条件
func (this *Service) GetPrePermissionsWithId(ctx int64, permissionId int64) (result []*PrePermission, err error) {
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

// GetRoles 获取角色列表
//
// 如果参数 isGrantedToTarget 的值不为空字符串，则返回的角色数据中将包含该角色（通过 Granted 判断）是否已授权给 isGrantedToTarget
//
// 如果参数 limitedInTarget 的值不为空字符串， 则返回的角色数据将限定在 limitedInTarget 已拥有的角色及其子角色范围内
//
// 返回的角色数据的 Granted 字段参照的是 isGrantedToTarget
//
// 返回的角色数据的 Accessible 字段参照的是 limitedInTarget
func (this *Service) GetRoles(ctx int64, status Status, keywords, isGrantedToTarget, limitedInTarget string) (result []*Role, err error) {
	if limitedInTarget == "" {
		return this.repo.GetRoles(ctx, -1, status, keywords, isGrantedToTarget)
	}
	return this.repo.GetRolesInTarget(ctx, limitedInTarget, status, keywords, isGrantedToTarget)
}

// GetRolesWithParent 获取角色列表
//
// 如果参数 isGrantedToTarget 的值不为空字符串，则返回的角色数据中将包含该角色（通过 Granted 判断）是否已授权给 isGrantedToTarget
func (this *Service) GetRolesWithParent(ctx int64, parentRoleName string, status Status, keywords, isGrantedToTarget string) (result []*Role, err error) {
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

// GetRolesWithParentId 获取角色列表
//
// 如果参数 isGrantedToTarget 的值不为空字符串，则返回的角色数据中将包含该角色（通过 Granted 判断）是否已授权给 isGrantedToTarget
func (this *Service) GetRolesWithParentId(ctx, parentRoleId int64, status Status, keywords, isGrantedToTarget string) (result []*Role, err error) {
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

// GetRole 根据 roleName 获取角色信息
func (this *Service) GetRole(ctx int64, name string) (result *Role, err error) {
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

// GetRoleWithId 根据 roleId 获取角色信息
func (this *Service) GetRoleWithId(ctx, roleId int64) (result *Role, err error) {
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

// AddRole 添加角色
func (this *Service) AddRole(ctx int64, roleName, aliasName, description string, status Status) (result int64, err error) {
	return this.AddRoleWithParentId(ctx, 0, roleName, aliasName, description, status)
}

// AddRoleWithParent 添加角色，新添加的角色将作为 parentRoleName 的子角色
//
// 调用时应该确认操作者是否有访问 parentRoleName 的权限，即 parentRoleName 是否为当前操作者拥有的角色及其子角色
func (this *Service) AddRoleWithParent(ctx int64, parentRoleName, roleName, aliasName, description string, status Status) (result int64, err error) {
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

// AddRoleWithParentId 添加角色，新添加的角色将作为 parentRoleId 的子角色
//
// 调用时应该确认操作者是否有访问 parentRoleId 的权限，即 parentRoleId 是否为当前操作者拥有的角色及其子角色
func (this *Service) AddRoleWithParentId(ctx, parentRoleId int64, roleName, aliasName, description string, status Status) (result int64, err error) {
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

// UpdateRole 根据 roleName 更新角色信息
func (this *Service) UpdateRole(ctx int64, roleName, aliasName, description string, status Status) (err error) {
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

// UpdateRoleWithId 根据 roleId 更新角色信息
func (this *Service) UpdateRoleWithId(ctx, roleId int64, aliasName, description string, status Status) (err error) {
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

// UpdateRoleStatus 根据 roleName 更新角色的状态
func (this *Service) UpdateRoleStatus(ctx int64, roleName string, status Status) (err error) {
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

// UpdateRoleStatusWithId 根据 roleId 更新角色的状态
func (this *Service) UpdateRoleStatusWithId(ctx, roleId int64, status Status) (err error) {
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

// GrantRole 授权角色给 target
func (this *Service) GrantRole(ctx int64, target string, roleNames ...string) (err error) {
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

// GrantRoleWithId 授权角色给 target
func (this *Service) GrantRoleWithId(ctx int64, target string, roleIds ...int64) (err error) {
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

// ReGrantRole 授权角色给 target，会将原有的角色授权先取消掉
func (this *Service) ReGrantRole(ctx int64, target string, roleNames ...string) (err error) {
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

// ReGrantRoleWithId 授权角色给 target，会将原有的角色授权先取消掉
func (this *Service) ReGrantRoleWithId(ctx int64, target string, roleIds ...int64) (err error) {
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

// RevokeRole 取消对 target 的角色授权
func (this *Service) RevokeRole(ctx int64, target string, roleNames ...string) (err error) {
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

// RevokeRoleWithId 取消对 target 的角色授权
func (this *Service) RevokeRoleWithId(ctx int64, target string, roleIds ...int64) (err error) {
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

// RevokeAllRole 取消对 target 的所有角色授权
func (this *Service) RevokeAllRole(ctx int64, target string) (err error) {
	return this.repo.RevokeAllRole(ctx, target)
}

// AddRoleMutex 添加角色互斥关系
func (this *Service) AddRoleMutex(ctx int64, roleName string, mutexRoleNames ...string) (err error) {
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

// AddRoleMutexWithId 添加角色互斥关系
func (this *Service) AddRoleMutexWithId(ctx, roleId int64, mutexRoleIds ...int64) (err error) {
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

// RemoveRoleMutex 删除角色互斥关系
func (this *Service) RemoveRoleMutex(ctx int64, roleName string, mutexRoleNames ...string) (err error) {
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

// RemoveRoleMutexWithId 删除角色互斥关系
func (this *Service) RemoveRoleMutexWithId(ctx, roleId int64, mutexRoleIds ...int64) (err error) {
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

// RemoveAllRoleMutex 删除该角色所有的互斥关系
func (this *Service) RemoveAllRoleMutex(ctx int64, roleName string) (err error) {
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

// RemoveAllRoleMutexWithId 删除该角色所有的互斥关系
func (this *Service) RemoveAllRoleMutexWithId(ctx, roleId int64) (err error) {
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

// GetMutexRoles 获取与该角色互斥的角色列表
func (this *Service) GetMutexRoles(ctx int64, roleName string) (result []*RoleMutex, err error) {
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

// GetMutexRolesWithId 获取与该角色互斥的角色列表
func (this *Service) GetMutexRolesWithId(ctx, roleId int64) (result []*RoleMutex, err error) {
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

// CheckRoleMutex 验证两个角色是否互斥
func (this *Service) CheckRoleMutex(ctx int64, roleName, mutexRoleName string) bool {
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

// CheckRoleMutexWithId 验证两个角色是否互斥
func (this *Service) CheckRoleMutexWithId(ctx, roleId, mutexRoleId int64) bool {
	return this.repo.CheckRoleMutex(ctx, roleId, mutexRoleId)
}

// AddPreRole 添加授予该角色时需要的先决条件
func (this *Service) AddPreRole(ctx int64, roleName string, preRoleNames ...string) (err error) {
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

// AddPreRoleWithId 添加授予该角色时需要的先决条件
func (this *Service) AddPreRoleWithId(ctx, roleId int64, preRoleIds ...int64) (err error) {
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

// RemovePreRole 删除授予该角色时需要的先决条件
func (this *Service) RemovePreRole(ctx int64, roleName string, preRoleNames ...string) (err error) {
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

// RemovePreRoleWithId 删除授予该角色时需要的先决条件
func (this *Service) RemovePreRoleWithId(ctx, roleId int64, preRoleIds ...int64) (err error) {
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

// RemoveAllPreRole 删除授予该角色时需要的所有先决条件
func (this *Service) RemoveAllPreRole(ctx int64, roleName string) (err error) {
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

// RemoveAllPreRoleWithId 删除授予该角色时需要的所有先决条件
func (this *Service) RemoveAllPreRoleWithId(ctx, roleId int64) (err error) {
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

// GetPreRoles 获取授予该角色时需要的所有先决条件
func (this *Service) GetPreRoles(ctx int64, roleName string) (result []*PreRole, err error) {
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

// GetPreRolesWithId 获取授予该角色时需要的所有先决条件
func (this *Service) GetPreRolesWithId(ctx, roleId int64) (result []*PreRole, err error) {
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

// GetGrantedRoles 获取已授权给 target 的角色列表
func (this *Service) GetGrantedRoles(ctx int64, target string) (result []*Role, err error) {
	return this.repo.GetGrantedRoles(ctx, target, false)
}

// GetRolesWithTarget 获取已授权给 target 的角色，及其角色的子角色
func (this *Service) GetRolesWithTarget(ctx int64, target string) (result []*Role, err error) {
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

// CheckRole 验证 target 是否拥有指定角色
func (this *Service) CheckRole(ctx int64, target string, roleName string) bool {
	return this.repo.CheckRole(ctx, target, roleName)
}

// CheckRoleWithId 验证 target 是否拥有指定角色
func (this *Service) CheckRoleWithId(ctx int64, target string, roleId int64) bool {
	return this.repo.CheckRoleWithId(ctx, target, roleId)
}

// CheckRoleAccessible 验证 target 是否拥有操作访问 roleName 的权限
func (this *Service) CheckRoleAccessible(ctx int64, target string, roleName string) bool {
	return this.repo.CheckRoleAccessible(ctx, target, roleName)
}

// CheckRoleAccessibleWithId 验证 target 是否拥有操作访问 roleId 的权限
func (this *Service) CheckRoleAccessibleWithId(ctx int64, target string, roleId int64) bool {
	return this.repo.CheckRoleAccessibleWithId(ctx, target, roleId)
}

// GetPermissionsWithRole 获取已授权给 roleName 的权限列表
func (this *Service) GetPermissionsWithRole(ctx int64, roleName string) (result []*Permission, err error) {
	role, err := this.repo.GetRoleWithName(ctx, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, ErrRoleNotExist
	}
	return this.repo.GetPermissionsWithRoleId(ctx, role.Id)
}

// GetPermissionsWithRoleId 获取已授权给 roleId 的权限列表
func (this *Service) GetPermissionsWithRoleId(ctx int64, roleId int64) (result []*Permission, err error) {
	role, err := this.repo.GetRoleWithId(ctx, roleId)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, ErrRoleNotExist
	}
	return this.repo.GetPermissionsWithRoleId(ctx, role.Id)
}

// GetGrantedPermissions 获取已授权给 target 的权限列表
func (this *Service) GetGrantedPermissions(ctx int64, target string) (result []*Permission, err error) {
	return this.repo.GetGrantedPermissions(ctx, target)
}

// GetPermissionsTreeWithRole 获取权限组列表，组中包含该组所有的权限信息
//
// 如果参数 roleName 的值不为空字符串，则返回的权限数据中将附带该权限是否已授权给该 roleName
//
// 如果参数 limitedInParentRole 的值为 true，并且 roleName 的值不为空字符串，则返回的权限数据将限定在 roleName 的父角色拥有的权限范围内.
func (this *Service) GetPermissionsTreeWithRole(ctx int64, roleName string, status Status, limitedInParentRole bool) (result []*Group, err error) {
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

// GetPermissionsTreeWithRoleId 获取权限组列表，组中包含该组所有的权限信息
//
// 如果参数 roleId 的值大于 0，则返回的权限数据中将附带该权限是否已授权给该 roleId
//
// 如果参数 limitedInParentRole 的值为 true，并且 roleId 的值大于 0，则返回的权限数据将限定在 roleId 的父角色拥有的权限范围内
func (this *Service) GetPermissionsTreeWithRoleId(ctx, roleId int64, status Status, limitedInParentRole bool) (result []*Group, err error) {
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

// CheckPermission 验证 target 是否拥有指定权限
func (this *Service) CheckPermission(ctx int64, target string, permissionName string) bool {
	return this.repo.CheckPermission(ctx, target, permissionName)
}

// CheckPermissionWithId 验证 target 是否拥有指定权限
func (this *Service) CheckPermissionWithId(ctx int64, target string, permissionId int64) bool {
	return this.repo.CheckPermissionWithId(ctx, target, permissionId)
}

// CheckRolePermission 验证角色是否拥有指定权限
func (this *Service) CheckRolePermission(ctx int64, roleName, permissionName string) bool {
	return this.repo.CheckRolePermission(ctx, roleName, permissionName)
}

// CheckRolePermissionWithId 验证角色是否拥有指定权限
func (this *Service) CheckRolePermissionWithId(ctx, roleId, permissionId int64) bool {
	return this.repo.CheckRolePermissionWithId(ctx, roleId, permissionId)
}

// CleanCache 清除缓存，如果 target 为空字符串或者 target 的值为星号(*)，则会清空所有缓存
func (this *Service) CleanCache(ctx int64, target string) {
	this.repo.CleanCache(ctx, target)
}
