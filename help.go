package main

const mainUsage = `
Usage: 	jangle COMMAND [ARGS]

Manage Mac OS Keychain secrets

Commands:
  set                           Set a secret with a name and value
  get                           Get a secret value by name
  update                        Update a secret with a name and value
  delete                        Delete a secret by name
  ls                            List all jangle secrets

Run 'jangle <command> --help' for more information.`

const setUsage = `
Usage: 	jangle set <NAME> <VALUE>

Set a secret with a name and value

Example:
  jangle set SECRET_TOKEN 50m3t0k3n

Run 'jangle --help' for more information.`

const getUsage = `
Usage: 	jangle get <NAME>

Get a secret value by name

Example:
  jangle get SECRET_TOKEN

Run 'jangle --help' for more information.`

const updateUsage = `
Usage: 	jangle get <NAME>

Update a secret with a name and value

Example:
  jangle update SECRET_TOKEN 50m3t0k3n

Run 'jangle --help' for more information.`

const deleteUsage = `
Usage: 	jangle delete <NAME>

Delete a secret by name

Example:
  jangle delete SECRET_TOKEN

Run 'jangle --help' for more information.`

const listUsage = `
Usage: 	jangle ls

List all jangle secrets

Example:
  jangle ls

Run 'jangle --help' for more information.`
