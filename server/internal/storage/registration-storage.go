// internal/storage/registration.go
package storage

import (
    "context"
    "fmt"
    "log"
    "time"
    "swing-society-website/server/internal/api/models"
)

// RegistrationStorage defines the interface for registration data storage
type RegistrationStorage interface {
    StoreRegistration(reg *models.RegistrationForm) error
    GetRegistration(email string) (*models.RegistrationForm, error)
}

// SimpleRegistrationStorage implements RegistrationStorage using basic storage
type SimpleRegistrationStorage struct {
    // Could be extended to use a database later
    registrations map[string]*models.RegistrationForm
    // Add Firebase client
    firebase      *FirebaseClient
}

func NewSimpleRegistrationStorage()*SimpleRegistrationStorage {
    // Initialize Firebase client
    firebase, err := NewFirebaseClient()
    if err != nil {
        // Log the error but don't fail - we can fall back to in-memory storage
        log.Printf("Warning: Firebase initialization failed: %v", err)
    }

    return &SimpleRegistrationStorage{
        registrations: make(map[string]*models.RegistrationForm),
        firebase:      firebase,
    }
}

func (s *SimpleRegistrationStorage) StoreRegistration(reg *models.RegistrationForm) error {
    // Store in memory map
    s.registrations[reg.Email] = reg

    // Also store in Firebase if available
    if s.firebase != nil {
        // First, check if user already exists by email
        userId, err := s.firebase.GetUserByEmail(reg.Email)
        if err != nil {
            log.Printf("Error checking user by email: %v", err)
        }

        // Convert registration to map for submission data
        submissionData := map[string]interface{}{
            "timestamp":     reg.Timestamp,
            "userId":        userId, // Will be null if new user
            "name":          reg.Name,
            "email":         reg.Email,
            "phone":         reg.Phone,
            "role":          "",
            "signupMethod":  "",
            "partnerEmail":  reg.PartnerName,
            "signupSource":  reg.Source,
            "message":       reg.Message,
            "termsAccepted": true,
            "gdprAccepted":  true,
        }

        // Set role and signup method
        if len(reg.Roles) > 0 {
            submissionData["role"] = reg.Roles[0]
        }
        
        if reg.RegisterAlone == "yes" {
            submissionData["signupMethod"] = "alone"
        } else {
            submissionData["signupMethod"] = "withPartner"
        }
        
        // Store in Firebase under "submissions/registrations" collection
        submissionId, err := s.firebase.SaveForm("registrations", submissionData)
        if err != nil {
            log.Printf("Error saving registration submission: %v", err)
            return err
        }
        
        // Create or update user record
        if userId == "" {
            // Create new user ID based on timestamp
            userId = fmt.Sprintf("user_%d", time.Now().UnixNano())
            
            // Create user from registration
            user := models.NewUserFromRegistration(reg)
            
            // Add message if present
            user.AddMessageFromRegistration(reg)
            
            // Save user to Firebase
            if err := s.firebase.SaveUser(userId, user); err != nil {
                log.Printf("Error saving user: %v", err)
                return err
            }
            
            // Update indexes
            if err := s.firebase.UpdateEmailIndex(reg.Email, userId); err != nil {
                log.Printf("Error updating email index: %v", err)
            }
            
            if reg.Phone != "" {
                if err := s.firebase.UpdatePhoneIndex(reg.Phone, userId); err != nil {
                    log.Printf("Error updating phone index: %v", err)
                }
            }
            
            // Update submission with userId
            updateData := map[string]interface{}{
                "userId": userId,
            }
            submissionRef := s.firebase.db.NewRef("submissions/registrations").Child(submissionId)
            if err := submissionRef.Update(context.Background(), updateData); err != nil {
                log.Printf("Error updating submission with userId: %v", err)
            }
        } else {
            // User exists, just add the message if any
            if reg.Message != "" {
                messageRef := s.firebase.db.NewRef("users").Child(userId).Child("messages")
                messageId := time.Now().Format("20060102150405")
                messageData := map[string]interface{}{
                    "content":   reg.Message,
                    "timestamp": time.Now().Unix(),
                    "type":      "registration",
                }
                if err := messageRef.Child(messageId).Set(context.Background(), messageData); err != nil {
                    log.Printf("Error adding message to existing user: %v", err)
                }
            }
            
            // Update some user fields
            updateData := map[string]interface{}{
                "profile/updatedAt": time.Now().Unix(),
            }
            userRef := s.firebase.db.NewRef("users").Child(userId)
            if err := userRef.Update(context.Background(), updateData); err != nil {
                log.Printf("Error updating user fields: %v", err)
            }
        }
    }

    return nil
}

func (s *SimpleRegistrationStorage) GetRegistration(email string) (*models.RegistrationForm, error) {
    // Try in-memory first
    if reg, exists := s.registrations[email]; exists {
        return reg, nil
    }
    
    // If Firebase available, check there
    if s.firebase != nil {
        userId, err := s.firebase.GetUserByEmail(email)
        if err != nil {
            return nil, err
        }
        
        if userId != "" {
            // User exists, extract registration info
            userRef := s.firebase.db.NewRef("users").Child(userId).Child("profile")
            var profile map[string]interface{}
            if err := userRef.Get(context.Background(), &profile); err != nil {
                return nil, err
            }
            
            // Convert to registration form
            reg := &models.RegistrationForm{
                Name:  profile["name"].(string),
                Email: profile["email"].(string),
            }
            
            if phone, ok := profile["phone"].(string); ok {
                reg.Phone = phone
            }
            
            return reg, nil
        }
    }
    
    return nil, nil
}