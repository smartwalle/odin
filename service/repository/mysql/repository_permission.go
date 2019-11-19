package mysql

import (
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/odin"
	"time"
)

func (this *odinRepository) GetPermissionTree(ctx, roleId int64, status odin.Status, name string) (result []*odin.Group, err error) {
	var tx = dbs.MustTx(this.db)

	if result, err = this.getGroupList(tx, ctx, odin.GroupTypeOfPermission, status, name); err != nil {
		return nil, err
	}

	var gMap = make(map[int64]*odin.Group)
	var gIdList []int64
	for _, group := range result {
		gMap[group.Id] = group
		gIdList = append(gIdList, group.Id)
	}

	pList, err := this.getPermissionListWithGroupIdList(tx, ctx, roleId, gIdList, status, "")
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

func (this *odinRepository) GetPermissionList(ctx int64, groupIdList []int64, status odin.Status, keyword string) (result []*odin.Permission, err error) {
	var tx = dbs.MustTx(this.db)
	if result, err = this.getPermissionListWithGroupIdList(tx, ctx, 0, groupIdList, status, keyword); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *odinRepository) getPermissionListWithGroupIdList(tx dbs.TX, ctx, roleId int64, groupIdList []int64, status odin.Status, keyword string) (result []*odin.Permission, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("p.id", "p.ctx", "p.group_id", "p.name", "p.identifier", "p.status", "p.created_on")
	sb.From(this.tblPermission, "AS p")
	if roleId > 0 {
		sb.Selects("IF(rp.role_id IS NULL, false , true) AS granted")
		sb.LeftJoin(this.tblRolePermission, "AS rp ON rp.permission_id = p.id AND rp.role_id = ?", roleId)
	}

	sb.Where("(p.ctx = ? OR p.ctx = ?)", 0, ctx)

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
	sb.OrderBy("p.ctx", "p.id")
	if err = sb.Scan(tx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *odinRepository) GetPermissionWithIdList(ctx int64, idList []int64) (result []*odin.Permission, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("p.id", "p.ctx", "p.group_id", "p.name", "p.identifier", "p.status", "p.created_on")
	sb.From(this.tblPermission, "AS p")
	sb.Where("(p.ctx = ? OR p.ctx = ?)", 0, ctx)
	if len(idList) > 0 {
		sb.Where(dbs.IN("p.id", idList))
	}
	sb.Limit(int64(len(idList)))

	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *odinRepository) GetPermissionWithId(ctx, id int64) (result *odin.Permission, err error) {
	var tx = dbs.MustTx(this.db)
	if result, err = this.getPermission(tx, ctx, id, "", ""); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, err
}

func (this *odinRepository) GetPermissionWithName(ctx int64, name string) (result *odin.Permission, err error) {
	var tx = dbs.MustTx(this.db)
	if result, err = this.getPermission(tx, ctx, 0, name, ""); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, err
}

func (this *odinRepository) GetPermissionWithIdentifier(ctx int64, identifier string) (result *odin.Permission, err error) {
	var tx = dbs.MustTx(this.db)
	if result, err = this.getPermission(tx, ctx, 0, "", identifier); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, err
}

func (this *odinRepository) AddPermission(ctx int64, groupId int64, name, identifier string, status odin.Status) (result *odin.Permission, err error) {
	var tx = dbs.MustTx(this.db)
	var newPermissionId int64 = 0
	if newPermissionId, err = this.insertPermission(tx, ctx, groupId, status, name, identifier); err != nil {
		return nil, err
	}
	if result, err = this.getPermission(tx, ctx, newPermissionId, "", ""); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, err
}

func (this *odinRepository) insertPermission(tx dbs.TX, ctx, groupId int64, status odin.Status, name, identifier string) (id int64, err error) {
	var ib = dbs.NewInsertBuilder()
	ib.Table(this.tblPermission)
	ib.Columns("ctx", "group_id", "status", "name", "identifier", "created_on", "updated_on")
	ib.Values(ctx, groupId, status, name, identifier, time.Now(), time.Now())
	if result, err := ib.Exec(tx); err != nil {
		return 0, err
	} else {
		id, _ = result.LastInsertId()
	}
	return id, err
}

func (this *odinRepository) UpdatePermission(ctx, id, groupId int64, name, identifier string, status odin.Status) (err error) {
	var ub = dbs.NewUpdateBuilder()
	ub.Table(this.tblPermission)
	ub.SET("group_id", groupId)
	ub.SET("name", name)
	ub.SET("identifier", identifier)
	ub.SET("status", status)
	ub.SET("updated_on", time.Now())
	ub.Where("id = ?", id)
	ub.Where("ctx = ?", ctx)
	ub.Limit(1)
	_, err = ub.Exec(this.db)
	return err
}

func (this *odinRepository) UpdatePermissionStatus(ctx, id int64, status odin.Status) (err error) {
	var ub = dbs.NewUpdateBuilder()
	ub.Table(this.tblPermission)
	ub.SET("status", status)
	ub.SET("updated_on", time.Now())
	ub.Where("id = ?", id)
	ub.Where("ctx = ?", ctx)
	ub.Limit(1)
	_, err = ub.Exec(this.db)
	return err
}

func (this *odinRepository) getPermission(tx dbs.TX, ctx, id int64, name, identifier string) (result *odin.Permission, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("p.id", "p.ctx", "p.group_id", "p.name", "p.identifier", "p.status", "p.created_on", "p.updated_on")
	sb.From(this.tblPermission, "AS p")
	if id > 0 {
		sb.Where("p.id = ?", id)
	}
	sb.Where("(p.ctx = ? OR p.ctx = ?)", 0, ctx)
	if name != "" {
		sb.Where("p.name = ?", name)
	}
	if identifier != "" {
		sb.Where("p.identifier = ?", identifier)
	}
	sb.Limit(1)
	if err = sb.Scan(tx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *odinRepository) GetPermissionListWithRoleId(ctx, roleId int64) (result []*odin.Permission, err error) {
	var tx = dbs.MustTx(this.db)
	if result, err = this.getPermissionListWithRole(tx, ctx, roleId); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, err
}

func (this *odinRepository) getPermissionListWithRole(tx dbs.TX, ctx, roleId int64) (result []*odin.Permission, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("p.id", "p.ctx", "p.group_id", "p.name", "p.identifier", "p.status", "p.created_on", "p.updated_on")
	sb.Selects("IF(rp.role_id IS NULL, false, true) AS granted")
	sb.From(this.tblPermission, "AS p")
	sb.LeftJoin(this.tblRolePermission, "AS rp ON rp.permission_id = p.id")
	sb.Where("(p.ctx = ? OR p.ctx = ?)", 0, ctx)
	sb.Where("rp.role_id = ?", roleId)
	if err = sb.Scan(tx, &result); err != nil {
		return nil, err
	}
	return result, err
}

func (this *odinRepository) GetGrantedPermissionList(ctx int64, target string) (result []*odin.Permission, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("p.id", "p.ctx", "p.group_id", "p.name", "p.identifier", "p.status", "p.created_on", "p.updated_on")
	sb.Selects("IF(rp.role_id IS NULL, false, true) AS granted")
	sb.From(this.tblPermission, "AS p")
	sb.LeftJoin(this.tblRolePermission, "AS rp ON rp.permission_id = p.id")
	sb.LeftJoin(this.tblGrant, "AS rg ON rg.role_id = rp.role_id")
	sb.Where("rg.target = ? AND p.status = ?", target, odin.StatusOfEnable)
	sb.Where("(p.ctx = ? OR p.ctx = ?)", 0, ctx)
	sb.GroupBy("p.id")
	if err := sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, err
}

func (this *odinRepository) GrantPermission(ctx, roleId int64, permissionIdList []int64) (err error) {
	if len(permissionIdList) > 0 {
		var now = time.Now()
		var ib = dbs.NewInsertBuilder()
		ib.Table(this.tblRolePermission)
		ib.Options("IGNORE")
		ib.Columns("ctx", "role_id", "permission_id", "created_on")
		for _, pId := range permissionIdList {
			ib.Values(ctx, roleId, pId, now)
		}
		if _, err = ib.Exec(this.db); err != nil {
			return err
		}
	}
	return nil
}

func (this *odinRepository) RevokePermission(ctx, roleId int64, permissionIdList []int64) (err error) {
	var rb = dbs.NewDeleteBuilder()
	rb.Table(this.tblRolePermission)
	rb.Where("role_id = ?", roleId)
	//rb.Where("(ctx = ? OR ctx = ?)", 0, ctx)
	rb.Where("ctx = ?", ctx)
	rb.Where(dbs.IN("permission_id", permissionIdList))
	if _, err = rb.Exec(this.db); err != nil {
		return err
	}
	return nil
}

func (this *odinRepository) ReGrantPermission(ctx, roleId int64, permissionIdList []int64) (err error) {
	var tx = dbs.MustTx(this.db)
	var now = time.Now()

	var rb = dbs.NewDeleteBuilder()
	rb.Table(this.tblRolePermission)
	rb.Where("role_id = ?", roleId)
	rb.Where("ctx = ?", ctx)
	if _, err = rb.Exec(tx); err != nil {
		return err
	}

	if len(permissionIdList) > 0 {
		var ib = dbs.NewInsertBuilder()
		ib.Table(this.tblRolePermission)
		ib.Options("IGNORE")
		ib.Columns("ctx", "role_id", "permission_id", "created_on")
		for _, pId := range permissionIdList {
			ib.Values(ctx, roleId, pId, now)
		}
		if _, err = ib.Exec(tx); err != nil {
			return err
		}
	}
	tx.Commit()
	return nil
}
