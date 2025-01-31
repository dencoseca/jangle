package help

import (
	"fmt"
	"os"
)

var MainUsage = `
Usage: 	jangle COMMAND [ARGS]

Manage Mac OS Keychain secrets

Commands:
  set                           Set a secret with a name and value
  get                           Get a secret value by name
  update                        Update a secret with a name and value
  delete                        Delete a secret by name
  ls                            List all jangle secrets

Run 'jangle <command> --help' for more information.`

var SetUsage = `
Usage: 	jangle set <NAME> <VALUE>

Set a secret with a name and value

Example:
  jangle set SECRET_TOKEN 50m3t0k3n

Run 'jangle --help' for more information.`

var GetUsage = `
Usage: 	jangle get <NAME>

Get a secret value by name

Example:
  jangle get SECRET_TOKEN

Run 'jangle --help' for more information.`

var UpdateUsage = `
Usage: 	jangle get <NAME>

Update a secret with a name and value

Example:
  jangle update SECRET_TOKEN 50m3t0k3n

Run 'jangle --help' for more information.`

var DeleteUsage = `
Usage: 	jangle delete <NAME>

Delete a secret by name

Example:
  jangle delete SECRET_TOKEN

Run 'jangle --help' for more information.`

var ListUsage = `
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

// PrintHelpAndExit prints the usage information for a given command and exits
// the program with a specified exit code.
func PrintHelpAndExit(usage Usage, code ...int) {
	exitCode := 0
	if len(code) > 0 {
		exitCode = code[0]
	}

	output := ""
	switch usage {
	case Main:
		output = MainUsage
	case Set:
		output = SetUsage
	case Get:
		output = GetUsage
	case Update:
		output = UpdateUsage
	case Delete:
		output = DeleteUsage
	case List:
		output = ListUsage
	}

	fmt.Println(output)
	os.Exit(exitCode)
}
