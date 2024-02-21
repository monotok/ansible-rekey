package ansible

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

var (
	vaultHeader = []byte("$ANSIBLE_VAULT;1.1;AES256")
)

type YamlEditor interface {
	run(path string, yml map[string]yaml.Node)
}

type Yaml struct{}

func (y Yaml) run(path string, yml map[string]yaml.Node) {
	fmt.Println(path)
	fmt.Println(yaml.Marshal(yml))
}

func Walk(root string, y YamlEditor) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && d.Name() != root && strings.HasPrefix(d.Name(), ".") {
			return filepath.SkipDir
		}

		if !slices.Contains([]string{".yaml", ".yml"}, filepath.Ext(d.Name())) {
			return nil
		}

		yml, err := openAndParseFile(path)
		if err != nil {
			return err
		}
		if yml == nil {
			return nil
		}
		y.run(path, yml)
		return nil
	})
}

func openAndParseFile(path string) (map[string]yaml.Node, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return ParseFile(file)
}

func ParseFile(file io.Reader) (map[string]yaml.Node, error) {
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	if bytes.HasPrefix(content, vaultHeader) {
		return nil, nil
	}

	yml := make(map[string]yaml.Node)
	err = yaml.Unmarshal(content, &yml)
	return yml, err
}
