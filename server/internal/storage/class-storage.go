package storage

import (
    "swing-society-website/server/internal/api/models"
)

// ClassStorage defines the interface for class inquiry storage
type ClassStorage interface {
    StoreInquiry(inquiry *models.ClassInquiry) error
    GetInquiriesByType(classType string) ([]*models.ClassInquiry, error)
    GetInquiriesByEmail(email string) ([]*models.ClassInquiry, error)
}

// SimpleClassStorage implements ClassStorage using basic storage
type SimpleClassStorage struct {
    // Could be extended to use a database later
    inquiries []*models.ClassInquiry
}

func NewSimpleClassStorage() *SimpleClassStorage {
    return &SimpleClassStorage{
        inquiries: make([]*models.ClassInquiry, 0),
    }
}

func (s *SimpleClassStorage) StoreInquiry(inquiry *models.ClassInquiry) error {
    s.inquiries = append(s.inquiries, inquiry)
    return nil
}

func (s *SimpleClassStorage) GetInquiriesByType(classType string) ([]*models.ClassInquiry, error) {
    var result []*models.ClassInquiry
    for _, inquiry := range s.inquiries {
        if inquiry.ClassType == classType {
            result = append(result, inquiry)
        }
    }
    return result, nil
}

func (s *SimpleClassStorage) GetInquiriesByEmail(email string) ([]*models.ClassInquiry, error) {
    var result []*models.ClassInquiry
    for _, inquiry := range s.inquiries {
        if inquiry.Email == email {
            result = append(result, inquiry)
        }
    }
    return result, nil
}