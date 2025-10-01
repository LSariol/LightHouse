package scanner

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/LSariol/LightHouse/internal/models"
)

var WatchList []models.WatchedRepo

func AddNewRepo(name string, url string) ([]models.WatchedRepo, error) {
	var watchList []models.WatchedRepo = readWatchList()

	// Check if new repo is already being watched
	exists, conflictVal, conflictType := repoExists(watchList, name, url)
	if exists {
		fmt.Printf("Conflict detected: %s is already in use as a %s\n", conflictVal, conflictType)
		return watchList, nil
	}

	apiURL, downloadURL, err := getURLs(url)
	if err != nil {
		return watchList, err
	}

	newRepo := models.NewWatchedRepo(name, url, apiURL, downloadURL)

	watchList = append(watchList, newRepo)

	storeWatchList(watchList)

	return watchList, nil

}

func RemoveRepo(toRemove string) []models.WatchedRepo {
	watchList := readWatchList()

	indexToRemove := -1

	for index, existingRepo := range watchList {
		if existingRepo.Name == toRemove {
			indexToRemove = index
			break
		}
	}

	if indexToRemove != -1 {
		watchList = append(watchList[:indexToRemove], watchList[indexToRemove+1:]...)
		fmt.Println("Watcher - RemoveFromWatchList: " + toRemove + " has been removed.")
		storeWatchList(watchList)
		return watchList
	}
	fmt.Println("Watcher - RemoveFromWatchList: Unable to remove " + toRemove + ".")

	return watchList

}

func changeRepoName(currentName string, name string) ([]models.WatchedRepo, error) {
	watchList := readWatchList()
	updated := false

	if checkNamingConflicts(name, currentName, watchList) {
		return watchList, fmt.Errorf("This name is already being used to watch a different repo.")
	}

	for i := range watchList {
		if watchList[i].Name == currentName {
			watchList[i].Name = name
			lastModified := time.Now()
			watchList[i].Stats.Meta.LastModifiedAt = &lastModified
			updated = true
			break
		}
	}

	if !updated {
		return nil, fmt.Errorf("Lighthouse - UpdateRepo: Unable to update " + currentName + ". Does not exist.")
	}

	storeWatchList(watchList)

	return watchList, nil
}

func changeRepoURL(name string, newURL string) ([]models.WatchedRepo, error) {
	watchList := readWatchList()
	updated := false

	if checkURLConflicts(name, newURL, watchList) {
		return watchList, fmt.Errorf("This url is already being watched under a different name.")
	}

	for i := range watchList {
		if watchList[i].Name == name {

			apiURL, downloadURL, err := getURLs(newURL)
			if err != nil {
				return watchList, fmt.Errorf("changeRepoURL: %v", err)
			}

			watchList[i].URL = newURL
			lastModified := time.Now()
			watchList[i].Stats.Meta.LastModifiedAt = &lastModified
			watchList[i].APIURL = apiURL
			watchList[i].DownloadURL = downloadURL
			updated = true
			break
		}
	}

	if !updated {
		return nil, fmt.Errorf("Lighthouse - UpdateRepo: Unable to change the URL of " + name + ". Does not exist.")
	}

	storeWatchList(watchList)

	return watchList, nil
}

//Helper Functions

// Returns a boolean if repo exists
func repoExists(watchList []models.WatchedRepo, newRepoName string, newRepoURL string) (bool, string, string) {

	for _, existingRepo := range watchList {
		if existingRepo.URL == newRepoURL {
			return true, "URL", existingRepo.URL
		}

		if existingRepo.Name == newRepoName {
			return true, "Name", existingRepo.Name
		}
	}

	return false, "", ""

}

func checkNamingConflicts(name string, currentName string, watchList []models.WatchedRepo) bool {

	for _, repo := range watchList {
		if repo.Name == name && repo.Name != currentName {
			return true
		}
	}

	return false
}

func checkURLConflicts(name string, currentURL string, watchList []models.WatchedRepo) bool {

	for _, repo := range watchList {
		if repo.URL == currentURL && repo.Name != name {
			return true
		}
	}

	return false
}

func readWatchList() []models.WatchedRepo {
	var watchList []models.WatchedRepo

	//read json file
	data, err := os.ReadFile("config/repos.json")
	if err != nil {
		fmt.Println("Watcher - LoadWatchList: Failed to load repos.json")
		fmt.Println(err)
	}

	//Unmarshal repos.json into watchList
	err = json.Unmarshal([]byte(data), &watchList)
	if err != nil {
		fmt.Println("Watcher - LoadWatchList: Failed to Unmarshal json into WatchedRepos.")
		fmt.Println(err)
	}

	return watchList
}

func storeWatchList(watchList []models.WatchedRepo) {

	updatedData, err := json.MarshalIndent(watchList, "", "	")
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
func DisplayWatchList() {

	fmt.Printf("%-20s | %-40s | %-20s | %-15s\n", "Name", "URL", "Started Watching", "Query Count")
	fmt.Println(strings.Repeat("-", 20) + "-+-" + strings.Repeat("-", 40) + "-+-" + strings.Repeat("-", 20) + "-+-" + strings.Repeat("-", 15))

	for _, repo := range WatchList {
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
