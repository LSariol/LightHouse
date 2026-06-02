package cli

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/LSariol/LightHouse/internal/watcher"
)

type CLI struct {
	Watcher *watcher.Watcher
}

func NewCLI(w *watcher.Watcher) *CLI {
	return &CLI{
		Watcher: w,
	}
}

func (c *CLI) Run() {
	ioScanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("LightHouse CLI> ")
		if !ioScanner.Scan() {
			break
		}
		input := ioScanner.Text()
		c.parseCLI(strings.Fields(input))
	}
}

func (c *CLI) parseCLI(args []string) {
	if len(args) == 0 {
		return
	}

	switch args[0] {

	case "add", "a":
		if len(args) != 3 {
			fmt.Println("add requires 3 total arguments.")
			fmt.Println("add <DisplayName> <repoURL>")
			return
		}
		err := c.Watcher.AddNewRepo(args[1], args[2])
		if err != nil {
			fmt.Printf("Failed adding new repo: %v\n", err)
			return
		}
		log.Printf("%s is now being watched.\n", args[1])

	case "remove", "r":
		if len(args) != 2 {
			fmt.Println("remove requires 2 total arguments.")
			fmt.Println("remove <repoName>")
			return
		}
		err := c.Watcher.RemoveRepo(args[1])
		if err != nil {
			fmt.Printf("Failed removing repo: %v\n", err)
		}

	case "change", "c":
		if len(args) != 3 {
			fmt.Println("change requires 3 total arguments.")
			fmt.Println("change <repoName> <newURL>")
			return
		}
		err := c.Watcher.UpdateRepo(args[1], args[2])
		if err != nil {
			fmt.Printf("Failed to update repository URL: %v\n", err)
			return
		}
		fmt.Println("URL has been updated.")

	case "update", "u":
		if len(args) < 3 {
			fmt.Println("update requires 3 total arguments.")
			fmt.Println("update url <repoName> <newURL>")
			fmt.Println("update name <repoName> <newName>")
			return
		}
		switch args[1] {
		case "url", "URL":
			if len(args) != 4 {
				fmt.Println("update url requires 4 total arguments.")
				fmt.Println("update url <repoName> <newURL>")
				return
			}
			err := c.Watcher.UpdateRepo(args[2], args[3])
			if err != nil {
				fmt.Printf("Failed to update URL: %v\n", err)
				return
			}
			fmt.Println("URL has been updated.")
		case "name", "NAME":
			if len(args) != 4 {
				fmt.Println("update name requires 4 total arguments.")
				fmt.Println("update name <currentName> <newName>")
				return
			}
			err := c.Watcher.ChangeRepoName(args[2], args[3])
			if err != nil {
				fmt.Printf("Failed to change name: %v\n", err)
				return
			}
			fmt.Println("Name has been updated.")
		default:
			fmt.Printf("Unknown update subcommand: %q\n", args[1])
			fmt.Println("Use: update url | update name")
		}

	case "start", "START":
		if len(args) != 2 {
			fmt.Println("start requires 2 total arguments.")
			fmt.Println("start <projectName|ALL>")
			return
		}
		if strings.EqualFold(args[1], "all") {
			if err := c.Watcher.Builder.StartAllContainers(); err != nil {
				fmt.Printf("Error starting all containers: %v\n", err)
				return
			}
			fmt.Println("All containers started.")
			return
		}
		if err := c.Watcher.Builder.StartContainer(args[1]); err != nil {
			fmt.Printf("Error starting %q: %v\n", args[1], err)
			return
		}
		fmt.Printf("%s has been started.\n", args[1])

	case "stop", "STOP":
		if len(args) != 2 {
			fmt.Println("stop requires 2 total arguments.")
			fmt.Println("stop <projectName|ALL>")
			return
		}
		if strings.EqualFold(args[1], "all") {
			if err := c.Watcher.Builder.StopAllContainers(); err != nil {
				fmt.Printf("Error stopping all containers: %v\n", err)
				return
			}
			fmt.Println("All containers stopped.")
			return
		}
		if err := c.Watcher.Builder.StopContainer(args[1]); err != nil {
			fmt.Printf("Error stopping %q: %v\n", args[1], err)
			return
		}
		fmt.Printf("%s has been stopped.\n", args[1])

	case "restart", "RESTART":
		if len(args) != 2 {
			fmt.Println("restart requires 2 total arguments.")
			fmt.Println("restart <projectName|ALL>")
			return
		}
		if strings.EqualFold(args[1], "all") {
			for _, repo := range c.Watcher.WatchList {
				if err := c.Watcher.Builder.RestartContainer(strings.ToLower(repo.ContainerName)); err != nil {
					fmt.Printf("Error restarting %q: %v\n", repo.ContainerName, err)
				}
			}
			fmt.Println("All containers restarted.")
			return
		}
		if err := c.Watcher.Builder.RestartContainer(args[1]); err != nil {
			fmt.Printf("Error restarting %q: %v\n", args[1], err)
			return
		}
		fmt.Printf("%s has been restarted.\n", args[1])

	case "rebuild", "REBUILD":
		if len(args) != 2 {
			fmt.Println("rebuild requires 2 total arguments.")
			fmt.Println("rebuild <projectName|ALL>")
			return
		}
		if strings.EqualFold(args[1], "all") {
			for _, repo := range c.Watcher.WatchList {
				fmt.Printf("Rebuilding %s...\n", repo.DisplayName)
				if err := c.Watcher.Builder.Build(repo); err != nil {
					fmt.Printf("Rebuild failed for %q: %v\n", repo.DisplayName, err)
				}
			}
			fmt.Println("All rebuilds complete.")
			return
		}
		repo, ok := c.Watcher.FindRepo(args[1])
		if !ok {
			fmt.Printf("No watched repo named %q.\n", args[1])
			return
		}
		fmt.Printf("Rebuilding %s...\n", repo.DisplayName)
		if err := c.Watcher.Builder.Build(repo); err != nil {
			fmt.Printf("Rebuild failed: %v\n", err)
			return
		}
		fmt.Printf("%s has been rebuilt and is running.\n", repo.DisplayName)

	case "logs", "LOGS":
		if len(args) < 2 || len(args) > 3 {
			fmt.Println("logs requires 2 or 3 total arguments.")
			fmt.Println("logs <projectName> [lines]")
			return
		}
		tail := 50
		if len(args) == 3 {
			n, err := strconv.Atoi(args[2])
			if err != nil || n <= 0 {
				fmt.Printf("Invalid line count %q — must be a positive integer.\n", args[2])
				return
			}
			tail = n
		}
		output, err := c.Watcher.Builder.GetContainerLogs(args[1], tail)
		if err != nil {
			fmt.Printf("Error fetching logs for %q: %v\n", args[1], err)
			return
		}
		fmt.Print(output)

	case "status", "STATUS":
		fmt.Printf("%-20s | %-20s | %s\n", "Name", "Container", "Status")
		fmt.Println(strings.Repeat("-", 20) + "-+-" + strings.Repeat("-", 20) + "-+-" + strings.Repeat("-", 10))
		for _, repo := range c.Watcher.WatchList {
			running, err := c.Watcher.Builder.IsContainerRunning(strings.ToLower(repo.ContainerName))
			state := "Stopped"
			if err != nil {
				state = "Unknown"
			} else if running {
				state = "Running"
			}
			fmt.Printf("%-20s | %-20s | %s\n", repo.DisplayName, repo.ContainerName, state)
		}

	case "pause", "PAUSE":
		if c.Watcher.IsPaused() {
			fmt.Println("Auto-scan is already paused.")
			return
		}
		c.Watcher.Pause()
		fmt.Println("Auto-scan paused. Use 'resume' to restart.")

	case "resume", "RESUME":
		if !c.Watcher.IsPaused() {
			fmt.Println("Auto-scan is not paused.")
			return
		}
		c.Watcher.Resume()
		fmt.Println("Auto-scan resumed.")

	case "scan", "SCAN":
		c.Watcher.Scan()

	case "list", "LIST", "l", "L":
		c.Watcher.DisplayWatchList()

	case "help", "h":
		printHelp()

	case "exit", "quit", "q":
		if len(args) == 1 {
			fmt.Println("Shutting down LightHouse...")
			os.Exit(0)
			return
		}
		switch args[1] {
		case "all", "a":
			fmt.Println("Shutting down LightHouse and all containers...")
			if err := c.Watcher.Builder.StopAllContainers(); err != nil {
				fmt.Printf("Error while shutting down containers: %v\n", err)
			}
			os.Exit(0)
		}

	default:
		fmt.Printf("Unknown command: %q\n", args[0])
		printHelp()
	}
}

