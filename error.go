package odin

import "errors"

var (
	ErrPermissionIdentifierExists = errors.New("权限已经存在")
	ErrRoleNotExists              = errors.New("角色不存在")
	ErrObjectNotAllowed           = errors.New("不合法的 Object")
	ErrGrantFailed                = errors.New("授权失败")
)
