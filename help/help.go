package help

import (
	"fmt"
	"os"
)

const MainUsage = `
Usage: 	jangle COMMAND [ARGS]

Manage Mac OS Keychain secrets

Commands:
  set                           Set a secret with a name and value
  get                           Get a secret value by name
  update                        Update a secret with a name and value
  delete                        Delete a secret by name
  ls                            List all jangle secrets

Run 'jangle <command> --help' for more information.`

const SetUsage = `
Usage: 	jangle set <NAME> <VALUE>

Set a secret with a name and value

Example:
  jangle set SECRET_TOKEN 50m3t0k3n

Run 'jangle --help' for more information.`

const GetUsage = `
Usage: 	jangle get <NAME>

Get a secret value by name

Example:
  jangle get SECRET_TOKEN

Run 'jangle --help' for more information.`

const UpdateUsage = `
Usage: 	jangle get <NAME>

Update a secret with a name and value

Example:
  jangle update SECRET_TOKEN 50m3t0k3n

Run 'jangle --help' for more information.`

const DeleteUsage = `
Usage: 	jangle delete <NAME>

Delete a secret by name

Example:
  jangle delete SECRET_TOKEN

Run 'jangle --help' for more information.`

const ListUsage = `
Usage: 	jangle ls

List all jangle secrets

Example:
  jangle ls

Run 'jangle --help' for more information.`

type Usage int

const (
	Main Usage = iota
	Set
	Get
	Update
	Delete
	List
)

// PrintUsageAndExit prints the usage information for a given command and exits
// the program with a specified exit code.
func PrintUsageAndExit(usage Usage, code ...int) {
	exitCode := 0
	if len(code) > 0 {
		exitCode = code[0]
	}

	switch usage {
	case Main:
		fmt.Println(MainUsage)
		os.Exit(exitCode)
	case Set:
		fmt.Println(SetUsage)
		os.Exit(exitCode)
	case Get:
		fmt.Println(GetUsage)
		os.Exit(exitCode)
	case Update:
		fmt.Println(UpdateUsage)
		os.Exit(exitCode)
	case Delete:
		fmt.Println(DeleteUsage)
		os.Exit(exitCode)
	case List:
		fmt.Println(ListUsage)
		os.Exit(exitCode)
	}
}
