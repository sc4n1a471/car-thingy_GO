package controllers

import (
	"Go_Thingy_GO/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// MARK: GetCar
func GetCar(ctx *gin.Context) {
	var car models.Car
	var returnData []models.Car

	isAccessGranted, error := GetAuthenticatedClient(ctx.Request)
	if error != nil || !isAccessGranted {
		ctx.IndentedJSON(http.StatusUnauthorized, models.Response{
			Status:  "fail",
			Message: "Access denied!",
		})
		return
	}

	car.ID = ctx.Param("license-plate")
	result := DB.Preload("Accidents").Preload("Mileage").Preload("Restrictions").First(&car)
	if result.Error != nil {
		SendError(result.Error.Error(), ctx)
		return
	}
	var inspections = GetInspectionsHelper(ctx, car.ID)
	car.Inspections = &inspections

	returnData = append(returnData, car)
	SendData(returnData, ctx)
}

// MARK: GetCars
func GetCars(ctx *gin.Context) {
	isAccessGranted, error := GetAuthenticatedClient(ctx.Request)
	if error != nil || !isAccessGranted {
		ctx.IndentedJSON(http.StatusUnauthorized, models.Response{
			Status:  "fail",
			Message: "Access denied!",
		})
		return
	}

	var returnCars []models.Car

	result := DB.Find(&returnCars)
	if result.Error != nil {
		SendError(result.Error.Error(), ctx)
		return
	}

	for i, car := range returnCars {
		if car.Mileage == nil {
			returnCars[i].Mileage = &[]models.Mileage{}
			// Empty array instead of null in JSON
			// Required because of the chart in SwiftUI
		}
	}

	// for _, licensePlate := range allLicensePlates {
	// 	var car models.CarResult
	// 	var coordinates models.Coordinate

	// 	car.LicensePlate = licensePlate

	// 	car.Specs = GetSpecs(ctx, licensePlate.LicensePlate, true)

	// 	// car.Accidents = []models.Accident{}
	// 	// car.Restrictions = []models.Restriction{}
	// 	car.Mileage = []models.Mileage{}
	// 	// car.Inspections = []models.InspectionResult{}

	// 	result := DB.Find(&coordinates, "car_id = ?", licensePlate.LicensePlate)
	// 	if result.Error != nil {
	// 		SendError(result.Error.Error(), ctx)
	// 		return
	// 	}
	// 	car.Coordinates = coordinates

	// 	returnCars = append(returnCars, car)
	// }

	if returnCars == nil {
		returnCars = []models.Car{}
	}

	SendData(returnCars, ctx)
}

// Returns all cars in the database with all information including inspection images
// func GetCarsAllData(ctx *gin.Context) {
// 	var allLicensePlates []models.LicensePlate
// 	var returnCars []models.CarResult
// 	result := DB.Find(&allLicensePlates)
// 	if result.Error != nil {
// 		SendError(result.Error.Error(), ctx)
// 		return
// 	}
// 	for _, licensePlate := range allLicensePlates {
// 		var car models.CarResult
// 		var coordinates models.Coordinate
// 		car.LicensePlate = licensePlate
// 		car.Specs = GetSpecs(ctx, licensePlate.LicensePlate, false)
// 		car.Accidents = getAccidents(ctx, licensePlate.LicensePlate)
// 		if car.Accidents == nil {
// 			return
// 		}
// 		car.Restrictions = GetRestrictions(ctx, licensePlate.LicensePlate)
// 		if car.Restrictions == nil {
// 			return
// 		}
// 		car.Mileage = GetMileages(ctx, licensePlate.LicensePlate)
// 		if car.Mileage == nil {
// 			return
// 		}
// 		result := DB.Find(&coordinates, "car_id = ?", licensePlate.LicensePlate)
// 		if result.Error != nil {
// 			SendError(result.Error.Error(), ctx)
// 			return
// 		}
// 		car.Coordinates = coordinates
// 		car.Inspections = []models.InspectionResult{}
// 		car.Inspections = GetInspectionsHelper(ctx, car.Specs.LicensePlate)
// 		if car.Inspections == nil {
// 			return
// 		}
// 		returnCars = append(returnCars, car)
// 	}
// 	if returnCars == nil {
// 		returnCars = []models.CarResult{}
// 	}
// 	SendData(returnCars, ctx)
// }

// MARK: CreateCar
func CreateCar(ctx *gin.Context) {
	isAccessGranted, error := GetAuthenticatedClient(ctx.Request)
	if error != nil || !isAccessGranted {
		ctx.IndentedJSON(http.StatusUnauthorized, models.Response{
			Status:  "fail",
			Message: "Access denied!",
		})
		return
	}

	var newCar models.Car
	var newAccidents []models.Accident
	var newRestrictions []models.Restriction
	var newMileages []models.Mileage
	var newInspections []models.Inspection

	if err := ctx.BindJSON(&newCar); err != nil {
		SendError(err.Error(), ctx)
		return
	}

	if newCar.Accidents != nil {
		newAccidents = *newCar.Accidents
	}
	if newCar.Restrictions != nil {
		newRestrictions = *newCar.Restrictions
	}
	if newCar.Mileage != nil {
		newMileages = *newCar.Mileage
	}
	if newCar.Inspections != nil {
		newInspections = *newCar.Inspections
	}

	fmts := "Creating car: %+v\n"
	fmt.Printf(fmts, newCar)
	fmts = "With accidents: %+v\n"
	fmt.Printf(fmts, newAccidents)
	fmts = "With restrictions: %+v\n"
	fmt.Printf(fmts, newRestrictions)
	fmts = "With mileages: %+v\n"
	fmt.Printf(fmts, newMileages)
	fmt.Printf("With inspections: %d items\n", len(newInspections))
	for _, insp := range newInspections {
		fmt.Printf(" - Inspection at %s\n", insp.Name)
	}

	tx := DB.Begin()

	var existingCar models.Car
	result := tx.First(&existingCar, "id = ?", newCar.ID)
	if result.RowsAffected == 0 {
		result := tx.Create(&newCar)
		if result.Error != nil {
			tx.Rollback()
			SendError(result.Error.Error(), ctx)
			return
		}
	} else {
		// Update existing car
		// In case of an unknown car but saved license plate and want to update its data, update only selected fields
		result := tx.
			Select(
				"Comment",
				"Brand",
				"Color",
				"EngineSize",
				"FirstReg",
				"FirstRegHun",
				"FuelType",
				"Gearbox",
				"Model",
				"NumOfOwners",
				"Performance",
				"Status",
				"TypeCode",
				"Year",
				"UpdatedAt").
			Save(&newCar)
		if result.Error != nil {
			tx.Rollback()
			SendError(Error.Error(), ctx)
			return
		}
	}

	// MARK: - Accidents
	for _, newAccident := range newAccidents {
		var existingAccident models.Accident
		checkResult := tx.Where(&models.Accident{
			CarID:        newAccident.CarID,
			AccidentDate: newAccident.AccidentDate,
		}).Find(&existingAccident)
		if checkResult.RowsAffected != 0 {
			continue
		}

		result := tx.Create(&newAccident)
		if result.Error != nil {
			tx.Rollback()
			SendError(result.Error.Error(), ctx)
			return
		}
	}

	// MARK: - Restrictions
	var existingRestrictions []models.Restriction
	result = DB.Find(&existingRestrictions, "car_id = ?", newCar.ID)
	if result.Error != nil {
		SendError(result.Error.Error(), ctx)
		return
	}

existingsLoop:
	for _, existingRestriction := range existingRestrictions {
		for _, newRestriction := range newRestrictions {
			if existingRestriction.Restriction == newRestriction.Restriction {
				continue existingsLoop
			}
		}
		tx.Model(&models.Restriction{}).
			Where(
				"car_id = ? AND restriction = ?",
				existingRestriction.CarID,
				existingRestriction.Restriction).
			Update("is_active", false)
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
			continue
			// 	tx.Rollback()
			// 	SendError(result.Error.Error(), ctx)
			// 	return
		}
	}

	// MARK: - Mileages
	for _, newMileage := range newMileages {
		checkResult := tx.Where("car_id = ? AND mileage = ?", newMileage.CarID, newMileage.Mileage).Find(&newMileage)
		if checkResult.RowsAffected != 0 {
			continue
		}

		result := tx.Create(&newMileage)
		if result.Error != nil {
			tx.Rollback()
			SendError(result.Error.Error(), ctx)
			return
		}
	}

	// MARK: - Inspections
	CreateInspectionHelper(ctx, newInspections, tx)

	tx.Commit()

	SendData("Car was uploaded successfully", ctx)
	return
}

