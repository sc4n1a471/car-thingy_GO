package controllers

import (
	"Go_Thingy_GO/models"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var Error error

func SetupDatabase() error {
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
		return Error
	}

	err := DB.AutoMigrate(
		&models.Car{},
		&models.Accident{},
		&models.Inspection{},
		&models.QueryInspection{},
		&models.Mileage{},
		&models.Restriction{},
	)
	if err != nil {
		fmt.Print(err.Error())
		return err
	}

	return nil
}

// Converts images in *imageLocation* directory to base64 format
// Returns an array of string containing the base64 images
// https://freshman.tech/snippets/go/image-to-base64/
func ConvertImagesToBase64(imageLocation string) []string {
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

func SendError(error string, ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusConflict, models.Response{
		Status:  "fail",
		Message: error,
	})
}

func SendData(message interface{}, ctx *gin.Context) {
	response := map[string]interface{}{
		"status":  "success",
		"message": message,
	}

	ctx.IndentedJSON(http.StatusOK, response)
}
