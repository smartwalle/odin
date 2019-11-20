package mysql

const odinGrantSQL = "" +
	"CREATE TABLE IF NOT EXISTS `%s` (" +
	"`ctx` bigint(20) DEFAULT 0," +
	"`target` varchar(64) DEFAULT NULL," +
	"`role_id` bigint(20) DEFAULT NULL," +
	"`created_on` datetime DEFAULT NULL," +
	"UNIQUE KEY `odin_grant_destination_id_role_id_pk` (`ctx`,`target`,`role_id`)," +
	"KEY `odin_grant_ctx_index` (`ctx`)," +
	"KEY `odin_grant_target_index` (`target`)," +
	"KEY `odin_grant_role_id_index` (`role_id`)" +
	");"

const odinGroupSQL = "" +
	"CREATE TABLE IF NOT EXISTS `%s` (" +
	"`id` bigint(20) NOT NULL AUTO_INCREMENT," +
	"`ctx` bigint(20) DEFAULT 0," +
	"`type` int(11) DEFAULT 0," +
	"`name` varchar(128) DEFAULT NULL," +
	"`status` int(11) DEFAULT 0," +
	"`created_on` datetime DEFAULT NULL," +
	"`updated_on` datetime DEFAULT NULL," +
	"PRIMARY KEY (`id`)," +
	"UNIQUE KEY `odin_role_group_id_uindex` (`id`)," +
	"KEY `odin_group_type_index` (`type`)," +
	"KEY `odin_group_ctx_index` (`ctx`)" +
	");"

const odinPermissionSQL = "" +
	"CREATE TABLE IF NOT EXISTS `%s` (" +
	"`id` bigint(20) NOT NULL AUTO_INCREMENT," +
	"`ctx` bigint(20) DEFAULT 0," +
	"`group_id` bigint(20) DEFAULT NULL," +
	"`name` varchar(128) DEFAULT NULL," +
	"`identifier` varchar(128) DEFAULT NULL," +
	"`status` int(11) DEFAULT 0," +
	"`created_on` datetime DEFAULT NULL," +
	"`updated_on` datetime DEFAULT NULL," +
	"PRIMARY KEY (`id`)," +
	"UNIQUE KEY `odin_permission_id_uindex` (`id`)," +
	"KEY `odin_permission_ctx_index` (`ctx`)," +
	"KEY `odin_permission_group_id_index` (`group_id`)," +
	"KEY `odin_permission_status_index` (`status`)" +
	");"

const odinRoleSQL = "" +
	"CREATE TABLE  IF NOT EXISTS`%s` (" +
	"`id` bigint(20) NOT NULL AUTO_INCREMENT," +
	"`ctx` bigint(20) DEFAULT 0," +
	"`group_id` int(11) DEFAULT NULL," +
	"`name` varchar(128) DEFAULT NULL," +
	"`status` int(11) DEFAULT 0," +
	"`created_on` datetime DEFAULT NULL," +
	"`updated_on` datetime DEFAULT NULL," +
	"PRIMARY KEY (`id`)," +
	"UNIQUE KEY `odin_role_id_uindex` (`id`)," +
	"KEY `odin_role_ctx_index` (`ctx`)," +
	"KEY `odin_role_group_id_index` (`group_id`)" +
	");"

const odinRolePermissionSQL = "" +
	"CREATE TABLE IF NOT EXISTS `%s` (" +
	"`ctx` bigint(20) DEFAULT 0," +
	"`role_id` bigint(20) NOT NULL," +
	"`permission_id` bigint(20) NOT NULL," +
	"`created_on` datetime DEFAULT NULL," +
	"UNIQUE KEY `odin_role_permission_role_id_permission_id_pk` (`ctx`,`role_id`,`permission_id`)," +
	"KEY `odin_role_permission_role_id_index` (`role_id`)," +
	"KEY `odin_role_permission_ctx_index` (`ctx`)," +
	"KEY `odin_role_permission_permission_id_index` (`permission_id`)" +
	");"
