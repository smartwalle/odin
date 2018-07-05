package odin

import (
	"github.com/smartwalle/dbs"
	"time"
)

func (this *manager) getRoleTree(objectId string, status int, name string) (result []*Group, err error) {
	var tx = dbs.MustTx(this.db)

	if result, err = this.getGroupList(tx, K_GROUP_TYPE_ROLE, status, name); err != nil {
		return nil, err
	}

	var gMap = make(map[int64]*Group)
	var gIdList []int64
	for _, group := range result {
		gMap[group.Id] = group
		gIdList = append(gIdList, group.Id)
	}

	rList, err := this.getRoles(tx, objectId, gIdList, status, "")
	if err != nil {
		return nil, err
	}

	for _, r := range rList {
		var group = gMap[r.GroupId]
		if group != nil {
			group.RoleList = append(group.RoleList, r)
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *manager) getRoleList(groupId int64, status int, keyword string) (result []*Role, err error) {
	var tx = dbs.MustTx(this.db)
	var groupIdList []int64
	if groupId > 0 {
		groupIdList = append(groupIdList, groupId)
	}
	if result, err = this.getRoles(tx, "", groupIdList, status, keyword); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *manager) getRoles(tx dbs.TX, objectId string, groupIdList []int64, status int, keyword string) (result []*Role, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("r.id", "r.group_id", "r.name", "r.status", "r.created_on", "r.updated_on")
	sb.From(this.roleTable, "AS r")
	if objectId != "" {
		sb.Selects("IF(rg.object_id IS NULL, false, true) AS granted")
		sb.LeftJoin(this.roleGrantTable, "AS rg ON rg.role_id = r.id AND rg.object_id = ?", objectId)
	}
	if len(groupIdList) > 0 {
		sb.Where(dbs.IN("r.group_id", groupIdList))
	}
	if status > 0 {
		sb.Where("r.status = ?", status)
	}
	if keyword != "" {
		var k = "%" + keyword + "%"
		sb.Where("r.name LIKE ?", k, k)
	}
	sb.OrderBy("r.id")

	//if err = tx.ExecSelectBuilder(sb, &result); err != nil {
	//	return nil, err
	//}
	if err = sb.ScanTx(tx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *manager) getRoleWithId(id int64, withPermissionList bool) (result *Role, err error) {
	var tx = dbs.MustTx(this.db)
	if result, err = this.getRole(tx, id, ""); err != nil {
		return nil, err
	}
	if withPermissionList {
		if result.PermissionList, err = this.getPermissionListWithRole(tx, result.Id); err != nil {
			return nil, err
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, err
}

func (this *manager) getRoleWithName(name string, withPermissionList bool) (result *Role, err error) {
	var tx = dbs.MustTx(this.db)
	if result, err = this.getRole(tx, 0, name); err != nil {
		return nil, err
	}
	if withPermissionList {
		if result.PermissionList, err = this.getPermissionListWithRole(tx, result.Id); err != nil {
			return nil, err
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, err
}

func (this *manager) addRole(groupId int64, name string, status int) (result *Role, err error) {
	var tx = dbs.MustTx(this.db)
	var newRoleId int64 = 0
	if newRoleId, err = this.insertRole(tx, groupId, status, name); err != nil {
		return nil, err
	}
	if result, err = this.getRole(tx, newRoleId, ""); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, err
}

func (this *manager) insertRole(tx dbs.TX, groupId int64, status int, name string) (id int64, err error) {
	var ib = dbs.NewInsertBuilder()
	ib.Table(this.roleTable)
	ib.Columns("group_id", "status", "name", "created_on", "updated_on")
	ib.Values(groupId, status, name, time.Now(), time.Now())
	//if result, err := tx.ExecInsertBuilder(ib); err != nil {
	//	return 0, err
	//} else {
	//	id, _ = result.LastInsertId()
	//}
	if result, err := ib.ExecTx(tx); err != nil {
		return 0, err
	} else {
		id, _ = result.LastInsertId()
	}
	return id, err
}

func (this *manager) updateRole(id, groupId int64, name string, status int) (err error) {
	var ub = dbs.NewUpdateBuilder()
	ub.Table(this.roleTable)
	ub.SET("group_id", groupId)
	ub.SET("name", name)
	ub.SET("status", status)
	ub.SET("updated_on", time.Now())
	ub.Where("id = ?", id)
	ub.Limit(1)
	_, err = ub.Exec(this.db)
	return err
}

func (this *manager) updateRoleStatus(id int64, status int) (err error) {
	var ub = dbs.NewUpdateBuilder()
	ub.Table(this.roleTable)
	ub.SET("status", status)
	ub.SET("updated_on", time.Now())
	ub.Where("id = ?", id)
	ub.Limit(1)
	_, err = ub.Exec(this.db)
	return err
}

func (this *manager) getRole(tx dbs.TX, id int64, name string) (result *Role, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("r.id", "r.group_id", "r.name", "r.status", "r.created_on", "r.updated_on")
	sb.From(this.roleTable, "AS r")
	if id > 0 {
		sb.Where("r.id = ?", id)
	}
	if name != "" {
		sb.Where("r.name = ?", name)
	}
	sb.Limit(1)
	//if err = tx.ExecSelectBuilder(sb, &result); err != nil {
	//	return nil, err
	//}
	if err = sb.ScanTx(tx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *manager) getRoleWithIdList(idList []int64) (result []*Role, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("r.id", "r.group_id", "r.name", "r.status", "r.created_on", "r.updated_on")
	sb.From(this.roleTable, "AS r")
	if len(idList) > 0 {
		sb.Where(dbs.IN("r.id", idList))
	}
	sb.Limit(uint64(len(idList)))

	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *manager) grantPermission(roleId int64, permissionIdList []int64) (err error) {
	var now = time.Now()
	var ib = dbs.NewInsertBuilder()
	ib.Table(this.rolePermissionTable)
	ib.Options("IGNORE")
	ib.Columns("role_id", "permission_id", "created_on")
	for _, pId := range permissionIdList {
		ib.Values(roleId, pId, now)
	}
	if _, err = ib.Exec(this.db); err != nil {
		return err
	}
	return nil
}

func (this *manager) revokePermission(roleId int64, permissionIdList []int64) (err error) {
	var rb = dbs.NewDeleteBuilder()
	rb.Table(this.rolePermissionTable)
	rb.Where("role_id = ?", roleId)
	rb.Where(dbs.IN("permission_id", permissionIdList))
	if _, err = rb.Exec(this.db); err != nil {
		return err
	}
	return nil
}

func (this *manager) grantRole(objectId string, roleIdList []int64) (err error) {
	var now = time.Now()
	var ib = dbs.NewInsertBuilder()
	ib.Table(this.roleGrantTable)
	ib.Options("IGNORE")
	ib.Columns("object_id", "role_id", "created_on")
	for _, rId := range roleIdList {
		ib.Values(objectId, rId, now)
	}
	if _, err = ib.Exec(this.db); err != nil {
		return err
	}
	return nil
}

func (this *manager) revokeRole(objectId string, roleIdList []int64) (err error) {
	var rb = dbs.NewDeleteBuilder()
	rb.Table(this.roleGrantTable)
	rb.Where("object_id = ?", objectId)
	rb.Where(dbs.IN("role_id", roleIdList))
	if _, err = rb.Exec(this.db); err != nil {
		return err
	}
	return nil
}

func (this *manager) check(objectId, identifier string) (result bool) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("rg.object_id", "rg.role_id", "rp.permission_id", "p.identifier")
	sb.From(this.roleGrantTable, "AS rg")
	sb.LeftJoin(this.rolePermissionTable, "AS rp ON rp.role_id = rg.role_id")
	sb.LeftJoin(this.permissionTable, "AS p ON p.id = rp.permission_id")
	sb.LeftJoin(this.roleTable, "AS r ON r.id = rg.role_id")
	sb.Where("rg.object_id = ? AND p.identifier = ? AND p.status = ? AND r.status = ?", objectId, identifier, K_STATUS_ENABLE, K_STATUS_ENABLE)
	sb.OrderBy("r.status", "p.status")
	sb.Limit(1)

	var grant *Grant
	if err := sb.Scan(this.db, &grant); err != nil || grant == nil {
		return false
	}
	if grant.Identifier == identifier && grant.ObjectId == objectId {
		return true
	}
	return false
}

func (this *manager) checkList(objectId string, identifiers ...string) (result map[string]bool) {
	result = make(map[string]bool)
	var sb = dbs.NewSelectBuilder()
	sb.Selects("rg.object_id", "rg.role_id", "rp.permission_id", "p.identifier")
	sb.From(this.roleGrantTable, "AS rg")
	sb.LeftJoin(this.rolePermissionTable, "AS rp ON rp.role_id = rg.role_id")
	sb.LeftJoin(this.permissionTable, "AS p ON p.id = rp.permission_id")
	sb.LeftJoin(this.roleTable, "AS r ON r.id = rg.role_id")


	sb.Where("rg.object_id = ?", objectId)
	if len(identifiers) > 0 {
		var or = dbs.OR()
		for _, identifier := range identifiers {
			or.Append("p.identifier = ?", identifier)
			result[identifier] = false
		}
		sb.Where(or)
	}
	sb.Where("p.status = ?", K_STATUS_ENABLE)
	sb.Where("r.status = ?", K_STATUS_ENABLE)

	sb.GroupBy("p.id")
	sb.OrderBy("r.status", "p.status")
	sb.Limit(uint64(len(identifiers)))

	var grantList []*Grant
	if err := sb.Scan(this.db, &grantList); err != nil || grantList == nil {
		return result
	}
	for _, grant := range grantList {
		result[grant.Identifier] = true
	}
	return result
}

func (this *manager) getGrantedRoleList(objectId string) (result []*Role, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("r.id", "r.group_id", "r.name", "r.status", "r.created_on", "r.updated_on")
	sb.From(this.roleTable, "AS r")
	sb.Selects("IF(rg.object_id IS NULL, false, true) AS granted")
	sb.LeftJoin(this.roleGrantTable, "AS rg ON rg.role_id = r.id")
	sb.Where("rg.object_id = ? AND r.status = ?", objectId, K_STATUS_ENABLE)

	if err := sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, err
}
