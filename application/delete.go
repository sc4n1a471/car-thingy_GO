package application

import (
	"Go_Thingy/models"
	"github.com/gin-gonic/gin"
)

func deleteCar(ctx *gin.Context) {
	var deletableSpecs models.Specs

	deletableSpecs.LicensePlate = ctx.Param("license_plate")

	result := DB.Where("license_plate = ?", deletableSpecs.LicensePlate).Delete(&deletableSpecs)

	if result.RowsAffected == 0 {
		sendError(result.Error.Error(), ctx)
		return
	}

	sendData("Car was deleted successfully", ctx)
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

	sendData("Inspections were deleted successfully", ctx)
}
