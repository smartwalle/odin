package loaddata

import (
	"encoding/json"
	"github.com/smartwalle/odin"
	"io/ioutil"
	"os"
)

type Role struct {
	Ctx              int64              `json:"ctx"`
	Name             string             `json:"name"`
	AliasName        string             `json:"alias_name"`
	Description      string             `json:"description"`
	Targets          []string           `json:"targets"`
	PermissionGroups []*PermissionGroup `json:"permission_groups"`
}

type PermissionGroup struct {
	Name        string        `json:"name"`
	AliasName   string        `json:"alias_name"`
	Permissions []*Permission `json:"permissions"`
}

type Permission struct {
	Name        string `json:"name"`
	AliasName   string `json:"alias_name"`
	Description string `json:"description"`
}

func LoadData(repo odin.Repository, configFile string) error {
	var file, err = os.Open(configFile)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	var role *Role
	if err = json.Unmarshal(data, &role); err != nil {
		return err
	}

	var service = odin.NewService(repo)

	if role != nil {
		service.AddRole(role.Ctx, role.Name, role.AliasName, role.Description, odin.Enable)
		for _, target := range role.Targets {
			service.GrantRole(role.Ctx, target, role.Name)
		}

		for _, group := range role.PermissionGroups {
			service.AddPermissionGroup(role.Ctx, group.Name, group.AliasName, odin.Enable)
			for _, permission := range group.Permissions {
				service.AddPermissionWithGroup(role.Ctx, group.Name, permission.Name, permission.AliasName, permission.Description, odin.Enable)
				service.GrantPermission(role.Ctx, role.Name, permission.Name)
			}
		}
	}

	return nil
}
