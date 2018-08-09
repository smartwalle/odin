package mysql

import (
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/odin"
	"github.com/smartwalle/odin/service"
	"strings"
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

func (this *odinRepository) ClearCache(ctx int64, objectId string) {
}
