package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/LSariol/LightHouse/internal/builder"
	"github.com/LSariol/LightHouse/internal/scanner"
)

func StartCLI() {

	ioScanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("LightHouse CLI> ")
		if !ioScanner.Scan() {
			break
		}
		input := ioScanner.Text()
		parseCLI(strings.Fields(input))
	}
}

func parseCLI(args []string) {

	if len(args) == 0 {
		return
	}

	switch args[0] {
	case "update", "u":

		switch args[1] {

		case "url", "URL":
			//do url update

		case "name", "NAME":
			//do name update

		}

	case "scan", "SCAN", "s", "S":
		scanner.Scan()
	case "list", "LIST", "l", "L":
		scanner.DisplayWatchList()

	case "exit", "quit", "q":

		if len(args) == 1 {
			fmt.Println("Shutting down Lighthouse...")
			os.Exit(0)
			return
		}

		switch args[1] {
		case "all", "a":
			fmt.Println("Shutting down Lighthouse and all containers...")
			if err := builder.StopAllContainers(scanner.WatchList); err != nil {
				fmt.Printf("Error while shutting down containers: %v", err)
			}
			os.Exit(0)
		}

	}

	// 	case "get", "g":

	// 		if len(args) != 2 {
	// 			yellowLog("Get requires 2 total arguments.")
	// 			yellowLog("get <secret>")
	// 			return
	// 		}

	// 		if !ok {
	// 			redLog(res + ": " + args[1])
	// 			return
	// 		}

	// 		greenLog("Secret has been retreived: " + res)

	// 	case "add", "a":

	// 		if len(args) != 3 {
	// 			yellowLog("add requires 3 total arguments.")
	// 			yellowLog("add <secretName> <value>")
	// 			return
	// 		}

	// 		res, ok := encryption.AddSecret(args[1], args[2])
	// 		if !ok {
	// 			redLog(res)
	// 		}

	// 		greenLog("Secret has been added")

	// 	case "remove", "r", "delete", "d":

	// 		if len(args) != 2 {
	// 			yellowLog("Get requires 2 total arguments.")
	// 			yellowLog("remove <secret>")
	// 			return
	// 		}

	// 		res, ok := encryption.RemoveSecret(args[1])
	// 		if !ok {
	// 			redLog(res)
	// 			return
	// 		}

	// 		greenLog("Secret has been removed")

	// 	case "update", "u":

	// 		if len(args) != 3 {
	// 			yellowLog("Update requires 3 total arguments.")
	// 			yellowLog("update <secretName> <newValue>")
	// 			return
	// 		}

	// 		res, ok := encryption.UpdateSecret(args[1], args[2])
	// 		if !ok {
	// 			redLog(res)
	// 			return
	// 		}

	// 		greenLog("Secret has been updated.")

	// 	case "list", "l":

	// 		if len(args) > 2 {
	// 			yellowLog("Update requires 1 or 2 arguments.")
	// 			yellowLog("list or list all")
	// 			return
	// 		}

	// 		if len(args) == 1 {
	// 			displayWatchedRepos()
	// 			return
	// 		}

	// 		if len(args) == 2 || args[2] == "all" {

	// 		}

	// 	}
	// }

	// func greenLog(s string) {
	// 	fmt.Println("\033[32mCove CLI> " + s + "\033[0m")
	// }

	// func yellowLog(s string) {
	// 	fmt.Println("\033[33mCove CLI> " + s + "\033[0m")
	// }

	// func redLog(s string) {
	// 	fmt.Println("\033[31mCove CLI> " + s + "\033[0m")
	// }
}
