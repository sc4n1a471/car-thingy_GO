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
	// MARK: Car endpoints
	router.GET("/cars/:license-plate", controllers.GetCar)
	router.GET("/cars", controllers.GetCars)
	// router.GET("/cars-all-data", controllers.GetCarsAllData) // For testing purposes only
	router.POST("/cars", controllers.CreateCar)
	router.PUT("/cars", controllers.UpdateCar)
	router.DELETE("/cars/:license-plate", controllers.DeleteCar)

	// MARK: License Plate endpoints
	router.POST("/license-plate", controllers.CreateLicensePlate)
	router.PUT("/license-plate/:license-plate", controllers.UpdateLicensePlate)

	// MARK: Inspection endpoints
	router.GET("/inspections/:license-plate", controllers.GetInspectionsHelper)
	router.GET("/query-inspections/:license-plate", controllers.GetQueryInspectionsWrapper)
	router.POST("/inspections", controllers.CreateQueryInspectionsWrapper)

	// MARK: Auth Key endpoints
	router.GET("/auth", controllers.CheckAuthKeyWrapper)
	router.POST("/auth", controllers.CreateAuthKeyWrapper)
	router.DELETE("/auth", controllers.DeleteAuthKey)

	// MARK: Statistics endpoints
	router.GET("/statistics", controllers.GetStatistics)

	http.ListenAndServe(":3000", router)
}
