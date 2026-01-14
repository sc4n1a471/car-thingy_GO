package controllers

import (
	"Go_Thingy_GO/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// MARK: CreateLicensePlate
func CreateLicensePlate(ctx *gin.Context) {
	isAccessGranted, error := CheckAuthKey(ctx.Request)
	if error != nil || !isAccessGranted {
		SendError("Access denied for creating license plate", http.StatusUnauthorized, ctx)
		return
	}

	var newCar models.Car

	if err := ctx.BindJSON(&newCar); err != nil {
		SendError("Could not parse JSON: "+err.Error(), http.StatusBadRequest, ctx)
		return
	}

	newCar.ID = strings.ReplaceAll(newCar.ID, " ", "")

	tx := DB.Begin()
	result := tx.First(&newCar)
	if result.RowsAffected == 0 {
		result := tx.Create(&newCar)
		if result.Error != nil {
			tx.Rollback()
			SendError("Error creating license plate: "+result.Error.Error(), http.StatusInternalServerError, ctx)
			return
		}
	}

	tx.Commit()

	SendData("License plate was uploaded successfully", ctx)
	return
}

// MARK: UpdateLicensePlate
func UpdateLicensePlate(ctx *gin.Context) {
	isAccessGranted, error := CheckAuthKey(ctx.Request)
	if error != nil || !isAccessGranted {
		SendError("Access denied for updating license plate", http.StatusUnauthorized, ctx)
		return
	}

	var updatedCar models.Car

	if err := ctx.BindJSON(&updatedCar); err != nil {
		SendError("Could not parse JSON: "+err.Error(), http.StatusBadRequest, ctx)
		return
	}

	var oldLicensePlate = ctx.Param("license-plate")

	updatedCar.ID = strings.ReplaceAll(updatedCar.ID, " ", "")
	oldLicensePlate = strings.ReplaceAll(oldLicensePlate, " ", "")

	tx := DB.Begin()

	// Update license plate (oldLicensePlate) with new license plate (updatedCar.LicensePlate)
	result := tx.
		Model(&models.Car{}).
		Where("id = ?", oldLicensePlate).
		Update("id", updatedCar.ID)
	if result.Error != nil {
		tx.Rollback()
		SendError("Could not update license plate: "+result.Error.Error(), http.StatusInternalServerError, ctx)
		return
	}

	// Update inspections with new license plate
	result = tx.
		Model(&models.Inspection{}).
		Where("car_id = ?", oldLicensePlate).
		Update("car_id", updatedCar.ID)
	if result.Error != nil {
		tx.Rollback()
		SendError("Could not update inspections: "+result.Error.Error(), http.StatusInternalServerError, ctx)
		return
	}

	tx.Commit()

	SendData("License Plate was updated successfully", ctx)
}
