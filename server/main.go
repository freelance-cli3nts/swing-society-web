package main

import (
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
	"swinng-society-website/server/internal"
)

func main() {
	if err := godotenv.Load(); err != nil {
			log.Printf("Warning: .env file not found")
	}

	internal.SetupRoutes()

	port := os.Getenv("PORT")
	if port == "" {
			port = "8080"
	}
	
	log.Printf("Server starting on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}


// // Serve static files (CSS, JS, Images)
// func serveStaticFiles() {
// 	fs := http.FileServer(http.Dir("static"))
// 	http.Handle("/static/", http.StripPrefix("/static/", fs))
// }


// // Handle all template requests
// func handleTemplateRequest(w http.ResponseWriter, r *http.Request) {
//   // Print request details
//   log.Printf("\n=== New Request ===")
//   log.Printf("Request URL: %s", r.URL.Path)
  
//   // Get absolute path to requested file
//   absPath, err := filepath.Abs(filepath.Join(".", r.URL.Path))
//   if err != nil {
//       log.Printf("Error getting absolute path: %v", err)
//       http.Error(w, "Server error", http.StatusInternalServerError)
//       return
//   }
//   log.Printf("Absolute path: %s", absPath)
  
//   // Check if file exists
//   _, err = os.Stat(absPath)
//   if os.IsNotExist(err) {
//       log.Printf("File does not exist at: %s", absPath)
//       http.NotFound(w, r)
//       return
//   }
  
//   // If we get here, file exists
//   log.Printf("File found, attempting to serve: %s", absPath)
  
//   // Set headers
//   w.Header().Set("Access-Control-Allow-Origin", "*")
//   w.Header().Set("Content-Type", "text/html")
  
//   // Serve the file
//   http.ServeFile(w, r, absPath)
//   log.Printf("File served successfully")
// }



