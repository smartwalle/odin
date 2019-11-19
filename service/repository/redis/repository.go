package redis

import (
	"fmt"
	"github.com/smartwalle/dbr"
	"github.com/smartwalle/odin/service"
)

type odinRepository struct {
	service.OdinRepository
	rPool   *dbr.Pool
	tPrefix string
}

func NewOdinRepository(rPool *dbr.Pool, tPrefix string, repo service.OdinRepository) service.OdinRepository {
	var r = &odinRepository{}
	r.rPool = rPool
	r.tPrefix = tPrefix
	r.OdinRepository = repo
	return r
}

func (this *odinRepository) buildGrantKey(ctx int64, targetId string, wild bool) (result string) {
	var cKey = ""
	var tKey = ""

	if ctx > 0 {
		cKey = fmt.Sprintf("%d", ctx)
	} else if wild {
		cKey = "*"
	}

	if targetId != "" {
		tKey = targetId
	} else if wild {
		tKey = "*"
	}
	return fmt.Sprintf("%s:odin:grant:%s:%s", this.tPrefix, cKey, tKey)
}

func (this *odinRepository) Check(ctx int64, targetId, identifier string) (result bool) {
	var s = this.rPool.GetSession()
	defer s.Close()

	var key = this.buildGrantKey(ctx, targetId, false)
	result = s.SISMEMBER(key, identifier).MustBool()
	if result == false {
		if s.EXISTS(key).MustBool() == false {
			pList, _ := this.OdinRepository.GetGrantedPermissionList(ctx, targetId)
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

func (this *odinRepository) CheckList(ctx int64, targetId string, identifiers ...string) (result map[string]bool) {
	var s = this.rPool.GetSession()
	defer s.Close()

	var key = this.buildGrantKey(ctx, targetId, false)

	if s.EXISTS(key).MustBool() == false {
		pList, _ := this.OdinRepository.GetGrantedPermissionList(ctx, targetId)
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

func (this *odinRepository) ClearCache(ctx int64, targetId string) {
	var s = this.rPool.GetSession()
	defer s.Close()

	var key = this.buildGrantKey(ctx, targetId, true)

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
