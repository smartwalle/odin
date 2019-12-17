package mysql

import (
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/odin"
	"time"
)

func (this *odinRepository) GetPermissions(ctx int64, status odin.Status, keywords string, groupIds []int64, limitedInRole, isGrantedToRole int64) (result []*odin.Permission, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("p.id", "p.group_id", "p.ctx", "p.name", "p.alias_name", "p.status", "p.description", "p.created_on", "p.updated_on")
	sb.From(this.tblPermission, "AS p")
	if isGrantedToRole > 0 {
		sb.Selects("IF(rp.role_id IS NULL, false , true) AS granted")
		sb.LeftJoin(this.tblRolePermission, "AS rp ON rp.permission_id = p.id AND rp.role_id = ?", isGrantedToRole)
	}
	if limitedInRole > 0 {
		sb.LeftJoin(this.tblRolePermission, "AS rpl ON rpl.permission_id = p.id")
		sb.Where("rpl.ctx = ?", ctx)
		sb.Where("rpl.role_id = ?", limitedInRole)
	}
	sb.Where("p.ctx = ?", ctx)
	if len(groupIds) > 0 {
		sb.Where(dbs.IN("p.group_id", groupIds))
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

func (this *odinRepository) GetPermissionsWithIds(ctx int64, permissionIds ...int64) (result []*odin.Permission, err error) {
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

func (this *odinRepository) GetPermissionsWithNames(ctx int64, names ...string) (result []*odin.Permission, err error) {
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

func (this *odinRepository) GetPermissionsWithRoleId(ctx int64, roleId int64) (result []*odin.Permission, err error) {
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

func (this *odinRepository) GrantPermissionWithIds(ctx, roleId int64, permissionIds []int64) (err error) {
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

func (this *odinRepository) RevokePermissionWithIds(ctx, roleId int64, permissionIds []int64) (err error) {
	if len(permissionIds) == 0 {
		return nil
	}
	// TODO 优化
	var rb = dbs.NewDeleteBuilder()
	rb.Alias("p")
	rb.Table(this.tblRolePermission, "AS p")
	rb.LeftJoin(this.tblRole, "AS r ON r.id = p.role_id")
	rb.LeftJoin(this.tblRole, "AS rp ON rp.left_value <= r.left_value AND r.right_value >= r.right_value")
	rb.Where("p.ctx = ?", ctx)
	rb.Where("rp.ctx = ?", ctx)
	rb.Where("rp.id = ?", roleId)
	rb.Where("r.ctx = ?", ctx)
	rb.Where(dbs.IN("permission_id", permissionIds))
	if _, err = rb.Exec(this.db); err != nil {
		return err
	}
	return nil
}

func (this *odinRepository) RevokeAllPermission(ctx, roleId int64) (err error) {
	// TODO 优化
	var rb = dbs.NewDeleteBuilder()
	rb.Alias("p")
	rb.Table(this.tblRolePermission, "AS p")
	rb.LeftJoin(this.tblRole, "AS r ON r.id = p.role_id")
	rb.LeftJoin(this.tblRole, "AS rp ON rp.left_value <= r.left_value AND r.right_value >= r.right_value")
	rb.Where("p.ctx = ?", ctx)
	rb.Where("rp.ctx = ?", ctx)
	rb.Where("rp.id = ?", roleId)
	rb.Where("r.ctx = ?", ctx)
	if _, err = rb.Exec(this.db); err != nil {
		return err
	}
	return nil
}

func (this *odinRepository) GetGrantedPermissions(ctx int64, target string) (result []*odin.Permission, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("p.id", "p.group_id", "p.ctx", "p.name", "p.alias_name", "p.status", "p.description", "p.created_on", "p.updated_on")
	sb.Selects("IF(p.id IS NULL, false, true) AS granted")
	sb.From(this.tblRolePermission, "AS rp")
	sb.LeftJoin(this.tblPermission, "AS p ON p.id = rp.permission_id")
	sb.LeftJoin(this.tblGrant, "AS g ON g.role_id = rp.role_id")
	sb.Where("rp.ctx = ?", ctx)
	sb.Where("g.ctx = ? AND g.target = ?", ctx, target)
	sb.Where("p.ctx = ? AND p.status = ?", ctx, odin.Enable)
	sb.GroupBy("p.id")
	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// AddPrePermission 添加授予权限的先决权限条件
func (this *odinRepository) AddPrePermission(ctx, permissionId int64, prePermissionIds []int64) (err error) {
	var now = time.Now()
	var ib = dbs.NewInsertBuilder()
	ib.Options("IGNORE")
	ib.Table(this.tblPrePermission)
	ib.Columns("ctx", "permission_id", "pre_permission_id", "auto_grant", "created_on")
	for _, prePermissionId := range prePermissionIds {
		ib.Values(ctx, permissionId, prePermissionId, false, now)
	}
	if _, err = ib.Exec(this.db); err != nil {
		return err
	}
	return nil
}

// RemovePrePermission 删除授予权限的先决权限条件
func (this *odinRepository) RemovePrePermission(ctx, permissionId int64, prePermissionIds []int64) (err error) {
	var rb = dbs.NewDeleteBuilder()
	rb.Table(this.tblPrePermission)
	rb.Where("ctx = ?", ctx)
	rb.Where("permission_id = ?", permissionId)
	rb.Where(dbs.IN("pre_permission_id", prePermissionIds))
	rb.Limit(int64(len(prePermissionIds)))
	if _, err = rb.Exec(this.db); err != nil {
		return err
	}
	return nil
}

// CleanPrePermission 清除授予权限的先决权限条件
func (this *odinRepository) CleanPrePermission(ctx, permissionId int64) (err error) {
	var rb = dbs.NewDeleteBuilder()
	rb.Table(this.tblPrePermission)
	rb.Where("ctx = ?", ctx)
	rb.Where("permission_id = ?", permissionId)
	if _, err = rb.Exec(this.db); err != nil {
		return err
	}
	return nil
}

// GetPrePermissions 获取授予权限的先决权限条件
func (this *odinRepository) GetPrePermissions(ctx, permissionId int64) (result []*odin.PrePermission, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("pp.ctx", "pp.permission_id", "pp.pre_permission_id", "pp.auto_grant", "pp.created_on")
	sb.Selects("p.name AS permission_name", "p.alias_name AS permission_alias_name")
	sb.Selects("ppp.name AS pre_permission_name", "ppp.alias_name AS pre_permission_alias_name")
	sb.From(this.tblPrePermission, "AS pp")
	sb.LeftJoin(this.tblPermission, "AS p ON p.id = pp.permission_id")
	sb.LeftJoin(this.tblPermission, "AS ppp ON ppp.id = pp.pre_permission_id")
	sb.Where("p.ctx = ?", ctx)
	sb.Where("pp.permission_id = ?", permissionId)
	sb.Where("pp.ctx = ?", ctx)
	sb.Where("ppp.ctx = ?", ctx)
	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetPrePermissionsWithIds 获取授予权限列表的先决权限条件
func (this *odinRepository) GetPrePermissionsWithIds(ctx int64, permissionIds []int64) (result []*odin.PrePermission, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("pp.ctx", "pp.permission_id", "pp.pre_permission_id", "pp.auto_grant", "pp.created_on")
	sb.Selects("p.name AS permission_name", "p.alias_name AS permission_alias_name")
	sb.Selects("ppp.name AS pre_permission_name", "ppp.alias_name AS pre_permission_alias_name")
	sb.From(this.tblPrePermission, "AS pp")
	sb.LeftJoin(this.tblPermission, "AS p ON p.id = pp.permission_id")
	sb.LeftJoin(this.tblPermission, "AS ppp ON ppp.id = pp.pre_permission_id")
	sb.Where("p.ctx = ?", ctx)
	sb.Where(dbs.IN("pp.permission_id", permissionIds))
	sb.Where("pp.ctx = ?", ctx)
	sb.Where("ppp.ctx = ?", ctx)
	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, nil
}