// MARK: UpdateCar
func UpdateCar(ctx *gin.Context) {
	isAccessGranted, error := GetAuthenticatedClient(ctx.Request)
	if error != nil || !isAccessGranted {
		ctx.IndentedJSON(http.StatusUnauthorized, models.Response{
			Status:  "fail",
			Message: "Access denied!",
		})
		return
	}

	var updatedCar models.Car

	if err := ctx.BindJSON(&updatedCar); err != nil {
		SendError(err.Error(), ctx)
		return
	}

	tx := DB.Begin()

	result := tx.Save(&updatedCar)
	if result.Error != nil {
		tx.Rollback()
		SendError(Error.Error(), ctx)
		return
	}

	tx.Commit()

	SendData("Car was updated successfully", ctx)
	return
}

// MARK: DeleteCar
func DeleteCar(ctx *gin.Context) {
	isAccessGranted, error := GetAuthenticatedClient(ctx.Request)
	if error != nil || !isAccessGranted {
		ctx.IndentedJSON(http.StatusUnauthorized, models.Response{
			Status:  "fail",
			Message: "Access denied!",
		})
		return
	}

	var deletableLicensePlate models.Car

	deletableLicensePlate.ID = ctx.Param("license-plate")

	success := DeleteQueryInspections(ctx, deletableLicensePlate.ID, true)

	if !success {
		SendData("Inspections were not deleted successfully", ctx)
		return
	}

	result := DB.Delete(&deletableLicensePlate)

	if result.RowsAffected == 0 {
		SendError(result.Error.Error(), ctx)
		return
	}

	SendData("Car was deleted successfully", ctx)
}

