package odin

import (
	"github.com/smartwalle/dbr"
	"github.com/smartwalle/going/xid"
	"strings"
)

const (
	k_ODIN_ROLE_PREFIX                 = "odin_ro_"
	k_ODIN_ROLE_LIST                   = "odin_ro_list"
	k_ODIN_ROLE_PERMISSION_LIST_PREFIX = "odin_rp_"
	k_ODIN_GRANT_PREFIX                = "odin_grant_"
	k_ODIN_GRANT_LIST                  = "grant_list"
)

func getRoleKey(id string) string {
	return k_ODIN_ROLE_PREFIX + id
}

func getRolePermissionListKey(id string) string {
	return k_ODIN_ROLE_PERMISSION_LIST_PREFIX + id
}

func getGrantKey(id string) string {
	return k_ODIN_GRANT_PREFIX + id
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

	var params []interface{}
	params = append(params, getRolePermissionListKey(r.Id))
	for _, id := range permissionIds {
		params = append(params, id)
	}
	if r := s.Send("SADD", params...); r.Error != nil {
		return "", r.Error
	}

	if r := s.Send("SADD", k_ODIN_ROLE_LIST, r.Id); r.Error != nil {
		return "", r.Error
	}

	if r := s.Do("EXEC"); r.Error != nil {
		return "", r.Error
	}

	id = r.Id
	return id, err
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

	roleIds, err := s.SMEMBERS(k_ODIN_ROLE_LIST).Strings()
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

	if r := s.Send("SREM", k_ODIN_ROLE_LIST, id); r.Error != nil {
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