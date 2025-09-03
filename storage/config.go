package storage

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	OpenAIKey string `yaml:"openai_api_key"`
}

var (
	configPath string = ""
	config     Config = Config{}
)

func ConfigFilePath() (string, error) {
	base, err := os.UserConfigDir() // honors XDG_CONFIG_HOME
	if err != nil {
		return "", err
	}
	dir := filepath.Join(base, "aidoc")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	configPath = filepath.Join(dir, "config.yaml")
	return configPath, nil
}

func LoadConfig() (*Config, error) {
	if v := os.Getenv("OPENAI_API_KEY"); v != "" {
		config.OpenAIKey = v
		return &config, nil
	}

	path, err := ConfigFilePath()
	if err != nil {
		return nil, err
	}
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &config, nil
		}
		return nil, err
	}
	if err := yaml.Unmarshal(b, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func SaveConfig() error {
	b, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	path, err := ConfigFilePath()
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func SaveAPIKey(key string) error {
	config.OpenAIKey = key
	return SaveConfig()
}
