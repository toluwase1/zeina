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

func (s *Server) WebookTest() gin.HandlerFunc {
	return func(c *gin.Context) {
		todaysTime := time.Now().String()
		message := "This is Zeina: " + todaysTime
		webhookData := models.Webhook{}
		if err := decode(c, &webhookData); err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, err)
			return
		}
		log.Println(time.Now(), " ", webhookData)
		response.JSON(c, message, http.StatusOK, webhookData, nil)
	}
}
