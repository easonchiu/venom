/*
 * @Author: zhaozhida zhaozhida@qiniu.com
 * @Date: 2023-07-26 11:11:40
 * @LastEditors: zhaozhida zhaozhida@qiniu.com
 * @LastEditTime: 2023-07-26 14:36:39
 * @Description:
 */
package venom

import (
	"context"
	"fmt"

	"github.com/qiniu/qmgo"
)

type QmgoPlugin struct {
	client *qmgo.QmgoClient
	config *QmgoConfig
}

type QmgoConfig struct {
	Name        string
	URI         string
	Database    string
	MinPoolSize uint64
	MaxPoolSize uint64
}

const DefaultQmgoClientName = "default"

// qmgo 的 client list
var qmgoClients = make(map[string]*qmgo.QmgoClient)

// init plugin
func InitQmgoPlugin(config *QmgoConfig) *QmgoPlugin {
	return &QmgoPlugin{config: config}
}

// 启动
func (plugin *QmgoPlugin) OnStart(config *Config) {
	if plugin.config.Name == "" {
		plugin.config.Name = DefaultQmgoClientName
	}

	// client重名
	if _, exist := qmgoClients[plugin.config.Name]; exist {
		panic(fmt.Errorf("[PLUGIN] Qmgo named %v already exists", plugin.config.Name))
	}

	client, err := qmgo.Open(context.Background(), &qmgo.Config{
		Uri:         plugin.config.URI,
		Database:    plugin.config.Database,
		MinPoolSize: &plugin.config.MinPoolSize,
		MaxPoolSize: &plugin.config.MaxPoolSize,
	})

	if err != nil {
		panic(err)
	}

	plugin.client = client
	qmgoClients[plugin.config.Name] = client

	fmt.Printf("[PLUGIN] Qmgo <%v> start ok...\n", plugin.config.Name)
}

// 卸载
func (plugin *QmgoPlugin) OnDestroy(config *Config) {
	if plugin.config.Name != "" && plugin != nil {
		if _, exist := qmgoClients[plugin.config.Name]; exist {
			_ = plugin.client.Close(context.Background())
			delete(qmgoClients, plugin.config.Name)
		}
	}
}

// 获取 qmgo 的 client
func (plugin *QmgoPlugin) GetClient(name ...string) *qmgo.QmgoClient {
	k := DefaultQmgoClientName

	if name != nil {
		k = name[0]
	}

	if plugin.config != nil && plugin.config.Name == k && plugin.client != nil {
		return plugin.client
	}

	if qmgoClients == nil {
		return nil
	}

	return qmgoClients[k]
}

func GetQmgoClient(client ...string) *qmgo.QmgoClient {
	plugin := new(QmgoPlugin)
	return plugin.GetClient(client...)
}

func GetQmgoDB(client ...string) *qmgo.Database {
	plugin := new(QmgoPlugin)
	return plugin.GetClient(client...).Database
}

// 检验是否实现了plugin interface
var _ IPlugin = (*QmgoPlugin)(nil)
