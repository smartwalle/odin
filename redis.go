package odin

import "github.com/smartwalle/dbr"

var pool *dbr.Pool

func getPool() *dbr.Pool {
	if pool == nil {
		pool = dbr.NewRedis("127.0.0.1:6379", "", 2, 30, 10)
	}
	return pool
}

func getSession() *dbr.Session {
	var s = getPool().GetSession()
	return s
}