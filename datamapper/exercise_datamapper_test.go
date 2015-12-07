package datamapper

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/richterrettich/jsonpatch"
	"github.com/richterrettich/lecture-service/lecturepatch"
	"github.com/richterrettich/lecture-service/paginator"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestCompleteExercise(t *testing.T) {
	mapper, err := prepareMapper()
	assert.Nil(t, err)
	assert.Nil(t, resetDatabase(mapper))
	err = mapper.CompleteExercise("f7c21557-03fc-3e99-bdff-7b065f58b39d", "233804c6-55b8-3807-9733-9c090d75decf")
	assert.Nil(t, err)
	pr := paginator.PageRequest{-1, -1, nil}
	data, err := mapper.GetExerciseHistory("233804c6-55b8-3807-9733-9c090d75decf", pr, "")
	assert.Nil(t, err)
	histories := make([]interface{}, 0)
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&histories)
	assert.Nil(t, err)
	historySet := toSet(histories, "exercise_id")
	for _, v := range []string{"f7c21557-03fc-3e99-bdff-7b065f58b39d"} {
		assert.True(t, historySet[v])
	}
}

func TestGetHint(t *testing.T) {
	mapper, err := prepareMapper()
	assert.Nil(t, err)
	assert.Nil(t, resetDatabase(mapper))
	_, err = mapper.GetHint("3be2043d-9e70-3212-8fcd-42a7ae38c8a2", "f20919fa-08bd-3a8d-9e3c-e3406c680162")
	assert.Nil(t, err)
	_, err = mapper.GetHint("164186cb-1252-3672-a015-e8128b999bb4", "f20919fa-08bd-3a8d-9e3c-e3406c680162")
	_, ok := err.(PaymentRequiredError)
	assert.True(t, ok)
}

func TestPurchaseHint(t *testing.T) {
	mapper, err := prepareMapper()
	assert.Nil(t, err)
	assert.Nil(t, resetDatabase(mapper))
	err = mapper.PurchaseHint("164186cb-1252-3672-a015-e8128b999bb4", "f20919fa-08bd-3a8d-9e3c-e3406c680162")
	assert.Nil(t, err)

	err = mapper.PurchaseHint("164186cb-1252-3672-a015-e8128b999bb4", "f20919fa-08bd-3a8d-9e3c-e3406c680162")
	_, ok := err.(AlreadyPurchasedError)
	assert.True(t, ok)
	err = mapper.PurchaseHint("1bcfd6c3-b269-392a-a7c0-2206f9aefcb6", "f20919fa-08bd-3a8d-9e3c-e3406c680162")
	_, ok = err.(InsufficientPointsError)
	assert.True(t, ok)
	err = mapper.PurchaseHint("f20919fa-08bd-3a8d-9e3c-e3406c680162", "1bcfd6c3-b269-392a-a7c0-2206f9aefcb6")
	_, ok = err.(HintNotFoundError)
	assert.True(t, ok)
}

func TestAddRemoveAlterHint(t *testing.T) {
	mapper, err := prepareMapper()
	assert.Nil(t, err)
	assert.Nil(t, resetDatabase(mapper))
	hintId := uuid.NewV4().String()
	p := jsonpatch.Patch{
		Version: 1,
		Operations: []jsonpatch.Operation{
			jsonpatch.Operation{
				Type: jsonpatch.ADD,
				Path: "/hints",
				Value: map[string]interface{}{
					"id":       hintId,
					"position": 1,
					"content":  "ein neuer hint",
					"cost":     100,
				},
			},
		},
	}
	compiler := lecturepatch.ForExercises()
	err = mapper.ApplyPatch("f7c21557-03fc-3e99-bdff-7b065f58b39d", &p, compiler)
	assert.Nil(t, err)
	ex, err := getExercise(mapper, "f7c21557-03fc-3e99-bdff-7b065f58b39d")
	assert.Nil(t, err)
	hints := ex["hint_ids"].([]interface{})
	assert.Equal(t, hintId, hints[0].(string))
	hint, err := getHint(mapper, hintId)
	assert.Nil(t, err)
	assert.Equal(t, int(hint["cost"].(float64)), int(100))
	assert.Equal(t, hint["content"], "ein neuer hint")
	p = jsonpatch.Patch{
		Version: 2,
		Operations: []jsonpatch.Operation{
			jsonpatch.Operation{
				Type: jsonpatch.MOVE,
				From: "/hints/1",
				Path: "/hints/2",
			},
			jsonpatch.Operation{
				Type: jsonpatch.REMOVE,
				Path: "/hints/3",
			},
			jsonpatch.Operation{
				Type:  jsonpatch.REPLACE,
				Path:  "/hints/3/content",
				Value: "das ist der neue dritte hint",
			},
			jsonpatch.Operation{
				Type:  jsonpatch.REPLACE,
				Path:  "/hints/3/cost",
				Value: 200,
			},
		},
	}
	err = mapper.ApplyPatch("f7c21557-03fc-3e99-bdff-7b065f58b39d", &p, compiler)
	assert.Nil(t, err)
	ex, err = getExercise(mapper, "f7c21557-03fc-3e99-bdff-7b065f58b39d")
	hints = ex["hint_ids"].([]interface{})
	assert.Equal(t, hintId, hints[1].(string))
	assert.Equal(t, 3, len(hints))
	set := arrayToSet(hints)
	assert.False(t, set["164186cb-1252-3672-a015-e8128b999bb4"])
	assert.Equal(t, hints[2].(string), "1bcfd6c3-b269-392a-a7c0-2206f9aefcb6")
	hint, err = getHint(mapper, "1bcfd6c3-b269-392a-a7c0-2206f9aefcb6")
	assert.Nil(t, err)
	assert.Equal(t, "das ist der neue dritte hint", hint["content"])
	assert.Equal(t, 200, int(hint["cost"].(float64)))
}

