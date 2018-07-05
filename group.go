package odin

import (
	"github.com/smartwalle/dbs"
	"time"
)

func (this *manager) getGroupListWithType(ctxId int64, gType, status int, name string) (result []*Group, err error) {
	var tx = dbs.MustTx(this.db)
	if result, err = this.getGroupList(tx, ctxId, gType, status, name); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, err
}

func (this *manager) getGroupList(tx dbs.TX, ctxId int64, gType, status int, name string) (result []*Group, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("g.id", "g.ctx_id", "g.type", "g.name", "g.status", "g.created_on", "g.updated_on")
	sb.From(this.groupTable, "AS g")

	sb.Where("(g.ctx_id = ? OR g.ctx_id = ?)", 0, ctxId)

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
	sb.OrderBy("g.ctx_id", "g.id")

	if err = sb.ScanTx(tx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *manager) getGroupWithId(ctxId, id int64, gType int) (result *Group, err error) {
	var tx = dbs.MustTx(this.db)
	if result, err = this.getGroup(tx, ctxId, id, gType, ""); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, err
}

func (this *manager) getGroupWithName(ctxId int64, name string, gType int) (result *Group, err error) {
	var tx = dbs.MustTx(this.db)
	if result, err = this.getGroup(tx, ctxId, 0, gType, name); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, err
}

func (this *manager) addGroup(ctxId int64, gType int, name string, status int) (result *Group, err error) {
	var tx = dbs.MustTx(this.db)
	var newGroupId int64 = 0
	if newGroupId, err = this.insertGroup(tx, ctxId, gType, status, name); err != nil {
		return nil, err
	}
	if result, err = this.getGroup(tx, ctxId, newGroupId, 0, ""); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, err
}

func (this *manager) insertGroup(tx dbs.TX, ctxId int64, gType, status int, name string) (id int64, err error) {
	var ib = dbs.NewInsertBuilder()
	ib.Table(this.groupTable)
	ib.Columns("ctx_id", "type", "status", "name", "created_on", "updated_on")
	ib.Values(ctxId, gType, status, name, time.Now(), time.Now())
	if result, err := ib.ExecTx(tx); err != nil {
		return 0, err
	} else {
		id, _ = result.LastInsertId()
	}
	return id, err
}

func (this *manager) updateGroup(ctxId, id int64, name string, status int) (err error) {
	var ub = dbs.NewUpdateBuilder()
	ub.Table(this.groupTable)
	ub.SET("name", name)
	ub.SET("status", status)
	ub.SET("updated_on", time.Now())
	ub.Where("id = ?", id)
	ub.Where("ctx_id = ?", ctxId)
	ub.Limit(1)
	_, err = ub.Exec(this.db)
	return err
}

func (this *manager) updateGroupStatus(ctxId, id int64, status int) (err error) {
	var ub = dbs.NewUpdateBuilder()
	ub.Table(this.groupTable)
	ub.SET("status", status)
	ub.SET("updated_on", time.Now())
	ub.Where("id = ?", id)
	ub.Where("ctx_id = ?", ctxId)
	ub.Limit(1)
	_, err = ub.Exec(this.db)
	return err
}

func (this *manager) removeGroup(ctxId, id int64) (err error) {
	var rb = dbs.NewDeleteBuilder()
	rb.Table(this.groupTable)
	rb.Where("id = ?", id)
	rb.Where("ctx_id = ?", ctxId)
	rb.Limit(1)
	_, err = rb.Exec(this.db)
	return err
}

func (this *manager) getGroup(tx dbs.TX, ctxId, id int64, gType int, name string) (result *Group, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("g.id", "g.type", "g.status", "g.name", "g.created_on", "g.updated_on")
	sb.From(this.groupTable, "AS g")
	if id > 0 {
		sb.Where("g.id = ?", id)
	}

	sb.Where("(g.ctx_id = ? OR g.ctx_id = ?)", 0, ctxId)

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
