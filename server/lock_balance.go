package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"zeina/models"
	"zeina/server/response"
)

func (s *Server) LockFunds() gin.HandlerFunc {
	return func(c *gin.Context) {
		var locker models.LockFunds
		if err := decode(c, &locker); err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, err)
			return
		}
		errr := s.WalletService.LockBalance(c, locker)
		if errr != nil {
			response.JSON(c, "", http.StatusInternalServerError, nil, errr)
			return
		}
		message := fmt.Sprintf("Amount of %v locked for %v days \n It will be automatically released after %v days", locker.Amount, locker.Days, locker.Days)
		response.JSON(c, message, http.StatusCreated, nil, nil)
	}
}
