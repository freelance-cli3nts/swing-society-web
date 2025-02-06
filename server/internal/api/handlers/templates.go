	package handlers

	import (
			"log"
			"net/http"
			"os"
			"strings"
			"path/filepath"
			"swing-society-website/server/internal/config"
	)

	func HandleTemplate(w http.ResponseWriter, r *http.Request, templatePath string) {
    // Remove the leading "/templates" from the path since it's already included in TemplatesDir
    cleanPath := strings.TrimPrefix(templatePath, "/templates")
    fullPath := filepath.Join(config.AppPaths.TemplatesDir, cleanPath)
    
    log.Printf("Serving template: %s", fullPath)
    
    if _, err := os.Stat(fullPath); os.IsNotExist(err) {
        log.Printf("Template not found: %s", fullPath)
        http.NotFound(w, r)
        return
    }
    
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Content-Type", "text/html")
    
    http.ServeFile(w, r, fullPath)
}