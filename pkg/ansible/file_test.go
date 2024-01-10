package ansible

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"os"
	"testing"
)

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
		err := os.WriteFile(tmpDir+"/values.yml", []byte("my_var: value"), 777)
		require.NoError(t, err)

		var executedFiles []string
		err = walk(tmpDir, func(path string, yml map[string]yaml.Node) {
			executedFiles = append(executedFiles, path)
		})
		require.NoError(t, err)
		assert.Equal(t, executedFiles[0]+"/values.yml", tmpDir+"/values.yml")

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
