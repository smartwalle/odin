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
	fmt.Println(s.GetGroupList(0, 0, ""))

	fmt.Println(s.GetPermissionList(0, 0, ""))

	fmt.Println(s.GetRoleWithId(1))

	//fmt.Println(s.addGroup(odin.K_GROUP_TYPE_ROLE, "ttt", odin.K_STATUS_ENABLE))
}
