package odin

import (
	"fmt"
	"github.com/smartwalle/dbr"
)

type redisManager struct {
	r       *dbr.Pool
	tPrefix string
}

func (this *redisManager) buildKey(ctx int64, objectId string) (result string) {
	return fmt.Sprintf("%s_odin_g_%d_%s", this.tPrefix, ctx, objectId)
}

func (this *redisManager) exists(ctx int64, objectId string) (result bool) {
	var s = this.r.GetSession()
	defer s.Close()
	return s.EXISTS(this.buildKey(ctx, objectId)).MustBool()
}

func (this *redisManager) check(ctx int64, objectId, identifier string) (result bool) {
	var s = this.r.GetSession()
	defer s.Close()
	return s.SISMEMBER(this.buildKey(ctx, objectId), identifier).MustBool()
}

func (this *redisManager) checkList(ctx int64, objectId string, identifiers ...string) (result map[string]bool) {
	var s = this.r.GetSession()
	defer s.Close()

	result = make(map[string]bool)
	key := this.buildKey(ctx, objectId)
	for _, identifier := range identifiers {
		result[identifier] = s.SISMEMBER(key, identifier).MustBool()
	}
	return result
}

func (this *redisManager) grantPermissions(ctx int64, objectId string, identifier []interface{}) {
	var s = this.r.GetSession()
	defer s.Close()

	var key = this.buildKey(ctx, objectId)

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

func (this *redisManager) clear(ctx int64, objectId string) {
	var s = this.r.GetSession()
	defer s.Close()

	var cKey = "*"
	var oKey = "*"
	if ctx > 0 {
		cKey = fmt.Sprintf("%d", ctx)
	}
	if objectId != "" {
		oKey = objectId
	}
	var key = fmt.Sprintf("%s_odin_g_%s_%s", this.tPrefix, cKey, oKey)

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