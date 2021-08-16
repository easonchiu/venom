package venom

import (
  "context"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
  "time"
)

type MongoClient struct {
  *mongo.Database
}

type mongodbClients map[string]*MongoClient

const DefaultMongoClientKey = "default"

var mgclients mongodbClients = make(map[string]*MongoClient)

func initMongoClient(key string, config MongoConfig) *MongoClient {
  client := new(MongoClient)

  if key == "" || config.URI == "" || config.Disabled {
    return client
  }

  if c := client.GetClient(key); c != nil {
    return c
  }

  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()

  opts := options.Client()
  opts.ApplyURI(config.URI)
  mgo, err := mongo.Connect(ctx, opts)

  if err != nil {
    panic(err)
  }

  client = &MongoClient{
    mgo.Database(opts.Auth.AuthSource),
  }

  mgclients[key] = client

  return client
}

func (m *MongoClient) GetClient(key ...string) *MongoClient {
  k := DefaultMongoClientKey

  if key != nil && len(key) > 0 {
    k = key[0]
  }

  if mgclients == nil {
    return nil
  }

  return mgclients[k]
}

func (m *MongoClient) GetDefaultClient() *MongoClient {
  return m.GetClient(DefaultMongoClientKey)
}
