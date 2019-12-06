package odin

import (
	"github.com/smartwalle/dbs"
)

type Repository interface {
	BeginTx() (dbs.TX, Repository)

	WithTx(tx dbs.TX) Repository

	// InitTable 初始化数据库表
	InitTable() error

	// group
	GetGroups(ctx int64, gType GroupType, status Status, keywords string) (result []*Group, err error)

	GetGroupWithId(ctx int64, gType GroupType, groupId int64) (result *Group, err error)

	GetGroupWithName(ctx int64, gType GroupType, name string) (result *Group, err error)

	AddGroup(ctx int64, gType GroupType, name, aliasName string, status Status) (result int64, err error)

	UpdateGroup(ctx int64, gType GroupType, groupId int64, aliasName string, status Status) (err error)

	UpdateGroupStatus(ctx int64, gType GroupType, groupId int64, status Status) (err error)

	// permission

	// GetPermissions 获取角色列表
	// 如果参数 limitedInRole 的值大于 0，则返回的权限数据将限定在已授权给 limitedInRole 的权限范围之内
	// 如果参数 isGrantedToRole 的值大于 0，则返回的权限数据中将附带该权限是否已授权给该 isGrantedToRole
	GetPermissions(ctx int64, status Status, keywords string, groupIds []int64, limitedInRole, isGrantedToRole int64) (result []*Permission, err error)

	GetPermissionsWithIds(ctx int64, permissionIds ...int64) (result []*Permission, err error)

	GetPermissionsWithNames(ctx int64, names ...string) (result []*Permission, err error)

	GetPermissionsWithRoleId(ctx int64, roleId int64) (result []*Permission, err error)

	GetPermissionWithId(ctx, permissionId int64) (result *Permission, err error)

	GetPermissionWithName(ctx int64, name string) (result *Permission, err error)

	AddPermission(ctx, groupId int64, name, aliasName, description string, status Status) (result int64, err error)

	UpdatePermission(ctx, permissionId, groupId int64, aliasName, description string, status Status) (err error)

	UpdatePermissionStatus(ctx, permissionId int64, status Status) (err error)

	GrantPermissionWithIds(ctx, roleId int64, permissionIds []int64) (err error)

	RevokePermissionWithIds(ctx, roleId int64, permissionIds []int64) (err error)

	RevokeAllPermission(ctx, roleId int64) (err error)

	GetGrantedPermissions(ctx int64, target string) (result []*Permission, err error)

	// role

	// GetRoles 获取角色列表
	// 如果参数 parentId 的值大于等于 0，则表示查询 parentId 的子角色列表
	// 如果参数 isGrantedToTarget 的值不为空字符串，则返回的角色数据中将包含该角色（通过 Granted 判断）是否已授权给 isGrantedToTarget
	GetRoles(ctx int64, parentId int64, status Status, keywords string, isGrantedToTarget string) (result []*Role, err error)

	GetRolesWithIds(ctx int64, roleIds ...int64) (result []*Role, err error)

	GetRolesWithNames(ctx int64, names ...string) (result []*Role, err error)

	GetRoleWithId(ctx, roleId int64) (result *Role, err error)

	GetRoleWithName(ctx int64, name string) (result *Role, err error)

	AddRole(ctx int64, parent *Role, name, aliasName, description string, status Status) (result int64, err error)

	UpdateRole(ctx, roleId int64, aliasName, description string, status Status) (err error)

	UpdateRoleStatus(ctx, roleId int64, status Status) (err error)

	GetGrantedRoles(ctx int64, target string, withChildren bool) (result []*Role, err error)

	GrantRoleWithIds(ctx int64, target string, roleIds ...int64) (err error)

	RevokeRoleWithIds(ctx int64, target string, roleIds ...int64) (err error)

	RevokeAllRole(ctx int64, target string) (err error)

	CheckPermission(ctx int64, target string, permissionName string) bool

	CheckPermissionWithId(ctx int64, target string, permissionId int64) bool

	CheckRole(ctx int64, target string, roleName string) bool

	CheckRoleWithId(ctx int64, target string, roleId int64) bool

	CheckRoleAccessible(ctx int64, target string, roleName string) bool

	CheckRoleAccessibleWithId(ctx int64, target string, roleId int64) bool

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

// group

func (this *odinService) GetPermissionGroups(ctx int64, status Status, keywords string) (result []*Group, err error) {
	return this.repo.GetGroups(ctx, GroupPermission, status, keywords)
}

func (this *odinService) GetPermissionGroupWithId(ctx, groupId int64) (result *Group, err error) {
	return this.repo.GetGroupWithId(ctx, GroupPermission, groupId)
}

func (this *odinService) GetPermissionGroup(ctx int64, groupName string) (result *Group, err error) {
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

func (this *odinService) GetPermissions(ctx int64, status Status, keywords string, groupIds []int64) (result []*Permission, err error) {
	return this.repo.GetPermissions(ctx, status, keywords, groupIds, 0, 0)
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
				return ErrPermissionDenied
			}
		}
	}

	permissionList, err := nRepo.GetPermissionsWithIds(ctx, permissionIds...)
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

	if err = nRepo.GrantPermissionWithIds(ctx, roleId, nIds); err != nil {
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
				return ErrPermissionDenied
			}
		}
	}

	permissionList, err := nRepo.GetPermissionsWithNames(ctx, permissionNames...)
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
				return ErrPermissionDenied
			}
		}
	}

	permissions, err := nRepo.GetPermissionsWithIds(ctx, permissionIds...)
	if err != nil {
		return err
	}

	var nIds = make([]int64, 0, len(permissions))
	var nPermissionMap = make(map[int64]struct{})
	for _, role := range permissions {
		nIds = append(nIds, role.Id)
		nPermissionMap[role.Id] = struct{}{}
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
		if _, ok := nPermissionMap[p.Id]; ok == false {
			revokeIds = append(revokeIds, p.Id)
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
				return ErrPermissionDenied
			}
		}
	}

	permissionList, err := nRepo.GetPermissionsWithNames(ctx, permissionNames...)
	if err != nil {
		return err
	}

	var nIds = make([]int64, 0, len(permissionList))
	var nPermissionMap = make(map[int64]struct{})
	for _, role := range permissionList {
		nIds = append(nIds, role.Id)
		nPermissionMap[role.Id] = struct{}{}
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
		if _, ok := nPermissionMap[p.Id]; ok == false {
			revokeIds = append(revokeIds, p.Id)
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

func (this *odinService) RevokePermissionWithId(ctx int64, roleId int64, permissionIds ...int64) (err error) {
	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err = nRepo.RevokePermissionWithIds(ctx, roleId, permissionIds); err != nil {
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

	permissionList, err := nRepo.GetPermissionsWithNames(ctx, permissionNames...)
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

	if err = nRepo.RevokePermissionWithIds(ctx, role.Id, nIds); err != nil {
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

func (this *odinService) GetRoles(ctx int64, status Status, keywords, isGrantedToTarget string) (result []*Role, err error) {
	return this.repo.GetRoles(ctx, -1, status, keywords, isGrantedToTarget)
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

func (this *odinService) GetRoleWithId(ctx, roleId int64) (result *Role, err error) {
	result, err = this.repo.GetRoleWithId(ctx, roleId)
	if err != nil {
		return nil, err
	}
	if result != nil {
		result.PermissionList, err = this.repo.GetPermissionsWithRoleId(ctx, result.Id)
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
		result.PermissionList, err = this.repo.GetPermissionsWithRoleId(ctx, result.Id)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (this *odinService) AddRole(ctx int64, roleName, aliasName, description string, status Status) (result int64, err error) {
	return this.AddRoleWithParentId(ctx, 0, roleName, aliasName, description, status)
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

	var nIds = make([]int64, 0, len(roleList))
	for _, role := range roleList {
		nIds = append(nIds, role.Id)
	}
	if len(nIds) == 0 {
		return ErrGrantFailed
	}

	if err = nRepo.GrantRoleWithIds(ctx, target, nIds...); err != nil {
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

	var nIds = make([]int64, 0, len(roleList))
	for _, role := range roleList {
		nIds = append(nIds, role.Id)
	}
	if len(nIds) == 0 {
		return ErrGrantFailed
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
	for _, role := range roleList {
		nIds = append(nIds, role.Id)
	}
	if len(nIds) == 0 {
		return ErrGrantFailed
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
	for _, role := range roleList {
		nIds = append(nIds, role.Id)
	}
	if len(nIds) == 0 {
		return ErrGrantFailed
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

func (this *odinService) RevokeRoleWithId(ctx int64, target string, roleIds ...int64) (err error) {
	if len(roleIds) == 0 {
		return ErrRoleNotExist
	}
	if target == "" {
		return ErrTargetNotAllowed
	}
	return this.repo.RevokeRoleWithIds(ctx, target, roleIds...)
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

	roleList, err := nRepo.GetRolesWithNames(ctx, roleNames...)
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

	if err = nRepo.RevokeRoleWithIds(ctx, target, nIds...); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) RevokeAllRole(ctx int64, target string) (err error) {
	return this.repo.RevokeAllRole(ctx, target)
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

func (this *odinService) CheckRoleWithId(ctx int64, target string, roleId int64) bool {
	return this.repo.CheckRoleWithId(ctx, target, roleId)
}

func (this *odinService) CheckRole(ctx int64, target string, roleName string) bool {
	return this.repo.CheckRole(ctx, target, roleName)
}

func (this *odinService) CheckRoleAccessibleWithId(ctx int64, target string, roleId int64) bool {
	return this.repo.CheckRoleAccessibleWithId(ctx, target, roleId)
}

func (this *odinService) CheckRoleAccessible(ctx int64, target string, roleName string) bool {
	return this.repo.CheckRoleAccessible(ctx, target, roleName)
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

func (this *odinService) GetGrantedPermissions(ctx int64, target string) (result []*Permission, err error) {
	return this.repo.GetGrantedPermissions(ctx, target)
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

func (this *odinService) CheckPermissionWithId(ctx int64, target string, permissionId int64) bool {
	return this.repo.CheckPermissionWithId(ctx, target, permissionId)
}

func (this *odinService) CheckPermission(ctx int64, target string, permissionName string) bool {
	return this.repo.CheckPermission(ctx, target, permissionName)
}

func (this *odinService) CleanCache(ctx int64, target string) {
	this.repo.CleanCache(ctx, target)
}
