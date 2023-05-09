package server

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"zeina/server/response"
)

func (s *Server) GetBalance() gin.HandlerFunc {
	return func(c *gin.Context) {
		accountNumber := c.Param("account")
		log.Println(accountNumber)
		accountBalance, err := s.WalletService.GetBalance("8163608141")
		if err != nil {
			response.JSON(c, "error getting balance", http.StatusNotFound, nil, nil)
			return
		}
		response.JSON(c, "your balance", http.StatusOK, accountBalance, nil)
	}
}
