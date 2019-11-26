package mysql

import (
	"fmt"
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

	if err := r.initTable(); err != nil {
		panic(fmt.Sprintf("初始化 Odin 失败, 错误信息为: %v", err))
	}

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

func (this *odinRepository) initTable() error {
	var tx = dbs.MustTx(this.db)

	var cb *dbs.RawBuilder

	cb = dbs.NewBuilder("")
	cb.Format(odinGroupSQL, this.tblGroup)
	if _, err := cb.Exec(tx); err != nil {
		return err
	}

	cb = dbs.NewBuilder("")
	cb.Format(odinPermissionSQL, this.tblPermission)
	if _, err := cb.Exec(tx); err != nil {
		return err
	}

	cb = dbs.NewBuilder("")
	cb.Format(odinRoleSQL, this.tblRole)
	if _, err := cb.Exec(tx); err != nil {
		return err
	}

	cb = dbs.NewBuilder("")
	cb.Format(odinRolePermissionSQL, this.tblRolePermission)
	if _, err := cb.Exec(tx); err != nil {
		return err
	}

	cb = dbs.NewBuilder("")
	cb.Format(odinGrantSQL, this.tblGrant)
	if _, err := cb.Exec(tx); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func (this *odinRepository) CleanCache(ctx int64, targetId string) {
}
