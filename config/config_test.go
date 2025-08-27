package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected *Config
	}{
		{
			name:    "default values",
			envVars: map[string]string{},
			expected: &Config{
				Port:      "3000",
				JWTSecret: "your-secret-key",
				DBPath:    "users.db",
			},
		},
		{
			name: "custom values from env",
			envVars: map[string]string{
				"PORT":       "8080",
				"JWT_SECRET": "super-secret-key",
				"DB_PATH":    "/tmp/test.db",
			},
			expected: &Config{
				Port:      "8080",
				JWTSecret: "super-secret-key",
				DBPath:    "/tmp/test.db",
			},
		},
		{
			name: "partial env values",
			envVars: map[string]string{
				"PORT": "9000",
			},
			expected: &Config{
				Port:      "9000",
				JWTSecret: "your-secret-key",
				DBPath:    "users.db",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment variables
			os.Unsetenv("PORT")
			os.Unsetenv("JWT_SECRET")
			os.Unsetenv("DB_PATH")

			// Set test environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			// Test Load function
			config := Load()

			// Verify results
			if config.Port != tt.expected.Port {
				t.Errorf("Port = %v, want %v", config.Port, tt.expected.Port)
			}
			if config.JWTSecret != tt.expected.JWTSecret {
				t.Errorf("JWTSecret = %v, want %v", config.JWTSecret, tt.expected.JWTSecret)
			}
			if config.DBPath != tt.expected.DBPath {
				t.Errorf("DBPath = %v, want %v", config.DBPath, tt.expected.DBPath)
			}

			// Clean up
			for key := range tt.envVars {
				os.Unsetenv(key)
			}
		})
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "env variable exists",
			key:          "TEST_KEY",
			defaultValue: "default",
			envValue:     "env_value",
			expected:     "env_value",
		},
		{
			name:         "env variable does not exist",
			key:          "NON_EXISTENT_KEY",
			defaultValue: "default_value",
			envValue:     "",
			expected:     "default_value",
		},
		{
			name:         "empty env variable",
			key:          "EMPTY_KEY",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up
			os.Unsetenv(tt.key)

			// Set environment variable if provided
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
			}

			// Test getEnv function
			result := getEnv(tt.key, tt.defaultValue)

			// Verify result
			if result != tt.expected {
				t.Errorf("getEnv(%s, %s) = %v, want %v", tt.key, tt.defaultValue, result, tt.expected)
			}

			// Clean up
			os.Unsetenv(tt.key)
		})
	}
}

func TestConfigStruct(t *testing.T) {
	config := &Config{
		Port:      "8080",
		JWTSecret: "test-secret",
		DBPath:    "test.db",
	}

	if config.Port != "8080" {
		t.Errorf("Port = %v, want %v", config.Port, "8080")
	}
	if config.JWTSecret != "test-secret" {
		t.Errorf("JWTSecret = %v, want %v", config.JWTSecret, "test-secret")
	}
	if config.DBPath != "test.db" {
		t.Errorf("DBPath = %v, want %v", config.DBPath, "test.db")
	}
}
