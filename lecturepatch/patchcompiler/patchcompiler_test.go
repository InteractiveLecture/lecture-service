package patchcompiler

import (
	"testing"

	"github.com/richterrettich/lecture-service/lecturepatch"
	"github.com/stretchr/testify/assert"
)

func createExercisePatch() *lecturepatch.Patch {
	patch := &lecturepatch.Patch{}
	patch.Version = 1
	patch.Operations = []lecturepatch.Operation{
		lecturepatch.Operation{
			Type: lecturepatch.ADD,
			Path: "/hints",
			Value: map[string]interface{}{
				"id":       "999",
				"position": 1,
				"cost":     100,
				"content":  "Dies ist der erste Hint",
			},
		},
		lecturepatch.Operation{
			Type: lecturepatch.REMOVE,
			Path: "/hints/1",
		},
		lecturepatch.Operation{
			Type: lecturepatch.MOVE,
			Path: "/hints/2",
			From: "/hints/1",
		},
		lecturepatch.Operation{
			Type:  lecturepatch.REPLACE,
			Path:  "/hints/1/content",
			Value: "Dies ist der neue erste Hint",
		},
		lecturepatch.Operation{
			Type:  lecturepatch.REPLACE,
			Path:  "/hints/1/cost",
			Value: 200,
		},
		lecturepatch.Operation{
			Type: lecturepatch.ADD,
			Path: "/tasks",
			Value: map[string]interface{}{
				"position": 1,
				"content":  "Dies ist der erste Task",
			},
		},
		lecturepatch.Operation{
			Type: lecturepatch.REMOVE,
			Path: "/tasks/1",
		},

		lecturepatch.Operation{
			Type: lecturepatch.MOVE,
			Path: "/tasks/2",
			From: "/tasks/1",
		},
		lecturepatch.Operation{
			Type:  lecturepatch.REPLACE,
			Path:  "/tasks/1/content",
			Value: "Dies ist immer noch der erste Task",
		},
	}
	return patch
}

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
			Path: "/video/555",
		},

		lecturepatch.Operation{
			Type:  lecturepatch.ADD,
			Path:  "/script",
			Value: "555",
		},
		lecturepatch.Operation{
			Type: lecturepatch.REMOVE,
			Path: "/script/555",
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
			Path:  "/modules/222/parents",
			Value: []string{"555", "666"},
		},
		lecturepatch.Operation{
			Type:  lecturepatch.REPLACE,
			Path:  "/modules/222/parents/tree",
			Value: []string{"555", "666"},
		},

		lecturepatch.Operation{
			Type: lecturepatch.REMOVE,
			Path: "/modules/222",
		},

		lecturepatch.Operation{
			Type: lecturepatch.REMOVE,
			Path: "/modules/222/tree",
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
	assert.Equal(t, "SELECT check_version($1,$2,$3)", list.Commands[1].statement)

	assert.Equal(t, "SELECT increment_version($1,$2)", list.Commands[len(list.Commands)-2].statement)
	assert.Equal(t, "REFRESH MATERIALIZED VIEW module_trees", list.Commands[len(list.Commands)-1].statement)

	assert.Equal(t, "SELECT replace_module_description($1,$2)", list.Commands[2].statement)
	assert.Equal(t, "123", list.Commands[2].parameters[0])
	assert.Equal(t, "Hallo Welt", list.Commands[2].parameters[1])

	assert.Equal(t, "SELECT add_module_recommendation($1,$2)", list.Commands[3].statement)
	assert.Equal(t, "123", list.Commands[3].parameters[0])
	assert.Equal(t, "111", list.Commands[3].parameters[1])

	assert.Equal(t, "SELECT remove_module_recommendation($1,$2)", list.Commands[4].statement)
	assert.Equal(t, "123", list.Commands[4].parameters[0])
	assert.Equal(t, "111", list.Commands[4].parameters[1])

	assert.Equal(t, "SELECT add_exercise($1,$2,$3)", list.Commands[5].statement)
	assert.Equal(t, "333", list.Commands[5].parameters[0])
	assert.Equal(t, "123", list.Commands[5].parameters[1])
	assert.Equal(t, "Java", list.Commands[5].parameters[2])

	assert.Equal(t, "SELECT remove_exercise($1,$2)", list.Commands[6].statement)
	assert.Equal(t, "123", list.Commands[6].parameters[0])
	assert.Equal(t, "333", list.Commands[6].parameters[1])

	assert.Equal(t, "SELECT add_module_video($1,$2)", list.Commands[7].statement)
	assert.Equal(t, "123", list.Commands[7].parameters[0])
	assert.Equal(t, "555", list.Commands[7].parameters[1])

	assert.Equal(t, "SELECT remove_module_video($1,$2)", list.Commands[8].statement)
	assert.Equal(t, "123", list.Commands[8].parameters[0])
	assert.Equal(t, "555", list.Commands[8].parameters[1])

	assert.Equal(t, "SELECT add_module_script($1,$2)", list.Commands[9].statement)
	assert.Equal(t, "123", list.Commands[9].parameters[0])
	assert.Equal(t, "555", list.Commands[9].parameters[1])

	assert.Equal(t, "SELECT remove_module_script($1,$2)", list.Commands[10].statement)
	assert.Equal(t, "123", list.Commands[10].parameters[0])
	assert.Equal(t, "555", list.Commands[10].parameters[1])

}

