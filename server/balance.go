package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zeina/server/response"
)

func (s *Server) GetBalance() gin.HandlerFunc {
	return func(c *gin.Context) {
		//accountNumber := c.Param("account")
		accountBalance, err := s.WalletService.GetBalance("8163608141")
		if err != nil {
			response.JSON(c, "Signup successful", http.StatusNotFound, nil, nil)
			return
		}
		response.JSON(c, "Signup successful", http.StatusOK, accountBalance, nil)
	}
}
