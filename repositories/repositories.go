package repositories

import (
	"github.com/richterrettich/lecture-service/models"
	"github.com/richterrettich/lecture-service/paginator"
)

type TopicRepository interface {
	GetOne(id string) (*models.Topic, error)
	GetAll(paginator paginator.PageRequest) ([]*models.Topic, error)
	Create(topic *models.Topic) (id string, err error)
	Delete(ids ...string) error
	Update(id string, newValues map[string]interface{}) error
	Close()
	AddOfficers(id string, officers ...string) error
	AddAssistants(id string, assistants ...string) error
	RemoveOfficers(id string, officers ...string) error
	RemoveAssistants(id string, assistants ...string) error
}

type TopicRepositoryFactory interface {
	CreateRepository() TopicRepository
}

type ModuleRepository interface {
	GetOne(id string) (*models.Module, error)
	GetByLectureId(lectureId string, dr paginator.DepthRequest) ([]*models.Module, error)
	Create(module *models.Module) error
	GetChildren(id string) ([]*models.Module, error)
	Close()
}

type ModuleRepositoryFactory interface {
	CreateRepository() ModuleRepository
}
