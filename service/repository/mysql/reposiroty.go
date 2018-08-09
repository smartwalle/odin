package mysql

import (
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/odin"
	"github.com/smartwalle/odin/service"
	"strings"
	"time"
)

type odinRepository struct {
	db                  dbs.DB
	groupTable          string
	permissionTable     string
	roleTable           string
	rolePermissionTable string
	roleGrantTable      string
}

func NewOdinRepository(db dbs.DB, tablePrefix string) service.OdinRepository {
	var r = &odinRepository{}
	r.db = db

	tablePrefix = strings.TrimSpace(tablePrefix)
	if tablePrefix == "" {
		tablePrefix = "odin"
	}

	r.db = db
	r.groupTable = tablePrefix + "_group"
	r.permissionTable = tablePrefix + "_permission"
	r.roleTable = tablePrefix + "_role"
	r.rolePermissionTable = tablePrefix + "_role_permission"
	r.roleGrantTable = tablePrefix + "_grant"

	return r
}

func (this *odinRepository) GetGroupListWithType(ctx int64, gType, status int, name string) (result []*odin.Group, err error) {
	var tx = dbs.MustTx(this.db)
	if result, err = this.getGroupList(tx, ctx, gType, status, name); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, err
}

func (this *odinRepository) getGroupList(tx dbs.TX, ctx int64, gType, status int, name string) (result []*odin.Group, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("g.id", "g.ctx", "g.type", "g.name", "g.status", "g.created_on", "g.updated_on")
	sb.From(this.groupTable, "AS g")

	sb.Where("(g.ctx = ? OR g.ctx = ?)", 0, ctx)

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
	sb.OrderBy("g.ctx", "g.id")

	if err = sb.Scan(tx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *odinRepository) GetGroupWithId(ctx, id int64, gType int) (result *odin.Group, err error) {
	var tx = dbs.MustTx(this.db)
	if result, err = this.getGroup(tx, ctx, id, gType, ""); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, err
}

func (this *odinRepository) GetGroupWithName(ctx int64, name string, gType int) (result *odin.Group, err error) {
	var tx = dbs.MustTx(this.db)
	if result, err = this.getGroup(tx, ctx, 0, gType, name); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, err
}

func (this *odinRepository) AddGroup(ctx int64, gType int, name string, status int) (result *odin.Group, err error) {
	var tx = dbs.MustTx(this.db)
	var newGroupId int64 = 0
	if newGroupId, err = this.insertGroup(tx, ctx, gType, status, name); err != nil {
		return nil, err
	}
	if result, err = this.getGroup(tx, ctx, newGroupId, 0, ""); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, err
}

func (this *odinRepository) insertGroup(tx dbs.TX, ctx int64, gType, status int, name string) (id int64, err error) {
	var ib = dbs.NewInsertBuilder()
	ib.Table(this.groupTable)
	ib.Columns("ctx", "type", "status", "name", "created_on", "updated_on")
	ib.Values(ctx, gType, status, name, time.Now(), time.Now())
	if result, err := ib.Exec(tx); err != nil {
		return 0, err
	} else {
		id, _ = result.LastInsertId()
	}
	return id, err
}

func (this *odinRepository) UpdateGroup(ctx, id int64, name string, status int) (err error) {
	var ub = dbs.NewUpdateBuilder()
	ub.Table(this.groupTable)
	ub.SET("name", name)
	ub.SET("status", status)
	ub.SET("updated_on", time.Now())
	ub.Where("id = ?", id)
	ub.Where("ctx = ?", ctx)
	ub.Limit(1)
	_, err = ub.Exec(this.db)
	return err
}

func (this *odinRepository) UpdateGroupStatus(ctx, id int64, status int) (err error) {
	var ub = dbs.NewUpdateBuilder()
	ub.Table(this.groupTable)
	ub.SET("status", status)
	ub.SET("updated_on", time.Now())
	ub.Where("id = ?", id)
	ub.Where("ctx = ?", ctx)
	ub.Limit(1)
	_, err = ub.Exec(this.db)
	return err
}

func (this *odinRepository) RemoveGroup(ctx, id int64) (err error) {
	var rb = dbs.NewDeleteBuilder()
	rb.Table(this.groupTable)
	rb.Where("id = ?", id)
	rb.Where("ctx = ?", ctx)
	rb.Limit(1)
	_, err = rb.Exec(this.db)
	return err
}

func (this *odinRepository) getGroup(tx dbs.TX, ctx, id int64, gType int, name string) (result *odin.Group, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("g.id", "g.type", "g.status", "g.name", "g.created_on", "g.updated_on")
	sb.From(this.groupTable, "AS g")
	if id > 0 {
		sb.Where("g.id = ?", id)
	}

	sb.Where("(g.ctx = ? OR g.ctx = ?)", 0, ctx)

	if gType > 0 {
		sb.Where("g.type = ?", gType)
	}
	if name != "" {
		sb.Where("g.name = ?", name)
	}
	sb.Limit(1)
	if err = sb.Scan(tx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *odinRepository) GetPermissionTree(ctx, roleId int64, status int, name string) (result []*odin.Group, err error) {
	var tx = dbs.MustTx(this.db)

	if result, err = this.getGroupList(tx, ctx, odin.K_GROUP_TYPE_PERMISSION, status, name); err != nil {
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

func (this *odinRepository) GetPermissionList(ctx int64, groupIdList []int64, status int, keyword string) (result []*odin.Permission, err error) {
	var tx = dbs.MustTx(this.db)
	if result, err = this.getPermissionListWithGroupIdList(tx, ctx, 0, groupIdList, status, keyword); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *odinRepository) getPermissionListWithGroupIdList(tx dbs.TX, ctx, roleId int64, groupIdList []int64, status int, keyword string) (result []*odin.Permission, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("p.id", "p.ctx", "p.group_id", "p.name", "p.identifier", "p.status", "p.created_on")
	sb.From(this.permissionTable, "AS p")
	if roleId > 0 {
		sb.Selects("IF(rp.role_id IS NULL, false , true) AS granted")
		sb.LeftJoin(this.rolePermissionTable, "AS rp ON rp.permission_id = p.id AND rp.role_id = ?", roleId)
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
	sb.From(this.permissionTable, "AS p")
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

func (this *odinRepository) AddPermission(ctx int64, groupId int64, name, identifier string, status int) (result *odin.Permission, err error) {
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

func (this *odinRepository) insertPermission(tx dbs.TX, ctx, groupId int64, status int, name, identifier string) (id int64, err error) {
	var ib = dbs.NewInsertBuilder()
	ib.Table(this.permissionTable)
	ib.Columns("ctx", "group_id", "status", "name", "identifier", "created_on", "updated_on")
	ib.Values(ctx, groupId, status, name, identifier, time.Now(), time.Now())
	if result, err := ib.Exec(tx); err != nil {
		return 0, err
	} else {
		id, _ = result.LastInsertId()
	}
	return id, err
}

func (this *odinRepository) UpdatePermission(ctx, id, groupId int64, name, identifier string, status int) (err error) {
	var ub = dbs.NewUpdateBuilder()
	ub.Table(this.permissionTable)
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

func (this *odinRepository) UpdatePermissionStatus(ctx, id int64, status int) (err error) {
	var ub = dbs.NewUpdateBuilder()
	ub.Table(this.permissionTable)
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
	sb.From(this.permissionTable, "AS p")
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
	sb.From(this.permissionTable, "AS p")
	sb.LeftJoin(this.rolePermissionTable, "AS rp ON rp.permission_id = p.id")
	sb.Where("(p.ctx = ? OR p.ctx = ?)", 0, ctx)
	sb.Where("rp.role_id = ?", roleId)
	if err = sb.Scan(tx, &result); err != nil {
		return nil, err
	}
	return result, err
}

func (this *odinRepository) GetGrantedPermissionList(ctx int64, objectId string) (result []*odin.Permission, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("p.id", "p.ctx", "p.group_id", "p.name", "p.identifier", "p.status", "p.created_on", "p.updated_on")
	sb.Selects("IF(rp.role_id IS NULL, false, true) AS granted")
	sb.From(this.permissionTable, "AS p")
	sb.LeftJoin(this.rolePermissionTable, "AS rp ON rp.permission_id = p.id")
	sb.LeftJoin(this.roleGrantTable, "AS rg ON rg.role_id = rp.role_id")
	sb.Where("rg.object_id = ?", objectId)
	sb.Where("p.status = ?", odin.K_STATUS_ENABLE)
	sb.Where("(p.ctx = ? OR p.ctx = ?)", 0, ctx)
	sb.GroupBy("p.id")
	if err := sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, err
}

func (this *odinRepository) GetRoleTree(ctx int64, objectId string, status int, name string) (result []*odin.Group, err error) {
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

	rList, err := this.getRoles(tx, ctx, objectId, gIdList, status, "")
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

func (this *odinRepository) getRoles(tx dbs.TX, ctx int64, objectId string, groupIdList []int64, status int, keyword string) (result []*odin.Role, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("r.id", "r.ctx", "r.group_id", "r.name", "r.status", "r.created_on", "r.updated_on")
	sb.From(this.roleTable, "AS r")
	if objectId != "" {
		sb.Selects("IF(rg.object_id IS NULL, false, true) AS granted")
		sb.LeftJoin(this.roleGrantTable, "AS rg ON rg.role_id = r.id AND rg.object_id = ?", objectId)
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
	ib.Table(this.roleTable)
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
	ub.Table(this.roleTable)
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
	ub.Table(this.roleTable)
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
	sb.From(this.roleTable, "AS r")
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
	sb.From(this.roleTable, "AS r")
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

func (this *odinRepository) GrantPermission(ctx, roleId int64, permissionIdList []int64) (err error) {
	if len(permissionIdList) > 0 {
		var now = time.Now()
		var ib = dbs.NewInsertBuilder()
		ib.Table(this.rolePermissionTable)
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
	rb.Table(this.rolePermissionTable)
	rb.Where("role_id = ?", roleId)
	rb.Where("(ctx = ? OR ctx = ?)", 0, ctx)
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
	rb.Table(this.rolePermissionTable)
	rb.Where("role_id = ?", roleId)
	if _, err = rb.Exec(tx); err != nil {
		return err
	}

	if len(permissionIdList) > 0 {
		var ib = dbs.NewInsertBuilder()
		ib.Table(this.rolePermissionTable)
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

func (this *odinRepository) GrantRole(ctx int64, objectId string, roleIdList []int64) (err error) {
	var now = time.Now()
	var ib = dbs.NewInsertBuilder()
	ib.Table(this.roleGrantTable)
	ib.Options("IGNORE")
	ib.Columns("ctx", "object_id", "role_id", "created_on")
	for _, rId := range roleIdList {
		ib.Values(ctx, objectId, rId, now)
	}
	if _, err = ib.Exec(this.db); err != nil {
		return err
	}
	return nil
}

func (this *odinRepository) RevokeRole(ctx int64, objectId string, roleIdList []int64) (err error) {
	var rb = dbs.NewDeleteBuilder()
	rb.Table(this.roleGrantTable)
	rb.Where("object_id = ?", objectId)
	rb.Where("(ctx = ? OR ctx = ?)", 0, ctx)
	rb.Where(dbs.IN("role_id", roleIdList))
	if _, err = rb.Exec(this.db); err != nil {
		return err
	}
	return nil
}

func (this *odinRepository) ReGrantRole(ctx int64, objectId string, roleIdList []int64) (err error) {
	var tx = dbs.MustTx(this.db)
	var now = time.Now()

	var rb = dbs.NewDeleteBuilder()
	rb.Table(this.roleGrantTable)
	rb.Where("object_id = ?", objectId)
	if _, err = rb.Exec(tx); err != nil {
		return err
	}

	if len(roleIdList) > 0 {
		var ib = dbs.NewInsertBuilder()
		ib.Table(this.roleGrantTable)
		ib.Options("IGNORE")
		ib.Columns("ctx", "object_id", "role_id", "created_on")
		for _, rId := range roleIdList {
			ib.Values(ctx, objectId, rId, now)
		}
		if _, err = ib.Exec(tx); err != nil {
			return err
		}
	}
	tx.Commit()
	return nil
}

func (this *odinRepository) Check(ctx int64, objectId, identifier string) (result bool) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("rg.object_id", "rg.role_id", "rp.permission_id", "p.identifier")
	sb.From(this.roleGrantTable, "AS rg")
	sb.LeftJoin(this.rolePermissionTable, "AS rp ON rp.role_id = rg.role_id")
	sb.LeftJoin(this.permissionTable, "AS p ON p.id = rp.permission_id")
	sb.LeftJoin(this.roleTable, "AS r ON r.id = rg.role_id")
	sb.Where("(rg.ctx = ? OR rg.ctx = ?)", 0, ctx)
	sb.Where("rg.object_id = ? AND p.identifier = ? AND p.status = ? AND r.status = ?", objectId, identifier, odin.K_STATUS_ENABLE, odin.K_STATUS_ENABLE)
	sb.OrderBy("r.status", "p.status")
	sb.Limit(1)

	var grant *odin.Grant
	if err := sb.Scan(this.db, &grant); err != nil || grant == nil {
		return false
	}
	if grant.Identifier == identifier && grant.ObjectId == objectId {
		return true
	}
	return false
}

func (this *odinRepository) CheckList(ctx int64, objectId string, identifiers ...string) (result map[string]bool) {
	result = make(map[string]bool)
	var sb = dbs.NewSelectBuilder()
	sb.Selects("rg.object_id", "rg.role_id", "rp.permission_id", "p.identifier")
	sb.From(this.roleGrantTable, "AS rg")
	sb.LeftJoin(this.rolePermissionTable, "AS rp ON rp.role_id = rg.role_id")
	sb.LeftJoin(this.permissionTable, "AS p ON p.id = rp.permission_id")
	sb.LeftJoin(this.roleTable, "AS r ON r.id = rg.role_id")

	sb.Where("rg.object_id = ?", objectId)
	sb.Where("(rg.ctx = ? OR rg.ctx = ?)", 0, ctx)
	if len(identifiers) > 0 {
		var or = dbs.OR()
		for _, identifier := range identifiers {
			or.Append("p.identifier = ?", identifier)
			result[identifier] = false
		}
		sb.Where(or)
	}
	sb.Where("p.status = ?", odin.K_STATUS_ENABLE)
	sb.Where("r.status = ?", odin.K_STATUS_ENABLE)

	sb.GroupBy("p.id")
	sb.OrderBy("r.status", "p.status")
	sb.Limit(int64(len(identifiers)))

	var grantList []*odin.Grant
	if err := sb.Scan(this.db, &grantList); err != nil || grantList == nil {
		return result
	}
	for _, grant := range grantList {
		result[grant.Identifier] = true
	}
	return result
}

func (this *odinRepository) GetGrantedRoleList(ctx int64, objectId string) (result []*odin.Role, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("r.id", "r.group_id", "r.name", "r.status", "r.created_on", "r.updated_on")
	sb.From(this.roleTable, "AS r")
	sb.Selects("IF(rg.object_id IS NULL, false, true) AS granted")
	sb.LeftJoin(this.roleGrantTable, "AS rg ON rg.role_id = r.id")
	sb.Where("rg.object_id = ? AND r.status = ?", objectId, odin.K_STATUS_ENABLE)
	sb.Where("(rg.ctx = ? OR rg.ctx = ?)", 0, ctx)
	if err := sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, err
}

func (this *odinRepository) ClearCache(ctx int64, objectId string) {
}
