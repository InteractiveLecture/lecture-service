package lecturepatch

import (
	"testing"

	"github.com/richterrettich/jsonpatch"
	"github.com/stretchr/testify/assert"
)

func createExercisePatch() *jsonpatch.Patch {
	patch := &jsonpatch.Patch{}
	patch.Version = 1
	patch.Operations = []jsonpatch.Operation{
		jsonpatch.Operation{
			Type: jsonpatch.ADD,
			Path: "/tasks/1/hints",
			Value: map[string]interface{}{
				"id":       "999",
				"position": 1,
				"cost":     100,
				"content":  "Dies ist der erste Hint",
			},
		},
		jsonpatch.Operation{
			Type: jsonpatch.REMOVE,
			Path: "/tasks/1/hints/1",
		},
		jsonpatch.Operation{
			Type: jsonpatch.MOVE,
			Path: "/tasks/1/hints/2",
			From: "/tasks/1/hints/1",
		},
		jsonpatch.Operation{
			Type:  jsonpatch.REPLACE,
			Path:  "/tasks/1/hints/1/content",
			Value: "Dies ist der neue erste Hint",
		},
		jsonpatch.Operation{
			Type:  jsonpatch.REPLACE,
			Path:  "/tasks/1/hints/1/cost",
			Value: 200,
		},
		jsonpatch.Operation{
			Type: jsonpatch.ADD,
			Path: "/tasks",
			Value: map[string]interface{}{
				"id":       "888",
				"position": 1,
				"content":  "Dies ist der erste Task",
			},
		},
		jsonpatch.Operation{
			Type: jsonpatch.REMOVE,
			Path: "/tasks/1",
		},

		jsonpatch.Operation{
			Type: jsonpatch.MOVE,
			Path: "/tasks/2",
			From: "/tasks/1",
		},
		jsonpatch.Operation{
			Type:  jsonpatch.REPLACE,
			Path:  "/tasks/1/content",
			Value: "Dies ist immer noch der erste Task",
		},
	}
	return patch
}

func createModulePatch() *jsonpatch.Patch {

	patch := &jsonpatch.Patch{}

	patch.Version = 1
	patch.Operations = []jsonpatch.Operation{
		jsonpatch.Operation{
			Type:  jsonpatch.REPLACE,
			Path:  "/description",
			Value: "Hallo Welt",
		},
		jsonpatch.Operation{
			Type:  jsonpatch.ADD,
			Path:  "/recommendations",
			Value: "111",
		},
		jsonpatch.Operation{
			Type: jsonpatch.REMOVE,
			Path: "/recommendations/111",
		},
		jsonpatch.Operation{
			Type: jsonpatch.ADD,
			Path: "/exercises",
			Value: map[string]interface{}{
				"id":      "333",
				"backend": "Java",
				"tasks":   []string{"urf urf", "bla bla"},
			},
		},
		jsonpatch.Operation{
			Type: jsonpatch.REMOVE,
			Path: "/exercises/333",
		},
		jsonpatch.Operation{
			Type:  jsonpatch.ADD,
			Path:  "/video",
			Value: "555",
		},
		jsonpatch.Operation{
			Type: jsonpatch.REMOVE,
			Path: "/video/555",
		},

		jsonpatch.Operation{
			Type:  jsonpatch.ADD,
			Path:  "/script",
			Value: "555",
		},
		jsonpatch.Operation{
			Type: jsonpatch.REMOVE,
			Path: "/script/555",
		},
	}
	return patch

}

func createTopicPatch() *jsonpatch.Patch {

	patch := &jsonpatch.Patch{}

	patch.Version = 1
	patch.Operations = []jsonpatch.Operation{
		jsonpatch.Operation{
			Type:  jsonpatch.REPLACE,
			Path:  "/description",
			Value: "Hallo Welt",
		},
		jsonpatch.Operation{
			Type:  jsonpatch.ADD,
			Path:  "/assistants",
			Value: "111",
		},
		jsonpatch.Operation{
			Type: jsonpatch.REMOVE,
			Path: "/assistants/111",
		},
		jsonpatch.Operation{
			Type: jsonpatch.ADD,
			Path: "/modules",
			Value: map[string]interface{}{
				"id":          "222",
				"description": "Hugo",
				"video_id":    "333",
				"script_id":   "444",
				"parents":     []string{"555", "666"},
			},
		},

		jsonpatch.Operation{
			Type:  jsonpatch.REPLACE,
			Path:  "/modules/222/parents",
			Value: []string{"555", "666"},
		},
		jsonpatch.Operation{
			Type:  jsonpatch.REPLACE,
			Path:  "/modules/222/parents/tree",
			Value: []string{"555", "666"},
		},

		jsonpatch.Operation{
			Type: jsonpatch.REMOVE,
			Path: "/modules/222",
		},

		jsonpatch.Operation{
			Type: jsonpatch.REMOVE,
			Path: "/modules/222/tree",
		},
	}
	return patch

}

