package controllers

import (
	"Go_Thingy_GO/models"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MARK: Get query inspections
func GetQueryInspectionsWrapper(ctx *gin.Context) {
	isAccessGranted, error := CheckAuthKey(ctx.Request)
	if error != nil || !isAccessGranted {
		SendError("Access denied for getting query inspections", http.StatusUnauthorized, ctx)
		return
	}

	var inspections []models.QueryInspection
	licensePlate := ctx.Param("license-plate")

	inspections = GetQueryInspections(ctx, licensePlate)
	if inspections == nil {
		return
	}

	SendData(inspections, ctx)
}

// Returns all inspections for a given license plate
func GetQueryInspections(ctx *gin.Context, licensePlate string) []models.QueryInspection {
	isAccessGranted, error := CheckAuthKey(ctx.Request)
	if error != nil || !isAccessGranted {
		SendError("Access denied for getting query inspections", http.StatusUnauthorized, ctx)
		return nil
	}

	slog.Info("Getting query inspections to " + ctx.ClientIP())

	var inspections []models.QueryInspection

	result := DB.Find(&inspections, "car_id = ?", licensePlate)
	if result.Error != nil {
		SendError("Could not find inspections: "+result.Error.Error(), http.StatusInternalServerError, ctx)
		return nil
	} else if result.RowsAffected == 0 {
		return []models.QueryInspection{}
	}

	slog.Info("Found ", result.RowsAffected, " query inspections for car ", licensePlate)

	for i, inspection := range inspections {
		inspections[i].Base64 = ConvertImagesToBase64(inspection.ImageLocation)
	}
	slog.Info("Converted images to Base64 for car ", licensePlate)

	return inspections
}

// MARK: Create query inspections
func CreateQueryInspectionsWrapper(ctx *gin.Context) {
	isAccessGranted, error := CheckAuthKey(ctx.Request)
	if error != nil || !isAccessGranted {
		SendError("Access denied for creating query inspections", http.StatusUnauthorized, ctx)
		return
	}

	var newInspections []models.QueryInspection

	if err := ctx.BindJSON(&newInspections); err != nil {
		SendError("Could not parse JSON: "+err.Error(), http.StatusBadRequest, ctx)
		return
	}

	tx := DB.Begin()

	successful := CreateQueryInspection(ctx, newInspections, tx)

	if !successful {
		return
	}

	tx.Commit()

	SendData("Inspections were uploaded successfully", ctx)
}
func CreateQueryInspection(ctx *gin.Context, newInspections []models.QueryInspection, tx *gorm.DB) bool {
	isAccessGranted, error := CheckAuthKey(ctx.Request)
	if error != nil || !isAccessGranted {
		SendError("Access denied for creating query inspections", http.StatusUnauthorized, ctx)
		return false
	}
	slog.Info("Creating ", len(newInspections), " query inspections")

	var successfulCreations int = 0

	for _, newInspection := range newInspections {
		slog.Info("Creating query inspection: " + newInspection.Name)
		checkResult := tx.Where("name = ? and car_id = ?", newInspection.Name, newInspection.CarID).First(&newInspection)
		if checkResult.RowsAffected != 0 {
			slog.Info("Query inspection already exists: " + newInspection.Name + " for car " + newInspection.CarID)
			continue
		}

		result := tx.Create(&newInspection)
		if result.Error != nil {
			tx.Rollback()
			SendError("Could not create query inspection: "+result.Error.Error(), http.StatusInternalServerError, ctx)
			return false
		}
		successfulCreations += 1
	}
	slog.Info("Successfully created ", successfulCreations, " query inspections")
	return true
}

// MARK: Delete query inspections and their images
// Deletes all query inspections and their images for a given license plate
func DeleteQueryInspections(ctx *gin.Context, licensePlate string, imagesOnly bool) bool {
	isAccessGranted, error := CheckAuthKey(ctx.Request)
	if error != nil || !isAccessGranted {
		SendError("Access denied for deleting query inspections", http.StatusUnauthorized, ctx)
		return false
	}

	var inspections []models.QueryInspection

	result := DB.Find(&inspections, "car_id = ?", licensePlate)
	if result.RowsAffected == 0 {
		return true
	}

	for _, inspection := range inspections {
		errorResult := os.RemoveAll(inspection.ImageLocation)
		if errorResult != nil {
			SendError("Could not delete inspection images: "+errorResult.Error(), http.StatusInternalServerError, ctx)
			return false
		}
	}

	result = DB.Where("car_id = ?", licensePlate).Delete(&inspections)
	if result.Error != nil {
		SendError("Could not delete query inspections: "+result.Error.Error(), http.StatusInternalServerError, ctx)
		return false
	}
	return true
}

// Delete all older query inspections and their images that were not saved, used as cleanup function only
// MARK: Cleanup function
func DeleteOldQueryInspections() bool {
	slog.Info("Deleting old queries...")
	var inspections []models.QueryInspection
	var deletedSuccessfully int64 = 0

	// SELECT car_id FROM `query_inspections` qi where (select count(*) from inspections where car_id = qi.car_id) = 0 group by car_id;

	var queryInspectionCarIds []string
	result := DB.Table("query_inspections").Select("car_id").Where("car_id NOT IN (SELECT car_id FROM inspections GROUP BY car_id)").Group("car_id").Scan(&queryInspectionCarIds)
	if result.Error != nil {
		slog.Error("Error fetching old query inspections: " + result.Error.Error())
		return false
	}

	for _, id := range queryInspectionCarIds {
		slog.Info("To be deleted query: " + id)
		result = DB.Find(&inspections, "car_id = ?", id)
		if result.RowsAffected == 0 {
			return true
		}
		if result.Error != nil {
			slog.Error("Error fetching old query inspections: " + result.Error.Error())
			return false
		}

		for _, inspection := range inspections {
			errorResult := os.RemoveAll(inspection.ImageLocation)
			if errorResult != nil {
				slog.Error("Error deleting old query inspection images: " + errorResult.Error())
				return false
			}
		}

		result = DB.Where("car_id = ?", id).Delete(&inspections)
		if result.Error != nil {
			slog.Error("Error deleting old query inspections: " + result.Error.Error())
			return false
		}
		deletedSuccessfully += 1
	}
	slog.Info("Deleted ", deletedSuccessfully, " query inspections")
	return true
}
