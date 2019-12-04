package mysql

const odinGroupSQL = "" +
	"CREATE TABLE IF NOT EXISTS `%s` (" +
	"`id` bigint(20) NOT NULL AUTO_INCREMENT," +
	"`ctx` bigint(20) DEFAULT NULL," +
	"`type` int(2) DEFAULT NULL," +
	"`name` varchar(64) DEFAULT NULL," +
	"`alias_name` varchar(255) DEFAULT NULL," +
	"`status` int(2) DEFAULT NULL," +
	"`created_on` datetime DEFAULT NULL," +
	"`updated_on` datetime DEFAULT NULL," +
	"PRIMARY KEY (`id`)," +
	"UNIQUE KEY `odin_group_id_uindex` (`id`)," +
	"UNIQUE KEY `odin_group_pk` (`ctx`,`type`,`name`)" +
	") ENGINE=InnoDB DEFAULT CHARSET=utf8;"

const odinPermissionSQL = "" +
	"CREATE TABLE IF NOT EXISTS `%s` (" +
	"`id` bigint(20) NOT NULL AUTO_INCREMENT," +
	"`group_id` bigint(20) DEFAULT NULL," +
	"`ctx` bigint(20) DEFAULT NULL," +
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
	") ENGINE=InnoDB DEFAULT CHARSET=utf8;"

const odinRoleSQL = "" +
	"CREATE TABLE IF NOT EXISTS `%s` (" +
	"`id` bigint(20) NOT NULL AUTO_INCREMENT," +
	"`ctx` bigint(20) DEFAULT NULL," +
	"`name` varchar(64) DEFAULT NULL," +
	"`alias_name` varchar(255) DEFAULT NULL," +
	"`status` int(2) DEFAULT '1'," +
	"`description` varchar(1024) DEFAULT NULL," +
	"`parent_id` bigint(20) DEFAULT NULL," +
	"`created_on` datetime DEFAULT NULL," +
	"`updated_on` datetime DEFAULT NULL," +
	"PRIMARY KEY (`id`)," +
	"UNIQUE KEY `odin_role_id_uindex` (`id`)," +
	"UNIQUE KEY `odin_role_ctx_name_uindex` (`ctx`,`name`)," +
	"KEY `odin_role_ctx_index` (`ctx`)," +
	"KEY `odin_role_ctx_parent_id_index` (`ctx`,`parent_id`)" +
	") ENGINE=InnoDB DEFAULT CHARSET=utf8;"

const odinRolePermissionSQL = "" +
	"CREATE TABLE IF NOT EXISTS `%s` (" +
	"`ctx` bigint(20) DEFAULT NULL," +
	"`role_id` bigint(20) DEFAULT NULL," +
	"`permission_id` bigint(20) DEFAULT NULL," +
	"`created_on` datetime DEFAULT NULL," +
	"UNIQUE KEY `odin_role_permission_pk` (`ctx`,`role_id`,`permission_id`)" +
	") ENGINE=InnoDB DEFAULT CHARSET=utf8;"

const odinGrantSQL = "" +
	"CREATE TABLE IF NOT EXISTS `%s` (" +
	"`ctx` bigint(20) DEFAULT NULL," +
	"`role_id` bigint(20) DEFAULT NULL," +
	"`target_id` varchar(64) DEFAULT NULL," +
	"`created_on` datetime DEFAULT NULL," +
	"UNIQUE KEY `odin_grant_pk` (`ctx`,`role_id`,`target_id`)," +
	"KEY `odin_grant_role_id_index` (`role_id`)," +
	"KEY `odin_grant_target_id_index` (`target_id`)" +
	") ENGINE=InnoDB DEFAULT CHARSET=utf8;"
