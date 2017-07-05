package main

import (
	"github.com/smartwalle/odin"
	"fmt"
)

func main() {
	odin.Init("localhost:6379", "", 2, 10, 2)

	fmt.Println(odin.NewPermission("产品", "添加产品", "POST", "/api/product"))
	fmt.Println(odin.NewPermission("产品", "修改产品", "PUT", "/api/product"))

	var pList, _ = odin.GetPermissionList()
	for _, p := range pList {
		fmt.Println("权限", p.Id, p.Identifier, p.Group, p.Name)
	}

	odin.AddPermissionsToRole("5b195208b07c", "3b6ebd50650d523009874b9128e33d31")

	var rList, _ = odin.GetRoleList()
	for _, r := range rList {
		fmt.Println("角色", r.Id, r.Group, r.Name, r.PermissionIdList)
	}

	fmt.Println(odin.Check("111", "PUT", "/api/product"))
	fmt.Println(odin.Check("111", "POST", "/api/product"))
	fmt.Println(odin.Check("111", "GET", "/api/product"))

	fmt.Println(odin.GetGrantPermissionList("111"))

}