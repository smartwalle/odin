package odin

import (
	"github.com/smartwalle/dbr"
)

type redisManager struct {
	r       *dbr.Pool
	tPrefix string
}

func (this *redisManager) buildKey(objectId string) (result string) {
	return this.tPrefix + "_odin_g_" + objectId
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

func (this *redisManager) checkList(objectId string, identifiers ...string) (result map[string]bool) {
	var s = this.r.GetSession()
	defer s.Close()

	result = make(map[string]bool)
	key := this.buildKey(objectId)
	for _, identifier := range identifiers {
		result[identifier] = s.SISMEMBER(key, identifier).MustBool()
	}
	return result
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
