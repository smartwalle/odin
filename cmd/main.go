package main

import (
	"fmt"
	"github.com/smartwalle/dbr"
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/odin"
	"github.com/smartwalle/odin/service/repository/mysql"
	"github.com/smartwalle/odin/service/repository/postgresql"
	"github.com/smartwalle/odin/service/repository/redis"
	//_ "github.com/go-sql-driver/mysql"
	//_ "github.com/lib/pq"
)

// 如需要测试，在 mod 中添加以下依赖
// github.com/go-sql-driver/mysql v1.4.1
// github.com/lib/pq v1.3.0
//
// import 中添加以下包导入
// _ "github.com/go-sql-driver/mysql"
// _ "github.com/lib/pq"

func main() {
	dbs.SetLogger(nil)

	//testMySQL()
	testPostgreSQL()
}

func testMySQL() {
	var db, _ = dbs.NewSQL("mysql", "root:yangfeng@tcp(127.0.0.1:3306)/test?parseTime=true", 10, 1)
	var repo = mysql.NewRepository(db, "xx")
	testRepo(repo)
}

func testPostgreSQL() {
	var db, _ = dbs.NewSQL("postgres", "host=localhost port=5432 user=postgres password=yangfeng dbname=postgres sslmode=disable", 10, 1)
	var repo = postgresql.NewRepository(db, "xx")
	testRepo(repo)
}

