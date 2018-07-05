package odin

import (
	"fmt"
	"github.com/smartwalle/dbr"
)

type redisManager struct {
	r       *dbr.Pool
	tPrefix string
}

func (this *redisManager) buildKey(ctxId int64, objectId string) (result string) {
	return fmt.Sprintf("%s_odin_g_%d_%s", this.tPrefix, ctxId, objectId)
}

func (this *redisManager) exists(ctxId int64, objectId string) (result bool) {
	var s = this.r.GetSession()
	defer s.Close()
	return s.EXISTS(this.buildKey(ctxId, objectId)).MustBool()
}

func (this *redisManager) check(ctxId int64, objectId, identifier string) (result bool) {
	var s = this.r.GetSession()
	defer s.Close()
	return s.SISMEMBER(this.buildKey(ctxId, objectId), identifier).MustBool()
}

func (this *redisManager) checkList(ctxId int64, objectId string, identifiers ...string) (result map[string]bool) {
	var s = this.r.GetSession()
	defer s.Close()

	result = make(map[string]bool)
	key := this.buildKey(ctxId, objectId)
	for _, identifier := range identifiers {
		result[identifier] = s.SISMEMBER(key, identifier).MustBool()
	}
	return result
}

func (this *redisManager) grantPermissions(ctxId int64, objectId string, identifier []interface{}) {
	var s = this.r.GetSession()
	defer s.Close()

	var key = this.buildKey(ctxId, objectId)

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

func (this *redisManager) clear() {
	var s = this.r.GetSession()
	defer s.Close()

	var keys = s.KEYS(this.tPrefix + "_odin_g_*").MustStrings()

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
