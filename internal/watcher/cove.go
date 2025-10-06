package watcher

import (
	"fmt"
	"os"

	"github.com/LSariol/coveclient"
)

func (w *Watcher) loadGitCredentials() error {

	gitToken, err := w.CC.GetSecret("LIGHTHOUSE_GITHUB_PAT")
	if err != nil {
		return fmt.Errorf("loadGitCredentials: %v", err)
	}

	w.GitToken = gitToken
	return nil
}

func NewCoveClient() *coveclient.Client {

	clientSecret := os.Getenv("COVE_CLIENT_SECRET")
	var coveClient *coveclient.Client = coveclient.New("http://localhost:"+os.Getenv("COVE_PORT"), clientSecret)

	if clientSecret == "" {
		clientSecret, err := coveClient.Bootstrap()
		if err != nil {
			panic(err)
		}
		coveClient.ClientSecret = clientSecret
	}

	return coveClient
}
