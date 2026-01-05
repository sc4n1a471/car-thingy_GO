package controllers

import (
	"Go_Thingy_GO/models"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MARK: Normal inspections
func GetInspectionsHelper(ctx *gin.Context) {
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

	inspections = GetInspections(ctx, licensePlate)
	if inspections == nil {
		return
	}

	SendData(inspections, ctx)
}

// Returns all inspections for a given license plate
func GetInspections(ctx *gin.Context, licensePlate string) []models.Inspection {
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
		checkResult := tx.Where("name = ? and car_id = ?", newInspection.Name, newInspection.CarID).Find(&newInspection)
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

func GetQueryInspectionsHelper(ctx *gin.Context) {
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

	inspections = GetQueryInspections(ctx, licensePlate)
	if inspections == nil {
		return
	}

	SendData(inspections, ctx)
}

// Returns all inspections for a given license plate
func GetQueryInspections(ctx *gin.Context, licensePlate string) []models.QueryInspection {
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

func CreateQueryInspectionsHelper(ctx *gin.Context) {
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

	successful := CreateQueryInspection(ctx, newInspections, tx)

	if !successful {
		return
	}

	tx.Commit()

	SendData("Inspections were uploaded successfully", ctx)
}

func CreateQueryInspection(ctx *gin.Context, newInspections []models.QueryInspection, tx *gorm.DB) bool {
	isAccessGranted, error := GetAuthenticatedClient(ctx.Request)
	if error != nil || !isAccessGranted {
		ctx.IndentedJSON(http.StatusUnauthorized, models.Response{
			Status:  "fail",
			Message: "Access denied!",
		})
		return false
	}
	fmt.Println("Creating ", len(newInspections), " query inspections")

	var successfulCreations int = 0

	for _, newInspection := range newInspections {
		fmt.Println("Creating query inspection: ", newInspection.Name)
		checkResult := tx.Where("name = ? and car_id = ?", newInspection.Name, newInspection.CarID).First(&newInspection)
		if checkResult.RowsAffected != 0 {
			fmt.Println("Query inspection already exists: ", newInspection.Name, " for car ", newInspection.CarID)
			continue
		}

		result := tx.Create(&newInspection)
		if result.Error != nil {
			tx.Rollback()
			SendError(result.Error.Error(), ctx)
			return false
		}
		successfulCreations += 1
	}
	fmt.Println("Successfully created ", successfulCreations, " query inspections")
	return true
}

// MARK: Delete query inspections and their images
// Deletes all query inspections and their images for a given license plate
func DeleteQueryInspections(ctx *gin.Context, licensePlate string, imagesOnly bool) bool {
	isAccessGranted, error := GetAuthenticatedClient(ctx.Request)
	if error != nil || !isAccessGranted {
		ctx.IndentedJSON(http.StatusUnauthorized, models.Response{
			Status:  "fail",
			Message: "Access denied!",
		})
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

// Delete all older query inspections and their images that were not saved
// MARK: Cleanup function
func DeleteOldQueryInspections() bool {
	fmt.Println("Deleting old queries...")
	var inspections []models.QueryInspection
	var deletedSuccessfully int64 = 0

	// SELECT car_id FROM `query_inspections` qi where (select count(*) from inspections where car_id = qi.car_id) = 0 group by car_id;

	var queryInspectionCarIds []string
	result := DB.Table("query_inspections").Select("car_id").Where("car_id NOT IN (SELECT car_id FROM inspections GROUP BY car_id)").Group("car_id").Scan(&queryInspectionCarIds)
	if result.Error != nil {
		fmt.Println("Error fetching old query inspections: ", result.Error.Error())
		return false
	}

	for _, id := range queryInspectionCarIds {
		fmt.Println("To be deleted query: ", id)
		result = DB.Find(&inspections, "car_id = ?", id)
		if result.RowsAffected == 0 {
			return true
		}
		if result.Error != nil {
			fmt.Println("Error fetching old query inspections: ", result.Error.Error())
			return false
		}

		for _, inspection := range inspections {
			errorResult := os.RemoveAll(inspection.ImageLocation)
			if errorResult != nil {
				fmt.Println("Error deleting old query inspection images: ", errorResult.Error())
				return false
			}
		}

		result = DB.Where("car_id = ?", id).Delete(&inspections)
		if result.Error != nil {
			fmt.Println("Error deleting old query inspections: ", result.Error.Error())
			return false
		}
		deletedSuccessfully += 1
	}
	fmt.Println("Deleted ", deletedSuccessfully, " query inspections")
	return true
}
