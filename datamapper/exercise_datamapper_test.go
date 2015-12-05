package datamapper

import (
	"bytes"
	"encoding/json"
	"log"
	"testing"

	"github.com/richterrettich/lecture-service/paginator"
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
	log.Println(histories)
	historySet := toSet(histories, "exercise_id")
	for _, v := range []string{"233804c6-55b8-3807-9733-9c090d75decf"} {
		assert.True(t, historySet[v])
	}

}

/*
func TestAddRemoveRecommendation(t *testing.T) {
	mapper, err := prepareMapper()
	assert.Nil(t, err)
	assert.Nil(t, resetDatabase(mapper))
	modules := getModules(t, mapper.db)
	p := jsonpatch.Patch{
		Version: 1,
		Operations: []jsonpatch.Operation{
			jsonpatch.Operation{
				Type:  jsonpatch.ADD,
				Path:  "/recommendations",
				Value: modules["bazz"].Id,
			},
		},
	}
	compiler := lecturepatch.ForModules()
	err = mapper.ApplyPatch(modules["foo"].Id, &p, compiler)
	assert.Nil(t, err)
	foo, err := getModule(mapper, "foo")
	assert.Nil(t, err)
	recommendations := toSet(foo["recommendations"].([]interface{}), "id")
	for _, v := range []string{modules["bazz"].Id} {
		assert.True(t, recommendations[v])
	}
	p = jsonpatch.Patch{
		Version: 2,
		Operations: []jsonpatch.Operation{
			jsonpatch.Operation{
				Type: jsonpatch.REMOVE,
				Path: "/recommendations/" + modules["bazz"].Id,
			},
		},
	}
	err = mapper.ApplyPatch(modules["foo"].Id, &p, compiler)
	foo, err = getModule(mapper, "foo")
	assert.Nil(t, err)
	recommendations = toSet(foo["recommendations"].([]interface{}), "id")
	for _, v := range []string{modules["bazz"].Id} {
		assert.False(t, recommendations[v])
	}
}

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

func getModule(mapper *DataMapper, description string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	data, err := mapper.queryIntoBytes("SELECT d.details from module_details d inner join modules m on m.id = d.id where m.description = $1", description)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&result)
	return result, err
}*/
