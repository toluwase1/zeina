package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zeina/server/response"
)

func (s *Server) Ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		response.JSON(c, "This is Zeina", http.StatusOK, nil, nil)
	}
}
