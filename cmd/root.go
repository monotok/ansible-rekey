package cmd

import (
	"github.com/spf13/cobra"
)

type Cli interface {
	Rekey(ansibleDirectory, vaultFile string) error
}

func NewCliCommand(cli Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use: "ansible-utils <command>",
	}
	cmd.AddCommand(
		newRekeyCommand(cli),
	)
	return cmd
}
