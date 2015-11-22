package handler

import (
	"encoding/json"
	"net/http"

	"github.com/InteractiveLecture/id-extractor"
	"github.com/richterrettich/lecture-service/models"
	"github.com/richterrettich/lecture-service/paginator"
	"github.com/richterrettich/lecture-service/repositories"
)

func ModulesTreeHandler(factory repositories.ModuleRepositoryFactory, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		repository := factory.CreateRepository()
		defer repository.Close()
		id, err := extractor(r)
		if err != nil {
			return http.StatusBadRequest
		}
		dr, err := paginator.ParseDepth(r.URL)
		if err != nil {
			return http.StatusInternalServerError
		}
		result, err := repository.GetByLectureId(id, dr)
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

func ModulesDeleteHandler(factory repositories.ModuleRepositoryFactory, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		repository := factory.CreateRepository()
		defer repository.Close()
		id := extractor(r)
		var children = make([]models.Module, 0)
		return -1
	}
	return createHandler(handlerFunc)
}

func ModulesGetHandler(factory repositories.ModuleRepositoryFactory, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		repo := factory.CreateRepository()
		defer repo.Close()
		id, err := extractor(r)
		if err != nil {
			return http.StatusInternalServerError
		}
		module, err := repo.GetOne(id)
		if err != nil {
			return http.StatusNotFound
		}
		err = json.NewEncoder(w).Encode(module)
		if err != nil {
			return http.StatusInternalServerError
		}
		return -1
	}
	return createHandler(handlerFunc)
}

func ModulesCreateHandler(factory repositories.ModuleRepositoryFactory, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		repository := factory.CreateRepository()
		defer repository.Close()
		m := &models.Module{}
		err := json.NewDecoder(r.Body).Decode(m)
		if err != nil {
			return http.StatusBadRequest
		}
		err = models.Validate(m)
		if err != nil {
			return http.StatusBadRequest
		}
		err = repository.Create(m)
		if err != nil {
			return http.StatusInternalServerError
		}
		return -1
	}
	return createHandler(handlerFunc)
}
