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

	"github.com/richterrettich/jsonpatch"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

var tokens = make(map[string]string)

func TestGetTopics(t *testing.T) {

	//INSERT USER
	userId, userUsername := RegisterNewUser(t, "user")
	officerId, officerUsername := RegisterNewUser(t, "officer")
	_, _ = RegisterNewUser(t, "assistant")

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
	newAssistantId := uuid.NewV4().String()
	//ADD an assistant
	operations = append(operations, jsonpatch.Operation{
		Type:  jsonpatch.ADD,
		Path:  "/assistants",
		Value: newAssistantId,
	})

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
	//moduleId := modules[0]["id"].(string)
	/*
		//TEST TOPICS
		checkUnauthorized(t, path)
		resp = getAuthorized(t, userUsername, path)
		require.Equal(t, 200, resp.StatusCode)
		topics := readArrayJsonResult(t, resp)
		require.Equal(t, 2, len(topics))
		require.Equal(t, topics[0]["name"].(string), "Grundlagen der Programmierung mit Java")
		topicId := topics[0]["id"].(string)

		// the topic id should be in the balances gathered in last section
		found := false
		for _, balance := range balances {
			if balance["topic_id"].(string) == topicId {
				found = true
			}
			require.Equal(t, float64(100), balance["amount"])
		}
		require.True(t, found)

		//TEST GET ONE MODULE
		path = "/lecture-service/modules/" + moduleId
		checkUnauthorized(t, path)
		resp = getAuthorized(t, userUsername, path)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		module := readSingleJsonResult(t, resp)
		require.Equal(t, "foo", module["description"].(string))
		exercises := module["exercises"].([]interface{})
		require.NotZero(t, len(exercises))
		exercise := exercises[0].(map[string]interface{})
		exerciseTwo := exercises[1].(map[string]interface{})
		exerciseThree := exercises[2].(map[string]interface{})
		tasks := exercise["tasks"].([]interface{})
		require.NotZero(t, tasks)
		task := tasks[0].(map[string]interface{})
		require.Equal(t, "do something", task["content"].(string))
		hints := task["hints"].([]interface{})
		require.NotZero(t, len(hints))
		hintId := hints[0].(string)
		secondHintId := hints[1].(string)

		//START THE MODULE
		path = "/lecture-service/modules/" + moduleId + "/start"
		PostAuthorizedAndCheckStatusCode(t, userUsername, path, "", http.StatusOK)

		path = "/lecture-service/users/" + userId + "/modules"
		checkUnauthorized(t, path)
		resp = getAuthorized(t, "user1", path)
		moduleHistory := readArrayJsonResult(t, resp)
		require.Equal(t, 1, len(moduleHistory))
		require.Equal(t, moduleHistory[0]["event_type"], "BEGIN")

		// TEST GET HINTS
		path = "/lecture-service/hints/" + hintId
		checkUnauthorized(t, path)
		GetAuthorizedAndCheckStatusCode(t, userUsername, path, 402)
		PostAuthorizedAndCheckStatusCode(t, userUsername, path, "", http.StatusOK)
		PostAuthorizedAndCheckStatusCode(t, userUsername, path, "", http.StatusConflict)

		path = "/lecture-service/hints/" + secondHintId
		GetAuthorizedAndCheckStatusCode(t, userUsername, path, 402)
		PostAuthorizedAndCheckStatusCode(t, userUsername, path, "", 420)

		path = "/lecture-service/hints/" + uuid.NewV4().String()
		PostAuthorizedAndCheckStatusCode(t, userUsername, path, "", http.StatusNotFound)

		// TEST COMPLETE TASK

		path = "/nats-remote/task-backend.task-finished"
		for _, ta := range tasks {
			task = ta.(map[string]interface{})
			result := map[string]interface{}{
				"userId": userId,
				"taskId": task["id"],
			}
			re, err := json.Marshal(result)
			require.Nil(t, err)
			PostAuthorizedAndCheckStatusCode(t, userUsername, path, string(re), 200)
			time.Sleep(500 * time.Millisecond)
		}
		for _, ta := range exerciseTwo["tasks"].([]interface{}) {
			task = ta.(map[string]interface{})
			result := map[string]interface{}{
				"userId": userId,
				"taskId": task["id"],
			}
			re, err := json.Marshal(result)
			require.Nil(t, err)
			PostAuthorizedAndCheckStatusCode(t, userUsername, path, string(re), 200)
			time.Sleep(500 * time.Millisecond)
		}
		for _, ta := range exerciseThree["tasks"].([]interface{}) {
			task = ta.(map[string]interface{})
			result := map[string]interface{}{
				"userId": userId,
				"taskId": task["id"],
			}
			re, err := json.Marshal(result)
			require.Nil(t, err)
			PostAuthorizedAndCheckStatusCode(t, userUsername, path, string(re), 200)
			time.Sleep(500 * time.Millisecond)
		}

		path = "/lecture-service/users/" + userId + "/exercises"
		resp = getAuthorized(t, userUsername, path)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		exerciseHistory := readArrayJsonResult(t, resp)
		require.NotZero(t, exerciseHistory)
		require.Equal(t, "FINISH", exerciseHistory[len(exerciseHistory)-1]["description"])

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
		require.Equal(t, float64(600), balances[0]["amount"])
	*/
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
	jsonBytes, _ := json.Marshal(newModule)
	*operations = append(*operations, jsonpatch.Operation{
		Type:  jsonpatch.ADD,
		Path:  "/modules",
		Value: string(jsonBytes),
	})
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
