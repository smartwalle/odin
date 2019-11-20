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

func (this *odinRepository) buildGrantKey(ctx int64, target string, wild bool) (result string) {
	var cKey = "0"
	var tKey = ""

	if ctx > 0 {
		cKey = fmt.Sprintf("%d", ctx)
	} else if wild {
		cKey = "*"
	}

	if target != "" {
		tKey = target
	} else if wild {
		tKey = "*"
	}
	return fmt.Sprintf("%s:odin:grant:%s:%s", this.tPrefix, cKey, tKey)
}

func (this *odinRepository) Check(ctx int64, target, identifier string) (result bool) {
	var s = this.rPool.GetSession()
	defer s.Close()

	var key = this.buildGrantKey(ctx, target, false)
	result = s.SISMEMBER(key, identifier).MustBool()
	if result == false {
		if s.EXISTS(key).MustBool() == false {
			pList, _ := this.Repository.GetGrantedPermissionList(ctx, target)
			var identifierList []interface{}
			for _, p := range pList {
				if p.Identifier == identifier {
					result = true
				}
				identifierList = append(identifierList, p.Identifier)
			}
			this.grantPermissions(key, identifierList)
		}
	}
	return result
}

func (this *odinRepository) CheckList(ctx int64, target string, identifiers ...string) (result map[string]bool) {
	var s = this.rPool.GetSession()
	defer s.Close()

	var key = this.buildGrantKey(ctx, target, false)

	if s.EXISTS(key).MustBool() == false {
		pList, _ := this.Repository.GetGrantedPermissionList(ctx, target)
		var identifierList []interface{}
		for _, p := range pList {
			identifierList = append(identifierList, p.Identifier)
		}
		this.grantPermissions(key, identifierList)
	}

	result = make(map[string]bool)
	for _, identifier := range identifiers {
		result[identifier] = s.SISMEMBER(key, identifier).MustBool()
	}

	return result
}

func (this *odinRepository) grantPermissions(key string, identifier []interface{}) {
	var s = this.rPool.GetSession()
	defer s.Close()

	if len(identifier) == 0 {
		s.DEL(key)
		return
	}
	if r := s.Send("MULTI"); r.Error != nil {
		return
	}

	var ps = []interface{}{key}
	ps = append(ps, identifier...)

	s.Send("DEL", key)
	s.Send("SADD", ps...)
	s.Send("EXPIRE", key, 3600)
	if r := s.Do("EXEC"); r.Error != nil {
		return
	}
}

func (this *odinRepository) CleanCache(ctx int64, target string) {
	var s = this.rPool.GetSession()
	defer s.Close()

	var key = this.buildGrantKey(ctx, target, true)

	var keys = s.KEYS(key).MustStrings()

	if r := s.Send("MULTI"); r.Error != nil {
		return
	}
	for _, key := range keys {
		s.Send("DEL", key)
	}
	if r := s.Do("EXEC"); r.Error != nil {
		return
	}
}
