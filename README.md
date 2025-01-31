# Jangle

## Overview

Jangle is a command line tool to provide an easy way to manage secrets and export them in a shell environment.

## Getting started

### Prerequisites

Mac OS is required as the Mac OS Keychain is currently the only supported secret store.

Go SDK installed:

```shell
brew install go
```

GOPATH/bin folder in your system PATH variable:

```shell
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Installation

```shell
go install github.com/dencoseca/jangle@latest
```

### Usage

```shell
jangle --help
```
