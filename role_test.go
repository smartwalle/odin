package odin

import (
	"testing"
	"fmt"
)

//func TestNewGroup(t *testing.T) {
//	var id, err = NewGroup(k_ODIN_GROUP_TYPE_PERMISSION, "haha")
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println(id)
//}

func TestGetPermissionGroupList(t *testing.T) {
	var groupList, err = GetPermissionGroupList()
	if err != nil {
		fmt.Println(err)
	}

	for _, group := range groupList {
		fmt.Println(group.Id, group.Type, group.Name)
	}
}

func TestGetGroupWithId(t *testing.T) {
	var group, err = GetGroupWithId("5a34a7d63173")
	if err != nil {
		fmt.Println(err)
	}
	if group != nil {
		fmt.Println(group.Id, group.Type, group.Name)
	}
}

func TestUpdateGroup(t *testing.T) {
	fmt.Println(UpdateGroup("5a34a7d63173", "gege"))
}