package postgresql

import (
	"github.com/smartwalle/dbs"
	"time"
)

func (this *repository) AddPreRole(ctx, roleId int64, preRoleIds []int64) (err error) {
	var now = time.Now()
	var ib = dbs.NewInsertBuilder()
	ib.UseDialect(this.Dialect())
	ib.Table(this.TablePreRole())
	ib.Columns("ctx", "role_id", "pre_role_id", "created_on")
	for _, preRoleId := range preRoleIds {
		ib.Values(ctx, roleId, preRoleId, now)
	}
	if _, err = ib.Exec(this.DB()); err != nil {
		return err
	}
	return nil
}
