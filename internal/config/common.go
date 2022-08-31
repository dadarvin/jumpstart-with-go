package config

func (c *Config) initCommonConfig(cfg *configIni) {
	c.AppName = cfg.AppName
	c.Environment = cfg.Environment
	c.HttpPort = cfg.HttpPort
}
