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
	var db, _ = sql.Open("mysql", "root:youle@tcp(192.168.1.77:3306)/tt?parseTime=true")
	var r = dbr.NewRedis("192.168.1.99:6379", 30, 10, dbr.DialDatabase(1))

	dbs.SetLogger(nil)

	var sRepo = mysql.NewRepository(db, "odin")
	var rRepo = redis.NewRepository(r, "odin", sRepo)
	var s = odin.NewService(rRepo)

	if g, err := s.AddPermissionGroup(1, "pg1", odin.StatusOfEnable); err == nil && g != nil {
		s.AddPermission(1, g.Id, "权限g11", "pg1-p1", odin.StatusOfEnable)
		s.AddPermission(1, g.Id, "权限g12", "pg1-p2", odin.StatusOfEnable)
		s.AddPermission(1, g.Id, "权限g13", "pg1-p3", odin.StatusOfEnable)
	} else {
		fmt.Println(err)
	}
	if g, err := s.AddPermissionGroup(2, "pg2", odin.StatusOfEnable); err == nil && g != nil {
		s.AddPermission(2, g.Id, "权限g21", "pg2-p1", odin.StatusOfEnable)
		s.AddPermission(2, g.Id, "权限g22", "pg2-p2", odin.StatusOfEnable)
		s.AddPermission(2, g.Id, "权限g23", "pg2-p3", odin.StatusOfEnable)
	} else {
		fmt.Println(err)
	}

	if g, err := s.AddRoleGroup(1, "rg1", odin.StatusOfEnable); err == nil && g != nil {
		s.AddRole(1, g.Id, "角色g11", odin.StatusOfEnable)
		s.AddRole(1, g.Id, "角色g12", odin.StatusOfEnable)
		s.AddRole(1, g.Id, "角色g13", odin.StatusOfEnable)
	} else {
		fmt.Println(err)
	}

	if g, err := s.AddRoleGroup(2, "rg2", odin.StatusOfEnable); err == nil && g != nil {
		s.AddRole(2, g.Id, "角色g21", odin.StatusOfEnable)
		s.AddRole(2, g.Id, "角色g22", odin.StatusOfEnable)
		s.AddRole(2, g.Id, "角色g23", odin.StatusOfEnable)
	} else {
		fmt.Println(err)
	}

	fmt.Println(s.GrantPermission(1, 1, 1))
	fmt.Println(s.GrantPermission(1, 1, 4))
	fmt.Println(s.GrantRole(1, "1", 1))
	fmt.Println(s.GrantRole(0, "1", 1))
}
