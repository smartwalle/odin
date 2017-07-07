package odin

import (
	"github.com/smartwalle/dbr"
	"github.com/smartwalle/going/xid"
	"time"
)

const (
	k_ODIN_ROLE_PREFIX                 = "odin_ro_"
	k_ODIN_ROLE_LIST                   = "odin_ro_list"
	k_ODIN_ROLE_PERMISSION_LIST_PREFIX = "odin_rp_"
)

func getRoleKey(id string) string {
	return k_ODIN_ROLE_PREFIX + id
}

func getRolePermissionListKey(id string) string {
	return k_ODIN_ROLE_PERMISSION_LIST_PREFIX + id
}

////////////////////////////////////////////////////////////////////////////////
// NewRole 添加新的角色信息, 添加成功将会返回新角色的 id.
func NewRole(group, name string, permissionIds ...string) (id string, err error) {
	var s = getRedisSession()
	defer s.Close()

	if r := s.Send("MULTI"); r.Error != nil {
		return "", r.Error
	}

	var r = &Role{}
	r.Name = name
	r.Group = group
	r.Id = xid.NewMID().Hex()

	if r := s.Send("HMSET", dbr.StructToArgs(getRoleKey(r.Id), r)...); r.Error != nil {
		return "", r.Error
	}

	if len(permissionIds) > 0 {
		var params []interface{}
		params = append(params, getRolePermissionListKey(r.Id))
		for _, pId := range permissionIds {
			params = append(params, pId)
		}
		if r := s.Send("SADD", params...); r.Error != nil {
			return "", r.Error
		}
	}

	if r := s.Send("ZADD", k_ODIN_ROLE_LIST, time.Now().Unix(), r.Id); r.Error != nil {
		return "", r.Error
	}

	if r := s.Do("EXEC"); r.Error != nil {
		return "", r.Error
	}

	id = r.Id
	return id, err
}

////////////////////////////////////////////////////////////////////////////////
// UpdateRole 更新角色信息，如果角色信息不存在，则会根据信息创建新的角色信息。
func UpdateRole(id, group, name string, permissionIds ...string) (err error) {
	var s = getRedisSession()
	defer s.Close()

	if r := s.Send("MULTI"); r.Error != nil {
		return r.Error
	}

	var r = &Role{}
	r.Name = name
	r.Group = group
	r.Id = id

	if r := s.Send("HMSET", dbr.StructToArgs(getRoleKey(r.Id), r)...); r.Error != nil {
		return r.Error
	}

	s.Send("DEL", getRolePermissionListKey(r.Id))

	if len(permissionIds) > 0 {
		var params []interface{}
		params = append(params, getRolePermissionListKey(r.Id))
		for _, pId := range permissionIds {
			params = append(params, pId)
		}
		if r := s.Send("SADD", params...); r.Error != nil {
			return r.Error
		}
	}

	if r := s.Send("ZADD", k_ODIN_ROLE_LIST, time.Now().Unix(), r.Id); r.Error != nil {
		return r.Error
	}

	if r := s.Do("EXEC"); r.Error != nil {
		return r.Error
	}

	return err
}

////////////////////////////////////////////////////////////////////////////////
// AddPermissionsToRole 向角色添加指定的权限信息。
func AddPermissionsToRole(id string, permissionIds ...string) (err error) {
	var s = getRedisSession()
	defer s.Close()

	if r := s.Send("MULTI"); r.Error != nil {
		return r.Error
	}

	var params []interface{}
	params = append(params, getRolePermissionListKey(id))
	for _, pId := range permissionIds {
		params = append(params, pId)
	}
	if r := s.Send("SADD", params...); r.Error != nil {
		return r.Error
	}
	if r := s.Do("EXEC"); r.Error != nil {
		return r.Error
	}
	return err
}

// AddPermissionToRole 向角色添加指定的权限信息。
func AddPermissionToRole(id string, identifiers ...string) (err error) {
	var s = getRedisSession()
	defer s.Close()

	if r := s.Send("MULTI"); r.Error != nil {
		return r.Error
	}

	var params []interface{}
	params = append(params, getRolePermissionListKey(id))
	for _, identifier := range identifiers {
		params = append(params, md5String(identifier))
	}

	if r := s.Send("SADD", params...); r.Error != nil {
		return r.Error
	}
	if r := s.Do("EXEC"); r.Error != nil {
		return r.Error
	}
	return err
}

