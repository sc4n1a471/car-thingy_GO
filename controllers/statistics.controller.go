package controllers

import (
	"Go_Thingy_GO/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetStatistics(ctx *gin.Context) {
	isAccessGranted, error := GetAuthenticatedClient(ctx.Request)
	if error != nil || !isAccessGranted {
		ctx.IndentedJSON(http.StatusUnauthorized, models.Response{
			Status:  "fail",
			Message: "Access denied!",
		})
		return
	}

	var statistics models.Statistics

	var carCount int64
	// Get total number of cars
	result := DB.Table("cars").Count(&carCount)
	if result.Error != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, models.Response{
			Status:  "fail",
			Message: "Failed to get car count!",
		})
		return
	}
	statistics.CarCount = int(carCount)

	var knownCars int64
	// Get total number of known cars
	result = DB.Table("cars").Where("brand IS NOT NULL").Count(&knownCars)
	if result.Error != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, models.Response{
			Status:  "fail",
			Message: "Failed to get known car count!",
		})
		return
	}
	statistics.KnownCars = int(knownCars)

	var unknownCars int64
	// Get total number of unknown cars
	result = DB.Table("cars").Where("brand IS NULL").Count(&unknownCars)
	if result.Error != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, models.Response{
			Status:  "fail",
			Message: "Failed to get unknown car count!",
		})
		return
	}
	statistics.UnknownCars = int(unknownCars)

	var brandStats []models.BrandStatistics
	// Get brand statistics
	rows, err := DB.Table("cars").Select("brand, COUNT(*) as count").Group("brand").Where("brand is not NULL").Rows()
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, models.Response{
			Status:  "fail",
			Message: "Failed to get brand statistics!",
		})
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
			ctx.IndentedJSON(http.StatusInternalServerError, models.Response{
				Status:  "fail",
				Message: "Failed to get model statistics!",
			})
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
