package redis

import (
	"github.com/smartwalle/dbr"
	"github.com/smartwalle/dbs"
	"github.com/smartwalle/odin"
)

type odinRepository struct {
	odin.Repository
	rPool   dbr.Pool
	tPrefix string
}

func NewRepository(rPool dbr.Pool, tPrefix string, repo odin.Repository) odin.Repository {
	var r = &odinRepository{}
	r.rPool = rPool
	r.tPrefix = tPrefix
	r.Repository = repo
	return r
}

func (this *odinRepository) BeginTx() (dbs.TX, odin.Repository) {
	var nRepo = *this
	var tx dbs.TX
	tx, nRepo.Repository = this.Repository.BeginTx()
	return tx, &nRepo
}

func (this *odinRepository) WithTx(tx dbs.TX) odin.Repository {
	var nRepo = *this
	nRepo.Repository = this.Repository.WithTx(tx)
	return &nRepo
}
