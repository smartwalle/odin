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
	var db, _ = sql.Open("mysql", "root:yangfeng@tcp(192.168.1.99:3306)/tt?parseTime=true")
	var r = dbr.NewRedis("192.168.1.99:6379", 30, 10, dbr.DialDatabase(1))

	dbs.SetLogger(nil)

	var sRepo = mysql.NewRepository(db, "v2")
	var rRepo = redis.NewRepository(r, "v2", sRepo)
	var s = odin.NewService(rRepo)

	fmt.Println(s.AddPermissionGroup(1, "pg1", "权限组1", odin.Enable))
	fmt.Println(s.AddPermissionGroup(1, "pg2", "权限组1", odin.Enable))
}
