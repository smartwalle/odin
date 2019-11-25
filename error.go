package odin

import "errors"

var (
	ErrRoleNameExists       = errors.New("角色名已存在")
	ErrParentRoleNotExist   = errors.New("父角色不存在")
	ErrRoleNotExist         = errors.New("角色不存在")
	ErrTargetNotAllowed     = errors.New("不合法的 Target Id")
	ErrGrantFailed          = errors.New("授权失败")
	ErrPermissionNotExist   = errors.New("权限不存在")
	ErrPermissionNameExists = errors.New("权限名已存在")
)
