package ansible

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"os"
	"testing"
)

type MockExecutor struct {
	mock.Mock
	executedFiles []string
}

func (ye *MockExecutor) Run(_, _ string, node map[string]yaml.Node) []byte {
	for _, v := range node {
		ye.executedFiles = append(ye.executedFiles, v.Value)
	}
	return []byte{}
}

func TestFile(t *testing.T) {
	t.Run("marshal file into yaml object successfully", func(t *testing.T) {
		var buffer bytes.Buffer
		buffer.WriteString("my_val: value")
		content, err := ParseFile(&buffer)

		expected := map[string]yaml.Node{"my_val": {
			Kind:        0x8,
			Style:       0x0,
			Tag:         "!!str",
			Value:       "value",
			Anchor:      "",
			Alias:       (*yaml.Node)(nil),
			Content:     []*yaml.Node(nil),
			HeadComment: "",
			LineComment: "",
			FootComment: "",
			Line:        1,
			Column:      9}}

		require.NoError(t, err)
		assert.Equal(t, expected, content)
	})

	t.Run("walk will execute a function on each parsed yaml file", func(t *testing.T) {
		tmpDir := t.TempDir()
		fileNames := []string{tmpDir + "/values.yml", tmpDir + "/values2.yml"}

		err := os.WriteFile(fileNames[0], []byte("my_var: value"), 0777)
		require.NoError(t, err)

		err = os.WriteFile(fileNames[1], []byte("my_var: value2"), 0777)
		require.NoError(t, err)

		mockYamlEditor := MockExecutor{}
		err = Walk(tmpDir, "", "", &mockYamlEditor)
		require.NoError(t, err)
		assert.Equal(t, mockYamlEditor.executedFiles[0], "value")
		assert.Equal(t, mockYamlEditor.executedFiles[1], "value2")

	})

	t.Run("walking the directory should ignore non yaml files", func(t *testing.T) {
		tmpDir := t.TempDir()

		err := os.WriteFile(tmpDir+"/values.txt", []byte("my_var: value"), 777)
		require.NoError(t, err)

		mockYamlEditor := MockExecutor{}
		err = Walk(tmpDir, "", "", &mockYamlEditor)
		require.NoError(t, err)
		assert.Empty(t, mockYamlEditor.executedFiles)

	})

	t.Run("walking the directory should ignore hidden directories", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", ".DS_STORE")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		err = os.WriteFile(tmpDir+"/values.txt", []byte("my_var: value"), 777)
		require.NoError(t, err)

		mockYamlEditor := MockExecutor{}
		err = Walk(tmpDir, "", "", &mockYamlEditor)
		require.NoError(t, err)
		assert.Empty(t, mockYamlEditor.executedFiles)

	})
}

func TestFile_Errors(t *testing.T) {
	t.Run("marshal invalid yaml file into yaml object throws error", func(t *testing.T) {
		var buffer bytes.Buffer
		buffer.WriteString("my_val value")
		content, err := ParseFile(&buffer)

		require.Error(t, err)
		assert.Empty(t, content)
	})
}
