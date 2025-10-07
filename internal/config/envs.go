package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

func Load() (string, error) {

	if err := godotenv.Load(".env"); err == nil {
		return ".env", nil
	}

	if err := godotenv.Load("/app/vault/.env"); err == nil {
		return "/app/vault/.env", nil
	}

	return "", fmt.Errorf("no .env file found")
}

func SaveClientSecret(envPath string, clientSecret string) error {

	var myEnv map[string]string
	myEnv, err := godotenv.Read(envPath)
	if err != nil {
		return fmt.Errorf("unable to read .env file for save")
	}

	var needWrite bool = false
	if myEnv["APP_ENV_PATH"] == "" {
		myEnv["APP_ENV_PATH"] = envPath
		needWrite = true
	}

	if myEnv["COVE_CLIENT_SECRET"] == "" {
		myEnv["COVE_CLIENT_SECRET"] = clientSecret
		needWrite = true
	}

	if needWrite {
		err := godotenv.Write(myEnv, envPath)
		if err != nil {
			return fmt.Errorf("unable to write .env: %w", err)
		}
	}

	godotenv.Load(envPath)
	return nil
}
