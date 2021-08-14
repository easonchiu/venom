package venom

import (
  "context"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
  "time"
)

type mongodbClients map[string]*mongo.Database

var mgclients mongodbClients = make(map[string]*mongo.Database)

func initMongoClient(config MongoConfig) *mongo.Database {
  if config.URI == "" || config.Disabled {
    return nil
  }

  if client := GetMongoClient(); client != nil {
    return client
  }

  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()

  opts := options.Client()
  opts.ApplyURI(config.URI)
  client, err := mongo.Connect(ctx, opts)

  if err != nil {
    panic(err)
  }

  db := client.Database(opts.Auth.AuthSource)

  mgclients["default"] = db

  return db
}

func GetMongoClient() *mongo.Database {
  if mgclients != nil && mgclients["default"] != nil {
    return mgclients["default"]
  }

  return nil
}
