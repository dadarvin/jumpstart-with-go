package config

// Config type of config file
type Config struct {
	AppName  string              `yaml:"appname"`
	HttpPort string              `yaml:"httpport"`
	Database map[string]DBConfig `yaml:"database"`
}

type DBConfig struct {
	MasterDSN string `yaml:"master""`
	SlaveDSN  string `yaml:"slave"`
}
