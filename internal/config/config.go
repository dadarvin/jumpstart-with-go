package config

import (
	"gopkg.in/ini.v1"
	"log"
)

var (
	appConfig *Config
)

// Config config struct
type Config struct {
	Environment string
	AppName     string
	HttpPort    string

	AuthConfig *AuthConfig
	DBMaster   *sqlDB
	DBSlave    *sqlDB
	Redis      *RedisConfig
}

// Config config struct for .ini mapping
type configIni struct {
	// General config
	AppName     string `ini:"appname"`
	Environment string `ini:"environment"`
	HttpPort    string `ini:"httpport"`

	// Authentication config
	JWTSecret string `ini:"jwtsecret"`

	// Database config
	DBMasterUser string `ini:"dbmaster_user"`
	DBMasterPass string `ini:"dbmaster_pass"`
	DBMasterHost string `ini:"dbmaster_host"`
	DBMasterPort string `ini:"dbmaster_port"`
	DBMasterName string `ini:"dbmaster_name"`
	DBSlaveUser  string `ini:"dbslave_user"`
	DBSlavePass  string `ini:"dbslave_pass"`
	DBSlaveHost  string `ini:"dbslave_host"`
	DBSlavePort  string `ini:"dbslave_port"`
	DBSlaveName  string `ini:"d	bslave_name"`

	// Redis config
	RedisAddr            string `ini:"redismaster_addr"`
	RedisPass            string `ini:"redusmaster_pass"`
	RedisMaxIdle         int    `ini:"redis_maxidle"`
	RedisMaxActive       int    `ini:"redis_maxactive"`
	RedisIdleTimeout     int64  `ini:"redis_idletimeoutsec"`
	RedisMaxConnLifetime int64  `ini:"redis_maxconnlifetimesec"`
	RedisWait            bool   `ini:"redis_wait"`
}

// Init init the config mapping
func Init() {
	c := &configIni{}

	cIni, err := ini.Load("./conf/app.ini")
	if err != nil {
		log.Fatalf("[Init] failed to read config, %+v\n", err)
	}
	err = cIni.MapTo(c)
	if err != nil {
		log.Fatalf("[Init] failed to map config, %+v\n", err)
	}

	// Init the config
	appConfig = &Config{}
	appConfig.initCommonConfig(c)
	appConfig.initAuthConfig(c)
	appConfig.initSqlDBConfig(c)
	appConfig.initRedisConfig(c)
}

// Get getting config data
func Get() *Config {
	return appConfig
}
