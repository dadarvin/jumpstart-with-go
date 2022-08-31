package config

import "fmt"

const (
	dbStringConnection = "%s:%s@tcp(%s:%s)/%s"
)

type sqlDB struct {
	ConnectionString string
}

func (c *Config) initSqlDBConfig(cfg *configIni) {
	appConfig.DBMaster = &sqlDB{
		ConnectionString: fmt.Sprintf(dbStringConnection,
			cfg.DBMasterUser,
			cfg.DBMasterPass,
			cfg.DBMasterHost,
			cfg.DBMasterPort,
			cfg.DBMasterName,
		),
	}

	appConfig.DBSlave = &sqlDB{
		ConnectionString: fmt.Sprintf(dbStringConnection,
			cfg.DBSlaveUser,
			cfg.DBSlavePass,
			cfg.DBSlaveHost,
			cfg.DBSlavePort,
			cfg.DBSlaveName,
		),
	}
}
