package odin

import "github.com/smartwalle/dbs"

type manager struct {
	db dbs.DB

	groupTable      string
	permissionTable string
}
