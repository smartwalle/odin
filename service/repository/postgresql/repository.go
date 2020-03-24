package postgresql

import (
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/odin"
	"github.com/smartwalle/odin/service/repository/internal/sql"
	"strings"
	"time"
)

type repository struct {
	*sql.Repository
}

func NewRepository(db dbs.DB, tablePrefix string) odin.Repository {
	var r = &repository{}
	r.Repository = sql.NewRepository(db, dbs.DialectPostgreSQL, tablePrefix)
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
create table if not exists odin_group
(
	id         bigint not null,
	ctx        integer,
	type       integer,
	name       varchar(64),
	alias_name varchar(255),
	status     integer,
	created_on timestamp with time zone,
	updated_on timestamp with time zone,
	constraint odin_group_pk
		primary key (id),
	constraint odin_group_pk_2
		unique (ctx, type, name)
);

create table if not exists odin_permission
(
	id          bigint not null,
	group_id    bigint,
	ctx         integer,
	name        varchar(255),
	alias_name  varchar(255),
	status      integer default 1,
	description varchar(1024),
	created_on  timestamp with time zone,
	updated_on  timestamp with time zone,
	constraint odin_permission_pk
		primary key (id),
	constraint odin_permission_pk_2
		unique (ctx, name)
);

create unique index if not exists odin_permission_id_uindex
	on odin_permission (id);

create index if not exists odin_permission_ctx_group_id_index
	on odin_permission (ctx, group_id);

create index if not exists odin_permission_ctx_index
	on odin_permission (ctx);

create unique index if not exists odin_permission_ctx_name_uindex
	on odin_permission (ctx, name);

create table if not exists odin_role
(
	id          bigint not null,
	ctx         integer,
	name        varchar(64),
	alias_name  varchar(255),
	status      integer default 1,
	description varchar(1024),
	parent_id   bigint,
	left_value  integer,
	right_value integer,
	depth       integer,
	created_on  timestamp with time zone,
	updated_on  timestamp with time zone,
	constraint odin_role_pk
		primary key (id),
	constraint odin_role_pk_2
		unique (ctx, name)
);

create unique index if not exists odin_role_id_uindex
	on odin_role (id);

create unique index if not exists odin_role_ctx_name_uindex
	on odin_role (ctx, name);

create index if not exists odin_role_ctx_index
	on odin_role (ctx);

create index if not exists odin_role_ctx_left_value_index
	on odin_role (ctx, left_value);

create index if not exists odin_role_ctx_parent_id_index
	on odin_role (ctx, parent_id);

create index if not exists odin_role_ctx_right_value_index
	on odin_role (ctx, right_value);

create table if not exists odin_grant
(
	ctx        integer     not null,
	role_id    bigint      not null,
	target     varchar(64) not null,
	created_on timestamp with time zone,
	constraint odin_grant_pk
		primary key (ctx, role_id, target)
);

create index if not exists odin_grant_role_id_index
	on odin_grant (role_id);

create index if not exists odin_grant_target_index
	on odin_grant (target);

create index if not exists odin_grant_ctx_target_index
	on odin_grant (ctx, target);

create table if not exists odin_role_permission
(
	ctx           integer not null,
	role_id       bigint  not null,
	permission_id bigint  not null,
	created_on    timestamp with time zone,
	constraint odin_role_permission_pk
		primary key (ctx, role_id, permission_id)
);

create table if not exists odin_role_mutex
(
	ctx           integer not null,
	role_id       bigint  not null,
	mutex_role_id bigint  not null,
	created_on    timestamp with time zone,
	constraint odin_role_mutex_pk
		primary key (ctx, role_id, mutex_role_id)
);

create index if not exists odin_role_mutex_ctx_mutex_role_id_index
	on odin_role_mutex (ctx, mutex_role_id);

create index if not exists odin_role_mutex_ctx_role_id_index
	on odin_role_mutex (ctx, role_id);

create table if not exists odin_pre_role
(
	ctx         integer not null,
	role_id     bigint  not null,
	pre_role_id bigint  not null,
	created_on  timestamp with time zone,
	constraint odin_pre_role_pk
		primary key (ctx, role_id, pre_role_id)
);

create index if not exists odin_pre_role_ctx_role_id_index
	on odin_pre_role (ctx, role_id);

create table if not exists odin_pre_permission
(
	ctx               integer not null,
	permission_id     bigint  not null,
	pre_permission_id bigint  not null,
	auto_grant        boolean default false,
	created_on        timestamp with time zone,
	constraint odin_pre_permission_pk
		primary key (ctx, permission_id, pre_permission_id)
);

create index if not exists odin_pre_permission_ctx_permission_id_index
	on odin_pre_permission (ctx, permission_id);
`

	// 添加 rule
	rawText = rawText + `
create or replace rule odin_role_permission_pk_rule as on insert to odin_role_permission where exists (
select 1 from odin_role_permission where ctx = NEW.ctx and role_id = NEW.role_id and permission_id = NEW.permission_id
) do instead nothing;
`
	rawText = rawText + `
create or replace rule odin_pre_permission_pk_rule as on insert to odin_pre_permission where exists (
select 1 from odin_pre_permission where ctx = NEW.ctx and permission_id = NEW.permission_id and pre_permission_id = NEW.pre_permission_id
) do instead nothing;
`
	rawText = rawText + `
create or replace rule odin_role_mutex_pk_rule as on insert to odin_role_mutex where exists (
select 1 from odin_role_mutex where ctx = NEW.ctx and role_id = NEW.role_id and mutex_role_id = NEW.mutex_role_id
) do instead nothing;
`
	rawText = rawText + `
create or replace rule odin_pre_role_pk_rule as on insert to odin_pre_role where exists (
select 1 from odin_pre_role where ctx = NEW.ctx and role_id = NEW.role_id and pre_role_id = NEW.pre_role_id
) do instead nothing;
`

	rawText = rawText + `
create or replace rule odin_grant_pk_rule as on insert to odin_grant where exists (
select 1 from odin_grant where ctx = NEW.ctx and role_id = NEW.role_id and target = NEW.target
) do instead nothing;
`
	var sql = strings.ReplaceAll(rawText, "odin", this.TablePrefix())
	var rb = dbs.NewBuilder(sql)
	if _, err := rb.Exec(this.DB()); err != nil {
		return err
	}

	return nil
}

func (this *repository) GrantPermissionWithIds(ctx, roleId int64, permissionIds []int64) (err error) {
	if len(permissionIds) == 0 {
		return nil
	}
	var now = time.Now()
	var ib = dbs.NewInsertBuilder()
	ib.UseDialect(this.Dialect())
	ib.Table(this.TableRolePermission())
	ib.Columns("ctx", "role_id", "permission_id", "created_on")
	for _, permissionId := range permissionIds {
		ib.Values(ctx, roleId, permissionId, now)
	}
	if _, err = ib.Exec(this.DB()); err != nil {
		return err
	}
	return nil
}

func (this *repository) AddPrePermission(ctx, permissionId int64, prePermissionIds []int64) (err error) {
	var now = time.Now()
	var ib = dbs.NewInsertBuilder()
	ib.UseDialect(this.Dialect())
	ib.Table(this.TablePrePermission())
	ib.Columns("ctx", "permission_id", "pre_permission_id", "auto_grant", "created_on")
	for _, prePermissionId := range prePermissionIds {
		ib.Values(ctx, permissionId, prePermissionId, false, now)
	}
	if _, err = ib.Exec(this.DB()); err != nil {
		return err
	}
	return nil
}
