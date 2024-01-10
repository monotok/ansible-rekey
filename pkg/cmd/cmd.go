package cmd

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
)

type CommandExecutor interface {
	Execute(vaultFile string) error
}

type DefaultCommandExecutor struct{}

func (ce *DefaultCommandExecutor) Execute(vaultFile string) error {
	dir := "."
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}

	contents, err := os.ReadFile(vaultFile)
	if err != nil {
		return err
	}
	password := strings.TrimSuffix(string(contents), "\n")

	fmt.Println("DIRECTORY " + dir)
	fmt.Println("Old password: " + password)
	fmt.Println("Enter the new password: ")
	pass, _ := terminal.ReadPassword(0)
	fmt.Println("You entered " + string(pass))

	return nil
}
