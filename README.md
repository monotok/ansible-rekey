# ansible-utils

**STATUS**

Very early beta. Use at your own risk.

A collection of ansible utility functions. Currently only supports `rekey` which recursively rekeys encrypted variables within a directory.

## Install

```
go install github.com/monotok/ansible-utils@latest
```

## Usage

```
Usage:
  ansible-utils [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  rekey       Easily rekey all encrypted string variables within an ansible project

Flags:
  -h, --help   help for ansible-utils

Use "ansible-utils [command] --help" for more information about a command.

```

### Example

```console
$ ansible-find rekey testdata -v ~/ansible-project/.vault
```

