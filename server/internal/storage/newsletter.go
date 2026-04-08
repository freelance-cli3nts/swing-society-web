package storage

import (
    "context"
    "fmt"
    "log"
    "time"
    "swing-society-website/server/internal/api/models"
)

// NewsletterStorage defines the interface for newsletter subscription storage
type NewsletterStorage interface {
    Subscribe(subscription *models.Newsletter) error
    Unsubscribe(email string) error
    IsSubscribed(email string) (bool, error)
}

type SimpleNewsletterStorage struct {
    firebase *FirebaseClient
}

func NewSimpleNewsletterStorage(firebase *FirebaseClient) *SimpleNewsletterStorage {
    return &SimpleNewsletterStorage{
        firebase: firebase,
    }
}

func (s *SimpleNewsletterStorage) Subscribe(subscription *models.Newsletter) error {
    if s.firebase != nil {
        // Check if user already exists by email
        userId, err := s.firebase.GetUserByEmail(subscription.Email)
        if err != nil {
            log.Printf("Error checking user by email for newsletter: %v", err)
        }

        now := time.Now().Unix()
        firstName := subscription.Name
        if firstName == "" {
            firstName = "Subscriber"
        }

        // Create submission data
        submissionData := map[string]interface{}{
            "timestamp":    now,
            "email":        subscription.Email,
            "firstName":    firstName,
            "frequency":    "weekly", // Default
            "termsAccepted": true,
            "gdprAccepted":  true,
        }

        // Store in Firebase under "submissions/newsletters" collection
        submissionId, err := s.firebase.SaveForm("newsletters", submissionData)
        if err != nil {
            log.Printf("Error saving newsletter submission: %v", err)
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
                    "email":         subscription.Email,
                    "createdAt":     now,
                    "updatedAt":     now,
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
                "preferences": map[string]interface{}{
                    "contactPreference": "email",
                },
            }

            // Save user to Firebase
            if err := s.firebase.SaveUser(userId, userData); err != nil {
                log.Printf("Error saving user from newsletter: %v", err)
                return err
            }

            // Update indexes
            if err := s.firebase.UpdateEmailIndex(subscription.Email, userId); err != nil {
                log.Printf("Error updating email index: %v", err)
            }

            // Update submission with userId
            updateData := map[string]interface{}{
                "userId": userId,
            }
            submissionRef := s.firebase.db.NewRef("submissions/newsletters").Child(submissionId)
            if err := submissionRef.Update(context.Background(), updateData); err != nil {
                log.Printf("Error updating submission with userId: %v", err)
            }
        } else {
            // User exists, update subscription status
            updateData := map[string]interface{}{
                "subscriptions/newsletter/subscribed": true,
                "subscriptions/newsletter/frequency":  "weekly",
                "profile/updatedAt":                   now,
            }
            userRef := s.firebase.db.NewRef("users").Child(userId)
            if err := userRef.Update(context.Background(), updateData); err != nil {
                log.Printf("Error updating user subscription status: %v", err)
            }
        }
    }

    return nil
}

func (s *SimpleNewsletterStorage) Unsubscribe(email string) error {
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
                "subscriptions/newsletter/subscribed": false,
                "profile/updatedAt":                   time.Now().Unix(),
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

func (s *SimpleNewsletterStorage) IsSubscribed(email string) (bool, error) {
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
            ref := s.firebase.db.NewRef("users").Child(userId).Child("subscriptions/newsletter")
            if err := ref.Get(context.Background(), &subscriptionStatus); err != nil {
                log.Printf("Error getting subscription status: %v", err)
                return false, err
            }

            return subscriptionStatus.Subscribed, nil
        }
    }

    return false, nil
}