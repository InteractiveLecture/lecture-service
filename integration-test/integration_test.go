package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/richterrettich/jsonpatch"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

var tokens = make(map[string]string)

//TODO test moving/removing root modules

func TestGetTopics(t *testing.T) {

	//INSERT USER
	userId, userUsername := RegisterNewUser(t, "user")
	officerId, officerUsername := RegisterNewUser(t, "officer")
	assistantId, assistantUsername := RegisterNewUser(t, "assistant")
	assistantId2, assistantUsername2 := RegisterNewUser(t, "assistant")

	path := "/lecture-service/users/" + userId + "/balances"
	resp := getAuthorized(t, userUsername, path)
	require.Equal(t, 200, resp.StatusCode)
	balances := readArrayJsonResult(t, resp)
	path = "/lecture-service/topics?page=0&size=100000"
	topics := readArrayJsonResult(t, getAuthorized(t, userUsername, path))
	require.Equal(t, len(topics), len(balances))

	// admin creates topic
	path = "/lecture-service/topics"
	newTopicId := uuid.NewV4().String()
	newTopic := `{
		"id": "` + newTopicId + `",
		"name" : "datenbanken",
		"description": "Eine Einführung in SQL-Datenbanken",
		"officers" : ["` + officerId + `"]
	}`
	PostAuthorizedAndCheckStatusCode(t, "admin", path, newTopic, 201)
	// add modules
	path = "/lecture-service/topics/" + newTopicId
	operations := make([]jsonpatch.Operation, 0)
	newModuleIds := make([]string, 0)
	newModuleIds = append(newModuleIds, addNewModule(&operations, "NewModule1"))
	newModuleIds = append(newModuleIds, addNewModule(&operations, "NewModule2", newModuleIds[0]))
	newModuleIds = append(newModuleIds, addNewModule(&operations, "NewModule3", newModuleIds[1]))
	newModuleIds = append(newModuleIds, addNewModule(&operations, "NewModule4", newModuleIds[2]))
	newModuleIds = append(newModuleIds, addNewModule(&operations, "NewModule5", newModuleIds[2]))                  //same parent
	newModuleIds = append(newModuleIds, addNewModule(&operations, "NewModule6", newModuleIds[3], newModuleIds[4])) //multiple parents
	//Patch topic
	operations = append(operations,
		jsonpatch.Operation{
			Type:  jsonpatch.ADD,
			Path:  "/assistants",
			Value: assistantId,
		},
		jsonpatch.Operation{
			Type:  jsonpatch.ADD,
			Path:  "/assistants",
			Value: assistantId2,
		},
		jsonpatch.Operation{
			Type: jsonpatch.REMOVE,
			Path: "/assistants/" + assistantId2,
		},
		jsonpatch.Operation{
			Type:  jsonpatch.REPLACE,
			Path:  "/modules/" + newModuleIds[2] + "/parents/tree",
			Value: []string{newModuleIds[0]},
		},
		jsonpatch.Operation{
			Type:  jsonpatch.REPLACE,
			Path:  "/modules/" + newModuleIds[4] + "/parents",
			Value: []string{newModuleIds[1]},
		},
	)

	topicPatch := jsonpatch.Patch{
		Operations: operations,
		Version:    1,
	}

	patchJson, _ := json.Marshal(topicPatch)
	PatchAuthorizedAndCheckStatusCode(t, officerUsername, path, string(patchJson), 200)

	//TEST MODULES
	path = "/lecture-service/topics/" + newTopicId + "/modules"
	checkUnauthorized(t, path)
	resp = getAuthorized(t, userUsername, path)
	modules := readArrayJsonResult(t, resp)
	require.Equal(t, len(newModuleIds), len(modules))
	testModule := findLocalById(t, modules, newModuleIds[2], "id")
	require.Equal(t, "/"+newModuleIds[0]+"/"+newModuleIds[2], testModule["paths"].([]interface{})[0])
	testModule = findLocalById(t, modules, newModuleIds[3], "id")
	require.Equal(t, "/"+newModuleIds[0]+"/"+newModuleIds[2]+"/"+newModuleIds[3], testModule["paths"].([]interface{})[0])
	testModule = findLocalById(t, modules, newModuleIds[4], "id")
	require.Equal(t, fmt.Sprintf("/%s/%s/%s", newModuleIds[0], newModuleIds[1], newModuleIds[4]), testModule["paths"].([]interface{})[0])

	//Patch topic again
	path = "/lecture-service/topics/" + newTopicId
	operations = make([]jsonpatch.Operation, 0)
	newModuleIds = append(newModuleIds, addNewModule(&operations, "NewModule7", newModuleIds[1]))
	newModuleIds = append(newModuleIds, addNewModule(&operations, "NewModule8", newModuleIds[6]))
	newModuleIds = append(newModuleIds, addNewModule(&operations, "NewModule9", newModuleIds[6]))
	newModuleIds = append(newModuleIds, addNewModule(&operations, "NewModule10", newModuleIds[7], newModuleIds[8]))
	operations = append(operations,
		jsonpatch.Operation{
			Type: jsonpatch.REMOVE,
			Path: "/modules/" + newModuleIds[6],
		},
	)
	topicPatch = jsonpatch.Patch{
		Operations: operations,
		Version:    2,
	}

	patchJson, _ = json.Marshal(topicPatch)
	PatchAuthorizedAndCheckStatusCode(t, officerUsername, path, string(patchJson), 200)

	//check topic modules again
	path = "/lecture-service/topics/" + newTopicId + "/modules"
	resp = getAuthorized(t, userUsername, path)
	modules = readArrayJsonResult(t, resp)
	testModule = findLocalById(t, modules, newModuleIds[9], "id")
	require.Equal(t, len(newModuleIds)-1, len(modules))
	require.Equal(t, 2, len(testModule["paths"].([]interface{})))
	require.Contains(t, testModule["paths"].([]interface{}), fmt.Sprintf("/%s/%s/%s/%s", newModuleIds[0], newModuleIds[1], newModuleIds[7], newModuleIds[9]))
	require.Contains(t, testModule["paths"].([]interface{}), fmt.Sprintf("/%s/%s/%s/%s", newModuleIds[0], newModuleIds[1], newModuleIds[8], newModuleIds[9]))
	//Patch one module
	path = "/lecture-service/modules/" + newModuleIds[0]
	operations = make([]jsonpatch.Operation, 0)
	newExerciseId := addNewExercise(&operations)
	toRemoveExerciseId := addNewExercise(&operations)
	newVideoId := uuid.NewV4().String()
	newScriptId := uuid.NewV4().String()
	operations = append(operations,
		jsonpatch.Operation{
			Type: jsonpatch.REMOVE,
			Path: "/exercises/" + toRemoveExerciseId,
		},
		jsonpatch.Operation{
			Type:  jsonpatch.REPLACE,
			Path:  "/description",
			Value: "hugo",
		},
		jsonpatch.Operation{
			Type:  jsonpatch.ADD,
			Path:  "/recommendations",
			Value: newModuleIds[1],
		},
		jsonpatch.Operation{
			Type: jsonpatch.REMOVE,
			Path: "/recommendations/" + newModuleIds[1],
		},
		jsonpatch.Operation{
			Type:  jsonpatch.ADD,
			Path:  "/video",
			Value: newVideoId,
		},
		jsonpatch.Operation{
			Type: jsonpatch.REMOVE,
			Path: "/video/" + newVideoId,
		},
		jsonpatch.Operation{
			Type:  jsonpatch.ADD,
			Path:  "/script",
			Value: newScriptId,
		},
		jsonpatch.Operation{
			Type: jsonpatch.REMOVE,
			Path: "/script/" + newScriptId,
		},
	)
	modulePatch := jsonpatch.Patch{
		Operations: operations,
		Version:    1,
	}
	patchJson, _ = json.Marshal(modulePatch)
	PatchAuthorizedAndCheckStatusCode(t, officerUsername, path, string(patchJson), 200)
	// Assistants and users should not be able to do some of the above operations.
	PatchAuthorizedAndCheckStatusCode(t, assistantUsername, path, string(patchJson), 401)
	PatchAuthorizedAndCheckStatusCode(t, userUsername, path, string(patchJson), 401)

	//ADD tasks to exercise
	path = "/lecture-service/exercises/" + newExerciseId
	operations = make([]jsonpatch.Operation, 0)
	_ = addNewTask(&operations, "bla blubb", 1)
	_ = addNewTask(&operations, "asdf asdf", 2)
	_ = addNewTask(&operations, "one last task", 3)
	exercisePatch := jsonpatch.Patch{
		Operations: operations,
		Version:    1,
	}
	patchJson, _ = json.Marshal(exercisePatch)
	PatchAuthorizedAndCheckStatusCode(t, officerUsername, path, string(patchJson), 200)

	//  add additional exercises
	path = "/lecture-service/modules/" + newModuleIds[0]
	operations = make([]jsonpatch.Operation, 0)
	toRemoveExerciseId = addNewExercise(&operations)
	newExerciseId = addNewExercise(&operations)
	operations = append(operations,
		jsonpatch.Operation{
			Type: jsonpatch.REMOVE,
			Path: "/exercises/" + toRemoveExerciseId,
		},
	)
	modulePatch = jsonpatch.Patch{
		Operations: operations,
		Version:    2,
	}
	patchJson, _ = json.Marshal(modulePatch)
	PatchAuthorizedAndCheckStatusCode(t, assistantUsername, path, string(patchJson), 200)

	path = "/lecture-service/exercises/" + newExerciseId
	operations = make([]jsonpatch.Operation, 0)
	_ = addNewTask(&operations, "bla blubb", 1)
	_ = addNewTask(&operations, "asdf asdf", 2)
	_ = addNewTask(&operations, "one last task", 3)

	operations = append(operations,
		jsonpatch.Operation{
			Type: jsonpatch.MOVE,
			Path: "/tasks/2",
			From: "/tasks/1",
		},
		jsonpatch.Operation{
			Type: jsonpatch.REMOVE,
			Path: "/tasks/2",
		},
		jsonpatch.Operation{
			Type:  jsonpatch.REPLACE,
			Path:  "/tasks/1/content",
			Value: "urf",
		},
	)
	_ = addNewHint(&operations, "a hint", 1, 1)
	_ = addNewHint(&operations, "another hint", 1, 2)
	_ = addNewHint(&operations, "one last hint", 1, 3)
	operations = append(operations,
		jsonpatch.Operation{
			Type: jsonpatch.MOVE,
			Path: "/tasks/1/hints/2",
			From: "/tasks/1/hints/1",
		},
		jsonpatch.Operation{
			Type: jsonpatch.REMOVE,
			Path: "/tasks/1/hints/1",
		},
	)
	exercisePatch = jsonpatch.Patch{
		Operations: operations,
		Version:    1,
	}
	patchJson, _ = json.Marshal(exercisePatch)
	PatchAuthorizedAndCheckStatusCode(t, assistantUsername, path, string(patchJson), 200)
	// assistant 2 should not have accessrights
	PatchAuthorizedAndCheckStatusCode(t, assistantUsername2, path, string(patchJson), 401)

	// lets TEST one module
	path = "/lecture-service/modules/" + newModuleIds[0]
	resp = getAuthorized(t, userUsername, path)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	module := readSingleJsonResult(t, resp)
	require.Nil(t, module["script_id"])
	require.Nil(t, module["video_id"])

	exercises := module["exercises"].([]interface{})
	require.Equal(t, 2, len(exercises))
	exercise := findRawLocalById(t, exercises, newExerciseId, "id")
	require.Equal(t, "java", exercise["backend"].(string))

	path = "/lecture-service/exercises/" + exercise["id"].(string) + "/start"
	PostAuthorizedAndCheckStatusCode(t, userUsername, path, "", http.StatusOK)
	path = "/lecture-service/users/" + userId + "/exercises"
	checkUnauthorized(t, path)
	resp = getAuthorized(t, userUsername, path)
	exerciseHistory := readArrayJsonResult(t, resp)
	require.Equal(t, 1, len(exerciseHistory))
	require.Equal(t, exerciseHistory[0]["event_type"], "BEGIN")
	path = "/lecture-service/users/" + userId + "/modules"
	checkUnauthorized(t, path)
	resp = getAuthorized(t, userUsername, path)
	moduleHistory := readArrayJsonResult(t, resp)
	require.Equal(t, 1, len(moduleHistory))
	require.Equal(t, moduleHistory[0]["event_type"], "BEGIN")

	tasks := exercise["tasks"].([]interface{})
	require.Equal(t, 2, len(tasks))
	task := tasks[0].(map[string]interface{})
	require.Equal(t, "urf", task["content"].(string))
	require.Equal(t, "one last task", tasks[1].(map[string]interface{})["content"].(string))
	hints := task["hints"].([]interface{})
	require.Equal(t, 2, len(hints))
	hintId := hints[0].(string)
	secondHintId := hints[1].(string)

	path = "/lecture-service/hints/" + hintId
	checkUnauthorized(t, path)
	GetAuthorizedAndCheckStatusCode(t, userUsername, path, 402)
	PostAuthorizedAndCheckStatusCode(t, userUsername, path, "", http.StatusOK)
	PostAuthorizedAndCheckStatusCode(t, userUsername, path, "", http.StatusConflict)
	resp = getAuthorized(t, userUsername, path)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	realHint := readSingleJsonResult(t, resp)
	require.Equal(t, "a hint", realHint["content"].(string))

	path = "/lecture-service/hints/" + secondHintId
	GetAuthorizedAndCheckStatusCode(t, userUsername, path, 402)
	PostAuthorizedAndCheckStatusCode(t, userUsername, path, "", 420)

	path = "/lecture-service/hints/" + uuid.NewV4().String()
	PostAuthorizedAndCheckStatusCode(t, userUsername, path, "", http.StatusNotFound)

	// TEST COMPLETE TASK

	path = "/nats-remote/task-backend.task-finished"
	for _, v := range exercises {
		exercise := v.(map[string]interface{})
		for _, ta := range exercise["tasks"].([]interface{}) {
			task := ta.(map[string]interface{})
			result := map[string]interface{}{
				"userId": userId,
				"taskId": task["id"],
			}
			re, err := json.Marshal(result)
			require.Nil(t, err)
			PostAuthorizedAndCheckStatusCode(t, userUsername, path, string(re), 200)
			time.Sleep(500 * time.Millisecond)
		}
	}

	path = "/lecture-service/users/" + userId + "/exercises"
	resp = getAuthorized(t, userUsername, path)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	exerciseHistory = readArrayJsonResult(t, resp)
	require.NotZero(t, exerciseHistory)
	require.Equal(t, "FINISH", exerciseHistory[len(exerciseHistory)-1]["event_type"])

	path = "/lecture-service/users/" + userId + "/modules"
	resp = getAuthorized(t, userUsername, path)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	moduleHistory = readArrayJsonResult(t, resp)
	require.NotZero(t, exerciseHistory)
	require.Equal(t, "FINISH", moduleHistory[len(moduleHistory)-1]["event_type"])

	path = "/lecture-service/users/" + userId + "/balances"
	resp = getAuthorized(t, userUsername, path)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	balances = readArrayJsonResult(t, resp)
	balance := findLocalById(t, balances, newTopicId, "topic_id")
	require.Equal(t, float64(500), balance["amount"])
}

