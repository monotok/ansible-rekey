package cmd

import (
	"ansible-rekey/ansible"
	"ansible-rekey/rekey"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
	"strings"
)

type RekeyCli struct{}

func NewRekeyCli() *RekeyCli {
	return &RekeyCli{}
}

type cmdFlags struct {
	vaultFile string
}

func (cli *RekeyCli) Rekey(ansibleDirectory, vaultFile string) error {
	contents, err := os.ReadFile(vaultFile)
	if err != nil {
		return err
	}
	password := strings.TrimSuffix(string(contents), "\n")

	fmt.Println("DIRECTORY " + ansibleDirectory)
	fmt.Println("Old password: " + password)
	fmt.Println("Enter the new password: ")
	pass, _ := terminal.ReadPassword(0)
	e := rekey.Execute{}
	ansible.Walk(ansibleDirectory, password, string(pass), &e)

	return nil
}

func newRekeyCommand(cli Cli) *cobra.Command {
	flags := &cmdFlags{}
	cmd := &cobra.Command{
		Use:   "rekey DIRECTORY",
		Short: "Easily rekey all encrypted string variables within an ansible project",
		Long: `This tool allows you to rekey all the encrypted strings within an ansible project.

		It will search for and find any encrypted string and rekey to the new password.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := RunRekey(cli, args[0], flags.vaultFile); err != nil {
				return err
			}
			return nil
		},
	}
	configureFlags(cmd, flags)
	return cmd
}

func configureFlags(command *cobra.Command, flags *cmdFlags) {
	command.Flags().StringVarP(&flags.vaultFile, "vault", "v", ".vault", "current ansible vault password file")
	err := command.MarkFlagRequired("vault")
	if err != nil {
		log.Fatalf("error marking flag as required: %v", err)
	}
}

func RunRekey(cli Cli, ansibleDirectory, vaultFile string) error {
	err := cli.Rekey(ansibleDirectory, vaultFile)
	if err != nil {
		return err
	}
	return nil
}
