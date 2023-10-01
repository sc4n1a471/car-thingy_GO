package main

import (
	"Go_Thingy/models"
	"fmt"
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

//var cars = []models.Specs{
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

func getCar(ctx *gin.Context) {
	var requested models.Specs

	requested.LicensePlate = ctx.Param("license_plate")
	result := DB.First(&requested)
	if result.Error != nil {
		ctx.IndentedJSON(http.StatusConflict, models.Response{
			Status:  "fail",
			Message: result.Error.Error(),
		})
		return
	}

	var car models.Car
	var accidents []models.Accident
	var restrictions []models.Restriction
	var mileages []models.Mileage

	car.Specs = requested

	accidentResults := DB.Find(&accidents, "license_plate = ?", car.Specs.LicensePlate)
	if accidentResults.Error != nil {
		ctx.IndentedJSON(http.StatusConflict, models.Response{
			Status:  "fail",
			Message: accidentResults.Error.Error(),
		})
		return
	}
	car.Accidents = accidents

	restrictionResult := DB.Find(&restrictions, "license_plate = ? AND active = true", car.Specs.LicensePlate)
	if restrictionResult.Error != nil {
		ctx.IndentedJSON(http.StatusConflict, models.Response{
			Status:  "fail",
			Message: restrictionResult.Error.Error(),
		})
		return
	}
	car.Restrictions = restrictions

	mileageResult := DB.Find(&mileages, "license_plate = ?", car.Specs.LicensePlate)
	if mileageResult.Error != nil {
		ctx.IndentedJSON(http.StatusConflict, models.Response{
			Status:  "fail",
			Message: mileageResult.Error.Error(),
		})
		return
	}
	car.Mileage = mileages

	ctx.IndentedJSON(http.StatusOK, car)
}

func getCars(ctx *gin.Context) {
	var allSpecs []models.Specs

	var returnCars []models.Car

	result := DB.Find(&allSpecs)
	if result.Error != nil {
		ctx.IndentedJSON(http.StatusConflict, models.Response{
			Status:  "fail",
			Message: result.Error.Error(),
		})
		return
	}

	for _, specs := range allSpecs {
		var car models.Car
		var accidents []models.Accident
		var restrictions []models.Restriction
		var mileages []models.Mileage

		car.Specs = specs

		result := DB.Find(&accidents, "license_plate = ?", specs.LicensePlate)
		if result.Error != nil {
			ctx.IndentedJSON(http.StatusConflict, models.Response{
				Status:  "fail",
				Message: result.Error.Error(),
			})
			return
		}
		car.Accidents = accidents

		result = DB.Find(&restrictions, "license_plate = ?", specs.LicensePlate)
		if result.Error != nil {
			ctx.IndentedJSON(http.StatusConflict, models.Response{
				Status:  "fail",
				Message: result.Error.Error(),
			})
			return
		}
		car.Restrictions = restrictions

		result = DB.Find(&mileages, "license_plate = ?", specs.LicensePlate)
		if result.Error != nil {
			ctx.IndentedJSON(http.StatusConflict, models.Response{
				Status:  "fail",
				Message: result.Error.Error(),
			})
			return
		}
		car.Mileage = mileages

		returnCars = append(returnCars, car)
	}

	ctx.IndentedJSON(http.StatusOK, returnCars)
}

func createCar(ctx *gin.Context) {
	var newSpecs models.Specs
	var newAccidents []models.Accident
	var newRestrictions []models.Restriction
	var newMileages []models.Mileage
	var newCar models.Car

	if err := ctx.BindJSON(&newCar); err != nil {
		ctx.IndentedJSON(http.StatusConflict, models.Response{
			Status:  "fail",
			Message: err.Error(),
		})
		return
	}

	newSpecs = newCar.Specs
	newAccidents = newCar.Accidents
	newRestrictions = newCar.Restrictions
	newMileages = newCar.Mileage

	tx := DB.Begin()
	result := DB.First(&newSpecs)
	if result.RowsAffected == 0 {
		result := DB.Create(&newSpecs)
		if result.Error != nil {
			tx.Rollback()
			ctx.IndentedJSON(http.StatusConflict, models.Response{
				Status:  "fail",
				Message: result.Error.Error(),
			})
			return
		}
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

existingsLoop:
	for _, existingRestriction := range existingRestrictions {
		for _, newRestriction := range newRestrictions {
			if existingRestriction.Restriction == newRestriction.Restriction {
				fmt.Println(existingRestriction.Restriction)
				fmt.Println(newRestriction.Restriction)
				fmt.Println(existingRestriction.Restriction == newRestriction.Restriction)
				continue existingsLoop
			}
		}
		fmt.Println(existingRestriction)
		DB.Model(&models.Restriction{}).
			Where(
				"license_plate = ? AND restriction = ?",
				existingRestriction.LicensePlate,
				existingRestriction.Restriction).
			Update("active", false)
	}

newsLoop:
	for _, newRestriction := range newRestrictions {
		for _, existingRestriction := range existingRestrictions {
			if existingRestriction.Restriction == newRestriction.Restriction {
				continue newsLoop
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
		Message: "Car was uploaded successfully",
	})
	return
}

func deleteCar(ctx *gin.Context) {
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
		Message: "Car was deleted successfully",
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
	router.GET("/cars/:license_plate", getCar)
	router.GET("/cars", getCars)
	router.POST("/cars", createCar)
	router.DELETE("/cars/:license_plate", deleteCar)

	router.Run("localhost:3000")
}
