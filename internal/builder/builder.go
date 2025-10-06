package builder

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/LSariol/LightHouse/internal/models"
	"github.com/docker/docker/client"
)

type Builder struct {
	Docker    *client.Client
	Ctx       context.Context
	WatchList []models.WatchedRepo
	BasePath  string
}

func NewBuilder(dh *client.Client, ctx context.Context) *Builder {
	return &Builder{
		Docker: dh,
		Ctx:    ctx,
	}
}

func (b *Builder) Build(repo models.WatchedRepo) error {

	fmt.Println("----Building " + repo.Name + " ----")

	err := cleanUp()
	if err != nil {
		wError := "Cleaning Failed for " + repo.Name + " " + err.Error()
		fmt.Println(wError)
		return err
	}

	// Prepare Repo for build
	err = downloadNewCommit(repo.DownloadURL, repo.Name)
	if err != nil {
		wError := "Download Failed for " + repo.Name + " " + err.Error()
		return fmt.Errorf(wError)
	}

	err = b.StopContainer(repo.Name)
	if err != nil {
		return fmt.Errorf("build failed to stop container: %w", err)
	}

	err = unpackNewProject(repo.Name)
	if err != nil {
		wError := "Unzip Failed for " + repo.Name + " " + err.Error()
		fmt.Println(wError)
		return err
	}

	err = b.createContainer(strings.ToLower(repo.Name))
	if err != nil {
		return err
	}

	err = cleanUp()
	if err != nil {
		wError := "Cleaning Failed for " + repo.Name + " " + err.Error()
		fmt.Println(wError)
		return err
	}

	fmt.Println("Clean Complete")

	return nil
}

func ErrorHandler() {

}

// Run containers if they already exist.
func (b *Builder) InitilizeContainers(watchList []models.WatchedRepo) error {

	for _, model := range watchList {
		// If container is running, good
		containerName := strings.ToLower(model.Name)
		status, err := b.IsContainerRunning(containerName)
		if err != nil {
			return err
		}
		if status {
			fmt.Println(containerName + " is already running.")
			return nil
		}

		err = b.StartContainer(containerName)
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

func (b *Builder) LoadPaths() error {

	b.BasePath = os.Getenv("BASE_PATH")

	return nil
}
