package models

import (
    "regexp"
    "strings"
		
)

type RegistrationForm struct {
    Name  string `json:"name"`
    Email string `json:"email"`
		Phone string `json:"phone,omitempty"`
		Roles        []string `json:"roles,omitempty"`
    RegisterAlone string   `json:"registerAlone,omitempty"`
    PartnerName  string   `json:"partnerName,omitempty"`
    Source       string   `json:"source,omitempty"`
    Message      string   `json:"message,omitempty"`
    Timestamp    int64    `json:"timestamp,omitempty"`
}

func (f *RegistrationForm) Validate() map[string]string {
    errors := make(map[string]string)

    if strings.TrimSpace(f.Name) == "" {
        errors["name"] = "Моля, въведете вашето име"
    } else if len(f.Name) < 3 {
        errors["name"] = "Името трябва да е поне 3 символа"
    }

    emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
    if strings.TrimSpace(f.Email) == "" {
        errors["email"] = "Моля, въведете вашия имейл"
    } else if !emailRegex.MatchString(strings.ToLower(f.Email)) {
        errors["email"] = "Моля, въведете валиден имейл адрес"
    }

		phoneRegex := regexp.MustCompile(`^(\+359|0|00359)[0-9]{9}$`) // pattern="(\+359|00359|0)\d{9}$"
    if strings.TrimSpace(f.Phone) != "" && !phoneRegex.MatchString(strings.ToLower(f.Phone)) {
        errors["phone"] = "Моля, въведете валиден телефонен номер: 0888123456 / +359888123456 / 00359888123456"
    }

    return errors
}