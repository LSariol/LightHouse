package builder

import (
	"fmt"
	"os"
	"strings"

	"github.com/LSariol/LightHouse/internal/docker"
	"github.com/LSariol/LightHouse/internal/models"
)

type Builder struct {
	Paths         Paths
	DockerHandler *docker.Handler
}

type Paths struct {
	Original string
	Download string
	Staging  string
}

func NewBuilder(dh *docker.Handler) *Builder {
	return &Builder{
		DockerHandler: dh,
	}
}

func Build(repo models.WatchedRepo) (models.WatchedRepo, error) {

	fmt.Println("----Building " + repo.Name + " ----")

	err := cleanUp()
	if err != nil {
		wError := "Cleaning Failed for " + repo.Name + " " + err.Error()
		fmt.Println(wError)
		return repo, err
	}

	// Prepare Repo for build
	err = downloadNewCommit(repo.DownloadURL, repo.Name)
	if err != nil {
		wError := "Download Failed for " + repo.Name + " " + err.Error()
		fmt.Println(wError)
		return models.UpdateDownloadStats(repo, wError), err
	}
	repo = models.UpdateDownloadStats(repo, "Success")

	fmt.Println(stopContainer(strings.ToLower(repo.Name)))

	err = unpackNewProject(repo.Name)
	if err != nil {
		wError := "Unzip Failed for " + repo.Name + " " + err.Error()
		fmt.Println(wError)
		return repo, err
	}

	err = buildDockerImage(strings.ToLower(repo.Name))
	if err != nil {
		return repo, err
	}

	repo = models.UpdateBuildStats(repo, "Success")

	err = startDockerContainer(strings.ToLower(repo.Name))
	if err != nil {
		return repo, err
	}

	err = cleanUp()
	if err != nil {
		wError := "Cleaning Failed for " + repo.Name + " " + err.Error()
		fmt.Println(wError)
		return repo, err
	}

	fmt.Println("Clean Complete")

	return repo, nil
}

func ErrorHandler() {

}

// Run containers if they already exist.
func InitilizeContainers(watchList []models.WatchedRepo) error {

	for _, model := range watchList {
		// If container is running, good
		containerName := strings.ToLower(model.Name)
		status, err := getContainerStatus(containerName)
		if err != nil {
			return err
		}
		if status == "true" {
			fmt.Println(containerName + " is already running.")
			return nil
		}

		err = startDockerContainer(containerName)
		if err != nil {
			return err
		}
	}

	return nil
}

func InitilizeOriginalPath() string {
	originalPath, _ := os.Getwd()

	return originalPath
}