func TestTopicPatchCompiler(t *testing.T) {
	patch := createTopicPatch()
	compiler := ForTopics()
	list, err := compiler.Compile("123", patch)
	assert.Nil(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, "SET TRANSACTION ISOLATION LEVEL SERIALIZABLE", list.Commands[0].statement)
	assert.Equal(t, "SELECT check_version($1,$2,$3)", list.Commands[1].statement)
	assert.Equal(t, "SELECT replace_topic_description($1,$2)", list.Commands[2].statement)
	assert.Equal(t, "Hallo Welt", list.Commands[2].parameters[1])
	assert.Equal(t, "SELECT add_assistant($1,$2)", list.Commands[3].statement)
	assert.Equal(t, "123", list.Commands[3].parameters[0])
	assert.Equal(t, "111", list.Commands[3].parameters[1])

	assert.Equal(t, "SELECT remove_assistant($1,$2)", list.Commands[4].statement)
	assert.Equal(t, "123", list.Commands[4].parameters[0])
	assert.Equal(t, "111", list.Commands[4].parameters[1])

	assert.Equal(t, "SELECT add_module($1,$2,$3,$4,$5,$6,$7)", list.Commands[5].statement)
	assert.Equal(t, "222", list.Commands[5].parameters[0])
	assert.Equal(t, "123", list.Commands[5].parameters[1])
	assert.Equal(t, "Hugo", list.Commands[5].parameters[2])
	assert.Equal(t, "333", list.Commands[5].parameters[3])
	assert.Equal(t, "444", list.Commands[5].parameters[4])
	assert.Equal(t, "555", list.Commands[5].parameters[5])
	assert.Equal(t, "666", list.Commands[5].parameters[6])

	assert.Equal(t, "SELECT move_module($1,$2,$3,$4)", list.Commands[6].statement)
	assert.Equal(t, "123", list.Commands[6].parameters[0])
	assert.Equal(t, "222", list.Commands[6].parameters[1])
	assert.Equal(t, "555", list.Commands[6].parameters[2])
	assert.Equal(t, "666", list.Commands[6].parameters[3])

	assert.Equal(t, "SELECT move_module_tree($1,$2,$3,$4)", list.Commands[7].statement)
	assert.Equal(t, "123", list.Commands[7].parameters[0])
	assert.Equal(t, "222", list.Commands[7].parameters[1])
	assert.Equal(t, "555", list.Commands[7].parameters[2])
	assert.Equal(t, "666", list.Commands[7].parameters[3])

	assert.Equal(t, "SELECT remove_module($1,$2)", list.Commands[8].statement)
	assert.Equal(t, "123", list.Commands[8].parameters[0])
	assert.Equal(t, "222", list.Commands[8].parameters[1])

	assert.Equal(t, "SELECT remove_module_tree($1,$2)", list.Commands[9].statement)
	assert.Equal(t, "123", list.Commands[9].parameters[0])
	assert.Equal(t, "222", list.Commands[9].parameters[1])

	assert.Equal(t, "SELECT increment_version($1,$2)", list.Commands[len(list.Commands)-2].statement)
	assert.Equal(t, "REFRESH MATERIALIZED VIEW module_trees", list.Commands[len(list.Commands)-1].statement)
}

