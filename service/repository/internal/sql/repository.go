package sql

import (
	"errors"
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/odin"
	"strings"
)

type Repository struct {
	db                  dbs.DB
	dialect             dbs.Dialect
	idGenerator         dbs.IdGenerator
	tablePrefix         string
	tableGroup          string
	tablePermission     string
	tableRole           string
	tableRolePermission string
	tableGrant          string
	tableRoleMutex      string
	tablePreRole        string
	tablePrePermission  string
}

func NewRepository(db dbs.DB, dialect dbs.Dialect, tblPrefix string) *Repository {
	var r = &Repository{}
	r.db = db
	r.dialect = dialect
	r.idGenerator = dbs.GetIdGenerator()

	tblPrefix = strings.TrimSpace(tblPrefix)
	if tblPrefix == "" {
		tblPrefix = "odin"
	} else {
		tblPrefix = tblPrefix + "_odin"
	}
	r.tablePrefix = tblPrefix
	r.tableGroup = tblPrefix + "_group"
	r.tablePermission = tblPrefix + "_permission"
	r.tableRole = tblPrefix + "_role"
	r.tableRolePermission = tblPrefix + "_role_permission"
	r.tableGrant = tblPrefix + "_grant"
	r.tableRoleMutex = tblPrefix + "_role_mutex"
	r.tablePreRole = tblPrefix + "_pre_role"
	r.tablePrePermission = tblPrefix + "_pre_permission"
	return r
}

func (this *Repository) DB() dbs.DB {
	return this.db
}

func (this *Repository) Dialect() dbs.Dialect {
	return this.dialect
}

func (this *Repository) BeginTx() (dbs.TX, odin.Repository) {
	var tx = dbs.MustTx(this.db)
	var nRepo = *this
	nRepo.db = tx
	return tx, &nRepo
}

func (this *Repository) WithTx(tx dbs.TX) odin.Repository {
	var nRepo = *this
	nRepo.db = tx
	return &nRepo
}

func (this *Repository) ExBeginTx() (dbs.TX, *Repository) {
	var tx = dbs.MustTx(this.db)
	var nRepo = *this
	nRepo.db = tx
	return tx, &nRepo
}

func (this *Repository) ExWithTx(tx dbs.TX) *Repository {
	var nRepo = *this
	nRepo.db = tx
	return &nRepo
}

func (this *Repository) UseIdGenerator(g dbs.IdGenerator) {
	this.idGenerator = g
}

func (this *Repository) IdGenerator() dbs.IdGenerator {
	return this.idGenerator
}

func (this *Repository) TablePrefix() string {
	return this.tablePrefix
}

func (this *Repository) TableGroup() string {
	return this.tableGroup
}

func (this *Repository) TablePermission() string {
	return this.tablePermission
}

func (this *Repository) TableRole() string {
	return this.tableRole
}

func (this *Repository) TableRolePermission() string {
	return this.tableRolePermission
}

func (this *Repository) TableGrant() string {
	return this.tableGrant
}

func (this *Repository) TableRoleMutex() string {
	return this.tableRoleMutex
}

func (this *Repository) TablePreRole() string {
	return this.tablePreRole
}

func (this *Repository) TablePrePermission() string {
	return this.tablePrePermission
}

func (this *Repository) InitTable() error {
	return errors.New("odin: not implemented this method")
}

