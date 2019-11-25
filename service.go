package odin

import (
	"github.com/smartwalle/dbs"
)

type Repository interface {
	BeginTx() (dbs.TX, Repository)

	WithTx(tx dbs.TX) Repository

	// permission
	GetPermissionList(ctx int64, status Status, keywords string) (result []*Permission, err error)

	GetPermissionListWithIds(ctx int64, permissionIds ...int64) (result []*Permission, err error)

	GetPermissionListWithNames(ctx int64, names ...string) (result []*Permission, err error)

	GetPermissionListWithRoleId(ctx int64, roleId int64) (result []*Permission, err error)

	GetPermissionWithId(ctx, permission int64) (result *Permission, err error)

	GetPermissionWithName(ctx int64, name string) (result *Permission, err error)

	AddPermission(ctx int64, name, aliasName, description string, status Status) (result int64, err error)

	UpdatePermission(ctx, permissionId int64, name, aliasName, description string, status Status) (err error)

	UpdatePermissionStatus(ctx, permissionId int64, status Status) (err error)

	GrantPermissionWithIds(ctx, roleId int64, permissionIds ...int64) (err error)

	RevokePermissionWithIds(ctx, roleId int64, permissionIds ...int64) (err error)

	RevokeAllPermission(ctx, roleId int64) (err error)

	// role
	GetRoleList(ctx int64, targetId string, status Status, keywords string) (result []*Role, err error)

	GetRoleListWithIds(ctx int64, roleIds ...int64) (result []*Role, err error)

	GetRoleListWithNames(ctx int64, names ...string) (result []*Role, err error)

	GetRoleWithId(ctx, roleId int64) (result *Role, err error)

	GetRoleWithName(ctx int64, name string) (result *Role, err error)

	AddRole(ctx, parentId int64, name, aliasName, description string, status Status) (result int64, err error)

	UpdateRole(ctx, roleId int64, name, aliasName, description string, status Status) (err error)

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

func (this *odinService) GetPermissionList(ctx int64, status Status, keywords string) (result []*Permission, err error) {
	return this.repo.GetPermissionList(ctx, status, keywords)
}

func (this *odinService) GetPermissionWithId(ctx, permissionId int64) (result *Permission, err error) {
	result, err = this.repo.GetPermissionWithId(ctx, permissionId)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (this *odinService) GetPermissionWithName(ctx int64, name string) (result *Permission, err error) {
	result, err = this.repo.GetPermissionWithName(ctx, name)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (this *odinService) AddPermission(ctx int64, name, aliasName, description string, status Status) (result int64, err error) {
	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 验证 name 是否已经存在
	permission, err := nRepo.GetPermissionWithName(ctx, name)
	if err != nil {
		return 0, err
	}
	if permission != nil {
		return 0, ErrPermissionNameExists
	}

	if result, err = nRepo.AddPermission(ctx, name, aliasName, description, status); err != nil {
		return 0, err
	}

	tx.Commit()
	return result, nil
}

func (this *odinService) UpdatePermission(ctx, permissionId int64, name, aliasName, description string, status Status) (err error) {
	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 验证角色是否存在
	permission, err := nRepo.GetPermissionWithId(ctx, permissionId)
	if err != nil {
		return err
	}
	if permission == nil {
		return ErrPermissionNotExist
	}

	// 验证 name 是否已经存在
	permission, err = nRepo.GetPermissionWithName(ctx, name)
	if err != nil {
		return err
	}
	if permission != nil && permission.Id != permissionId {
		return ErrPermissionNameExists
	}

	if err = nRepo.UpdatePermission(ctx, permissionId, name, aliasName, description, status); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) UpdatePermissionStatus(ctx, permissionId int64, status Status) (err error) {
	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 验证权限是否存在
	role, err := nRepo.GetPermissionWithId(ctx, permissionId)
	if err != nil {
		return err
	}
	if role == nil {
		return ErrRoleNotExist
	}

	if err = nRepo.UpdatePermissionStatus(ctx, permissionId, status); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) GrantPermission(ctx int64, roleId int64, names ...string) (err error) {
	if len(names) == 0 {
		return ErrPermissionNotExist
	}

	if roleId <= 0 {
		return ErrTargetNotAllowed
	}

	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	permissionList, err := nRepo.GetPermissionListWithNames(ctx, names...)
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

	if err = nRepo.GrantPermissionWithIds(ctx, roleId, nIds...); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) GrantPermissionWithIds(ctx int64, roleId int64, permissionIds ...int64) (err error) {
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

func (this *odinService) ReGrantPermission(ctx int64, roleId int64, names ...string) (err error) {
	if len(names) == 0 {
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

	permissionList, err := nRepo.GetPermissionListWithNames(ctx, names...)
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

	if err = nRepo.RevokeAllPermission(ctx, roleId); err != nil {
		return err
	}

	if err = nRepo.GrantPermissionWithIds(ctx, roleId, nIds...); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) ReGrantPermissionWithIds(ctx int64, roleId int64, permissionIds ...int64) (err error) {
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

	if err = nRepo.RevokeAllPermission(ctx, roleId); err != nil {
		return err
	}

	if err = nRepo.GrantPermissionWithIds(ctx, roleId, nIds...); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) RevokePermission(ctx int64, roleId int64, names ...string) (err error) {
	if len(names) == 0 {
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

	permissionList, err := nRepo.GetPermissionListWithNames(ctx, names...)
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

	if err = nRepo.RevokePermissionWithIds(ctx, roleId, nIds...); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) RevokePermissionWithIds(ctx int64, roleId int64, permissionIds ...int64) (err error) {
	return this.repo.RevokePermissionWithIds(ctx, roleId, permissionIds...)
}

func (this *odinService) RevokeAllPermission(ctx, roleId int64) (err error) {
	return this.repo.RevokeAllPermission(ctx, roleId)
}

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

func (this *odinService) GetRoleWithName(ctx int64, name string) (result *Role, err error) {
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

func (this *odinService) AddRole(ctx, parentId int64, name, aliasName, description string, status Status) (result int64, err error) {
	var tx, nRepo = this.repo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 验证 parent id 是否存在
	if parentId > 0 {
		role, err := nRepo.GetRoleWithId(ctx, parentId)
		if err != nil {
			return 0, err
		}
		if role == nil {
			return 0, ErrParentRoleNotExist
		}
	}

	// 验证 name 是否已经存在
	role, err := nRepo.GetRoleWithName(ctx, name)
	if err != nil {
		return 0, err
	}
	if role != nil {
		return 0, ErrRoleNameExists
	}

	if result, err = nRepo.AddRole(ctx, parentId, name, aliasName, description, status); err != nil {
		return 0, err
	}

	tx.Commit()
	return result, nil
}

func (this *odinService) UpdateRole(ctx, roleId int64, name, aliasName, description string, status Status) (err error) {
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

	// 验证 name 是否已经存在
	role, err = nRepo.GetRoleWithName(ctx, name)
	if err != nil {
		return err
	}
	if role != nil && role.Id != roleId {
		return ErrRoleNameExists
	}

	if err = nRepo.UpdateRole(ctx, roleId, name, aliasName, description, status); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) UpdateRoleStatus(ctx, roleId int64, status Status) (err error) {
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

	if err = nRepo.UpdateRoleStatus(ctx, roleId, status); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinService) GetGrantedRoleList(ctx int64, targetId string) (result []*Role, err error) {
	return this.repo.GetGrantedRoleList(ctx, targetId)
}

func (this *odinService) GrantRole(ctx int64, targetId string, names ...string) (err error) {
	if len(names) == 0 {
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

	roleList, err := nRepo.GetRoleListWithNames(ctx, names...)
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

func (this *odinService) GrantRoleWithIds(ctx int64, targetId string, roleIds ...int64) (err error) {
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

func (this *odinService) ReGrantRole(ctx int64, targetId string, names ...string) (err error) {
	if len(names) == 0 {
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

	roleList, err := nRepo.GetRoleListWithNames(ctx, names...)
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

func (this *odinService) ReGrantRoleWithIds(ctx int64, targetId string, roleIds ...int64) (err error) {
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

func (this *odinService) RevokeRole(ctx int64, targetId string, names ...string) (err error) {
	if len(names) == 0 {
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

	roleList, err := nRepo.GetRoleListWithNames(ctx, names...)
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

func (this *odinService) RevokeRoleWithIds(ctx int64, targetId string, roleIds ...int64) (err error) {
	return this.repo.RevokeRoleWithIds(ctx, targetId, roleIds...)
}

func (this *odinService) RevokeAllRole(ctx int64, targetId string) (err error) {
	return this.repo.RevokeAllRole(ctx, targetId)
}
