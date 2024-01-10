package ansible

import (
	"bytes"
	"gopkg.in/yaml.v3"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var (
	vaultHeader = []byte("$ANSIBLE_VAULT;1.1;AES256")
)

func walk(root string, run func(path string, yml map[string]yaml.Node)) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && d.Name() != root && strings.HasPrefix(d.Name(), ".") {
			return filepath.SkipDir
		}

		yml, err := openAndParseFile(path)
		if yml == nil {
			return nil
		}
		if err != nil {
			return err
		}
		run(path, yml)
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
	if bytes.HasPrefix(content, vaultHeader) {
		return nil, nil
	}

	yml := make(map[string]yaml.Node)
	err = yaml.Unmarshal(content, &yml)
	return yml, err
}
