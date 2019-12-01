package redis

import (
	"github.com/smartwalle/dbr"
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/odin"
	"fmt"
)

type odinRepository struct {
	odin.Repository
	rPool   dbr.Pool
	tPrefix string
}

func NewRepository(rPool dbr.Pool, tPrefix string, repo odin.Repository) odin.Repository {
	var r = &odinRepository{}
	r.rPool = rPool
	r.tPrefix = tPrefix
	r.Repository = repo
	return r
}

func (this *odinRepository) BeginTx() (dbs.TX, odin.Repository) {
	var nRepo = *this
	var tx dbs.TX
	tx, nRepo.Repository = this.Repository.BeginTx()
	return tx, &nRepo
}

func (this *odinRepository) WithTx(tx dbs.TX) odin.Repository {
	var nRepo = *this
	nRepo.Repository = this.Repository.WithTx(tx)
	return &nRepo
}

func (this *odinRepository) buildGrantKey(ctx int64, target string) (result string) {
	return fmt.Sprintf("%s:odin:grant:%d:%s", this.tPrefix, ctx, target)
}

func (this *odinRepository) Check(ctx int64, targetId string, permissionName string) (bool) {
	var rSess = this.rPool.GetSession()
	defer rSess.Close()

	var key = this.buildGrantKey(ctx, targetId)
	var result = rSess.SISMEMBER(key, permissionName).MustBool()

	if result == false {
		if rSess.EXISTS(key).MustBool() == false {
			pList, err := this.Repository.GetGrantedPermissions(ctx, targetId)
			if err != nil {
				return false
			}

			var pNames = make([]interface{}, 0, len(pList))
			for _, p := range pList {
				if p.Name == permissionName {
					result = true
				}
				pNames = append(pNames, p.Name)
			}
			this.grantPermissions(key, pNames)
		}
	}

	return result
}

func (this *odinRepository) grantPermissions(key string, permissionNames []interface{}) {
	var rSess = this.rPool.GetSession()
	defer rSess.Close()

	if len(permissionNames) == 0 {
		rSess.DEL(key)
		return
	}

	if rSess.Send("MULTI").Error != nil {
		return
	}

	var ps = make([]interface{}, 0, 1+len(permissionNames))
	ps = append(ps, key)
	ps = append(ps, permissionNames...)

	rSess.Send("DEL", key)
	rSess.Send("SADD", ps...)
	rSess.Send("EXPIRE", key, 3600)
	rSess.Do("EXEC")
}

func (this *odinRepository) CleanCache(ctx int64, targetId string) {
	var rSess = this.rPool.GetSession()
	defer rSess.Close()

	var key = this.buildGrantKey(ctx, targetId)
	rSess.DEL(key)
}