func testRepo(repo odin.Repository) {
	var r = dbr.NewRedis("127.0.0.1:6379", 30, 10, dbr.DialDatabase(1))
	var rRepo = redis.NewRepository(r, "xx", repo)

	var s = odin.NewService(rRepo)

	fmt.Println("Init", s.Init())

	// 添加权限组
	fmt.Println(s.AddPermissionGroup(1, "yf", "研发组", odin.Enable))
	fmt.Println(s.AddPermissionGroup(1, "yx", "营销组", odin.Enable))
	fmt.Println(s.AddPermissionGroup(1, "cw", "财务组", odin.Enable))

	// 添加权限信息
	fmt.Println(s.AddPermissionWithGroup(1, "yf", "yf1", "研发权限1", "", odin.Enable))
	fmt.Println(s.AddPermissionWithGroup(1, "yf", "yf2", "研发权限2", "", odin.Enable))
	fmt.Println(s.AddPermissionWithGroup(1, "yf", "yf3", "研发权限3", "", odin.Enable))
	fmt.Println(s.AddPermissionWithGroup(1, "yf", "yf4", "研发权限4", "", odin.Enable))
	fmt.Println(s.AddPermissionWithGroup(1, "yf", "yf5", "研发权限5", "", odin.Enable))
	fmt.Println(s.AddPermissionWithGroup(1, "yf", "yf6", "研发权限6", "", odin.Enable))
	fmt.Println(s.AddPermissionWithGroup(1, "yf", "yf7", "研发权限7", "", odin.Enable))
	fmt.Println(s.AddPermissionWithGroup(1, "yf", "yf8", "研发权限8", "", odin.Enable))
	fmt.Println(s.AddPermissionWithGroup(1, "yf", "yf9", "研发权限9", "", odin.Enable))

	fmt.Println(s.AddPermissionWithGroup(1, "yx", "yx1", "营销权限1", "", odin.Enable))
	fmt.Println(s.AddPermissionWithGroup(1, "yx", "yx2", "营销权限2", "", odin.Enable))
	fmt.Println(s.AddPermissionWithGroup(1, "yx", "yx3", "营销权限3", "", odin.Enable))
	fmt.Println(s.AddPermissionWithGroup(1, "yx", "yx4", "营销权限4", "", odin.Enable))
	fmt.Println(s.AddPermissionWithGroup(1, "yx", "yx5", "营销权限5", "", odin.Enable))
	fmt.Println(s.AddPermissionWithGroup(1, "yx", "yx6", "营销权限6", "", odin.Enable))
	fmt.Println(s.AddPermissionWithGroup(1, "yx", "yx7", "营销权限7", "", odin.Enable))
	fmt.Println(s.AddPermissionWithGroup(1, "yx", "yx8", "营销权限8", "", odin.Enable))
	fmt.Println(s.AddPermissionWithGroup(1, "yx", "yx9", "营销权限9", "", odin.Enable))

	fmt.Println(s.AddPermissionWithGroup(1, "cw", "cw1", "财务权限1", "", odin.Enable))
	fmt.Println(s.AddPermissionWithGroup(1, "cw", "cw2", "财务权限2", "", odin.Enable))
	fmt.Println(s.AddPermissionWithGroup(1, "cw", "cw3", "财务权限3", "", odin.Enable))
	fmt.Println(s.AddPermissionWithGroup(1, "cw", "cw4", "财务权限4", "", odin.Enable))
	fmt.Println(s.AddPermissionWithGroup(1, "cw", "cw5", "财务权限5", "", odin.Enable))
	fmt.Println(s.AddPermissionWithGroup(1, "cw", "cw6", "财务权限6", "", odin.Enable))
	fmt.Println(s.AddPermissionWithGroup(1, "cw", "cw7", "财务权限7", "", odin.Enable))
	fmt.Println(s.AddPermissionWithGroup(1, "cw", "cw8", "财务权限8", "", odin.Enable))
	fmt.Println(s.AddPermissionWithGroup(1, "cw", "cw9", "财务权限9", "", odin.Enable))

	fmt.Println(s.AddRoleWithParent(1, "", "admin", "管理员", "", odin.Enable))
	fmt.Println(s.GrantPermission(1, "admin", "yf1", "yf2", "yf3", "yf4", "yf5", "yf6", "yf7", "yf8", "yf9"))

	fmt.Println(s.GrantPermission(1, "admin", "yx1", "yx2", "yx3", "yx4", "yx5", "yx6", "yx7", "yx8", "yx9"))
	fmt.Println(s.GrantPermission(1, "admin", "cw1", "cw2", "cw3", "cw4", "cw5", "cw6", "cw7", "cw8", "cw9"))

	fmt.Println(s.AddRoleWithParent(1, "admin", "yfzj", "研发总监", "", odin.Enable))
	fmt.Println(s.AddRoleWithParent(1, "yfzj", "yfjl", "研发经理", "", odin.Enable))
	fmt.Println(s.AddRoleWithParent(1, "yfjl", "yfzg", "研发主管", "", odin.Enable))
	fmt.Println(s.AddRoleWithParent(1, "yfzg", "yfry", "研发人员", "", odin.Enable))
	fmt.Println(s.GrantPermission(1, "yfzj", "yf1", "yf2", "yf3", "yf4", "yf5", "yf6", "yf7", "yf8", "yf9"))
	fmt.Println(s.GrantPermission(1, "yfjl", "yf1", "yf2", "yf3", "yf4", "yf5", "yf6", "yf7", "yf8", "yf9"))
	fmt.Println(s.GrantPermission(1, "yfzg", "yf1", "yf2", "yf3", "yf4", "yf5", "yf6", "yf7", "yf8", "yf9"))
	fmt.Println(s.GrantPermission(1, "yfry", "yf1", "yf2", "yf3", "yf4", "yf5", "yf6", "yf7", "yf8", "yf9"))

	fmt.Println(s.AddRoleWithParent(1, "admin", "yxzj", "营销总监", "", odin.Enable))
	fmt.Println(s.AddRoleWithParent(1, "yxzj", "yxjl", "营销经理", "", odin.Enable))
	fmt.Println(s.AddRoleWithParent(1, "yxjl", "yxzg", "营销主管", "", odin.Enable))
	fmt.Println(s.AddRoleWithParent(1, "yxzg", "yxry", "营销人员", "", odin.Enable))
	fmt.Println(s.GrantPermission(1, "yxzj", "yx1", "yx2", "yx3", "yx4", "yx5", "yx6", "yx7", "yx8", "yx9"))
	//
	fmt.Println(s.AddRoleWithParent(1, "admin", "cwzj", "财务总监", "", odin.Enable))
	fmt.Println(s.AddRoleWithParent(1, "cwzj", "cwjl", "财务经理", "", odin.Enable))
	fmt.Println(s.AddRoleWithParent(1, "cwjl", "cwzg", "财务主管", "", odin.Enable))
	fmt.Println(s.AddRoleWithParent(1, "cwzg", "cwry", "财务人员", "", odin.Enable))
	fmt.Println(s.GrantPermission(1, "cwzj", "cw1", "cw2", "cw3", "cw4", "cw5", "cw6", "cw7", "cw8", "cw9"))

	s.AddRoleMutex(1, "yfzj", "yxzj", "cwzj")
	s.AddPreRole(1, "yfzj", "yfjl")
	s.AddPreRole(1, "yfjl", "yfzg")
	s.AddPreRole(1, "yfzg", "yfry")

	fmt.Println("3", s.GrantRole(1, "t1", "yfry"))
	fmt.Println("3", s.GrantRole(1, "t1", "yfzg"))
	fmt.Println("2", s.GrantRole(1, "t1", "yfjl"))
	fmt.Println("1", s.GrantRole(1, "t1", "yfzj"))
	fmt.Println("4", s.GrantRole(1, "t2", "yfzg"))

	fmt.Println("========= GetRolesWithTarget - t1")
	roles, _ := s.GetRolesWithTarget(1, "t1")
	printRoles(0, roles)

	fmt.Println("========= GetRolesWithTarget - t2")
	roles, _ = s.GetRolesWithTarget(1, "t2")
	printRoles(0, roles)

	fmt.Println("========= GetRoles - t2")
	roles, _ = s.GetRoles(1, 0, "", "t2", "")
	printRoles(0, roles)

	fmt.Println("========= GetGrantedRoles - t1")
	roles, _ = s.GetGrantedRoles(1, "t1")
	printRoles(0, roles)

	s.AddRole(1, "test", "测试权限", "xx", odin.Enable)
	fmt.Println(s.GrantPermission(1, "test", "yf3", "yf2", "yf1"))
}

func printRoles(level int, roles []*odin.Role) {
	for _, role := range roles {
		for i := 0; i < role.Depth; i++ {
			fmt.Print("-")
		}

		fmt.Println("Id:", role.Id, "Alias name:", role.AliasName, "Granted:", role.Granted, "Accessible:", role.Accessible)
	}
}
