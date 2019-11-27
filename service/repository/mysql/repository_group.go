package mysql

import (
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/odin"
	"time"
)

func (this *odinRepository) GetGroupList(ctx int64, gType odin.GroupType, status odin.Status, keywords string) (result []*odin.Group, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("g.id", "g.ctx", "g.type", "g.name", "g.alias_name", "g.status", "g.created_on", "g.updated_on")
	sb.From(this.tblGroup, "AS g")
	sb.Where("g.ctx = ?", ctx)
	if status != 0 {
		sb.Where("g.status = ?", status)
	}
	if keywords != "" {
		var or = dbs.OR()
		or.Append(dbs.Like("g.name", "%", keywords, "%"))
		or.Append(dbs.Like("g.alias_name", "%", keywords, "%"))
		sb.Where(or)
	}
	sb.OrderBy("g.ctx", "g.id")
	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *odinRepository) getGroup(ctx int64, gType odin.GroupType, groupId int64, name string) (result *odin.Group, err error) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("g.id", "g.ctx", "g.type", "g.name", "g.alias_name", "g.status", "g.created_on", "g.updated_on")
	sb.From(this.tblGroup, "AS g")
	sb.Where("g.ctx = ?", ctx)
	sb.Where("g.type = ?", gType)
	if groupId > 0 {
		sb.Where("g.id = ?", groupId)
	}
	if name != "" {
		sb.Where("g.name = ?", name)
	}
	sb.Limit(1)
	if err = sb.Scan(this.db, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *odinRepository) GetGroupWithId(ctx int64, gType odin.GroupType, groupId int64) (result *odin.Group, err error) {
	return this.getGroup(ctx, gType, groupId, "")
}

func (this *odinRepository) GetGroupWithName(ctx int64, gType odin.GroupType, name string) (result *odin.Group, err error) {
	return this.getGroup(ctx, gType, 0, name)
}

func (this *odinRepository) AddGroup(ctx int64, gType odin.GroupType, name, aliasName string, status odin.Status) (result int64, err error) {
	var now = time.Now()
	var ib = dbs.NewInsertBuilder()
	ib.Table(this.tblGroup)
	ib.Columns("ctx", "type", "name", "alias_name", "status", "created_on", "updated_on")
	ib.Values(ctx, gType, name, aliasName, status, now, now)
	rResult, err := ib.Exec(this.db)
	if err != nil {
		return 0, err
	}
	if result, err = rResult.LastInsertId(); err != nil {
		return 0, err
	}
	return result, nil
}

func (this *odinRepository) UpdateGroup(ctx int64, gType odin.GroupType, groupId int64, aliasName string, status odin.Status) (err error) {
	var now = time.Now()
	var ub = dbs.NewUpdateBuilder()
	ub.Table(this.tblGroup)
	ub.SET("alias_name", aliasName)
	ub.SET("status", status)
	ub.SET("updated_on", now)
	ub.Where("id = ?", groupId)
	ub.Where("ctx = ?", ctx)
	ub.Where("type = ?", gType)
	ub.Limit(1)
	_, err = ub.Exec(this.db)
	return err
}

func (this *odinRepository) UpdateGroupStatus(ctx int64, gType odin.GroupType, groupId int64, status odin.Status) (err error) {
	var now = time.Now()
	var ub = dbs.NewUpdateBuilder()
	ub.Table(this.tblGroup)
	ub.SET("status", status)
	ub.SET("updated_on", now)
	ub.Where("id = ?", groupId)
	ub.Where("ctx = ?", ctx)
	ub.Where("type = ?", gType)
	ub.Limit(1)
	_, err = ub.Exec(this.db)
	return err
}
