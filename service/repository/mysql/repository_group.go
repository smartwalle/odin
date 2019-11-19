package mysql

import (
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/odin"
	"time"
)

func (this *odinRepository) GetGroupList(ctx int64, gType odin.GroupType, status odin.Status, name string) (result []*odin.Group, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("g.id", "g.ctx", "g.type", "g.name", "g.status", "g.created_on", "g.updated_on")
	sb.From(this.tblGroup, "AS g")

	sb.Where("g.ctx = ?", ctx)

	if gType > 0 {
		sb.Where("g.type = ?", gType)
	}
	if status > 0 {
		sb.Where("g.status = ?", status)
	}
	if name != "" {
		sb.Where(dbs.Like("g.name", "%", name, "%"))
	}
	sb.OrderBy("g.ctx", "g.id")

	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *odinRepository) GetGroupWithId(ctx, id int64, gType odin.GroupType) (result *odin.Group, err error) {
	return this.getGroup(ctx, id, gType, "")
}

func (this *odinRepository) GetGroupWithName(ctx int64, name string, gType odin.GroupType) (result *odin.Group, err error) {
	return this.getGroup(ctx, 0, gType, name)
}

func (this *odinRepository) getGroup(ctx, id int64, gType odin.GroupType, name string) (result *odin.Group, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("g.id", "g.ctx", "g.type", "g.name", "g.status", "g.created_on", "g.updated_on")
	sb.From(this.tblGroup, "AS g")
	if id > 0 {
		sb.Where("g.id = ?", id)
	}

	sb.Where("g.ctx = ?", ctx)

	if gType > 0 {
		sb.Where("g.type = ?", gType)
	}
	if name != "" {
		sb.Where("g.name = ?", name)
	}
	sb.Limit(1)
	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *odinRepository) AddGroup(ctx int64, gType odin.GroupType, name string, status odin.Status) (result int64, err error) {
	var ib = dbs.NewInsertBuilder()
	ib.Table(this.tblGroup)
	ib.Columns("ctx", "type", "status", "name", "created_on", "updated_on")
	ib.Values(ctx, gType, status, name, time.Now(), time.Now())

	rResult, err := ib.Exec(this.db)
	if err != nil {
		return 0, err
	}
	result, _ = rResult.LastInsertId()
	return result, nil
}

func (this *odinRepository) UpdateGroup(ctx, id int64, name string, status odin.Status) (err error) {
	var ub = dbs.NewUpdateBuilder()
	ub.Table(this.tblGroup)
	ub.SET("name", name)
	ub.SET("status", status)
	ub.SET("updated_on", time.Now())
	ub.Where("id = ?", id)
	ub.Where("ctx = ?", ctx)
	ub.Limit(1)
	_, err = ub.Exec(this.db)
	return err
}

func (this *odinRepository) UpdateGroupStatus(ctx, id int64, status odin.Status) (err error) {
	var ub = dbs.NewUpdateBuilder()
	ub.Table(this.tblGroup)
	ub.SET("status", status)
	ub.SET("updated_on", time.Now())
	ub.Where("id = ?", id)
	ub.Where("ctx = ?", ctx)
	ub.Limit(1)
	_, err = ub.Exec(this.db)
	return err
}

func (this *odinRepository) RemoveGroup(ctx, id int64) (err error) {
	var rb = dbs.NewDeleteBuilder()
	rb.Table(this.tblGroup)
	rb.Where("id = ?", id)
	rb.Where("ctx = ?", ctx)
	rb.Limit(1)
	_, err = rb.Exec(this.db)
	return err
}
