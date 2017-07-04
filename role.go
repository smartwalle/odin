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
	var s = getSession()
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
	var s = getSession()
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
	var s = getSession()
	defer s.Close()

	if r := s.Send("MULTI"); r.Error != nil {
		return r.Error
	}

	var params []interface{}
	params = append(params, getRolePermissionListKey(id))
	params = append(params, MD5String(strings.Join(identifiers, "-")))

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
	var s = getSession()
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
	var s = getSession()
	defer s.Close()

	if r := s.Send("MULTI"); r.Error != nil {
		return r.Error
	}

	var params []interface{}
	params = append(params, getRolePermissionListKey(id))
	params = append(params, MD5String(strings.Join(identifiers, "-")))

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
	var s = getSession()
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
	var s = getSession()
	defer s.Close()
	return getRole(s, id)
}

////////////////////////////////////////////////////////////////////////////////
func RemoveRoleWithId(id string) (err error) {
	var s = getSession()
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

////////////////////////////////////////////////////////////////////////////////
func Grant(destinationId string, roleIds ...string) (err error) {
	var s = getSession()
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
	var s = getSession()
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
	var s = getSession()
	defer s.Close()

	var roleIds = s.SMEMBERS(getGrantKey(destinationId)).MustStrings()
	for _, roleId := range roleIds {
		var pIdList = s.SMEMBERS(getRolePermissionListKey(roleId)).MustStrings()
		results = append(results, pIdList...)
	}

	return results, err
}

func Check(destinationId string, identifiers ...string) (bool) {
	var s = getSession()
	defer s.Close()

	var id = MD5String(strings.Join(identifiers, "-"))

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