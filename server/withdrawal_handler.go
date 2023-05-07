package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"zeina/models"
	"zeina/server/response"
)

func (s *Server) HandleWithdrawal() gin.HandlerFunc {
	return func(c *gin.Context) {
		var withdrawal models.TransactionRequest
		if err := decode(c, &withdrawal); err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, err)
			return
		}
		errr := s.WalletService.Withdrawal(c, withdrawal)
		if errr != nil {
			response.JSON(c, "", http.StatusInternalServerError, nil, errr)
			return
		}
		message := fmt.Sprintf("Withdrawal of %v successful", withdrawal.Amount)
		response.JSON(c, message, http.StatusCreated, nil, nil)
	}
}
