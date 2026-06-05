package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// 日志中间件 Logger 自定义日志中间件：记录请求耗时，以及状态码 并且还有是谁操作的
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		//记录开始时间
		startTime := time.Now()
		//处理请求
		c.Next()

		//记录结束时间
		endTime := time.Now()
		latency := endTime.Sub(startTime)

		//打印日志
		log.Printf("| %13v | %s | %s | %s |", latency, c.ClientIP(), c.Request.Method, c.Request.URL.Path)
	}
}

// Recovery 自定义全局异常恢复中间件 (防止程序 panic 导致崩溃)
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic Recovered: %v", err)
				c.JSON(500, gin.H{
					"code": 500,
					"msg":  "服务器内部错误",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
