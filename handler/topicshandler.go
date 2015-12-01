package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/InteractiveLecture/id-extractor"
	"github.com/richterrettich/lecture-service/lecturepatch"
	"github.com/richterrettich/lecture-service/models"
	"github.com/richterrettich/lecture-service/paginator"
)

func TopicCollectionHandler(factory repositories.TopicRepositoryFactory) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		repository := factory.CreateRepository()
		defer repository.Close()
		pageRequest, err := paginator.ParsePages(r.URL)
		if err != nil {
			return http.StatusInternalServerError
		}
		result, err := repository.GetAll(pageRequest)
		if err != nil {
			return http.StatusInternalServerError
		}
		reader := bytes.NewReader(result)
		_, err := io.Copy(w, reader)
		if err != nil {
			log.Println(err)
		}
		return -1
	}
	return createHandler(handlerFunc)
}

func TopicFindHandler(factory repositories.TopicRepositoryFactory, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		repository := factory.CreateRepository()
		defer repository.Close()
		id, err := extractor(r)
		if err != nil {
			return http.StatusInternalServerError
		}
		result, err := repository.GetOne(id)
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

func TopicCreateHandler(factory repositories.TopicRepositoryFactory) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		repository := factory.CreateRepository()
		defer repository.Close()
		var topic *models.Topic
		err := json.NewDecoder(r.Body).Decode(topic)
		if err != nil {
			return http.StatusBadRequest
		}
		err = models.Validate(topic)
		if err != nil {
			return http.StatusBadRequest
		}

		id, err := repository.Create(topic)
		if err != nil {
			return http.StatusInternalServerError
		}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Location", id)
		return -1
	}
	return createHandler(handlerFunc)
}

func TopicPatchHandler(factory repositories.TopicRepositoryFactory, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		repository := factory.CreateRepository()
		defer repository.Close()
		id, err := extractor(r)
		if err != nil {
			return http.StatusInternalServerError
		}
		patch, err := lecturepatch.Decode(r.Body)
		if err != nil {
			return http.StatusBadRequest
		}
		err = repository.Update(id, patch)
		if err != nil {
			return http.StatusInternalServerError
		}
		return -1
	}
	return createHandler(handlerFunc)
}

func TopicAddOfficersHandler(factory repositories.TopicRepositoryFactory, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		repository := factory.CreateRepository()
		defer repository.Close()
		id, err := extractor(r)
		if err != nil {
			return http.StatusInternalServerError
		}
		var officer string
		err = json.NewDecoder(r.Body).Decode(officer)
		if err != nil {
			return http.StatusBadRequest
		}
		err = repository.AddOfficer(id, officer)
		if err != nil {
			return http.StatusInternalServerError
		}
		return -1
	}
	return createHandler(handlerFunc)
}

func TopicRemoveOfficersHandler(factory repositories.TopicRepositoryFactory, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		repository := factory.CreateRepository()
		defer repository.Close()
		id, err := extractor(r)
		if err != nil {
			return http.StatusInternalServerError
		}
		var officer string
		err = json.NewDecoder(r.Body).Decode(officer)
		if err != nil {
			return http.StatusBadRequest
		}
		err = repository.RemoveOfficers(id, officer)
		if err != nil {
			return http.StatusInternalServerError
		}
		return -1
	}
	return createHandler(handlerFunc)
}
