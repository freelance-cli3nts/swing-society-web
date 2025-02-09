package config

import (
    "log"
    "os"
    "path/filepath"
    "strings"
)

type Paths struct {
    RootDir      string
    TemplatesDir string
    StaticDir    string
}

var AppPaths *Paths

func InitPaths() {
    var rootDir string

    // Check if we're in Docker/Cloud Run
    if os.Getenv("DOCKER_CONTAINER") == "true" {
        rootDir = "/app"
    } else {
        // Development environment
        serverDir, err := os.Getwd()
        if err != nil {
            log.Fatalf("Failed to get current working directory: %v", err)
        }
        
        // Find project root
        rootDir = findProjectRoot(serverDir)
    }

    AppPaths = &Paths{
        RootDir:      rootDir,
        TemplatesDir: filepath.Join(rootDir, "templates"),
        StaticDir:    filepath.Join(rootDir, "static"),
    }

    // Log paths for debugging
    log.Printf("Root Dir: %s", AppPaths.RootDir)
    log.Printf("Templates Dir: %s", AppPaths.TemplatesDir)
    log.Printf("Static Dir: %s", AppPaths.StaticDir)

    // Create directories if they don't exist in development mode
    if os.Getenv("DOCKER_CONTAINER") != "true" {
        os.MkdirAll(AppPaths.TemplatesDir, 0755)
        os.MkdirAll(AppPaths.StaticDir, 0755)
    }

    // Verify directories exist
    verifyDir(AppPaths.TemplatesDir, "Templates")
    verifyDir(AppPaths.StaticDir, "Static")
}

// findProjectRoot looks for go.mod file to determine project root
func findProjectRoot(dir string) string {
    // Check if this is the project directory (contains go.mod)
    if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
        return dir
    }

    // If we're in the server directory, go up one level
    if strings.HasSuffix(dir, "server") {
        parentDir := filepath.Dir(dir)
        // Check if parent directory has go.mod
        if _, err := os.Stat(filepath.Join(parentDir, "go.mod")); err == nil {
            return parentDir
        }
    }

    // If we can't find go.mod, use the directory above server/
    if strings.HasSuffix(dir, "server") {
        return filepath.Dir(dir)
    }

    // Default to current directory if we can't determine project root
    return dir
}

func verifyDir(path string, name string) {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        log.Printf("Warning: %s directory not found at: %s", name, path)
    } else {
        log.Printf("%s directory found at: %s", name, path)
    }
}

// package config

// import (
//     "log"
//     "os"
//     "path/filepath"
// )

// type Paths struct {
//     RootDir      string
//     TemplatesDir string
//     StaticDir    string
// }

// var AppPaths *Paths

// func InitPaths() {
//     var rootDir string

//     // Check if we're in Docker/Cloud Run
//     if os.Getenv("DOCKER_CONTAINER") == "true" {
//         rootDir = "/app"
//     } else {
//         // Development environment
//         serverDir, err := os.Getwd()
//         if err != nil {
//             log.Fatalf("Failed to get current working directory: %v", err)
//         }
//         rootDir = filepath.Dir(serverDir) // Go up one level from server/
//     }

//     AppPaths = &Paths{
//         RootDir:      rootDir,
//         TemplatesDir: filepath.Join(rootDir, "templates"),
//         StaticDir:    filepath.Join(rootDir, "static"),
//     }

//     // Log paths for debugging
//     log.Printf("Root Dir: %s", AppPaths.RootDir)
//     log.Printf("Templates Dir: %s", AppPaths.TemplatesDir)
//     log.Printf("Static Dir: %s", AppPaths.StaticDir)

//     // Verify directories exist
//     verifyDir(AppPaths.TemplatesDir, "Templates")
//     verifyDir(AppPaths.StaticDir, "Static")
// }

// func verifyDir(path string, name string) {
//     if _, err := os.Stat(path); os.IsNotExist(err) {
//         log.Printf("Warning: %s directory not found at: %s", name, path)
//     } else {
//         log.Printf("%s directory found at: %s", name, path)
//     }
// }


// // findProjectRoot looks for go.mod file to determine project root
// func findProjectRoot(dir string) string {
// 	// Check if this is the project directory (contains go.mod)
// 	if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
// 			return dir
// 	}

// 	// If we're in the server directory, go up one level
// 	if strings.HasSuffix(dir, "server") {
// 			return filepath.Dir(dir)
// 	}

// 	// Default to current directory if we can't find project root
// 	return dir
// }