func printHelp() {
	fmt.Println()
	fmt.Println("LightHouse CLI — available commands:")
	fmt.Println()
	fmt.Printf("  %-35s %s\n", "list  |  l", "List all watched repositories and stats")
	fmt.Printf("  %-35s %s\n", "status", "Show running state of all watched containers")
	fmt.Printf("  %-35s %s\n", "add <name> <url>", "Add a GitHub repository to the watchlist")
	fmt.Printf("  %-35s %s\n", "remove <name>", "Remove a repository from the watchlist")
	fmt.Printf("  %-35s %s\n", "change <name> <new-url>", "Update a repository's URL")
	fmt.Printf("  %-35s %s\n", "update url <name> <new-url>", "Update a repository's URL")
	fmt.Printf("  %-35s %s\n", "update name <name> <new-name>", "Rename a repository")
	fmt.Printf("  %-35s %s\n", "start <name|ALL>", "Start a container (or all containers)")
	fmt.Printf("  %-35s %s\n", "stop <name|ALL>", "Stop a container (or all containers)")
	fmt.Printf("  %-35s %s\n", "restart <name|ALL>", "Restart a container (or all containers)")
	fmt.Printf("  %-35s %s\n", "rebuild <name|ALL>", "Force a full repull and rebuild")
	fmt.Printf("  %-35s %s\n", "logs <name> [lines]", "Print last N lines of container logs (default 50)")
	fmt.Printf("  %-35s %s\n", "scan", "Manually trigger one scan cycle")
	fmt.Printf("  %-35s %s\n", "pause", "Pause the automatic scan loop")
	fmt.Printf("  %-35s %s\n", "resume", "Resume the automatic scan loop")
	fmt.Printf("  %-35s %s\n", "help  |  h", "Show this help message")
	fmt.Printf("  %-35s %s\n", "exit [all]", "Shut down LightHouse (exit all also stops containers)")
	fmt.Println()
}
