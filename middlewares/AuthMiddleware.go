package middlewares

import (
	"log"
	"net/http"

	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/clients"
	"github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/utils"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authClient clients.AuthClientInterface
}

func NewAuthMiddleware(authClient clients.AuthClientInterface) *AuthMiddleware {
	return &AuthMiddleware{
		authClient: authClient,
	}
}

func (auth *AuthMiddleware) ValidateToken(c *gin.Context) {
	log.Println("Validating token")
	authToken := c.GetHeader("Authorization")
	if authToken == "" {
		log.Println("Token not found")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token no encontrado"})
		return
	}
	user, err := auth.authClient.GetUserInfo(authToken)
	if err != nil {
		log.Println("User not authenticated")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}
	log.Println("Setting user info in context")
	utils.SetUserInContext(c, user)

	c.Next()
}
