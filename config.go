package venom

type Mode int

const (
	ProductionMode Mode = iota
	DevelopmentMode
)

type Config struct {
	Address       string
	Port          string
	SuccessFormat func(obj interface{}) interface{}
	FailFormat    func(errCode interface{}, errMessage string, obj interface{}) interface{}
}

// var config = new(Config)

// func SetConfig(c *Config) {
// 	config = c
// }

// func GetConfig() *Config {
// 	return config
// }
