package mysql

import (
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/odin"
	"time"
)

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
