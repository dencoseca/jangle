package main

import (
	"fmt"
	"github.com/dencoseca/jangle/help"
	"github.com/dencoseca/jangle/stores"
	"github.com/dencoseca/jangle/styles"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		help.PrintHelpAndExit(help.Main, 1)
	}

	command := os.Args[1]

	name := getArgSafely(2)
	value := getArgSafely(3)

	store := stores.NewMacOSKeychainStore()
	exportFile := stores.NewExportFile()

	switch command {
	case "--help", "help":
		help.PrintHelpAndExit(help.Main)
	case "get":
		if name == "" {
			help.PrintHelpAndExit(help.Get, 1)
		}

		s, err := store.Get(name)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(s)
	case "set":
		if name == "" || value == "" {
			help.PrintHelpAndExit(help.Set, 1)
		}

		if err := store.Set(name, value); err != nil {
			log.Fatal(err)
		}

		if err := exportFile.Set(name); err != nil {
			log.Fatal(err)
		}

		styles.Green("Successfully added '%s'.\n", name)
		fmt.Println("Source your terminal configuration or restart your shell to use the environment variable.")
	case "update":
		if name == "" || value == "" {
			help.PrintHelpAndExit(help.Update, 1)
		}

		if err := store.Update(name, value); err != nil {
			log.Fatal(err)
		}

		styles.Green("Successfully updated '%s'.\n", name)
		fmt.Println("Source your terminal configuration or restart your shell to use the updated environment variable.")
	case "ls":
		if len(os.Args[2:]) > 0 {
			help.PrintHelpAndExit(help.List, 1)
		}

		secretNames, err := store.List()
		if err != nil {
			log.Fatal(err)
		}

		if len(secretNames) == 0 {
			styles.Red("No secrets found for user '%s'", os.Getenv("USER"))
			return
		}

		fmt.Println(fmt.Sprintf("Secrets for user '%s':\n", os.Getenv("USER")))
		for _, s := range secretNames {
			fmt.Println("- " + s)
		}
	case "delete":
		if name == "" {
			help.PrintHelpAndExit(help.Delete, 1)
		}

		if err := store.Delete(name); err != nil {
			log.Fatal(err)
		}

		if err := exportFile.Delete(name); err != nil {
			log.Fatal(err)
		}

		styles.Green("Successfully removed '%s'.\n", name)
		fmt.Println(fmt.Sprintf("To delete the environment variable restart your terminal or run: unset %s", name))
	default:
		help.PrintHelpAndExit(help.Main, 1)
	}
}

// getArgSafely retrieves the command-line argument at the specified index if it
// exists; otherwise, returns an empty string.
func getArgSafely(index int) string {
	if len(os.Args) > index {
		return os.Args[index]
	}

	return ""
}
