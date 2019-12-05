package redis

import (
	"fmt"
	"github.com/smartwalle/dbr"
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/odin"
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

func (this *odinRepository) buildGrantListKey(ctx int64) (result string) {
	return fmt.Sprintf("%s:odin:grant:%d:list", this.tPrefix, ctx)
}

func (this *odinRepository) buildTargetKey(ctx int64, target string) (result string) {
	return fmt.Sprintf("%s:odin:grant:%d:%s", this.tPrefix, ctx, target)
}

func (this *odinRepository) CheckPermission(ctx int64, target string, permissionName string) bool {
	var rSess = this.rPool.GetSession()
	defer rSess.Close()

	var key = this.buildTargetKey(ctx, target)
	var result = rSess.SISMEMBER(key, permissionName).MustBool()

	if result == false {
		if rSess.EXISTS(key).MustBool() == false {
			pList, err := this.Repository.GetGrantedPermissions(ctx, target)
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
			this.grantPermissions(ctx, key, pNames)
		}
	}

	return result
}

func (this *odinRepository) grantPermissions(ctx int64, key string, permissionNames []interface{}) {
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
	rSess.Send("SADD", this.buildGrantListKey(ctx), key)
	rSess.Send("EXPIRE", key, 3600)
	rSess.Do("EXEC")
}

func (this *odinRepository) CleanCache(ctx int64, target string) {
	var rSess = this.rPool.GetSession()
	defer rSess.Close()

	if target == "" || target == "*" {
		var key = this.buildGrantListKey(ctx)
		var items = rSess.SMEMBERS(key).MustStrings()
		if rSess.Send("MULTI").Error != nil {
			return
		}
		for _, item := range items {
			rSess.Send("DEL", item)
			rSess.Send("SREM", key, item)
		}
		rSess.Do("EXEC")
	} else {
		var key = this.buildTargetKey(ctx, target)
		rSess.DEL(key)
	}
}
