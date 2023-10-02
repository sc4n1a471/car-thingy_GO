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

	var car models.CarResult
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

	car.Inspections = getInspectionsHelper(ctx, requested.LicensePlate)
	if car.Inspections == nil {
		return
	}

	ctx.IndentedJSON(http.StatusOK, car)
}

func getCars(ctx *gin.Context) {
	var allSpecs []models.Specs

	var returnCars []models.CarResult

	result := DB.Find(&allSpecs)
	if result.Error != nil {
		sendError(result.Error.Error(), ctx)
		return
	}

	for _, specs := range allSpecs {
		var car models.CarResult
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

		car.Inspections = getInspectionsHelper(ctx, car.Specs.LicensePlate)
		if car.Inspections == nil {
			return
		}

		returnCars = append(returnCars, car)
	}

	ctx.IndentedJSON(http.StatusOK, returnCars)
}

func getInspectionsHelper(ctx *gin.Context, licensePlate string) []models.InspectionResult {
	var inspections []models.Inspection
	var inspectionResults []models.InspectionResult

	result := DB.Find(&inspections, "license_plate = ?", licensePlate)
	if result.Error != nil {
		sendError(result.Error.Error(), ctx)
		return nil
	}

	for _, inspection := range inspections {
		var inspectionResult models.InspectionResult
		inspectionResult.LicensePlate = inspection.LicensePlate
		inspectionResult.Name = inspection.Name
		inspectionResult.Base64 = inspection.LicensePlate + inspection.LicensePlate
		// TODO: Convert image to base64

		inspectionResults = append(inspectionResults, inspectionResult)
	}

	return inspectionResults
}

func getInspections(ctx *gin.Context) {
	var inspectionResults []models.InspectionResult
	licensePlate := ctx.Param("license_plate")

	inspectionResults = getInspectionsHelper(ctx, licensePlate)
	if inspectionResults == nil {
		return
	}

	ctx.IndentedJSON(http.StatusOK, inspectionResults)
}

func createCar(ctx *gin.Context) {
	var newCar models.Car
	var newSpecs models.Specs
	var newAccidents []models.Accident
	var newRestrictions []models.Restriction
	var newMileages []models.Mileage
	var newGeneral models.General
	var newInspections []models.Inspection

	if err := ctx.BindJSON(&newCar); err != nil {
		sendError(err.Error(), ctx)
		return
	}

	newSpecs = newCar.Specs
	newAccidents = newCar.Accidents
	newRestrictions = newCar.Restrictions
	newMileages = newCar.Mileage
	newGeneral = newCar.General
	newInspections = newCar.Inspections

	tx := DB.Begin()
	result := tx.First(&newSpecs)
	if result.RowsAffected == 0 {
		result := tx.Create(&newSpecs)
		if result.Error != nil {
			tx.Rollback()
			sendError(result.Error.Error(), ctx)
			return
		}
	}

	for _, newAccident := range newAccidents {

		var existingAccident models.Accident
		checkResult := tx.Where(&models.Accident{
			LicensePlate: newAccident.LicensePlate,
			AccidentDate: newAccident.AccidentDate,
		}).Find(&existingAccident)
		if checkResult.RowsAffected != 0 {
			continue
		}

		result := tx.Create(&newAccident)
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
				continue existingsLoop
			}
		}
		fmt.Println(existingRestriction)
		tx.Model(&models.Restriction{}).
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

		result := tx.Create(&newRestriction)
		if result.Error != nil {
			tx.Rollback()
			sendError(result.Error.Error(), ctx)
			return
		}
	}

	for _, newMileage := range newMileages {

		var tempMileage models.Mileage
		checkResult := tx.Where(&models.Mileage{
			LicensePlate: newMileage.LicensePlate,
			Mileage:      newMileage.Mileage,
		}).Find(&tempMileage)
		if checkResult.RowsAffected != 0 {
			continue
		}

		result := tx.Create(&newMileage)
		if result.Error != nil {
			tx.Rollback()
			sendError(result.Error.Error(), ctx)
			return
		}
	}

	result = tx.Find(&newGeneral, "license_plate = ?", newSpecs.LicensePlate)
	if result.RowsAffected == 0 {
		result := tx.Create(&newGeneral)
		if result.Error != nil {
			tx.Rollback()
			sendError(result.Error.Error(), ctx)
			return
		}
	}

	createInspectionHelper(ctx, newInspections, tx)

	tx.Commit()
	ctx.IndentedJSON(http.StatusCreated, models.Response{
		Status:  "success",
		Message: "Car was uploaded successfully",
	})
	return
}

