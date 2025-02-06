package main

import (
    "log"
    "net/http"
    "os"
    "github.com/joho/godotenv"
    "swing-society-website/server/internal"
    "swing-society-website/server/internal/config"
)

func main() {
    // Load environment variables
    if err := godotenv.Load(); err != nil {
        log.Printf("Warning: .env file not found")
    }

    // Initialize paths
    config.InitPaths()

    // Create router with Google Cloud project ID
    projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
    if projectID == "" {
        projectID = "swingsociety-backend" // fallback to your project ID
    }

    router, err := internal.NewRouter(projectID)
    if err != nil {
        log.Fatalf("Failed to create router: %v", err)
    }

    // Setup routes
    router.SetupRoutes()

    // Get port from environment
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    log.Printf("Server starting on port %s", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}


