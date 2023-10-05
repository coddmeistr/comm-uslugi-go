package controllers

import (
	"encoding/json"
	"golang-uslugi-server/m/initializers"
	"golang-uslugi-server/m/models"
	"golang-uslugi-server/m/utilities"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func Find(a []uint, x uint) int {
	for i, n := range a {
		if x == n {
			return i
		}
	}
	return len(a)
}

func CreateRequest(c *gin.Context) {

	// Handle preflight OPTION request to allow cors
	utilities.HandleOptions(c)

	// Authorized user
	user, _ := c.Get("user")

	// Get request data
	var body struct {
		Address   string
		WorkType  string
		WorkScale string
		Time      time.Time
	}

	if c.Bind(&body) != nil {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body.",
			"code":  1,
		})

		return
	}

	// Create the request
	const initStatusMessage = "Ваша заявка находится в обработке, ожидайте подтверждения."
	request := models.Request{UserID: user.(models.User).ID, Address: body.Address, WorkType: body.WorkType, WorkScale: body.WorkScale,
		Time: body.Time, Status: "pending", StatusMessage: initStatusMessage, Workers: `{"Workers":[]}`}
	result := initializers.DB.Create(&request)

	if result.Error != nil {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create request.",
			"code":  1,
		})

		return
	}

	// Respond
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, gin.H{"code": 0})
}

func RejectRequest(c *gin.Context) {

	// Handle preflight OPTION request to allow cors
	utilities.HandleOptions(c)

	// Authorized user
	user, _ := c.Get("user")
	// Check admin access
	if user.(models.User).AccessType != "admin" {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusNotAcceptable, gin.H{
			"error": "access denied.",
			"code":  2,
		})

		return
	}

	// Get request data
	var body struct {
		RequestID     uint
		RejectMessage string
	}

	if c.Bind(&body) != nil {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body.",
			"code":  1,
		})

		return
	}

	// Get our request out of db
	var request models.Request
	initializers.DB.First(&request, body.RequestID)

	if request.ID == 0 {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "This request doesn't exist.",
			"code":  0,
		})

		return
	}

	// Check if request is in pending status
	if request.Status != "pending" {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No need to apply actions to this request.",
			"code":  1,
		})

		return
	}

	// Update fields
	result := initializers.DB.Model(&request).Updates(models.Request{Status: "rejected", StatusMessage: body.RejectMessage})
	if result.Error != nil {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to reject request(db error).",
			"code":  1,
		})

		return
	}

	// Send message
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully rejected.",
		"code":    0,
	})

}

func ApproveRequest(c *gin.Context) {

	// Handle preflight OPTION request to allow cors
	utilities.HandleOptions(c)

	// Authorized user
	user, _ := c.Get("user")
	// Check admin access
	if user.(models.User).AccessType != "admin" {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusNotAcceptable, gin.H{
			"error": "access denied.",
			"code":  2,
		})

		return
	}

	// Get request data
	var body struct {
		RequestID      uint
		ApproveMessage string
		Workers        []uint
	}

	if c.Bind(&body) != nil {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body.",
			"code":  1,
		})

		return
	}

	// Get our request out of db
	var request models.Request
	initializers.DB.First(&request, body.RequestID)

	if request.ID == 0 {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "This request doesn't exist.",
			"code":  1,
		})

		return
	}

	// Check if request is in pending status
	if request.Status != "pending" {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No need to apply actions to this request.",
			"code":  1,
		})

		return
	}

	// Check and update all included workers
	for i := 0; i < len(body.Workers); i++ {
		workerID := body.Workers[i]

		// Get current worker out of db
		var worker models.Worker
		initializers.DB.First(&worker, workerID)

		if worker.ID == 0 {
			c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Worker doesn't exist.",
				"code":  1,
			})

			return
		}

		// Update his current works
		var WorkerRequests struct {
			Requests []uint `json:"Requests"`
		}

		// bind json
		err := json.Unmarshal([]byte(worker.CurrentWork), &WorkerRequests)

		if err != nil {
			c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed while unmarshalling json.",
				"code":  1,
			})

			return
		}

		// add work
		WorkerRequests.Requests = append(WorkerRequests.Requests, body.RequestID)

		// Marshal new json
		requestsJSON, err := json.Marshal(WorkerRequests)

		if err != nil {
			c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Falied to marshal JSON.",
				"code":  1,
			})

			return
		}

		// Update workers fields
		result := initializers.DB.Model(&worker).Update("current_work", string(requestsJSON))

		if result.Error != nil {
			c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to approve request(db error).",
				"code":  1,
			})

			return
		}

	}

	// Generate new JSON of Workers
	type Workers struct {
		Workers []uint `json:"Workers"`
	}
	workers := &Workers{Workers: body.Workers}
	workersJSON, err := json.Marshal(workers)

	if err != nil {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Falied to marshal JSON.",
			"code":  1,
		})

		return
	}

	// Update fields of request
	result := initializers.DB.Model(&request).Updates(models.Request{Status: "approved", StatusMessage: body.ApproveMessage,
		Workers: string(workersJSON)})
	if result.Error != nil {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to approve request(db error).",
			"code":  1,
		})

		return
	}

	// Send message
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully approved.",
		"code":    0,
	})

}

