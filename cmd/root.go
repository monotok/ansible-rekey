/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"ansible-rekey/pkg/cmd"
	"github.com/spf13/cobra"
	"io"
	"log"
)

type cmdFlags struct {
	vaultFile string
}

type RootCmd struct {
	cobraCmd *cobra.Command
	config   *cmdFlags
	executor cmd.CommandExecutor
}

func NewRootCmd() *RootCmd {
	flags := cmdFlags{
		vaultFile: "",
	}

	cobraCmd := newCobraCommand()
	configureFlags(cobraCmd, &flags)

	rootCmd := &RootCmd{
		executor: &cmd.DefaultCommandExecutor{},
		config:   &flags,
		cobraCmd: cobraCmd,
	}

	cobraCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return rootCmd.executor.Execute(rootCmd.config.vaultFile)
	}

	return rootCmd
}

func (c *RootCmd) WithExecutor(executor cmd.CommandExecutor) *RootCmd {
	c.executor = executor
	return c
}

func (c *RootCmd) SetArgs(args []string) {
	c.cobraCmd.SetArgs(args)
}

func (c *RootCmd) SetOut(writer io.Writer) {
	c.cobraCmd.SetOut(writer)
}

func (c *RootCmd) SetErr(writer io.Writer) {
	c.cobraCmd.SetErr(writer)
}

func (c *RootCmd) Execute() error {
	return c.cobraCmd.Execute()
}

func configureFlags(command *cobra.Command, flags *cmdFlags) {
	command.Flags().StringVarP(&flags.vaultFile, "vault", "v", ".vault", "current ansible vault password file")
	err := command.MarkFlagRequired("vault")
	if err != nil {
		log.Fatalf("error marking flag as required: %v", err)
	}
}

func newCobraCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "ansible-rekey DIRECTORY",
		Short: "Easily rekey all encrypted string variables within an ansible project",
		Long: `This tool allows you to rekey all the encrypted strings within an ansible project.

It will search for and find any encrypted string and rekey to the new password.`,
	}
}
