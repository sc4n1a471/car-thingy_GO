package application

import (
	"Go_Thingy/models"
	"github.com/gin-gonic/gin"
)

func updateCar(ctx *gin.Context) {
	var updatedCar models.Car
	var updatedLicensePlate models.LicensePlate
	var updatedCoordinates models.Coordinate

	if err := ctx.BindJSON(&updatedCar); err != nil {
		sendError(err.Error(), ctx)
		return
	}

	updatedLicensePlate = updatedCar.LicensePlate
	updatedCoordinates = updatedCar.Coordinates

	tx := DB.Begin()

	result := tx.Save(&updatedLicensePlate)
	if result.Error != nil {
		tx.Rollback()
		sendError(Error.Error(), ctx)
		return
	}

	result = tx.
		Model(&updatedCoordinates).
		Select("latitude", "longitude").
		Updates(models.Coordinate{
			Latitude:  updatedCoordinates.Latitude,
			Longitude: updatedCoordinates.Longitude,
		})
	if result.Error != nil {
		tx.Rollback()
		sendError(result.Error.Error(), ctx)
		return
	}

	tx.Commit()

	sendData("Car was updated successfully", ctx)
	return
}

func updateLicensePLate(ctx *gin.Context) {
	var updatedLicensePlate models.LicensePlate

	if err := ctx.BindJSON(&updatedLicensePlate); err != nil {
		sendError(err.Error(), ctx)
		return
	}

	var oldLicensePlate = ctx.Param("license_plate")

	tx := DB.Begin()

	// Update license plate (oldLicensePlate) with new license plate (updatedLicensePlate.LicensePlate)
	result := tx.
		Model(&models.LicensePlate{}).
		Where("license_plate = ?", oldLicensePlate).
		Update("license_plate", updatedLicensePlate.LicensePlate)
	if result.Error != nil {
		tx.Rollback()
		sendError(result.Error.Error(), ctx)
		return
	}

	// Update inspections with new license plate
	result = tx.
		Model(&models.Inspection{}).
		Where("license_plate = ?", oldLicensePlate).
		Update("license_plate", updatedLicensePlate.LicensePlate)
	if result.Error != nil {
		tx.Rollback()
		sendError(result.Error.Error(), ctx)
		return
	}

	tx.Commit()

	sendData("License Plate was updated successfully", ctx)
	return
}
