package server

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) defineRoutes(router *gin.Engine) {
	router.GET("/ping", s.Ping())
	router.POST("/webhook/live", s.WebhookLive())
	router.POST("/webhook/test", s.WebhookTest())

	apirouter := router.Group("/api/v1")
	apirouter.POST("/auth/signup", s.HandleSignup())
	apirouter.POST("/auth/login", s.handleLogin())
	apirouter.POST("/deposit", s.HandleDeposit())
	apirouter.POST("/withdraw", s.HandleWithdrawal())
	apirouter.POST("/lock-funds", s.LockFunds())
	apirouter.GET("/balance", s.GetBalance())

	//authorized := apirouter.Group("/")
	//authorized.Use(s.Authorize())
	//authorized.GET("/logout", s.handleLogout())

}

func (s *Server) setupRouter() *gin.Engine {
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "test" {
		r := gin.New()
		s.defineRoutes(r)
		return r
	}

	r := gin.New()

	// LoggerWithFormatter middleware will write the logs to gin.DefaultWriter
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	r.Use(gin.Recovery())
	// setup cors
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	s.defineRoutes(r)

	return r
}
