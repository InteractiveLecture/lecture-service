package lecturepatch

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPatch(t *testing.T) {

	patchStr := `{"version" : 1, "lecture_id":"123","operations": [{"op": "add", "path": "/biscuits/1", "value": {"name": "Ginger Nut"}}, {"op": "remove", "path": "/biscuits"}]}`

	reader := strings.NewReader(patchStr)
	patch, err := Decode(reader)
	assert.Nil(t, err)
	assert.Equal(t, ADD, patch.Operations[0].Type)
	assert.Equal(t, 1, patch.Version)
	assert.Equal(t, "/biscuits/1", patch.Operations[0].Path)
	assert.Equal(t, 2, len(patch.Operations))

	assert.Equal(t, "123", patch.LectureID)
}
