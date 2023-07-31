/*
 * @Author: zhaozhida zhaozhida@qiniu.com
 * @Date: 2023-07-26 11:11:40
 * @LastEditors: zhaozhida zhaozhida@qiniu.com
 * @LastEditTime: 2023-07-26 14:36:39
 * @Description:
 */
package venom

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LoggerMiddleware struct {
	client *logrus.Logger
	config *LoggerConfig
}

type LoggerConfig struct {
	Filename   string
	MaxSize    int
	MaxBackups int
	Level      logrus.Level
	MaxAge     int
}

const DefaultLoggerClientName = "default"

// logger 的 client list
var loggerClients = make(map[string]*logrus.Logger)

// init middleware
func InitLoggerMiddleware(config *LoggerConfig) *LoggerMiddleware {
	return &LoggerMiddleware{config: config}
}

// 启动
func (mw *LoggerMiddleware) OnStart(config *Config) {
	mw.client = logrus.New()
	mw.client.SetOutput(&lumberjack.Logger{
		Filename:   mw.config.Filename,
		MaxSize:    mw.config.MaxSize,
		MaxBackups: mw.config.MaxBackups,
		MaxAge:     mw.config.MaxAge,
		Compress:   true,
	})
	mw.client.SetLevel(mw.config.Level)

	loggerClients[DefaultLoggerClientName] = mw.client

	fmt.Printf("[MIDDLEWARE] Logger start ok...\n")
}

// 卸载
func (mw *LoggerMiddleware) OnDestroy(config *Config) {
	delete(loggerClients, DefaultLoggerClientName)
}

// 获取 logger 的 client
func (mw *LoggerMiddleware) GetClient() *logrus.Logger {
	return mw.client
}

// 获取 gin 的中间件
func (mw *LoggerMiddleware) GetGinMiddleware(config *Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := mw.client

		if logger == nil {
			return
		}

		hostname, err := os.Hostname()
		if err != nil {
			hostname = "unknow"
		}

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

func GetLoggerClient() *logrus.Logger {
	return loggerClients[DefaultLoggerClientName]
}

// 检验是否实现了middleware interface
var _ IMiddleware = (*LoggerMiddleware)(nil)
