package odin

import (
	"strings"
	"time"
)

const (
	k_ODIN_GRANT_PREFIX                = "odin_grant_"
	k_ODIN_GRANT_LIST                  = "odin_grant_list"
)

func getGrantKey(id string) string {
	return k_ODIN_GRANT_PREFIX + id
}

func RemoveAllGrant() (error){
	var s = getRedisSession()
	defer s.Close()

	gList, err := s.ZREVRANGE(k_ODIN_GRANT_LIST, 0, -1).Strings()
	if err != nil {
		return err
	}

	if r := s.Send("MULTI"); r.Error != nil {
		return r.Error
	}
	for _, g := range gList {
		s.Send("DEL", getGrantKey(g))
	}
	s.Send("DEL", k_ODIN_GRANT_LIST)

	if r := s.Do("EXEC"); r.Error != nil {
		return r.Error
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
func GetGrantList() (results []*GrantInfo, err error) {
	var s = getRedisSession()
	defer s.Close()

	gList, err := s.ZREVRANGE(k_ODIN_GRANT_LIST, 0, -1).Strings()
	if err != nil {
		return nil, err
	}

	for _, gId := range gList {
		var gInfo = &GrantInfo{}
		gInfo.DestinationId = gId

		gInfo.RoleIdList = s.SMEMBERS(getGrantKey(gId)).MustStrings()

		results = append(results, gInfo)
	}

	return results, err
}

////////////////////////////////////////////////////////////////////////////////
func Grant(destinationId string, roleIds ...string) (err error) {
	var s = getRedisSession()
	defer s.Close()

	if r := s.Send("MULTI"); r.Error != nil {
		return r.Error
	}

	s.Send("ZADD", k_ODIN_GRANT_LIST, time.Now().Unix(), destinationId)

	if len(roleIds) > 0 {
		var key = getGrantKey(destinationId)
		var params []interface{}
		params = append(params, key)
		for _, id := range roleIds {
			params = append(params, id)
		}
		s.Send("SADD", params...)
	}

	if r := s.Do("EXEC"); r.Error != nil {
		return r.Error
	}
	return err
}

func CancelGrant(destinationId string, roleIds ...string) (err error) {
	var s = getRedisSession()
	defer s.Close()

	if r := s.Send("MULTI"); r.Error != nil {
		return r.Error
	}

	var key = getGrantKey(destinationId)
	var params []interface{}
	params = append(params, key)
	for _, id := range roleIds {
		params = append(params, id)
	}

	s.Send("SREM", params...)

	if r := s.Do("EXEC"); r.Error != nil {
		return r.Error
	}
	return err
}

func GetGrantPermissionList(destinationId string) (results []string, err error) {
	var s = getRedisSession()
	defer s.Close()

	var roleIds = s.SMEMBERS(getGrantKey(destinationId)).MustStrings()
	for _, roleId := range roleIds {
		var pIdList = s.SMEMBERS(getRolePermissionListKey(roleId)).MustStrings()
		results = append(results, pIdList...)
	}

	return results, err
}

func GetGrantRoleList(destinationId string) (results []string, err error) {
	var s = getRedisSession()
	defer s.Close()

	results = s.SMEMBERS(getGrantKey(destinationId)).MustStrings()
	return results, err
}

func Check(destinationId string, identifiers ...string) (bool) {
	var s = getRedisSession()
	defer s.Close()

	var id = md5String(strings.Join(identifiers, "-"))

	var roleIds = s.SMEMBERS(getGrantKey(destinationId)).MustStrings()
	for _, roleId := range roleIds {
		if s.SISMEMBER(getRolePermissionListKey(roleId), id).MustBool() {
			return true
		}
	}

	if s.SISMEMBER(k_ODIN_PERMISSION_LIST, id).MustBool() == false {
		return true
	}

	return false
}