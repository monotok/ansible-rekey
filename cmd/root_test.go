package cmd

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"strings"
	"testing"
)

type MockCmdExecutor struct {
	mock.Mock
}

func (ce *MockCmdExecutor) Execute(vaultFile string) error {
	args := ce.Called(vaultFile)
	return args.Error(0)
}

func Test_ExecuteCommand(t *testing.T) {
	tests := map[string]struct {
		args             string
		expectedCmdFlags any
		returnArgs       any
	}{
		"valid short hand args": {
			args:             "-v .vault",
			expectedCmdFlags: ".vault",
		},
		"valid long hand args": {
			args:             "--vault .vault",
			expectedCmdFlags: ".vault",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			bStdOut, bStdErr := bytes.NewBufferString(""), bytes.NewBufferString("")

			rootCmd := setupRootCmd(strings.Split(tc.args, " "), bStdOut, bStdErr)

			mockCmdExecutor := new(MockCmdExecutor)
			mockCmdExecutor.On("Execute", tc.expectedCmdFlags).Return(tc.returnArgs)
			defer mockCmdExecutor.AssertExpectations(t)

			rootCmd.WithExecutor(mockCmdExecutor)

			err := rootCmd.Execute()
			require.NoError(t, err)

			stdOut, err := io.ReadAll(bStdOut)
			require.NoError(t, err)

			stdErr, err := io.ReadAll(bStdErr)
			require.NoError(t, err)

			assert.Empty(t, string(stdOut))
			assert.Empty(t, string(stdErr))

			mockCmdExecutor.AssertExpectations(t)
		})
	}
}

func Test_ExecuteCommand_Errors(t *testing.T) {
	defaultMockBehaviourSetup := func(m *MockCmdExecutor) {}
	executeNotCalledAssertion := func(m *MockCmdExecutor) {
		m.AssertNotCalled(t, "Execute")
	}

	tests := map[string]struct {
		args                   string
		expectedStdErr         string
		mockBehaviourSetup     func(m *MockCmdExecutor)
		mockBehaviourAssertion func(m *MockCmdExecutor)
	}{
		"no flags supplied": {
			args:                   "",
			expectedStdErr:         `Error: required flag(s) "vault" not set`,
			mockBehaviourSetup:     defaultMockBehaviourSetup,
			mockBehaviourAssertion: executeNotCalledAssertion,
		},
		"unrecognised flag": {
			args:                   "-v .vault --foo",
			expectedStdErr:         `Error: unknown flag: --foo`,
			mockBehaviourSetup:     defaultMockBehaviourSetup,
			mockBehaviourAssertion: executeNotCalledAssertion,
		},
		"vault flag supplied but with no value": {
			args:                   "-v",
			expectedStdErr:         `Error: flag needs an argument: 'v' in -v`,
			mockBehaviourSetup:     defaultMockBehaviourSetup,
			mockBehaviourAssertion: executeNotCalledAssertion,
		},
		"valid args - execution error": {
			args:           "-v .vault",
			expectedStdErr: `Error: command execution error`,
			mockBehaviourSetup: func(m *MockCmdExecutor) {
				m.On("Execute", ".vault").Return(errors.New("command execution error"))
			},
			mockBehaviourAssertion: func(m *MockCmdExecutor) {
				m.AssertCalled(t, "Execute", ".vault")
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			bStdOut, bStdErr := bytes.NewBufferString(""), bytes.NewBufferString("")

			rootCmd := setupRootCmd(strings.Split(tc.args, " "), bStdOut, bStdErr)

			mockCmdExecutor := new(MockCmdExecutor)
			tc.mockBehaviourSetup(mockCmdExecutor)
			defer tc.mockBehaviourAssertion(mockCmdExecutor)

			rootCmd.WithExecutor(mockCmdExecutor)

			err := rootCmd.Execute()
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

func setupRootCmd(args []string, stdOut io.Writer, stdErr io.Writer) *RootCmd {
	rootCmd := NewRootCmd()

	rootCmd.SetOut(stdOut)
	rootCmd.SetErr(stdErr)
	rootCmd.SetArgs(args)

	return rootCmd
}
