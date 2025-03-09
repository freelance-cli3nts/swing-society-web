package storage

import (
	"context"
	"os"
	
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


	// Check for credentials in environment variable
	credentials := os.Getenv("GOOGLE_CREDENTIALS")

	if credentials != "" {
		// Initialize with credentials from environment variable
		opt := option.WithCredentialsJSON([]byte(credentials))
		app, err = firebase.NewApp(ctx, nil, opt)
	} else {
		// Fall back to Application Default Credentials
		app, err = firebase.NewApp(ctx, nil)
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
func (fc *FirebaseClient) SaveForm(formType string, data map[string]interface{}) error {
	ctx := context.Background()
	
	// Store data in a collection based on form type
	ref := fc.db.NewRef(formType)
	_, err := ref.Push(ctx, data)
	
	return err
}

// GetForms retrieves form submissions
func (fc *FirebaseClient) GetForms(formType string) (map[string]interface{}, error) {
	ctx := context.Background()
	
	// Get reference to form type collection
	ref := fc.db.NewRef(formType)
	
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