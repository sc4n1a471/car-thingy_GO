package application

import (
	"Go_Thingy/models"
	"github.com/gin-gonic/gin"
)

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

	sendData("Car was updated successfully", ctx)
	return
}
