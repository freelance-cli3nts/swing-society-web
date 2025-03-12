package storage

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// TestFirebaseSchemaHandler is a handler for testing the Firebase schema implementation
func TestFirebaseSchemaHandler(w http.ResponseWriter, r *http.Request) {
	// Initialize Firebase client
	firebase, err := NewFirebaseClient()
	if err != nil {
		http.Error(w, "Failed to initialize Firebase: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Test user creation
	testUserId := "test_" + fmt.Sprintf("%d", time.Now().Unix())
	testUser := map[string]interface{}{
		"profile": map[string]interface{}{
			"name":          "Test User",
			"email":         "test@example.com",
			"phone":         "+1234567890",
			"createdAt":     time.Now().Unix(),
			"updatedAt":     time.Now().Unix(),
			"termsAccepted": true,
			"gdprAccepted":  true,
		},
		"auth": map[string]interface{}{
			"hasAccount": false,
		},
		"subscriptions": map[string]interface{}{
			"newsletter": map[string]interface{}{
				"subscribed": true,
				"frequency":  "weekly",
			},
		},
	}
	
	if err := firebase.SaveUser(testUserId, testUser); err != nil {
		http.Error(w, "Failed to test user creation: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Test email index
	if err := firebase.UpdateEmailIndex("test@example.com", testUserId); err != nil {
		http.Error(w, "Failed to update email index: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Test phone index
	if err := firebase.UpdatePhoneIndex("+1234567890", testUserId); err != nil {
		http.Error(w, "Failed to update phone index: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Test submission creation
	submissionData := map[string]interface{}{
		"timestamp":     time.Now().Unix(),
		"email":         "test@example.com",
		"name":          "Test User",
		"phone":         "+1234567890",
		"role":          "follower",
		"signupMethod":  "alone",
		"termsAccepted": true,
		"gdprAccepted":  true,
		"userId":        testUserId,
	}
	
	submissionId, err := firebase.SaveForm("registrations", submissionData)
	if err != nil {
		http.Error(w, "Failed to save test submission: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Write success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Firebase schema verified",
		"data": map[string]interface{}{
			"testUserId":       testUserId,
			"testSubmissionId": submissionId,
		},
	})
}