func TestModulePatchCompiler(t *testing.T) {
	patch := createModulePatch()
	compiler := ForModules()
	options := map[string]interface{}{
		"id":     "123",
		"userId": "444",
	}
	list, err := compiler.Compile(patch, options)

	assert.Nil(t, err)

	assert.NotNil(t, list)

	assert.Equal(t, "SET TRANSACTION ISOLATION LEVEL SERIALIZABLE", list.Commands[0].(*SqlCommandContainer).statement)
	assert.Equal(t, "SELECT check_version($1,$2,$3)", list.Commands[1].(*SqlCommandContainer).statement)

	assert.Equal(t, "SELECT increment_version($1,$2)", list.Commands[len(list.Commands)-2].(*SqlCommandContainer).statement)
	assert.Equal(t, "REFRESH MATERIALIZED VIEW module_trees", list.Commands[len(list.Commands)-1].(*SqlCommandContainer).statement)

	assert.Equal(t, "SELECT replace_module_description($1,$2)", list.Commands[2].(*SqlCommandContainer).statement)
	assert.Equal(t, "123", list.Commands[2].(*SqlCommandContainer).parameters[0])
	assert.Equal(t, "Hallo Welt", list.Commands[2].(*SqlCommandContainer).parameters[1])

	assert.Equal(t, "SELECT add_module_recommendation($1,$2)", list.Commands[3].(*SqlCommandContainer).statement)
	assert.Equal(t, "123", list.Commands[3].(*SqlCommandContainer).parameters[0])
	assert.Equal(t, "111", list.Commands[3].(*SqlCommandContainer).parameters[1])

	assert.Equal(t, "SELECT remove_module_recommendation($1,$2)", list.Commands[4].(*SqlCommandContainer).statement)
	assert.Equal(t, "123", list.Commands[4].(*SqlCommandContainer).parameters[0])
	assert.Equal(t, "111", list.Commands[4].(*SqlCommandContainer).parameters[1])

	assert.Equal(t, "SELECT add_exercise($1,$2,$3)", list.Commands[5].(*SqlCommandContainer).statement)
	assert.Equal(t, "333", list.Commands[5].(*SqlCommandContainer).parameters[0])
	assert.Equal(t, "123", list.Commands[5].(*SqlCommandContainer).parameters[1])
	assert.Equal(t, "Java", list.Commands[5].(*SqlCommandContainer).parameters[2])

	assert.Equal(t, "SELECT remove_exercise($1,$2)", list.Commands[6].(*SqlCommandContainer).statement)
	assert.Equal(t, "123", list.Commands[6].(*SqlCommandContainer).parameters[0])
	assert.Equal(t, "333", list.Commands[6].(*SqlCommandContainer).parameters[1])

	assert.Equal(t, "SELECT add_module_video($1,$2)", list.Commands[7].(*SqlCommandContainer).statement)
	assert.Equal(t, "123", list.Commands[7].(*SqlCommandContainer).parameters[0])
	assert.Equal(t, "555", list.Commands[7].(*SqlCommandContainer).parameters[1])

	assert.Equal(t, "SELECT remove_module_video($1,$2)", list.Commands[8].(*SqlCommandContainer).statement)
	assert.Equal(t, "123", list.Commands[8].(*SqlCommandContainer).parameters[0])
	assert.Equal(t, "555", list.Commands[8].(*SqlCommandContainer).parameters[1])

	assert.Equal(t, "SELECT add_module_script($1,$2)", list.Commands[9].(*SqlCommandContainer).statement)
	assert.Equal(t, "123", list.Commands[9].(*SqlCommandContainer).parameters[0])
	assert.Equal(t, "555", list.Commands[9].(*SqlCommandContainer).parameters[1])

	assert.Equal(t, "SELECT remove_module_script($1,$2)", list.Commands[10].(*SqlCommandContainer).statement)
	assert.Equal(t, "123", list.Commands[10].(*SqlCommandContainer).parameters[0])
	assert.Equal(t, "555", list.Commands[10].(*SqlCommandContainer).parameters[1])
}

