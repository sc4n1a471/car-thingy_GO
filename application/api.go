package application

import (
	"Go_Thingy_GO/controllers"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Api() {
	err := controllers.SetupDatabase()
	if err != nil {
		fmt.Print(err.Error())
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
	router.PUT("/license-plate/:license-plate", controllers.UpdateLicensePLate)

	router.GET("/inspections/:license-plate", controllers.GetInspections)
	router.GET("/query-inspections/:license-plate", controllers.GetQueryInspections)
	router.POST("/inspections", controllers.CreateQueryInspections)
	router.DELETE("/query-inspections/:license-plate", controllers.DeleteQueryInspections)

	// router.GET("/coordinates", controllers.GetCoordinates)

	router.GET("/auth", controllers.CheckAuthKey)
	router.POST("/auth", controllers.CreateAuthKeyWrapper)
	router.DELETE("/auth", controllers.DeleteAuthKey)

	router.GET("/statistics", controllers.GetStatistics)

	//router.Run("localhost:3000")
	http.ListenAndServe(":3000", router)
}
