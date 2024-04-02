package rekey

import (
	"ansible-rekey/common"
	"bytes"
	log "github.com/sirupsen/logrus"
	vault "github.com/sosedoff/ansible-vault-go"
	"gopkg.in/yaml.v3"
)

type Execute struct{}

func (e Execute) Run(currentPassword, newPassword string, yml map[string]yaml.Node) []byte {
	for k := range yml {
		result, valid := reEncryptVariable(yml[k], currentPassword, newPassword)
		if valid {
			yml[k] = result
		}
	}
	yamlFile, err := yaml.Marshal(&yml)
	if err != nil {
		log.Fatal(err)
	}
	return yamlFile
}

func reEncryptVariable(v yaml.Node, password, newPassword string) (yaml.Node, bool) {
	if bytes.HasPrefix([]byte(v.Value), common.VaultHeader) {
		decrypted, err := vault.Decrypt(v.Value, password)
		if err != nil {
			log.Fatal(err)
		}
		encrypted, err := vault.Encrypt(decrypted, newPassword)
		if err != nil {
			log.Fatal(err)
		}
		v.SetString(encrypted)
		return v, true
	}
	return yaml.Node{}, false
}
