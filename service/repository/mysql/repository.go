package mysql

import (
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/odin"
	"github.com/smartwalle/odin/service/repository/internal/sql"
	"strings"
)

type repository struct {
	sql.Repository
}

func NewRepository(db dbs.DB, tablePrefix string) odin.Repository {
	var r = &repository{}
	r.Repository = sql.NewRepository(db, dbs.DialectMySQL, tablePrefix)
	return r
}

func (this *repository) BeginTx() (dbs.TX, odin.Repository) {
	var nRepo = *this
	var tx dbs.TX
	tx, nRepo.Repository = this.Repository.ExBeginTx()
	return tx, &nRepo
}

func (this *repository) WithTx(tx dbs.TX) odin.Repository {
	var nRepo = *this
	nRepo.Repository = this.Repository.ExWithTx(tx)
	return &nRepo
}

func (this *repository) InitTable() error {
	var rawText = "" +
		"CREATE TABLE IF NOT EXISTS `odin_grant` (" +
		"  `ctx` bigint(20) DEFAULT NULL," +
		"  `role_id` bigint(20) DEFAULT NULL," +
		"  `target` varchar(64) DEFAULT NULL," +
		"  `created_on` datetime DEFAULT NULL," +
		"  UNIQUE KEY `odin_grant_pk` (`ctx`,`role_id`,`target`)," +
		"  KEY `odin_grant_ctx_target_index` (`ctx`,`target`)," +
		"  KEY `odin_grant_role_id_index` (`role_id`)," +
		"  KEY `odin_grant_target_index` (`target`)" +
		") ENGINE=InnoDB;" +
		"" +
		"CREATE TABLE IF NOT EXISTS `odin_group` (" +
		"  `id` bigint(20) NOT NULL," +
		"  `ctx` bigint(20) DEFAULT NULL," +
		"  `type` int(2) DEFAULT NULL," +
		"  `name` varchar(64) DEFAULT NULL," +
		"  `alias_name` varchar(255) DEFAULT NULL," +
		"  `status` int(2) DEFAULT NULL," +
		"  `created_on` datetime DEFAULT NULL," +
		"  `updated_on` datetime DEFAULT NULL," +
		"  PRIMARY KEY (`id`)," +
		"  UNIQUE KEY `odin_group_id_uindex` (`id`)," +
		"  UNIQUE KEY `odin_group_pk` (`ctx`,`type`,`name`)" +
		") ENGINE=InnoDB;" +
		"" +
		"CREATE TABLE IF NOT EXISTS `odin_permission` (" +
		"  `id` bigint(20) NOT NULL," +
		"  `group_id` bigint(20) DEFAULT NULL," +
		"  `ctx` bigint(20) DEFAULT NULL," +
		"  `name` varchar(128) DEFAULT NULL," +
		"  `alias_name` varchar(255) DEFAULT NULL," +
		"  `status` int(2) DEFAULT '1'," +
		"  `description` varchar(1024) DEFAULT NULL," +
		"  `created_on` datetime DEFAULT NULL," +
		"  `updated_on` datetime DEFAULT NULL," +
		"  PRIMARY KEY (`id`)," +
		"  UNIQUE KEY `odin_permission_id_uindex` (`id`)," +
		"  UNIQUE KEY `odin_permission_ctx_name_uindex` (`ctx`,`name`)," +
		"  KEY `odin_permission_ctx_group_id_index` (`ctx`,`group_id`)," +
		"  KEY `odin_permission_ctx_index` (`ctx`)" +
		") ENGINE=InnoDB;" +
		"" +
		"CREATE TABLE IF NOT EXISTS `odin_pre_permission` (" +
		"  `ctx` bigint(20) NOT NULL," +
		"  `permission_id` bigint(20) NOT NULL," +
		"  `pre_permission_id` bigint(20) NOT NULL," +
		"  `auto_grant` tinyint(1) DEFAULT '0'," +
		"  `created_on` datetime DEFAULT NULL," +
		"  PRIMARY KEY (`ctx`,`permission_id`,`pre_permission_id`)," +
		"  KEY `odin_pre_permission_ctx_permission_id_index` (`ctx`,`permission_id`)," +
		"  KEY `odin_pre_permission_ctx_permission_id_pre_permission_id_index` (`ctx`,`permission_id`,`pre_permission_id`)" +
		") ENGINE=InnoDB;" +
		"" +
		"CREATE TABLE IF NOT EXISTS `odin_pre_role` (" +
		"  `ctx` bigint(20) NOT NULL," +
		"  `role_id` bigint(20) NOT NULL," +
		"  `pre_role_id` bigint(20) NOT NULL," +
		"  `created_on` datetime DEFAULT NULL," +
		"  PRIMARY KEY (`ctx`,`role_id`,`pre_role_id`)," +
		"  KEY `odin_pre_role_ctx_role_id_index` (`ctx`,`role_id`)," +
		"  KEY `odin_pre_role_ctx_role_id_pre_role_id_index` (`ctx`,`role_id`,`pre_role_id`)" +
		") ENGINE=InnoDB;" +
		"" +
		"CREATE TABLE IF NOT EXISTS `odin_role` (" +
		"  `id` bigint(20) NOT NULL," +
		"  `ctx` bigint(20) DEFAULT NULL," +
		"  `name` varchar(64) DEFAULT NULL," +
		"  `alias_name` varchar(255) DEFAULT NULL," +
		"  `status` int(2) DEFAULT '1'," +
		"  `description` varchar(1024) DEFAULT NULL," +
		"  `parent_id` bigint(20) DEFAULT NULL," +
		"  `left_value` bigint(20) DEFAULT NULL," +
		"  `right_value` bigint(20) DEFAULT NULL," +
		"  `depth` int(11) DEFAULT NULL," +
		"  `created_on` datetime DEFAULT NULL," +
		"  `updated_on` datetime DEFAULT NULL," +
		"  PRIMARY KEY (`id`)," +
		"  UNIQUE KEY `odin_role_id_uindex` (`id`)," +
		"  UNIQUE KEY `odin_role_ctx_name_uindex` (`ctx`,`name`)," +
		"  KEY `odin_role_ctx_index` (`ctx`)," +
		"  KEY `odin_role_ctx_left_value_index` (`ctx`,`left_value`)," +
		"  KEY `odin_role_ctx_parent_id_index` (`ctx`,`parent_id`)," +
		"  KEY `odin_role_ctx_right_value_index` (`ctx`,`right_value`)" +
		") ENGINE=InnoDB;" +
		"" +
		"CREATE TABLE IF NOT EXISTS `odin_role_mutex` (" +
		"  `ctx` bigint(20) NOT NULL," +
		"  `role_id` bigint(20) NOT NULL," +
		"  `mutex_role_id` bigint(20) NOT NULL," +
		"  `created_on` datetime DEFAULT NULL," +
		"  PRIMARY KEY (`ctx`,`role_id`,`mutex_role_id`)," +
		"  KEY `odin_role_mutex_ctx_mutex_role_id_index` (`ctx`,`mutex_role_id`)," +
		"  KEY `odin_role_mutex_ctx_mutex_role_id_role_id_index` (`ctx`,`mutex_role_id`,`role_id`)," +
		"  KEY `odin_role_mutex_ctx_role_id_index` (`ctx`,`role_id`)" +
		") ENGINE=InnoDB;" +
		"" +
		"CREATE TABLE IF NOT EXISTS `odin_role_permission` (" +
		"  `ctx` bigint(20) DEFAULT NULL," +
		"  `role_id` bigint(20) DEFAULT NULL," +
		"  `permission_id` bigint(20) DEFAULT NULL," +
		"  `created_on` datetime DEFAULT NULL," +
		"  UNIQUE KEY `odin_role_permission_pk` (`ctx`,`role_id`,`permission_id`)" +
		") ENGINE=InnoDB;"

	var sqlList = strings.Split(strings.ReplaceAll(rawText, "odin", this.TablePrefix()), ";")
	for _, sql := range sqlList {
		if sql == "" {
			continue
		}
		var rb = dbs.NewBuilder(sql)
		if _, err := rb.Exec(this.DB()); err != nil {
			return err
		}
	}
	return nil
}
