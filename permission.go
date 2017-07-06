package odin

import (
	"github.com/smartwalle/dbr"
	"time"
)

const (
	k_ODIN_PERMISSION_PREFIX = "odin_pn_"
	k_ODIN_PERMISSION_LIST   = "odin_pn_list"
)

func getPermissionKey(id string) string {
	return k_ODIN_PERMISSION_PREFIX + id
}

////////////////////////////////////////////////////////////////////////////////
// NewPermission 创建新的权限信息.
// 权限的 id 是由 identifier 决定的，所以不能添加重复的 identifier.
func NewPermission(group, name, identifier string) (id string, err error) {
	var s = getRedisSession()
	defer s.Close()

	if r := s.Send("MULTI"); r.Error != nil {
		return "", r.Error
	}

	var p = &Permission{}
	p.Identifier = identifier
	p.Name = name
	p.Group = group
	p.Id = md5String(p.Identifier)

	if r := s.Send("HMSET", dbr.StructToArgs(getPermissionKey(p.Id), p)...); r.Error != nil {
		return "", r.Error
	}
	if r := s.Send("ZADD", k_ODIN_PERMISSION_LIST, time.Now().Unix(), p.Id); r.Error != nil {
		return "", r.Error
	}
	if r := s.Do("EXEC"); r.Error != nil {
		return "", r.Error
	}
	id = p.Id
	return id, err
}

////////////////////////////////////////////////////////////////////////////////
// UpdatePermission 更新指定权限的信息，由于权限的 id 是由 identifier 决定的，所以当权限的 identifier 发生改变之后，其 id 也会改变。
// 如果权限的 id 有发生变化，将会对角色的权限列表进行更新，把老的权限 id 从相关的角色中删除，将新的权限 id 添加到相关的角色中.
// 通常情况，不建议使用本方法.
func UpdatePermission(id, group, name, identifier string) (string, error) {
	var s = getRedisSession()
	defer s.Close()

	if r := s.Send("MULTI"); r.Error != nil {
		return "", r.Error
	}

	var p = &Permission{}
	p.Identifier = identifier
	p.Name = name
	p.Group = group
	p.Id = md5String(p.Identifier)

	if r := s.Send("HMSET", dbr.StructToArgs(getPermissionKey(p.Id), p)...); r.Error != nil {
		return "", r.Error
	}
	if r := s.Send("ZADD", k_ODIN_PERMISSION_LIST, time.Now().Unix(), p.Id); r.Error != nil {
		return "", r.Error
	}
	if id != p.Id {
		s.Send("DEL", getPermissionKey(id))
	}
	if r := s.Do("EXEC"); r.Error != nil {
		return "", r.Error
	}

	if id != p.Id {
		roleIds := s.ZRANGE(k_ODIN_ROLE_LIST, 0, -1).MustStrings()
		for _, rId := range roleIds {
			var key = getRolePermissionListKey(rId)
			if s.SISMEMBER(key, id).MustBool() {
				s.SREM(key, id)
				s.SADD(key, p.Id)
			}
		}
	}

	id = p.Id
	return id, nil
}

////////////////////////////////////////////////////////////////////////////////
// GetPermissionList 获取所有的权限列表.
func GetPermissionList() (results []*Permission, err error) {
	var s = getRedisSession()
	defer s.Close()

	pIdList, err := s.ZREVRANGE(k_ODIN_PERMISSION_LIST, 0, -1).Strings()
	if err != nil {
		return nil, err
	}
	for _, id := range pIdList {
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

// GetPermissionWithId 获取指定的权限信息.
func GetPermissionWithId(id string) (results *Permission, err error) {
	var s = getRedisSession()
	defer s.Close()
	return getPermission(s, id)
}

// GetPermission 获取指定的权限信息.
func GetPermission(identifier string) (results *Permission, err error) {
	var s = getRedisSession()
	defer s.Close()

	var id = md5String(identifier)
	return getPermission(s, id)
}

////////////////////////////////////////////////////////////////////////////////
// RemovePermissionWithId 移除指定的权限信息.
func RemovePermissionWithId(id string) (err error) {
	var s = getRedisSession()
	defer s.Close()

	if r := s.Send("MULTI"); r.Error != nil {
		return r.Error
	}

	if r := s.Send("ZREM", k_ODIN_PERMISSION_LIST, id); r.Error != nil {
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

// RemovePermission 移除指定的权限信息.
func RemovePermission(identifier string) (err error) {
	var id = md5String(identifier)
	return RemovePermissionWithId(id)
}

// RemoveAllPermission 移除所有的权限信息.
func RemoveAllPermission() (error){
	var s = getRedisSession()
	defer s.Close()

	pIdList, err := s.ZREVRANGE(k_ODIN_PERMISSION_LIST, 0, -1).Strings()
	if err != nil {
		return err
	}

	if r := s.Send("MULTI"); r.Error != nil {
		return r.Error
	}
	for _, pId := range pIdList {
		s.Send("DEL", getPermissionKey(pId))
	}
	s.Send("DEL", k_ODIN_PERMISSION_LIST)

	if r := s.Do("EXEC"); r.Error != nil {
		return r.Error
	}
	return nil
}