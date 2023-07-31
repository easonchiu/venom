/*
 * @Author: zhaozhida zhaozhida@qiniu.com
 * @Date: 2023-07-26 10:32:09
 * @LastEditors: zhaozhida zhaozhida@qiniu.com
 * @LastEditTime: 2023-07-26 14:32:02
 * @Description:
 */
package venom

type Mode int

const (
	ProductionMode Mode = iota
	DevelopmentMode
)

type Config struct {
	Address       string
	Port          string
	Mode          Mode
	SuccessFormat func(obj interface{}) interface{}
	FailFormat    func(errCode interface{}, errMessage string, obj interface{}) interface{}
	Plugins       []IPlugin
	Middlewares   []IMiddleware
}

var config = new(Config)

func SetConfig(c *Config) {
	config = c
}

func GetConfig() *Config {
	return config
}
