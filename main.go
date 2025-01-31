package main

import (
	"fmt"
	"github.com/dencoseca/jangle/help"
	"github.com/dencoseca/jangle/stores"
	"github.com/dencoseca/jangle/styles"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(help.MainUsage)
		os.Exit(1)
	}

	command := os.Args[1]

	name := getArg(2)
	value := getArg(3)

	store := stores.NewMacOSKeychainStore()

	switch command {
	case "--help", "help":
		fmt.Println(help.MainUsage)
	case "get":
		if name == "" {
			fmt.Println(help.GetUsage)
			os.Exit(1)
		}
		value, err := store.Get(name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Print(value)
	case "set":
		if name == "" || value == "" {
			fmt.Println(help.SetUsage)
			os.Exit(1)
		}
		err := store.Set(name, value)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		styles.Green("Successfully added '%s'.\n", name)
		fmt.Println("Source your terminal configuration or restart your shell to use the environment variable.")
	case "update":
		if name == "" || value == "" {
			fmt.Println(help.UpdateUsage)
			os.Exit(1)
		}
		err := store.Update(name, value)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		styles.Green("Successfully updated '%s'.\n", name)
		fmt.Println("Source your terminal configuration or restart your shell to use the updated environment variable.")
	case "ls":
		if len(os.Args[2:]) > 0 {
			fmt.Println(help.ListUsage)
			os.Exit(1)
		}
		secretNames, err := store.List()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if len(secretNames) == 0 {
			styles.Red("no secrets found for user '%s'", os.Getenv("USER"))
			os.Exit(0)
		}
		fmt.Println(fmt.Sprintf("Secrets for user '%s':\n", os.Getenv("USER")))
		for _, s := range secretNames {
			fmt.Println("- " + s)
		}
	case "delete":
		if name == "" {
			fmt.Println(help.DeleteUsage)
			os.Exit(1)
		}
		err := store.Delete(name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		styles.Green("Successfully removed '%s'.\n", name)
		fmt.Println(fmt.Sprintf("To Delete the environment variable restart your terminal or run: unset %s", name))
	default:
		fmt.Println(help.MainUsage)
		os.Exit(1)
	}
}

func getArg(index int) string {
	if len(os.Args) > index {
		return os.Args[index]
	}

	return ""
}
