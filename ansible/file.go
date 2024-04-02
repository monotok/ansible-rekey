package ansible

import (
	"ansible-rekey/common"
	"bytes"
	"gopkg.in/yaml.v3"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type Executor interface {
	Run(currentPassword, newPassword string, yml map[string]yaml.Node) []byte
}

func Walk(root, currentPassword, newPassword string, e Executor) error {
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
		result := e.Run(currentPassword, newPassword, yml)
		err = os.WriteFile(path, result, 0644)
		if err != nil {
			log.Fatal(err)
		}
		return nil
	})
}

func openAndParseFile(path string) (map[string]yaml.Node, error) {
	file, err := os.Open(path)
	defer file.Close()
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
	if bytes.HasPrefix(content, common.VaultHeader) {
		return nil, nil
	}

	yml := make(map[string]yaml.Node)
	err = yaml.Unmarshal(content, yml)
	return yml, err
}
