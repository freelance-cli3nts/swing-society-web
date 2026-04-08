package storage

import (
	"context"
	"log"
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
		log.Printf("Error getting user by email %s: %v", email, err)
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
		log.Printf("Error updating email index for %s: %v", email, err)
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


// RESTful API 4 read/write to Firebase
// func writeData(w http.ResponseWriter, r *http.Request) {
// 	// Write data to Firebase
// 	ref := client.NewRef("data")
// 	if err := ref.Set(ctx, map[string]interface{}{
// 			"key": "value",
// 	}); err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 	}
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("Data written successfully!"))
// }

// DeepSeek suggestion
// import (
//     "context"
//     "log"
//     "net/http"
//     "firebase.google.com/go"
//     "google.golang.org/api/option"
// )

// var client *firebase.Database

// func main() {
//     // Initialize Firebase Admin SDK
//     opt := option.WithCredentialsFile("path/to/serviceAccountKey.json") // Path to your Firebase service account key
//     app, err := firebase.NewApp(context.Background(), nil, opt)
//     if err != nil {
//         log.Fatalf("error initializing Firebase app: %v\n", err)
//     }

//     // Initialize Realtime Database client
//     client, err = app.Database(context.Background())
//     if err != nil {
//         log.Fatalf("error initializing database client: %v\n", err)
//     }

//     // Set up HTTP routes
//     http.HandleFunc("/submit-form", handleFormSubmission)

//     // Start the server
//     log.Println("Server started on :8080")
//     log.Fatal(http.ListenAndServe(":8080", nil))
// }

// func handleFormSubmission(w http.ResponseWriter, r *http.Request) {
//     // Parse form data
//     name := r.FormValue("name")
//     email := r.FormValue("email")
//     message := r.FormValue("message")

//     // Validate data (example)
//     if name == "" || email == "" || message == "" {
//         http.Error(w, "All fields are required", http.StatusBadRequest)
//         return
//     }

//     // Save data to Firebase Realtime Database
//     ref := client.NewRef("formSubmissions/")
//     err := ref.Push(context.Background(), map[string]interface{}{
//         "name":    name,
//         "email":   email,
//         "message": message,
//     })
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }

//     // Respond to the client
//     w.WriteHeader(http.StatusOK)
//     w.Write([]byte("Form submitted successfully!"))
// }