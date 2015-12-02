package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/InteractiveLecture/id-extractor"
	"github.com/richterrettich/lecture-service/datamapper"
	"github.com/richterrettich/lecture-service/paginator"
)

func HintHistoryHandler(mapper *datamapper.DataMapper, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		pageRequest, err := paginator.ParsePages(r.URL)
		if err != nil {
			return http.StatusInternalServerError
		}
		id, err := extractor(r)
		if err != nil {
			return http.StatusInternalServerError
		}
		exerciseId := ""
		ids, ok := r.URL.Query()["exercise_id"]
		if ok {
			exerciseId = ids[0]
		}
		result, err := mapper.GetHintHistory(id, pageRequest, exerciseId)
		if err != nil {
			log.Println(err)
			return http.StatusNotFound
		}
		reader := bytes.NewReader(result)
		_, err = io.Copy(w, reader)
		if err != nil {
			log.Println(err)
		}
		return -1
	}
	return createHandler(handlerFunc)
}

func ModuleHistoryHandler(mapper *datamapper.DataMapper, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		pageRequest, err := paginator.ParsePages(r.URL)
		if err != nil {
			return http.StatusInternalServerError
		}
		id, err := extractor(r)
		if err != nil {
			return http.StatusInternalServerError
		}
		topicId := ""
		ids, ok := r.URL.Query()["topic_id"]
		if ok {
			topicId = ids[0]
		}
		result, err := mapper.GetModuleHistory(id, pageRequest, topicId)
		if err != nil {
			return http.StatusNotFound
		}
		reader := bytes.NewReader(result)
		_, err = io.Copy(w, reader)
		if err != nil {
			log.Println(err)
		}
		return -1
	}
	return createHandler(handlerFunc)
}

func ExerciseHistoryHandler(mapper *datamapper.DataMapper, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		pageRequest, err := paginator.ParsePages(r.URL)
		if err != nil {
			return http.StatusInternalServerError
		}
		id, err := extractor(r)
		if err != nil {
			return http.StatusInternalServerError
		}
		moduleId := ""
		ids, ok := r.URL.Query()["module_id"]
		if ok {
			moduleId = ids[0]
		}
		result, err := mapper.GetExerciseHistory(id, pageRequest, moduleId)
		if err != nil {
			return http.StatusNotFound
		}
		reader := bytes.NewReader(result)
		_, err = io.Copy(w, reader)
		if err != nil {
			log.Println(err)
		}
		return -1
	}
	return createHandler(handlerFunc)
}

func NextModulesForUserHandler(mapper *datamapper.DataMapper, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		id, err := extractor(r)
		if err != nil {
			return http.StatusInternalServerError
		}
		result, err := mapper.GetNextModulesForUser(id)
		if err != nil {
			return http.StatusNotFound
		}
		reader := bytes.NewReader(result)
		_, err = io.Copy(w, reader)
		if err != nil {
			log.Println(err)
		}
		return -1
	}
	return createHandler(handlerFunc)
}

func TopicBalanceHandler(mapper *datamapper.DataMapper, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		id, err := extractor(r)
		if err != nil {
			return http.StatusInternalServerError
		}
		result, err := mapper.GetTopicBalances(id)
		if err != nil {
			return http.StatusNotFound
		}
		reader := bytes.NewReader(result)
		_, err = io.Copy(w, reader)
		if err != nil {
			log.Println(err)
		}
		return -1
	}
	return createHandler(handlerFunc)
}

func ModuleStartHandler(mapper *datamapper.DataMapper, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		id, err := extractor(r)
		if err != nil {
			return http.StatusInternalServerError
		}
		var moduleId string
		err = json.NewDecoder(r.Body).Decode(moduleId)
		if err != nil {
			return http.StatusBadRequest
		}
		err = mapper.StartModule(id, moduleId)
		if err != nil {
			return http.StatusInternalServerError
		}
		return -1
	}
	return createHandler(handlerFunc)
}

func ExerciseStartHandler(mapper *datamapper.DataMapper, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		id, err := extractor(r)
		if err != nil {
			return http.StatusInternalServerError
		}
		var exerciseId string
		err = json.NewDecoder(r.Body).Decode(exerciseId)
		if err != nil {
			return http.StatusBadRequest
		}
		err = mapper.StartExercise(id, exerciseId)
		if err != nil {
			return http.StatusNotFound
		}
		return -1
	}
	return createHandler(handlerFunc)
}
