package application

import (
	"Go_Thingy/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func createCar(ctx *gin.Context) {
	var newCar models.Car
	var newLicensePlate models.LicensePlate
	var newSpecs models.Specs
	var newAccidents []models.Accident
	var newRestrictions []models.Restriction
	var newMileages []models.Mileage
	var newCoordinates models.Coordinate
	var newInspections []models.Inspection

	if err := ctx.BindJSON(&newCar); err != nil {
		sendError(err.Error(), ctx)
		return
	}

	newLicensePlate = newCar.LicensePlate
	newSpecs = newCar.Specs
	newAccidents = newCar.Accidents
	newRestrictions = newCar.Restrictions
	newMileages = newCar.Mileage
	newCoordinates = newCar.Coordinates
	newInspections = newCar.Inspections

	tx := DB.Begin()

	result := tx.First(&newLicensePlate)
	if result.RowsAffected == 0 {
		result := tx.Create(&newLicensePlate)
		if result.Error != nil {
			tx.Rollback()
			sendError(result.Error.Error(), ctx)
			return
		}
	} else {
		result = tx.
			Model(&newLicensePlate).
			Select("updated_at").
			Updates(models.LicensePlate{
				UpdatedAt: newCar.LicensePlate.UpdatedAt,
			})
		if result.Error != nil {
			tx.Rollback()
			sendError(result.Error.Error(), ctx)
			return
		}

	}

	result = tx.First(&newSpecs)
	if result.RowsAffected == 0 {
		result := tx.Create(&newSpecs)
		if result.Error != nil {
			tx.Rollback()
			sendError(result.Error.Error(), ctx)
			return
		}
	} else {
		result := tx.Save(&newSpecs)
		if result.Error != nil {
			tx.Rollback()
			sendError(Error.Error(), ctx)
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
	result = DB.Find(&existingRestrictions, "license_plate = ?", newLicensePlate.LicensePlate)
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

	result = tx.Find(&newCoordinates, "license_plate = ?", newLicensePlate.LicensePlate)
	if result.RowsAffected == 0 {
		result := tx.Create(&newCoordinates)
		if result.Error != nil {
			tx.Rollback()
			sendError(result.Error.Error(), ctx)
			return
		}
	}

	createInspectionHelper(ctx, newInspections, tx)

	tx.Commit()

	sendData("Car was uploaded successfully", ctx)
	return
}

func createLicensePlate(ctx *gin.Context) {
	var newCar models.Car
	var newLicensePlate models.LicensePlate
	var newCoordinates models.Coordinate

	if err := ctx.BindJSON(&newCar); err != nil {
		sendError(err.Error(), ctx)
		return
	}

	newLicensePlate = newCar.LicensePlate
	newCoordinates = newCar.Coordinates

	tx := DB.Begin()
	result := tx.First(&newLicensePlate)
	if result.RowsAffected == 0 {
		result := tx.Create(&newLicensePlate)
		if result.Error != nil {
			tx.Rollback()
			sendError(result.Error.Error(), ctx)
			return
		}
	}

	result = tx.Find(&newCoordinates, "license_plate = ?", newLicensePlate.LicensePlate)
	if result.RowsAffected == 0 {
		result := tx.Create(&newCoordinates)
		if result.Error != nil {
			tx.Rollback()
			sendError(result.Error.Error(), ctx)
			return
		}
	}

	tx.Commit()

	sendData("License plate was uploaded successfully", ctx)
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

	sendData("Inspections were uploaded successfully", ctx)
	return
}
