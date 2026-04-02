package models

import (
	"regexp"
	"strings"
)

type EventNotification struct {
	Email             string   `json:"email"`
	FirstName         string   `json:"firstName,omitempty"`
	Phone             string   `json:"phone,omitempty"`
	ContactPreference string   `json:"contactPreference,omitempty"`
	EventTypes        []string `json:"eventTypes,omitempty"`
	Frequency         string   `json:"frequency,omitempty"`
	TermsAccepted     bool     `json:"termsAccepted"`
	GdprAccepted      bool     `json:"gdprAccepted"`
	Timestamp         int64    `json:"timestamp,omitempty"`
}

func (n *EventNotification) Validate() map[string]string {
	errors := make(map[string]string)

	// Email validation
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if strings.TrimSpace(n.Email) == "" {
		errors["email"] = "Моля, въведете вашия имейл"
	} else if !emailRegex.MatchString(strings.ToLower(n.Email)) {
		errors["email"] = "Моля, въведете валиден имейл адрес"
	}

	// Phone validation (if provided)
	phoneRegex := regexp.MustCompile(`^(\+359|0|00359)[0-9]{9}$`)
	if strings.TrimSpace(n.Phone) != "" && !phoneRegex.MatchString(strings.ToLower(n.Phone)) {
		errors["phone"] = "Моля, въведете валиден телефонен номер: 0888123456 / +359888123456 / 00359888123456"
	}

	nameRegex := regexp.MustCompile(`^[a-zA-Z\s]+$`)
	if len(strings.TrimSpace(n.FirstName)) < 3 || !nameRegex.MatchString(n.FirstName) {
		errors["firstName"] = "Моля, въведете валидно име (поне 3 символа)"
	}

	// GDPR validation
	if !n.GdprAccepted {
		errors["gdpr"] = "Трябва да се съгласите с GDPR условията"
	}

	// Terms validation
	if !n.TermsAccepted {
		errors["regulations"] = "Трябва да се съгласите с правилата и условията"
	}

	return errors
}