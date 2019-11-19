package mysql

import (
	"fmt"
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/odin"
	"github.com/smartwalle/odin/service"
	"strings"
)

type odinRepository struct {
	db                dbs.DB
	tblGroup          string
	tblPermission     string
	tblRole           string
	tblRolePermission string
	tblGrant          string
}

func NewOdinRepository(db dbs.DB, tblPrefix string) service.OdinRepository {
	var r = &odinRepository{}
	r.db = db

	tblPrefix = strings.TrimSpace(tblPrefix)
	if tblPrefix != "" {
		tblPrefix = tblPrefix + "_odin"
	} else {
		tblPrefix = "odin"
	}

	r.db = db
	r.tblGroup = tblPrefix + "_group"
	r.tblPermission = tblPrefix + "_permission"
	r.tblRole = tblPrefix + "_role"
	r.tblRolePermission = tblPrefix + "_role_permission"
	r.tblGrant = tblPrefix + "_grant"

	if err := r.initTable(); err != nil {
		panic(fmt.Sprintf("初始化 Odin 失败, 错误信息为: %v", err))
	}

	return r
}

func (this *odinRepository) initTable() error {
	var tx = dbs.MustTx(this.db)

	var cb = dbs.NewBuilder("")
	cb.Format(odinGroupSQL, this.tblGroup)
	if _, err := cb.Exec(tx); err != nil {
		return err
	}

	cb = dbs.NewBuilder("")
	cb.Format(odinPermissionSQL, this.tblPermission)
	if _, err := cb.Exec(tx); err != nil {
		return err
	}

	cb = dbs.NewBuilder("")
	cb.Format(odinRoleSQL, this.tblRole)
	if _, err := cb.Exec(tx); err != nil {
		return err
	}

	cb = dbs.NewBuilder("")
	cb.Format(odinRolePermissionSQL, this.tblRolePermission)
	if _, err := cb.Exec(tx); err != nil {
		return err
	}

	cb = dbs.NewBuilder("")
	cb.Format(odinGrantSQL, this.tblGrant)
	if _, err := cb.Exec(tx); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinRepository) Check(ctx int64, targetId, identifier string) (result bool) {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("rg.target_id", "rg.role_id", "rp.permission_id", "p.identifier")
	sb.From(this.tblGrant, "AS rg")
	sb.LeftJoin(this.tblRolePermission, "AS rp ON rp.role_id = rg.role_id")
	sb.LeftJoin(this.tblPermission, "AS p ON p.id = rp.permission_id")
	sb.LeftJoin(this.tblRole, "AS r ON r.id = rg.role_id")
	sb.Where("(rg.ctx = ? OR rg.ctx = ?)", 0, ctx)
	sb.Where("rg.target_id = ? AND p.identifier = ? AND p.status = ? AND r.status = ?", targetId, identifier, odin.K_STATUS_ENABLE, odin.K_STATUS_ENABLE)
	sb.OrderBy("r.status", "p.status")
	sb.Limit(1)

	var grant *odin.Grant
	if err := sb.Scan(this.db, &grant); err != nil || grant == nil {
		return false
	}
	if grant.Identifier == identifier && grant.TargetId == targetId {
		return true
	}
	return false
}

func (this *odinRepository) CheckList(ctx int64, targetId string, identifiers ...string) (result map[string]bool) {
	result = make(map[string]bool)
	var sb = dbs.NewSelectBuilder()
	sb.Selects("rg.target_id", "rg.role_id", "rp.permission_id", "p.identifier")
	sb.From(this.tblGrant, "AS rg")
	sb.LeftJoin(this.tblRolePermission, "AS rp ON rp.role_id = rg.role_id")
	sb.LeftJoin(this.tblPermission, "AS p ON p.id = rp.permission_id")
	sb.LeftJoin(this.tblRole, "AS r ON r.id = rg.role_id")

	sb.Where("rg.target_id = ?", targetId)
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

func (this *odinRepository) ClearCache(ctx int64, targetId string) {
}
