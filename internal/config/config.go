package config

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



