package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/smartwalle/dbr"
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/odin"
	"github.com/smartwalle/odin/service/repository/mysql"
	"github.com/smartwalle/odin/service/repository/redis"
)

func main() {
	var db, _ = sql.Open("mysql", "root:yangfeng@tcp(127.0.0.1:3306)/test?parseTime=true")
	var r = dbr.NewRedis("127.0.0.1:6379", 30, 10, dbr.DialDatabase(1))

	dbs.SetLogger(nil)

	var sRepo = mysql.NewRepository(db, "v2")
	var rRepo = redis.NewRepository(r, "v2", sRepo)
	var s = odin.NewService(rRepo)

	// 添加权限组
	fmt.Println(s.AddPermissionGroup(1, "pg1", "权限组1", odin.Enable))
	fmt.Println(s.AddPermissionGroup(1, "pg2", "权限组2", odin.Enable))
	fmt.Println(s.AddPermissionGroup(1, "pg3", "权限组3", odin.Enable))

	// 添加权限信息
	fmt.Println(s.AddPermission(1, "pg1", "pg1-p1", "权限组1下的权限1", "", odin.Enable))
	fmt.Println(s.AddPermission(1, "pg1", "pg1-p2", "权限组1下的权限2", "", odin.Enable))
	fmt.Println(s.AddPermission(1, "pg1", "pg1-p3", "权限组1下的权限3", "", odin.Enable))
	fmt.Println(s.AddPermission(1, "pg1", "pg1-p4", "权限组1下的权限4", "", odin.Enable))

	fmt.Println(s.AddPermission(1, "pg2", "pg2-p1", "权限组2下的权限1", "", odin.Enable))
	fmt.Println(s.AddPermission(1, "pg2", "pg2-p2", "权限组2下的权限2", "", odin.Enable))
	fmt.Println(s.AddPermission(1, "pg2", "pg2-p3", "权限组2下的权限3", "", odin.Enable))
	fmt.Println(s.AddPermission(1, "pg2", "pg2-p4", "权限组2下的权限4", "", odin.Enable))

	fmt.Println(s.AddPermission(1, "pg3", "pg3-p1", "权限组3下的权限1", "", odin.Enable))
	fmt.Println(s.AddPermission(1, "pg3", "pg3-p2", "权限组3下的权限2", "", odin.Enable))
	fmt.Println(s.AddPermission(1, "pg3", "pg3-p3", "权限组3下的权限3", "", odin.Enable))
	fmt.Println(s.AddPermission(1, "pg3", "pg3-p4", "权限组3下的权限4", "", odin.Enable))

	// 添加角色信息
	fmt.Println(s.AddRole(1, "r1", "角色1", "", odin.Enable))
	fmt.Println(s.AddRole(1, "r2", "角色2", "", odin.Enable))
	fmt.Println(s.AddRole(1, "r3", "角色3", "", odin.Enable))
	fmt.Println(s.AddRole(1, "r4", "角色4", "", odin.Enable))

	// 授予权限给角色
	fmt.Println(s.GrantPermission(1, "r1", "pg1-p1", "pg2-p2", "pg3-p3"))
	fmt.Println(s.GrantPermission(1, "r4", "pg1-p1", "pg2-p2", "pg3-p3"))

	fmt.Println("---")
	groupList, _ := s.GetPermissionTreeWithRole(1, "r1", odin.Enable)
	for _, group := range groupList {
		fmt.Println("-", group.Name, group.AliasName)
		for _, p := range group.PermissionList {
			fmt.Println("--", p.Name, p.AliasName, p.Granted)
		}
	}
	fmt.Println("---")
	groupList, _ = s.GetPermissionTreeWithRoleId(1, 4, odin.Enable)
	for _, group := range groupList {
		fmt.Println("-", group.Name, group.AliasName)
		for _, p := range group.PermissionList {
			fmt.Println("--", p.Name, p.AliasName, p.Granted)
		}
	}

	s.GrantRole(1, "1", "r1", "r2", "r4")

	permissions, err := s.GetGrantedPermissions(1, "1")
	fmt.Println(err)
	for _, p := range permissions {
		fmt.Println(p.Name, p.AliasName, p.Granted)
	}

	roles, err := s.GetGrantedRoles(1, "1")
	for _, r := range roles {
		fmt.Println(r.Name, r.AliasName, r.Granted)
	}

	fmt.Println(s.Check(1, "1", "pg1-p1"))
	fmt.Println(s.Check(1, "2", "pg1-p1"))
	fmt.Println(s.Check(1, "1", "pg2-p2"))
	fmt.Println(s.Check(1, "1", "pg3-p3"))
	fmt.Println(s.Check(1, "1", "pg3-p1"))
	fmt.Println(s.Check(1, "1", "pg3-p2"))

	s.CleanCache(1, "*")
}
