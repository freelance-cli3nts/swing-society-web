package models

import (
    "regexp"
    "strings"
		
)

type RegistrationForm struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

func (f *RegistrationForm) Validate() map[string]string {
    errors := make(map[string]string)

    if strings.TrimSpace(f.Name) == "" {
        errors["name"] = "Моля, въведете вашето име"
    } else if len(f.Name) < 2 {
        errors["name"] = "Името трябва да е поне 2 символа"
    }

    emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
    if strings.TrimSpace(f.Email) == "" {
        errors["email"] = "Моля, въведете вашия имейл"
    } else if !emailRegex.MatchString(strings.ToLower(f.Email)) {
        errors["email"] = "Моля, въведете валиден имейл адрес"
    }

    return errors
}