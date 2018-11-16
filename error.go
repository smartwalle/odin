package odin

import (
	"errors"
)

var (
	ErrPermissionIdentifierExists = errors.New("权限已经存在")
	ErrRoleNotExist               = errors.New("角色不存在")
	ErrObjectNotAllowed           = errors.New("不合法的 Object")
	ErrGrantFailed                = errors.New("授权失败")
	ErrGroupNotExist              = errors.New("组不存在")
	ErrRemoveGroupNotAllowed      = errors.New("组不能被删除")
	ErrGroupExists                = errors.New("组名已经存在")
	ErrPermissionNameExists       = errors.New("权限名称已经存在")
	ErrRoleNameExists             = errors.New("角色名称已经存在")
)
