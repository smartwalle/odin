package odin

import "github.com/smartwalle/dbr"

var rPool *dbr.Pool

func Init(url, password string, dbIndex, maxActive, maxIdle int) {
	rPool = dbr.NewRedis(url, password, dbIndex, maxActive, maxIdle)
}

func getRedisSession() *dbr.Session {
	var s = rPool.GetSession()
	return s
}