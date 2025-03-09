package internal

import (
    "log"
		"mime"
    "net/http"
		"path/filepath"
    "swing-society-website/server/internal/config"
)

func init() {
	// Register common MIME types
	mime.AddExtensionType(".css", "text/css")
	mime.AddExtensionType(".js", "application/javascript")
	mime.AddExtensionType(".html", "text/html")
	mime.AddExtensionType(".png", "image/png")
	mime.AddExtensionType(".jpg", "image/jpeg")
	mime.AddExtensionType(".jpeg", "image/jpeg")
	mime.AddExtensionType(".gif", "image/gif")
	mime.AddExtensionType(".svg", "image/svg+xml")
}


func ServeFiles() {
	fs := http.FileServer(http.Dir(config.AppConfig.Paths.StaticDir))
	
	log.Printf("Entered ServeFiles function: %s", config.AppConfig.Paths.StaticDir)


	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := filepath.Clean(r.URL.Path)
			ext := filepath.Ext(path)
			
			if mimeType := mime.TypeByExtension(ext); mimeType != "" {
					w.Header().Set("Content-Type", mimeType)
			}

			log.Printf("Entered handler function. Path: %s", path)
			// Cache control based on environment
			if config.AppConfig.Environment == "development" {
					w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
					w.Header().Set("Pragma", "no-cache")
					w.Header().Set("Expires", "0")
			} else {
					w.Header().Set("Cache-Control", "public, max-age=31536000")
			}

			log.Printf("Serving static file: %s", path)
			fs.ServeHTTP(w, r)
	})

	http.Handle("/static/", http.StripPrefix("/static/", handler))
}