func TestTopicPatchCompiler(t *testing.T) {
	patch := createTopicPatch()
	compiler := ForTopics()
	options := map[string]interface{}{
		"id":     "123",
		"userId": "444",
	}
	list, err := compiler.Compile(patch, options)

	assert.Nil(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, "SET TRANSACTION ISOLATION LEVEL SERIALIZABLE", list.Commands[0].(*SqlCommandContainer).statement)
	assert.Equal(t, "SELECT check_version($1,$2,$3)", list.Commands[1].(*SqlCommandContainer).statement)
	assert.Equal(t, "SELECT replace_topic_description($1,$2)", list.Commands[2].(*SqlCommandContainer).statement)
	assert.Equal(t, "Hallo Welt", list.Commands[2].(*SqlCommandContainer).parameters[1])
	assert.Equal(t, "SELECT add_assistant($1,$2)", list.Commands[3].(*SqlCommandContainer).statement)
	assert.Equal(t, "123", list.Commands[3].(*SqlCommandContainer).parameters[0])
	assert.Equal(t, "111", list.Commands[3].(*SqlCommandContainer).parameters[1])

	assert.Equal(t, "SELECT remove_assistant($1,$2)", list.Commands[4].(*SqlCommandContainer).statement)
	assert.Equal(t, "123", list.Commands[4].(*SqlCommandContainer).parameters[0])
	assert.Equal(t, "111", list.Commands[4].(*SqlCommandContainer).parameters[1])

	assert.Equal(t, "SELECT add_module($1,$2,$3,$4,$5,$6,$7)", list.Commands[5].(*SqlCommandContainer).statement)
	assert.Equal(t, "222", list.Commands[5].(*SqlCommandContainer).parameters[0])
	assert.Equal(t, "123", list.Commands[5].(*SqlCommandContainer).parameters[1])
	assert.Equal(t, "Hugo", list.Commands[5].(*SqlCommandContainer).parameters[2])
	assert.Equal(t, "333", list.Commands[5].(*SqlCommandContainer).parameters[3])
	assert.Equal(t, "444", list.Commands[5].(*SqlCommandContainer).parameters[4])
	assert.Equal(t, "555", list.Commands[5].(*SqlCommandContainer).parameters[5])
	assert.Equal(t, "666", list.Commands[5].(*SqlCommandContainer).parameters[6])

	assert.Equal(t, "SELECT move_module($1,$2,$3,$4)", list.Commands[6].(*SqlCommandContainer).statement)
	assert.Equal(t, "123", list.Commands[6].(*SqlCommandContainer).parameters[0])
	assert.Equal(t, "222", list.Commands[6].(*SqlCommandContainer).parameters[1])
	assert.Equal(t, "555", list.Commands[6].(*SqlCommandContainer).parameters[2])
	assert.Equal(t, "666", list.Commands[6].(*SqlCommandContainer).parameters[3])

	assert.Equal(t, "SELECT move_module_tree($1,$2,$3,$4)", list.Commands[7].(*SqlCommandContainer).statement)
	assert.Equal(t, "123", list.Commands[7].(*SqlCommandContainer).parameters[0])
	assert.Equal(t, "222", list.Commands[7].(*SqlCommandContainer).parameters[1])
	assert.Equal(t, "555", list.Commands[7].(*SqlCommandContainer).parameters[2])
	assert.Equal(t, "666", list.Commands[7].(*SqlCommandContainer).parameters[3])

	assert.Equal(t, "SELECT remove_module($1,$2)", list.Commands[8].(*SqlCommandContainer).statement)
	assert.Equal(t, "123", list.Commands[8].(*SqlCommandContainer).parameters[0])
	assert.Equal(t, "222", list.Commands[8].(*SqlCommandContainer).parameters[1])

	assert.Equal(t, "SELECT remove_module_tree($1,$2)", list.Commands[9].(*SqlCommandContainer).statement)
	assert.Equal(t, "123", list.Commands[9].(*SqlCommandContainer).parameters[0])
	assert.Equal(t, "222", list.Commands[9].(*SqlCommandContainer).parameters[1])

	assert.Equal(t, "SELECT increment_version($1,$2)", list.Commands[len(list.Commands)-1].(*SqlCommandContainer).statement)
}