func (this *Repository) CheckPermission(ctx int64, target string, permissionName string) bool {
	var sb = dbs.NewSelectBuilder()
	sb.UseDialect(this.dialect)
	sb.Selects("g.ctx", "g.target", "g.role_id")
	sb.Selects("r.name AS role_name")
	sb.Selects("p.id AS permission_id", "p.name AS permission_name")
	sb.From(this.tableGrant, "AS g")
	sb.LeftJoin(this.tableRole, "AS r ON r.id = g.role_id")
	sb.LeftJoin(this.tableRolePermission, "AS rp ON rp.role_id = r.id")
	sb.LeftJoin(this.tablePermission, "AS p ON p.id = rp.permission_id")
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

func (this *Repository) CheckPermissionWithId(ctx int64, target string, permissionId int64) bool {
	var sb = dbs.NewSelectBuilder()
	sb.UseDialect(this.dialect)
	sb.Selects("g.ctx", "g.target", "g.role_id")
	sb.Selects("r.name AS role_name")
	sb.Selects("p.id AS permission_id", "p.name AS permission_name")
	sb.From(this.tableGrant, "AS g")
	sb.LeftJoin(this.tableRole, "AS r ON r.id = g.role_id")
	sb.LeftJoin(this.tableRolePermission, "AS rp ON rp.role_id = r.id")
	sb.LeftJoin(this.tablePermission, "AS p ON p.id = rp.permission_id")
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

func (this *Repository) CheckRole(ctx int64, target string, roleName string) bool {
	var sb = dbs.NewSelectBuilder()
	sb.UseDialect(this.dialect)
	sb.Selects("g.ctx", "g.target", "g.role_id")
	sb.Selects("r.name AS role_name")
	sb.From(this.tableGrant, "AS g")
	sb.LeftJoin(this.tableRole, "AS r ON r.id = g.role_id")
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

func (this *Repository) CheckRoleWithId(ctx int64, target string, roleId int64) bool {
	var sb = dbs.NewSelectBuilder()
	sb.UseDialect(this.dialect)
	sb.Selects("g.ctx", "g.target", "g.role_id")
	sb.Selects("r.name AS role_name")
	sb.From(this.tableGrant, "AS g")
	sb.LeftJoin(this.tableRole, "AS r ON r.id = g.role_id")
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

func (this *Repository) CheckRoleAccessible(ctx int64, target string, roleName string) bool {
	var sb = dbs.NewSelectBuilder()
	sb.UseDialect(this.dialect)
	sb.Selects("r.id", "r.ctx", "r.name", "r.alias_name", "r.status", "r.description", "r.parent_id", "r.left_value", "r.right_value", "r.depth", "r.created_on", "r.updated_on")
	sb.Selects("MAX(CASE WHEN rg.role_id = r.id THEN 0 ELSE 1 END) AS can_access")
	sb.From(this.tableRole, "AS r")
	sb.LeftJoin(this.tableRole, "AS rp ON rp.left_value < r.left_value AND rp.right_value > r.right_value")
	sb.LeftJoin(this.tableGrant, "AS rg ON rg.role_id = rp.id")
	sb.Where("rg.ctx = ?", ctx)
	sb.Where("rg.target = ?", target)
	sb.Where("rp.ctx = ?", ctx)
	sb.Where("rp.status = ?", odin.Enable)
	sb.Where("r.ctx = ?", ctx)
	sb.Where("r.status = ?", odin.Enable)
	sb.Where("r.name = ?", roleName)
	sb.GroupBy("r.ctx", "r.id")
	sb.OrderBy("r.ctx", "r.id")
	sb.Limit(1)
	var role *odin.Role
	if err := sb.Scan(this.db, &role); err != nil {
		return false
	}
	if role == nil {
		return false
	}
	return role.Accessible
}

func (this *Repository) CheckRoleAccessibleWithId(ctx int64, target string, roleId int64) bool {
	var sb = dbs.NewSelectBuilder()
	sb.UseDialect(this.dialect)
	sb.Selects("r.id", "r.ctx", "r.name", "r.alias_name", "r.status", "r.description", "r.parent_id", "r.left_value", "r.right_value", "r.depth", "r.created_on", "r.updated_on")
	sb.Selects("MAX(CASE WHEN rg.role_id = r.id THEN 0 ELSE 1 END) AS can_access")
	sb.From(this.tableRole, "AS r")
	sb.LeftJoin(this.tableRole, "AS rp ON rp.left_value < r.left_value AND rp.right_value > r.right_value")
	sb.LeftJoin(this.tableGrant, "AS rg ON rg.role_id = rp.id")
	sb.Where("rg.ctx = ?", ctx)
	sb.Where("rg.target = ?", target)
	sb.Where("rp.ctx = ?", ctx)
	sb.Where("rp.status = ?", odin.Enable)
	sb.Where("r.ctx = ?", ctx)
	sb.Where("r.status = ?", odin.Enable)
	sb.Where("r.id = ?", roleId)
	sb.GroupBy("r.ctx", "r.id")
	sb.OrderBy("r.ctx", "r.id")
	sb.Limit(1)
	var role *odin.Role
	if err := sb.Scan(this.db, &role); err != nil {
		return false
	}
	if role == nil {
		return false
	}
	return role.Accessible
}

func (this *Repository) CheckRolePermission(ctx int64, roleName, permissionName string) bool {
	var sb = dbs.NewSelectBuilder()
	sb.UseDialect(this.dialect)
	sb.Selects("rp.ctx", "rp.role_id", "rp.permission_id")
	sb.From(this.tableRolePermission, "AS rp")
	sb.LeftJoin(this.tableRole, "AS r ON r.id = rp.role_id")
	sb.LeftJoin(this.tablePermission, "AS p ON p.id = rp.permission_id")
	sb.Where("rp.ctx = ?", ctx)
	sb.Where("r.ctx = ? AND r.name = ?", ctx, roleName)
	sb.Where("p.ctx = ? AND p.name = ?", ctx, permissionName)
	sb.Limit(1)
	var rp *odin.RolePermission
	if err := sb.Scan(this.db, &rp); err != nil {
		return false
	}
	if rp == nil {
		return false
	}
	return true
}

func (this *Repository) CheckRolePermissionWithId(ctx, roleId, permissionId int64) bool {
	var sb = dbs.NewSelectBuilder()
	sb.UseDialect(this.dialect)
	sb.Selects("rp.ctx", "rp.role_id", "rp.permission_id")
	sb.From(this.tableRolePermission, "AS rp")
	sb.LeftJoin(this.tableRole, "AS r ON r.id = rp.role_id")
	sb.LeftJoin(this.tablePermission, "AS p ON p.id = rp.permission_id")
	sb.Where("rp.ctx = ?", ctx)
	sb.Where("r.ctx = ? AND r.id = ?", ctx, roleId)
	sb.Where("p.ctx = ? AND p.id = ?", ctx, permissionId)
	sb.Limit(1)
	var rp *odin.RolePermission
	if err := sb.Scan(this.db, &rp); err != nil {
		return false
	}
	if rp == nil {
		return false
	}
	return true
}

func (this *Repository) CleanCache(ctx int64, target string) {
}
