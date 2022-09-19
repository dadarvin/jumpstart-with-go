package config

type RedisConfig struct {
	Address         string
	Password        string
	MaxIdle         int
	MaxActive       int
	IdleTimeout     int64
	MaxConnLifetime int64
	Wait            bool
}

func (c *Config) initRedisConfig(cfg *configIni) {
	appConfig.Redis = &RedisConfig{
		Address:         cfg.RedisAddr,
		Password:        cfg.RedisPass,
		MaxIdle:         cfg.RedisMaxIdle,
		MaxActive:       cfg.RedisMaxActive,
		IdleTimeout:     cfg.RedisIdleTimeout,
		MaxConnLifetime: cfg.RedisMaxConnLifetime,
		Wait:            cfg.RedisWait,
	}
}
