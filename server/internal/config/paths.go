package config

import (
    "log"
    "os"
    "path/filepath"
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
        rootDir = filepath.Dir(serverDir) // Go up one level from server/
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

    // Verify directories exist
    verifyDir(AppPaths.TemplatesDir, "Templates")
    verifyDir(AppPaths.StaticDir, "Static")
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
//     var err error

//     // Get the current working directory
//     serverDir, err := os.Getwd()
//     if err != nil {
//         log.Fatalf("Failed to get current working directory: %v", err)
//     }
		
// 		// Get the project root directory (one level up)
//     rootDir := filepath.Dir(serverDir)


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
