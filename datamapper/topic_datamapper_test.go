package datamapper

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/richterrettich/jsonpatch"
	"github.com/richterrettich/lecture-service/lecturepatch"
	"github.com/richterrettich/lecture-service/paginator"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestQueryTopics(t *testing.T) {
	mapper, err := prepareMapper()
	assert.Nil(t, err)
	assert.Nil(t, resetDatabase(mapper))
	pr := paginator.PageRequest{0, 10, nil}
	result, err := mapper.GetTopicsPage(pr)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	var topics = make([]map[string]interface{}, 0)
	err = json.NewDecoder(bytes.NewReader(result)).Decode(&topics)
	assert.Nil(t, err)
}

func TestQuerySingleTopic(t *testing.T) {
	mapper, err := prepareMapper()
	assert.Nil(t, err)
	assert.Nil(t, resetDatabase(mapper))
	result, err := mapper.GetOneTopic("b8c98f3e-bb7c-39e7-a3ce-e479c7892882")
	assert.Nil(t, err)
	assert.NotNil(t, result)
	var topics = make(map[string]interface{})
	err = json.NewDecoder(bytes.NewReader(result)).Decode(&topics)
	assert.Nil(t, err)
}

func TestGetBalances(t *testing.T) {
	mapper, err := prepareMapper()
	assert.Nil(t, err)
	assert.Nil(t, resetDatabase(mapper))
	result, err := mapper.GetTopicBalances("f20919fa-08bd-3a8d-9e3c-e3406c680162")
	assert.Nil(t, err)
	assert.NotNil(t, result)
	var topics = make([]map[string]interface{}, 0)
	err = json.NewDecoder(bytes.NewReader(result)).Decode(&topics)
	assert.Nil(t, err)
}

func TestGetHintHistory(t *testing.T) {
	mapper, err := prepareMapper()
	assert.Nil(t, err)
	assert.Nil(t, resetDatabase(mapper))
	pr := paginator.PageRequest{0, 10, nil}
	result, err := mapper.GetHintHistory("f20919fa-08bd-3a8d-9e3c-e3406c680162", pr, "")
	assert.Nil(t, err)
	assert.NotNil(t, result)
	var topics = make([]map[string]interface{}, 0)
	err = json.NewDecoder(bytes.NewReader(result)).Decode(&topics)
	assert.Nil(t, err)

	result, err = mapper.GetHintHistory("f20919fa-08bd-3a8d-9e3c-e3406c680162", pr, "f7c21557-03fc-3e99-bdff-7b065f58b39d")
	assert.Nil(t, err)
	err = json.NewDecoder(bytes.NewReader(result)).Decode(&topics)
	assert.Nil(t, err)

}

func TestGetModuleHistory(t *testing.T) {
	mapper, err := prepareMapper()
	assert.Nil(t, err)
	assert.Nil(t, resetDatabase(mapper))
	pr := paginator.PageRequest{0, 10, nil}
	result, err := mapper.GetModuleHistory("f20919fa-08bd-3a8d-9e3c-e3406c680162", pr, "")
	assert.Nil(t, err)
	assert.NotNil(t, result)
	var topics = make([]map[string]interface{}, 0)
	err = json.NewDecoder(bytes.NewReader(result)).Decode(&topics)
	assert.Nil(t, err)
	result, err = mapper.GetModuleHistory("f20919fa-08bd-3a8d-9e3c-e3406c680162", pr, "b8c98f3e-bb7c-39e7-a3ce-e479c7892882")
	assert.Nil(t, err)
	err = json.NewDecoder(bytes.NewReader(result)).Decode(&topics)
	assert.Nil(t, err)
}

