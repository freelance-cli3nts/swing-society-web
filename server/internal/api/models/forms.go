package models

// import "time"

type ContactForm struct {
    Name    string `json:"name"`
    Email   string `json:"email"`
    Phone   string `json:"phone,omitempty"`
    Message string `json:"message"`
}

type ClassInquiry struct {
    Name      string `json:"name"`
    Email     string `json:"email"`
    Phone     string `json:"phone,omitempty"`
    ClassType string `json:"classType"`
    Level     string `json:"level"`
    Message   string `json:"message,omitempty"`
}

type Newsletter struct {
    Email string `json:"email"`
    Name  string `json:"name,omitempty"`
}