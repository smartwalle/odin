package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/smartwalle/odin"
	"os"
	"io/ioutil"
	"strings"
)

func main() {
	var host string
	var port string
	var password string
	var dbIndex int
	var path string
	var merge int
	var in int

	flag.StringVar(&host, "host", "localhost", "Redis 服务器地址.")
	flag.StringVar(&port, "port", "6379", "Redis 端口.")
	flag.StringVar(&password, "password", "", "Redis 链接密码.")
	flag.IntVar(&dbIndex, "db", 2, "Redis 数据库.")
	flag.StringVar(&path, "path", "odin.json", "备份文件保存路径.")
	flag.IntVar(&merge, "merge", 1, "导入数据的时候, 是否与原有数据合并: 1-需要合并, 其它-不需要合并.")
	flag.IntVar(&in, "in", 0, " 1-导入数据, 其它-导出数据.")
	flag.Parse()

	var url = strings.Join([]string{host, port}, ":")
	odin.Init(url, password, dbIndex, 1, 1)

	if in == 1 {
		importOdin(path, merge==1)
	} else {
		exportOdin(path)
	}
}

func exportOdin(path string) {
	var permissionList, _ = odin.GetPermissionList()

	var roleList, _ = odin.GetRoleList()

	var grantList, _ = odin.GetAllGrantRoleList()

	var exportData = &exportData{}
	exportData.PermissionList = permissionList
	exportData.RoleList = roleList
	exportData.GrantList = grantList

	bs, err := json.Marshal(exportData)
	if err != nil {
		fmt.Println(err)
		return
	}

	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = f.Write(bs)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("导出数据成功:", path)
}

func importOdin(path string, merge bool) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	bs, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println(err)
		return
	}

	var exportData *exportData
	err = json.Unmarshal(bs, &exportData)
	if err != nil {
		fmt.Println(err)
		return
	}

	if merge == false {
		if err = odin.RevokeAllGrant(); err != nil {
			fmt.Println("Remove GrantRole", err)
		}
		if err = odin.RemoveAllPermission(); err != nil {
			fmt.Println("Remove Permission", err)
		}
		if err = odin.RemoveAllRole(); err != nil {
			fmt.Println("Remove Role", err)
		}
	}

	if exportData != nil {
		for _, p := range exportData.PermissionList {
			if _, err = odin.UpdatePermission(p.Id, p.Group, p.Name, p.Identifier); err != nil {
				fmt.Println("UpdatePermission", err)
			}
		}
		for _, r := range exportData.RoleList {
			if err = odin.UpdateRole(r.Id, r.Group, r.Name, r.PermissionIdList...); err != nil {
				fmt.Println("UpdateRole",err)
			}
		}
		for _, g := range exportData.GrantList {
			if err = odin.GrantRole(g.DestinationId, g.RoleIdList...); err != nil {
				fmt.Println("GrantRole",err)
			}
		}
	}

	fmt.Println("导入数据成功:", path)
}

type exportData struct {
	PermissionList []*odin.Permission `json:"permission_list"`
	RoleList       []*odin.Role       `json:"role_list"`
	GrantList      []*odin.GrantInfo  `json:"grant_list"`
}
