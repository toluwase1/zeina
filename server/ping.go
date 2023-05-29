package server

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
	"zeina/models"
	"zeina/server/response"
)

func (s *Server) Ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		response.JSON(c, "This is Zeina", http.StatusOK, nil, nil)
	}
}

func (s *Server) WebhookLive() gin.HandlerFunc {
	return func(c *gin.Context) {
		todaysTime := time.Now().String()
		message := "This is LIVE: " + todaysTime
		webhookData := models.LedgerTransaction{}
		if err := decode(c, &webhookData); err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, err)
			return
		}
		log.Println(time.Now(), " LIVE DATA ", webhookData)
		response.JSON(c, message, http.StatusOK, webhookData, nil)
	}
}

func (s *Server) WebhookTest() gin.HandlerFunc {
	return func(c *gin.Context) {
		todaysTime := time.Now().String()
		message := "This is TEST: " + todaysTime
		webhookData := models.LedgerTransaction{}
		if err := decode(c, &webhookData); err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, err)
			return
		}
		log.Println(time.Now(), " TEST DATA ", webhookData)
		response.JSON(c, message, http.StatusOK, webhookData, nil)
	}
}