func findRawLocalById(t *testing.T, collection []interface{}, id, idField string) map[string]interface{} {
	var result map[string]interface{}
	for _, v := range collection {
		val := v.(map[string]interface{})
		if val[idField].(string) == id {
			result = val
			break
		}
	}
	require.NotNil(t, result)
	return result
}
func findLocalById(t *testing.T, collection []map[string]interface{}, id, idField string) map[string]interface{} {
	var result map[string]interface{}
	for _, v := range collection {
		if v[idField].(string) == id {
			result = v
			break
		}
	}
	require.NotNil(t, result)
	return result
}

func addNewHint(operations *[]jsonpatch.Operation, hint string, taskPosition, position int) string {
	newHint := map[string]interface{}{
		"id":       uuid.NewV4().String(),
		"content":  hint,
		"position": position,
		"cost":     100,
	}
	addSubEntity(operations, newHint, fmt.Sprintf("/tasks/%d/hints", taskPosition))
	return newHint["id"].(string)
}
func addNewTask(operations *[]jsonpatch.Operation, task string, position int) string {
	newTask := map[string]interface{}{
		"id":       uuid.NewV4().String(),
		"content":  task,
		"position": position,
	}
	addSubEntity(operations, newTask, "/tasks")
	return newTask["id"].(string)
}

