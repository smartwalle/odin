package odin

import (
	"github.com/smartwalle/dbs"
	"time"
)

func (this *manager) getGroupListWithType(ctx int64, gType, status int, name string) (result []*Group, err error) {
	var tx = dbs.MustTx(this.db)
	if result, err = this.getGroupList(tx, ctx, gType, status, name); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, err
}

func (this *manager) getGroupList(tx dbs.TX, ctx int64, gType, status int, name string) (result []*Group, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("g.id", "g.ctx", "g.type", "g.name", "g.status", "g.created_on", "g.updated_on")
	sb.From(this.groupTable, "AS g")

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

	if err = sb.ScanTx(tx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *manager) getGroupWithId(ctx, id int64, gType int) (result *Group, err error) {
	var tx = dbs.MustTx(this.db)
	if result, err = this.getGroup(tx, ctx, id, gType, ""); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, err
}

func (this *manager) getGroupWithName(ctx int64, name string, gType int) (result *Group, err error) {
	var tx = dbs.MustTx(this.db)
	if result, err = this.getGroup(tx, ctx, 0, gType, name); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, err
}

func (this *manager) addGroup(ctx int64, gType int, name string, status int) (result *Group, err error) {
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

func (this *manager) insertGroup(tx dbs.TX, ctx int64, gType, status int, name string) (id int64, err error) {
	var ib = dbs.NewInsertBuilder()
	ib.Table(this.groupTable)
	ib.Columns("ctx", "type", "status", "name", "created_on", "updated_on")
	ib.Values(ctx, gType, status, name, time.Now(), time.Now())
	if result, err := ib.ExecTx(tx); err != nil {
		return 0, err
	} else {
		id, _ = result.LastInsertId()
	}
	return id, err
}

func (this *manager) updateGroup(ctx, id int64, name string, status int) (err error) {
	var ub = dbs.NewUpdateBuilder()
	ub.Table(this.groupTable)
	ub.SET("name", name)
	ub.SET("status", status)
	ub.SET("updated_on", time.Now())
	ub.Where("id = ?", id)
	ub.Where("ctx = ?", ctx)
	ub.Limit(1)
	_, err = ub.Exec(this.db)
	return err
}

func (this *manager) updateGroupStatus(ctx, id int64, status int) (err error) {
	var ub = dbs.NewUpdateBuilder()
	ub.Table(this.groupTable)
	ub.SET("status", status)
	ub.SET("updated_on", time.Now())
	ub.Where("id = ?", id)
	ub.Where("ctx = ?", ctx)
	ub.Limit(1)
	_, err = ub.Exec(this.db)
	return err
}

func (this *manager) removeGroup(ctx, id int64) (err error) {
	var rb = dbs.NewDeleteBuilder()
	rb.Table(this.groupTable)
	rb.Where("id = ?", id)
	rb.Where("ctx = ?", ctx)
	rb.Limit(1)
	_, err = rb.Exec(this.db)
	return err
}

func (this *manager) getGroup(tx dbs.TX, ctx, id int64, gType int, name string) (result *Group, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("g.id", "g.type", "g.status", "g.name", "g.created_on", "g.updated_on")
	sb.From(this.groupTable, "AS g")
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
	if err = sb.ScanTx(tx, &result); err != nil {
		return nil, err
	}
	return result, nil
}
