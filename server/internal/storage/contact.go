// internal/storage/contact.go
package storage

import (
    "context"
    "fmt"
    "log"
    "time"
    "swing-society-website/server/internal/api/models"
)

// ContactStorage defines the interface for contact form storage
type ContactStorage interface {
    StoreContactForm(form *models.ContactForm) error
    GetContactForm(email string) (*models.ContactForm, error)
    GetAllContactForms() ([]*models.ContactForm, error)
}

type SimpleContactStorage struct {
    firebase *FirebaseClient
}

func NewSimpleContactStorage(firebase *FirebaseClient) *SimpleContactStorage {
    return &SimpleContactStorage{
        firebase: firebase,
    }
}

func (s *SimpleContactStorage) StoreContactForm(form *models.ContactForm) error {
    if s.firebase != nil {
        // Check if user already exists by email
        userId, err := s.firebase.GetUserByEmail(form.Email)
        if err != nil {
            log.Printf("Error checking user by email for contact: %v", err)
        }

        // Create submission data
        submissionData := map[string]interface{}{
            "timestamp":         time.Now().Unix(),
            "name":              form.Name,
            "email":             form.Email,
            "phone":             form.Phone,
            "message":           form.Message,
            "contactPreference": "email", // Default
        }

        // Store in Firebase under "submissions/contacts" collection
        submissionId, err := s.firebase.SaveForm("contacts", submissionData)
        if err != nil {
            log.Printf("Error saving contact submission: %v", err)
            return err
        }

        // Create or update user record
        if userId == "" {
            // Create new user ID
            userId = fmt.Sprintf("user_%d", time.Now().UnixNano())

            // Create user from contact form
            now := time.Now().Unix()
            userData := map[string]interface{}{
                "profile": map[string]interface{}{
                    "name":          form.Name,
                    "email":         form.Email,
                    "phone":         form.Phone,
                    "createdAt":     now,
                    "updatedAt":     now,
                    "termsAccepted": false,
                    "gdprAccepted":  false,
                },
                "auth": map[string]interface{}{
                    "hasAccount": false,
                },
                "preferences": map[string]interface{}{
                    "contactPreference": "email",
                },
                "messages": map[string]interface{}{
                    fmt.Sprintf("contact_%d", now): map[string]interface{}{
                        "content":   form.Message,
                        "timestamp": now,
                        "type":      "contact",
                    },
                },
            }

            // Save user to Firebase
            if err := s.firebase.SaveUser(userId, userData); err != nil {
                log.Printf("Error saving user from contact: %v", err)
                return err
            }

            // Update indexes
            if err := s.firebase.UpdateEmailIndex(form.Email, userId); err != nil {
                log.Printf("Error updating email index: %v", err)
            }

            if form.Phone != "" {
                if err := s.firebase.UpdatePhoneIndex(form.Phone, userId); err != nil {
                    log.Printf("Error updating phone index: %v", err)
                }
            }

            // Update submission with userId
            updateData := map[string]interface{}{
                "userId": userId,
            }
            submissionRef := s.firebase.db.NewRef("submissions/contacts").Child(submissionId)
            if err := submissionRef.Update(context.Background(), updateData); err != nil {
                log.Printf("Error updating submission with userId: %v", err)
            }
        } else {
            // User exists, just add the message
            messageRef := s.firebase.db.NewRef("users").Child(userId).Child("messages")
            messageId := fmt.Sprintf("contact_%d", time.Now().Unix())
            messageData := map[string]interface{}{
                "content":   form.Message,
                "timestamp": time.Now().Unix(),
                "type":      "contact",
            }
            if err := messageRef.Child(messageId).Set(context.Background(), messageData); err != nil {
                log.Printf("Error adding message to existing user: %v", err)
            }

            // Update user timestamp
            updateData := map[string]interface{}{
                "profile/updatedAt": time.Now().Unix(),
            }
            userRef := s.firebase.db.NewRef("users").Child(userId)
            if err := userRef.Update(context.Background(), updateData); err != nil {
                log.Printf("Error updating user timestamp: %v", err)
            }
        }
    }

    return nil
}

func (s *SimpleContactStorage) GetContactForm(email string) (*models.ContactForm, error) {
    if s.firebase != nil {
        // Get submissions for this email
        var submissions map[string]interface{}
        ref := s.firebase.db.NewRef("submissions/contacts")
        if err := ref.OrderByChild("email").EqualTo(email).LimitToLast(1).Get(context.Background(), &submissions); err != nil {
            log.Printf("Error fetching contact from Firebase: %v", err)
            return nil, err
        }
        
        // If submissions found, convert first one to ContactForm
        for _, data := range submissions {
            submissionData, ok := data.(map[string]interface{})
            if !ok {
                continue
            }
            
            form := &models.ContactForm{
                Name:    submissionData["name"].(string),
                Email:   submissionData["email"].(string),
                Message: submissionData["message"].(string),
            }
            
            if phone, ok := submissionData["phone"].(string); ok {
                form.Phone = phone
            }
            
            return form, nil
        }
    }
    
    return nil, nil
}

func (s *SimpleContactStorage) GetAllContactForms() ([]*models.ContactForm, error) {
    if s.firebase != nil {
        var submissions map[string]interface{}
        ref := s.firebase.db.NewRef("submissions/contacts")
        if err := ref.Get(context.Background(), &submissions); err != nil {
            log.Printf("Error fetching contacts from Firebase: %v", err)
            return nil, err
        }
        
        // Convert submissions to ContactForm array
        forms := make([]*models.ContactForm, 0, len(submissions))
        for _, data := range submissions {
            submissionData, ok := data.(map[string]interface{})
            if !ok {
                continue
            }
            
            form := &models.ContactForm{
                Name:    submissionData["name"].(string),
                Email:   submissionData["email"].(string),
                Message: submissionData["message"].(string),
            }
            
            if phone, ok := submissionData["phone"].(string); ok {
                form.Phone = phone
            }
            
            forms = append(forms, form)
        }
        
        return forms, nil
    }
    
    return nil, nil
}