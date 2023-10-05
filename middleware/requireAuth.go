package middleware

import (
	"fmt"
	"golang-uslugi-server/m/initializers"
	"golang-uslugi-server/m/models"
	"golang-uslugi-server/m/utilities"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func RequireAuth(c *gin.Context) {

	// Handle preflight OPTION request to allow cors
	utilities.HandleOptions(c)

	// Get the cookie out of req
	tokenString, err := c.Cookie("LoginToken")

	if tokenString == "0" {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if err != nil {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Decode/validate it
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check the exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		// Find the user with token sub
		var user models.User
		initializers.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		// Attach it to the req
		c.Set("user", user)

		// Continue
		c.Next()

	} else {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}
