package mysql

import (
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/odin"
	"strings"
)

type odinRepository struct {
	db dbs.DB

	tblGroup          string
	tblPermission     string
	tblRole           string
	tblRolePermission string
	tblGrant          string
}

func NewRepository(db dbs.DB, tblPrefix string) odin.Repository {
	var r = &odinRepository{}
	r.db = db

	tblPrefix = strings.TrimSpace(tblPrefix)
	if tblPrefix == "" {
		tblPrefix = "odin"
	} else {
		tblPrefix = tblPrefix + "_odin"
	}

	r.tblGroup = tblPrefix + "_group"
	r.tblPermission = tblPrefix + "_permission"
	r.tblRole = tblPrefix + "_role"
	r.tblRolePermission = tblPrefix + "_role_permission"
	r.tblGrant = tblPrefix + "_grant"

	return r
}

func (this *odinRepository) BeginTx() (dbs.TX, odin.Repository) {
	var tx = dbs.MustTx(this.db)
	var nRepo = *this
	nRepo.db = tx
	return tx, &nRepo
}

func (this *odinRepository) WithTx(tx dbs.TX) odin.Repository {
	var nRepo = *this
	nRepo.db = tx
	return &nRepo
}

func (this *odinRepository) InitTable() error {
	var tx = dbs.MustTx(this.db)

	var cb *dbs.RawBuilder

	cb = dbs.NewBuilder("")
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

func (this *odinRepository) Check(ctx int64, targetId string, permissionName string) bool {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("g.ctx", "g.target_id", "g.role_id")
	sb.Selects("r.name AS role_name")
	sb.Selects("p.id AS permission_id", "p.name AS permission_name")
	sb.From(this.tblGrant, "AS g")
	sb.LeftJoin(this.tblRole, "AS r ON r.id = g.role_id")
	sb.LeftJoin(this.tblRolePermission, "AS rp ON rp.role_id = r.id")
	sb.LeftJoin(this.tblPermission, "AS p ON p.id = rp.permission_id")
	sb.Where("g.ctx = ? AND g.target_id = ?", ctx, targetId)
	sb.Where("r.ctx = ? AND r.status = ?", ctx, odin.Enable)
	sb.Where("rp.ctx = ?", ctx)
	sb.Where("p.ctx = ? AND p.name = ? AND p.status = ?", ctx, permissionName, odin.Enable)
	sb.Limit(1)
	var grant *odin.Grant
	if err := sb.Scan(this.db, &grant); err != nil {
		return false
	}
	if grant != nil {
		return true
	}
	return false
}

func (this *odinRepository) CleanCache(ctx int64, targetId string) {
}
