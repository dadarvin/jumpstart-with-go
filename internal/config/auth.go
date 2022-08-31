package config

// AuthConfig represents configurations for this service's authentication.
type AuthConfig struct {
	JWTSecret string
}

func (c *Config) initAuthConfig(cfg *configIni) {
	appConfig.AuthConfig = &AuthConfig{
		JWTSecret: cfg.JWTSecret,
	}
}
