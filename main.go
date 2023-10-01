package main

import (
	"Go_Thingy/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

var DB *gorm.DB
var Error error

//var queries = []models.Specs{
//	{
//		LicensePlate: "TEST111",
//		Brand:        "BMW",
//		Color:        "FEKETE",
//		EngineSize:   4400,
//		FirstReg:     "1994.01.01",
//		FirstRegHun:  "1994.01.01",
//		FuelType:     "Benzin",
//		Gearbox:      "Automata",
//		Model:        "740i",
//		NumOfOwners:  3,
//		Performance:  300,
//		Status:       "Forgalomban van",
//		TypeCode:     "E38",
//		Year:         1994,
//		Comment:      "yas",
//	},
//	{
//		LicensePlate: "TEST112",
//		Brand:        "BMW",
//		Color:        "FEKETE",
//		EngineSize:   4400,
//		FirstReg:     "1994.01.01",
//		FirstRegHun:  "1994.01.01",
//		FuelType:     "Benzin",
//		Gearbox:      "Automata",
//		Model:        "740i",
//		NumOfOwners:  3,
//		Performance:  300,
//		Status:       "Forgalomban van",
//		TypeCode:     "E38",
//		Year:         1994,
//		Comment:      "yas",
//	},
//}

func getQuery(ctx *gin.Context) {
	//licensePlate := ctx.Param("license_plate")
	//fmt.Println(licensePlate)
	//
	//for _, query := range queries {
	//	if strings.ToLower(query.LicensePlate) == strings.ToLower(licensePlate) {
	//		ctx.IndentedJSON(http.StatusOK, query)
	//		return
	//	}
	//}
	//ctx.IndentedJSON(http.StatusNotFound, models.Response{
	//	Status:  "fail",
	//	Message: "not found",
	//})
}

func getQueries(ctx *gin.Context) {
	var allSpecs []models.Specs

	var returnQueries []models.Query

	result := DB.Find(&allSpecs)
	if result.Error != nil {
		ctx.IndentedJSON(http.StatusConflict, models.Response{
			Status:  "fail",
			Message: result.Error.Error(),
		})
		return
	}

	for _, specs := range allSpecs {
		var query models.Query
		var accidents []models.Accident
		var restrictions []models.Restriction
		var mileages []models.Mileage

		query.Specs = specs

		result := DB.Find(&accidents, "license_plate = ?", specs.LicensePlate)
		if result.Error != nil {
			ctx.IndentedJSON(http.StatusConflict, models.Response{
				Status:  "fail",
				Message: result.Error.Error(),
			})
			return
		}
		query.Accidents = accidents

		result = DB.Find(&restrictions, "license_plate = ?", specs.LicensePlate)
		if result.Error != nil {
			ctx.IndentedJSON(http.StatusConflict, models.Response{
				Status:  "fail",
				Message: result.Error.Error(),
			})
			return
		}
		query.Restrictions = restrictions

		result = DB.Find(&mileages, "license_plate = ?", specs.LicensePlate)
		if result.Error != nil {
			ctx.IndentedJSON(http.StatusConflict, models.Response{
				Status:  "fail",
				Message: result.Error.Error(),
			})
			return
		}
		query.Mileage = mileages

		returnQueries = append(returnQueries, query)
	}

	ctx.IndentedJSON(http.StatusOK, returnQueries)
}

func createQuery(ctx *gin.Context) {
	var newSpecs models.Specs
	var newAccidents []models.Accident
	var newRestrictions []models.Restriction
	var newMileages []models.Mileage
	var newQuery models.Query

	if err := ctx.BindJSON(&newQuery); err != nil {
		ctx.IndentedJSON(http.StatusConflict, models.Response{
			Status:  "fail",
			Message: err.Error(),
		})
		return
	}

	newSpecs = newQuery.Specs
	newAccidents = newQuery.Accidents
	newRestrictions = newQuery.Restrictions
	newMileages = newQuery.Mileage

	tx := DB.Begin()
	result := DB.Create(&newSpecs)
	if result.Error != nil {
		tx.Rollback()
		ctx.IndentedJSON(http.StatusConflict, models.Response{
			Status:  "fail",
			Message: result.Error.Error(),
		})
		return
	}

	for _, newAccident := range newAccidents {

		var existingAccident models.Accident
		checkResult := DB.Where(&models.Accident{
			LicensePlate: newAccident.LicensePlate,
			AccidentDate: newAccident.AccidentDate,
		}).Find(&existingAccident)
		if checkResult.RowsAffected != 0 {
			continue
		}

		result := DB.Create(&newAccident)
		if result.Error != nil {
			tx.Rollback()
			ctx.IndentedJSON(http.StatusConflict, models.Response{
				Status:  "fail",
				Message: result.Error.Error(),
			})
			return
		}
	}

	var existingRestrictions []models.Restriction
	result = DB.Find(&existingRestrictions, "license_plate = ?", newSpecs.LicensePlate)
	if result.Error != nil {
		ctx.IndentedJSON(http.StatusConflict, models.Response{
			Status:  "fail",
			Message: result.Error.Error(),
		})
		return
	}
	for _, existingRestriction := range existingRestrictions {
		for _, newRestriction := range newRestrictions {
			if existingRestriction.Restriction == newRestriction.Restriction {
				continue
			}
		}
		// set existing to inactive
	}
	for _, newRestriction := range newRestrictions {
		for _, existingRestriction := range existingRestrictions {
			if existingRestriction.Restriction == newRestriction.Restriction {
				continue
			}
		}

		result := DB.Create(&newRestriction)
		if result.Error != nil {
			tx.Rollback()
			ctx.IndentedJSON(http.StatusConflict, models.Response{
				Status:  "fail",
				Message: result.Error.Error(),
			})
			return
		}
	}

	for _, newMileage := range newMileages {

		var tempMileage models.Mileage
		checkResult := DB.Where(&models.Mileage{
			LicensePlate: newMileage.LicensePlate,
			Mileage:      newMileage.Mileage,
		}).Find(&tempMileage)
		if checkResult.RowsAffected != 0 {
			continue
		}

		result := DB.Create(&newMileage)
		if result.Error != nil {
			tx.Rollback()
			ctx.IndentedJSON(http.StatusConflict, models.Response{
				Status:  "fail",
				Message: result.Error.Error(),
			})
			return
		}
	}

	tx.Commit()
	ctx.IndentedJSON(http.StatusCreated, models.Response{
		Status:  "success",
		Message: "Specs was uploaded successfully",
	})
	return
}

func deleteQuery(ctx *gin.Context) {
	var deletableSpecs models.Specs

	deletableSpecs.LicensePlate = ctx.Param("license_plate")

	result := DB.Where("license_plate = ?", deletableSpecs.LicensePlate).Delete(&deletableSpecs)

	if result.RowsAffected == 0 {
		ctx.IndentedJSON(http.StatusConflict, models.Response{
			Status:  "fail",
			Message: result.Error.Error(),
		})
		return
	}

	ctx.IndentedJSON(http.StatusCreated, models.Response{
		Status:  "success",
		Message: "Query was deleted successfully",
	})
}

func main() {
	err := godotenv.Load("env.env")
	if err != nil {
		log.Fatalf("Error loading env.env file")
	}
	dsn := os.Getenv("DB_USERNAME") +
		":" +
		os.Getenv("DB_PASSWORD") +
		"@tcp(" +
		os.Getenv("DB_IP") +
		":" +
		os.Getenv("DB_PORT") +
		")/" +
		os.Getenv("DB_NAME") +
		"?parseTime=true"

	DB, Error = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if Error != nil {
		return
	}

	router := gin.Default()
	//router.GET("/queries/:license_plate", getQuery)
	router.GET("/queries", getQueries)
	router.POST("/queries", createQuery)
	router.DELETE("/queries/:license_plate", deleteQuery)

	router.Run("localhost:3000")
}