func addNewExercise(operations *[]jsonpatch.Operation) string {
	newExercise := map[string]interface{}{
		"id":      uuid.NewV4().String(),
		"backend": "java",
	}
	addSubEntity(operations, newExercise, "/exercises")
	return newExercise["id"].(string)
}

func addSubEntity(operations *[]jsonpatch.Operation, entity map[string]interface{}, path string) {
	*operations = append(*operations, jsonpatch.Operation{
		Type:  jsonpatch.ADD,
		Path:  path,
		Value: entity,
	})
}
func addNewModule(operations *[]jsonpatch.Operation, description string, parents ...string) string {
	if parents == nil {
		parents = []string{}
	}
	newModule := map[string]interface{}{
		"id":          uuid.NewV4().String(),
		"description": description,
		"video_id":    uuid.NewV4().String(),
		"script_id":   uuid.NewV4().String(),
		"parents":     parents,
	}
	addSubEntity(operations, newModule, "/modules")
	return newModule["id"].(string)
}

func RegisterNewUser(t *testing.T, authorities ...string) (string, string) {
	path := "/authentication-service/users"
	username := uuid.NewV4().String()
	userId := uuid.NewV4().String()
	user := `{
		"id" : "` + userId + `",
		"username": "` + username + `",
		"password": "` + username + `",
		"enabled": true,
		"authorities": [`
	for _, v := range authorities {
		user = user + `{"authority":"` + v + `",`
	}
	user = strings.TrimRight(user, ",")
	user = user + "}]}"
	PostUnauthorizedAndCheckStatusCode(t, path, user, 204)
	return userId, username
}

