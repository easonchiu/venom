package venom

import (
  "context"
  "github.com/qiniu/qmgo"
)

type IQmgoClient interface {
  GetClient(key ...string) *QmgoClient
  GetDefaultClient() *QmgoClient
  GetDB(name string) *qmgo.Database
  C(name string) *qmgo.Collection
  CloseAll()
}

type QmgoClient struct {
  *qmgo.QmgoClient
}

const DefaultQmgoClientKey = "default"

var qmgoclients = make(map[string]*QmgoClient)

func initQmgoClient(key string, config QmgoConfig) *QmgoClient {
  client := new(QmgoClient)

  if key == "" || config.URI == "" || config.Database == "" || config.Disabled {
    return client
  }

  cli, err := qmgo.Open(context.Background(), &qmgo.Config{
    Uri:         config.URI,
    Database:    config.Database,
    MinPoolSize: &config.MinPoolSize,
    MaxPoolSize: &config.MaxPoolSize,
  })

  if err != nil {
    panic(err)
  }

  client = &QmgoClient{
    cli,
  }

  qmgoclients[key] = client

  return client
}

func (m *QmgoClient) GetClient(key ...string) *QmgoClient {
  k := DefaultQmgoClientKey

  if key != nil && len(key) > 0 {
    k = key[0]
  }

  if qmgoclients == nil {
    return nil
  }

  return qmgoclients[k]
}

func (m *QmgoClient) GetDefaultClient() *QmgoClient {
  return m.GetClient(DefaultQmgoClientKey)
}

func (m *QmgoClient) GetDB(name string) *qmgo.Database {
  return m.Client.Database(name)
}

func (m *QmgoClient) C(name string) *qmgo.Collection {
  return m.QmgoClient.Database.Collection(name)
}

func (m *QmgoClient) CloseAll() {
  for k, c := range qmgoclients {
    if c != nil {
      _ = c.Close(context.Background())
      qmgoclients[k] = nil
    }
  }
}

var (
  _ IQmgoClient = &QmgoClient{}
)
