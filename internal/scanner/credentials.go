package scanner

import (
	"fmt"
	"log"
	"os"

	"github.com/LSariol/coveclient"
	"github.com/joho/godotenv"
)

// Pulls client secret from environment variables
func getClientSecret() string {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("LightHouse - getClientSecret: Error loading .env file")
	}
	return os.Getenv("COVE_CLIENT_SECRET")
}

func loadGitCredentials() (string, error) {
	client := coveclient.New("http://localhost:8081", getClientSecret())

	// gitToken, err := client.GetSecret("LIGHTHOUSE_GITHUB_PAT")
	gitToken, err := client.GetSecret("LIGHTHOUSE_GITHUB_PAT")
	if err != nil {
		return "", fmt.Errorf("loadGitCredentials failed: %v", err)
	}

	return gitToken, nil
}
