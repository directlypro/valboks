package config

import (
	"fmt"
	"os"
)

type Config string {
	AppKey string `json:"app_key"`
	AppSecret string `json:"app_secret"`
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

//this handles loading and saving file configs
type ConfigManager struct {
	configPath string
	config *Config
}

func NewConfigManager() (*ConfigManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("Error getting home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".config", "valboks-cli")
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		return nil, fmt.Errorf("error creating config directory: %w", err)
	}
}

