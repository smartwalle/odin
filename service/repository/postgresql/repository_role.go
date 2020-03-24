package postgresql

import (
	"github.com/smartwalle/dbs"
	"time"
)

func (this *repository) GrantRoleWithIds(ctx int64, target string, roleIds ...int64) (err error) {
	if len(roleIds) == 0 {
		return nil
	}
	var now = time.Now()
	var ib = dbs.NewInsertBuilder()
	ib.UseDialect(this.Dialect())
	ib.Table(this.TableGrant())
	ib.Columns("ctx", "role_id", "target", "created_on")
	for _, rId := range roleIds {
		ib.Values(ctx, rId, target, now)
	}
	if _, err = ib.Exec(this.DB()); err != nil {
		return err
	}
	return nil
}
