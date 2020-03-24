package mysql

import (
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/odin"
	"github.com/smartwalle/odin/service/repository/internal/sql"
	"strings"
)

type repository struct {
	*sql.Repository
}

func NewRepository(db dbs.DB, tablePrefix string) odin.Repository {
	var r = &repository{}
	r.Repository = sql.NewRepository(db, dbs.DialectMySQL, tablePrefix)
	return r
}

func (this *repository) BeginTx() (dbs.TX, odin.Repository) {
	var nRepo = *this
	var tx dbs.TX
	tx, nRepo.Repository = this.Repository.BeginTx()
	return tx, &nRepo
}

func (this *repository) WithTx(tx dbs.TX) odin.Repository {
	var nRepo = *this
	nRepo.Repository = this.Repository.WithTx(tx)
	return &nRepo
}

func (this *repository) InitTable() error {
	var rawText = `
create table if not exists odin_grant
(
	ctx        int         null,
	role_id    bigint      null,
	target     varchar(64) null,
	created_on datetime    null,
	constraint odin_grant_pk
		unique (ctx, role_id, target)
);

create index odin_grant_ctx_target_index
	on odin_grant (ctx, target);

create index odin_grant_role_id_index
	on odin_grant (role_id);

create index odin_grant_target_index
	on odin_grant (target);

create table if not exists odin_group
(
	id         bigint auto_increment,
	ctx        int          null,
	type       int(2)       null,
	name       varchar(64)  null,
	alias_name varchar(255) null,
	status     int(2)       null,
	created_on datetime     null,
	updated_on datetime     null,
	constraint odin_group_id_uindex
		unique (id),
	constraint odin_group_pk
		unique (ctx, type, name)
);

alter table odin_group
	add primary key (id);

create table if not exists odin_permission
(
	id          bigint auto_increment,
	group_id    bigint           null,
	ctx         int              null,
	name        varchar(255)     null,
	alias_name  varchar(255)     null,
	status      int(2) default 1 null,
	description varchar(1024)    null,
	created_on  datetime         null,
	updated_on  datetime         null,
	constraint odin_permission_ctx_name_uindex
		unique (ctx, name),
	constraint odin_permission_id_uindex
		unique (id)
);

create index odin_permission_ctx_group_id_index
	on odin_permission (ctx, group_id);

create index odin_permission_ctx_index
	on odin_permission (ctx);

alter table odin_permission
	add primary key (id);

create table if not exists odin_pre_permission
(
	ctx               int                  not null,
	permission_id     bigint               not null,
	pre_permission_id bigint               not null,
	auto_grant        tinyint(1) default 0 null,
	created_on        datetime             null,
	primary key (ctx, permission_id, pre_permission_id)
);

create index odin_pre_permission_ctx_permission_id_index
	on odin_pre_permission (ctx, permission_id);

create index odin_pre_permission_ctx_permission_id_pre_permission_id_index
	on odin_pre_permission (ctx, permission_id, pre_permission_id);

create table if not exists odin_pre_role
(
	ctx         int      not null,
	role_id     bigint   not null,
	pre_role_id bigint   not null,
	created_on  datetime null,
	primary key (ctx, role_id, pre_role_id)
);

create index odin_pre_role_ctx_role_id_index
	on odin_pre_role (ctx, role_id);

create index odin_pre_role_ctx_role_id_pre_role_id_index
	on odin_pre_role (ctx, role_id, pre_role_id);

create table if not exists odin_role
(
	id          bigint auto_increment,
	ctx         int              null,
	name        varchar(64)      null,
	alias_name  varchar(255)     null,
	status      int(2) default 1 null,
	description varchar(1024)    null,
	parent_id   bigint           null,
	left_value  int              null,
	right_value int              null,
	depth       int              null,
	created_on  datetime         null,
	updated_on  datetime         null,
	constraint odin_role_ctx_name_uindex
		unique (ctx, name),
	constraint odin_role_id_uindex
		unique (id)
);

create index odin_role_ctx_index
	on odin_role (ctx);

create index odin_role_ctx_left_value_index
	on odin_role (ctx, left_value);

create index odin_role_ctx_parent_id_index
	on odin_role (ctx, parent_id);

create index odin_role_ctx_right_value_index
	on odin_role (ctx, right_value);

alter table odin_role
	add primary key (id);

create table if not exists odin_role_mutex
(
	ctx           int      not null,
	role_id       bigint   not null,
	mutex_role_id bigint   not null,
	created_on    datetime null,
	primary key (ctx, role_id, mutex_role_id)
);

create index odin_role_mutex_ctx_mutex_role_id_index
	on odin_role_mutex (ctx, mutex_role_id);

create index odin_role_mutex_ctx_mutex_role_id_role_id_index
	on odin_role_mutex (ctx, mutex_role_id, role_id);

create index odin_role_mutex_ctx_role_id_index
	on odin_role_mutex (ctx, role_id);

create table if not exists odin_role_permission
(
	ctx           int      null,
	role_id       bigint   null,
	permission_id bigint   null,
	created_on    datetime null,
	constraint odin_role_permission_pk
		unique (ctx, role_id, permission_id)
);
`

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
