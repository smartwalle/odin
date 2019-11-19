package mysql

import (
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/odin"
	"time"
)

func (this *odinRepository) GetRoleTree(ctx int64, targetId string, status int, name string) (result []*odin.Group, err error) {
	var tx = dbs.MustTx(this.db)

	if result, err = this.getGroupList(tx, ctx, odin.K_GROUP_TYPE_ROLE, status, name); err != nil {
		return nil, err
	}

	var gMap = make(map[int64]*odin.Group)
	var gIdList []int64
	for _, group := range result {
		gMap[group.Id] = group
		gIdList = append(gIdList, group.Id)
	}

	rList, err := this.getRoles(tx, ctx, targetId, gIdList, status, "")
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

func (this *odinRepository) GetRoleList(ctx, groupId int64, status int, keyword string) (result []*odin.Role, err error) {
	var tx = dbs.MustTx(this.db)
	var groupIdList []int64
	if groupId > 0 {
		groupIdList = append(groupIdList, groupId)
	}
	if result, err = this.getRoles(tx, ctx, "", groupIdList, status, keyword); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *odinRepository) getRoles(tx dbs.TX, ctx int64, targetId string, groupIdList []int64, status int, keyword string) (result []*odin.Role, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("r.id", "r.ctx", "r.group_id", "r.name", "r.status", "r.created_on", "r.updated_on")
	sb.From(this.tblRole, "AS r")
	if targetId != "" {
		sb.Selects("IF(rg.target_id IS NULL, false, true) AS granted")
		sb.LeftJoin(this.tblGrant, "AS rg ON rg.role_id = r.id AND rg.target_id = ?", targetId)
	}
	//sb.Where("(r.ctx = ? OR r.ctx = ?)", 0, ctx)
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
	sb.OrderBy("r.ctx", "r.id")
	if err = sb.Scan(tx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *odinRepository) GetRoleWithId(ctx, id int64, withPermissionList bool) (result *odin.Role, err error) {
	var tx = dbs.MustTx(this.db)
	if result, err = this.getRole(tx, ctx, id, ""); err != nil {
		return nil, err
	}
	if withPermissionList {
		if result.PermissionList, err = this.getPermissionListWithRole(tx, ctx, result.Id); err != nil {
			return nil, err
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, err
}

func (this *odinRepository) GetRoleWithName(ctx int64, name string, withPermissionList bool) (result *odin.Role, err error) {
	var tx = dbs.MustTx(this.db)
	if result, err = this.getRole(tx, ctx, 0, name); err != nil {
		return nil, err
	}
	if withPermissionList {
		if result.PermissionList, err = this.getPermissionListWithRole(tx, ctx, result.Id); err != nil {
			return nil, err
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, err
}

func (this *odinRepository) AddRole(ctx, groupId int64, name string, status int) (result *odin.Role, err error) {
	var tx = dbs.MustTx(this.db)
	var newRoleId int64 = 0
	if newRoleId, err = this.insertRole(tx, ctx, groupId, status, name); err != nil {
		return nil, err
	}
	if result, err = this.getRole(tx, ctx, newRoleId, ""); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, err
}

func (this *odinRepository) insertRole(tx dbs.TX, ctx, groupId int64, status int, name string) (id int64, err error) {
	var ib = dbs.NewInsertBuilder()
	ib.Table(this.tblRole)
	ib.Columns("ctx", "group_id", "status", "name", "created_on", "updated_on")
	ib.Values(ctx, groupId, status, name, time.Now(), time.Now())
	if result, err := ib.Exec(tx); err != nil {
		return 0, err
	} else {
		id, _ = result.LastInsertId()
	}
	return id, err
}

func (this *odinRepository) UpdateRole(ctx, id, groupId int64, name string, status int) (err error) {
	var ub = dbs.NewUpdateBuilder()
	ub.Table(this.tblRole)
	ub.SET("group_id", groupId)
	ub.SET("name", name)
	ub.SET("status", status)
	ub.SET("updated_on", time.Now())
	ub.Where("id = ?", id)
	ub.Where("ctx = ?", ctx)
	ub.Limit(1)
	_, err = ub.Exec(this.db)
	return err
}

func (this *odinRepository) UpdateRoleStatus(ctx, id int64, status int) (err error) {
	var ub = dbs.NewUpdateBuilder()
	ub.Table(this.tblRole)
	ub.SET("status", status)
	ub.SET("updated_on", time.Now())
	ub.Where("id = ?", id)
	ub.Where("ctx = ?", ctx)
	ub.Limit(1)
	_, err = ub.Exec(this.db)
	return err
}

func (this *odinRepository) getRole(tx dbs.TX, ctx, id int64, name string) (result *odin.Role, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("r.id", "r.ctx", "r.group_id", "r.name", "r.status", "r.created_on", "r.updated_on")
	sb.From(this.tblRole, "AS r")
	if id > 0 {
		sb.Where("r.id = ?", id)
	}
	sb.Where("(r.ctx = ? OR r.ctx = ?)", 0, ctx)
	if name != "" {
		sb.Where("r.name = ?", name)
	}
	sb.Limit(1)
	if err = sb.Scan(tx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *odinRepository) GetRoleWithIdList(ctx int64, idList []int64) (result []*odin.Role, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("r.id", "r.ctx", "r.group_id", "r.name", "r.status", "r.created_on", "r.updated_on")
	sb.From(this.tblRole, "AS r")
	if len(idList) > 0 {
		sb.Where(dbs.IN("r.id", idList))
	}
	sb.Where("(r.ctx = ? OR r.ctx = ?)", 0, ctx)
	sb.Limit(int64(len(idList)))

	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *odinRepository) GrantRole(ctx int64, targetId string, roleIdList []int64) (err error) {
	var now = time.Now()
	var ib = dbs.NewInsertBuilder()
	ib.Table(this.tblGrant)
	ib.Options("IGNORE")
	ib.Columns("ctx", "target_id", "role_id", "created_on")
	for _, rId := range roleIdList {
		ib.Values(ctx, targetId, rId, now)
	}
	if _, err = ib.Exec(this.db); err != nil {
		return err
	}
	return nil
}

func (this *odinRepository) RevokeRole(ctx int64, targetId string, roleIdList []int64) (err error) {
	var rb = dbs.NewDeleteBuilder()
	rb.Table(this.tblGrant)
	rb.Where("target_id = ?", targetId)
	//rb.Where("(ctx = ? OR ctx = ?)", 0, ctx)
	rb.Where("ctx = ?", ctx)
	rb.Where(dbs.IN("role_id", roleIdList))
	if _, err = rb.Exec(this.db); err != nil {
		return err
	}
	return nil
}

func (this *odinRepository) ReGrantRole(ctx int64, targetId string, roleIdList []int64) (err error) {
	var tx = dbs.MustTx(this.db)
	var now = time.Now()

	var rb = dbs.NewDeleteBuilder()
	rb.Table(this.tblGrant)
	rb.Where("target_id = ?", targetId)
	rb.Where("ctx = ?", ctx)
	if _, err = rb.Exec(tx); err != nil {
		return err
	}

	if len(roleIdList) > 0 {
		var ib = dbs.NewInsertBuilder()
		ib.Table(this.tblGrant)
		ib.Options("IGNORE")
		ib.Columns("ctx", "target_id", "role_id", "created_on")
		for _, rId := range roleIdList {
			ib.Values(ctx, targetId, rId, now)
		}
		if _, err = ib.Exec(tx); err != nil {
			return err
		}
	}
	tx.Commit()
	return nil
}

func (this *odinRepository) GetGrantedRoleList(ctx int64, targetId string) (result []*odin.Role, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("r.id", "r.group_id", "r.name", "r.status", "r.created_on", "r.updated_on")
	sb.From(this.tblRole, "AS r")
	sb.Selects("IF(rg.target_id IS NULL, false, true) AS granted")
	sb.LeftJoin(this.tblGrant, "AS rg ON rg.role_id = r.id")
	sb.Where("rg.target_id = ? AND r.status = ?", targetId, odin.K_STATUS_ENABLE)
	sb.Where("(rg.ctx = ? OR rg.ctx = ?)", 0, ctx)
	if err := sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, err
}
