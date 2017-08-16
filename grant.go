package odin

import (
	"time"
	"github.com/smartwalle/dbr"
)

const (
	k_ODIN_GRANT_ROLE_PREFIX       = "odin_grant_"
	k_ODIN_GRANT_LIST              = "odin_grant_list"
	k_ODIN_GRANT_PERMISSION_PREFIX = "odin_grant_permission_"
	k_ODIN_GRANT_PERMISSION_LIST   = "odin_grant_permission_list"
)

func getGrantRoleKey(id string) string {
	return k_ODIN_GRANT_ROLE_PREFIX + id
}

func getGrantPermissionKey(id string) string {
	return k_ODIN_GRANT_PERMISSION_PREFIX + id
}

// RevokeAllGrant 清除所有的角色授权信息.
func RevokeAllGrant() error {
	var s = getRedisSession()
	defer s.Close()

	gIdList, err := s.ZREVRANGE(k_ODIN_GRANT_LIST, 0, -1).Strings()
	if err != nil {
		return err
	}

	pIdList, err := s.ZREVRANGE(k_ODIN_GRANT_PERMISSION_LIST, 0, -1).Strings()
	if err != nil {
		return err
	}

	if r := s.Send("MULTI"); r.Error != nil {
		return r.Error
	}

	for _, gId := range gIdList {
		s.Send("DEL", getGrantRoleKey(gId))
	}
	s.Send("DEL", k_ODIN_GRANT_LIST)

	for _, pId := range pIdList {
		s.Send("DEL",  getGrantPermissionKey(pId))
	}
	s.Send("DEL", k_ODIN_GRANT_PERMISSION_LIST)

	if r := s.Do("EXEC"); r.Error != nil {
		return r.Error
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// GetGrantedRoleList 获取所有的角色授权信息.
func GetGrantedRoleList() (results []*GrantInfo, err error) {
	var s = getRedisSession()
	defer s.Close()

	gIdList, err := s.ZREVRANGE(k_ODIN_GRANT_LIST, 0, -1).Strings()
	if err != nil {
		return nil, err
	}

	for _, gId := range gIdList {
		var gInfo = &GrantInfo{}
		gInfo.DestinationId = gId
		gInfo.RoleIdList = s.SMEMBERS(getGrantRoleKey(gId)).MustStrings()
		results = append(results, gInfo)
	}

	return results, err
}

////////////////////////////////////////////////////////////////////////////////
// GetGrantedStandalonePermissionList 获取所有的权限授权信息.
func GetGrantedStandalonePermissionList() (results []*GrantInfo, err error) {
	var s = getRedisSession()
	defer s.Close()

	gIdList, err := s.ZREVRANGE(k_ODIN_GRANT_PERMISSION_LIST, 0, -1).Strings()
	if err != nil {
		return nil, err
	}

	for _, gId := range gIdList {
		var gInfo = &GrantInfo{}
		gInfo.DestinationId = gId
		gInfo.PermissionList = s.SMEMBERS(getGrantPermissionKey(gId)).MustStrings()
		results = append(results, gInfo)
	}
	return results, err
}

////////////////////////////////////////////////////////////////////////////////
// GrantRole 为 destinationId 授权角色信息.
func GrantRole(destinationId string, roleIds ...string) (err error) {
	var s = getRedisSession()
	defer s.Close()
	return grantRole(s, destinationId, roleIds...)
}

func grantRole(s *dbr.Session, destinationId string, roleIds ...string) (err error) {
	var key = getGrantRoleKey(destinationId)

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

// GrantStandalonePermission 为 destinationId 授权某个独立的权限信息.
func GrantStandalonePermission(destinationId string, identifiers ...string) (err error) {
	var s = getRedisSession()
	defer s.Close()

	return grantStandalonePermission(s, destinationId, identifiers...)
}

func grantStandalonePermission(s *dbr.Session, destinationId string, identifiers ...string) (err error) {
	if r := s.Send("MULTI"); r.Error != nil {
		return r.Error
	}
	var key = getGrantPermissionKey(destinationId)
	if len(identifiers) > 0 {
		var params []interface{}
		params = append(params, key)
		for _, identifier := range identifiers {
			params = append(params, identifier)
		}
		s.Send("SADD", params...)

		s.Send("ZADD", k_ODIN_GRANT_PERMISSION_LIST, time.Now().Unix(), destinationId)
	}

	if r := s.Do("EXEC"); r.Error != nil {
		return r.Error
	}
	return err
}

// Grant 为 destinationId 授权角色信息和独立权限信息.
func Grant(destinationId string, roleIds, permissions []string) (err error) {
	var s = getRedisSession()
	defer s.Close()

	if len(roleIds) > 0 {
		err = grantRole(s, destinationId, roleIds...)
		if err != nil {
			return err
		}
	}

	if len(permissions) > 0 {
		err = grantStandalonePermission(s, destinationId, permissions...)
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
		var key = getGrantRoleKey(destinationId)
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

// RevokeStandalonePermission 取消对 destinationId 的独立权限授权.
func RevokeStandalonePermission(destinationId string, identifiers ...string) (err error) {
	var s = getRedisSession()
	defer s.Close()

	if r := s.Send("MULTI"); r.Error != nil {
		return r.Error
	}

	if len(identifiers) > 0 {
		var key = getGrantPermissionKey(destinationId)
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

	var key = getGrantRoleKey(destinationId)
	s.Send("DEL", key)
	s.Send("ZREM", k_ODIN_GRANT_LIST, destinationId)

	if r := s.Do("EXEC"); r.Error != nil {
		return r.Error
	}
	return err
}

// RevokeAllStandalonePermission 取消对 destinationId 所有独立权限授权.
func RevokeAllStandalonePermission(destinationId string) (err error) {
	var s = getRedisSession()
	defer s.Close()

	if r := s.Send("MULTI"); r.Error != nil {
		return r.Error
	}

	var key = getGrantPermissionKey(destinationId)
	s.Send("DEL", key)
	s.Send("ZREM", k_ODIN_GRANT_PERMISSION_LIST, destinationId)

	if r := s.Do("EXEC"); r.Error != nil {
		return r.Error
	}
	return err
}

// Revoke 取消对 destinationId 的所有角色、独立权限的授权.
func Revoke(destinationId string) (err error) {
	var s = getRedisSession()
	defer s.Close()

	if r := s.Send("MULTI"); r.Error != nil {
		return r.Error
	}

	// 取消独立权限授权
	var key = getGrantPermissionKey(destinationId)
	s.Send("DEL", key)
	s.Send("ZREM", k_ODIN_GRANT_PERMISSION_LIST, destinationId)

	// 取消角色授权
	key = getGrantRoleKey(destinationId)
	s.Send("DEL", key)
	s.Send("ZREM", k_ODIN_GRANT_LIST, destinationId)

	if r := s.Do("EXEC"); r.Error != nil {
		return r.Error
	}

	return err
}

////////////////////////////////////////////////////////////////////////////////
// GetGrantedStandalonePermissions 获取 destinationId 拥有的所有独立权限信息.
func GetGrantedStandalonePermissions(destinationId string) (results []string, err error) {
	var s = getRedisSession()
	defer s.Close()

	var pList = s.SMEMBERS(getGrantPermissionKey(destinationId)).MustStrings()
	results = append(results, pList...)

	return results, err
}

// GetGrantedPermissions 获取 destinationId 拥有的所有权限信息，包括独立权限和角色关联的权限.
func GetGrantedPermissions(destinationId string) (results []string, err error) {
	var s = getRedisSession()
	defer s.Close()

	var rIdList = s.SMEMBERS(getGrantRoleKey(destinationId)).MustStrings()
	for _, rId := range rIdList {
		var pIdList = s.SMEMBERS(getRolePermissionListKey(rId)).MustStrings()
		for _, pId := range pIdList {
			var p, _ = getPermission(s, pId)
			if p != nil {
				results = append(results, p.Identifier)
			}
		}
	}

	var pList = s.SMEMBERS(getGrantPermissionKey(destinationId)).MustStrings()
	results = append(results, pList...)

	return results, err
}

// GetGrantedRoles 获取 destinationId 拥有的所有角色信息.
func GetGrantedRoles(destinationId string) (results []*Role, err error) {
	var s = getRedisSession()
	defer s.Close()

	rIdList := s.SMEMBERS(getGrantRoleKey(destinationId)).MustStrings()
	for _, rId := range rIdList {
		var role, _ = GetRoleWithId(rId)
		if role != nil {
			results = append(results, role)
		}
	}

	return results, err
}

////////////////////////////////////////////////////////////////////////////////
// Check 验证 destinationId 是否有访问指定信息的权限.
func Check(destinationId, identifier string) bool {
	var s = getRedisSession()
	defer s.Close()

	var rIdList = s.SMEMBERS(getGrantRoleKey(destinationId)).MustStrings()
	if len(rIdList) > 0 {
		var pId = md5String(identifier)
		for _, rId := range rIdList {
			if s.SISMEMBER(getRolePermissionListKey(rId), pId).MustBool() {
				return true
			}
		}
	}

	if s.SISMEMBER(getGrantPermissionKey(destinationId), identifier).MustBool() {
		return true
	}

	// 如果验证一项不存在的权限信息，那么将返回 true.
	//if s.ZSCORE(k_ODIN_PERMISSION_LIST, pId).Data == nil {
	//	return true
	//}

	return false
}

// CheckRole 验证 destinationId 是否拥有指定角色.
func CheckRole(destinationId, roleId string) bool {
	var s = getRedisSession()
	defer s.Close()

	var key = getGrantRoleKey(destinationId)
	return s.SISMEMBER(key, roleId).MustBool()
}

func CheckList(destinationId string, identifiers ...string) (results map[string]bool) {
	var s = getRedisSession()
	defer s.Close()

	results = make(map[string]bool)

	var rIdList = s.SMEMBERS(getGrantRoleKey(destinationId)).MustStrings()
	for _, identifier := range identifiers {
		var pId = md5String(identifier)
		results[identifier] = false

		for _, rId := range rIdList {
			if s.SISMEMBER(getRolePermissionListKey(rId), pId).MustBool() {
				results[identifier] = true
				break
			}
		}

		if s.SISMEMBER(getGrantPermissionKey(destinationId), identifier).MustBool() {
			results[identifier] = true
		}
	}
	return results
}
