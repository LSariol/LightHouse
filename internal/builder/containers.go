package builder

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func startCoveContainer(projectName string) error {

	// Step 3: Build volume paths
	homeDir, _ := os.UserHomeDir()
	envPath := filepath.Join(homeDir, "Server", "Storage", projectName, ".env")
	vaultPath := filepath.Join(homeDir, "Server", "Storage", projectName, "vault.json")

	// Step 4: Start the container
	startCMD := exec.Command("docker", "run",
		"-d", // Detached mode
		"-p", "8081:8081",
		"--name", projectName,
		"-v", vaultPath+":/internal/encryption/vault.json",
		"-v", envPath+":/.env",
		projectName,
	)

	if err := startCMD.Run(); err != nil {
		return fmt.Errorf("containers.startCoveContainer: error starting container: %v", err)
	}

	// Change to the target build directory
	if err := os.Chdir(OriginalPath); err != nil {
		return fmt.Errorf("failed to change to build directory: %v", err)
	}

	return nil
}
