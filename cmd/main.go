package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/smartwalle/dbr"
	"github.com/smartwalle/odin/service"
	"github.com/smartwalle/odin/service/repository/mysql"
	"github.com/smartwalle/odin/service/repository/redis"
)

func main() {
	var db, _ = sql.Open("mysql", "root:yangfeng@(192.168.1.99:3306)/odin?parseTime=true")
	var r = dbr.NewRedis("192.168.1.99:6379", "", 1, 30, 10)

	var sRepo = mysql.NewOdinRepository(db, "odin")
	var rRepo = redis.NewOdinRepository(r, "odin", sRepo)
	var s = service.NewOdinService(rRepo)

	//s.AddPermissionGroup(0, "S-PG1", odin.K_STATUS_ENABLE)
	//s.AddPermissionGroup(0, "S-PG2", odin.K_STATUS_ENABLE)
	//s.AddPermissionGroup(0, "S-PG3", odin.K_STATUS_ENABLE)
	//
	//s.AddRoleGroup(0, "S-RG1", odin.K_STATUS_ENABLE)
	//s.AddRoleGroup(0, "S-RG2", odin.K_STATUS_ENABLE)
	//s.AddRoleGroup(0, "S-RG3", odin.K_STATUS_ENABLE)

	//s.AddPermissionGroup(2, "PG-2-2", odin.K_STATUS_ENABLE)
	//s.AddPermissionGroup(2, "PG-2-3", odin.K_STATUS_ENABLE)
	//
	//s.AddRoleGroup(2, "RG-2-1", odin.K_STATUS_ENABLE)
	//s.AddRoleGroup(2, "RG-2-2", odin.K_STATUS_ENABLE)
	//s.AddRoleGroup(2, "RG-2-3", odin.K_STATUS_ENABLE)

	//fmt.Println("----- 添加权限组 -----")
	//fmt.Println(s.AddPermissionGroup("用户管理", odin.K_STATUS_ENABLE))
	//fmt.Println(s.AddPermissionGroup("商品管理", odin.K_STATUS_ENABLE))
	//fmt.Println(s.AddPermissionGroup("订单管理", odin.K_STATUS_ENABLE))
	//fmt.Println(s.AddPermissionGroup("文件管理", odin.K_STATUS_ENABLE))
	//
	//fmt.Println("----- 添加角色组 -----")
	//s.AddRoleGroup("用户管理", odin.K_STATUS_ENABLE)
	//s.AddRoleGroup("商品管理", odin.K_STATUS_ENABLE)
	//s.AddRoleGroup("订单管理", odin.K_STATUS_ENABLE)
	//s.AddRoleGroup("文件管理", odin.K_STATUS_ENABLE)
	//s.AddRoleGroup("其它", odin.K_STATUS_ENABLE)
	//
	//fmt.Println("----- 添加权限信息 -----")
	//fmt.Println(s.AddPermission(1, "添加用户", "add_user", odin.K_STATUS_ENABLE))
	//fmt.Println(s.AddPermission(1, "修改用户", "update_user", odin.K_STATUS_ENABLE))
	//fmt.Println(s.AddPermission(1, "删除用户", "delete_user", odin.K_STATUS_ENABLE))
	//fmt.Println(s.AddPermission(2, "添加商品", "add_product", odin.K_STATUS_ENABLE))
	//fmt.Println(s.AddPermission(2, "修改商品", "update_product", odin.K_STATUS_ENABLE))
	//fmt.Println(s.AddPermission(2, "删除商品", "delete_product", odin.K_STATUS_ENABLE))
	//fmt.Println(s.AddPermission(3, "添加订单", "add_order", odin.K_STATUS_ENABLE))
	//fmt.Println(s.AddPermission(3, "修改订单", "update_order", odin.K_STATUS_ENABLE))
	//fmt.Println(s.AddPermission(3, "删除订单", "delete_order", odin.K_STATUS_ENABLE))
	//fmt.Println(s.AddPermission(4, "添加文件", "add_file", odin.K_STATUS_ENABLE))
	//fmt.Println(s.AddPermission(4, "修改文件", "update_file", odin.K_STATUS_ENABLE))
	//fmt.Println(s.AddPermission(4, "删除文件", "delete_file", odin.K_STATUS_ENABLE))
	//
	//fmt.Println("----- 添加角色信息 -----")
	//fmt.Println(s.AddRole(5, "用户管理员", odin.K_STATUS_ENABLE))
	//fmt.Println(s.AddRole(6, "商品管理员", odin.K_STATUS_ENABLE))
	//fmt.Println(s.AddRole(7, "订单管理员", odin.K_STATUS_ENABLE))
	//fmt.Println(s.AddRole(8, "文件管理员", odin.K_STATUS_ENABLE))
	//fmt.Println(s.AddRole(9, "管理员", odin.K_STATUS_ENABLE))
	//
	//fmt.Println("----- 授权角色信息 -----")
	//fmt.Println(s.GrantRole("1", 1, 2, 3, 4, 5, 6, 7, 8))
	//fmt.Println(s.GrantRole("2", 1, 2, 7 ,8))
	//
	//fmt.Println("----- 授权权限信息 -----")
	//fmt.Println(s.GrantPermission(1, 1, 2, 5, 6))
	//fmt.Println(s.GrantPermission(2, 1, 2))
	//
	//fmt.Println("----- 验证权限信息 -----")
	//fmt.Println(s.Check("1", "add_user"))
	//fmt.Println(s.Check("1", "update_user"))
	//fmt.Println(s.Check("1", "delete_user"))
	//
	//fmt.Println(s.Check("2", "add_user"))
	//fmt.Println(s.Check("2", "update_user"))
	//fmt.Println(s.Check("2", "delete_user"))

	fmt.Println("----- 获取权限组列表 -----")
	var gl, _ = s.GetPermissionTree(1, 7, 0, "")
	for _, g := range gl {
		fmt.Println(g.Name)
		for _, p := range g.PermissionList {
			fmt.Println("-", p.Name, p.Identifier, p.Granted)
		}
	}

	//pl, _ := s.GetGrantedPermissionList(1, "111")
	//for _, p := range pl {
	//	fmt.Println(p.Id, p.Identifier, p.Granted)
	//}

	//fmt.Println("----- 获取角色组列表 -----")
	//gl, _ := s.GetRoleTree(1, "1", 0, "")
	//for _, g := range gl {
	//	fmt.Println(g.Name)
	//	for _, r := range g.RoleList {
	//		fmt.Println("-", r.Name, r.Granted, r.Status)
	//		pl, _ := s.GetPermissionListWithRole(1, r.Id)
	//		for _, p := range pl {
	//			fmt.Println("--", p.Name, p.Identifier, p.Granted, p.Status)
	//		}
	//	}
	//}
	//

	//fmt.Println("----- 权限授权列表 -----")
	//pl, _ := s.GetGrantedRoleList("1")
	//for _, p := range pl {
	//	fmt.Println(p.Id, p.Name, p.Granted)
	//}

	//fmt.Println(s.CheckList("1", "update_product", "ssss"))

	//fmt.Println("----- 获取已授权权限列表 -----")
	//var pl, _ = s.GetGrantedPermissionList(1, "111")
	//for _, p := range pl {
	//	fmt.Println(p.Id, p.Name, p.Identifier, p.Granted)
	//}
}
