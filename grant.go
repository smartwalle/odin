package odin

import (
	"time"
)

const (
	k_ODIN_GRANT_PREFIX                = "odin_grant_"
	k_ODIN_GRANT_LIST                  = "odin_grant_list"
)

func getGrantKey(id string) string {
	return k_ODIN_GRANT_PREFIX + id
}

// RemoveAllGrant 清除所有的授权信息.
func RemoveAllGrant() (error){
	var s = getRedisSession()
	defer s.Close()

	gIdList, err := s.ZREVRANGE(k_ODIN_GRANT_LIST, 0, -1).Strings()
	if err != nil {
		return err
	}

	if r := s.Send("MULTI"); r.Error != nil {
		return r.Error
	}
	for _, gId := range gIdList {
		s.Send("DEL", getGrantKey(gId))
	}
	s.Send("DEL", k_ODIN_GRANT_LIST)

	if r := s.Do("EXEC"); r.Error != nil {
		return r.Error
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// GetGrantList 获取所有的授权信息列表.
func GetGrantList() (results []*GrantInfo, err error) {
	var s = getRedisSession()
	defer s.Close()

	gIdList, err := s.ZREVRANGE(k_ODIN_GRANT_LIST, 0, -1).Strings()
	if err != nil {
		return nil, err
	}

	for _, gId := range gIdList {
		var gInfo = &GrantInfo{}
		gInfo.DestinationId = gId
		gInfo.RoleIdList = s.SMEMBERS(getGrantKey(gId)).MustStrings()
		results = append(results, gInfo)
	}

	return results, err
}

////////////////////////////////////////////////////////////////////////////////
// Grant 向 destinationId 授予角色信息.
func Grant(destinationId string, roleIds ...string) (err error) {
	var s = getRedisSession()
	defer s.Close()

	if r := s.Send("MULTI"); r.Error != nil {
		return r.Error
	}

	s.Send("ZADD", k_ODIN_GRANT_LIST, time.Now().Unix(), destinationId)

	var key = getGrantKey(destinationId)

	if len(roleIds) > 0 {
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

// CancelGrant 取消对 destinationId 的指定角色授权.
func CancelGrant(destinationId string, roleIds ...string) (err error) {
	var s = getRedisSession()
	defer s.Close()

	if r := s.Send("MULTI"); r.Error != nil {
		return r.Error
	}

	var key = getGrantKey(destinationId)
	var params []interface{}
	params = append(params, key)
	for _, rId := range roleIds {
		params = append(params, rId)
	}

	s.Send("SREM", params...)

	if r := s.Do("EXEC"); r.Error != nil {
		return r.Error
	}
	return err
}

// CancelAllGrant 取消对 destinationId 所有角色授权.
func CancelAllGrant(destinationId string) (err error) {
	var s = getRedisSession()
	defer s.Close()

	var key = getGrantKey(destinationId)
	s.Send("DEL", key)

	s.Send("ZREM", k_ODIN_GRANT_LIST, destinationId)

	if r := s.Do("EXEC"); r.Error != nil {
		return r.Error
	}
	return err
}

// GetGrantPermissionList 获取 destinationId 拥有的所有权限信息.
func GetGrantPermissionList(destinationId string) (results []string, err error) {
	var s = getRedisSession()
	defer s.Close()

	var rIdList = s.SMEMBERS(getGrantKey(destinationId)).MustStrings()
	for _, rId := range rIdList {
		var pIdList = s.SMEMBERS(getRolePermissionListKey(rId)).MustStrings()
		results = append(results, pIdList...)
	}

	return results, err
}

// GetGrantRoleList 获取 destinationId 拥有的所有角色信息.
func GetGrantRoleList(destinationId string) (results []string, err error) {
	var s = getRedisSession()
	defer s.Close()

	results = s.SMEMBERS(getGrantKey(destinationId)).MustStrings()
	return results, err
}

// Check 验证 destinationId 是否有访问指定信息的权限.
func Check(destinationId, identifier string) (bool) {
	var s = getRedisSession()
	defer s.Close()

	var pId = md5String(identifier)

	var rIdList = s.SMEMBERS(getGrantKey(destinationId)).MustStrings()
	for _, rId := range rIdList {
		if s.SISMEMBER(getRolePermissionListKey(rId), pId).MustBool() {
			return true
		}
	}

	// 如果验证一项不存在的权限信息，那么将返回 true.
	//if s.ZSCORE(k_ODIN_PERMISSION_LIST, pId).Data == nil {
	//	return true
	//}

	return false
}