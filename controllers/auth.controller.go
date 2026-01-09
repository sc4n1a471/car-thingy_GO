package controllers

import (
	"Go_Thingy_GO/models"
	"crypto/rand"
	"log/slog"
	"math/big"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func CheckAuthKey(ctx *gin.Context) {
	isAccessGranted, error := GetAuthenticatedClient(ctx.Request)
	if error != nil || !isAccessGranted {
		SendError("Access denied for checking API key", http.StatusUnauthorized, ctx)
		return
	}
	ctx.IndentedJSON(http.StatusOK, models.Response{
		Status:  "success",
		Message: "Access granted!",
	})
}

// MARK: GetAuthenticatedClient
func GetAuthenticatedClient(r *http.Request) (bool, error) {
	var authKey models.AuthKey
	authKey.ID = r.Header.Get("x-api-key")
	authKey.IsActive = true

	if authKey.ID == "" {
		return false, nil
	}

	result := DB.First(&authKey)
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

// MARK: CreateAuthKeyWrapper
func CreateAuthKeyWrapper(ctx *gin.Context) {
	if ctx.Request.Header.Get("x-api-key") != os.Getenv("API_SECRET") {
		SendError("Access denied for creating new API key", http.StatusUnauthorized, ctx)
		return
	}

	var validAuthKey models.AuthKey
	result := DB.Select(&validAuthKey).Where("is_valid = ?", true)
	if result.Error != nil {
		SendError("Can't create multiple active API keys", http.StatusBadRequest, ctx)
		return
	}

	generatedKey, error := createAuthKey()
	if error != nil {
		SendError("Error creating auth key: "+error.Error(), http.StatusBadRequest, ctx)
		return
	}

	ctx.IndentedJSON(http.StatusCreated, models.Response{
		Status:  "success",
		Message: generatedKey,
	})
}

func createAuthKey() (string, error) {
	generatedKey, error := generateRandomString(64)
	if error != nil {
		slog.Error("Error generating random string: " + error.Error())
		return "", error
	}

	var newAuthKey models.AuthKey
	newAuthKey.ID = generatedKey
	newAuthKey.IsActive = true

	tx := DB.Begin()
	result := tx.Create(&newAuthKey)
	if result.Error != nil {
		return "", result.Error
	}
	tx.Commit()
	return generatedKey, nil
}

// https://gist.github.com/dopey/c69559607800d2f2f90b1b1ed4e550fb
func generateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}

// MARK: DeleteAuthKey
func DeleteAuthKey(ctx *gin.Context) {
	isAccessGranted, error := GetAuthenticatedClient(ctx.Request)
	if error != nil || !isAccessGranted {
		SendError("Access denied for deleting API key", http.StatusUnauthorized, ctx)
		return
	}
	var authKey models.AuthKey
	authKey.ID = ctx.Request.Header.Get("x-api-key")

	tx := DB.Begin()
	result := tx.Delete(&authKey)
	if result.Error != nil {
		SendError("Error deleting auth key: "+result.Error.Error(), http.StatusBadRequest, ctx)
		return
	}

	tx.Commit()
	ctx.IndentedJSON(http.StatusOK, models.Response{
		Status:  "success",
		Message: "AutKey deleted successfully!",
	})
}
