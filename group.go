package odin
//
//import (
//	"github.com/smartwalle/going/xid"
//	"github.com/smartwalle/dbr"
//)
//
//const (
//	k_ODIN_GROUP_LIST_PREFIX = "odin_gps_"
//	k_ODIN_GROUP_PREFIX = "odin_gp_"
//)
//
//func getGroupKey(id string) string {
//	return k_ODIN_GROUP_PREFIX + id
//}
//
//func getGroupListKey(gTyep string) string {
//	return k_ODIN_GROUP_LIST_PREFIX + gTyep
//}
//
//////////////////////////////////////////////////////////////////////////////////
//func getGroupList(gType string) (results []*Group, err error) {
//	var s = getSession()
//	defer s.Close()
//
//	var key = getGroupListKey(gType)
//	groupIds, err := s.SMEMBERS(key).Strings()
//	if err != nil {
//		return nil, err
//	}
//	for _, id := range groupIds {
//		var r = s.HGETALL(getGroupKey(id))
//		if r.Error != nil {
//			return nil, err
//		}
//		var group Group
//		if err = r.ScanStruct(&group); err != nil {
//			return nil, err
//		}
//		results = append(results, &group)
//	}
//	return results, err
//}
//
//func GetPermissionGroupList() (results []*Group, err error) {
//	return getGroupList(k_ODIN_GROUP_TYPE_PERMISSION)
//}
//
//func GetRoleGroupList() (results []*Group, err error) {
//	return getGroupList(k_ODIN_GROUP_TYPE_ROLE)
//}
//
//////////////////////////////////////////////////////////////////////////////////
//func GetGroupWithId(id string) (results *Group, err error) {
//	var s = getSession()
//	defer s.Close()
//
//	var r = s.HGETALL(getGroupKey(id))
//	if r.Error != nil {
//		return nil, err
//	}
//	var group Group
//	if err = r.ScanStruct(&group); err != nil {
//		return nil, err
//	}
//	results = &group
//	return results, err
//}
//
//////////////////////////////////////////////////////////////////////////////////
//func newGroup(gType, name string) (id string, err error) {
//	var s = getSession()
//	defer s.Close()
//
//	if r := s.Send("MULTI"); r.Error != nil {
//		return "", r.Error
//	}
//
//	var group = &Group{}
//	group.Id = xid.NewMID().Hex()
//	group.Type = gType
//	group.Name = name
//	if r := s.Send("HMSET", dbr.StructToArgs(getGroupKey(group.Id), group)...); r.Error != nil {
//		return "", r.Error
//	}
//	if r := s.Send("SADD", getGroupListKey(gType), group.Id); r.Error != nil {
//		return "", r.Error
//	}
//
//	if r := s.Do("EXEC"); r.Error != nil {
//		return "", r.Error
//	}
//	id = group.Id
//	return id, err
//}
//
//func NewPermissionGroup(name string) (id string, err error) {
//	return newGroup(k_ODIN_GROUP_TYPE_PERMISSION, name)
//}
//
//func NewRoleGroup(name string) (id string, err error) {
//	return newGroup(k_ODIN_GROUP_TYPE_ROLE, name)
//}
//
//////////////////////////////////////////////////////////////////////////////////
//func UpdateGroup(id, name string) (err error) {
//	var s = getSession()
//	defer s.Close()
//
//	var group *Group
//	if group, err = GetGroupWithId(id); err != nil {
//		return err
//	}
//	if group == nil {
//		return nil
//	}
//
//	if r := s.HMSET(getGroupKey(id), "name", name); r.Error != nil {
//		return r.Error
//	}
//	return nil
//}
//
//////////////////////////////////////////////////////////////////////////////////
//func removeGroup(gType string, id string) (err error) {
//	var s = getSession()
//	defer s.Close()
//
//	var key = getGroupListKey(gType)
//
//
//	if r := s.Send("MULTI"); r.Error != nil {
//		return r.Error
//	}
//
//	if r := s.Send("SREM", key, id); r.Error != nil {
//		return r.Error
//	}
//	if r := s.Send("DEL", getGroupKey(id)); r.Error != nil {
//		return r.Error
//	}
//	if r := s.Do("EXEC"); r.Error != nil {
//		return r.Error
//	}
//	return err
//}
//
//func RemovePremissionGroup(id string) (err error) {
//	return removeGroup(k_ODIN_GROUP_TYPE_PERMISSION, id)
//}
//
//func RemoveRoleGroup(id string) (err error) {
//	return removeGroup(k_ODIN_GROUP_TYPE_ROLE, id)
//}