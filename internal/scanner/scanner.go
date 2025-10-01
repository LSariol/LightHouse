package scanner

import (
	"fmt"
	"time"

	"github.com/LSariol/LightHouse/internal/builder"
	"github.com/LSariol/LightHouse/internal/models"
)

var GithubToken string

func Initilize() error {

	var err error

	fmt.Println("Initilizing: WatchList")
	WatchList = readWatchList()

	fmt.Println("Initilizing: Builder HomePath")
	builder.InitilizeOriginalPath()

	fmt.Println("Initilizing: Containers")
	if err := builder.InitilizeContainers(WatchList); err != nil {
		return err
	}

	fmt.Println("Initilizing: Git Credentials")
	if GithubToken, err = loadGitCredentials(); err != nil {
		return err
	}

	return nil
}

func Run() {

	for {

		if err := Scan(); err != nil {
			fmt.Printf("ERROR IN SCAN: %v", err)
		}
		time.Sleep(10 * time.Second)

	}

}

func Scan() error {

	for i, repo := range WatchList {
		repo = WatchList[i]
		currentHash, err := getLatestSHA(repo.APIURL, GithubToken)
		if err != nil {

			repo = models.UpdateErrorStats(repo, err.Error())
			repo = models.UpdateQueryStats(repo)
			WatchList[i] = repo
			storeWatchList(WatchList)
			return fmt.Errorf("scanner.scan() - getLatestSha: %v", err)

		}

		if repo.Stats.Updates.LastSeenCommitSha == nil || *repo.Stats.Updates.LastSeenCommitSha != currentHash {

			repo = models.UpdateUpdateStats(repo, currentHash)
			repo, err = builder.Build(repo)
			if err != nil {
				builder.ErrorHandler()
				return fmt.Errorf("scanner.scan() - error in build: %v", err)
			}

		}

		repo = models.UpdateQueryStats(repo)
		WatchList[i] = repo
		storeWatchList(WatchList)

		continue
	}

	return nil

}
