package config

var GlobalConfig *Config

// Config is application global config
type Config struct {
	Mode string `mapstructure:"mode"` // gin启动模式
	Port string `mapstructure:"port"` // 启动端口
}

// Load load config by path
func Load(configFilePath string) *Config {
	// TODO:
	return &Config{
		Mode: "debug",
		Port: "3000",
	}
}
