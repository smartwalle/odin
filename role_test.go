package odin

import (
	"testing"
	"fmt"
)

func TestNewGroup(t *testing.T) {
	var id, err = NewPermissionGroup("haha")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(id)
}

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
		fmt.Println("aaa", group.Id, group.Type, group.Name)
	}
}

func TestUpdateGroup(t *testing.T) {
	fmt.Println(UpdateGroup("5a34a7d63173", "gege"))
}

func TestRemovePremissionGroup(t *testing.T) {
	fmt.Println(RemovePremissionGroup("5af797a1f9c3"))
}