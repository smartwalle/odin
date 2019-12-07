package mysql

import (
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/odin"
	"time"
)

func (this *odinRepository) GetRoles(ctx int64, parentId int64, status odin.Status, keywords, isGrantedToTarget string) (result []*odin.Role, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("r.id", "r.ctx", "r.name", "r.alias_name", "r.status", "r.description", "r.parent_id", "r.left_value", "r.right_value", "r.depth", "r.created_on", "r.updated_on")
	sb.From(this.tblRole, "AS r")
	if isGrantedToTarget != "" {
		sb.Selects("IF(rg.target IS NULL, false, true) AS granted")
		sb.LeftJoin(this.tblGrant, "AS rg ON rg.role_id = r.id AND rg.target = ?", isGrantedToTarget)
	}
	if parentId >= 0 {
		sb.LeftJoin(this.tblRole, "AS rp ON rp.left_value < r.left_value AND rp.right_value > r.right_value")
		sb.Where("rp.ctx = ? AND rp.id = ?", ctx, parentId)
	}
	sb.Where("r.ctx = ?", ctx)
	if status != 0 {
		sb.Where("r.status = ?", status)
	}
	if keywords != "" {
		var or = dbs.OR()
		or.Append(dbs.Like("r.name", "%", keywords, "%"))
		or.Append(dbs.Like("r.alias_name", "%", keywords, "%"))
		sb.Where(or)
	}
	sb.OrderBy("r.ctx", "r.id")

	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (this *odinRepository) GetRolesInTarget(ctx int64, limitedInTarget string, status odin.Status, keywords, isGrantedToTarget string) (result []*odin.Role, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("r.id", "r.ctx", "r.name", "r.alias_name", "r.status", "r.description", "r.parent_id", "r.left_value", "r.right_value", "r.depth", "r.created_on", "r.updated_on")
	sb.From(this.tblRole, "AS r")

	if isGrantedToTarget != "" {
		sb.Selects("IF(rgg.target IS NULL, false, true) AS granted")
		sb.LeftJoin(this.tblGrant, "AS rgg ON rgg.role_id = r.id AND rgg.target = ?", isGrantedToTarget)
	}

	sb.Selects("MAX(IF(rg.role_id = r.id, false, true)) AS can_access")
	sb.LeftJoin(this.tblRole, "AS rp ON rp.left_value <= r.left_value AND rp.right_value >= r.right_value")
	sb.LeftJoin(this.tblGrant, "AS rg ON rg.role_id = rp.id")

	sb.Where("rg.ctx = ?", ctx)
	sb.Where("rg.target = ?", limitedInTarget)
	sb.Where("rp.ctx = ?", ctx)
	if status > 0 {
		sb.Where("rp.status = ?", status)
	}
	sb.Where("r.ctx = ?", ctx)
	if status > 0 {
		sb.Where("r.status = ?", status)
	}
	if keywords != "" {
		var or = dbs.OR()
		or.Append(dbs.Like("r.name", "%", keywords, "%"))
		or.Append(dbs.Like("r.alias_name", "%", keywords, "%"))
		sb.Where(or)
	}
	sb.GroupBy("r.ctx", "r.id")
	sb.OrderBy("r.ctx", "r.id")
	if isGrantedToTarget != "" {
		sb.GroupBy("rgg.target")
		sb.OrderBy("rgg.target")
	}
	if err := sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, err
}

func (this *odinRepository) GetRolesWithIds(ctx int64, roleIds ...int64) (result []*odin.Role, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("r.id", "r.ctx", "r.name", "r.alias_name", "r.status", "r.description", "r.parent_id", "r.left_value", "r.right_value", "r.depth", "r.created_on", "r.updated_on")
	sb.From(this.tblRole, "AS r")
	sb.Where("r.ctx = ?", ctx)
	sb.Where(dbs.IN("r.id", roleIds))
	sb.OrderBy("r.ctx", "r.id")
	sb.Limit(int64(len(roleIds)))
	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *odinRepository) GetRolesWithNames(ctx int64, names ...string) (result []*odin.Role, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("r.id", "r.ctx", "r.name", "r.alias_name", "r.status", "r.description", "r.parent_id", "r.left_value", "r.right_value", "r.depth", "r.created_on", "r.updated_on")
	sb.From(this.tblRole, "AS r")
	sb.Where("r.ctx = ?", ctx)
	sb.Where(dbs.IN("r.name", names))
	sb.OrderBy("r.ctx", "r.id")
	sb.Limit(int64(len(names)))
	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *odinRepository) getRole(ctx int64, roleId int64, name string) (result *odin.Role, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("r.id", "r.ctx", "r.name", "r.alias_name", "r.status", "r.description", "r.parent_id", "r.left_value", "r.right_value", "r.depth", "r.created_on", "r.updated_on")
	sb.From(this.tblRole, "AS r")
	sb.Where("r.ctx = ?", ctx)
	if roleId > 0 {
		sb.Where("r.id = ?", roleId)
	}
	if name != "" {
		sb.Where("r.name = ?", name)
	}
	sb.Limit(1)
	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *odinRepository) GetRoleWithId(ctx, roleId int64) (result *odin.Role, err error) {
	return this.getRole(ctx, roleId, "")
}

func (this *odinRepository) GetRoleWithName(ctx int64, name string) (result *odin.Role, err error) {
	return this.getRole(ctx, 0, name)
}

func (this *odinRepository) AddRole(ctx int64, parent *odin.Role, name, aliasName, description string, status odin.Status) (result int64, err error) {
	if parent == nil {
		if parent, err = this.getMaxRightRole(ctx); err != nil {
			return 0, err
		}
		if parent == nil {
			parent = &odin.Role{}
			parent.Id = 0
			parent.Ctx = ctx
			parent.LeftValue = 0
			parent.RightValue = 0
			parent.Depth = 1
		}
		return this.insertRoleToRoot(parent, name, aliasName, description, status)
	}
	return this.insertRoleToLast(parent, name, aliasName, description, status)
}

func (this *odinRepository) insertRoleToRoot(parent *odin.Role, name, aliasName, description string, status odin.Status) (result int64, err error) {
	return this.insertRole(parent.Ctx, parent.Id, parent.RightValue+1, parent.RightValue+2, parent.Depth, name, aliasName, description, status)
}

func (this *odinRepository) insertRoleToLast(parent *odin.Role, name, aliasName, description string, status odin.Status) (result int64, err error) {
	var ubLeft = dbs.NewUpdateBuilder()
	ubLeft.Table(this.tblRole)
	ubLeft.SET("left_value", dbs.SQL("left_value + 2"))
	ubLeft.SET("updated_on", time.Now())
	ubLeft.Where("ctx = ? AND left_value > ?", parent.Ctx, parent.RightValue)
	if _, err = ubLeft.Exec(this.db); err != nil {
		return 0, err
	}

	var ubRight = dbs.NewUpdateBuilder()
	ubRight.Table(this.tblRole)
	ubRight.SET("right_value", dbs.SQL("right_value + 2"))
	ubRight.SET("updated_on", time.Now())
	ubRight.Where("ctx = ? AND right_value >= ?", parent.Ctx, parent.RightValue)
	if _, err = ubRight.Exec(this.db); err != nil {
		return 0, err
	}
	return this.insertRole(parent.Ctx, parent.Id, parent.RightValue, parent.RightValue+1, parent.Depth+1, name, aliasName, description, status)
}

func (this *odinRepository) insertRole(ctx, parentId int64, leftValue, rightValue, depth int, name, aliasName, description string, status odin.Status) (result int64, err error) {
	var now = time.Now()
	var ib = dbs.NewInsertBuilder()
	ib.Table(this.tblRole)
	ib.Columns("ctx", "name", "alias_name", "status", "description", "parent_id", "left_value", "right_value", "depth", "created_on", "updated_on")
	ib.Values(ctx, name, aliasName, status, description, parentId, leftValue, rightValue, depth, now, now)
	rResult, err := ib.Exec(this.db)
	if err != nil {
		return 0, err
	}
	if result, err = rResult.LastInsertId(); err != nil {
		return 0, err
	}
	return result, nil
}

func (this *odinRepository) getMaxRightRole(ctx int64) (result *odin.Role, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("r.id", "r.ctx", "r.name", "r.alias_name", "r.status", "r.description", "r.parent_id", "r.left_value", "r.right_value", "r.depth", "r.created_on", "r.updated_on")
	sb.From(this.tblRole, "AS r")
	sb.Where("r.ctx = ?", ctx)
	sb.OrderBy("r.right_value DESC")
	sb.Limit(1)
	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *odinRepository) UpdateRole(ctx, roleId int64, aliasName, description string, status odin.Status) (err error) {
	var now = time.Now()
	var ub = dbs.NewUpdateBuilder()
	ub.Table(this.tblRole)
	ub.SET("alias_name", aliasName)
	ub.SET("status", status)
	ub.SET("description", description)
	ub.SET("updated_on", now)
	ub.Where("ctx = ? AND id = ?", ctx, roleId)
	ub.Limit(1)
	_, err = ub.Exec(this.db)
	return err
}

func (this *odinRepository) UpdateRoleStatus(ctx, roleId int64, status odin.Status) (err error) {
	var now = time.Now()
	var ub = dbs.NewUpdateBuilder()
	ub.Table(this.tblRole)
	ub.SET("status", status)
	ub.SET("updated_on", now)
	ub.Where("ctx = ? AND id = ?", ctx, roleId)
	ub.Limit(1)
	_, err = ub.Exec(this.db)
	return err
}

func (this *odinRepository) GetGrantedRoles(ctx int64, target string, withChildren bool) (result []*odin.Role, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("r.id", "r.ctx", "r.name", "r.alias_name", "r.status", "r.description", "r.parent_id", "r.left_value", "r.right_value", "r.depth", "r.created_on", "r.updated_on")
	sb.Selects("MAX(IF(rg.role_id <> r.id, false, true)) AS granted")
	sb.Selects("MAX(IF(rg.role_id = r.id, false, true)) AS can_access")
	sb.From(this.tblRole, "AS r")
	if withChildren {
		sb.LeftJoin(this.tblRole, "AS rp ON rp.left_value <= r.left_value AND rp.right_value >= r.right_value")
	} else {
		sb.LeftJoin(this.tblRole, "AS rp ON rp.left_value = r.left_value AND rp.right_value = r.right_value")
	}
	sb.LeftJoin(this.tblGrant, "AS rg ON rg.role_id = rp.id")
	sb.Where("rg.ctx = ?", ctx)
	sb.Where("rg.target = ?", target)
	sb.Where("rp.ctx = ?", ctx)
	sb.Where("rp.status = ?", odin.Enable)
	sb.Where("r.ctx = ?", ctx)
	sb.Where("r.status = ?", odin.Enable)
	sb.GroupBy("r.ctx", "r.id")
	sb.OrderBy("r.ctx", "r.id")
	if err := sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, err
}

func (this *odinRepository) GrantRoleWithIds(ctx int64, target string, roleIds ...int64) (err error) {
	if len(roleIds) == 0 {
		return nil
	}
	var now = time.Now()
	var ib = dbs.NewInsertBuilder()
	ib.Table(this.tblGrant)
	ib.Options("IGNORE")
	ib.Columns("ctx", "role_id", "target", "created_on")
	for _, rId := range roleIds {
		ib.Values(ctx, rId, target, now)
	}
	if _, err = ib.Exec(this.db); err != nil {
		return err
	}
	return nil
}

func (this *odinRepository) RevokeRoleWithIds(ctx int64, target string, roleIds ...int64) (err error) {
	if len(roleIds) == 0 {
		return nil
	}
	var rb = dbs.NewDeleteBuilder()
	rb.Table(this.tblGrant)
	rb.Where("ctx = ?", ctx)
	rb.Where("target = ?", target)
	rb.Where(dbs.IN("role_id", roleIds))
	rb.Limit(int64(len(roleIds)))
	if _, err = rb.Exec(this.db); err != nil {
		return err
	}
	return nil
}

func (this *odinRepository) RevokeAllRole(ctx int64, target string) (err error) {
	var rb = dbs.NewDeleteBuilder()
	rb.Table(this.tblGrant)
	rb.Where("ctx = ?", ctx)
	rb.Where("target = ?", target)
	if _, err = rb.Exec(this.db); err != nil {
		return err
	}
	return nil
}
