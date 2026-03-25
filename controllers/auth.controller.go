package controllers

import (
	"Go_Thingy_GO/models"
	"crypto/rand"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// MARK: CheckAuthKeyWrapper
func CheckAuthKeyWrapper(ctx *gin.Context) {
	isAccessGranted, error := CheckAuthKey(ctx.Request)
	if error != nil || !isAccessGranted {
		SendError("Access denied for checking API key", http.StatusUnauthorized, ctx)
		return
	}
	ctx.IndentedJSON(http.StatusOK, models.Response{
		Status:  "success",
		Message: "Access granted!",
	})
}

// MARK: CheckAuthKey
// Returns true if the auth key is valid and active
func CheckAuthKey(r *http.Request) (bool, error) {
	var authKey models.AuthKey
	authKey.ID = r.Header.Get("x-api-key")
	authKey.IsActive = true

	slog.Info("Check API key from host: " + r.RemoteAddr)

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
	isAccessGranted, error := CheckAuthKey(ctx.Request)
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

// MARK: Log query timestamp
func LogQueryTimestamp(ctx *gin.Context) {
	isAccessGranted, error := CheckAuthKey(ctx.Request)
	if error != nil || !isAccessGranted {
		SendError("Access denied for logging query timestamp", http.StatusUnauthorized, ctx)
		return
	}

	var queryLog models.QueryLog

	if err := ctx.BindJSON(&queryLog); err != nil {
		SendError("Could not parse JSON: "+err.Error(), http.StatusBadRequest, ctx)
		return
	}

	queryLog.User = ctx.Request.Header.Get("x-api-key")
	queryLog.QueryTimestamp = time.Now()

	tx := DB.Begin()
	result := tx.Create(&queryLog)
	if result.Error != nil {
		SendError("Error logging query timestamp: "+result.Error.Error(), http.StatusBadRequest, ctx)
		return
	}

	tx.Commit()
	SendData("Query timestamp logged successfully", ctx)
}

// MARK: Get last query timestamp
func GetLastLogQueryTimestamp(ctx *gin.Context) {
	isAccessGranted, error := CheckAuthKey(ctx.Request)
	if error != nil || !isAccessGranted {
		SendError("Access denied for getting last query timestamp", http.StatusUnauthorized, ctx)
		return
	}

	var queryLog models.QueryLog
	result := DB.Where("user = ?", ctx.Request.Header.Get("x-api-key")).Order("query_timestamp desc").First(&queryLog)
	if result.RowsAffected == 0 {
		SendData(0, ctx)
		return
	}
	if result.Error != nil {
		SendError("Error getting last query timestamp: "+result.Error.Error(), http.StatusBadRequest, ctx)
		return
	}

	// Calculate how many seconds ago the queryTimeStamp was and return that in the response
	var secondsAgo = int(time.Since(queryLog.QueryTimestamp).Seconds())
	var waitingTime = 0
	if secondsAgo >= 30 {
		waitingTime = 0
	} else {
		waitingTime = 30 - secondsAgo
		waitingTime = int(waitingTime) + 1
	}
	slog.Info("Previous query timestamp: " + queryLog.QueryTimestamp.String() + ", how many seconds ago: " + fmt.Sprint(secondsAgo) + ", waiting time: " + fmt.Sprint(waitingTime))
	SendData(waitingTime, ctx)
}
