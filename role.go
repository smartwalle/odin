package odin

import (
	"github.com/smartwalle/dbr"
	"github.com/smartwalle/going/xid"
	"strings"
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
		for _, id := range permissionIds {
			params = append(params, id)
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
		for _, id := range permissionIds {
			params = append(params, id)
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
func AddPermissionsToRole(id string, permissionIds ...string) (err error) {
	var s = getRedisSession()
	defer s.Close()

	if r := s.Send("MULTI"); r.Error != nil {
		return r.Error
	}

	var params []interface{}
	params = append(params, getRolePermissionListKey(id))
	for _, id := range permissionIds {
		params = append(params, id)
	}
	if r := s.Send("SADD", params...); r.Error != nil {
		return r.Error
	}
	if r := s.Do("EXEC"); r.Error != nil {
		return r.Error
	}
	return err
}

func AddPermissionToRole(id string, identifiers ...string) (err error) {
	var s = getRedisSession()
	defer s.Close()

	if r := s.Send("MULTI"); r.Error != nil {
		return r.Error
	}

	var params []interface{}
	params = append(params, getRolePermissionListKey(id))
	params = append(params, md5String(strings.Join(identifiers, "-")))

	if r := s.Send("SADD", params...); r.Error != nil {
		return r.Error
	}
	if r := s.Do("EXEC"); r.Error != nil {
		return r.Error
	}
	return err
}

////////////////////////////////////////////////////////////////////////////////
func RemovePermissionsFromRole(id string, permissionIds ...string) (err error) {
	var s = getRedisSession()
	defer s.Close()

	if r := s.Send("MULTI"); r.Error != nil {
		return r.Error
	}

	var params []interface{}
	params = append(params, getRolePermissionListKey(id))
	for _, id := range permissionIds {
		params = append(params, id)
	}
	if r := s.Send("SREM", params...); r.Error != nil {
		return r.Error
	}
	if r := s.Do("EXEC"); r.Error != nil {
		return r.Error
	}
	return err
}

func RemovePermissionFromRole(id string, identifiers ...string) (err error) {
	var s = getRedisSession()
	defer s.Close()

	if r := s.Send("MULTI"); r.Error != nil {
		return r.Error
	}

	var params []interface{}
	params = append(params, getRolePermissionListKey(id))
	params = append(params, md5String(strings.Join(identifiers, "-")))

	if r := s.Send("SREM", params...); r.Error != nil {
		return r.Error
	}
	if r := s.Do("EXEC"); r.Error != nil {
		return r.Error
	}
	return err
}

////////////////////////////////////////////////////////////////////////////////
func GetRoleList() (results []*Role, err error) {
	var s = getRedisSession()
	defer s.Close()

	roleIds, err := s.ZREVRANGE(k_ODIN_ROLE_LIST, 0, -1).Strings()
	if err != nil {
		return nil, err
	}
	for _, id := range roleIds {
		var r = s.HGETALL(getRoleKey(id))
		if r.Error != nil {
			return nil, err
		}
		role, err := getRole(s, id)
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

func GetRoleWithId(id string) (results *Role, err error) {
	var s = getRedisSession()
	defer s.Close()
	return getRole(s, id)
}

////////////////////////////////////////////////////////////////////////////////
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