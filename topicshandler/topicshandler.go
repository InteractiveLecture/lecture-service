package topicshandler

import (
	"encoding/json"
	"net/http"

	"github.com/InteractiveLecture/id-extractor"
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

func CollectionHandler(factory repositories.TopicRepositoryFactory) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) error {
		repository := factory.CreateRepository()
		defer repository.Close()
		pageRequest, err := paginator.ParsePages(r.URL)
		if err != nil {
			return err
		}
		result, err := repository.GetAll(pageRequest)
		if err != nil {
			return err
		}
		return json.NewEncoder(w).Encode(result)
	}
	return createHandler(handlerFunc)
}

func FindHandler(factory repositories.TopicRepositoryFactory, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) error {
		repository := factory.CreateRepository()
		defer repository.Close()
		id, err := extractor(r)
		if err != nil {
			return err
		}
		result, err := repository.GetOne(id)
		if err != nil {
			return err
		}
		return json.NewEncoder(w).Encode(result)
	}
	return createHandler(handlerFunc)
}

func CreateHandler(factory repositories.TopicRepositoryFactory) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) error {
		repository := factory.CreateRepository()
		defer repository.Close()
		var topic *models.Topic
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

func UpdateHandler(factory repositories.TopicRepositoryFactory, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) error {
		repository := factory.CreateRepository()
		defer repository.Close()
		id, err := extractor(r)
		if err != nil {
			return err
		}
		var newValues = make(map[string]interface{})
		err = json.NewDecoder(r.Body).Decode(newValues)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return nil
		}
		return repository.Update(id, newValues)
	}
	return createHandler(handlerFunc)
}

func AddOfficersHandler(factory repositories.TopicRepositoryFactory, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) error {
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
	handlerFunc := func(w http.ResponseWriter, r *http.Request) error {
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
		id, err := extracotr(r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return nil
		}
		result, err := repository.GetByLectureId(id)
		if err != nil {
			return err
		}
		return json.NewEncoder(w).Encode(result)
	}
	return createHandler(handlerFunc)
}

func CreateModuleHandler(factory repositories.ModuleRepositoryFactory, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		repository := factory.CreateRepository()
		repository.Close()

	}
}

func createHandler(handlerFunc errorHandler) http.Handler {
	return http.Handler(errorHandler(handlerFunc))
}
