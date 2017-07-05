package odin

import "strings"

////////////////////////////////////////////////////////////////////////////////
func Grant(destinationId string, roleIds ...string) (err error) {
	var s = getRedisSession()
	defer s.Close()

	if r := s.Send("MULTI"); r.Error != nil {
		return r.Error
	}

	s.Send("SADD", k_ODIN_GRANT_LIST, destinationId)

	var key = getGrantKey(destinationId)
	var params []interface{}
	params = append(params, key)
	for _, id := range roleIds {
		params = append(params, id)
	}

	s.Send("SADD", params...)

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