package venom

import (
  "context"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
  "path"
  "time"
)

type MongoClient struct {
  *mongo.Database
}

const DefaultMongoClientKey = "default"

var mgclients = make(map[string]*MongoClient)

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
  if config.MinPoolSize > 0 {
    opts.SetMinPoolSize(config.MinPoolSize)
  }
  if config.MaxPoolSize > 0 {
    opts.SetMaxPoolSize(config.MinPoolSize)
  }

  mgo, err := mongo.Connect(ctx, opts)
  if err != nil {
    panic(err)
  }

  databaseName := path.Base(config.URI)
  if opts.Auth != nil && opts.Auth.AuthSource != "" {
    databaseName = opts.Auth.AuthSource
  }
  if config.Database != "" {
    databaseName = config.Database
  }

  client = &MongoClient{
    mgo.Database(databaseName),
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

func (m *MongoClient) GetDB(name string) *mongo.Database {
  return m.Database.Client().Database(name)
}

func (m *MongoClient) C(name string) *mongo.Collection {
  return m.Collection(name)
}

func (m *MongoClient) CloseAll() {
  for k, c := range mgclients {
    if c != nil {
      _ = c.Client().Disconnect(context.Background())
      mgclients[k] = nil
    }
  }
}
