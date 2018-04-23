package odin

import (
	"github.com/smartwalle/dbr"
)

type redisManager struct {
	r *dbr.Pool
}

func (this *redisManager) buildKey(objectId string) (result string) {
	return "odin_grant_" + objectId
}

func (this *redisManager) exists(objectId string) (result bool) {
	var s = this.r.GetSession()
	defer s.Close()
	return s.EXISTS(this.buildKey(objectId)).MustBool()
}

func (this *redisManager) check(objectId, identifier string) (result bool) {
	var s = this.r.GetSession()
	defer s.Close()
	return s.SISMEMBER(this.buildKey(objectId), identifier).MustBool()
}

func (this *redisManager) grantPermissions(objectId string, identifier []interface{}) {
	var s = this.r.GetSession()
	defer s.Close()

	var key = this.buildKey(objectId)

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

	var keys = s.KEYS("odin_*").MustStrings()

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
