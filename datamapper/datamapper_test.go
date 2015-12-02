package datamapper

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/richterrettich/lecture-service/paginator"
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

func resetDatabase(mapper *DataMapper) error {
	result, err := ioutil.ReadFile("dummy_data.sql")
	if err != nil {
		return err
	}
	parts := strings.Split(string(result), ";")
	tx, err := mapper.db.Begin()
	for _, v := range parts {
		_, err = tx.Exec(v)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil

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
