package main

import (
	"github.com/smartwalle/odin"
	"fmt"
)

func main() {
	// 初始化数据库信息
	odin.Init("localhost:6379", "", 5, 10, 2)

	fmt.Println("--------------------------------------------------")
	// 创建权限信息
	fmt.Println("初始化权限信息...")
	var p11Id, _ = odin.NewPermission("产品", "添加产品", "POST-/api/product")
	var p12Id, _ = odin.NewPermission("产品", "修改产品", "PUT-/api/product")

	var p21Id, _ = odin.NewPermission("用户", "添加用户", "POST-/api/user")
	var p22Id, _ = odin.NewPermission("用户", "修改产品", "PUT-/api/user")
	fmt.Println("初始化权限信息完成...")


	var pList, _ = odin.GetPermissionList()
	fmt.Println("现有权限信息:")
	for _, p := range pList {
		fmt.Println(p.Id, p.Identifier, p.Group, p.Name)
	}

	fmt.Println("--------------------------------------------------")
	odin.RemoveAllRole()
	// 创建角色信息
	fmt.Println("初始化角色信息...")
	var r1Id, _ = odin.NewRole("产品组", "产品管理员", p11Id, p12Id)
	var r2Id, _ = odin.NewRole("用户组", "用户管理员", p21Id, p22Id)
	fmt.Println("初始化角色信息完成...")

	fmt.Println("现有角色信息:")
	var rList, _ = odin.GetRoleList()
	for _, r := range rList {
		fmt.Println(r.Id, r.Name, r.Group, r.PermissionIdList)
	}

	fmt.Println("--------------------------------------------------")
	// 为指定对象授权
	var userId1 = "user_id_001"
	var userId2 = "user_id_002"
	var userId3 = "user_id_003"

	// 因为每次都创建了新的角色信息，原有角色信息会被清楚，所以先取消原有授权信息
	odin.RevokeAllRole(userId1)
	odin.RevokeAllRole(userId2)
	odin.RevokeAllRole(userId3)

	odin.GrantRole(userId1, r1Id)
	odin.GrantRole(userId2, r2Id)
	odin.GrantRole(userId3, r1Id, r2Id)

	fmt.Println("授权信息:")
	var gList, _ = odin.GetAllGrantRoleList()
	for _, g := range gList {
		fmt.Println(g.DestinationId, g.RoleIdList)
	}

	fmt.Println("--------------------------------------------------")
	fmt.Println(userId1, "POST", "/api/product", odin.Check(userId1, "POST-/api/product"))
	fmt.Println(userId1, "PUT", "/api/product", odin.Check(userId1, "PUT-/api/product"))
	fmt.Println(userId1, "POST", "/api/user", odin.Check(userId1, "POST-/api/user"))
	fmt.Println(userId1, "PUT", "/api/user", odin.Check(userId1, "PUT-/api/user"))
	fmt.Println("")
	fmt.Println(userId2, "POST", "/api/product", odin.Check(userId2, "POST-/api/product"))
	fmt.Println(userId2, "PUT", "/api/product", odin.Check(userId2, "PUT-/api/product"))
	fmt.Println(userId2, "POST", "/api/user", odin.Check(userId2, "POST-/api/user"))
	fmt.Println(userId2, "PUT", "/api/user", odin.Check(userId2, "PUT-/api/user"))
	fmt.Println("")
	fmt.Println(userId3, "POST", "/api/product", odin.Check(userId3, "POST-/api/product"))
	fmt.Println(userId3, "PUT", "/api/product", odin.Check(userId3, "PUT-/api/product"))
	fmt.Println(userId3, "POST", "/api/user", odin.Check(userId3, "POST-/api/user"))
	fmt.Println(userId3, "PUT", "/api/user", odin.Check(userId3, "PUT-/api/user"))
	fmt.Println("")
	fmt.Println(userId1, "GET", "/api/users", odin.Check(userId1, "GET-/api/users"))
	fmt.Println(userId2, "GET", "/api/users", odin.Check(userId2, "GET-/api/users"))
	fmt.Println(userId3, "GET", "/api/users", odin.Check(userId3, "GET-/api/users"))
}
