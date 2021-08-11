package venom

type Mode int

const (
  ModeProduction Mode = iota
  ModeDevelopment
)

type Config struct {
  Address       string
  Port          string
  Mode          Mode
  SuccessFormat func(code int, obj interface{}) interface{}
  ErrorFormat   func(code int, errCode string, errMessage string, obj... interface{}) interface{}
  ErrorCodes    map[string]string
  Apollo        ApolloConfig
  Redis         RedisConfig
  Mongodb       MongodbConfig
  Logger        LoggerConfig
}

type ApolloConfig struct {
  IP        string
  ID        string
  Cluster   string
  Namespace string
  Disabled  bool
}

type RedisConfig struct {
  Address  string
  Port     string
  Password string
  DB       int
  Disabled bool
}

type MongodbConfig struct {
  Address  string
  Port     string
  Username string
  Password string
  DB       string
  Disabled bool
}

type LoggerConfig struct {
  Disabled bool
}
