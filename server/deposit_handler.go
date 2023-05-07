package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"zeina/models"
	"zeina/server/response"
)

func (s *Server) HandleDeposit() gin.HandlerFunc {
	return func(c *gin.Context) {
		var deposit models.TransactionRequest
		if err := decode(c, &deposit); err != nil {
			log.Println("11", err)
			response.JSON(c, "", http.StatusBadRequest, nil, err)
			return
		}
		errr := s.WalletService.Deposit(c, deposit)
		if errr != nil {
			response.JSON(c, "", http.StatusInternalServerError, nil, errr)
			return
		}
		message := fmt.Sprintf("Deposit of %v successful", deposit.Amount)
		response.JSON(c, message, http.StatusCreated, nil, nil)
	}
}
