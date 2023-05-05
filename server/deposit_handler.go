package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zeina/models"
	"zeina/server/response"
)

func (s *Server) HandleDeposit() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := decode(c, &user); err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, err)
			return
		}
		userResponse, err := s.AuthService.SignupUser(&user)
		if err != nil {
			err.Respond(c)
			return
		}
		response.JSON(c, "Signup successful", http.StatusCreated, userResponse, nil)
	}
}
