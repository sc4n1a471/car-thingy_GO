package application

import (
	"Go_Thingy/models"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var Error error

func Api() {
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
	router.GET("/cars/:license-plate", getCar)
	router.GET("/cars", getCars)
	router.GET("/cars-all-data", getCarsAllData)
	router.POST("/cars", createCar)
	router.PUT("/cars", updateCar)
	router.DELETE("/cars/:license-plate", deleteCar)

	router.POST("/license-plate", createLicensePlate)
	router.PUT("/license_plate/:license-plate", updateLicensePLate)

	router.GET("/inspections/:license-plate", getInspections)
	router.POST("/inspections", createInspections)
	router.DELETE("/inspections/:license-plate", deleteInspections)

	router.GET("/coordinates", getCoordinates)

	//router.Run("localhost:3000")
	http.ListenAndServe(":3000", router)
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
