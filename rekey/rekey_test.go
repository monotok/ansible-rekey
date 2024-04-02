package rekey

import (
	"ansible-rekey/common"
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"testing"
)

func Test_reEncryptVariable(t *testing.T) {
	const encryptedString = "$ANSIBLE_VAULT;1.1;AES256\n33393235313365663038333137376534353965353930383635356663663065353566333138323938\n6637613830653039363632383434343162326662653666650a333864653534343763323062666566\n36346565623261303530346635623430393431356563386665623436623531333164316130326331\n3433373132386330370a303536313664363939313133616163393339656561666265623365613533\n66346632346663633836386363643163646566303437633761366637346361613034\n"

	t.Run("re-encrypt an encrypted yaml object with a new password", func(t *testing.T) {
		password := "pa$$word"
		newPassword := "strongNewpa$$word"

		encryptedNode := map[string]yaml.Node{"my-var": {
			Kind:        0x8,
			Style:       0x9,
			Tag:         "!vault",
			Value:       encryptedString,
			Anchor:      "",
			Alias:       nil,
			Content:     nil,
			HeadComment: "",
			LineComment: "",
			FootComment: "",
			Line:        1,
			Column:      11,
		},
			"anotherNode": {
				Kind:        0x8,
				Style:       0x9,
				Tag:         "!vault",
				Value:       encryptedString,
				Anchor:      "",
				Alias:       nil,
				Content:     nil,
				HeadComment: "",
				LineComment: "",
				FootComment: "",
				Line:        1,
				Column:      11,
			}}
		execute := Execute{}
		result := execute.Run(password, newPassword, encryptedNode)
		newlyEncryptedNode := map[string]yaml.Node{}
		err := yaml.Unmarshal(result, &newlyEncryptedNode)
		require.NoError(t, err)

		assert.True(t, bytes.HasPrefix([]byte(newlyEncryptedNode["my-var"].Value), common.VaultHeader))
		assert.True(t, bytes.HasPrefix([]byte(newlyEncryptedNode["anotherNode"].Value), common.VaultHeader))
		assert.NotEqual(t, encryptedString, newlyEncryptedNode["my-var"].Value)
		assert.NotEqual(t, encryptedString, newlyEncryptedNode["anotherNode"].Value)
	})
}
