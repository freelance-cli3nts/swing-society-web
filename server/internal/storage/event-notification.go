package storage

import (
	"context"
	"fmt"
	"log"
	"time"
	"swing-society-website/server/internal/api/models"
)

// EventNotificationStorage defines the interface for event notification storage
type EventNotificationStorage interface {
	StoreEventNotification(notification *models.EventNotification) error
	IsSubscribed(email string) (bool, error)
	Unsubscribe(email string) error
}

type SimpleEventNotificationStorage struct {
	firebase *FirebaseClient
}

func NewSimpleEventNotificationStorage(firebase *FirebaseClient) *SimpleEventNotificationStorage {
	return &SimpleEventNotificationStorage{
		firebase: firebase,
	}
}

func (s *SimpleEventNotificationStorage) StoreEventNotification(notification *models.EventNotification) error {
	if s.firebase != nil {
		// Check if user already exists by email
		userId, err := s.firebase.GetUserByEmail(notification.Email)
		if err != nil {
			log.Printf("Error checking user by email for event notification: %v", err)
		}

		now := time.Now().Unix()
		firstName := notification.FirstName
		if firstName == "" {
			firstName = "Subscriber"
		}

		// Create submission data
		submissionData := map[string]interface{}{
			"timestamp":         now,
			"email":             notification.Email,
			"firstName":         firstName,
			"phone":             notification.Phone,
			"contactPreference": notification.ContactPreference,
			"eventTypes":        notification.EventTypes,
			"frequency":         notification.Frequency,
			"termsAccepted":     notification.TermsAccepted,
			"gdprAccepted":      notification.GdprAccepted,
		}

		// Store in Firebase under "submissions/eventNotifications" collection
		submissionId, err := s.firebase.SaveForm("eventNotifications", submissionData)
		if err != nil {
			log.Printf("Error saving event notification submission: %v", err)
			return err
		}

		// Create or update user record
		if userId == "" {
			// Create new user ID
			userId = fmt.Sprintf("user_%d", now)

			// Create user profile
			userData := map[string]interface{}{
				"profile": map[string]interface{}{
					"name":          firstName,
					"email":         notification.Email,
					"phone":         notification.Phone,
					"createdAt":     now,
					"updatedAt":     now,
					"termsAccepted": notification.TermsAccepted,
					"gdprAccepted":  notification.GdprAccepted,
				},
				"auth": map[string]interface{}{
					"hasAccount": false,
				},
				"subscriptions": map[string]interface{}{
					"eventNotifications": map[string]interface{}{
						"subscribed": true,
						"frequency":  notification.Frequency,
					},
				},
				"preferences": map[string]interface{}{
					"contactPreference": notification.ContactPreference,
					"eventTypes":       notification.EventTypes,
				},
			}

			// Save user to Firebase
			if err := s.firebase.SaveUser(userId, userData); err != nil {
				log.Printf("Error saving user from event notification: %v", err)
				return err
			}

			// Update indexes
			if err := s.firebase.UpdateEmailIndex(notification.Email, userId); err != nil {
				log.Printf("Error updating email index: %v", err)
			}

			if notification.Phone != "" {
				if err := s.firebase.UpdatePhoneIndex(notification.Phone, userId); err != nil {
					log.Printf("Error updating phone index: %v", err)
				}
			}

			// Update submission with userId
			updateData := map[string]interface{}{
				"userId": userId,
			}
			submissionRef := s.firebase.db.NewRef("submissions/eventNotifications").Child(submissionId)
			if err := submissionRef.Update(context.Background(), updateData); err != nil {
				log.Printf("Error updating submission with userId: %v", err)
			}
		} else {
			// User exists, update subscription status
			updateData := map[string]interface{}{
				"subscriptions/eventNotifications/subscribed": true,
				"subscriptions/eventNotifications/frequency":  notification.Frequency,
				"preferences/contactPreference":              notification.ContactPreference,
				"preferences/eventTypes":                     notification.EventTypes,
				"profile/updatedAt":                          now,
			}
			userRef := s.firebase.db.NewRef("users").Child(userId)
			if err := userRef.Update(context.Background(), updateData); err != nil {
				log.Printf("Error updating user subscription status: %v", err)
			}
		}
	}

	return nil
}

func (s *SimpleEventNotificationStorage) IsSubscribed(email string) (bool, error) {
	if s.firebase != nil {
		// Get user by email
		userId, err := s.firebase.GetUserByEmail(email)
		if err != nil {
			log.Printf("Error getting user by email for IsSubscribed: %v", err)
			return false, err
		}

		if userId != "" {
			// Check subscription status
			var subscriptionStatus struct {
				Subscribed bool `json:"subscribed"`
			}
			ref := s.firebase.db.NewRef("users").Child(userId).Child("subscriptions/eventNotifications")
			if err := ref.Get(context.Background(), &subscriptionStatus); err != nil {
				log.Printf("Error getting subscription status: %v", err)
				return false, err
			}

			return subscriptionStatus.Subscribed, nil
		}
	}

	return false, nil
}

func (s *SimpleEventNotificationStorage) Unsubscribe(email string) error {
	if s.firebase != nil {
		// Get user by email
		userId, err := s.firebase.GetUserByEmail(email)
		if err != nil {
			log.Printf("Error getting user by email for unsubscribe: %v", err)
			return err
		}

		if userId != "" {
			// Update subscription status to false
			updateData := map[string]interface{}{
				"subscriptions/eventNotifications/subscribed": false,
				"profile/updatedAt":                           time.Now().Unix(),
			}
			userRef := s.firebase.db.NewRef("users").Child(userId)
			if err := userRef.Update(context.Background(), updateData); err != nil {
				log.Printf("Error updating user subscription status: %v", err)
				return err
			}
		}
	}

	return nil
}