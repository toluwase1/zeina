package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"zeina/models"
	"zeina/server/response"
)

func (s *Server) LockBalance() gin.HandlerFunc {
	return func(c *gin.Context) {
		var locker models.LockFunds
		if err := decode(c, &locker); err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, err)
			return
		}
		errr := s.WalletService.LockBalance(c, locker)
		if errr != nil {
			response.JSON(c, "", http.StatusInternalServerError, nil, nil)
			return
		}
		message := fmt.Sprintf("Amount of %v locked", locker.Amount)
		response.JSON(c, message, http.StatusCreated, nil, nil)
	}
}
