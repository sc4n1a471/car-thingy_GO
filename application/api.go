package application

import (
	"Go_Thingy_GO/controllers"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Api() {
	err := controllers.SetupDatabase()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	router := gin.Default()
	router.GET("/cars/:license-plate", controllers.GetCar)
	router.GET("/cars", controllers.GetCars)
	// router.GET("/cars-all-data", controllers.GetCarsAllData)
	router.POST("/cars", controllers.CreateCar)
	router.PUT("/cars", controllers.UpdateCar)
	router.DELETE("/cars/:license-plate", controllers.DeleteCar)

	router.POST("/license-plate", controllers.CreateLicensePlate)
	router.PUT("/license-plate/:license-plate", controllers.UpdateLicensePlate)

	router.GET("/inspections/:license-plate", controllers.GetInspectionsHelper)
	router.GET("/query-inspections/:license-plate", controllers.GetQueryInspectionsHelper)
	router.POST("/inspections", controllers.CreateQueryInspectionsHelper)

	// router.GET("/coordinates", controllers.GetCoordinates)

	router.GET("/auth", controllers.CheckAuthKey)
	router.POST("/auth", controllers.CreateAuthKeyWrapper)
	router.DELETE("/auth", controllers.DeleteAuthKey)

	router.GET("/statistics", controllers.GetStatistics)

	//router.Run("localhost:3000")
	http.ListenAndServe(":3000", router)
}
