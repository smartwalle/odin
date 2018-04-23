package odin

import (
	"github.com/smartwalle/dbs"
	"time"
)

func (this *manager) getGroupListWithType(gType, status int, name string) (result []*Group, err error) {
	var tx = dbs.MustTx(this.db)
	if result, err = this.getGroupList(tx, gType, status, name); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, err
}

func (this *manager) getGroupList(tx *dbs.Tx, gType, status int, name string) (result []*Group, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("g.id", "g.type", "g.name", "g.status", "g.created_on")
	sb.From(this.groupTable, "AS g")
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
	sb.OrderBy("g.id")

	if err = tx.ExecSelectBuilder(sb, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *manager) getGroupWithId(id int64, gType int) (result *Group, err error) {
	var tx = dbs.MustTx(this.db)
	if result, err = this.getGroup(tx, id, gType, ""); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, err
}

func (this *manager) getGroupWithName(name string, gType int) (result *Group, err error) {
	var tx = dbs.MustTx(this.db)
	if result, err = this.getGroup(tx, 0, gType, name); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, err
}

func (this *manager) addGroup(gType int, name string, status int) (result *Group, err error) {
	var tx = dbs.MustTx(this.db)
	var newGroupId int64 = 0
	if newGroupId, err = this.insertGroup(tx, gType, status, name); err != nil {
		return nil, err
	}
	if result, err = this.getGroup(tx, newGroupId, 0, ""); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, err
}

func (this *manager) insertGroup(tx *dbs.Tx, gType, status int, name string) (id int64, err error) {
	var ib = dbs.NewInsertBuilder()
	ib.Table(this.groupTable)
	ib.Columns("type", "status", "name", "created_on")
	ib.Values(gType, status, name, time.Now())
	if result, err := tx.ExecInsertBuilder(ib); err != nil {
		return 0, err
	} else {
		id, _ = result.LastInsertId()
	}
	return id, err
}

func (this *manager) updateGroup(id int64, name string, status int) (err error) {
	var ub = dbs.NewUpdateBuilder()
	ub.Table(this.groupTable)
	ub.SET("name", name)
	ub.SET("status", status)
	ub.Where("id = ?", id)
	ub.Limit(1)
	_, err = ub.Exec(this.db)
	return err
}

func (this *manager) updateGroupStatus(id int64, status int) (err error) {
	var ub = dbs.NewUpdateBuilder()
	ub.Table(this.groupTable)
	ub.SET("status", status)
	ub.Where("id = ?", id)
	ub.Limit(1)
	_, err = ub.Exec(this.db)
	return err
}

func (this *manager) removeGroup(id int64) (err error) {
	var rb = dbs.NewDeleteBuilder()
	rb.Table(this.groupTable)
	rb.Where("id = ?", id)
	rb.Limit(1)
	_, err = rb.Exec(this.db)
	return err
}

func (this *manager) getGroup(tx *dbs.Tx, id int64, gType int, name string) (result *Group, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("g.id", "g.type", "g.status", "g.name", "g.created_on")
	sb.From(this.groupTable, "AS g")
	if id > 0 {
		sb.Where("g.id = ?", id)
	}
	if gType > 0 {
		sb.Where("g.type = ?", gType)
	}
	if name != "" {
		sb.Where("g.name = ?", name)
	}
	sb.Limit(1)
	if err = tx.ExecSelectBuilder(sb, &result); err != nil {
		return nil, err
	}
	return result, nil
}
