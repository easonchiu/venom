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

func initMongoClient(config MongoConfig) *MongoClient {
  if config.URI == "" || config.Disabled {
    return nil
  }

  client := new(MongoClient)

  if c:= client.GetDefaultClient(); c != nil {
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

  mgclients[DefaultMongoClientKey] = client

  return client
}

func (m *MongoClient) GetClient(key string) *MongoClient {
  if mgclients == nil {
    return nil
  }

  return mgclients[key]
}

func (m *MongoClient) GetDefaultClient() *MongoClient {
  return m.GetClient(DefaultMongoClientKey)
}
