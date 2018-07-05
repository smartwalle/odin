package odin

import (
	"github.com/smartwalle/dbs"
	"time"
)

func (this *manager) getPermissionTree(ctxId, roleId int64, status int, name string) (result []*Group, err error) {
	var tx = dbs.MustTx(this.db)

	if result, err = this.getGroupList(tx, ctxId, K_GROUP_TYPE_PERMISSION, status, name); err != nil {
		return nil, err
	}

	var gMap = make(map[int64]*Group)
	var gIdList []int64
	for _, group := range result {
		gMap[group.Id] = group
		gIdList = append(gIdList, group.Id)
	}

	pList, err := this.getPermissionListWithGroupIdList(tx, ctxId, roleId, gIdList, status, "")
	if err != nil {
		return nil, err
	}

	for _, p := range pList {
		var group = gMap[p.GroupId]
		if group != nil {
			group.PermissionList = append(group.PermissionList, p)
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *manager) getPermissionList(ctxId int64, groupIdList []int64, status int, keyword string) (result []*Permission, err error) {
	var tx = dbs.MustTx(this.db)
	if result, err = this.getPermissionListWithGroupIdList(tx, ctxId, 0, groupIdList, status, keyword); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *manager) getPermissionListWithGroupIdList(tx dbs.TX, ctxId, roleId int64, groupIdList []int64, status int, keyword string) (result []*Permission, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("p.id", "p.ctx_id", "p.group_id", "p.name", "p.identifier", "p.status", "p.created_on")
	sb.From(this.permissionTable, "AS p")
	if roleId > 0 {
		sb.Selects("IF(rp.role_id IS NULL, false , true) AS granted")
		sb.LeftJoin(this.rolePermissionTable, "AS rp ON rp.permission_id = p.id AND rp.role_id = ?", roleId)
	}

	sb.Where("(p.ctx_id = ? OR p.ctx_id = ?)", 0, ctxId)

	if len(groupIdList) > 0 {
		sb.Where(dbs.IN("p.group_id", groupIdList))
	}
	if status > 0 {
		sb.Where("p.status = ?", status)
	}
	if keyword != "" {
		var k = "%" + keyword + "%"
		sb.Where("(p.name LIKE ? OR p.identifier LIKE ?)", k, k)
	}
	sb.OrderBy("p.ctx_id", "p.id")
	if err = sb.ScanTx(tx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *manager) getPermissionWithIdList(ctxId int64, idList []int64) (result []*Permission, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("p.id", "p.ctx_id", "p.group_id", "p.name", "p.identifier", "p.status", "p.created_on")
	sb.From(this.permissionTable, "AS p")
	sb.Where("(p.ctx_id = ? OR p.ctx_id = ?)", 0, ctxId)
	if len(idList) > 0 {
		sb.Where(dbs.IN("p.id", idList))
	}
	sb.Limit(uint64(len(idList)))

	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *manager) getPermissionWithId(ctxId, id int64) (result *Permission, err error) {
	var tx = dbs.MustTx(this.db)
	if result, err = this.getPermission(tx, ctxId, id, "", ""); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, err
}

func (this *manager) getPermissionWithName(ctxId int64, name string) (result *Permission, err error) {
	var tx = dbs.MustTx(this.db)
	if result, err = this.getPermission(tx, ctxId, 0, name, ""); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, err
}

func (this *manager) getPermissionWithIdentifier(ctxId int64, identifier string) (result *Permission, err error) {
	var tx = dbs.MustTx(this.db)
	if result, err = this.getPermission(tx, ctxId, 0, "", identifier); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, err
}

func (this *manager) addPermission(ctxId int64, groupId int64, name, identifier string, status int) (result *Permission, err error) {
	var tx = dbs.MustTx(this.db)
	var newPermissionId int64 = 0
	if newPermissionId, err = this.insertPermission(tx, ctxId, groupId, status, name, identifier); err != nil {
		return nil, err
	}
	if result, err = this.getPermission(tx, ctxId, newPermissionId, "", ""); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, err
}

func (this *manager) insertPermission(tx dbs.TX, ctxId, groupId int64, status int, name, identifier string) (id int64, err error) {
	var ib = dbs.NewInsertBuilder()
	ib.Table(this.permissionTable)
	ib.Columns("ctx_id", "group_id", "status", "name", "identifier", "created_on", "updated_on")
	ib.Values(ctxId, groupId, status, name, identifier, time.Now(), time.Now())
	if result, err := ib.ExecTx(tx); err != nil {
		return 0, err
	} else {
		id, _ = result.LastInsertId()
	}
	return id, err
}

func (this *manager) updatePermission(ctxId, id, groupId int64, name, identifier string, status int) (err error) {
	var ub = dbs.NewUpdateBuilder()
	ub.Table(this.permissionTable)
	ub.SET("group_id", groupId)
	ub.SET("name", name)
	ub.SET("identifier", identifier)
	ub.SET("status", status)
	ub.SET("updated_on", time.Now())
	ub.Where("id = ?", id)
	ub.Where("ctx_id = ?", ctxId)
	ub.Limit(1)
	_, err = ub.Exec(this.db)
	return err
}

func (this *manager) updatePermissionStatus(ctxId, id int64, status int) (err error) {
	var ub = dbs.NewUpdateBuilder()
	ub.Table(this.permissionTable)
	ub.SET("status", status)
	ub.SET("updated_on", time.Now())
	ub.Where("id = ?", id)
	ub.Where("ctx_id = ?", ctxId)
	ub.Limit(1)
	_, err = ub.Exec(this.db)
	return err
}

func (this *manager) getPermission(tx dbs.TX, ctxId, id int64, name, identifier string) (result *Permission, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("p.id", "p.ctx_id", "p.group_id", "p.name", "p.identifier", "p.status", "p.created_on", "p.updated_on")
	sb.From(this.permissionTable, "AS p")
	if id > 0 {
		sb.Where("p.id = ?", id)
	}
	sb.Where("(p.ctx_id = ? OR p.ctx_id = ?)", 0, ctxId)
	if name != "" {
		sb.Where("p.name = ?", name)
	}
	if identifier != "" {
		sb.Where("p.identifier = ?", identifier)
	}
	sb.Limit(1)
	if err = sb.ScanTx(tx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *manager) getPermissionListWithRoleId(ctxId, roleId int64) (result []*Permission, err error) {
	var tx = dbs.MustTx(this.db)
	if result, err = this.getPermissionListWithRole(tx, ctxId, roleId); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, err
}

func (this *manager) getPermissionListWithRole(tx dbs.TX, ctxId, roleId int64) (result []*Permission, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("p.id", "p.ctx_id", "p.group_id", "p.name", "p.identifier", "p.status", "p.created_on", "p.updated_on")
	sb.Selects("IF(rp.role_id IS NULL, false, true) AS granted")
	sb.From(this.permissionTable, "AS p")
	sb.LeftJoin(this.rolePermissionTable, "AS rp ON rp.permission_id = p.id")
	sb.Where("(p.ctx_id = ? OR p.ctx_id = ?)", 0, ctxId)
	sb.Where("rp.role_id = ?", roleId)
	if err = sb.ScanTx(tx, &result); err != nil {
		return nil, err
	}
	return result, err
}

func (this *manager) getGrantedPermissionList(ctxId int64, objectId string) (result []*Permission, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("p.id", "p.ctx_id", "p.group_id", "p.name", "p.identifier", "p.status", "p.created_on", "p.updated_on")
	sb.Selects("IF(rp.role_id IS NULL, false, true) AS granted")
	sb.From(this.permissionTable, "AS p")
	sb.LeftJoin(this.rolePermissionTable, "AS rp ON rp.permission_id = p.id")
	sb.LeftJoin(this.roleGrantTable, "AS rg ON rg.role_id = rp.role_id")
	sb.Where("rg.object_id = ?", objectId)
	sb.Where("p.status = ?", K_STATUS_ENABLE)
	sb.Where("(p.ctx_id = ? OR p.ctx_id = ?)", 0, ctxId)
	sb.GroupBy("p.id")
	if err := sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, err
}
