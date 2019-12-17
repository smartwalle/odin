package odin

import "errors"

var (
	ErrRoleNameExists        = errors.New("角色名已存在")
	ErrParentRoleNotExist    = errors.New("父角色不存在")
	ErrRoleNotExist          = errors.New("角色不存在")
	ErrTargetNotAllowed      = errors.New("不合法的 Target Id")
	ErrGrantFailed           = errors.New("授权失败")
	ErrPermissionNotExist    = errors.New("权限不存在")
	ErrPermissionNameExists  = errors.New("权限名已存在")
	ErrGroupNotExist         = errors.New("组不存在")
	ErrGroupNameExists       = errors.New("组名已存在")
	ErrRevokeFailed          = errors.New("取消授权失败")
	ErrInvalidParentRole     = errors.New("无效的父角色")
	ErrPermissionOutOfParent = errors.New("授予权限超出父角色范围")
	ErrMutexRoleNotExist     = errors.New("互斥角色不存在")
	ErrPreRoleNotExist       = errors.New("前置角色不存在")
	ErrPrePermissionNotExist = errors.New("前置权限不存在")
)
