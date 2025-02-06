package models

import "strings"

type ClassInquiry struct {
    Name      string `json:"name"`
    Email     string `json:"email"`
    Phone     string `json:"phone,omitempty"`
    ClassType string `json:"classType"`
    Level     string `json:"level"`
    Message   string `json:"message,omitempty"`

}

func (c *ClassInquiry) Validate() map[string]string {
    errors := make(map[string]string)

    if strings.TrimSpace(c.Name) == "" {
        errors["name"] = "Моля, въведете вашето име"
    }
    if strings.TrimSpace(c.Email) == "" {
        errors["email"] = "Моля, въведете вашия имейл"
    }
    if strings.TrimSpace(c.ClassType) == "" {
        errors["classType"] = "Моля, изберете тип клас"
    }
    if strings.TrimSpace(c.Level) == "" {
        errors["level"] = "Моля, изберете ниво"
    }

    return errors
}