package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/smartwalle/odin"
)

func main() {
	var db, _ = sql.Open("mysql", "root:yangfeng09@(tw.smartwalle.tk:3306)/sm?parseTime=true")

	var s = odin.NewService(db)

	fmt.Println("----- 获取权限组列表 -----")
	var gl, _ = s.GetPermissionTree(0, "")
	for _, g := range gl {
		fmt.Println(g.Name)
		for _, p := range g.PermissionList {
			fmt.Println("-", p.Name, p.Identifier)
		}
	}

	fmt.Println("----- 获取角色组列表 -----")
	gl, _ = s.GetRoleTree(0, "")
	for _, g := range gl {
		fmt.Println(g.Name)
		for _, r := range g.RoleList {
			fmt.Println("-", r.Name)
			pl, _ := s.GetPermissionListWithRole(r.Id)
			for _, p := range pl {
				fmt.Println("--", p.Name, p.Identifier)
			}
		}
	}

	fmt.Println(s.Check("1", "pp1"))
	fmt.Println(s.Check("1", "pp2"))
	fmt.Println(s.Check("1", "pp3"))
	fmt.Println(s.Check("1", "pp4"))
	fmt.Println(s.Check("1", "pp5"))
	fmt.Println(s.Check("2", "pp1"))
}
