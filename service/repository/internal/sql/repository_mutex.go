package sql

import (
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/odin"
	"time"
)

// AddRoleMutex 添加互斥关系
func (this *Repository) AddRoleMutex(ctx, roleId int64, mutexRoleIds []int64) (err error) {
	var now = time.Now()
	var ib = dbs.NewInsertBuilder()
	ib.UseDialect(this.dialect)
	ib.Options("IGNORE")
	ib.Table(this.tableRoleMutex)
	ib.Columns("ctx", "role_id", "mutex_role_id", "created_on")
	for _, mutexRoleId := range mutexRoleIds {
		ib.Values(ctx, roleId, mutexRoleId, now)
		ib.Values(ctx, mutexRoleId, roleId, now)
	}
	if _, err = ib.Exec(this.db); err != nil {
		return err
	}
	return nil
}

// RemoveRoleMutex 删除互斥关系
func (this *Repository) RemoveRoleMutex(ctx, roleId int64, mutexRoleIds []int64) (err error) {
	var rb = dbs.NewDeleteBuilder()
	rb.UseDialect(this.dialect)
	rb.Table(this.tableRoleMutex)
	rb.Where("ctx = ?", ctx)
	rb.Where("role_id = ?", roleId)
	rb.Where(dbs.IN("mutex_role_id", mutexRoleIds))
	rb.Limit(int64(len(mutexRoleIds)))
	if _, err = rb.Exec(this.db); err != nil {
		return err
	}

	rb = dbs.NewDeleteBuilder()
	rb.UseDialect(this.dialect)
	rb.Table(this.tableRoleMutex)
	rb.Where("ctx = ?", ctx)
	rb.Where("mutex_role_id = ?", roleId)
	rb.Where(dbs.IN("role_id", mutexRoleIds))
	rb.Limit(int64(len(mutexRoleIds)))
	if _, err = rb.Exec(this.db); err != nil {
		return err
	}
	return nil
}

// CleanRoleMutex 清除互斥关系
func (this *Repository) CleanRoleMutex(ctx, roleId int64) (err error) {
	var rb = dbs.NewDeleteBuilder()
	rb.UseDialect(this.dialect)
	rb.Table(this.tableRoleMutex)
	rb.Where("ctx = ?", ctx)
	rb.Where("role_id = ?", roleId)
	if _, err = rb.Exec(this.db); err != nil {
		return err
	}

	rb = dbs.NewDeleteBuilder()
	rb.UseDialect(this.dialect)
	rb.Table(this.tableRoleMutex)
	rb.Where("ctx = ?", ctx)
	rb.Where("mutex_role_id = ?", roleId)
	if _, err = rb.Exec(this.db); err != nil {
		return err
	}
	return nil
}

// GetMutexRoles 获取互斥关系
func (this *Repository) GetMutexRoles(ctx, roleId int64) (result []*odin.RoleMutex, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.UseDialect(this.dialect)
	sb.Selects("m.ctx", "m.role_id", "m.mutex_role_id", "m.created_on")
	sb.Selects("r.name AS role_name", "r.alias_name AS role_alias_name")
	sb.Selects("rm.name AS mutex_role_name", "rm.alias_name AS mutex_role_alias_name")
	sb.From(this.tableRoleMutex, "AS m")
	sb.LeftJoin(this.tableRole, "AS r ON r.id = m.role_id")
	sb.LeftJoin(this.tableRole, "AS rm ON rm.id = m.mutex_role_id")
	sb.Where("m.ctx = ?", ctx)
	sb.Where("m.role_id = ?", roleId)
	sb.Where("r.ctx = ?", ctx)
	sb.Where("rm.ctx = ?", ctx)
	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetMutexRolesWithIds 获取角色之间的互斥关系
func (this *Repository) GetMutexRolesWithIds(ctx int64, roleIds []int64) (result []*odin.RoleMutex, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.UseDialect(this.dialect)
	sb.Selects("m.ctx", "m.role_id", "m.mutex_role_id", "m.created_on")
	sb.Selects("r.name AS role_name", "r.alias_name AS role_alias_name")
	sb.Selects("rm.name AS mutex_role_name", "rm.alias_name AS mutex_role_alias_name")
	sb.From(this.tableRoleMutex, "AS m")
	sb.LeftJoin(this.tableRole, "AS r ON r.id = m.role_id")
	sb.LeftJoin(this.tableRole, "AS rm ON rm.id = m.mutex_role_id")
	sb.Where("m.ctx = ?", ctx)
	sb.Where(dbs.IN("m.role_id", roleIds))
	sb.Where(dbs.IN("m.mutex_role_id", roleIds))
	sb.Where("r.ctx = ?", ctx)
	sb.Where("rm.ctx = ?", ctx)
	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// CheckRoleMutex 检测两个角色是否互斥
func (this *Repository) CheckRoleMutex(ctx, roleId, mutexRoleId int64) bool {
	var sb = dbs.NewSelectBuilder()
	sb.UseDialect(this.dialect)
	sb.Selects("m.ctx", "m.role_id", "m.mutex_role_id", "m.created_on")
	sb.From(this.tableRoleMutex, "AS m")
	sb.Where("m.ctx = ?", ctx)
	sb.Where("m.role_id = ?", roleId)
	sb.Where("m.mutex_role_id = ?", mutexRoleId)

	var mutex *odin.RoleMutex
	if err := sb.Scan(this.db, &mutex); err != nil {
		return true
	}
	if mutex == nil {
		return false
	}
	return true
}
