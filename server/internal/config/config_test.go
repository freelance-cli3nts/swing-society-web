// internal/config/config_test.go
package config

import (
    "os"
    "path/filepath"
		"testing"
)

// setupTestConfig creates a temporary config file for testing
func setupTestConfig(t *testing.T) (string, func()) {
    dir := t.TempDir()
    configPath := filepath.Join(dir, "config.json")
    
    configData := []byte(`{
        "server": {
            "port": "3001",
            "host": "localhost",
            "read_timeout": 15000000000,
            "write_timeout": 15000000000
        },
        "paths": {
            "templates_dir": "templates",
            "static_dir": "static",
            "data_dir": "data"
        },
        "security": {
            "allowed_origins": ["*"],
            "rate_limits": {
                "api": "100-M",
                "static": "1000-M",
                "default": "60-M"
            }
        },
        "external": {
            "project_id": "test-project"
        },
        "environment": "development"
    }`)

    err := os.WriteFile(configPath, configData, 0644)
    if err != nil {
        t.Fatalf("Failed to create test config file: %v", err)
    }

    return configPath, func() {
        os.RemoveAll(dir)
    }
}


func TestLoadConfig(t *testing.T) {
    configPath, cleanup := setupTestConfig(t)
    defer cleanup()

    // Set the config path for testing
    os.Setenv("CONFIG_PATH", configPath)
    defer os.Unsetenv("CONFIG_PATH")

    // Test normal config loading
    err := LoadConfig()
    if err != nil {
        t.Errorf("LoadConfig() failed: %v", err)
    }

    // Test config values after loading
    if AppConfig == nil {
        t.Fatal("AppConfig is nil after LoadConfig()")
    }

    // Test environment variable override
    os.Setenv("PORT", "8080")
    defer os.Unsetenv("PORT")
    
    err = LoadConfig()
    if err != nil {
        t.Errorf("LoadConfig() with PORT env failed: %v", err)
    }

    if AppConfig.Server.Port != "8080" {
        t.Errorf("Port environment variable not applied, expected 8080, got %s", AppConfig.Server.Port)
    }
}

func TestLoadConfigWithInvalidPath(t *testing.T) {
    // Test with invalid config path
    os.Setenv("CONFIG_PATH", "nonexistent/config.json")
    defer os.Unsetenv("CONFIG_PATH")
    
    if err := LoadConfig(); err == nil {
        t.Error("Expected error with invalid config path, got nil")
    }
}

func TestConfigWithProductionEnv(t *testing.T) {
    configPath, cleanup := setupTestConfig(t)
    defer cleanup()

    // Set up the test environment
    os.Setenv("CONFIG_PATH", configPath)
    defer os.Unsetenv("CONFIG_PATH")
    
    os.Setenv("ENVIRONMENT", "production")
    defer os.Unsetenv("ENVIRONMENT")
    
    err := LoadConfig()
    if err != nil {
        t.Errorf("LoadConfig() in production mode failed: %v", err)
    }

    // Verify production-specific settings
    if AppConfig.Server.Host != "0.0.0.0" {
        t.Errorf("Expected production host to be 0.0.0.0, got %s", AppConfig.Server.Host)
    }
}

func TestValidateConfig(t *testing.T) {
    t.Run("Empty Port", func(t *testing.T) {
        config := &Config{}
        config.Paths.RootDir = "test"  // Set required root dir
        if err := validateConfig(config); err == nil {
            t.Error("Expected error for empty port, got nil")
        }
    })

    t.Run("Empty Root Dir", func(t *testing.T) {
        config := &Config{}
        config.Server.Port = "3001"  // Set required port
        if err := validateConfig(config); err == nil {
            t.Error("Expected error for empty root dir, got nil")
        }
    })

    t.Run("Valid Config", func(t *testing.T) {
        config := &Config{}
        config.Server.Port = "3001"
        config.Paths.RootDir = "test"
        if err := validateConfig(config); err != nil {
            t.Errorf("Unexpected error for valid config: %v", err)
        }
    })
}