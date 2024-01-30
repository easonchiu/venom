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
	Key         string
	URI         string
	Database    string
	MinPoolSize uint64
	MaxPoolSize uint64
}

const DefaultQmgoClientKey = "default"

// qmgo 的 client list
var qmgoClients = make(map[string]*qmgo.QmgoClient)

// init plugin
func InitQmgoPlugin(config *QmgoConfig) *QmgoPlugin {
	return &QmgoPlugin{config: config}
}

// 启动
func (plugin *QmgoPlugin) OnStart(config *Config) {
	if plugin.config.Key == "" {
		plugin.config.Key = DefaultQmgoClientKey
	}

	// client重名
	if _, exist := qmgoClients[plugin.config.Key]; exist {
		panic(fmt.Errorf("[VENOM] PLUGIN - Qmgo key %v already exists", plugin.config.Key))
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
	qmgoClients[plugin.config.Key] = client

	fmt.Printf("[VENOM] PLUGIN - Qmgo <%v> start ok...\n", plugin.config.Key)
}

// 卸载
func (plugin *QmgoPlugin) OnDestroy(config *Config) {
	if plugin.config.Key != "" && plugin != nil {
		if _, exist := qmgoClients[plugin.config.Key]; exist {
			_ = plugin.client.Close(context.Background())
			delete(qmgoClients, plugin.config.Key)
		}
	}
}

// 获取 qmgo 的 client
func (plugin *QmgoPlugin) GetClient(key ...string) *qmgo.QmgoClient {
	k := DefaultQmgoClientKey

	if key != nil {
		k = key[0]
	}

	if plugin.config != nil && plugin.config.Key == k && plugin.client != nil {
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
