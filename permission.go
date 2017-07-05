package odin

import (
	"github.com/smartwalle/dbr"
	"strings"
)

const (
	k_ODIN_PERMISSION_PREFIX = "odin_pn_"
	k_ODIN_PERMISSION_LIST   = "odin_pn_list"
)

func getPermissionKey(id string) string {
	return k_ODIN_PERMISSION_PREFIX + id
}

////////////////////////////////////////////////////////////////////////////////
func NewPermission(group, name string, identifiers ...string) (id string, err error) {
	var s = getRedisSession()
	defer s.Close()

	if r := s.Send("MULTI"); r.Error != nil {
		return "", r.Error
	}

	var p = &Permission{}
	p.Identifier = strings.Join(identifiers, "-")
	p.Name = name
	p.Group = group
	p.Id = md5String(p.Identifier)

	if r := s.Send("HMSET", dbr.StructToArgs(getPermissionKey(p.Id), p)...); r.Error != nil {
		return "", r.Error
	}
	if r := s.Send("SADD", k_ODIN_PERMISSION_LIST, p.Id); r.Error != nil {
		return "", r.Error
	}
	if r := s.Do("EXEC"); r.Error != nil {
		return "", r.Error
	}
	id = p.Id
	return id, err
}

////////////////////////////////////////////////////////////////////////////////
func UpdatePermission(id, group, name string, identifiers ...string) (string, error) {
	var s = getRedisSession()
	defer s.Close()

	if r := s.Send("MULTI"); r.Error != nil {
		return "", r.Error
	}

	var p = &Permission{}
	p.Identifier = strings.Join(identifiers, "-")
	p.Name = name
	p.Group = group
	p.Id = md5String(p.Identifier)

	if r := s.Send("HMSET", dbr.StructToArgs(getPermissionKey(p.Id), p)...); r.Error != nil {
		return "", r.Error
	}
	if r := s.Send("SADD", k_ODIN_PERMISSION_LIST, p.Id); r.Error != nil {
		return "", r.Error
	}
	if id != p.Id {
		s.Send("DEL", getPermissionKey(id))
	}
	if r := s.Do("EXEC"); r.Error != nil {
		return "", r.Error
	}

	roleIds := s.SMEMBERS(k_ODIN_ROLE_LIST).MustStrings()
	for _, rId := range roleIds {
		var key = getRolePermissionListKey(rId)
		if s.SISMEMBER(key, id).MustBool() {
			s.SREM(key, id)
			s.SADD(key, p.Id)
		}
	}

	id = p.Id

	return id, nil
}

////////////////////////////////////////////////////////////////////////////////
func GetPermissionList() (results []*Permission, err error) {
	var s = getRedisSession()
	defer s.Close()

	groupIds, err := s.SMEMBERS(k_ODIN_PERMISSION_LIST).Strings()
	if err != nil {
		return nil, err
	}
	for _, id := range groupIds {
		p, err := getPermission(s, id)
		if err != nil {
			return nil, err
		}

		if p != nil {
			results = append(results, p)
		}
	}
	return results, err
}

////////////////////////////////////////////////////////////////////////////////
func getPermission(s *dbr.Session, id string) (results *Permission, err error) {
	var r = s.HGETALL(getPermissionKey(id))
	if r.Error != nil {
		return nil, err
	}
	var p Permission
	if err = r.ScanStruct(&p); err != nil {
		return nil, err
	}
	results = &p
	return results, err
}

func GetPermissionWithId(id string) (results *Permission, err error) {
	var s = getRedisSession()
	defer s.Close()
	return getPermission(s, id)
}

func GetPermission(identifiers ...string) (results *Permission, err error) {
	var s = getRedisSession()
	defer s.Close()

	var id = md5String(strings.Join(identifiers, "-"))
	return getPermission(s, id)
}

////////////////////////////////////////////////////////////////////////////////
func RemovePermissionWithId(id string) (err error) {
	var s = getRedisSession()
	defer s.Close()

	if r := s.Send("MULTI"); r.Error != nil {
		return r.Error
	}

	if r := s.Send("SREM", k_ODIN_PERMISSION_LIST, id); r.Error != nil {
		return r.Error
	}
	if r := s.Send("DEL", getPermissionKey(id)); r.Error != nil {
		return r.Error
	}
	if r := s.Do("EXEC"); r.Error != nil {
		return r.Error
	}
	return err
}

func RemovePermission(identifier ...string) (err error) {
	var id = md5String(strings.Join(identifier, "-"))
	return RemovePermissionWithId(id)
}