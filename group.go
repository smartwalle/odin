package odin

import (
	"github.com/smartwalle/going/xid"
	"github.com/smartwalle/dbr"
	"errors"
)

const (
	k_ODIN_GROUP_LIST_PREFIX = "odin_gps_"
	k_ODIN_GROUP_PREFIX = "odin_gp_"
)

func getGroupKey(id string) string {
	return k_ODIN_GROUP_PREFIX + id
}

func getGroupListKey(gTyep string) string {
	return k_ODIN_GROUP_LIST_PREFIX + gTyep
}

func GetGroupList(gType string) (results []*Group, err error) {
	var s = getSession()
	defer s.Close()

	var key = getGroupListKey(gType)
	groupIds, err := s.ZRANGE(key, 0, -1).Strings()
	if err != nil {
		return nil, err
	}
	for _, id := range groupIds {
		var r = s.HGETALL(getGroupKey(id))
		if r.Error != nil {
			return nil, err
		}
		if len(r.MustValues()) == 0 {
			continue
		}
		var group Group
		if err = r.ScanStruct(&group); err != nil {
			return nil, err
		}

		results = append(results, &group)
	}
	return results, err
}

func GetPermissionGroupList() (results []*Group, err error) {
	return GetGroupList(k_ODIN_GROUP_TYPE_PERMISSION)
}

func GetRoleGroupList() (results []*Group, err error) {
	return GetGroupList(k_ODIN_GROUP_TYPE_ROLE)
}

func GetGroupWithId(id string) (results *Group, err error) {
	var s = getSession()
	defer s.Close()

	var r = s.HGETALL(getGroupKey(id))
	if r.Error != nil {
		return nil, err
	}
	if len(r.MustValues()) == 0 {
		return nil, errors.New("Group not exists.")
	}

	var group Group
	if err = r.ScanStruct(&group); err != nil {
		return nil, err
	}

	results = &group
	return results, err
}

func NewGroup(gType, name string) (id string, err error) {
	var s = getSession()
	defer s.Close()

	if r := s.Send("MULTI"); r.Error != nil {
		return "", r.Error
	}

	var group = &Group{}
	group.Id = xid.NewMID().Hex()
	group.Type = gType
	group.Name = name
	if r := s.Send("HMSET", dbr.StructToArgs(getGroupKey(group.Id), group)...); r.Error != nil {
		return "", r.Error
	}
	if r := s.Send("ZADD", getGroupListKey(gType), 0, group.Id); r.Error != nil {
		return "", r.Error
	}

	if r := s.Do("EXEC"); r.Error != nil {
		return "", r.Error
	}
	id = group.Id
	return id, err
}

func UpdateGroup(id, name string) (err error) {
	var s = getSession()
	defer s.Close()

	if _, err = GetGroupWithId(id); err != nil {
		return err
	}

	if r := s.HMSET(getGroupKey(id), "name", name); r.Error != nil {
		return r.Error
	}
	return nil
}