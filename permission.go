package odin

import (
	"github.com/smartwalle/dbs"
	"time"
)

func (this *manager) getPermissionList(groupId int64, status int, keyword string) (result []*Permission, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("p.id", "p.group_id", "p.name", "p.identifier", "p.status", "p.created_on")
	sb.From(this.permissionTable, "AS p")
	if groupId > 0 {
		sb.Where("p.group_id = ?", groupId)
	}
	if status > 0 {
		sb.Where("p.status = ?", status)
	}
	if keyword != "" {
		var k = "%" + keyword + "%"
		sb.Where("(p.name LIKE ? OR p.identifier LIKE ?)", k, k)
	}
	sb.OrderBy("p.id")

	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *manager) getPermissionWithId(id int64) (result *Permission, err error) {
	var tx = dbs.MustTx(this.db)
	if result, err = this.getPermission(tx, id, "", ""); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, err
}

func (this *manager) getPermissionWithName(name string) (result *Permission, err error) {
	var tx = dbs.MustTx(this.db)
	if result, err = this.getPermission(tx, 0, name, ""); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, err
}

func (this *manager) getPermissionWithIdentifier(identifier string) (result *Permission, err error) {
	var tx = dbs.MustTx(this.db)
	if result, err = this.getPermission(tx, 0, "", identifier); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, err
}

func (this *manager) addPermission(groupId int64, name, identifier string, status int) (result *Permission, err error) {
	var tx = dbs.MustTx(this.db)
	var newPermissionId int64 = 0
	if newPermissionId, err = this.insertPermission(tx, groupId, status, name, identifier); err != nil {
		return nil, err
	}
	if result, err = this.getPermission(tx, newPermissionId, "", ""); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, err
}

func (this *manager) insertPermission(tx *dbs.Tx, groupId int64, status int, name, identifier string) (id int64, err error) {
	var ib = dbs.NewInsertBuilder()
	ib.Table(this.permissionTable)
	ib.Columns("group_id", "status", "name", "identifier", "created_on")
	ib.Values(groupId, status, name, identifier, time.Now())
	if result, err := tx.ExecInsertBuilder(ib); err != nil {
		return 0, err
	} else {
		id, _ = result.LastInsertId()
	}
	return id, err
}

func (this *manager) updatePermission(id, groupId int64, name, identifier string, status int) (err error) {
	var ub = dbs.NewUpdateBuilder()
	ub.Table(this.permissionTable)
	ub.SET("group_id", groupId)
	ub.SET("name", name)
	ub.SET("identifier", identifier)
	ub.SET("status", status)
	ub.Where("id = ?", id)
	ub.Limit(1)
	_, err = ub.Exec(this.db)
	return err
}

func (this *manager) updatePermissionStatus(id int64, status int) (err error) {
	var ub = dbs.NewUpdateBuilder()
	ub.Table(this.permissionTable)
	ub.SET("status", status)
	ub.Where("id = ?", id)
	ub.Limit(1)
	_, err = ub.Exec(this.db)
	return err
}

func (this *manager) getPermission(tx *dbs.Tx, id int64, name, identifier string) (result *Permission, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("p.id", "p.group_id", "p.name", "p.identifier", "p.status", "p.created_on")
	sb.From(this.permissionTable, "AS p")
	if id > 0 {
		sb.Where("p.id = ?", id)
	}
	if name != "" {
		sb.Where("p.name = ?", name)
	}
	if identifier != "" {
		sb.Where("p.identifier = ?", identifier)
	}
	sb.Limit(1)
	if err = tx.ExecSelectBuilder(sb, &result); err != nil {
		return nil, err
	}
	return result, nil
}