func TestExercicePatchCompiler(t *testing.T) {
	patch := createExercisePatch()
	compiler := ForExercises()
	list, err := compiler.Compile("123", patch)
	assert.Nil(t, err)
	assert.Equal(t, "SET TRANSACTION ISOLATION LEVEL SERIALIZABLE", list.Commands[0].statement)
	assert.Equal(t, "SELECT check_version($1,$2,$3)", list.Commands[1].statement)
	assert.Equal(t, "SELECT increment_version($1,$2)", list.Commands[len(list.Commands)-1].statement)

	assert.Equal(t, "SELECT add_hint($1,$2,$3,$4,$5)", list.Commands[2].statement)
	assert.Equal(t, "999", list.Commands[2].parameters[0])
	assert.Equal(t, "123", list.Commands[2].parameters[1])
	assert.Equal(t, 1, list.Commands[2].parameters[2])
	assert.Equal(t, "Dies ist der erste Hint", list.Commands[2].parameters[3])
	assert.Equal(t, 100, list.Commands[2].parameters[4])

	assert.Equal(t, "SELECT remove_hint($1,$2)", list.Commands[3].statement)
	assert.Equal(t, "123", list.Commands[3].parameters[0])
	assert.Equal(t, 1, list.Commands[3].parameters[1])

	assert.Equal(t, "SELECT move_hint($1,$2,$3)", list.Commands[4].statement)
	assert.Equal(t, "123", list.Commands[4].parameters[0])
	assert.Equal(t, 1, list.Commands[4].parameters[1])
	assert.Equal(t, 2, list.Commands[4].parameters[2])

	assert.Equal(t, "SELECT replace_hint_content($1,$2,$3)", list.Commands[5].statement)
	assert.Equal(t, "123", list.Commands[5].parameters[0])
	assert.Equal(t, 1, list.Commands[5].parameters[1])
	assert.Equal(t, "Dies ist der neue erste Hint", list.Commands[5].parameters[2])

	assert.Equal(t, "SELECT replace_hint_cost($1,$2,$3)", list.Commands[6].statement)
	assert.Equal(t, "123", list.Commands[6].parameters[0])
	assert.Equal(t, 1, list.Commands[6].parameters[1])
	assert.Equal(t, 200, list.Commands[6].parameters[2])

	assert.Equal(t, "SELECT add_task($1,$2,$3)", list.Commands[7].statement)
	assert.Equal(t, "123", list.Commands[7].parameters[0])
	assert.Equal(t, 1, list.Commands[7].parameters[1])
	assert.Equal(t, "Dies ist der erste Task", list.Commands[7].parameters[2])

	assert.Equal(t, "SELECT remove_task($1,$2)", list.Commands[8].statement)
	assert.Equal(t, "123", list.Commands[8].parameters[0])
	assert.Equal(t, 1, list.Commands[8].parameters[1])

	assert.Equal(t, "SELECT move_task($1,$2,$3)", list.Commands[9].statement)
	assert.Equal(t, "123", list.Commands[9].parameters[0])
	assert.Equal(t, 1, list.Commands[9].parameters[1])
	assert.Equal(t, 2, list.Commands[9].parameters[2])

	assert.Equal(t, "SELECT replace_task_content($1,$2,$3)", list.Commands[10].statement)
	assert.Equal(t, "123", list.Commands[10].parameters[1])
	assert.Equal(t, 1, list.Commands[10].parameters[2])
	assert.Equal(t, "Dies ist immer noch der erste Task", list.Commands[10].parameters[0])

}
