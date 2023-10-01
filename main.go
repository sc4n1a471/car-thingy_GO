package main

import (
	"Go_Thingy/models"
	"fmt"
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

func getCar(ctx *gin.Context) {
	var requested models.Specs

	requested.LicensePlate = ctx.Param("license_plate")
	result := DB.First(&requested)
	if result.Error != nil {
		sendError(result.Error.Error(), ctx)
		return
	}

	var car models.Car
	var accidents []models.Accident
	var restrictions []models.Restriction
	var mileages []models.Mileage
	var general models.General

	car.Specs = requested

	accidentResults := DB.Find(&accidents, "license_plate = ?", car.Specs.LicensePlate)
	if accidentResults.Error != nil {
		sendError(accidentResults.Error.Error(), ctx)
		return
	}
	car.Accidents = accidents

	restrictionResult := DB.Find(&restrictions, "license_plate = ? AND active = true", car.Specs.LicensePlate)
	if restrictionResult.Error != nil {
		sendError(restrictionResult.Error.Error(), ctx)
		return
	}
	car.Restrictions = restrictions

	mileageResult := DB.Find(&mileages, "license_plate = ?", car.Specs.LicensePlate)
	if mileageResult.Error != nil {
		sendError(mileageResult.Error.Error(), ctx)
		return
	}
	car.Mileage = mileages

	generalResult := DB.Find(&general, "license_plate = ?", car.Specs.LicensePlate)
	if generalResult.Error != nil {
		sendError(generalResult.Error.Error(), ctx)
		return
	}
	car.General = general

	ctx.IndentedJSON(http.StatusOK, car)
}

func getCars(ctx *gin.Context) {
	var allSpecs []models.Specs

	var returnCars []models.Car

	result := DB.Find(&allSpecs)
	if result.Error != nil {
		sendError(result.Error.Error(), ctx)
		return
	}

	for _, specs := range allSpecs {
		var car models.Car
		var accidents []models.Accident
		var restrictions []models.Restriction
		var mileages []models.Mileage
		var general models.General

		car.Specs = specs

		result := DB.Find(&accidents, "license_plate = ?", specs.LicensePlate)
		if result.Error != nil {
			sendError(result.Error.Error(), ctx)
			return
		}
		car.Accidents = accidents

		result = DB.Find(&restrictions, "license_plate = ?", specs.LicensePlate)
		if result.Error != nil {
			sendError(result.Error.Error(), ctx)
			return
		}
		car.Restrictions = restrictions

		result = DB.Find(&mileages, "license_plate = ?", specs.LicensePlate)
		if result.Error != nil {
			sendError(result.Error.Error(), ctx)
			return
		}
		car.Mileage = mileages

		generalResult := DB.Find(&general, "license_plate = ?", car.Specs.LicensePlate)
		if generalResult.Error != nil {
			sendError(generalResult.Error.Error(), ctx)
			return
		}
		car.General = general

		returnCars = append(returnCars, car)
	}

	ctx.IndentedJSON(http.StatusOK, returnCars)
}

func createCar(ctx *gin.Context) {
	var newSpecs models.Specs
	var newAccidents []models.Accident
	var newRestrictions []models.Restriction
	var newMileages []models.Mileage
	var newCar models.Car

	if err := ctx.BindJSON(&newCar); err != nil {
		ctx.IndentedJSON(http.StatusConflict, models.Response{
			Status:  "fail",
			Message: err.Error(),
		})
		return
	}

	newSpecs = newCar.Specs
	newAccidents = newCar.Accidents
	newRestrictions = newCar.Restrictions
	newMileages = newCar.Mileage

	tx := DB.Begin()
	result := DB.First(&newSpecs)
	if result.RowsAffected == 0 {
		result := DB.Create(&newSpecs)
		if result.Error != nil {
			tx.Rollback()
			sendError(result.Error.Error(), ctx)
			return
		}
	}

	for _, newAccident := range newAccidents {

		var existingAccident models.Accident
		checkResult := DB.Where(&models.Accident{
			LicensePlate: newAccident.LicensePlate,
			AccidentDate: newAccident.AccidentDate,
		}).Find(&existingAccident)
		if checkResult.RowsAffected != 0 {
			continue
		}

		result := DB.Create(&newAccident)
		if result.Error != nil {
			tx.Rollback()
			sendError(result.Error.Error(), ctx)
			return
		}
	}

	var existingRestrictions []models.Restriction
	result = DB.Find(&existingRestrictions, "license_plate = ?", newSpecs.LicensePlate)
	if result.Error != nil {
		sendError(result.Error.Error(), ctx)
		return
	}

existingsLoop:
	for _, existingRestriction := range existingRestrictions {
		for _, newRestriction := range newRestrictions {
			if existingRestriction.Restriction == newRestriction.Restriction {
				fmt.Println(existingRestriction.Restriction)
				fmt.Println(newRestriction.Restriction)
				fmt.Println(existingRestriction.Restriction == newRestriction.Restriction)
				continue existingsLoop
			}
		}
		fmt.Println(existingRestriction)
		DB.Model(&models.Restriction{}).
			Where(
				"license_plate = ? AND restriction = ?",
				existingRestriction.LicensePlate,
				existingRestriction.Restriction).
			Update("active", false)
	}

newsLoop:
	for _, newRestriction := range newRestrictions {
		for _, existingRestriction := range existingRestrictions {
			if existingRestriction.Restriction == newRestriction.Restriction {
				continue newsLoop
			}
		}

		result := DB.Create(&newRestriction)
		if result.Error != nil {
			tx.Rollback()
			sendError(result.Error.Error(), ctx)
			return
		}
	}

	for _, newMileage := range newMileages {

		var tempMileage models.Mileage
		checkResult := DB.Where(&models.Mileage{
			LicensePlate: newMileage.LicensePlate,
			Mileage:      newMileage.Mileage,
		}).Find(&tempMileage)
		if checkResult.RowsAffected != 0 {
			continue
		}

		result := DB.Create(&newMileage)
		if result.Error != nil {
			tx.Rollback()
			sendError(result.Error.Error(), ctx)
			return
		}
	}

	tx.Commit()
	ctx.IndentedJSON(http.StatusCreated, models.Response{
		Status:  "success",
		Message: "Car was uploaded successfully",
	})
	return
}

func deleteCar(ctx *gin.Context) {
	var deletableSpecs models.Specs

	deletableSpecs.LicensePlate = ctx.Param("license_plate")

	result := DB.Where("license_plate = ?", deletableSpecs.LicensePlate).Delete(&deletableSpecs)

	if result.RowsAffected == 0 {
		sendError(result.Error.Error(), ctx)
		return
	}

	ctx.IndentedJSON(http.StatusCreated, models.Response{
		Status:  "success",
		Message: "Car was deleted successfully",
	})
}

func sendError(error string, ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusConflict, models.Response{
		Status:  "fail",
		Message: error,
	})
}

func main() {
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
	router.DELETE("/cars/:license_plate", deleteCar)

	router.Run("localhost:3000")
}
