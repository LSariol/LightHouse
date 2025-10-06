package watcher

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/LSariol/LightHouse/internal/models"
)

func (w *Watcher) addNewRepo(name string, url string) error {

	// Check if new repo is already being watched
	exists, conflictVal, conflictType := w.repoExists(name, url)
	if exists {
		fmt.Printf("Conflict detected: %s is already in use as a %s\n", conflictVal, conflictType)
		return nil
	}

	apiURL, downloadURL, err := getURLs(url)
	if err != nil {
		return err
	}

	newRepo := models.NewWatchedRepo(name, url, apiURL, downloadURL)

	w.WatchList = append(w.WatchList, newRepo)

	w.storeWatchList()

	return nil

}

func (w *Watcher) removeRepo(toRemove string) {

	indexToRemove := -1

	for index, existingRepo := range w.WatchList {
		if existingRepo.Name == toRemove {
			indexToRemove = index
			break
		}
	}

	if indexToRemove != -1 {
		w.WatchList = append(w.WatchList[:indexToRemove], w.WatchList[indexToRemove+1:]...)
		fmt.Println("Watcher - RemoveFromWatchList: " + toRemove + " has been removed.")
		w.storeWatchList()
		return
	}

	fmt.Println("Watcher - RemoveFromWatchList: Unable to remove " + toRemove + ".")
	return

}

func (w *Watcher) changeRepoName(currentName string, name string) error {
	updated := false

	if w.checkNamingConflicts(name, currentName) {
		return fmt.Errorf("This name is already being used to watch a different repo.")
	}

	for i := range w.WatchList {
		if w.WatchList[i].Name == currentName {
			w.WatchList[i].Name = name
			lastModified := time.Now()
			w.WatchList[i].Stats.Meta.LastModifiedAt = &lastModified
			updated = true
			break
		}
	}

	if !updated {
		return fmt.Errorf("changeRepoName: " + currentName + " does not exist.")
	}

	w.storeWatchList()

	return nil
}

func (w *Watcher) changeRepoURL(name string, newURL string) error {
	updated := false

	if w.checkURLConflicts(name, newURL) {
		return fmt.Errorf("This url is already being watched under a different name.")
	}

	for i := range w.WatchList {
		if w.WatchList[i].Name == name {

			apiURL, downloadURL, err := getURLs(newURL)
			if err != nil {
				return fmt.Errorf("changeRepoURL: %w", err)
			}

			w.WatchList[i].URL = newURL
			lastModified := time.Now()
			w.WatchList[i].Stats.Meta.LastModifiedAt = &lastModified
			w.WatchList[i].APIURL = apiURL
			w.WatchList[i].DownloadURL = downloadURL
			updated = true
			break
		}
	}

	if !updated {
		return fmt.Errorf("changeRepoURL: " + name + " does not exist.")
	}

	w.storeWatchList()

	return nil
}

//Helper Functions

// Returns a boolean if repo exists
func (w *Watcher) repoExists(newRepoName string, newRepoURL string) (bool, string, string) {

	for _, existingRepo := range w.WatchList {
		if existingRepo.URL == newRepoURL {
			return true, "URL", existingRepo.URL
		}

		if existingRepo.Name == newRepoName {
			return true, "Name", existingRepo.Name
		}
	}

	return false, "", ""

}

func (w *Watcher) checkNamingConflicts(name string, currentName string) bool {

	for _, repo := range w.WatchList {
		if repo.Name == name && repo.Name != currentName {
			return true
		}
	}

	return false
}

func (w *Watcher) checkURLConflicts(name string, currentURL string) bool {

	for _, repo := range w.WatchList {
		if repo.URL == currentURL && repo.Name != name {
			return true
		}
	}

	return false
}

func (w *Watcher) loadWatchList() error {
	var watchList []models.WatchedRepo

	//read json file
	data, err := os.ReadFile("config/repos.json")
	if err != nil {
		fmt.Println("Watcher - LoadWatchList: Failed to load repos.json")
		return fmt.Errorf("loadWatchList: %w", err)
	}

	//Unmarshal repos.json into watchList
	err = json.Unmarshal([]byte(data), &watchList)
	if err != nil {
		fmt.Println("Watcher - LoadWatchList: Failed to Unmarshal json into WatchedRepos.")
		return fmt.Errorf("loadWatchList: %w", err)
	}

	w.WatchList = watchList
	return nil
}

func (w *Watcher) storeWatchList() {

	updatedData, err := json.MarshalIndent(w.WatchList, "", "	")
	if err != nil {
		fmt.Println("Watcher - storeWatchList: Failed to Marhsal json into UpdatedData.")
		return
	}

	err = os.WriteFile("config/repos.json", updatedData, 0644)
	if err != nil {
		fmt.Println("Watcher - storeWatchList: Failed to write to repos.json." + err.Error())
		return
	}
}

// Display WatchList in a nice format
func (w *Watcher) DisplayWatchList() {

	fmt.Printf("%-20s | %-40s | %-20s | %-15s\n", "Name", "URL", "Started Watching", "Query Count")
	fmt.Println(strings.Repeat("-", 20) + "-+-" + strings.Repeat("-", 40) + "-+-" + strings.Repeat("-", 20) + "-+-" + strings.Repeat("-", 15))

	for _, repo := range w.WatchList {
		fmt.Printf(
			"%-20s | %-40s | %-20s | %-15d \n",
			repo.Name,
			repo.URL,
			repo.Stats.Meta.StartedWatchingAt.Format("2006-01-02 15:04:05"),
			repo.Stats.Queries.QueryCount,
		)
	}
}

func getURLs(url string) (string, string, error) {

	trim := strings.TrimPrefix(url, "https://github.com/")
	parts := strings.Split(trim, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("Invalid Github repo URL: %s", url)
	}

	return "https://api.github.com/repos/" + parts[0] + "/" + parts[1], "https://github.com/" + parts[0] + "/" + parts[1] + "/archive/refs/heads/main.zip", nil

}
