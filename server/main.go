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

    // Load configuration
    if err := config.LoadConfig(); err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    // Create router with Google Cloud project ID from config
    projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
    if projectID == "" {
        projectID = config.AppConfig.External.ProjectID
    }

    router, err := internal.NewRouter(projectID)
    if err != nil {
        log.Fatalf("Failed to create router: %v", err)
    }

    // Setup routes
    if err := router.SetupRoutes(); err != nil {
        log.Fatalf("Failed to setup routes: %v", err)
    }

    // Get port from config
    port := config.AppConfig.Server.Port
    
    log.Printf("Server starting on port %s", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}