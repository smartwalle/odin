package mysql

import (
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/odin"
	"time"
)

func (this *odinRepository) GetPermissionList(ctx, groupId int64, status odin.Status, keywords string) (result []*odin.Permission, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("p.id", "p.group_id", "p.ctx", "p.name", "p.alias_name", "p.status", "p.description", "p.created_on", "p.updated_on")
	sb.From(this.tblPermission, "AS p")
	sb.Where("p.ctx = ?", ctx)
	if groupId > 0 {
		sb.Where("p.group_id = ?", groupId)
	}
	if status != 0 {
		sb.Where("p.status = ?", status)
	}
	if keywords != "" {
		var or = dbs.OR()
		or.Append(dbs.Like("p.name", "%", keywords, "%"))
		or.Append(dbs.Like("p.alias_name", "%", keywords, "%"))
		sb.Where(or)
	}
	sb.OrderBy("p.ctx", "p.id")
	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *odinRepository) GetPermissionListWithIds(ctx int64, permissionIds ...int64) (result []*odin.Permission, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("p.id", "p.group_id", "p.ctx", "p.name", "p.alias_name", "p.status", "p.description", "p.created_on", "p.updated_on")
	sb.From(this.tblPermission, "AS p")
	sb.Where("p.ctx = ?", ctx)
	sb.Where(dbs.IN("p.id", permissionIds))
	sb.OrderBy("p.ctx", "p.id")
	sb.Limit(int64(len(permissionIds)))
	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *odinRepository) GetPermissionListWithNames(ctx int64, names ...string) (result []*odin.Permission, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("p.id", "p.group_id", "p.ctx", "p.name", "p.alias_name", "p.status", "p.description", "p.created_on", "p.updated_on")
	sb.From(this.tblPermission, "AS p")
	sb.Where("p.ctx = ?", ctx)
	sb.Where(dbs.IN("p.name", names))
	sb.OrderBy("p.ctx", "p.id")
	sb.Limit(int64(len(names)))
	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *odinRepository) GetPermissionListWithRoleId(ctx int64, roleId int64) (result []*odin.Permission, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("p.id", "p.group_id", "p.ctx", "p.name", "p.alias_name", "p.status", "p.description", "p.created_on", "p.updated_on")
	sb.From(this.tblRolePermission, "AS rp")
	sb.LeftJoin(this.tblPermission, "AS p ON p.id = rp.permission_id")
	sb.Where("rp.ctx = ?", ctx)
	sb.Where("rp.role_id = ?", roleId)
	sb.Where("p.ctx = ?", ctx)
	sb.OrderBy("p.ctx", "p.id")
	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *odinRepository) getPermission(ctx int64, permissionId int64, name string) (result *odin.Permission, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("p.id", "p.group_id", "p.ctx", "p.name", "p.alias_name", "p.status", "p.description", "p.created_on", "p.updated_on")
	sb.From(this.tblPermission, "AS p")
	sb.Where("p.ctx = ?", ctx)
	if permissionId > 0 {
		sb.Where("p.id = ?", permissionId)
	}
	if name != "" {
		sb.Where("p.name = ?", name)
	}
	sb.Limit(1)
	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *odinRepository) GetPermissionWithId(ctx, permission int64) (result *odin.Permission, err error) {
	return this.getPermission(ctx, permission, "")
}

func (this *odinRepository) GetPermissionWithName(ctx int64, name string) (result *odin.Permission, err error) {
	return this.getPermission(ctx, 0, name)
}

func (this *odinRepository) AddPermission(ctx, groupId int64, name, aliasName, description string, status odin.Status) (result int64, err error) {
	var now = time.Now()
	var ib = dbs.NewInsertBuilder()
	ib.Table(this.tblPermission)
	ib.Columns("group_id", "ctx", "name", "alias_name", "status", "description", "created_on", "updated_on")
	ib.Values(groupId, ctx, name, aliasName, status, description, now, now)
	rResult, err := ib.Exec(this.db)
	if err != nil {
		return 0, err
	}
	if result, err = rResult.LastInsertId(); err != nil {
		return 0, err
	}
	return result, nil
}

func (this *odinRepository) UpdatePermission(ctx, permissionId, groupId int64, aliasName, description string, status odin.Status) (err error) {
	var now = time.Now()
	var ub = dbs.NewUpdateBuilder()
	ub.Table(this.tblPermission)
	ub.SET("group_id", groupId)
	ub.SET("alias_name", aliasName)
	ub.SET("status", status)
	ub.SET("description", description)
	ub.SET("updated_on", now)
	ub.Where("ctx = ? AND id = ?", ctx, permissionId)
	ub.Limit(1)
	_, err = ub.Exec(this.db)
	return err
}

func (this *odinRepository) UpdatePermissionStatus(ctx, permissionId int64, status odin.Status) (err error) {
	var now = time.Now()
	var ub = dbs.NewUpdateBuilder()
	ub.Table(this.tblPermission)
	ub.SET("status", status)
	ub.SET("updated_on", now)
	ub.Where("ctx = ? AND id = ?", ctx, permissionId)
	ub.Limit(1)
	_, err = ub.Exec(this.db)
	return err
}

func (this *odinRepository) GrantPermissionWithIds(ctx, roleId int64, permissionIds ...int64) (err error) {
	if len(permissionIds) == 0 {
		return nil
	}
	var now = time.Now()
	var ib = dbs.NewInsertBuilder()
	ib.Table(this.tblRolePermission)
	ib.Options("IGNORE")
	ib.Columns("ctx", "role_id", "permission_id", "created_on")
	for _, permissionId := range permissionIds {
		ib.Values(ctx, roleId, permissionId, now)
	}
	if _, err = ib.Exec(this.db); err != nil {
		return err
	}
	return nil
}

func (this *odinRepository) RevokePermissionWithIds(ctx, roleId int64, permissionIds ...int64) (err error) {
	if len(permissionIds) == 0 {
		return nil
	}
	var rb = dbs.NewDeleteBuilder()
	rb.Table(this.tblRolePermission)
	rb.Where("ctx = ?", ctx)
	rb.Where("role_id = ?", roleId)
	rb.Where(dbs.IN("permission_id", permissionIds))
	rb.Limit(int64(len(permissionIds)))
	if _, err = rb.Exec(this.db); err != nil {
		return err
	}
	return nil
}

func (this *odinRepository) RevokeAllPermission(ctx, roleId int64) (err error) {
	var rb = dbs.NewDeleteBuilder()
	rb.Table(this.tblRolePermission)
	rb.Where("ctx = ?", ctx)
	rb.Where("role_id = ?", roleId)
	if _, err = rb.Exec(this.db); err != nil {
		return err
	}
	return nil
}
