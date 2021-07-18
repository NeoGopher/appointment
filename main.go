package main

import (
	"appointment/domain"
	"appointment/handlers"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	domain.Repo.InitializeDB()
	defer domain.Repo.CloseDB()

	r := setupRouter()
	r.Run(port())
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	handlers.RegisterValidator()

	r.POST("/schedule", handlers.SetSchedule)
	r.POST("/book", handlers.BookAppointment)
	r.POST("/list", handlers.ListAppointments)
	r.POST("/cancel", handlers.CancelAppointment)
	r.POST("/signup", handlers.Signup)

	return r
}

// port gets the PORT Number to run the service on, from the environment
// defaults to 8080.
func port() string {
	port := os.Getenv("PORT")

	if len(port) == 0 {
		port = "8080"
	}

	return ":" + port
}
