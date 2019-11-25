package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/smartwalle/dbr"
	"github.com/smartwalle/odin"
	"github.com/smartwalle/odin/service/repository/mysql"
	"github.com/smartwalle/odin/service/repository/redis"
)

func main() {
	var db, _ = sql.Open("mysql", "root:yangfeng@tcp(192.168.1.99:3306)/tt?parseTime=true")
	var r = dbr.NewRedis("192.168.1.99:6379", 30, 10, dbr.DialDatabase(1))

	//dbs.SetLogger()

	var sRepo = mysql.NewRepository(db, "v2")
	var rRepo = redis.NewRepository(r, "v2", sRepo)
	var s = odin.NewService(rRepo)

	var roleList, err = s.GetRoleList(1, "1", 0, "")
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, role := range roleList {
		fmt.Println(role)
	}

	role, err := s.GetRoleWithName(1, "root")
	fmt.Println(role)

	role, err = s.GetRoleWithId(1, 1)
	fmt.Println(role)

	fmt.Println(s.UpdateRole(1, 8, "update_root2", "更新一下2", "haha2", odin.Enable))

	fmt.Println(s.GrantRoleWithIds(1, "2", 1, 7,8))
	fmt.Println(s.GetGrantedRoleList(1, "1"))
	fmt.Println(s.GetGrantedRoleList(1, "2"))
}
