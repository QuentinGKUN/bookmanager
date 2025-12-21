package middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

// CORS 跨域中间件
func CORS() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if string(c.Method()) == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next(ctx)
	}
}



