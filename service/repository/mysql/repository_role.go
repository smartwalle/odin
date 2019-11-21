package mysql

import (
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/odin"
	"time"
)

func (this *odinRepository) GetRoleList(ctx int64, target string, groupIdList []int64, status odin.Status, keyword string) (result []*odin.Role, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("r.id", "r.ctx", "r.group_id", "r.name", "r.status", "r.created_on", "r.updated_on")
	sb.From(this.tblRole, "AS r")
	if target != "" {
		sb.Selects("IF(rg.target IS NULL, false, true) AS granted")
		sb.LeftJoin(this.tblGrant, "AS rg ON rg.role_id = r.id AND rg.target = ?", target)
	}
	sb.Where("r.ctx = ?", ctx)
	if len(groupIdList) > 0 {
		sb.Where(dbs.IN("r.group_id", groupIdList))
	}
	if status > 0 {
		sb.Where("r.status = ?", status)
	}
	if keyword != "" {
		sb.Where(dbs.Like("r.name", "%", keyword, "%"))
	}
	sb.OrderBy("r.ctx", "r.id")
	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *odinRepository) GetRoleWithId(ctx, id int64, withPermissionList bool) (result *odin.Role, err error) {
	if result, err = this.getRole(ctx, id, ""); err != nil {
		return nil, err
	}
	if withPermissionList {
		if result.PermissionList, err = this.GetRolePermissionList(ctx, result.Id); err != nil {
			return nil, err
		}
	}
	return result, err
}

func (this *odinRepository) GetRoleWithName(ctx int64, name string, withPermissionList bool) (result *odin.Role, err error) {
	if result, err = this.getRole(ctx, 0, name); err != nil {
		return nil, err
	}
	if withPermissionList {
		if result.PermissionList, err = this.GetRolePermissionList(ctx, result.Id); err != nil {
			return nil, err
		}
	}
	return result, err
}

func (this *odinRepository) AddRole(ctx, groupId int64, name string, status odin.Status) (result int64, err error) {
	var ib = dbs.NewInsertBuilder()
	ib.Table(this.tblRole)
	ib.Columns("ctx", "group_id", "status", "name", "created_on", "updated_on")
	ib.Values(ctx, groupId, status, name, time.Now(), time.Now())
	rResult, err := ib.Exec(this.db)
	if err != nil {
		return 0, err
	}
	result, _ = rResult.LastInsertId()
	return result, nil
}

func (this *odinRepository) UpdateRole(ctx, id, groupId int64, name string, status odin.Status) (err error) {
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

func (this *odinRepository) UpdateRoleStatus(ctx, id int64, status odin.Status) (err error) {
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

func (this *odinRepository) getRole(ctx, id int64, name string) (result *odin.Role, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("r.id", "r.ctx", "r.group_id", "r.name", "r.status", "r.created_on", "r.updated_on")
	sb.From(this.tblRole, "AS r")
	if id > 0 {
		sb.Where("r.id = ?", id)
	}
	sb.Where("r.ctx = ?", ctx)
	if name != "" {
		sb.Where("r.name = ?", name)
	}
	sb.Limit(1)
	if err = sb.Scan(this.db, &result); err != nil {
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
	sb.Where("r.ctx = ?", ctx)
	sb.Limit(int64(len(idList)))

	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *odinRepository) GrantRole(ctx int64, target string, roleIdList []int64) (err error) {
	var now = time.Now()
	var ib = dbs.NewInsertBuilder()
	ib.Table(this.tblGrant)
	ib.Options("IGNORE")
	ib.Columns("ctx", "target", "role_id", "created_on")
	for _, rId := range roleIdList {
		ib.Values(ctx, target, rId, now)
	}
	if _, err = ib.Exec(this.db); err != nil {
		return err
	}
	return nil
}

func (this *odinRepository) RevokeRole(ctx int64, target string, roleIdList []int64) (err error) {
	var rb = dbs.NewDeleteBuilder()
	rb.Table(this.tblGrant)
	rb.Where("target = ?", target)
	rb.Where("ctx = ?", ctx)
	rb.Where(dbs.IN("role_id", roleIdList))
	if _, err = rb.Exec(this.db); err != nil {
		return err
	}
	return nil
}

func (this *odinRepository) ReGrantRole(ctx int64, target string, roleIdList []int64) (err error) {
	var now = time.Now()

	var rb = dbs.NewDeleteBuilder()
	rb.Table(this.tblGrant)
	rb.Where("target = ?", target)
	rb.Where("ctx = ?", ctx)
	if _, err = rb.Exec(this.db); err != nil {
		return err
	}

	if len(roleIdList) > 0 {
		var ib = dbs.NewInsertBuilder()
		ib.Table(this.tblGrant)
		ib.Options("IGNORE")
		ib.Columns("ctx", "target", "role_id", "created_on")
		for _, rId := range roleIdList {
			ib.Values(ctx, target, rId, now)
		}
		if _, err = ib.Exec(this.db); err != nil {
			return err
		}
	}
	return nil
}

func (this *odinRepository) GetGrantedRoleList(ctx int64, target string) (result []*odin.Role, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("r.id", "r.group_id", "r.name", "r.status", "r.created_on", "r.updated_on")
	sb.From(this.tblRole, "AS r")
	sb.Selects("IF(rg.target IS NULL, false, true) AS granted")
	sb.LeftJoin(this.tblGrant, "AS rg ON rg.role_id = r.id")
	sb.Where("rg.target = ? AND r.status = ?", target, odin.Enable)
	sb.Where("rg.ctx = ?", ctx)
	if err := sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, err
}
