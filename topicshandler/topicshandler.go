package topicshandler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/richterrettich/lecture-service/models"
	"github.com/richterrettich/lecture-service/paginator"
	"github.com/richterrettich/lecture-service/repositories"
)

type errorHandler func(http.ResponseWriter, *http.Request) error

func (handler errorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := handler(w, r); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func NewCollectionHandler(reposiroty repositories.TopicRepository) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) error {
		pageRequest, err := paginator.ParsePages(r.URL)
		if err != nil {
			return err
		}
		defer reposiroty.Close()
		result, err := repo.GetAll(pageRequest)
		if err != nil {
			return err
		}
		return json.NewEncoder(w).Encode(result)
	}
	return createHandler(handlerFunc)
}

func NewFindHandler(repository repositories.TopicRepository) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) error {
		vars := mux.Vars(r)
		id := vars["id"]
		defer repository.Close()
		result, err := repo.GetOne(id)
		if err != nil {
			return err
		}
		return json.NewEncoder(w).Encode(result)
	}
	return createHandler(handlerFunc)
}

func NewCreateHandler(repository repositories.TopicRepository) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) error {
		defer reposiroty.Close()
		var topic models.Topic
		err := json.NewDecoder(r.Body).Decode(topic)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return nil
		}
		err = topic.Validate()
		if err != nil {
			http.Error(w, "validation failed", http.StatusBadRequest)
			return nil
		}
		return repository.Create(topic)
	}
	return createHandler(handlerFunc)
}

func NewUpdateHandler(repository repositories.TopicRepository) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) error {
		defer repository.Close()
		id := mux.Vars(r)["id"]
		var newValues = make(map[string]interface{})
		err := json.NewDecoder(r.Body).Decode(newValues)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return nil
		}
		return repository.Update(id, newValues)
	}
	return createHandler(handlerFunc)
}

func createHandler(handlerFunc errorHandler) http.Handler {
	return http.Handler(errorHandler(handlerFunc))
}
