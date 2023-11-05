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
