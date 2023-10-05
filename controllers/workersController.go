package controllers

import (
	"golang-uslugi-server/m/initializers"
	"golang-uslugi-server/m/models"
	"golang-uslugi-server/m/utilities"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateWorker(c *gin.Context) {

	// Handle preflight OPTION request to allow cors
	utilities.HandleOptions(c)

	// Authorized user
	user, _ := c.Get("user")
	// Check admin access
	if user.(models.User).AccessType != "admin" {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusNotAcceptable, gin.H{
			"error": "access denied.",
		})

		return
	}

	var body struct {
		Name           string
		Specialization string
	}

	if c.Bind(&body) != nil {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body.",
		})

		return
	}

	// Create the worker
	worker := models.Worker{Name: body.Name, Specialization: body.Specialization, CurrentWork: `{"Requests":[]}`}
	result := initializers.DB.Create(&worker)

	if result.Error != nil {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create worker.",
		})

		return
	}

	// Respond
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, gin.H{})

}

func GetWorkersAll(c *gin.Context) {

	// Handle preflight OPTION request to allow cors
	utilities.HandleOptions(c)

	// Authorized user
	user, _ := c.Get("user")
	// Check admin access
	if user.(models.User).AccessType != "admin" {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusNotAcceptable, gin.H{
			"error": "access denied.",
		})

		return
	}

	// Getting all workers
	var workers []models.Worker
	result := initializers.DB.Find(&workers)

	if result.Error != nil {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error while getting workers.",
		})

		return
	}

	// sending it as resp
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, gin.H{
		"workers": workers,
		"code":    0,
	})

}
