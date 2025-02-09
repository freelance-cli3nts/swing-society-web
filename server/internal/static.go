package internal

import (
    "log"
		"mime"
		"os"
    "net/http"
		"strings"
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
    fs := http.FileServer(http.Dir(config.AppPaths.StaticDir))

    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Clean the path to prevent directory traversal
        path := filepath.Clean(r.URL.Path)
        
        // Get file extension
        ext := strings.ToLower(filepath.Ext(path))
        
        // Set appropriate MIME type
        if mimeType := mime.TypeByExtension(ext); mimeType != "" {
            w.Header().Set("Content-Type", mimeType)
        }

				if os.Getenv("DOCKER_CONTAINER") != "true" {
					// Development environment: disable caching
					w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
					w.Header().Set("Pragma", "no-cache")
					w.Header().Set("Expires", "0")
				} else {
						// Production environment: enable caching
						w.Header().Set("Cache-Control", "public, max-age=31536000")
				}

        log.Printf("Serving static file: %s with MIME type: %s", path, w.Header().Get("Content-Type"))
        fs.ServeHTTP(w, r)
    })

    http.Handle("/static/", http.StripPrefix("/static/", handler))
}
