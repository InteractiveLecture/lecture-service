package modulepatch

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPatch(t *testing.T) {

	patchStr := `{"lock_date" : 1448216365638, "operations": [{"op": "add", "path": "/biscuits/1", "value": {"name": "Ginger Nut"}}, {"op": "remove", "path": "/biscuits"}]}`

	reader := strings.NewReader(patchStr)
	patch, err := Decode(reader)
	assert.Nil(t, err)
	assert.Equal(t, ADD, patch.Operations[0].Type)
	assert.Equal(t, uint64(1448216365638), patch.LockDate)
	assert.Equal(t, "/biscuits/1", patch.Operations[0].Path)
	assert.Equal(t, 2, len(patch.Operations))
}
