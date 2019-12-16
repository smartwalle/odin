package mysql

import (
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/odin"
	"time"
)

// AddPreRole 添加授予角色的先决角色条件
func (this *odinRepository) AddPreRole(ctx, roleId int64, preRoleIds []int64) (err error) {
	var now = time.Now()
	var ib = dbs.NewInsertBuilder()
	ib.Options("IGNORE")
	ib.Table(this.tblPreRole)
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
func (this *odinRepository) RemovePreRole(ctx, roleId int64, preRoleIds []int64) (err error) {
	var rb = dbs.NewDeleteBuilder()
	rb.Table(this.tblPreRole)
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
func (this *odinRepository) CleanPreRole(ctx, roleId int64) (err error) {
	var rb = dbs.NewDeleteBuilder()
	rb.Table(this.tblPreRole)
	rb.Where("ctx = ?", ctx)
	rb.Where("role_id = ?", roleId)
	if _, err = rb.Exec(this.db); err != nil {
		return err
	}
	return nil
}

// GetPreRoles 获取先决角色条件
func (this *odinRepository) GetPreRoles(ctx, roleId int64) (result []*odin.PreRole, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("p.ctx", "p.role_id", "p.pro_role_id", "p.created_on")
	sb.Selects("r.name AS role_name", "r.alias_name AS role_alias_name")
	sb.Selects("pr.name AS pro_role_name", "pr.alias_name AS pro_role_alias_name")
	sb.From(this.tblPreRole, "AS p")
	sb.LeftJoin(this.tblRole, "AS r ON r.id = p.role_id")
	sb.LeftJoin(this.tblRole, "AS pr ON pr.id = p.pro_role_id")
	sb.Where("p.ctx = ?", ctx)
	sb.Where("p.role_id = ?", roleId)
	sb.Where("r.ctx = ?", ctx)
	sb.Where("pr.ctx = ?", ctx)
	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, nil
}
