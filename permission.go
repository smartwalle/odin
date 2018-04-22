package odin

import (
	"github.com/smartwalle/dbs"
	"time"
)

func (this *manager) getPermissionTree(roleId int64, status int, name string) (result []*Group, err error) {
	var tx = dbs.MustTx(this.db)

	if result, err = this.getGroupList(tx, K_GROUP_TYPE_PERMISSION, status, name); err != nil {
		return nil, err
	}

	var gMap = make(map[int64]*Group)
	var gIdList []int64
	for _, group := range result {
		gMap[group.Id] = group
		gIdList = append(gIdList, group.Id)
	}

	pList, err := this.getPermissionListWithGroupIdList(tx, roleId, gIdList, status, "")
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

func (this *manager) getPermissionList(groupIdList []int64, status int, keyword string) (result []*Permission, err error) {
	var tx = dbs.MustTx(this.db)
	if result, err = this.getPermissionListWithGroupIdList(tx, 0, groupIdList, status, keyword); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *manager) getPermissionListWithGroupIdList(tx *dbs.Tx, roleId int64, groupIdList []int64, status int, keyword string) (result []*Permission, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("p.id", "p.group_id", "p.name", "p.identifier", "p.status", "p.created_on")
	sb.From(this.permissionTable, "AS p")
	if roleId > 0 {
		sb.Selects("IF(rp.role_id IS NULL, false , true) AS grant_to_role")
		sb.LeftJoin(this.rolePermissionTable, "AS rp ON rp.permission_id = p.id AND rp.role_id = ?", roleId)
	}

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
	sb.OrderBy("p.id")

	if err = tx.ExecSelectBuilder(sb, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *manager) getPermissionWithIdList(idList []int64) (result []*Permission, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("p.id", "p.group_id", "p.name", "p.identifier", "p.status", "p.created_on")
	sb.From(this.permissionTable, "AS p")
	if len(idList) > 0 {
		sb.Where(dbs.IN("r.id", idList))
	}
	sb.Limit(uint64(len(idList)))

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
