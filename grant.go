package odin

import (
	"time"
)

const (
	k_ODIN_GRANT_PREFIX         = "odin_grant_"
	k_ODIN_GRANT_LIST           = "odin_grant_list"
	k_ODIN_GRANT_PRIVATE_PREFIX = "odin_grant_private_"
	k_ODIN_GRANT_PRIVATE_LIST   = "odin_grant_private_list"
)

func getGrantKey(id string) string {
	return k_ODIN_GRANT_PREFIX + id
}

func getGrantPrivateKey(id string) string {
	return k_ODIN_GRANT_PRIVATE_PREFIX + id
}

// RevokeAllGrant 清除所有的角色授权信息.
func RevokeAllGrant() error {
	var s = getRedisSession()
	defer s.Close()

	gIdList, err := s.ZREVRANGE(k_ODIN_GRANT_LIST, 0, -1).Strings()
	if err != nil {
		return err
	}

	pIdList, err := s.ZREVRANGE(k_ODIN_GRANT_PRIVATE_LIST, 0, -1).Strings()
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

	for _, pId := range pIdList {
		s.Send("DEL",  getGrantPrivateKey(pId))
	}
	s.Send("DEL", k_ODIN_GRANT_PRIVATE_LIST)

	if r := s.Do("EXEC"); r.Error != nil {
		return r.Error
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// GetAllGrantRoleList 获取所有的角色授权信息列表.
func GetAllGrantRoleList() (results []*GrantInfo, err error) {
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
// GrantRole 为 destinationId 授权角色信息.
func GrantRole(destinationId string, roleIds ...string) (err error) {
	var s = getRedisSession()
	defer s.Close()

	var key = getGrantKey(destinationId)

	var params []interface{}
	if len(roleIds) > 0 {
		params = append(params, key)
		for _, roleId := range roleIds {
			// 判断 role 是否存在
			if s.ZSCORE(k_ODIN_ROLE_LIST, roleId).Data != nil {
				params = append(params, roleId)
			}
		}
	}

	if r := s.Send("MULTI"); r.Error != nil {
		return r.Error
	}

	s.Send("ZADD", k_ODIN_GRANT_LIST, time.Now().Unix(), destinationId)

	if len(params) > 0 {
		s.Send("SADD", params...)
	}

	if r := s.Do("EXEC"); r.Error != nil {
		return r.Error
	}
	return err
}

// GrantPermission 为 destinationId 授权某个独立的权限信息.
func GrantPermission(destinationId string, identifiers ...string) (err error) {
	var s = getRedisSession()
	defer s.Close()

	if r := s.Send("MULTI"); r.Error != nil {
		return r.Error
	}
	var key = getGrantPrivateKey(destinationId)
	if len(identifiers) > 0 {
		var params []interface{}
		params = append(params, key)
		for _, identifier := range identifiers {
			params = append(params, identifier)
		}
		s.Send("SADD", params...)

		s.Send("ZADD", k_ODIN_GRANT_PRIVATE_LIST, time.Now().Unix(), destinationId)
	}

	if r := s.Do("EXEC"); r.Error != nil {
		return r.Error
	}
	return err
}

// RevokeRole 取消对 destinationId 的指定角色授权.
func RevokeRole(destinationId string, roleIds ...string) (err error) {
	var s = getRedisSession()
	defer s.Close()

	if r := s.Send("MULTI"); r.Error != nil {
		return r.Error
	}

	if len(roleIds) > 0 {
		var key = getGrantKey(destinationId)
		var params []interface{}
		params = append(params, key)
		for _, rId := range roleIds {
			params = append(params, rId)
		}
		s.Send("SREM", params...)
	}

	if r := s.Do("EXEC"); r.Error != nil {
		return r.Error
	}
	return err
}

// RevokePermission 取消对 destinationId 的独立权限授权.
func RevokePermission(destinationId string, identifiers ...string) (err error) {
	var s = getRedisSession()
	defer s.Close()

	if r := s.Send("MULTI"); r.Error != nil {
		return r.Error
	}

	if len(identifiers) > 0 {
		var key = getGrantPrivateKey(destinationId)
		var params []interface{}
		params = append(params, key)
		for _, identifier := range identifiers {
			params = append(params, identifier)
		}
		s.Send("SREM", params...)
	}

	if r := s.Do("EXEC"); r.Error != nil {
		return r.Error
	}
	return err
}

// RevokeAllRole 取消对 destinationId 所有角色授权.
func RevokeAllRole(destinationId string) (err error) {
	var s = getRedisSession()
	defer s.Close()

	if r := s.Send("MULTI"); r.Error != nil {
		return r.Error
	}

	var key = getGrantKey(destinationId)
	s.Send("DEL", key)
	s.Send("ZREM", k_ODIN_GRANT_LIST, destinationId)

	if r := s.Do("EXEC"); r.Error != nil {
		return r.Error
	}
	return err
}

// RevokeAllPermission 取消对 destinationId 所有独立权限授权.
func RevokeAllPermission(destinationId string) (err error) {
	var s = getRedisSession()
	defer s.Close()

	if r := s.Send("MULTI"); r.Error != nil {
		return r.Error
	}

	var key = getGrantKey(destinationId)
	s.Send("DEL", key)
	s.Send("ZREM", k_ODIN_GRANT_PRIVATE_LIST, destinationId)

	if r := s.Do("EXEC"); r.Error != nil {
		return r.Error
	}
	return err
}

////////////////////////////////////////////////////////////////////////////////
// GetGrantPermissions 获取 destinationId 拥有的所有权限信息.
func GetGrantPermissions(destinationId string) (results []string, err error) {
	var s = getRedisSession()
	defer s.Close()

	var rIdList = s.SMEMBERS(getGrantKey(destinationId)).MustStrings()
	for _, rId := range rIdList {
		var pIdList = s.SMEMBERS(getRolePermissionListKey(rId)).MustStrings()
		results = append(results, pIdList...)
	}

	var pList = s.SMEMBERS(getGrantPrivateKey(destinationId)).MustStrings()
	results = append(results, pList...)

	return results, err
}

// GetGrantRoles 获取 destinationId 拥有的所有角色信息.
func GetGrantRoles(destinationId string) (results []string, err error) {
	var s = getRedisSession()
	defer s.Close()

	results = s.SMEMBERS(getGrantKey(destinationId)).MustStrings()
	return results, err
}

////////////////////////////////////////////////////////////////////////////////
// Check 验证 destinationId 是否有访问指定信息的权限.
func Check(destinationId, identifier string) bool {
	var s = getRedisSession()
	defer s.Close()

	var pId = md5String(identifier)

	var rIdList = s.SMEMBERS(getGrantKey(destinationId)).MustStrings()
	for _, rId := range rIdList {
		if s.SISMEMBER(getRolePermissionListKey(rId), pId).MustBool() {
			return true
		}
	}

	if s.SISMEMBER(getGrantPrivateKey(destinationId), identifier).MustBool() {
		return true
	}

	// 如果验证一项不存在的权限信息，那么将返回 true.
	//if s.ZSCORE(k_ODIN_PERMISSION_LIST, pId).Data == nil {
	//	return true
	//}

	return false
}

func CheckList(destinationId string, identifiers ...string) (results map[string]bool) {
	var s = getRedisSession()
	defer s.Close()

	results = make(map[string]bool)

	var rIdList = s.SMEMBERS(getGrantKey(destinationId)).MustStrings()
	for _, identifier := range identifiers {
		var pId = md5String(identifier)
		results[identifier] = false

		for _, rId := range rIdList {
			if s.SISMEMBER(getRolePermissionListKey(rId), pId).MustBool() {
				results[identifier] = true
				break
			}
		}

		if s.SISMEMBER(getGrantPrivateKey(destinationId), identifier).MustBool() {
			results[identifier] = true
		}
	}
	return results
}
