package storage

import (
	"context"
	"os"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"google.golang.org/api/option"
)

// FirebaseClient handles interactions with Firebase
type FirebaseClient struct {
	app *firebase.App
	db  *db.Client
}

// NewFirebaseClient creates a new Firebase client
func NewFirebaseClient() (*FirebaseClient, error) {
	ctx := context.Background()
	
	var app *firebase.App
	var err error
	var config *firebase.Config

	// Set database URL for the Realtime Database
	config = &firebase.Config{
		DatabaseURL: "https://swing-society-realtime-data-default-rtdb.europe-west1.firebasedatabase.app",
	}

	if credentials := os.Getenv("GOOGLE_CREDENTIALS"); credentials != "" {
		// Local dev: explicig JSON via env var
		opt := option.WithCredentialsJSON([]byte(credentials))
		app, err = firebase.NewApp(ctx, config, opt)
	} else {
		app, err = firebase.NewApp(ctx, config)
	}
	
	if err != nil {
		return nil, err
	}
	
	// Get Database client
	dbClient, err := app.Database(ctx)
	if err != nil {
		return nil, err
	}
	
	return &FirebaseClient{
		app: app,
		db:  dbClient,
	}, nil
}

// SaveForm saves form data to Firebase
func (fc *FirebaseClient) SaveForm(formType string, data map[string]interface{}) (string, error) {
	ctx := context.Background()
	
	// Store data in a collection based on form type
	ref := fc.db.NewRef("submissions").Child(formType)
	newRef, err := ref.Push(ctx, data)
	
	if err != nil {
		return "", err
	}
	
	// Return the generated key
	key := newRef.Key
	return key, nil
}

// SaveUser saves or updates a user in Firebase
func (fc *FirebaseClient) SaveUser(userId string, user interface{}) error {
	ctx := context.Background()
	
	// Store user data
	ref := fc.db.NewRef("users").Child(userId)
	return ref.Set(ctx, user)
}

// GetUserByEmail retrieves a user by email using the email index
func (fc *FirebaseClient) GetUserByEmail(email string) (string, error) {
	ctx := context.Background()
	
	// Clean the email for use as a key (Firebase doesn't allow '.' in keys)
	cleanEmail := strings.ReplaceAll(email, ".", "_dot_")
	
	// Use the email index to get user ID
	ref := fc.db.NewRef("indexes/emailToUser").Child(cleanEmail)
	
	var userId string
	if err := ref.Get(ctx, &userId); err != nil {
		return "", err
	}
	
	if userId == "" {
		return "", nil // User not found
	}
	
	return userId, nil
}

// UpdateEmailIndex updates the email-to-user index
func (fc *FirebaseClient) UpdateEmailIndex(email, userId string) error {
	ctx := context.Background()
	
	// Clean the email for use as a key (Firebase doesn't allow '.' in keys)
	cleanEmail := strings.ReplaceAll(email, ".", "_dot_")
	
	// Update email index
	ref := fc.db.NewRef("indexes/emailToUser").Child(cleanEmail)
	err := ref.Set(ctx, userId)
	if err != nil {
	}
	return err
}

// UpdatePhoneIndex updates the phone-to-user index
func (fc *FirebaseClient) UpdatePhoneIndex(phone, userId string) error {
	if phone == "" {
		return nil // Skip if phone is empty
	}
	
	ctx := context.Background()
	
	// Update phone index
	ref := fc.db.NewRef("indexes/phoneToUser").Child(phone)
	return ref.Set(ctx, userId)
}

// GetForms retrieves form submissions
func (fc *FirebaseClient) GetForms(formType string) (map[string]interface{}, error) {
	ctx := context.Background()
	
	// Get reference to form type collection
	ref := fc.db.NewRef("submissions").Child(formType)
	
	var data map[string]interface{}
	err := ref.Get(ctx, &data)
	
	return data, err
}


// TestConnection verifies the Firebase connection is working
func (fc *FirebaseClient) TestConnection(ctx context.Context) (map[string]interface{}, error) {
	// Try to read from a simple test path
	ref := fc.db.NewRef("connectionTest")
	
	var data map[string]interface{}
	if err := ref.Get(ctx, &data); err != nil {
		return nil, err
	}
	
	return data, nil
}