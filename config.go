package venom

type Mode int

const (
	ProductionMode Mode = iota
	DevelopmentMode
)

type Config struct {
	Address          string
	Port             string
	Mode             Mode
	SuccessFormat    func(obj interface{}) interface{}
	FailFormat       func(errCode interface{}, errMessage string, obj interface{}) interface{}
	Plugins          []IPlugin
	Middlewares      map[string]IMiddleware
	MiddlewarePrefix map[string]string
	Routers          []Router
}

var config = new(Config)

func SetConfig(c *Config) {
	config = c
}

func GetConfig() *Config {
	return config
}