func DoneRequest(c *gin.Context) {

	// Handle preflight OPTION request to allow cors
	utilities.HandleOptions(c)

	// Authorized user
	user, _ := c.Get("user")
	// Check admin access
	if user.(models.User).AccessType != "admin" {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusNotAcceptable, gin.H{
			"error": "access denied.",
			"code":  2,
		})

		return
	}

	// Get request data
	var body struct {
		RequestID   uint
		DoneMessage string
	}

	if c.Bind(&body) != nil {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body.",
			"code":  1,
		})

		return
	}

	// Get our request out of db
	var request models.Request
	initializers.DB.First(&request, body.RequestID)

	if request.ID == 0 {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "This request doesn't exist.",
			"code":  1,
		})

		return
	}

	// Check if request is in pending or rejected status
	if request.Status != "approved" {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "This request cannot be done.",
			"code":  1,
		})

		return
	}

	// Get workers array
	var Workers struct {
		Workers []uint
	}
	err := json.Unmarshal([]byte(request.Workers), &Workers)

	if err != nil {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to unmarshal json.",
			"code":  1,
		})

		return
	}

	for i := 0; i < len(Workers.Workers); i++ {
		workerID := Workers.Workers[i]

		// Get current worker out of db
		var worker models.Worker
		initializers.DB.First(&worker, workerID)

		if worker.ID == 0 {
			c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Worker doesn't exist.",
				"code":  1,
			})

			return
		}

		// Update his current works
		var WorkerRequests struct {
			Requests []uint `json:"Requests"`
		}

		// bind json
		err := json.Unmarshal([]byte(worker.CurrentWork), &WorkerRequests)

		if err != nil {
			c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed while unmarshalling json.",
				"code":  1,
			})

			return
		}

		// delete work
		index := Find(WorkerRequests.Requests, body.RequestID)
		if index == len(WorkerRequests.Requests) {
			c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to manage worker's works",
				"code":  1,
			})

			return
		}
		WorkerRequests.Requests = append(WorkerRequests.Requests[:index], WorkerRequests.Requests[index+1:]...)

		// Marshal new json
		requestsJSON, err := json.Marshal(WorkerRequests)

		if err != nil {
			c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Falied to marshal JSON.",
				"code":  1,
			})

			return
		}

		// Update workers fields
		result := initializers.DB.Model(&worker).Update("current_work", string(requestsJSON))

		if result.Error != nil {
			c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to approve request(db error).",
				"code":  1,
			})

			return
		}

	}

	// Update fields of request
	result := initializers.DB.Model(&request).Updates(models.Request{Status: "done", StatusMessage: body.DoneMessage})
	if result.Error != nil {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to done request(db error).",
			"code":  1,
		})

		return
	}

	// Send message
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully made done.",
		"code":    0,
	})

}

func GetRequestsAll(c *gin.Context) {

	// Handle preflight OPTION request to allow cors
	utilities.HandleOptions(c)

	// Authorized user
	user, _ := c.Get("user")
	// Check admin access
	if user.(models.User).AccessType != "admin" {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusNotAcceptable, gin.H{
			"error": "access denied.",
			"code":  2,
		})

		return
	}

	// Getting all requests
	var requests []models.Request
	result := initializers.DB.Find(&requests)

	if result.Error != nil {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error while getting requests.",
			"code":  1,
		})

		return
	}

	// sending it as resp
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, gin.H{
		"requests": requests,
		"code":     0,
	})

}

func GetRequestsOfUser(c *gin.Context) {

	// Handle preflight OPTION request to allow cors
	utilities.HandleOptions(c)

	// Authorized user
	user, _ := c.Get("user")

	// Get he user id from query param
	queryParamID, err := strconv.Atoi(c.Query("userID"))

	if err != nil {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid query param.",
			"code":  1,
		})

		return
	}

	// Get unit id in the var
	queryID := uint(queryParamID)

	// Check the access
	if user.(models.User).AccessType != "admin" && queryID != user.(models.User).ID {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "access denied.",
			"code":  2,
		})

		return
	}

	// Getting all user's requests
	var requests []models.Request
	result := initializers.DB.Where("user_id", queryID).Find(&requests)

	if result.Error != nil {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error while getting requests.",
			"code":  1,
		})

		return
	}

	// sending it as resp
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, gin.H{
		"requests": requests,
		"code":     0,
	})

}