func createInspectionHelper(ctx *gin.Context, newInspections []models.Inspection, tx *gorm.DB) bool {
	for _, newInspection := range newInspections {

		var existingInspection models.Inspection
		checkResult := tx.Where(&models.Inspection{
			LicensePlate:  newInspection.LicensePlate,
			Name:          newInspection.Name,
			ImageLocation: newInspection.ImageLocation,
		}).Find(&existingInspection)
		if checkResult.RowsAffected != 0 {
			continue
		}

		result := tx.Create(&newInspection)
		if result.Error != nil {
			tx.Rollback()
			sendError(result.Error.Error(), ctx)
			return false
		}
	}
	return true
}

func createInspections(ctx *gin.Context) {
	var newInspections []models.Inspection

	if err := ctx.BindJSON(&newInspections); err != nil {
		sendError(err.Error(), ctx)
		return
	}

	tx := DB.Begin()

	successful := createInspectionHelper(ctx, newInspections, tx)

	if !successful {
		return
	}

	tx.Commit()
	ctx.IndentedJSON(http.StatusCreated, models.Response{
		Status:  "success",
		Message: "Inspections were uploaded successfully",
	})
	return
}

func updateCar(ctx *gin.Context) {
	var updatedCar models.Car
	var updatedSpecs models.Specs
	var updatedGeneral models.General

	if err := ctx.BindJSON(&updatedCar); err != nil {
		sendError(err.Error(), ctx)
		return
	}

	updatedSpecs = updatedCar.Specs
	updatedGeneral = updatedCar.General

	tx := DB.Begin()

	result := tx.Save(&updatedSpecs)
	if result.Error != nil {
		tx.Rollback()
		sendError(Error.Error(), ctx)
		return
	}

	result = tx.
		Model(&updatedGeneral).
		Select("latitude", "longitude", "created_at").
		Updates(models.General{
			Latitude:  updatedGeneral.Latitude,
			Longitude: updatedGeneral.Longitude,
			CreatedAt: updatedGeneral.CreatedAt,
		})
	if result.Error != nil {
		tx.Rollback()
		sendError(result.Error.Error(), ctx)
		return
	}

	tx.Commit()
	ctx.IndentedJSON(http.StatusCreated, models.Response{
		Status:  "success",
		Message: "Car was updated successfully",
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

func deleteInspectionsHelper(ctx *gin.Context, licensePlate string) bool {
	var inspections []models.Inspection
	result := DB.Where("license_plate = ?", licensePlate).Delete(&inspections)

	if result.RowsAffected == 0 {
		sendError(result.Error.Error(), ctx)
		return false
	}
	return true
}

func deleteInspections(ctx *gin.Context) {
	licensePlate := ctx.Param("license_plate")

	success := deleteInspectionsHelper(ctx, licensePlate)

	if !success {
		return
	}
	ctx.IndentedJSON(http.StatusCreated, models.Response{
		Status:  "success",
		Message: "Inspections were deleted successfully",
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
	router.PUT("/cars", updateCar)
	router.DELETE("/cars/:license_plate", deleteCar)

	router.GET("/inspections/:license_plate", getInspections)
	router.POST("/inspections", createInspections)
	router.DELETE("/inspections/:license_plate", deleteInspections)

	router.Run("localhost:3000")
}