func PatchAuthorizedAndCheckStatusCode(t *testing.T, user, path, body string, expecedCode int, headers ...string) {
	headers = append(headers, "Authorization", "Bearer "+getToken(user))
	resp := PatchAuthorized(t, path, body, headers...)
	defer resp.Body.Close()
	require.Equal(t, expecedCode, resp.StatusCode)
}

func PatchAuthorized(t *testing.T, path, body string, headers ...string) *http.Response {
	return sendRequest(t, "PATCH", path, strings.NewReader(body), headers...)
}

func PostUnauthorizedAndCheckStatusCode(t *testing.T, path, body string, expecedCode int, headers ...string) {
	resp := postUnauthorized(t, path, body, headers...)
	defer resp.Body.Close()
	require.Equal(t, expecedCode, resp.StatusCode)
}

func PostAuthorizedAndCheckStatusCode(t *testing.T, user, path, body string, expecedCode int, headers ...string) {
	resp := postAuthorized(t, user, path, body, headers...)
	defer resp.Body.Close()
	require.Equal(t, expecedCode, resp.StatusCode)
}

func GetAuthorizedAndCheckStatusCode(t *testing.T, user, path string, expecedCode int, headers ...string) {
	resp := getAuthorized(t, user, path, headers...)
	defer resp.Body.Close()
	require.Equal(t, expecedCode, resp.StatusCode)
}

