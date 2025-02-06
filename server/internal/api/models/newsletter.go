package models

import (
	"regexp"
	"strings"
)

type Newsletter struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

func (n *Newsletter) Validate() map[string]string {
	errors := make(map[string]string)

	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if strings.TrimSpace(n.Email) == "" {
		errors["email"] = "Моля, въведете вашия имейл"
	} else if !emailRegex.MatchString(strings.ToLower(n.Email)) {
		errors["email"] = "Моля, въведете валиден имейл адрес"
	}

	return errors
}
