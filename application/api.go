package application

import (
	"Go_Thingy/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

var DB *gorm.DB
var Error error

func Api() {
	err := godotenv.Load("env.env")
	if err != nil {
		log.Fatalf("Error loading env.env file")
	}
	dsn := os.Getenv("DB_USERNAME") +
		":" +
		os.Getenv("DB_PASSWORD") +
		"@tcp(" +
		os.Getenv("DB_IP") +
		":" +
		os.Getenv("DB_PORT") +
		")/" +
		os.Getenv("DB_NAME") +
		"?parseTime=true"

	DB, Error = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if Error != nil {
		return
	}

	router := gin.Default()
	router.GET("/cars/:license_plate", getCar)
	router.GET("/cars", getCars)
	router.POST("/cars", createCar)
	router.PUT("/cars", updateCar)
	router.DELETE("/cars/:license_plate", deleteCar)

	router.GET("/inspections/:license_plate", getInspections)
	router.POST("/inspections", createInspections)
	router.DELETE("/inspections/:license_plate", deleteInspections)

	router.Run("localhost:3000")
}

func sendError(error string, ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusConflict, models.Response{
		Status:  "fail",
		Message: error,
	})
}

func sendData(message interface{}, ctx *gin.Context) {
	response := map[string]interface{}{
		"status":  "success",
		"message": message,
	}

	ctx.IndentedJSON(http.StatusOK, response)
}
