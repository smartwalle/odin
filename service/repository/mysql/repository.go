package mysql

import (
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/odin"
	"strings"
)

type odinRepository struct {
	db dbs.DB

	tblPrefix         string
	tblGroup          string
	tblPermission     string
	tblRole           string
	tblRolePermission string
	tblGrant          string
	tblRoleMutex      string
	tblPreRole        string
	tblPrePermission  string
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
	r.tblPrefix = tblPrefix
	r.tblGroup = tblPrefix + "_group"
	r.tblPermission = tblPrefix + "_permission"
	r.tblRole = tblPrefix + "_role"
	r.tblRolePermission = tblPrefix + "_role_permission"
	r.tblGrant = tblPrefix + "_grant"
	r.tblRoleMutex = tblPrefix + "_role_mutex"
	r.tblPreRole = tblPrefix + "_pre_role"
	r.tblPrePermission = tblPrefix + "_pre_permission"
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
	return this.initMySQLTable()
}

func (this *odinRepository) initMySQLTable() error {
	var rawText = "" +
		"CREATE TABLE IF NOT EXISTS `odin_group` (" + // odin_group
		"`id` bigint(20) NOT NULL AUTO_INCREMENT," +
		"`ctx` int(11) DEFAULT NULL," +
		"`type` int(2) DEFAULT NULL," +
		"`name` varchar(64) DEFAULT NULL," +
		"`alias_name` varchar(255) DEFAULT NULL," +
		"`status` int(2) DEFAULT NULL," +
		"`created_on` datetime DEFAULT NULL," +
		"`updated_on` datetime DEFAULT NULL," +
		"PRIMARY KEY (`id`)," +
		"UNIQUE KEY `odin_group_id_uindex` (`id`)," +
		"UNIQUE KEY `odin_group_pk` (`ctx`,`type`,`name`)" +
		");" +
		"" +
		"CREATE TABLE IF NOT EXISTS `odin_permission` (" + // odin_permission
		"`id` bigint(20) NOT NULL AUTO_INCREMENT," +
		"`group_id` bigint(20) DEFAULT NULL," +
		"`ctx` int(11) DEFAULT NULL," +
		"`name` varchar(255) DEFAULT NULL," +
		"`alias_name` varchar(255) DEFAULT NULL," +
		"`status` int(2) DEFAULT '1'," +
		"`description` varchar(1024)," +
		"`created_on` datetime DEFAULT NULL," +
		"`updated_on` datetime DEFAULT NULL," +
		"PRIMARY KEY (`id`)," +
		"UNIQUE KEY `odin_permission_id_uindex` (`id`)," +
		"UNIQUE KEY `odin_permission_ctx_name_uindex` (`ctx`,`name`)," +
		"KEY `odin_permission_ctx_index` (`ctx`)," +
		"KEY `odin_permission_ctx_group_id_index` (`ctx`,`group_id`)" +
		");" +
		"" +
		"CREATE TABLE IF NOT EXISTS `odin_role` (" + // odin_role
		"`id` bigint(20) NOT NULL AUTO_INCREMENT," +
		"`ctx` int(11) DEFAULT NULL," +
		"`name` varchar(64) DEFAULT NULL," +
		"`alias_name` varchar(255) DEFAULT NULL," +
		"`status` int(2) DEFAULT '1'," +
		"`description` varchar(1024) DEFAULT NULL," +
		"`parent_id` bigint(20) DEFAULT NULL," +
		"`left_value` int(11) DEFAULT NULL," +
		"`right_value` int(11) DEFAULT NULL," +
		"`depth` int(11) DEFAULT NULL," +
		"`created_on` datetime DEFAULT NULL," +
		"`updated_on` datetime DEFAULT NULL," +
		"PRIMARY KEY (`id`)," +
		"UNIQUE KEY `odin_role_id_uindex` (`id`)," +
		"UNIQUE KEY `odin_role_ctx_name_uindex` (`ctx`,`name`)," +
		"KEY `odin_role_ctx_index` (`ctx`)," +
		"KEY `odin_role_ctx_parent_id_index` (`ctx`,`parent_id`)," +
		"KEY `odin_role_ctx_left_value_index` (`ctx`,`left_value`)," +
		"KEY `odin_role_ctx_right_value_index` (`ctx`,`right_value`)" +
		");" +
		"" +
		"CREATE TABLE IF NOT EXISTS `odin_role_permission` (" + // odin_role_permission
		"`ctx` int(11) DEFAULT NULL," +
		"`role_id` bigint(20) DEFAULT NULL," +
		"`permission_id` bigint(20) DEFAULT NULL," +
		"`created_on` datetime DEFAULT NULL," +
		"UNIQUE KEY `odin_role_permission_pk` (`ctx`,`role_id`,`permission_id`)" +
		");" +
		"" +
		"CREATE TABLE IF NOT EXISTS `odin_grant` (" + // odin_grant
		"`ctx` int(11) DEFAULT NULL," +
		"`role_id` bigint(20) DEFAULT NULL," +
		"`target` varchar(64) DEFAULT NULL," +
		"`created_on` datetime DEFAULT NULL," +
		"UNIQUE KEY `odin_grant_pk` (`ctx`,`role_id`,`target`)," +
		"KEY `odin_grant_ctx_target_index` (`ctx`,`target`)," +
		"KEY `odin_grant_role_id_index` (`role_id`)," +
		"KEY `odin_grant_target_index` (`target`)" +
		");" +
		"" +
		"CREATE TABLE IF NOT EXISTS `odin_role_mutex` (" + // odin_role_mutex
		"`ctx` int(11) NOT NULL," +
		"`role_id` bigint(20) NOT NULL," +
		"`mutex_role_id` bigint(20) NOT NULL," +
		"`created_on` datetime DEFAULT NULL," +
		"PRIMARY KEY (`ctx`,`role_id`,`mutex_role_id`)," +
		"KEY `odin_role_mutex_ctx_mutex_role_id_index` (`ctx`,`mutex_role_id`)," +
		"KEY `odin_role_mutex_ctx_role_id_index` (`ctx`,`role_id`)," +
		"KEY `odin_role_mutex_ctx_mutex_role_id_role_id_index` (`ctx`,`mutex_role_id`,`role_id`)" +
		");" +
		"" +
		"CREATE TABLE IF NOT EXISTS `odin_pre_role` (" + // odin_pre_role
		"`ctx` int(11) NOT NULL," +
		"`role_id` bigint(20) NOT NULL," +
		"`pre_role_id` bigint(20) NOT NULL," +
		"`created_on` datetime DEFAULT NULL," +
		"PRIMARY KEY (`ctx`,`role_id`,`pre_role_id`)," +
		"KEY `odin_pre_role_ctx_role_id_pre_role_id_index` (`ctx`,`role_id`,`pre_role_id`)," +
		"KEY `odin_pre_role_ctx_role_id_index` (`ctx`,`role_id`)" +
		");" +
		"" +
		"CREATE TABLE IF NOT EXISTS `odin_pre_permission` (" + // odin_pre_permission
		"`ctx` int(11) NOT NULL," +
		"`permission_id` bigint(20) NOT NULL," +
		"`pre_permission_id` bigint(20) NOT NULL," +
		"`auto_grant` tinyint(1) DEFAULT '0'," +
		"`created_on` datetime DEFAULT NULL," +
		"PRIMARY KEY (`ctx`,`permission_id`,`pre_permission_id`)," +
		"KEY `odin_pre_permission_ctx_permission_id_index` (`ctx`,`permission_id`)," +
		"KEY `odin_pre_permission_ctx_permission_id_pre_permission_id_index` (`ctx`,`permission_id`,`pre_permission_id`)" +
		");"

	var sqlList = strings.Split(strings.ReplaceAll(rawText, "odin", this.tblPrefix), ";")
	for _, sql := range sqlList {
		if sql == "" {
			continue
		}
		var rb = dbs.NewBuilder(sql)
		if _, err := rb.Exec(this.db); err != nil {
			return err
		}
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

func (this *odinRepository) CheckRoleAccessible(ctx int64, target string, roleName string) bool {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("r.id", "r.ctx", "r.name", "r.alias_name", "r.status", "r.description", "r.parent_id", "r.left_value", "r.right_value", "r.depth", "r.created_on", "r.updated_on")
	sb.Selects("MAX(IF(rg.role_id = r.id, false, true)) AS can_access")
	sb.From(this.tblRole, "AS r")
	sb.LeftJoin(this.tblRole, "AS rp ON rp.left_value < r.left_value AND rp.right_value > r.right_value")
	sb.LeftJoin(this.tblGrant, "AS rg ON rg.role_id = rp.id")
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

func (this *odinRepository) CheckRoleAccessibleWithId(ctx int64, target string, roleId int64) bool {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("r.id", "r.ctx", "r.name", "r.alias_name", "r.status", "r.description", "r.parent_id", "r.left_value", "r.right_value", "r.depth", "r.created_on", "r.updated_on")
	sb.Selects("MAX(IF(rg.role_id = r.id, false, true)) AS can_access")
	sb.From(this.tblRole, "AS r")
	sb.LeftJoin(this.tblRole, "AS rp ON rp.left_value < r.left_value AND rp.right_value > r.right_value")
	sb.LeftJoin(this.tblGrant, "AS rg ON rg.role_id = rp.id")
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

func (this *odinRepository) CheckRolePermission(ctx int64, roleName, permissionName string) bool {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("rp.ctx", "rp.role_id", "rp.permission_id")
	sb.From(this.tblRolePermission, "AS rp")
	sb.LeftJoin(this.tblRole, "AS r ON r.id = rp.role_id")
	sb.LeftJoin(this.tblPermission, "AS p ON p.id = rp.permission_id")
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

func (this *odinRepository) CheckRolePermissionWithId(ctx, roleId, permissionId int64) bool {
	var sb = dbs.NewSelectBuilder()
	sb.Selects("rp.ctx", "rp.role_id", "rp.permission_id")
	sb.From(this.tblRolePermission, "AS rp")
	sb.LeftJoin(this.tblRole, "AS r ON r.id = rp.role_id")
	sb.LeftJoin(this.tblPermission, "AS p ON p.id = rp.permission_id")
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

func (this *odinRepository) CleanCache(ctx int64, target string) {
}
