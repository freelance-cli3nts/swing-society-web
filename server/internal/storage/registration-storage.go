// internal/storage/registration.go
package storage

import (
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
}

func NewSimpleRegistrationStorage() *SimpleRegistrationStorage {
    return &SimpleRegistrationStorage{
        registrations: make(map[string]*models.RegistrationForm),
    }
}

func (s *SimpleRegistrationStorage) StoreRegistration(reg *models.RegistrationForm) error {
    s.registrations[reg.Email] = reg
    return nil
}

func (s *SimpleRegistrationStorage) GetRegistration(email string) (*models.RegistrationForm, error) {
    if reg, exists := s.registrations[email]; exists {
        return reg, nil
    }
    return nil, nil
}