func TestAddRemoveAlterTask(t *testing.T) {
	mapper, err := prepareMapper()
	assert.Nil(t, err)
	assert.Nil(t, resetDatabase(mapper))
	taskId := uuid.NewV4().String()
	p := jsonpatch.Patch{
		Version: 1,
		Operations: []jsonpatch.Operation{
			jsonpatch.Operation{
				Type: jsonpatch.ADD,
				Path: "/tasks",
				Value: map[string]interface{}{
					"position": 1,
					"content":  "ein neuer erster task",
				},
			},
		},
	}
	compiler := lecturepatch.ForExercises()
	err = mapper.ApplyPatch("f7c21557-03fc-3e99-bdff-7b065f58b39d", &p, compiler)
	assert.Nil(t, err)
	ex, err := getExercise(mapper, "f7c21557-03fc-3e99-bdff-7b065f58b39d")
	assert.Nil(t, err)
	hints := ex["tasks"].([]interface{})
	assert.Equal(t, hintId, hints[0].(string))
	hint, err := getHint(mapper, hintId)
	assert.Nil(t, err)
	assert.Equal(t, int(hint["cost"].(float64)), int(100))
	assert.Equal(t, hint["content"], "ein neuer hint")
	p = jsonpatch.Patch{
		Version: 2,
		Operations: []jsonpatch.Operation{
			jsonpatch.Operation{
				Type: jsonpatch.MOVE,
				From: "/hints/1",
				Path: "/hints/2",
			},
			jsonpatch.Operation{
				Type: jsonpatch.REMOVE,
				Path: "/hints/3",
			},
			jsonpatch.Operation{
				Type:  jsonpatch.REPLACE,
				Path:  "/hints/3/content",
				Value: "das ist der neue dritte hint",
			},
			jsonpatch.Operation{
				Type:  jsonpatch.REPLACE,
				Path:  "/hints/3/cost",
				Value: 200,
			},
		},
	}
	err = mapper.ApplyPatch("f7c21557-03fc-3e99-bdff-7b065f58b39d", &p, compiler)
	assert.Nil(t, err)
	ex, err = getExercise(mapper, "f7c21557-03fc-3e99-bdff-7b065f58b39d")
	hints = ex["hint_ids"].([]interface{})
	assert.Equal(t, hintId, hints[1].(string))
	assert.Equal(t, 3, len(hints))
	set := arrayToSet(hints)
	assert.False(t, set["164186cb-1252-3672-a015-e8128b999bb4"])
	assert.Equal(t, hints[2].(string), "1bcfd6c3-b269-392a-a7c0-2206f9aefcb6")
	hint, err = getHint(mapper, "1bcfd6c3-b269-392a-a7c0-2206f9aefcb6")
	assert.Nil(t, err)
	assert.Equal(t, "das ist der neue dritte hint", hint["content"])
	assert.Equal(t, 200, int(hint["cost"].(float64)))
}

func getHint(mapper *DataMapper, hintId string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	data, err := mapper.queryIntoBytes("SELECT get_hint($1,$2)", "233804c6-55b8-3807-9733-9c090d75decf", hintId)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&result)
	return result, err
}

func getHint(mapper *DataMapper, hintId string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	data, err := mapper.queryIntoBytes("SELECT get_hint($1,$2)", "233804c6-55b8-3807-9733-9c090d75decf", hintId)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&result)
	return result, err
}

func getExercise(mapper *DataMapper, exerciseId string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	data, err := mapper.queryIntoBytes("SELECT get_one_exercise_as_json($1)", exerciseId)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&result)
	return result, err
}

func arrayToSet(slice []interface{}) map[string]bool {
	result := make(map[string]bool)
	for _, v := range slice {
		result[v.(string)] = true
	}
	return result
}

