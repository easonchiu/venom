package venom

import (
  "github.com/gin-gonic/gin"
  "github.com/sirupsen/logrus"
  "gopkg.in/natefinch/lumberjack.v2"
  "math"
  "net/http"
  "os"
  "time"
)

// 日志记录到文件
func LoggerMiddleware(config LoggerConfig) gin.HandlerFunc {
  // 实例化
  logger := logrus.New()
  logger.SetOutput(&lumberjack.Logger{
    Filename:   config.Filename,
    MaxSize:    config.MaxSize, // megabytes
    MaxBackups: config.MaxBackups,
    MaxAge:     config.MaxAge, // days
    Compress:   true,          // disabled by default
  })
  logger.SetLevel(config.Level)

  return loggerHandle(logger)
}

// Logger is the logrus logger handler
func loggerHandle(logger logrus.FieldLogger) gin.HandlerFunc {
  hostname, err := os.Hostname()
  if err != nil {
    hostname = "unknow"
  }

  return func(c *gin.Context) {
    // other handler can change c.Path so:
    p := c.Request.URL.Path
    start := time.Now()
    c.Next()
    stop := time.Since(start)
    latency := int(math.Ceil(float64(stop.Nanoseconds()) / 1000000.0))
    statusCode := c.Writer.Status()
    clientIP := c.ClientIP()
    clientUserAgent := c.Request.UserAgent()
    referer := c.Request.Referer()
    dataLength := c.Writer.Size()
    if dataLength < 0 {
      dataLength = 0
    }

    entry := logger.WithFields(logrus.Fields{
      "hostname":   hostname,
      "statusCode": statusCode,
      "latency":    latency, // time to process
      "clientIP":   clientIP,
      "method":     c.Request.Method,
      "path":       p,
      "referer":    referer,
      "dataLength": dataLength,
      "userAgent":  clientUserAgent,
    })

    if len(c.Errors) > 0 {
      entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
    } else {
      if statusCode >= http.StatusInternalServerError {
        entry.Error()
      } else if statusCode >= http.StatusBadRequest {
        entry.Warn()
      } else {
        entry.Info()
      }
    }
  }
}
