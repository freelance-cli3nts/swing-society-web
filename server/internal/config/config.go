// internal/config/config.go
package config

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "time"
)

// Config holds all configuration for the application
type Config struct {
    Server struct {
        Port         string        `json:"port"`
        Host         string        `json:"host"`
        ReadTimeout  time.Duration `json:"read_timeout"`
        WriteTimeout time.Duration `json:"write_timeout"`
    } `json:"server"`

    Paths struct {
        RootDir      string `json:"root_dir"`
        TemplatesDir string `json:"templates_dir"`
        StaticDir    string `json:"static_dir"`
        DataDir      string `json:"data_dir"`
    } `json:"paths"`

    Security struct {
        AllowedOrigins []string          `json:"allowed_origins"`
        RateLimits     map[string]string `json:"rate_limits"`
        CSPDirectives  string            `json:"csp_directives"`
    } `json:"security"`

    External struct {
        GoogleSheetsURL string `json:"google_sheets_url"`
        ProjectID       string `json:"project_id"`
    } `json:"external"`

    Environment string `json:"environment"`
}

var (
    // AppConfig holds the global configuration
    AppConfig *Config
)

// LoadConfig loads configuration from file and environment variables
func LoadConfig() error {
    config := &Config{}

    // Set default configuration
    setDefaults(config)

    // Load configuration from file
    if err := loadConfigFile(config); err != nil {
        // Only proceed with defaults if we're in development mode
        if os.Getenv("ENVIRONMENT") != "development" {
            return fmt.Errorf("failed to load config file: %v", err)
        }
        log.Printf("Warning: Could not load config file: %v", err)
    }

    // Override with environment variables
    overrideWithEnv(config)

    // Initialize paths
    if err := initializePaths(config); err != nil {
        return fmt.Errorf("failed to initialize paths: %v", err)
    }

    // Validate configuration
    if err := validateConfig(config); err != nil {
        return fmt.Errorf("invalid configuration: %v", err)
    }

    AppConfig = config
    return nil
}

// setDefaults initializes default configuration values
func setDefaults(config *Config) {
    config.Server.Port = "3001"
    config.Server.Host = "localhost"
    config.Server.ReadTimeout = 15 * time.Second
    config.Server.WriteTimeout = 15 * time.Second

    config.Security.AllowedOrigins = []string{"*"}
    config.Security.RateLimits = map[string]string{
        "api":      "100-M",  // 100 requests per minute
        "static":   "1000-M", // 1000 requests per minute
        "default":  "60-M",   // 60 requests per minute
    }
    config.Security.CSPDirectives = "default-src 'self' https:; " +
        "script-src 'self' 'unsafe-inline' https://unpkg.com; " +
        "style-src 'self' 'unsafe-inline' https:; " +
        "img-src 'self' https: data:; " +
        "frame-src 'self' https://www.youtube.com; " +
        "connect-src 'self' https:;"

    config.Environment = "development"
}

// loadConfigFile loads configuration from `config.json`
func loadConfigFile(config *Config) error {
    configPath := filepath.Join(findProjectRoot(), "config.json")

    if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
        configPath = envPath
    }

    file, err := os.Open(configPath)
    if err != nil {
        return fmt.Errorf("could not open config file: %v", err)
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    if err := decoder.Decode(config); err != nil {
        return fmt.Errorf("could not decode config file: %v", err)
    }

    // Override with production settings if applicable
    if env := os.Getenv("ENVIRONMENT"); env == "production" {
        config.Server.Host = "0.0.0.0"
        config.Security.AllowedOrigins = []string{"https://swingsociety.bg"}
    }

    return nil
}

// overrideWithEnv overrides configuration values with environment variables
func overrideWithEnv(config *Config) {
    if port := os.Getenv("PORT"); port != "" {
        config.Server.Port = port
    }
    if host := os.Getenv("HOST"); host != "" {
        config.Server.Host = host
    }
    if projectID := os.Getenv("GOOGLE_CLOUD_PROJECT"); projectID != "" {
        config.External.ProjectID = projectID
    }
    if sheetsURL := os.Getenv("GOOGLE_SHEETS_CSV_URL"); sheetsURL != "" {
        config.External.GoogleSheetsURL = sheetsURL
    }
    if env := os.Getenv("ENVIRONMENT"); env != "" {
        config.Environment = env
    }
}

// initializePaths sets up file paths for the application
func initializePaths(config *Config) error {
    // Determine project root
    config.Paths.RootDir = findProjectRoot()

    // Set paths relative to root
    config.Paths.TemplatesDir = filepath.Join(config.Paths.RootDir, "templates")
    config.Paths.StaticDir = filepath.Join(config.Paths.RootDir, "static")
    config.Paths.DataDir = filepath.Join(config.Paths.RootDir, "static/data")

    // Create directories if they don't exist (development mode only)
    if config.Environment == "development" {
        dirs := []string{
            config.Paths.TemplatesDir,
            config.Paths.StaticDir,
            config.Paths.DataDir,
        }
        for _, dir := range dirs {
            if err := os.MkdirAll(dir, 0755); err != nil {
                return fmt.Errorf("failed to create directory %s: %v", dir, err)
            }
        }
    }

    return nil
}

// validateConfig ensures required fields are set
func validateConfig(config *Config) error {
    if config.Server.Port == "" {
        return fmt.Errorf("server port is required")
    }
    if config.Paths.RootDir == "" {
        return fmt.Errorf("root directory is required")
    }
    return nil
}

// findProjectRoot dynamically determines the root directory
func findProjectRoot() string {
    cwd, err := os.Getwd()
    if err != nil {
        log.Fatalf("Error getting current directory: %v", err)
    }

    for {
        if _, err := os.Stat(filepath.Join(cwd, "go.mod")); err == nil {
            return cwd
        }
        parent := filepath.Dir(cwd)
        if parent == cwd { // Reached filesystem root
            break
        }
        cwd = parent
    }

    log.Fatalf("Could not determine project root. Make sure 'go.mod' exists in the root directory.")
    return ""
}
