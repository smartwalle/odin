package postgresql

import (
	"github.com/smartwalle/dbs"
	"time"
)

func (this *repository) AddRoleMutex(ctx, roleId int64, mutexRoleIds []int64) (err error) {
	var now = time.Now()
	var ib = dbs.NewInsertBuilder()
	ib.UseDialect(this.Dialect())
	ib.Table(this.TableRoleMutex())
	ib.Columns("ctx", "role_id", "mutex_role_id", "created_on")
	for _, mutexRoleId := range mutexRoleIds {
		ib.Values(ctx, roleId, mutexRoleId, now)
		ib.Values(ctx, mutexRoleId, roleId, now)
	}
	if _, err = ib.Exec(this.DB()); err != nil {
		return err
	}
	return nil
}
