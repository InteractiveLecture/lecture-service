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

func TestGetModule(t *testing.T) {
	mapper, err := prepareMapper()
	assert.Nil(t, err)
	assert.Nil(t, resetDatabase(mapper))
	result, err := mapper.GetOneModule("98bf99f7-3fed-3fd0-b43e-0b0f376b3607")
	assert.Nil(t, err)
	assert.NotNil(t, result)
	var topics = make(map[string]interface{})
	err = json.NewDecoder(bytes.NewReader(result)).Decode(&topics)
	assert.Nil(t, err)
}

func TestGetModuleTree(t *testing.T) {
	mapper, err := prepareMapper()
	assert.Nil(t, err)
	assert.Nil(t, resetDatabase(mapper))
	dr := paginator.DepthRequest{0, -1, -1}
	result, err := mapper.GetModuleRange("b8c98f3e-bb7c-39e7-a3ce-e479c7892882", dr)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	var topics = make([]map[string]interface{}, 0)
	err = json.NewDecoder(bytes.NewReader(result)).Decode(&topics)
	assert.Nil(t, err)
}

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
}
