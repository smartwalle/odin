package odin

import (
	"github.com/smartwalle/dbs"
	"time"
)

func (this *manager) getRoleTree(status int, name string) (result []*Group, err error) {
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

	rList, err := this.getRoleListWithGroupIdList(tx, gIdList, status, "")
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
	if result, err = this.getRoleListWithGroupIdList(tx, groupIdList, status, keyword); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *manager) getRoleListWithGroupIdList(tx *dbs.Tx, groupIdList []int64, status int, keyword string) (result []*Role, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("r.id", "r.group_id", "r.name", "r.status", "r.created_on")
	sb.From(this.roleTable, "AS r")
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

	if err = tx.ExecSelectBuilder(sb, &result); err != nil {
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

func (this *manager) insertRole(tx *dbs.Tx, groupId int64, status int, name string) (id int64, err error) {
	var ib = dbs.NewInsertBuilder()
	ib.Table(this.roleTable)
	ib.Columns("group_id", "status", "name", "created_on")
	ib.Values(groupId, status, name, time.Now())
	if result, err := tx.ExecInsertBuilder(ib); err != nil {
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
	ub.Where("id = ?", id)
	ub.Limit(1)
	_, err = ub.Exec(this.db)
	return err
}

func (this *manager) updateRoleStatus(id int64, status int) (err error) {
	var ub = dbs.NewUpdateBuilder()
	ub.Table(this.roleTable)
	ub.SET("status", status)
	ub.Where("id = ?", id)
	ub.Limit(1)
	_, err = ub.Exec(this.db)
	return err
}

func (this *manager) getRole(tx *dbs.Tx, id int64, name string) (result *Role, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("r.id", "r.group_id", "r.name", "r.status", "r.created_on")
	sb.From(this.roleTable, "AS r")
	if id > 0 {
		sb.Where("r.id = ?", id)
	}
	if name != "" {
		sb.Where("r.name = ?", name)
	}
	sb.Limit(1)
	if err = tx.ExecSelectBuilder(sb, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *manager) getRoleWithIdList(idList []int64) (result []*Role, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("r.id", "r.group_id", "r.name", "r.status", "r.created_on")
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

func (this *manager) getPermissionListWithRoleId(roleId int64) (result []*Permission, err error) {
	var tx = dbs.MustTx(this.db)
	if result, err = this.getPermissionListWithRole(tx, roleId); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, err
}

func (this *manager) getPermissionListWithRole(tx *dbs.Tx, roleId int64) (result []*Permission, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("p.id", "p.group_id", "p.name", "p.identifier", "p.status", "p.created_on")
	sb.From(this.permissionTable, "AS p")
	sb.LeftJoin(this.rolePermissionTable, "AS rp ON rp.permission_id = p.id")
	sb.Where("rp.role_id = ?", roleId)
	if err = tx.ExecSelectBuilder(sb, &result); err != nil {
		return nil, err
	}
	return result, err
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

func (this *manager) Check(objectId, identifier string) (result bool) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("rg.object_id", "rg.role_id", "rp.permission_id", "p.identifier")
	sb.From(this.roleGrantTable, "AS rg")
	sb.LeftJoin(this.rolePermissionTable, "AS rp ON rp.role_id = rg.role_id")
	sb.LeftJoin(this.permissionTable, "AS p ON p.id = rp.permission_id")
	sb.Where("rg.object_id = ? AND p.identifier = ?", objectId, identifier)
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
