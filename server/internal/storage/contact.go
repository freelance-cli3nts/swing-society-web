// internal/storage/contact.go
package storage

import (
    "swing-society-website/server/internal/api/models"
)

// ContactStorage defines the interface for contact form storage
type ContactStorage interface {
    StoreContactForm(form *models.ContactForm) error
    GetContactForm(email string) (*models.ContactForm, error)
    GetAllContactForms() ([]*models.ContactForm, error)
}

type SimpleContactStorage struct {
    contacts []*models.ContactForm
}

func NewSimpleContactStorage() *SimpleContactStorage {
    return &SimpleContactStorage{
        contacts: make([]*models.ContactForm, 0),
    }
}

func (s *SimpleContactStorage) StoreContactForm(form *models.ContactForm) error {
    s.contacts = append(s.contacts, form)
    return nil
}

func (s *SimpleContactStorage) GetContactForm(email string) (*models.ContactForm, error) {
    for _, contact := range s.contacts {
        if contact.Email == email {
            return contact, nil
        }
    }
    return nil, nil
}

func (s *SimpleContactStorage) GetAllContactForms() ([]*models.ContactForm, error) {
    return s.contacts, nil
}