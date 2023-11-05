package application

import (
	"Go_Thingy/models"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func getCar(ctx *gin.Context) {
	var requested models.LicensePlate
	var car models.CarResult
	var coordinates models.Coordinate
	var returnData []models.CarResult

	requested.LicensePlate = ctx.Param("license_plate")
	result := DB.First(&requested)
	if result.Error != nil {
		sendError(result.Error.Error(), ctx)
		return
	}

	car.LicensePlate = requested

	car.Specs = getSpecs(ctx, requested.LicensePlate)
	//if car.Specs == (models.Specs{}) {
	//	sendError("Specs is null?", ctx)
	//	return
	//}

	car.Accidents = getAccidents(ctx, requested.LicensePlate)
	if car.Accidents == nil {
		return
	}

	car.Restrictions = getRestrictions(ctx, requested.LicensePlate)
	if car.Restrictions == nil {
		return
	}

	car.Mileage = getMileages(ctx, requested.LicensePlate)
	if car.Mileage == nil {
		return
	}

	result = DB.Find(&coordinates, "license_plate = ?", requested.LicensePlate)
	if result.Error != nil {
		sendError(result.Error.Error(), ctx)
		return
	}
	car.Coordinates = coordinates

	car.Inspections = getInspectionsHelper(ctx, requested.LicensePlate)
	if car.Inspections == nil {
		return
	}

	returnData = append(returnData, car)
	sendData(returnData, ctx)
}

func getCars(ctx *gin.Context) {
	var allLicensePlates []models.LicensePlate

	var returnCars []models.CarResult

	result := DB.Find(&allLicensePlates)
	if result.Error != nil {
		sendError(result.Error.Error(), ctx)
		return
	}

	for _, licensePlate := range allLicensePlates {
		var car models.CarResult
		var coordinates models.Coordinate

		car.LicensePlate = licensePlate

		car.Specs = getSpecs(ctx, licensePlate.LicensePlate)
		//if car.Specs == (models.Specs{}) {
		//	sendError("Specs is null?", ctx)
		//	return
		//}

		car.Accidents = getAccidents(ctx, licensePlate.LicensePlate)
		if car.Accidents == nil {
			return
		}

		car.Restrictions = getRestrictions(ctx, licensePlate.LicensePlate)
		if car.Restrictions == nil {
			return
		}

		car.Mileage = getMileages(ctx, licensePlate.LicensePlate)
		if car.Mileage == nil {
			return
		}

		result := DB.Find(&coordinates, "license_plate = ?", licensePlate.LicensePlate)
		if result.Error != nil {
			sendError(result.Error.Error(), ctx)
			return
		}
		car.Coordinates = coordinates

		// Temporarily disabled
		car.Inspections = []models.InspectionResult{}
		//car.Inspections = getInspectionsHelper(ctx, car.Specs.LicensePlate)
		//if car.Inspections == nil {
		//	return
		//}

		returnCars = append(returnCars, car)
	}

	if returnCars == nil {
		returnCars = []models.CarResult{}
	}

	sendData(returnCars, ctx)
}

func getInspectionsHelper(ctx *gin.Context, licensePlate string) []models.InspectionResult {
	var inspections []models.Inspection
	var inspectionResults []models.InspectionResult

	result := DB.Find(&inspections, "license_plate = ?", licensePlate)
	if result.Error != nil {
		sendError(result.Error.Error(), ctx)
		return nil
	} else if result.RowsAffected == 0 {
		return []models.InspectionResult{}
	}

	for _, inspection := range inspections {
		var inspectionResult models.InspectionResult
		inspectionResult.LicensePlate = inspection.LicensePlate
		inspectionResult.Name = inspection.Name
		inspectionResult.Base64 = convertImagesToBase64(inspection.ImageLocation)
		inspectionResults = append(inspectionResults, inspectionResult)
	}

	return inspectionResults
}

// Converts images in *imageLocation* directory to base64 format
// Returns an array of string containing the base64 images
// https://freshman.tech/snippets/go/image-to-base64/
func convertImagesToBase64(imageLocation string) []string {
	var convertedImages []string

	inspectionLocation := imageLocation
	files, err := os.ReadDir(inspectionLocation)
	if err != nil {
		log.Println(err)
		return nil
	}

	for _, file := range files {
		bytes, err := os.ReadFile(inspectionLocation + file.Name())
		if err != nil {
			log.Println(err)
			return nil
		}

		var base64Encoding string

		// Determine the content type of the image file
		//mimeType := http.DetectContentType(bytes)

		// Prepend the appropriate URI scheme header depending
		// on the MIME type
		//switch mimeType {
		//case "image/jpeg":
		//	base64Encoding += "data:image/jpeg;base64,"
		//case "image/png":
		//	base64Encoding += "data:image/png;base64,"
		//case "image/jpg":
		//	base64Encoding += "data:image/jpg;base64,"
		//}

		// Append the base64 encoded output
		base64Encoding += base64.StdEncoding.EncodeToString(bytes)

		// Print the full base64 representation of the image
		convertedImages = append(convertedImages, base64Encoding)
	}
	return convertedImages
}

func getInspections(ctx *gin.Context) {
	var inspectionResults []models.InspectionResult
	licensePlate := ctx.Param("license_plate")

	inspectionResults = getInspectionsHelper(ctx, licensePlate)
	if inspectionResults == nil {
		return
	}

	sendData(inspectionResults, ctx)
}

func getAccidents(ctx *gin.Context, licensePlate string) []models.Accident {
	var accidents []models.Accident
	result := DB.Find(&accidents, "license_plate = ?", licensePlate)
	if result.Error != nil {
		sendError(result.Error.Error(), ctx)
		return nil
	}
	return accidents
}

func getRestrictions(ctx *gin.Context, licensePlate string) []models.Restriction {
	var restrictions []models.Restriction
	result := DB.Find(&restrictions, "license_plate = ? AND active = true", licensePlate)
	if result.Error != nil {
		sendError(result.Error.Error(), ctx)
		return nil
	}
	return restrictions
}

func getMileages(ctx *gin.Context, licensePlate string) []models.Mileage {
	var mileages []models.Mileage
	result := DB.Find(&mileages, "license_plate = ?", licensePlate)
	if result.Error != nil {
		sendError(result.Error.Error(), ctx)
		return nil
	}
	return mileages
}

func getSpecs(ctx *gin.Context, licensePlate string) models.Specs {
	var specs models.Specs

	result := DB.Find(&specs, "license_plate = ?", licensePlate)
	if result.Error != nil {
		return models.Specs{}
	}

	return specs
}
