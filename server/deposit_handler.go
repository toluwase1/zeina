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
		//_, user, err := GetValuesFromContext(c)
		//if err != nil {
		//	log.Println("11", err)
		//	err.Respond(c)
		//	return
		//}
		log.Println("11")
		var deposit models.DepositRequest
		if err := decode(c, &deposit); err != nil {
			log.Println("11", err)
			response.JSON(c, "", http.StatusBadRequest, nil, err)
			return
		}
		log.Println("22")
		errr := s.WalletService.Deposit(c, &models.User{}, deposit)
		if errr != nil {
			log.Println("33", errr)
			response.JSON(c, "", http.StatusInternalServerError, nil, errr)
			return
		}
		message := fmt.Sprintf("Deposit of %v successful", deposit.Amount)
		response.JSON(c, message, http.StatusCreated, nil, nil)
	}
}
