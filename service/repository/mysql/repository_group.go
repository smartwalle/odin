package mysql

import (
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/odin"
	"time"
)

func (this *odinRepository) GetGroupListWithType(ctx int64, gType, status int, name string) (result []*odin.Group, err error) {
	var tx = dbs.MustTx(this.db)
	if result, err = this.getGroupList(tx, ctx, gType, status, name); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, err
}

func (this *odinRepository) getGroupList(tx dbs.TX, ctx int64, gType, status int, name string) (result []*odin.Group, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("g.id", "g.ctx", "g.type", "g.name", "g.status", "g.created_on", "g.updated_on")
	sb.From(this.tblGroup, "AS g")

	sb.Where("(g.ctx = ? OR g.ctx = ?)", 0, ctx)

	if gType > 0 {
		sb.Where("g.type = ?", gType)
	}
	if status > 0 {
		sb.Where("g.status = ?", status)
	}
	if name != "" {
		var keyword = "%" + name + "%"
		sb.Where("g.name LIKE ?", keyword)
	}
	sb.OrderBy("g.ctx", "g.id")

	if err = sb.Scan(tx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *odinRepository) GetGroupWithId(ctx, id int64, gType int) (result *odin.Group, err error) {
	var tx = dbs.MustTx(this.db)
	if result, err = this.getGroup(tx, ctx, id, gType, ""); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, err
}

func (this *odinRepository) GetGroupWithName(ctx int64, name string, gType int) (result *odin.Group, err error) {
	var tx = dbs.MustTx(this.db)
	if result, err = this.getGroup(tx, ctx, 0, gType, name); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, err
}

func (this *odinRepository) AddGroup(ctx int64, gType int, name string, status int) (result *odin.Group, err error) {
	var tx = dbs.MustTx(this.db)
	var newGroupId int64 = 0
	if newGroupId, err = this.insertGroup(tx, ctx, gType, status, name); err != nil {
		return nil, err
	}
	if result, err = this.getGroup(tx, ctx, newGroupId, 0, ""); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, err
}

func (this *odinRepository) insertGroup(tx dbs.TX, ctx int64, gType, status int, name string) (id int64, err error) {
	var ib = dbs.NewInsertBuilder()
	ib.Table(this.tblGroup)
	ib.Columns("ctx", "type", "status", "name", "created_on", "updated_on")
	ib.Values(ctx, gType, status, name, time.Now(), time.Now())
	if result, err := ib.Exec(tx); err != nil {
		return 0, err
	} else {
		id, _ = result.LastInsertId()
	}
	return id, err
}

func (this *odinRepository) UpdateGroup(ctx, id int64, name string, status int) (err error) {
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

func (this *odinRepository) UpdateGroupStatus(ctx, id int64, status int) (err error) {
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

func (this *odinRepository) getGroup(tx dbs.TX, ctx, id int64, gType int, name string) (result *odin.Group, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("g.id", "g.type", "g.status", "g.name", "g.created_on", "g.updated_on")
	sb.From(this.tblGroup, "AS g")
	if id > 0 {
		sb.Where("g.id = ?", id)
	}

	sb.Where("(g.ctx = ? OR g.ctx = ?)", 0, ctx)

	if gType > 0 {
		sb.Where("g.type = ?", gType)
	}
	if name != "" {
		sb.Where("g.name = ?", name)
	}
	sb.Limit(1)
	if err = sb.Scan(tx, &result); err != nil {
		return nil, err
	}
	return result, nil
}
