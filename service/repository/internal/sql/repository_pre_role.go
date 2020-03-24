package sql

import (
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/odin"
	"time"
)

// AddPreRole 添加授予角色的先决角色条件
func (this *Repository) AddPreRole(ctx, roleId int64, preRoleIds []int64) (err error) {
	var now = time.Now()
	var ib = dbs.NewInsertBuilder()
	ib.UseDialect(this.dialect)
	ib.Options("IGNORE")
	ib.Table(this.tablePreRole)
	ib.Columns("ctx", "role_id", "pre_role_id", "created_on")
	for _, preRoleId := range preRoleIds {
		ib.Values(ctx, roleId, preRoleId, now)
	}
	if _, err = ib.Exec(this.db); err != nil {
		return err
	}
	return nil
}

// RemovePreRole 删除授予角色的先决角色条件
func (this *Repository) RemovePreRole(ctx, roleId int64, preRoleIds []int64) (err error) {
	var rb = dbs.NewDeleteBuilder()
	rb.UseDialect(this.dialect)
	rb.Table(this.tablePreRole)
	rb.Where("ctx = ?", ctx)
	rb.Where("role_id = ?", roleId)
	rb.Where(dbs.IN("pre_role_id", preRoleIds))
	rb.Limit(int64(len(preRoleIds)))
	if _, err = rb.Exec(this.db); err != nil {
		return err
	}
	return nil
}

// CleanPreRole 清除授予角色的先决角色条件
func (this *Repository) CleanPreRole(ctx, roleId int64) (err error) {
	var rb = dbs.NewDeleteBuilder()
	rb.UseDialect(this.dialect)
	rb.Table(this.tablePreRole)
	rb.Where("ctx = ?", ctx)
	rb.Where("role_id = ?", roleId)
	if _, err = rb.Exec(this.db); err != nil {
		return err
	}
	return nil
}

// GetPreRoles 获取授予角色的先决角色条件
func (this *Repository) GetPreRoles(ctx, roleId int64) (result []*odin.PreRole, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.UseDialect(this.dialect)
	sb.Selects("p.ctx", "p.role_id", "p.pre_role_id", "p.created_on")
	sb.Selects("r.name AS role_name", "r.alias_name AS role_alias_name")
	sb.Selects("pr.name AS pre_role_name", "pr.alias_name AS pre_role_alias_name")
	sb.From(this.tablePreRole, "AS p")
	sb.LeftJoin(this.tableRole, "AS r ON r.id = p.role_id")
	sb.LeftJoin(this.tableRole, "AS pr ON pr.id = p.pre_role_id")
	sb.Where("p.ctx = ?", ctx)
	sb.Where("p.role_id = ?", roleId)
	sb.Where("r.ctx = ?", ctx)
	sb.Where("pr.ctx = ?", ctx)
	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetPreRolesWithIds 获取授予角色列表的先决角色条件
func (this *Repository) GetPreRolesWithIds(ctx int64, roleIds []int64) (result []*odin.PreRole, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.UseDialect(this.dialect)
	sb.Selects("p.ctx", "p.role_id", "p.pre_role_id", "p.created_on")
	sb.Selects("r.name AS role_name", "r.alias_name AS role_alias_name")
	sb.Selects("pr.name AS pre_role_name", "pr.alias_name AS pre_role_alias_name")
	sb.From(this.tablePreRole, "AS p")
	sb.LeftJoin(this.tableRole, "AS r ON r.id = p.role_id")
	sb.LeftJoin(this.tableRole, "AS pr ON pr.id = p.pre_role_id")
	sb.Where("p.ctx = ?", ctx)
	sb.Where(dbs.IN("p.role_id", roleIds))
	sb.Where("r.ctx = ?", ctx)
	sb.Where("pr.ctx = ?", ctx)
	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, nil
}
