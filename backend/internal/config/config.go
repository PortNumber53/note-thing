package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"gopkg.in/ini.v1"
)

const defaultConfigPath = "/etc/note-thing/config.ini"

// Load tries to load environment variables from `.env` and then from an INI file.
// Environment variables already set take precedence over the INI values.
// If the INI file does not exist, Load is a no-op.
func Load() error {
	_ = godotenv.Load()

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = defaultConfigPath
	}

	absoluteConfigPath, err := filepath.Abs(configPath)
	if err != nil {
		absoluteConfigPath = configPath
	}

	if _, err := os.Stat(absoluteConfigPath); err != nil {
		return nil
	}

	iniFile, err := ini.Load(absoluteConfigPath)
	if err != nil {
		return err
	}

	section := iniFile.Section("")
	setIfMissing("DATABASE_URL", section.Key("DATABASE_URL").String())
	setIfMissing("GOOGLE_CLIENT_ID", section.Key("GOOGLE_CLIENT_ID").String())
	setIfMissing("GOOGLE_CLIENT_SECRET", section.Key("GOOGLE_CLIENT_SECRET").String())
	setIfMissing("JWT_SECRET", section.Key("JWT_SECRET").String())
	setIfMissing("PORT", section.Key("PORT").String())

	return nil
}

func setIfMissing(key string, value string) {
	if value == "" {
		return
	}
	if os.Getenv(key) != "" {
		return
	}
	_ = os.Setenv(key, value)
}
