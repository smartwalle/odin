package mysql

import (
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/odin"
	"time"
)

func (this *odinRepository) GetPermissionList(ctx int64, groupIdList []int64, status odin.Status, keyword string, roleId int64) (result []*odin.Permission, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("p.id", "p.ctx", "p.group_id", "p.name", "p.identifier", "p.status", "p.created_on", "p.updated_on")
	sb.From(this.tblPermission, "AS p")
	if roleId > 0 {
		sb.Selects("IF(rp.role_id IS NULL, false , true) AS granted")
		sb.LeftJoin(this.tblRolePermission, "AS rp ON rp.permission_id = p.id AND rp.role_id = ?", roleId)
	}

	sb.Where("p.ctx = ?", ctx)

	if len(groupIdList) > 0 {
		sb.Where(dbs.IN("p.group_id", groupIdList))
	}
	if status > 0 {
		sb.Where("p.status = ?", status)
	}
	if keyword != "" {
		var or = dbs.OR()
		or.Append("p.name", "%", keyword, "%")
		or.Append("p.identifier", "%", keyword, "%")
		sb.Where(or)
	}
	sb.OrderBy("p.ctx", "p.id")
	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *odinRepository) GetPermissionListWithIds(ctx int64, idList []int64) (result []*odin.Permission, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("p.id", "p.ctx", "p.group_id", "p.name", "p.identifier", "p.status", "p.created_on", "p.updated_on")
	sb.From(this.tblPermission, "AS p")
	sb.Where("p.ctx = ?", ctx)
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
	return this.getPermission(ctx, id, "", "")
}

func (this *odinRepository) GetPermissionWithName(ctx int64, name string) (result *odin.Permission, err error) {
	return this.getPermission(ctx, 0, name, "")
}

func (this *odinRepository) GetPermissionWithIdentifier(ctx int64, identifier string) (result *odin.Permission, err error) {
	return this.getPermission(ctx, 0, "", identifier)
}

func (this *odinRepository) getPermission(ctx, id int64, name, identifier string) (result *odin.Permission, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("p.id", "p.ctx", "p.group_id", "p.name", "p.identifier", "p.status", "p.created_on", "p.updated_on")
	sb.From(this.tblPermission, "AS p")
	if id > 0 {
		sb.Where("p.id = ?", id)
	}
	sb.Where("p.ctx = ?", ctx)
	if name != "" {
		sb.Where("p.name = ?", name)
	}
	if identifier != "" {
		sb.Where("p.identifier = ?", identifier)
	}
	sb.Limit(1)
	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *odinRepository) AddPermission(ctx int64, groupId int64, name, identifier string, status odin.Status) (result int64, err error) {
	var ib = dbs.NewInsertBuilder()
	ib.Table(this.tblPermission)
	ib.Columns("ctx", "group_id", "status", "name", "identifier", "created_on", "updated_on")
	ib.Values(ctx, groupId, status, name, identifier, time.Now(), time.Now())
	rResult, err := ib.Exec(this.db)
	if err != nil {
		return 0, err
	}
	result, _ = rResult.LastInsertId()
	return result, err
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

func (this *odinRepository) GetRolePermissionList(ctx, roleId int64) (result []*odin.Permission, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("p.id", "p.ctx", "p.group_id", "p.name", "p.identifier", "p.status", "p.created_on", "p.updated_on")
	sb.Selects("IF(rp.role_id IS NULL, false, true) AS granted")
	sb.From(this.tblPermission, "AS p")
	sb.LeftJoin(this.tblRolePermission, "AS rp ON rp.permission_id = p.id")
	sb.Where("p.ctx = ?", ctx)
	sb.Where("rp.role_id = ?", roleId)
	if err = sb.Scan(this.db, &result); err != nil {
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
	sb.Where("p.ctx = ?", ctx)
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
	rb.Where("ctx = ?", ctx)
	rb.Where("ctx = ?", ctx)
	rb.Where(dbs.IN("permission_id", permissionIdList))
	if _, err = rb.Exec(this.db); err != nil {
		return err
	}
	return nil
}

func (this *odinRepository) ReGrantPermission(ctx, roleId int64, permissionIdList []int64) (err error) {
	var now = time.Now()
	var rb = dbs.NewDeleteBuilder()
	rb.Table(this.tblRolePermission)
	rb.Where("role_id = ?", roleId)
	rb.Where("ctx = ?", ctx)
	if _, err = rb.Exec(this.db); err != nil {
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
		if _, err = ib.Exec(this.db); err != nil {
			return err
		}
	}
	return nil
}
