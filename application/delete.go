package application

import (
	"Go_Thingy/models"
	"github.com/gin-gonic/gin"
	"os"
)

func deleteCar(ctx *gin.Context) {
	var deletableLicensePlate models.LicensePlate

	deletableLicensePlate.LicensePlate = ctx.Param("license_plate")

	result := DB.Where("license_plate = ?", deletableLicensePlate.LicensePlate).Delete(&deletableLicensePlate)

	if result.RowsAffected == 0 {
		sendError(result.Error.Error(), ctx)
		return
	}

	sendData("Car was deleted successfully", ctx)
}

func deleteInspectionsHelper(ctx *gin.Context, licensePlate string) bool {
	var inspections []models.Inspection

	result := DB.Find(&inspections, "license_plate = ?", licensePlate)
	if result.RowsAffected == 0 {
		sendError(result.Error.Error(), ctx)
		return false
	}

	for _, inspection := range inspections {
		errorResult := os.RemoveAll(inspection.ImageLocation)
		if errorResult != nil {
			sendError(errorResult.Error(), ctx)
			return false
		}
	}

	result = DB.Where("license_plate = ?", licensePlate).Delete(&inspections)
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
