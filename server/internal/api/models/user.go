package models

import (
	"time"
)

// User represents a user profile in the system
type User struct {
	Profile struct {
		Name          string `json:"name"`
		Email         string `json:"email"`
		Phone         string `json:"phone"`
		Role          string `json:"role"`
		SignupMethod  string `json:"signupMethod"`
		SignupSource  string `json:"signupSource"`
		PartnerId     string `json:"partnerId,omitempty"`
		CreatedAt     int64  `json:"createdAt"`
		UpdatedAt     int64  `json:"updatedAt"`
		TermsAccepted bool   `json:"termsAccepted"`
		GdprAccepted  bool   `json:"gdprAccepted"`
	} `json:"profile"`

	Auth struct {
		HasAccount       bool  `json:"hasAccount"`
		AccountCreatedAt int64 `json:"accountCreatedAt,omitempty"`
	} `json:"auth"`

	Subscriptions struct {
		Newsletter struct {
			Subscribed bool   `json:"subscribed"`
			Frequency  string `json:"frequency,omitempty"`
		} `json:"newsletter"`
		EventNotifications struct {
			Subscribed bool   `json:"subscribed"`
			Frequency  string `json:"frequency,omitempty"`
		} `json:"eventNotifications"`
		Paid struct {
			Active    bool    `json:"active"`
			Plan      *string `json:"plan,omitempty"`
			StartDate *int64  `json:"startDate,omitempty"`
			EndDate   *int64  `json:"endDate,omitempty"`
			AutoRenew bool    `json:"autoRenew"`
		} `json:"paid"`
	} `json:"subscriptions"`

	Preferences struct {
		ContactPreference string `json:"contactPreference,omitempty"`
	} `json:"preferences"`

	Relationships struct {
		Partner string                      `json:"partner,omitempty"`
		Friends map[string]FriendConnection `json:"friends,omitempty"`
	} `json:"relationships"`

	Messages map[string]Message `json:"messages,omitempty"`
}

// FriendConnection represents a connection between users
type FriendConnection struct {
	Since  int64  `json:"since"`
	Status string `json:"status"` // confirmed, pending, requested
}

// Message represents a message from a user
type Message struct {
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
	Type      string `json:"type"`
}

// NewUserFromRegistration creates a user from registration form data
func NewUserFromRegistration(reg *RegistrationForm) *User {
	now := time.Now().Unix()
	user := &User{}
	
	// Set profile data
	user.Profile.Name = reg.Name
	user.Profile.Email = reg.Email
	user.Profile.Phone = reg.Phone
	user.Profile.CreatedAt = now
	user.Profile.UpdatedAt = now
	user.Profile.TermsAccepted = true
	user.Profile.GdprAccepted = true
	
	// Set roles if provided
	if len(reg.Roles) > 0 {
		user.Profile.Role = reg.Roles[0] // Use first selected role as primary
	}
	
	// Set signup method
	if reg.RegisterAlone == "yes" {
		user.Profile.SignupMethod = "alone"
	} else {
		user.Profile.SignupMethod = "withPartner"
	}
	
	// Set signup source
	user.Profile.SignupSource = reg.Source
	
	// Set auth status
	user.Auth.HasAccount = false
	
	return user
}

// AddMessageFromRegistration adds a message from registration form
func (u *User) AddMessageFromRegistration(reg *RegistrationForm) {
	if reg.Message == "" {
		return
	}
	
	if u.Messages == nil {
		u.Messages = make(map[string]Message)
	}
	
	messageId := time.Now().Format("20060102150405")
	u.Messages[messageId] = Message{
		Content:   reg.Message,
		Timestamp: time.Now().Unix(),
		Type:      "registration",
	}
}