func postUnauthorized(t *testing.T, path, body string, headers ...string) *http.Response {
	headers = append(headers, "Content-Type", "application/json;charset=UTF-8")
	return sendRequest(t, "POST", path, strings.NewReader(body), headers...)
}

func postAuthorized(t *testing.T, user, path, body string, headers ...string) *http.Response {
	headers = append(headers, "Authorization", "Bearer "+getToken(user), "Content-Type", "application/json;charset=UTF-8")
	return sendRequest(t, "POST", path, strings.NewReader(body), headers...)
}

func checkUnauthorized(t *testing.T, path string) {
	resp := getUnauthorized(t, path)
	defer resp.Body.Close()
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func readSingleJsonResult(t *testing.T, resp *http.Response) map[string]interface{} {
	defer resp.Body.Close()
	result := make(map[string]interface{})
	err := json.NewDecoder(resp.Body).Decode(&result)
	require.Nil(t, err)
	return result
}

func readArrayJsonResult(t *testing.T, resp *http.Response) []map[string]interface{} {
	defer resp.Body.Close()
	result := make([]map[string]interface{}, 0)
	err := json.NewDecoder(resp.Body).Decode(&result)
	require.Nil(t, err)
	require.NotZero(t, len(result))
	return result
}

func getUnauthorized(t *testing.T, path string) *http.Response {
	return sendRequest(t, "GET", path, nil)
}

func getAuthorized(t *testing.T, user, path string, headers ...string) *http.Response {
	headers = append(headers, "Authorization", "Bearer "+getToken(user))
	return sendRequest(t, "GET", path, nil, headers...)
}

func sendRequest(t *testing.T, requestType, path string, body io.Reader, headers ...string) *http.Response {
	host := getHost()
	req, err := http.NewRequest(requestType, "http://"+host+path, body)
	require.Nil(t, err)
	if len(headers)%2 != 0 {
		panic(fmt.Errorf("wrong number of header arguments!"))
	}
	for i := 0; i < len(headers); i = i + 2 {
		req.Header.Add(headers[i], headers[i+1])
	}
	client := http.Client{}
	resp, err := client.Do(req)
	require.Nil(t, err)
	return resp
}

func getToken(user string) string {
	if val, ok := tokens[user]; ok {
		return val
	}
	host := getHost()
	authString := "client_id=user-web-client&client_secret=user-web-client-secret&grant_type=password&username=" + user + "&password=" + user
	reader := strings.NewReader(authString)
	result, err := http.Post("http://"+host+"/authentication-service/oauth/token", "application/x-www-form-urlencoded", reader)
	if err != nil {
		panic(err)
	}
	defer result.Body.Close()
	if result.StatusCode != 200 {
		panic(errors.New("expected statuscode 200 from authentication-service, but got: " + result.Status))
	}
	token := make(map[string]interface{})
	err = json.NewDecoder(result.Body).Decode(&token)
	if err != nil {
		panic(err)
	}
	tokens[user] = token["access_token"].(string)
	return tokens[user]
}

func getHost() string {
	return os.Getenv("DH")
}
