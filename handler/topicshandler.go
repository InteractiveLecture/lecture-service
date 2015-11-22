package handler

import (
	"encoding/json"
	"net/http"

	"github.com/InteractiveLecture/id-extractor"
	"github.com/richterrettich/lecture-service/models"
	"github.com/richterrettich/lecture-service/paginator"
	"github.com/richterrettich/lecture-service/repositories"
)

func CollectionHandler(factory repositories.TopicRepositoryFactory) http.Handler {
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

func FindHandler(factory repositories.TopicRepositoryFactory, extractor idextractor.Extractor) http.Handler {
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

func CreateHandler(factory repositories.TopicRepositoryFactory) http.Handler {
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

		id, err = repository.Create(topic)
		if err != nil {
			return http.StatusInternalServerError
		}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Location", id)
		return -1
	}
	return createHandler(handlerFunc)
}

func UpdateHandler(factory repositories.TopicRepositoryFactory, extractor idextractor.Extractor) http.Handler {
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
	}
	return createHandler(handlerFunc)
}

func AddOfficersHandler(factory repositories.TopicRepositoryFactory, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		repository := factory.CreateRepository()
		defer repository.Close()
		id, err := extractor(r)
		if err != nil {
			return err
		}
		var officers = make([]string, 0)
		err = json.NewDecoder(r.Body).Decode(officers)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return nil
		}
		return repository.AddOfficers(id, officers...)
	}
	return createHandler(handlerFunc)
}

func AddAssistantsHandler(factory repositories.TopicRepositoryFactory, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		repository := factory.CreateRepository()
		defer repository.Close()
		id, err := extractor(r)
		if err != nil {
			return err
		}
		var assistants = make([]string, 0)
		err = json.NewDecoder(r.Body).Decode(assistants)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return nil
		}
		return repository.AddAssistants(id, assistants...)
	}
	return createHandler(handlerFunc)
}
func RemoveOfficersHandler(factory repositories.TopicRepositoryFactory, extractor idextractor.Extractor) http.Handler {

	handlerFunc := func(w http.ResponseWriter, r *http.Request) error {
		repository := factory.CreateRepository()
		defer repository.Close()
		defer r.Body.Close()
		id, err := extractor(r)
		if err != nil {
			return err
		}
		var officers = make([]string, 0)
		err = json.NewDecoder(r.Body).Decode(officers)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return nil
		}
		return repository.RemoveOfficers(id, officers...)
	}
	return createHandler(handlerFunc)
}

func RemoveAssistantsHandler(factory repositories.TopicRepositoryFactory, extractor idextractor.Extractor) http.Handler {

	handlerFunc := func(w http.ResponseWriter, r *http.Request) error {
		repository := factory.CreateRepository()
		defer repository.Close()
		defer r.Body.Close()
		id, err := extractor(r)
		if err != nil {
			return err
		}
		var assistants = make([]string, 0)
		err = json.NewDecoder(r.Body).Decode(assistants)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return nil
		}
		return repository.RemoveAssistants(id, assistants...)
	}
	return createHandler(handlerFunc)
}

func GetModulesHandler(factory repositories.ModuleRepositoryFactory, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) error {
		repository := factory.CreateRepository()
		defer repository.Close()
		id, err := extractor(r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return nil
		}
		dr, err := paginator.ParseDepth(r.URL)
		if err != nil {
			return err
		}
		result, err := repository.GetByLectureId(id, dr)
		if err != nil {
			return err
		}
		return json.NewEncoder(w).Encode(result)
	}
	return createHandler(handlerFunc)
}

func CreateModuleHandler(factory repositories.ModuleRepositoryFactory, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) error {
		repository := factory.CreateRepository()
		defer repository.Close()
		m := &models.Module{}
		err := json.NewDecoder(r.Body).Decode(m)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return nil
		}
		err = models.Validate(m)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return nil
		}
		return repository.Create(m)
	}
	return createHandler(handlerFunc)
}
