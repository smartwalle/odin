package mysql

const odinGroupSQL = "" +
	"CREATE TABLE `%s_group` (" +
	"`id` bigint(20) NOT NULL AUTO_INCREMENT," +
	"`ctx` bigint(20) DEFAULT NULL," +
	"`type` int(2) DEFAULT NULL," +
	"`name` varchar(64) DEFAULT NULL," +
	"`alias_name` varchar(255) DEFAULT NULL," +
	"`status` int(2) DEFAULT NULL," +
	"`created_on` datetime DEFAULT NULL," +
	"`updated_on` datetime DEFAULT NULL," +
	"PRIMARY KEY (`id`)," +
	"UNIQUE KEY `%s_group_id_uindex` (`id`)," +
	"UNIQUE KEY `%s_group_pk` (`ctx`,`type`,`name`)" +
	") ENGINE=InnoDB DEFAULT CHARSET=utf8;"

const odinPermissionSQL = "" +
	"CREATE TABLE IF NOT EXISTS `%s` (" +
	"`id` bigint(20) NOT NULL AUTO_INCREMENT," +
	"`ctx` bigint(20) DEFAULT NULL," +
	"`name` varchar(255) DEFAULT NULL," +
	"`alias_name` varchar(255) DEFAULT NULL," +
	"`status` int(2) DEFAULT '1'," +
	"`description` varchar(1024)," +
	"`created_on` datetime DEFAULT NULL," +
	"`updated_on` datetime DEFAULT NULL," +
	"PRIMARY KEY (`id`)," +
	"UNIQUE KEY `%s_id_uindex` (`id`)," +
	"UNIQUE KEY `%s_ctx_name_uindex` (`ctx`,`name`)," +
	"KEY `%s_ctx_index` (`ctx`)" +
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
	"UNIQUE KEY `%s_id_uindex` (`id`)," +
	"UNIQUE KEY `%s_ctx_name_uindex` (`ctx`,`name`)," +
	"KEY `%s_ctx_index` (`ctx`)" +
	") ENGINE=InnoDB DEFAULT CHARSET=utf8;"

const odinRolePermissionSQL = "" +
	"CREATE TABLE IF NOT EXISTS `%s` (" +
	"`ctx` bigint(20) DEFAULT NULL," +
	"`role_id` bigint(20) DEFAULT NULL," +
	"`permission_id` bigint(20) DEFAULT NULL," +
	"`created_on` datetime DEFAULT NULL," +
	"UNIQUE KEY `%s_pk` (`ctx`,`role_id`,`permission_id`)" +
	") ENGINE=InnoDB DEFAULT CHARSET=utf8;"

const odinGrantSQL = "" +
	"CREATE TABLE IF NOT EXISTS `%s` (" +
	"`ctx` bigint(20) DEFAULT NULL," +
	"`role_id` bigint(20) DEFAULT NULL," +
	"`target_id` varchar(64) DEFAULT NULL," +
	"`created_on` datetime DEFAULT NULL," +
	"UNIQUE KEY `%s_pk` (`ctx`,`role_id`,`target_id`)," +
	"KEY `%s_role_id_index` (`role_id`)," +
	"KEY `%s_target_id_index` (`target_id`)" +
	") ENGINE=InnoDB DEFAULT CHARSET=utf8;"
