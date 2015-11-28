package patchcompiler

import (
	"testing"

	"github.com/richterrettich/lecture-service/lecturepatch"
	"github.com/stretchr/testify/assert"
)

func createTopicPatch() *lecturepatch.Patch {

	patch := &lecturepatch.Patch{}

	patch.Version = 1
	patch.ModelID = "123"
	patch.Operations = []lecturepatch.Operation{
		lecturepatch.Operation{
			Type:  lecturepatch.REPLACE,
			Path:  "/description",
			Value: "Hallo Welt",
		},
		lecturepatch.Operation{
			Type:  lecturepatch.ADD,
			Path:  "/assistants",
			Value: "111",
		},
		lecturepatch.Operation{
			Type: lecturepatch.REMOVE,
			Path: "/assistants/111",
		},
		lecturepatch.Operation{
			Type: lecturepatch.ADD,
			Path: "/modules",
			Value: map[string]interface{}{
				"id":          "222",
				"description": "Hugo",
				"video_id":    "333",
				"script_id":   "444",
				"parents":     []string{"555", "666"},
			},
		},

		lecturepatch.Operation{
			Type:  lecturepatch.REPLACE,
			Path:  "/modules/111/parents",
			Value: []string{"555", "666"},
		},
		lecturepatch.Operation{
			Type:  lecturepatch.REPLACE,
			Path:  "/modules/111/parents/tree",
			Value: []string{"555", "666"},
		},

		lecturepatch.Operation{
			Type:  lecturepatch.REPLACE,
			Path:  "/modules/111/parents",
			Value: []string{"555", "666"},
		},

		lecturepatch.Operation{
			Type: lecturepatch.REMOVE,
			Path: "/modules/111",
		},

		lecturepatch.Operation{
			Type: lecturepatch.REMOVE,
			Path: "/modules/111/tree",
		},
	}
	return patch

}

func TestTopicPatchCompiler(t *testing.T) {
	patch := createTopicPatch()
	compiler := ForTopics()
	list, err := compiler.Compile("123", patch)

	assert.Nil(t, err)

	assert.NotNil(t, list)

	assert.Equal(t, "SET TRANSACTION ISOLATION LEVEL SERIALIZABLE", list.Commands[0].statement)
	assert.Equal(t, "SELECT check_topic_version($1,$2)", list.Commands[1].statement)

	assert.Equal(t, "SELECT update_topic_description($1,$2)", list.Commands[2].statement)
	assert.Equal(t, "Hallo Welt", list.Commands[2].parameters[1])

}
