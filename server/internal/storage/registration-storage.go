// internal/storage/registration.go
package storage

import (
		"context"
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
        // Consider adding a logger here instead of returning the error
    }

    return &SimpleRegistrationStorage{
        registrations: make(map[string]*models.RegistrationForm),
				firebase:      firebase,
    }
}

func (s *SimpleRegistrationStorage) StoreRegistration(reg *models.RegistrationForm) error {
    s.registrations[reg.Email] = reg

		// Also store in Firebase if available
    if s.firebase != nil {
        // Convert registration to map for Firebase
        data := map[string]interface{}{
            "name":  reg.Name,
            "email": reg.Email,
            "phone": reg.Phone,
            // Add timestamp
            "timestamp": context.Background().Value("requestTime"),
        }
        
        // Store in Firebase under "registrations" collection
        if err := s.firebase.SaveForm("registrations", data); err != nil {
            // Log the error but don't fail the registration
            return nil	
        }
    }

    return nil
}

func (s *SimpleRegistrationStorage) GetRegistration(email string) (*models.RegistrationForm, error) {
    if reg, exists := s.registrations[email]; exists {
        return reg, nil
    }
    return nil, nil
}