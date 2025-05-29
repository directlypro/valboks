package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
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

	configPath := filepath.Join(configDir, "config.json")

	return &ConfigManager{
		configPath: configPath,
		config: &Config{},
	}, nil
}

func (m *ConfigManager) Load() error {
	if _, err := os.Stat(m.configPath); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(m.configPath)
	if err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	err = json.Unmarshal(data, m.config)
	if err != nil {
		return fmt.Errorf("error parsing config fiel: %w", err)
	}

	return nil
}

func (m *ConfigManager) Save() error {
	data, err := json.MarshalIndent(m.config, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializing config: %w", err)
	}

	err = os.WriteFile(m.configPath, data, 0600)
	if err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
}

func (m *ConfigManager) GetConfig() *Config {
	return m.config
}

func (m *ConfigManager) SetCredentials(appKey, appSecret, accessToken string) {
	m.config.AppKey = appKey
	m.config.AppSecret = appSecret
	m.config.AccessToken = accessToken
}

func (m *ConfigManager) SetTokens(accessToken, refreshToken string) {
	m.config.AccessToken = accessToken
	m.config.RefreshToken = refreshToken
}

func (m *ConfigManager) IsConfigured() bool {
	return m.config.AccessToken != ""
}