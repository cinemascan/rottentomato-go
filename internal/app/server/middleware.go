package server

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func mwLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		if query != "" {
			path = path + "?" + query
		}

		c.Next()

		sc := c.Writer.Status()

		mwLogger := zap.L().With(
			zap.String("latency", time.Since(startTime).String()),
			zap.String("method", c.Request.Method),
			zap.Int("code", sc),
			zap.String("path", path),
		)

		switch {
		case sc >= 400:
			if c.Errors.Last() != nil {
				mwLogger.Error("request failed", zap.Error(c.Errors.Last()))
			}
		default:
			mwLogger.Info("request succeeded")
		}
	}
}