// MARK: CreateLicensePlate
func CreateLicensePlate(ctx *gin.Context) {
	isAccessGranted, error := GetAuthenticatedClient(ctx.Request)
	if error != nil || !isAccessGranted {
		ctx.IndentedJSON(http.StatusUnauthorized, models.Response{
			Status:  "fail",
			Message: "Access denied!",
		})
		return
	}

	var newCar models.Car

	if err := ctx.BindJSON(&newCar); err != nil {
		SendError(err.Error(), ctx)
		return
	}

	tx := DB.Begin()
	result := tx.First(&newCar)
	if result.RowsAffected == 0 {
		result := tx.Create(&newCar)
		if result.Error != nil {
			tx.Rollback()
			SendError(result.Error.Error(), ctx)
			return
		}
	}

	tx.Commit()

	SendData("License plate was uploaded successfully", ctx)
	return
}

// MARK: UpdateLicensePlate
func UpdateLicensePLate(ctx *gin.Context) {
	isAccessGranted, error := GetAuthenticatedClient(ctx.Request)
	if error != nil || !isAccessGranted {
		ctx.IndentedJSON(http.StatusUnauthorized, models.Response{
			Status:  "fail",
			Message: "Access denied!",
		})
		return
	}

	var updatedCar models.Car

	if err := ctx.BindJSON(&updatedCar); err != nil {
		SendError(err.Error(), ctx)
		return
	}

	var oldLicensePlate = ctx.Param("license-plate")

	tx := DB.Begin()

	// Update license plate (oldLicensePlate) with new license plate (updatedCar.LicensePlate)
	result := tx.
		Model(&models.Car{}).
		Where("id = ?", oldLicensePlate).
		Update("id", updatedCar.ID)
	if result.Error != nil {
		tx.Rollback()
		SendError(result.Error.Error(), ctx)
		return
	}

	// Update inspections with new license plate
	result = tx.
		Model(&models.Inspection{}).
		Where("car_id = ?", oldLicensePlate).
		Update("car_id", updatedCar.ID)
	if result.Error != nil {
		tx.Rollback()
		SendError(result.Error.Error(), ctx)
		return
	}

	tx.Commit()

	SendData("License Plate was updated successfully", ctx)
}