/*
TOPIC 1:
			FOO
			 |
			BAR
			 |
			BLI
			/ \
  BLUBB  BLA
	  \    /
		 BAZZ
		 TOPIC 2:
		 FOOBARBAZZ
*/
func TestMoveSingle(t *testing.T) {
	mapper, err := prepareMapper()
	defer mapper.db.Close()
	assert.Nil(t, err)
	assert.Nil(t, resetDatabase(mapper))
	db := mapper.db
	modules := getModules(t, db)
	assert.Equal(t, len(modules), 7)
	p := jsonpatch.Patch{
		Version: 1,
		Operations: []jsonpatch.Operation{
			jsonpatch.Operation{
				Type:  jsonpatch.REPLACE,
				Path:  fmt.Sprintf("/modules/%s/parents", modules["bazz"].Id),
				Value: []string{modules["foo"].Id},
			},
			jsonpatch.Operation{
				Type:  jsonpatch.REPLACE,
				Path:  fmt.Sprintf("/modules/%s/parents", modules["foo"].Id),
				Value: []string{modules["bla"].Id},
			},
		},
	}
	compiler := lecturepatch.ForTopics()
	err = mapper.ApplyPatch(modules["foo"].topicId, &p, compiler)
	assert.Nil(t, err)
	modules = getModules(t, db)
	assert.Equal(t, len(modules), 7) //lenght shouldnt have changed.
	assert.Equal(t, 1, modules["bazz"].level)
	modules = getModules(t, db)
	assert.Equal(t, 3, modules["foo"].level)
	assert.Equal(t, 0, modules["bar"].level)
	for k, v := range modules {
		if k != "foobarbazz" {
			assert.True(t, strings.HasPrefix(v.paths[0], "/"+modules["bar"].Id))
		}
	}
	_, err = db.Exec(`SELECT check_version($1,$2,$3)`, modules["foo"].topicId, "topics", 2)
	assert.Nil(t, err)
}

func TestMoveTree(t *testing.T) {
	mapper, err := prepareMapper()
	defer mapper.db.Close()
	assert.Nil(t, err)
	assert.Nil(t, resetDatabase(mapper))
	db := mapper.db
	modules := getModules(t, db)
	p := jsonpatch.Patch{
		Version: 1,
		Operations: []jsonpatch.Operation{
			jsonpatch.Operation{
				Type:  jsonpatch.REPLACE,
				Path:  fmt.Sprintf("/modules/%s/parents/tree", modules["bli"].Id),
				Value: []string{modules["foo"].Id},
			},
		},
	}
	compiler := lecturepatch.ForTopics()
	err = mapper.ApplyPatch(modules["foo"].topicId, &p, compiler)
	assert.Nil(t, err)

	assert.Nil(t, err)
	modules = getModules(t, db)
	for i := 0; i < 2; i++ {
		assert.False(t, strings.Contains(modules["bazz"].paths[i], modules["bar"].Id))
	}
	assert.Equal(t, 1, modules["bar"].level)
	_, err = db.Exec(`SELECT move_module_tree($1,$2,$3)`, modules["bli"].topicId, modules["bli"].Id, modules["foo"].Id)
	assert.Nil(t, err)
	assert.Equal(t, modules["bli"].Id, getDirectParents(modules["bla"])[0])
}

func TestInsertModule(t *testing.T) {
	mapper, err := prepareMapper()
	defer mapper.db.Close()
	assert.Nil(t, err)
	assert.Nil(t, resetDatabase(mapper))
	db := mapper.db
	modules := getModules(t, db)
	p := jsonpatch.Patch{
		Version: 1,
		Operations: []jsonpatch.Operation{
			jsonpatch.Operation{
				Type: jsonpatch.ADD,
				Path: "/modules",
				Value: map[string]interface{}{
					"id":          uuid.NewV4(),
					"description": "hugo",
					"video_id":    uuid.NewV4(),
					"script_id":   uuid.NewV4(),
					"parents":     []string{modules["foo"].Id},
				},
			},
		},
	}

	compiler := lecturepatch.ForTopics()
	err = mapper.ApplyPatch(modules["foo"].topicId, &p, compiler)
	assert.Nil(t, err)

	modules = getModules(t, db)
	val, ok := modules["hugo"]
	assert.True(t, ok)
	assert.Equal(t, modules["foo"].Id, getDirectParents(val)[0])
	assert.Equal(t, fmt.Sprintf("/%s/%s", modules["foo"].Id, val.Id), val.paths[0])
	assert.Equal(t, modules["foo"].Id, getDirectParents(modules["bar"])[0])
}

