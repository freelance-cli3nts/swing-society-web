package storage

import (
    "swing-society-website/server/internal/api/models"
)

// NewsletterStorage defines the interface for newsletter subscription storage
type NewsletterStorage interface {
    Subscribe(subscription *models.Newsletter) error
    Unsubscribe(email string) error
    IsSubscribed(email string) (bool, error)
}

type SimpleNewsletterStorage struct {
    subscriptions map[string]*models.Newsletter
}

func NewSimpleNewsletterStorage() *SimpleNewsletterStorage {
    return &SimpleNewsletterStorage{
        subscriptions: make(map[string]*models.Newsletter),
    }
}

func (s *SimpleNewsletterStorage) Subscribe(subscription *models.Newsletter) error {
    s.subscriptions[subscription.Email] = subscription
    return nil
}

func (s *SimpleNewsletterStorage) Unsubscribe(email string) error {
    delete(s.subscriptions, email)
    return nil
}

func (s *SimpleNewsletterStorage) IsSubscribed(email string) (bool, error) {
    _, exists := s.subscriptions[email]
    return exists, nil
}