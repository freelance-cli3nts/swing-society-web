package models

import "strings"

type ContactForm struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone,omitempty"`
	Message string `json:"message"`
}

func (f *ContactForm) Validate() map[string]string {
	errors := make(map[string]string)

	if strings.TrimSpace(f.Name) == "" {
		errors["name"] = "Моля, въведете вашето име"
	}
	if strings.TrimSpace(f.Email) == "" {
		errors["email"] = "Моля, въведете вашия имейл"
	}
	if strings.TrimSpace(f.Message) == "" {
		errors["message"] = "Моля, въведете съобщение"
	}

	return errors
}
