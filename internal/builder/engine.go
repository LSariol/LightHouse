package builder

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/LSariol/LightHouse/internal/models"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var DownloadPath = "C:/Users/Lu/Server/Download"
var StagingPath = "C:/Users/Lu/Server/Staging"
var OriginalPath = ""

func downloadNewCommit(URL string, projectName string) error {

	fmt.Println("Downloading " + projectName)

	resp, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = os.MkdirAll(filepath.Join(DownloadPath), 0755)
	if err != nil {
		return err
	}

	out, err := os.Create(filepath.Join(DownloadPath, projectName+".zip"))
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func stopContainer(projectName string) error {

	stopCMD := exec.Command("docker", "stop", projectName)
	output, err := stopCMD.CombinedOutput()

	if err != nil {
		return fmt.Errorf("failed to stop container '%s': %v\nOutput: %s", projectName, err, string(output))
	}

	rmCMD := exec.Command("docker", "rm", projectName)
	output, err = rmCMD.CombinedOutput()

	if err != nil {
		return fmt.Errorf("failed to remove container '%s': %v\nOutput: %s", projectName, err, string(output))
	}

	fmt.Printf("Container '%s' has stopped and has been successfully deleted.\n", projectName)
	return nil
}

func unpackNewProject(projectName string) error {

	r, err := zip.OpenReader(filepath.Join(DownloadPath, projectName+".zip"))
	if err != nil {
		return err
	}
	defer r.Close()

	for _, file := range r.File {
		filePath := filepath.Join(StagingPath, file.Name)

		// Check for zip slip (Check for malicious files)
		if !strings.HasPrefix(filePath, filepath.Clean(StagingPath)+string(os.PathSeparator)) {
			return os.ErrPermission
		}

		if file.FileInfo().IsDir() {
			err := os.MkdirAll(filePath, os.ModePerm)
			if err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		outFile, err := os.Create(filePath)
		if err != nil {
			rc.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}

	}
	return nil
}

func buildDockerImage(projectName string) error {

	ctx := context.Background()

	buildPath := filepath.Join(StagingPath, projectName+"-main")

	tarBuf := new(bytes.Buffer)
	tw := tar.NewWriter(tarBuf)
	err := filepath.Walk(buildPath, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Open files for reading
		if fi.Mode().IsRegular() {
			relPath, err := filepath.Rel(buildPath, file)
			if err != nil {
				return err
			}
			header, err := tar.FileInfoHeader(fi, "")
			if err != nil {
				return err
			}
			header.Name = relPath

			if err := tw.WriteHeader(header); err != nil {
				return err
			}

			f, err := os.Open(file)
			if err != nil {
				return err
			}
			defer f.Close()

			if _, err := io.Copy(tw, f); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to create build tar: %v", err)
	}

	if err := tw.Close(); err != nil {
		return err
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("failed to create docker client: %v", err)
	}

	buildOpts := types.ImageBuildOptions{
		Tags:       []string{projectName},
		Remove:     true,
		Dockerfile: "Dockerfile",
	}

	buildResp, err := cli.ImageBuild(ctx, bytes.NewReader(tarBuf.Bytes()), buildOpts)
	if err != nil {
		return fmt.Errorf("failed to build docker image: %v", err)
	}
	defer buildResp.Body.Close()

	_, err = io.Copy(os.Stdout, buildResp.Body)
	if err != nil {
		return fmt.Errorf("failed to read build response: %v", err)
	}

	fmt.Printf("Docker image for %s built successfully.\n", projectName)

	return nil
}

func startDockerContainer(projectName string) error {
	cmd := exec.Command("docker", "compose", "up", "-d", "--build", "--remove-orphans")
	cmd.Dir = "C:/Users/Lu/Server/Staging/Cove-main" // <- relative paths in compose resolve here
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	return cmd.Run()
}

func cleanUp() error {

	// Clean Staging Folder
	err := os.RemoveAll(StagingPath)
	if err != nil {
		return fmt.Errorf("failed to clean staging area at %s: %w", StagingPath, err)
	}
	err = os.MkdirAll(StagingPath+"/Working", 0755)
	if err != nil {
		return fmt.Errorf("failed to recreate staging area at %s: %w", StagingPath, err)
	}

	//Clean Download Folder
	err = os.RemoveAll(DownloadPath)
	if err != nil {
		return fmt.Errorf("failed to clean download area at %s: %w", DownloadPath, err)
	}
	err = os.MkdirAll(DownloadPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to recreate download area at %s: %w", DownloadPath, err)
	}

	return nil
}

func getContainerStatus(containerName string) (string, error) {

	statusCMD := exec.Command("docker", "inspect", "-f", "{{.State.Running}}", containerName)
	output, err := statusCMD.CombinedOutput()
	if err != nil {
		return "false", nil
	}

	return strings.TrimSpace(string(output)), nil
}

func StopAllContainers(WatchList []models.WatchedRepo) error {

	for _, repo := range WatchList {

		projectName := strings.ToLower(repo.Name)
		if err := stopContainer(projectName); err != nil {
			return fmt.Errorf("error stopping %s: %v", projectName, err)
		}
	}

	return nil
}
