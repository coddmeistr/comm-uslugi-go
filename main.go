package main

import (
	"golang-uslugi-server/m/controllers"
	"golang-uslugi-server/m/initializers"
	"golang-uslugi-server/m/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()

	r.POST("/signup", controllers.SignUp)
	r.OPTIONS("/signup", controllers.SignUp)
	r.DELETE("/logout", middleware.RequireAuth, controllers.LogOut)
	r.OPTIONS("/logout", middleware.RequireAuth, controllers.LogOut)
	r.POST("/login", controllers.Login)
	r.OPTIONS("/login", controllers.Login)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)
	r.OPTIONS("/validate", middleware.RequireAuth, controllers.Validate)
	r.POST("/newrequest", middleware.RequireAuth, controllers.CreateRequest)
	r.OPTIONS("/newrequest", middleware.RequireAuth, controllers.CreateRequest)
	r.PUT("/rejectrequest", middleware.RequireAuth, controllers.RejectRequest)
	r.OPTIONS("/rejectrequest", middleware.RequireAuth, controllers.RejectRequest)
	r.PUT("/approverequest", middleware.RequireAuth, controllers.ApproveRequest)
	r.OPTIONS("/approverequest", middleware.RequireAuth, controllers.ApproveRequest)
	r.PUT("/donerequest", middleware.RequireAuth, controllers.DoneRequest)
	r.OPTIONS("/donerequest", middleware.RequireAuth, controllers.DoneRequest)
	r.GET("/getallrequests", middleware.RequireAuth, controllers.GetRequestsAll)
	r.OPTIONS("/getallrequests", middleware.RequireAuth, controllers.GetRequestsAll)
	r.GET("/getalluserrequests", middleware.RequireAuth, controllers.GetRequestsOfUser)
	r.OPTIONS("/getalluserrequests", middleware.RequireAuth, controllers.GetRequestsOfUser)
	r.POST("/newworker", middleware.RequireAuth, controllers.CreateWorker)
	r.OPTIONS("/newworker", middleware.RequireAuth, controllers.CreateWorker)
	r.GET("/getallworkers", middleware.RequireAuth, controllers.GetWorkersAll)
	r.OPTIONS("/getallworkers", middleware.RequireAuth, controllers.GetWorkersAll)

	r.Run()
}
