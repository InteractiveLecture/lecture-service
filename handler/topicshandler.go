package handler

import (
	"encoding/json"
	"net/http"

	"github.com/InteractiveLecture/id-extractor"
	"github.com/richterrettich/lecture-service/models"
	"github.com/richterrettich/lecture-service/paginator"
	"github.com/richterrettich/lecture-service/repositories"
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
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			return http.StatusInternalServerError
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
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			return http.StatusInternalServerError
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

func TopicUpdateHandler(factory repositories.TopicRepositoryFactory, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		repository := factory.CreateRepository()
		defer repository.Close()
		id, err := extractor(r)
		if err != nil {
			return http.StatusInternalServerError
		}
		var newValues = make(map[string]interface{})
		//TODO validate new values.
		err = json.NewDecoder(r.Body).Decode(newValues)
		if err != nil {
			return http.StatusBadRequest
		}
		err = repository.Update(id, newValues)
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
		var officers = make([]string, 0)
		err = json.NewDecoder(r.Body).Decode(officers)
		if err != nil {
			return http.StatusBadRequest
		}
		err = repository.AddOfficers(id, officers...)
		if err != nil {
			return http.StatusInternalServerError
		}
		return -1
	}
	return createHandler(handlerFunc)
}

func TopicAddAssistantsHandler(factory repositories.TopicRepositoryFactory, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		repository := factory.CreateRepository()
		defer repository.Close()
		id, err := extractor(r)
		if err != nil {
			return http.StatusInternalServerError
		}
		var assistants = make([]string, 0)
		err = json.NewDecoder(r.Body).Decode(assistants)
		if err != nil {
			return http.StatusBadRequest
		}
		err = repository.AddAssistants(id, assistants...)
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
		var officers = make([]string, 0)
		err = json.NewDecoder(r.Body).Decode(officers)
		if err != nil {
			return http.StatusBadRequest
		}
		err = repository.RemoveOfficers(id, officers...)
		if err != nil {
			return http.StatusInternalServerError
		}
		return -1
	}
	return createHandler(handlerFunc)
}

func TopicRemoveAssistantsHandler(factory repositories.TopicRepositoryFactory, extractor idextractor.Extractor) http.Handler {

	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		repository := factory.CreateRepository()
		defer repository.Close()
		id, err := extractor(r)
		if err != nil {
			return http.StatusInternalServerError
		}
		var assistants = make([]string, 0)
		err = json.NewDecoder(r.Body).Decode(assistants)
		if err != nil {
			return http.StatusBadRequest
		}
		err = repository.RemoveAssistants(id, assistants...)
		if err != nil {
			return http.StatusInternalServerError
		}
		return -1
	}
	return createHandler(handlerFunc)
}
