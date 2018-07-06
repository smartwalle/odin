
### 创建表

首先创建如下结构的表

```sql
CREATE TABLE `odin_grant` (
  `ctx` bigint(20) DEFAULT 0,
  `object_id` varchar(64) DEFAULT NULL,
  `role_id` int(11) DEFAULT NULL,
  `created_on` datetime DEFAULT NULL,
  UNIQUE KEY `odin_grant_destination_id_role_id_pk` (`ctx`,`object_id`,`role_id`),
  KEY `odin_grant_ctx_index` (`ctx`),
  KEY `odin_grant_object_id_index` (`object_id`),
  KEY `odin_grant_role_id_index` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8
```

```sql
CREATE TABLE `odin_group` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `ctx` bigint(20) DEFAULT 0,
  `type` int(11) DEFAULT 0,
  `name` varchar(128) DEFAULT NULL,
  `status` int(11) DEFAULT 0,
  `created_on` datetime DEFAULT NULL,
  `updated_on` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `odin_role_group_id_uindex` (`id`),
  KEY `odin_group_type_index` (`type`),
  KEY `odin_group_ctx_index` (`ctx`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8
```

```sql
CREATE TABLE `odin_permission` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `ctx` bigint(20) DEFAULT 0,
  `group_id` int(11) DEFAULT NULL,
  `name` varchar(128) DEFAULT NULL,
  `identifier` varchar(128) DEFAULT NULL,
  `status` int(11) DEFAULT 0,
  `created_on` datetime DEFAULT NULL,
  `updated_on` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `odin_permission_id_uindex` (`id`),
  KEY `odin_permission_ctx_index` (`ctx`),
  KEY `odin_permission_group_id_index` (`group_id`),
  KEY `odin_permission_status_index` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8

```

```sql
CREATE TABLE `odin_role` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `ctx` bigint(20) DEFAULT 0,
  `group_id` int(11) DEFAULT NULL,
  `name` varchar(128) DEFAULT NULL,
  `status` int(11) DEFAULT 0,
  `created_on` datetime DEFAULT NULL,
  `updated_on` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `odin_role_id_uindex` (`id`),
  KEY `odin_role_ctx_index` (`ctx`),
  KEY `odin_role_group_id_index` (`group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8
```

```sql
CREATE TABLE `odin_role_permission` (
  `ctx` bigint(20) DEFAULT 0,
  `role_id` int(11) NOT NULL,
  `permission_id` int(11) NOT NULL,
  `created_on` datetime DEFAULT NULL,
  UNIQUE KEY `odin_role_permission_role_id_permission_id_pk` (`ctx`,`role_id`,`permission_id`),
  KEY `odin_role_permission_role_id_index` (`role_id`),
  KEY `odin_role_permission_ctx_index` (`ctx`),
  KEY `odin_role_permission_permission_id_index` (`permission_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8
```