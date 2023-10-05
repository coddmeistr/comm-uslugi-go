package controllers

import (
	"golang-uslugi-server/m/initializers"
	"golang-uslugi-server/m/models"
	"golang-uslugi-server/m/utilities"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {

	// Handle preflight OPTION request to allow cors
	utilities.HandleOptions(c)

	// Get the login/email/password
	var body struct {
		Login    string
		Password string
		Email    string
	}

	if c.Bind(&body) != nil {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body.",
			"code":  1,
		})

		return
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to hash password.",
			"code":  1,
		})

		return
	}

	// Create the user
	user := models.User{Login: body.Login, Password: string(hash), Email: body.Email, AccessType: "user"}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "Duplicate entry") {
			c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Login already exists.",
				"code":  1,
			})
		} else {
			c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "failed to create user.",
				"code":  1,
			})
		}

		return
	}

	// Respond
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, gin.H{"code": 0})

}

func Login(c *gin.Context) {

	// Handle preflight OPTION request to allow cors
	utilities.HandleOptions(c)

	// Get the login and pass of req body
	// Get the login/email/password
	var body struct {
		Login    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body.",
			"code":  1,
		})

		return
	}

	// Check if user exists
	var user models.User
	initializers.DB.First(&user, "login = ?", body.Login)

	if user.ID == 0 {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid login.",
			"code":  1,
		})

		return
	}

	// Compare sent pass with hash pass
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid password.",
			"code":  1,
		})

		return
	}

	// Generate jwt
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token.",
			"code":  1,
		})

		return
	}

	// Send it back
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("LoginToken", tokenString, 3600*24, "", "", false, true)

	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, gin.H{"code": 0})
}

func LogOut(c *gin.Context) {

	// Handle preflight OPTION request to allow cors
	utilities.HandleOptions(c)

	// Send it back
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("LoginToken", "0", 1, "", "", false, true)

	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, gin.H{"code": 0})

}

func Validate(c *gin.Context) {

	// Handle preflight OPTION request to allow cors
	utilities.HandleOptions(c)

	user, _ := c.Get("user")
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, gin.H{
		"user": user,
		"code": 0,
	})
}
