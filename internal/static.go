package static

import "net/http"

// ServeFiles sets up static file serving
func ServeFiles() {
    fs := http.FileServer(http.Dir("static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))
}