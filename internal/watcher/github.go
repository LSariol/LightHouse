package watcher

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func (w *Watcher) getLatestSHA(URL string, PAT string) (string, error) {

	req, err := http.NewRequest("GET", URL+"/commits", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "token "+PAT)
	// req.Header.Set("Accept", "application/vdn.github+json")

	resp, err := w.HTTP.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("GitHub API Error: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	var commits []map[string]interface{}
	err = json.Unmarshal(body, &commits)
	if err != nil {
		log.Fatal("Failed to unmarshal:", err)
	}
	shaHash := commits[0]["sha"].(string)
	return shaHash, nil
}
