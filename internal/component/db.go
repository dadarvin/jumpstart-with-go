package component

import (
	"database/sql"
	"entry_task/internal/config"
	"log"
)

type DB struct {
	Master *sql.DB
	Slave  *sql.DB
}

func InitDatabase() *DB {
	conf := config.Get()

	if conf.DBMaster == nil {
		log.Fatalf("failed to get DB config")
	}

	var (
		db  = &DB{}
		err error
	)

	db.Master, err = sql.Open("mysql", conf.DBMaster.ConnectionString)
	if err != nil {
		log.Fatalf("failed to open DB master connection. %+v", err)
	}
	err = db.Master.Ping()
	if err != nil {
		log.Fatalf("failed to ping DB master. %+v", err)
	}

	db.Slave, err = sql.Open("mysql", conf.DBSlave.ConnectionString)
	if err != nil {
		log.Fatalf("failed to open DB slave connection. %+v", err)
	}
	err = db.Slave.Ping()
	if err != nil {
		log.Fatalf("failed to ping DB slave. %+v", err)
	}

	return db
}