/*
func TestReplaceDescription(t *testing.T) {
	mapper, err := prepareMapper()
	assert.Nil(t, err)
	assert.Nil(t, resetDatabase(mapper))
	modules := getModules(t, mapper.db)
	p := jsonpatch.Patch{
		Version: 1,
		Operations: []jsonpatch.Operation{
			jsonpatch.Operation{
				Type:  jsonpatch.REPLACE,
				Path:  "/description",
				Value: "urf urf urf",
			},
		},
	}
	compiler := lecturepatch.ForModules()
	err = mapper.ApplyPatch(modules["foo"].Id, &p, compiler)
	assert.Nil(t, err)
	urf, err := getModule(mapper, "urf urf urf")
	assert.Nil(t, err)
	assert.Equal(t, "urf urf urf", urf["description"].(string))
}

func TestAddRemoveVideo(t *testing.T) {
	mapper, err := prepareMapper()
	assert.Nil(t, err)
	assert.Nil(t, resetDatabase(mapper))
	modules := getModules(t, mapper.db)
	videoId := uuid.NewV4().String()
	p := jsonpatch.Patch{
		Version: 1,
		Operations: []jsonpatch.Operation{
			jsonpatch.Operation{
				Type:  jsonpatch.ADD,
				Path:  "/video",
				Value: videoId,
			},
		},
	}
	compiler := lecturepatch.ForModules()
	err = mapper.ApplyPatch(modules["foo"].Id, &p, compiler)
	assert.Nil(t, err)
	foo, err := getModule(mapper, "foo")
	assert.Nil(t, err)
	assert.Equal(t, videoId, foo["video_id"].(string))
	p = jsonpatch.Patch{
		Version: 2,
		Operations: []jsonpatch.Operation{
			jsonpatch.Operation{
				Type: jsonpatch.REMOVE,
				Path: "/video/" + videoId,
			},
		},
	}
	err = mapper.ApplyPatch(modules["foo"].Id, &p, compiler)
	assert.Nil(t, err)
	foo, err = getModule(mapper, "foo")
	assert.Nil(t, err)
	assert.Nil(t, foo["video_id"])
}

func TestAddRemoveScript(t *testing.T) {
	mapper, err := prepareMapper()
	assert.Nil(t, err)
	assert.Nil(t, resetDatabase(mapper))
	modules := getModules(t, mapper.db)
	scriptId := uuid.NewV4().String()
	p := jsonpatch.Patch{
		Version: 1,
		Operations: []jsonpatch.Operation{
			jsonpatch.Operation{
				Type:  jsonpatch.ADD,
				Path:  "/script",
				Value: scriptId,
			},
		},
	}
	compiler := lecturepatch.ForModules()
	err = mapper.ApplyPatch(modules["foo"].Id, &p, compiler)
	assert.Nil(t, err)
	foo, err := getModule(mapper, "foo")
	assert.Nil(t, err)
	assert.Equal(t, scriptId, foo["script_id"].(string))
	p = jsonpatch.Patch{
		Version: 2,
		Operations: []jsonpatch.Operation{
			jsonpatch.Operation{
				Type: jsonpatch.REMOVE,
				Path: "/script/" + scriptId,
			},
		},
	}
	err = mapper.ApplyPatch(modules["foo"].Id, &p, compiler)
	assert.Nil(t, err)
	foo, err = getModule(mapper, "foo")
	assert.Nil(t, err)
	assert.Nil(t, foo["script_id"])
}
func TestAddRemoveExercise(t *testing.T) {
	mapper, err := prepareMapper()
	assert.Nil(t, err)
	assert.Nil(t, resetDatabase(mapper))
	modules := getModules(t, mapper.db)
	exerciseId := uuid.NewV4().String()
	p := jsonpatch.Patch{
		Version: 1,
		Operations: []jsonpatch.Operation{
			jsonpatch.Operation{
				Type: jsonpatch.ADD,
				Path: "/exercises",
				Value: map[string]interface{}{
					"id":      exerciseId,
					"backend": "java",
					"tasks": []string{
						"do something awesome",
						"do something more awesome",
					},
				},
			},
		},
	}
	compiler := lecturepatch.ForModules()
	err = mapper.ApplyPatch(modules["foo"].Id, &p, compiler)
	assert.Nil(t, err)
	foo, err := getModule(mapper, "foo")
	assert.Nil(t, err)
	exercises := toSet(foo["exercises"].([]interface{}), "id")
	for _, v := range []string{exerciseId} {
		assert.True(t, exercises[v])
	}
	p = jsonpatch.Patch{
		Version: 2,
		Operations: []jsonpatch.Operation{
			jsonpatch.Operation{
				Type: jsonpatch.REMOVE,
				Path: "/exercises/" + exerciseId,
			},
		},
	}
	err = mapper.ApplyPatch(modules["foo"].Id, &p, compiler)
	assert.Nil(t, err)
	foo, err = getModule(mapper, "foo")
	assert.Nil(t, err)
	exercises = toSet(foo["exercises"].([]interface{}), "id")
	for _, v := range []string{exerciseId} {
		assert.False(t, exercises[v])
	}
}
*/