func TestExercicePatchCompiler(t *testing.T) {
	patch := createExercisePatch()
	compiler := ForExercises()
	options := map[string]interface{}{
		"id":     "123",
		"userId": "444",
	}
	list, err := compiler.Compile(patch, options)

	assert.Nil(t, err)
	assert.Equal(t, "SET TRANSACTION ISOLATION LEVEL SERIALIZABLE", list.Commands[0].(*SqlCommandContainer).statement)
	assert.Equal(t, "SELECT check_version($1,$2,$3)", list.Commands[1].(*SqlCommandContainer).statement)
	assert.Equal(t, "SELECT increment_version($1,$2)", list.Commands[len(list.Commands)-1].(*SqlCommandContainer).statement)

	assert.Equal(t, "SELECT add_hint($1,$2,$3,$4,$5,$6)", list.Commands[2].(*SqlCommandContainer).statement)
	assert.Equal(t, "123", list.Commands[2].(*SqlCommandContainer).parameters[0])
	assert.Equal(t, 1, list.Commands[2].(*SqlCommandContainer).parameters[1])
	assert.Equal(t, "999", list.Commands[2].(*SqlCommandContainer).parameters[2])
	assert.Equal(t, 1, list.Commands[2].(*SqlCommandContainer).parameters[3])
	assert.Equal(t, "Dies ist der erste Hint", list.Commands[2].(*SqlCommandContainer).parameters[4])
	assert.Equal(t, 100, list.Commands[2].(*SqlCommandContainer).parameters[5])

	assert.Equal(t, "SELECT remove_hint($1,$2,$3)", list.Commands[3].(*SqlCommandContainer).statement)
	assert.Equal(t, "123", list.Commands[3].(*SqlCommandContainer).parameters[0])
	assert.Equal(t, 1, list.Commands[3].(*SqlCommandContainer).parameters[1])
	assert.Equal(t, 1, list.Commands[3].(*SqlCommandContainer).parameters[2])

	assert.Equal(t, "SELECT move_hint($1,$2,$3,$4,$5)", list.Commands[4].(*SqlCommandContainer).statement)
	assert.Equal(t, "123", list.Commands[4].(*SqlCommandContainer).parameters[0])
	assert.Equal(t, 1, list.Commands[4].(*SqlCommandContainer).parameters[1])
	assert.Equal(t, 1, list.Commands[4].(*SqlCommandContainer).parameters[2])
	assert.Equal(t, 1, list.Commands[4].(*SqlCommandContainer).parameters[3])
	assert.Equal(t, 2, list.Commands[4].(*SqlCommandContainer).parameters[4])

	assert.Equal(t, "SELECT replace_hint_content($1,$2,$3,$4)", list.Commands[5].(*SqlCommandContainer).statement)
	assert.Equal(t, "123", list.Commands[5].(*SqlCommandContainer).parameters[0])
	assert.Equal(t, 1, list.Commands[5].(*SqlCommandContainer).parameters[1])
	assert.Equal(t, 1, list.Commands[5].(*SqlCommandContainer).parameters[2])
	assert.Equal(t, "Dies ist der neue erste Hint", list.Commands[5].(*SqlCommandContainer).parameters[3])

	assert.Equal(t, "SELECT replace_hint_cost($1,$2,$3,$4)", list.Commands[6].(*SqlCommandContainer).statement)
	assert.Equal(t, "123", list.Commands[6].(*SqlCommandContainer).parameters[0])
	assert.Equal(t, 1, list.Commands[6].(*SqlCommandContainer).parameters[1])
	assert.Equal(t, 1, list.Commands[6].(*SqlCommandContainer).parameters[2])
	assert.Equal(t, 200, list.Commands[6].(*SqlCommandContainer).parameters[3])

	assert.Equal(t, "SELECT add_task($1,$2,$3,$4)", list.Commands[7].(*SqlCommandContainer).statement)
	assert.Equal(t, "123", list.Commands[7].(*SqlCommandContainer).parameters[0])
	assert.Equal(t, "888", list.Commands[7].(*SqlCommandContainer).parameters[1])
	assert.Equal(t, 1, list.Commands[7].(*SqlCommandContainer).parameters[2])
	assert.Equal(t, "Dies ist der erste Task", list.Commands[7].(*SqlCommandContainer).parameters[3])

	assert.Equal(t, "SELECT remove_task($1,$2)", list.Commands[8].(*SqlCommandContainer).statement)
	assert.Equal(t, "123", list.Commands[8].(*SqlCommandContainer).parameters[0])
	assert.Equal(t, 1, list.Commands[8].(*SqlCommandContainer).parameters[1])

	assert.Equal(t, "SELECT move_task($1,$2,$3)", list.Commands[9].(*SqlCommandContainer).statement)
	assert.Equal(t, "123", list.Commands[9].(*SqlCommandContainer).parameters[0])
	assert.Equal(t, 1, list.Commands[9].(*SqlCommandContainer).parameters[1])
	assert.Equal(t, 2, list.Commands[9].(*SqlCommandContainer).parameters[2])

	assert.Equal(t, "SELECT replace_task_content($1,$2,$3)", list.Commands[10].(*SqlCommandContainer).statement)
	assert.Equal(t, "123", list.Commands[10].(*SqlCommandContainer).parameters[0])
	assert.Equal(t, 1, list.Commands[10].(*SqlCommandContainer).parameters[1])
	assert.Equal(t, "Dies ist immer noch der erste Task", list.Commands[10].(*SqlCommandContainer).parameters[2])

}
