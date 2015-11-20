package repositories

import (
	"github.com/richterrettich/lecture-service/models"
	"github.com/richterrettich/lecture-service/paginator"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Repository interface {
	GetOne(id string) (*models.Topic, error)
	GetAll(paginator paginator.PageRequest) ([]models.Topic, error)
	Create(topic *models.Topic) error
	Delete(ids ...string) error
	Update(id string, newValues map[string]interface{}) error
	Close()
}

type TopicRepository interface {
	Repository
}

const (
	COL = "topics"
)

var originalSession *mgo.Session

func init() {
	var err error
	originalSession, err = mgo.Dial("mongo") //TODO config file oder sowas für host
	//TODO timeout und retries hinzufügen.
	if err != nil {
		panic(err)
	}
}

type MongoTopicRepository struct {
	*mgo.Session
}

func NewTopicRepository() TopicRepository {
	return &MongoTopicRepository{originalSession.Clone()}
}

func (t *MongoTopicRepository) col() *mgo.Collection {
	return t.DB("lecture").C("topics")
}

func (t *MongoTopicRepository) GetOne(id string) (topic *models.Topic, err error) {
	err = t.col().FindId(id).One(topic)
	return
}

func (t *MongoTopicRepository) Create(topic *models.Topic) error {
	return t.col().Insert(topic)
}

func (t *MongoTopicRepository) GetAll(page paginator.PageRequest) ([]models.Topic, error) {
	var result = make([]models.Topic, 0)
	return result, ApplyPagination(t.col().Find(nil), page, result)
}

func (t *MongoTopicRepository) Delete(ids ...string) error {
	mongoIds := make([]bson.M, len(ids))
	for i, v := range ids {
		mongoIds[i] = bson.M{"_id": v}
	}
	return t.col().Remove(mongoIds)
}

func (t *MongoTopicRepository) Update(id string, newValues map[string]interface{}) error {
	return t.col().Update(bson.M{"_id": id}, newValues)
}

func ApplyPagination(query *mgo.Query, page paginator.PageRequest, result interface{}) error {
	query = query.Skip(page.Number * page.Size).Limit(page.Size)
	if len(page.Sorts) > 0 {
		query = query.Sort(sortsToString(page.Sorts)...)
	}
	return query.All(result)
}

func sortsToString(sorts []paginator.Sort) []string {
	var result = make([]string, 0)
	for _, v := range sorts {
		prefix := ""
		if v.Direction == paginator.DESC {
			prefix = "-"
		}
		result = append(result, prefix+v.Name)
	}
	return result
}
