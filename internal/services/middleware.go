package services

import (
	"fmt"
	"log"
	"strings"

	"github.com/MiguelCVA/mhook-backend/pkg"
	"github.com/gin-gonic/gin"
)

func abortWithUnauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(401, gin.H{"message": "Unauthorized"})
}

func AuthMiddleware(c *gin.Context) {
	authorizationHeader := c.Request.Header.Get("Authorization")

	if authorizationHeader == "" || !strings.HasPrefix(authorizationHeader, "Bearer ") {
		abortWithUnauthorized(c)
		return
	}

	token := strings.TrimPrefix(authorizationHeader, "Bearer ")
	payload, err := pkg.DecodeJWT(token)

	if err != nil {
		log.Printf("Erro ao decodificar token: %v", err)
		abortWithUnauthorized(c)
		return
	}

	fmt.Println(payload)

	c.Next()
}
