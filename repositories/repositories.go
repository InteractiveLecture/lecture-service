package repositories

import (
	"github.com/richterrettich/lecture-service/models"
	"github.com/richterrettich/lecture-service/paginator"
)

type TopicRepository interface {
	GetOne(id string) ([]byte, error)
	GetAll(paginator paginator.PageRequest) ([]byte, error)
	Create(*models.Topic) (id string, err error)
	Delete(id string) error
	Update(id string, newValues map[string]interface{}) error
	AddAuthority(id, authority, kind string) error
	RemoveAuthority(id string, authority string) error
}

type TopicRepositoryFactory interface {
	CreateRepository() TopicRepository
}

type ModuleRepository interface {
	GetOne(id string) ([]byte, error)
	GetByLectureId(lectureId string, dr paginator.DepthRequest) ([]byte, error)
	Create(*models.Module) error
	GetChildren(id string) ([]byte, error)
}

type ModuleRepositoryFactory interface {
	CreateRepository() ModuleRepository
}