////////////////////////////////////////////////////////////////////////////////
// RemovePermissionsFromRole 移除指定角色的指令权限.
func RemovePermissionsFromRole(id string, permissionIds ...string) (err error) {
	var s = getRedisSession()
	defer s.Close()

	if r := s.Send("MULTI"); r.Error != nil {
		return r.Error
	}

	var params []interface{}
	params = append(params, getRolePermissionListKey(id))
	for _, pId := range permissionIds {
		params = append(params, pId)
	}
	if r := s.Send("SREM", params...); r.Error != nil {
		return r.Error
	}
	if r := s.Do("EXEC"); r.Error != nil {
		return r.Error
	}
	return err
}

// RemovePermissionFromRole 移除指定角色的指令权限.
func RemovePermissionFromRole(id string, identifiers ...string) (err error) {
	var s = getRedisSession()
	defer s.Close()

	if r := s.Send("MULTI"); r.Error != nil {
		return r.Error
	}

	var params []interface{}
	params = append(params, getRolePermissionListKey(id))
	for _, identifier := range identifiers {
		params = append(params, md5String(identifier))
	}

	if r := s.Send("SREM", params...); r.Error != nil {
		return r.Error
	}
	if r := s.Do("EXEC"); r.Error != nil {
		return r.Error
	}
	return err
}

////////////////////////////////////////////////////////////////////////////////
// GetRoleList 获取所有的角色列表.
func GetRoleList() (results []*Role, err error) {
	var s = getRedisSession()
	defer s.Close()

	rIdList, err := s.ZREVRANGE(k_ODIN_ROLE_LIST, 0, -1).Strings()
	if err != nil {
		return nil, err
	}
	for _, rId := range rIdList {
		role, err := getRole(s, rId)
		if err != nil {
			return nil, err
		}

		if role != nil {
			results = append(results, role)
		}
	}
	return results, err
}

////////////////////////////////////////////////////////////////////////////////
func getRole(s *dbr.Session, id string) (results *Role, err error) {
	var r = s.HGETALL(getRoleKey(id))
	if r.Error != nil {
		return nil, err
	}
	var role Role
	if err = r.ScanStruct(&role); err != nil {
		return nil, err
	}

	if r := s.SMEMBERS(getRolePermissionListKey(id)); r.Error != nil {
		return nil, err
	} else {
		role.PermissionIdList = r.MustStrings()
	}
	results = &role
	return results, err
}

// GetRoleWithId 获取指定的角色信息.
func GetRoleWithId(id string) (results *Role, err error) {
	var s = getRedisSession()
	defer s.Close()
	return getRole(s, id)
}

////////////////////////////////////////////////////////////////////////////////
// RemoveRoleWithId 移除指定的角色信息.
func RemoveRoleWithId(id string) (err error) {
	var s = getRedisSession()
	defer s.Close()

	if r := s.Send("MULTI"); r.Error != nil {
		return r.Error
	}

	if r := s.Send("ZREM", k_ODIN_ROLE_LIST, id); r.Error != nil {
		return r.Error
	}
	if r := s.Send("DEL", getRoleKey(id)); r.Error != nil {
		return r.Error
	}
	if r := s.Send("DEL", getRolePermissionListKey(id)); r.Error != nil {
		return r.Error
	}
	if r := s.Do("EXEC"); r.Error != nil {
		return r.Error
	}
	return err
}

// RemoveAllRole 移除所有角色信息.
func RemoveAllRole() (error){
	var s = getRedisSession()
	defer s.Close()

	rIdList, err := s.ZREVRANGE(k_ODIN_ROLE_LIST, 0, -1).Strings()
	if err != nil {
		return err
	}

	if r := s.Send("MULTI"); r.Error != nil {
		return r.Error
	}
	for _, rId := range rIdList {
		s.Send("DEL", getRoleKey(rId))
		s.Send("DEL", getRolePermissionListKey(rId))
	}
	s.Send("DEL", k_ODIN_ROLE_LIST)

	if r := s.Do("EXEC"); r.Error != nil {
		return r.Error
	}
	return nil
}