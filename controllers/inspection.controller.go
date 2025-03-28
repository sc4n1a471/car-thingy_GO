package controllers

import (
	"Go_Thingy_GO/models"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MARK: Normal inspections
func GetInspections(ctx *gin.Context) {
	isAccessGranted, error := GetAuthenticatedClient(ctx.Request)
	if error != nil || !isAccessGranted {
		ctx.IndentedJSON(http.StatusUnauthorized, models.Response{
			Status:  "fail",
			Message: "Access denied!",
		})
		return
	}

	var inspections []models.Inspection
	licensePlate := ctx.Param("license-plate")

	inspections = GetInspectionsHelper(ctx, licensePlate)
	if inspections == nil {
		return
	}

	SendData(inspections, ctx)
}

// Returns all inspections for a given license plate
func GetInspectionsHelper(ctx *gin.Context, licensePlate string) []models.Inspection {
	isAccessGranted, error := GetAuthenticatedClient(ctx.Request)
	if error != nil || !isAccessGranted {
		ctx.IndentedJSON(http.StatusUnauthorized, models.Response{
			Status:  "fail",
			Message: "Access denied!",
		})
		return nil
	}

	var inspections []models.Inspection

	result := DB.Find(&inspections, "car_id = ?", licensePlate)
	if result.Error != nil {
		SendError(result.Error.Error(), ctx)
		return nil
	} else if result.RowsAffected == 0 {
		return []models.Inspection{}
	}

	for i, inspection := range inspections {
		inspections[i].Base64 = ConvertImagesToBase64(inspection.ImageLocation)
	}

	return inspections
}

func CreateInspections(ctx *gin.Context) {
	isAccessGranted, error := GetAuthenticatedClient(ctx.Request)
	if error != nil || !isAccessGranted {
		ctx.IndentedJSON(http.StatusUnauthorized, models.Response{
			Status:  "fail",
			Message: "Access denied!",
		})
		return
	}

	var newInspections []models.Inspection

	if err := ctx.BindJSON(&newInspections); err != nil {
		SendError(err.Error(), ctx)
		return
	}

	tx := DB.Begin()

	successful := CreateInspectionHelper(ctx, newInspections, tx)

	if !successful {
		return
	}

	tx.Commit()

	SendData("Inspections were uploaded successfully", ctx)
	return
}

func CreateInspectionHelper(ctx *gin.Context, newInspections []models.Inspection, tx *gorm.DB) bool {
	isAccessGranted, error := GetAuthenticatedClient(ctx.Request)
	if error != nil || !isAccessGranted {
		ctx.IndentedJSON(http.StatusUnauthorized, models.Response{
			Status:  "fail",
			Message: "Access denied!",
		})
		return false
	}

	for _, newInspection := range newInspections {
		checkResult := tx.Where("name = ?", newInspection.Name).Find(&newInspection)
		if checkResult.RowsAffected != 0 {
			continue
		}

		result := tx.Create(&newInspection)
		if result.Error != nil {
			tx.Rollback()
			SendError(result.Error.Error(), ctx)
			return false
		}
	}
	return true
}

// MARK: Query inspections

func GetQueryInspections(ctx *gin.Context) {
	isAccessGranted, error := GetAuthenticatedClient(ctx.Request)
	if error != nil || !isAccessGranted {
		ctx.IndentedJSON(http.StatusUnauthorized, models.Response{
			Status:  "fail",
			Message: "Access denied!",
		})
		return
	}

	var inspections []models.QueryInspection
	licensePlate := ctx.Param("license-plate")

	inspections = GetQueryInspectionsHelper(ctx, licensePlate)
	if inspections == nil {
		return
	}

	SendData(inspections, ctx)
}

// Returns all inspections for a given license plate
func GetQueryInspectionsHelper(ctx *gin.Context, licensePlate string) []models.QueryInspection {
	isAccessGranted, error := GetAuthenticatedClient(ctx.Request)
	if error != nil || !isAccessGranted {
		ctx.IndentedJSON(http.StatusUnauthorized, models.Response{
			Status:  "fail",
			Message: "Access denied!",
		})
		return nil
	}

	var inspections []models.QueryInspection

	result := DB.Find(&inspections, "car_id = ?", licensePlate)
	if result.Error != nil {
		SendError(result.Error.Error(), ctx)
		return nil
	} else if result.RowsAffected == 0 {
		return []models.QueryInspection{}
	}

	for i, inspection := range inspections {
		inspections[i].Base64 = ConvertImagesToBase64(inspection.ImageLocation)
	}

	return inspections
}

func CreateQueryInspections(ctx *gin.Context) {
	isAccessGranted, error := GetAuthenticatedClient(ctx.Request)
	if error != nil || !isAccessGranted {
		ctx.IndentedJSON(http.StatusUnauthorized, models.Response{
			Status:  "fail",
			Message: "Access denied!",
		})
		return
	}

	var newInspections []models.QueryInspection

	if err := ctx.BindJSON(&newInspections); err != nil {
		SendError(err.Error(), ctx)
		return
	}

	tx := DB.Begin()

	successful := CreateQueryInspectionHelper(ctx, newInspections, tx)

	if !successful {
		return
	}

	tx.Commit()

	SendData("Inspections were uploaded successfully", ctx)
	return
}

func CreateQueryInspectionHelper(ctx *gin.Context, newInspections []models.QueryInspection, tx *gorm.DB) bool {
	isAccessGranted, error := GetAuthenticatedClient(ctx.Request)
	if error != nil || !isAccessGranted {
		ctx.IndentedJSON(http.StatusUnauthorized, models.Response{
			Status:  "fail",
			Message: "Access denied!",
		})
		return false
	}

	for _, newInspection := range newInspections {
		checkResult := tx.Where("name = ?", newInspection.Name).First(&newInspection)
		if checkResult.RowsAffected != 0 {
			continue
		}

		result := tx.Create(&newInspection)
		if result.Error != nil {
			tx.Rollback()
			SendError(result.Error.Error(), ctx)
			return false
		}
	}
	return true
}

func DeleteQueryInspections(ctx *gin.Context) {
	isAccessGranted, error := GetAuthenticatedClient(ctx.Request)
	if error != nil || !isAccessGranted {
		ctx.IndentedJSON(http.StatusUnauthorized, models.Response{
			Status:  "fail",
			Message: "Access denied!",
		})
		return
	}

	licensePlate := ctx.Param("license-plate")

	success := DeleteQueryInspectionsHelper(ctx, licensePlate, false)

	if !success {
		return
	}

	SendData("Inspections were deleted successfully", ctx)
}

func DeleteQueryInspectionsHelper(ctx *gin.Context, licensePlate string, imagesOnly bool) bool {
	isAccessGranted, error := GetAuthenticatedClient(ctx.Request)
	if error != nil || !isAccessGranted {
		ctx.IndentedJSON(http.StatusUnauthorized, models.Response{
			Status:  "fail",
			Message: "Access denied!",
		})
		return false
	}

	isQuerySaved := ctx.Query("isQuerySaved")

	var inspections []models.QueryInspection

	if isQuerySaved == "true" {
		result := DB.Where("car_id = ?", licensePlate).Delete(&inspections)
		if result.Error != nil {
			SendError(result.Error.Error(), ctx)
			return false
		}
		return true
	}

	if imagesOnly {
		for _, inspection := range inspections {
			errorResult := os.RemoveAll(inspection.ImageLocation)
			if errorResult != nil {
				SendError(errorResult.Error(), ctx)
				return false
			}
		}
		return true
	}

	result := DB.Find(&inspections, "car_id = ?", licensePlate)
	if result.RowsAffected == 0 {
		return true
	}

	for _, inspection := range inspections {
		errorResult := os.RemoveAll(inspection.ImageLocation)
		if errorResult != nil {
			SendError(errorResult.Error(), ctx)
			return false
		}
	}

	result = DB.Where("car_id = ?", licensePlate).Delete(&inspections)
	if result.Error != nil {
		SendError(result.Error.Error(), ctx)
		return false
	}
	return true
}
