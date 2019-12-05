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
	var rb *dbs.RawBuilder

	rb = dbs.NewBuilder("")
	rb.Format(odinGroupSQL, this.tblGroup)
	if _, err := rb.Exec(this.db); err != nil {
		return err
	}

	rb = dbs.NewBuilder("")
	rb.Format(odinPermissionSQL, this.tblPermission)
	if _, err := rb.Exec(this.db); err != nil {
		return err
	}

	rb = dbs.NewBuilder("")
	rb.Format(odinRoleSQL, this.tblRole)
	if _, err := rb.Exec(this.db); err != nil {
		return err
	}

	rb = dbs.NewBuilder("")
	rb.Format(odinRolePermissionSQL, this.tblRolePermission)
	if _, err := rb.Exec(this.db); err != nil {
		return err
	}

	rb = dbs.NewBuilder("")
	rb.Format(odinGrantSQL, this.tblGrant)
	if _, err := rb.Exec(this.db); err != nil {
		return err
	}
	return nil
}

func (this *odinRepository) CheckPermission(ctx int64, target string, permissionName string) bool {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("g.ctx", "g.target", "g.role_id")
	sb.Selects("r.name AS role_name")
	sb.Selects("p.id AS permission_id", "p.name AS permission_name")
	sb.From(this.tblGrant, "AS g")
	sb.LeftJoin(this.tblRole, "AS r ON r.id = g.role_id")
	sb.LeftJoin(this.tblRolePermission, "AS rp ON rp.role_id = r.id")
	sb.LeftJoin(this.tblPermission, "AS p ON p.id = rp.permission_id")
	sb.Where("g.ctx = ? AND g.target = ?", ctx, target)
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

func (this *odinRepository) CheckPermissionWithId(ctx int64, target string, permissionId int64) bool {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("g.ctx", "g.target", "g.role_id")
	sb.Selects("r.name AS role_name")
	sb.Selects("p.id AS permission_id", "p.name AS permission_name")
	sb.From(this.tblGrant, "AS g")
	sb.LeftJoin(this.tblRole, "AS r ON r.id = g.role_id")
	sb.LeftJoin(this.tblRolePermission, "AS rp ON rp.role_id = r.id")
	sb.LeftJoin(this.tblPermission, "AS p ON p.id = rp.permission_id")
	sb.Where("g.ctx = ? AND g.target = ?", ctx, target)
	sb.Where("r.ctx = ? AND r.status = ?", ctx, odin.Enable)
	sb.Where("rp.ctx = ?", ctx)
	sb.Where("p.ctx = ? AND p.id = ? AND p.status = ?", ctx, permissionId, odin.Enable)
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

func (this *odinRepository) CheckRole(ctx int64, target string, roleName string) bool {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("g.ctx", "g.target", "g.role_id")
	sb.Selects("r.name AS role_name")
	sb.From(this.tblGrant, "AS g")
	sb.LeftJoin(this.tblRole, "AS r ON r.id = g.role_id")
	sb.Where("g.ctx = ? AND g.target = ?", ctx, target)
	sb.Where("r.ctx = ? AND r.name = ? AND r.status = ?", ctx, roleName, odin.Enable)
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

func (this *odinRepository) CheckRoleWithId(ctx int64, target string, roleId int64) bool {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("g.ctx", "g.target", "g.role_id")
	sb.Selects("r.name AS role_name")
	sb.From(this.tblGrant, "AS g")
	sb.LeftJoin(this.tblRole, "AS r ON r.id = g.role_id")
	sb.Where("g.ctx = ? AND g.role_id = ? AND g.target = ?", ctx, roleId, target)
	sb.Where("r.ctx = ? AND r.status = ?", ctx, odin.Enable)
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

func (this *odinRepository) CleanCache(ctx int64, target string) {
}
