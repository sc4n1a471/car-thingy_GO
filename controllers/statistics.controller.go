package controllers

import (
	"Go_Thingy_GO/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetStatistics(ctx *gin.Context) {
	isAccessGranted, error := CheckAuthKey(ctx.Request)
	if error != nil || !isAccessGranted {
		SendError("Access denied for getting statistics", http.StatusUnauthorized, ctx)
		return
	}

	var statistics models.Statistics

	// MARK: Get total number of cars
	var carCount int64
	result := DB.Table("cars").Count(&carCount)
	if result.Error != nil {
		SendError("Failed to get car count: "+result.Error.Error(), http.StatusInternalServerError, ctx)
		return
	}
	statistics.CarCount = int(carCount)

	// MARK: Get known and unknown cars
	var knownCars int64
	result = DB.Table("cars").Where("brand IS NOT NULL").Count(&knownCars)
	if result.Error != nil {
		SendError("Failed to get known car count: "+result.Error.Error(), http.StatusInternalServerError, ctx)
		return
	}
	statistics.KnownCars = int(knownCars)

	var unknownCars int64
	result = DB.Table("cars").Where("brand IS NULL").Count(&unknownCars)
	if result.Error != nil {
		SendError("Failed to get unknown car count: "+result.Error.Error(), http.StatusInternalServerError, ctx)
		return
	}
	statistics.UnknownCars = int(unknownCars)

	// MARK: Get brand and model statistics
	var brandStats []models.BrandStatistics
	rows, err := DB.Table("cars").Select("brand, COUNT(*) as count").Group("brand").Where("brand is not NULL").Rows()
	if err != nil {
		SendError("Failed to get brand statistics: "+err.Error(), http.StatusInternalServerError, ctx)
		return
	}
	for rows.Next() {
		var brandStat models.BrandStatistics
		DB.ScanRows(rows, &brandStat)
		brandStats = append(brandStats, brandStat)
	}
	statistics.BrandStats = brandStats

	for i, brandStat := range brandStats {
		var modelStats []models.ModelStatistics
		// Get model statistics
		rows, err := DB.Table("cars").Select("model, COUNT(*) as count").Where("brand = ?", brandStat.Brand).Group("model").Rows()
		if err != nil {
			SendError("Failed to get model statistics: "+err.Error(), http.StatusInternalServerError, ctx)
			return
		}

		for rows.Next() {
			var modelStat models.ModelStatistics
			DB.ScanRows(rows, &modelStat)
			modelStats = append(modelStats, modelStat)
		}
		statistics.BrandStats[i].Models = modelStats
	}

	SendData(statistics, ctx)
}
