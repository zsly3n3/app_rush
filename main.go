package main

import (
	"app/conf"
	"app/event"
	"app/routes"
	"app/tools"
	"net/http"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/gin-gonic/gin"
)

var eventHandler *event.EventHandler

//跨域
func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		// c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		// c.Writer.Header().Add("Access-Control-Allow-Headers", "appversion,apptoken")
		// c.Next()
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Appversion, Apptoken, IsPc,Platform")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		//处理请求
		c.Next()
	}
}

func main() {
	tools.CreateIdWorker()
	eventHandler = event.CreateEventHandler()
	r := gin.Default()
	var mode string
	switch conf.Common.Mode {
	case conf.Debug:
		mode = gin.DebugMode
	case conf.Test:
		mode = gin.TestMode
	case conf.Release:
		mode = gin.ReleaseMode
	}
	gin.SetMode(mode)
	r.Use(cors())
	routes.Register(r, eventHandler)
	routes.WebRegister(r, eventHandler)

	server := &http.Server{Addr: conf.Server.HttpServer, Handler: r}
	gracehttp.Serve(server)

	// r.Run(conf.Server.HttpServer) //listen and serve on 0.0.0.0:8080 asd
}