func TestDeleteModule(t *testing.T) {
	mapper, err := prepareMapper()
	defer mapper.db.Close()
	assert.Nil(t, err)
	assert.Nil(t, resetDatabase(mapper))
	db := mapper.db
	modules := getModules(t, db)
	p := jsonpatch.Patch{
		Version: 1,
		Operations: []jsonpatch.Operation{
			jsonpatch.Operation{
				Type: jsonpatch.REMOVE,
				Path: "/modules/" + modules["bli"].Id,
			},
			jsonpatch.Operation{
				Type: jsonpatch.REMOVE,
				Path: "/modules/" + modules["foo"].Id,
			},
			jsonpatch.Operation{
				Type: jsonpatch.REMOVE,
				Path: "/modules/" + modules["bar"].Id,
			},
		},
	}
	compiler := lecturepatch.ForTopics()
	err = mapper.ApplyPatch(modules["foo"].topicId, &p, compiler)
	assert.Nil(t, err)
	modules = getModules(t, db)
	assert.Equal(t, 4, len(modules))
	assert.Equal(t, modules["bla"].Id, getDirectParents(modules["blubb"])[0])
	assert.Equal(t, fmt.Sprintf("/%s/%s/%s", modules["bla"].Id, modules["blubb"].Id, modules["bazz"].Id), modules["bazz"].paths[0])
}

func TestDeleteModuleTree(t *testing.T) {
	mapper, err := prepareMapper()
	assert.Nil(t, err)
	assert.Nil(t, resetDatabase(mapper))
	db := mapper.db
	modules := getModules(t, db)
	p := jsonpatch.Patch{
		Version: 1,
		Operations: []jsonpatch.Operation{
			jsonpatch.Operation{
				Type: jsonpatch.ADD,
				Path: "/modules",
				Value: map[string]interface{}{
					"id":          uuid.NewV4(),
					"description": "hugo",
					"video_id":    uuid.NewV4(),
					"script_id":   uuid.NewV4(),
					"parents":     []string{modules["foo"].Id},
				},
			},
			jsonpatch.Operation{
				Type: jsonpatch.REMOVE,
				Path: "/modules/" + modules["bar"].Id + "/tree",
			},
		},
	}
	compiler := lecturepatch.ForTopics()
	err = mapper.ApplyPatch(modules["foo"].topicId, &p, compiler)
	assert.Nil(t, err)
	modules = getModules(t, db)
	assert.Equal(t, 3, len(modules))
	for _, v := range []string{"foo", "hugo"} {
		_, ok := modules[v]
		assert.True(t, ok)
	}
	assert.Equal(t, fmt.Sprintf("/%s/%s", modules["foo"].Id, modules["hugo"].Id), modules["hugo"].paths[0])
}

func TestAddOfficer(t *testing.T) {
	mapper, err := prepareMapper()
	assert.Nil(t, err)
	assert.Nil(t, resetDatabase(mapper))
	officerId := uuid.NewV4().String()
	mapper.AddOfficer("b8c98f3e-bb7c-39e7-a3ce-e479c7892882", officerId)
	data, err := mapper.GetOneTopic("b8c98f3e-bb7c-39e7-a3ce-e479c7892882")
	assert.Nil(t, err)
	var topic = make(map[string]interface{})
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&topic)
	assert.Nil(t, err)
	authorities := toSet(topic["authorities"].([]interface{}), "user_id")
	for _, v := range []string{officerId} {
		assert.True(t, authorities[v])
	}
}

