package patchcompiler

import (
	"testing"

	"github.com/richterrettich/lecture-service/lecturepatch"
	"github.com/stretchr/testify/assert"
)

func createModulePatch() *lecturepatch.Patch {

	patch := &lecturepatch.Patch{}

	patch.Version = 1
	patch.Operations = []lecturepatch.Operation{
		lecturepatch.Operation{
			Type:  lecturepatch.REPLACE,
			Path:  "/description",
			Value: "Hallo Welt",
		},
		lecturepatch.Operation{
			Type:  lecturepatch.ADD,
			Path:  "/recommendations",
			Value: "111",
		},
		lecturepatch.Operation{
			Type: lecturepatch.REMOVE,
			Path: "/recommendations/111",
		},
		lecturepatch.Operation{
			Type: lecturepatch.ADD,
			Path: "/exercises",
			Value: map[string]interface{}{
				"id":      "333",
				"backend": "Java",
			},
		},
		lecturepatch.Operation{
			Type: lecturepatch.REMOVE,
			Path: "/exercises/333",
		},
		lecturepatch.Operation{
			Type:  lecturepatch.ADD,
			Path:  "/video",
			Value: "555",
		},
		lecturepatch.Operation{
			Type: lecturepatch.REMOVE,
			Path: "/video/567",
		},

		lecturepatch.Operation{
			Type:  lecturepatch.ADD,
			Path:  "/script",
			Value: "555",
		},
		lecturepatch.Operation{
			Type: lecturepatch.REMOVE,
			Path: "/script/667",
		},
	}
	return patch

}

func createTopicPatch() *lecturepatch.Patch {

	patch := &lecturepatch.Patch{}

	patch.Version = 1
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

func TestModulePatchCompiler(t *testing.T) {
	patch := createModulePatch()
	compiler := ForModules()
	list, err := compiler.Compile("123", patch)

	assert.Nil(t, err)

	assert.NotNil(t, list)

	assert.Equal(t, "SET TRANSACTION ISOLATION LEVEL SERIALIZABLE", list.Commands[0].statement)
	assert.Equal(t, "SELECT check_module_version($1,$2)", list.Commands[1].statement)

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
