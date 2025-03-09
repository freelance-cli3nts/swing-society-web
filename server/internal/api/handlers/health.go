// internal/api/handlers/health.go
package handlers

import (
	"context"
	"net/http"
	"swing-society-website/server/internal/api/response"
	"swing-society-website/server/internal/config"
	"swing-society-website/server/internal/storage"
	"time"
)

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	status := "ok"
	dbStatus := "not_checked"
	
	// Check if Firebase test is requested
	if r.URL.Query().Get("check_firebase") == "true" {
		// Create a context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		// Initialize Firebase client
		firebaseClient, err := storage.NewFirebaseClient()
		if err != nil {
			dbStatus = "error_initializing: " + err.Error()
		} else {
			// Test connection by reading a simple path
			testData, err := firebaseClient.TestConnection(ctx)
			if err != nil {
				dbStatus = "error_connecting: " + err.Error()
			} else {
				dbStatus = "connected"
				if testData != nil {
					dbStatus += " (data available)"
				}
			}
		}
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"status":      status,
		"environment": config.AppConfig.Environment,
		"firebase":    dbStatus,
		"timestamp":   time.Now().Format(time.RFC3339),
	})
}