func TestRemoveOfficer(t *testing.T) {
	mapper, err := prepareMapper()
	assert.Nil(t, err)
	assert.Nil(t, resetDatabase(mapper))
	officerId := "233804c6-55b8-3807-9733-9c090d75decf"
	err = mapper.RemoveOfficer("b8c98f3e-bb7c-39e7-a3ce-e479c7892882", officerId)
	assert.Nil(t, err)
	data, err := mapper.GetOneTopic("b8c98f3e-bb7c-39e7-a3ce-e479c7892882")
	assert.Nil(t, err)
	var topic = make(map[string]interface{})
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&topic)
	assert.Nil(t, err)

	authorities := toSet(topic["authorities"].([]interface{}), "user_id")
	for _, v := range []string{officerId} {
		assert.False(t, authorities[v])
	}
}

func TestAddTopic(t *testing.T) {
	mapper, err := prepareMapper()
	assert.Nil(t, err)
	assert.Nil(t, resetDatabase(mapper))

	topic := make(map[string]interface{})

	topicId := uuid.NewV4().String()
	firstOfficerId := uuid.NewV4().String()
	secondOfficerId := uuid.NewV4().String()
	topic["id"] = topicId
	topic["name"] = "Interne Unternehmensrechnung"
	topic["description"] = "lorem ipsum dolor"
	topic["officers"] = []string{firstOfficerId, secondOfficerId}
	err = mapper.CreateTopic(topic)
	assert.Nil(t, err)
	data, err := mapper.GetOneTopic(topicId)
	assert.Nil(t, err)
	var newTopic = make(map[string]interface{})
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&newTopic)
	assert.Nil(t, err)
	authorities := toSet(newTopic["authorities"].([]interface{}), "user_id")
	for _, v := range []string{firstOfficerId, secondOfficerId} {
		assert.True(t, authorities[v])
	}
}

func toSet(input []interface{}, idField string) map[string]bool {
	result := make(map[string]bool)
	for _, v := range input {
		castedEntity := v.(map[string]interface{})
		result[castedEntity[idField].(string)] = true
	}
	return result
}

//Projection functions

func getDirectParents(m module) []string {
	var result = make([]string, 0)
	for _, v := range m.paths {
		parts := strings.Split(v, "/")
		result = append(result, parts[len(parts)-2])
	}
	return result
}

func getModules(t *testing.T, db *sql.DB) map[string]module {
	rows, err := db.Query(`SELECT id,description,level,paths, topic_id FROM module_trees order by level`)
	assert.Nil(t, err)
	defer rows.Close()
	var id, description, paths, topicId string
	var level int
	var modules = make(map[string]module)
	for rows.Next() {
		err = rows.Scan(&id, &description, &level, &paths, &topicId)
		assert.Nil(t, err)
		modules[description] = module{id, description, parseArray(paths), topicId, level}
	}
	return modules
}

type module struct {
	Id          string
	description string
	paths       []string
	topicId     string
	level       int
}

func parseArray(arr string) []string {
	step1 := arr[1 : len(arr)-1]
	return strings.Split(step1, ",")
}

func prepareMapper() (*DataMapper, error) {
	config := DefaultConfig()
	host := os.Getenv("PGHOST")
	if host != "" {
		config.Host = host
	}
	mapper, err := New(config)

	if err != nil {
		log.Println("error with database connection")
		return nil, err
	}
	return mapper, nil
}

func resetDatabase(mapper *DataMapper) error {
	result, err := ioutil.ReadFile("dummy_data.sql")
	if err != nil {
		return err
	}
	parts := strings.Split(string(result), ";")
	tx, err := mapper.db.Begin()
	for _, v := range parts {
		_, err := tx.Exec(v)
		if err != nil {
			tx.Rollback()
			log.Println(err)
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	//_, err := exec.Command("./reset_database.sh").Output()

	return err

}
