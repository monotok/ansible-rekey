package cmd

import (
	"bytes"
	"errors"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"strings"
	"testing"
)

type MockCli struct {
	mock.Mock
}

func NewMockCli() *MockCli {
	return &MockCli{}
}

func (cli *MockCli) Rekey(ansibleDirectory, vaultFile string) error {
	args := cli.Called(ansibleDirectory, vaultFile)
	return args.Error(0)
}

func Test_ExecuteCommand(t *testing.T) {
	tests := map[string]struct {
		args             string
		expectedCmdFlags any
		expectedCmdArgs  any
		returnArgs       any
	}{
		"valid short hand args": {
			args:             "rekey mydir -v .vault",
			expectedCmdArgs:  "mydir",
			expectedCmdFlags: ".vault",
		},
		"valid long hand args": {
			args:             "rekey mydir --vault .vault",
			expectedCmdArgs:  "mydir",
			expectedCmdFlags: ".vault",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockCli := NewMockCli()
			cmd := NewCliCommand(mockCli)

			bStdOut, bStdErr := bytes.NewBufferString(""), bytes.NewBufferString("")

			setupRootCmd(cmd, strings.Split(tc.args, " "), bStdOut, bStdErr)

			mockCli.On("Rekey", tc.expectedCmdArgs, tc.expectedCmdFlags).Return(tc.returnArgs)
			defer mockCli.AssertExpectations(t)

			err := cmd.Execute()
			require.NoError(t, err)

			stdOut, err := io.ReadAll(bStdOut)
			require.NoError(t, err)

			stdErr, err := io.ReadAll(bStdErr)
			require.NoError(t, err)

			assert.Empty(t, string(stdOut))
			assert.Empty(t, string(stdErr))
		})
	}
}

func Test_ExecuteCommand_Errors(t *testing.T) {
	defaultMockBehaviourSetup := func(m *MockCli) {}
	executeNotCalledAssertion := func(m *MockCli) {
		m.AssertNotCalled(t, "Rekey")
	}

	tests := map[string]struct {
		args                   string
		expectedStdErr         string
		mockBehaviourSetup     func(m *MockCli)
		mockBehaviourAssertion func(m *MockCli)
	}{
		"no flags supplied": {
			args:                   "rekey mydir",
			expectedStdErr:         `Error: required flag(s) "vault" not set`,
			mockBehaviourSetup:     defaultMockBehaviourSetup,
			mockBehaviourAssertion: executeNotCalledAssertion,
		},
		"unrecognised flag": {
			args:                   "rekey mydir  -v .vault --foo",
			expectedStdErr:         `Error: unknown flag: --foo`,
			mockBehaviourSetup:     defaultMockBehaviourSetup,
			mockBehaviourAssertion: executeNotCalledAssertion,
		},
		"vault flag supplied but with no value": {
			args:                   "rekey mydir  -v",
			expectedStdErr:         `Error: flag needs an argument: 'v' in -v`,
			mockBehaviourSetup:     defaultMockBehaviourSetup,
			mockBehaviourAssertion: executeNotCalledAssertion,
		},
		"directory not supplied but flag provided": {
			args:                   "rekey -v .vault",
			expectedStdErr:         `Error: accepts 1 arg(s), received 0`,
			mockBehaviourSetup:     defaultMockBehaviourSetup,
			mockBehaviourAssertion: executeNotCalledAssertion,
		},
		"valid args - execution error": {
			args:           "rekey mydir -v .vault",
			expectedStdErr: `Error: command execution error`,
			mockBehaviourSetup: func(m *MockCli) {
				m.On("Rekey", "mydir", ".vault").Return(errors.New("command execution error"))
			},
			mockBehaviourAssertion: func(m *MockCli) {
				m.AssertCalled(t, "Rekey", "mydir", ".vault")
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockCli := NewMockCli()
			cmd := NewCliCommand(mockCli)

			bStdOut, bStdErr := bytes.NewBufferString(""), bytes.NewBufferString("")

			setupRootCmd(cmd, strings.Split(tc.args, " "), bStdOut, bStdErr)

			tc.mockBehaviourSetup(mockCli)
			defer tc.mockBehaviourAssertion(mockCli)

			err := cmd.Execute()
			require.Error(t, err)

			stdOut, err := io.ReadAll(bStdOut)
			require.NoError(t, err)

			stdErr, err := io.ReadAll(bStdErr)
			require.NoError(t, err)

			assert.True(t, strings.HasPrefix(string(stdOut), "Usage:"))
			assert.Equal(t, tc.expectedStdErr, strings.ReplaceAll(string(stdErr), "\n", ""))
		})
	}
}

func setupRootCmd(cmd *cobra.Command, args []string, stdOut io.Writer, stdErr io.Writer) {
	cmd.SetOut(stdOut)
	cmd.SetErr(stdErr)
	cmd.SetArgs(args)
}
