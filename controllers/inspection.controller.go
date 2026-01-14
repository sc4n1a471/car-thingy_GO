package controllers

import (
	"Go_Thingy_GO/models"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MARK: Normal inspections
func GetInspectionsHelper(ctx *gin.Context) {
	isAccessGranted, error := CheckAuthKey(ctx.Request)
	if error != nil || !isAccessGranted {
		SendError("Access denied for getting inspections", http.StatusUnauthorized, ctx)
		return
	}

	var inspections []models.Inspection
	licensePlate := ctx.Param("license-plate")

	inspections = GetInspections(ctx, licensePlate)
	if inspections == nil {
		return
	}

	SendData(inspections, ctx)
}

// Returns all inspections for a given license plate
func GetInspections(ctx *gin.Context, licensePlate string) []models.Inspection {
	isAccessGranted, error := CheckAuthKey(ctx.Request)
	if error != nil || !isAccessGranted {
		SendError("Access denied for getting inspections", http.StatusUnauthorized, ctx)
		return nil
	}

	slog.Info("Getting inspections to " + ctx.ClientIP())

	var inspections []models.Inspection

	result := DB.Find(&inspections, "car_id = ?", licensePlate)
	if result.Error != nil {
		SendError("Could not find inspections: "+result.Error.Error(), http.StatusInternalServerError, ctx)
		return nil
	} else if result.RowsAffected == 0 {
		return []models.Inspection{}
	}

	slog.Info("Found ", result.RowsAffected, " inspections for car ", licensePlate)

	for i, inspection := range inspections {
		inspections[i].Base64 = ConvertImagesToBase64(inspection.ImageLocation)
	}

	slog.Info("Converted images to Base64 for car ", licensePlate)

	return inspections
}

// MARK: Create inspections
func CreateInspectionHelper(ctx *gin.Context, newInspections []models.Inspection, tx *gorm.DB) bool {
	isAccessGranted, error := CheckAuthKey(ctx.Request)
	if error != nil || !isAccessGranted {
		SendError("Access denied for creating inspections", http.StatusUnauthorized, ctx)
		return false
	}

	for _, newInspection := range newInspections {
		checkResult := tx.Where("name = ? and car_id = ?", newInspection.Name, newInspection.CarID).Find(&newInspection)
		if checkResult.RowsAffected != 0 {
			continue
		}

		result := tx.Create(&newInspection)
		if result.Error != nil {
			tx.Rollback()
			SendError("Error creating new inspection: "+result.Error.Error(), http.StatusInternalServerError, ctx)
			return false
		}
	}
	return true
}
func CreateInspections(ctx *gin.Context) {
	isAccessGranted, error := CheckAuthKey(ctx.Request)
	if error != nil || !isAccessGranted {
		SendError("Access denied for creating inspections", http.StatusUnauthorized, ctx)
		return
	}

	var newInspections []models.Inspection

	if err := ctx.BindJSON(&newInspections); err != nil {
		SendError("Could not parse JSON: "+err.Error(), http.StatusBadRequest, ctx)
		return
	}

	tx := DB.Begin()

	successful := CreateInspectionHelper(ctx, newInspections, tx)

	if !successful {
		return
	}

	tx.Commit()

	SendData("Inspections were uploaded successfully", ctx)
}
