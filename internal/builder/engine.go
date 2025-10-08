package builder

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/errdefs"
)

func downloadNewCommit(URL string, projectName string) error {

	fmt.Println("Downloading " + projectName)

	resp, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = os.MkdirAll(filepath.Join(os.Getenv("DOWNLOAD_PATH")), 0755)
	if err != nil {
		return err
	}

	out, err := os.Create(filepath.Join(os.Getenv("DOWNLOAD_PATH"), projectName+".zip"))
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

func unpackNewProject(projectName string) error {

	r, err := zip.OpenReader(filepath.Join(os.Getenv("DOWNLOAD_PATH"), projectName+".zip"))
	if err != nil {
		return err
	}
	defer r.Close()

	for _, file := range r.File {
		filePath := filepath.Join(os.Getenv("STAGING_PATH"), file.Name)

		// Check for zip slip (Check for malicious files)
		if !strings.HasPrefix(filePath, filepath.Clean(os.Getenv("STAGING_PATH"))+string(os.PathSeparator)) {
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

func cleanUp() error {

	// Clean Staging Folder
	err := cleanupAll("STAGING_PATH")
	if err != nil {
		return fmt.Errorf("failed to clean staging area at %s: %w", os.Getenv("STAGING_PATH"), err)
	}
	err = os.MkdirAll(os.Getenv("STAGING_PATH")+"/Working", 0755)
	if err != nil {
		return fmt.Errorf("failed to recreate staging area at %s: %w", os.Getenv("STAGING_PATH"), err)
	}

	//Clean Download Folder
	err = cleanupAll("DOWNLOAD_PATH")
	if err != nil {
		return fmt.Errorf("failed to clean download area at %s: %w", os.Getenv("DOWNLOAD_PATH"), err)
	}
	err = os.MkdirAll(os.Getenv("DOWNLOAD_PATH"), 0755)
	if err != nil {
		return fmt.Errorf("failed to recreate download area at %s: %w", os.Getenv("DOWNLOAD_PATH"), err)
	}

	return nil
}

func cleanupAll(base string) error {
	base = os.Getenv(base)
	if base == "" {
		return fmt.Errorf("STAGING_PATH not set")
	}

	// Read all entries inside the staging path
	entries, err := os.ReadDir(base)
	if err != nil {
		return fmt.Errorf("failed to read staging dir: %w", err)
	}

	for _, e := range entries {
		p := filepath.Join(base, e.Name())
		if err := os.RemoveAll(p); err != nil {
			return fmt.Errorf("failed to remove %s: %w", p, err)
		}
	}

	return nil
}

func (b *Builder) createContainer(projectName string) error {

	cmd := exec.Command("docker", "compose", "up", "-d", "--build", "--remove-orphans")
	cmd.Dir = os.Getenv("STAGING_PATH") + projectName + "-main" // <- relative paths in compose resolve here
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr

	if projectName == "cloudflared" {
		token, err := b.CC.GetSecret("CLOUDFLARE_TUNNEL_TOKEN")
		if err != nil {
			return err
		}

		cmd.Env = append(os.Environ(), "TUNNEL_TOKEN="+token)
	}

	if projectName == "lighthousedb" {
		pg_db, err := b.CC.GetSecret("PG_DB")
		if err != nil {
			return err
		}

		pg_user, err := b.CC.GetSecret("PG_USER")
		if err != nil {
			return err
		}

		pg_password, err := b.CC.GetSecret("PG_PASSWORD")
		if err != nil {
			return err
		}
		cmd.Env = append(os.Environ(), "POSTGRES_DB="+pg_db, "POSTGRES_USER="+pg_user, "POSTGRES_PASSWORD="+pg_password)
	}

	return cmd.Run()
}

func (b *Builder) StartContainer(name string) error {
	return b.Docker.ContainerStart(b.Ctx, name, container.StartOptions{})
}

func (b *Builder) StopContainer(name string) error {
	return b.Docker.ContainerStop(b.Ctx, name, container.StopOptions{})
}

func (b *Builder) RestartContainer(name string) error {
	return b.Docker.ContainerRestart(b.Ctx, name, container.StopOptions{})
}

func (b *Builder) GetAllContainers() ([]types.Container, error) {

	containers, err := b.Docker.ContainerList(b.Ctx, container.ListOptions{
		All: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	return containers, nil
}

func (b *Builder) GetRunningContainers() ([]types.Container, error) {

	containers, err := b.Docker.ContainerList(b.Ctx, container.ListOptions{
		All: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	return containers, nil
}

func (b *Builder) IsContainerRunning(nameOrId string) (bool, error) {

	info, err := b.Docker.ContainerInspect(b.Ctx, nameOrId)
	if err != nil {
		if errdefs.IsNotFound(err) {
			return false, nil
		}
		return false, fmt.Errorf("inspect %q: %w", nameOrId, err)
	}

	if info.State == nil {
		return false, fmt.Errorf("no state for %q", nameOrId)
	}

	return info.State.Running, nil
}

func (b *Builder) StartAllContainers() error {

	for _, repo := range b.WatchList {
		name := strings.ToLower(repo.Name)

		err := b.StartContainer(name)
		if err != nil {
			return fmt.Errorf("starting all containers: %s: %w", name, err)
		}
	}

	return nil
}

func (b *Builder) StopAllContainers() error {

	for _, repo := range b.WatchList {
		name := strings.ToLower(repo.Name)

		err := b.StopContainer(name)
		if err != nil {
			return fmt.Errorf("starting all containers: %s: %w", name, err)
		}
	}

	return nil
}
