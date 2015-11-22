package repositories

import (
	"github.com/richterrettich/lecture-service/models"
	"github.com/richterrettich/lecture-service/paginator"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoRepositoryFactory struct {
	originalSession *mgo.Session
}

func (f *MongoRepositoryFactory) CreateRepository() TopicRepository {
	return &MongoTopicRepository{f.originalSession.Clone()}
}

func (r *MongoTopicRepository) topicToBson(t *models.Topic) *bson.M {
	return &bson.M{
		"_id":         t.ID,
		"name":        t.Name,
		"description": t.Description,
		"officers":    t.Officers,
		"assistants":  t.Assistants,
	}
}

func (r *MongoTopicRepository) MapToTopic(values map[string]interface{}) *models.Topic {
	return &models.Topic{
		ID:          values["_id"].(string),
		Name:        values["name"].(string),
		Description: values["description"].(string),
		Officers:    values["officers"].([]string),
		Assistants:  values["assistants"].([]string),
	}
}

type MongoTopicRepository struct {
	*mgo.Session
}

func (t *MongoTopicRepository) col() *mgo.Collection {
	return t.DB("lecture").C("topics")
}

func (r *MongoTopicRepository) AddOfficers(id string, officers ...string) error {
	return r.addSlice(id, "officers", officers)
}

func (r *MongoTopicRepository) addSlice(id, arrayName string, slice interface{}) error {
	return r.col().Update(bson.M{"_id": id}, bson.M{
		"$push": bson.M{
			arrayName: bson.M{
				"$each": slice,
			},
		},
	})
}

func (r *MongoTopicRepository) RemoveAssistants(id string, assistants ...string) error {
	return r.removeSlice(id, "assistants", assistants)
}

func (r *MongoTopicRepository) RemoveOfficers(id string, officers ...string) error {
	return r.removeSlice(id, "officers", officers)
}

func (r *MongoTopicRepository) removeSlice(id, sliceName string, data interface{}) error {
	return r.col().Update(bson.M{"_id": id}, bson.M{
		"$pull": bson.M{
			sliceName: bson.M{
				"$in": data,
			},
		},
	})
}

func (r *MongoTopicRepository) AddAssistants(id string, assistants ...string) error {
	return r.addSlice(id, "assistants", assistants)
}

func (t *MongoTopicRepository) GetOne(id string) (topic *models.Topic, err error) {
	m := bson.M{}
	err = t.col().FindId(id).One(m)
	topic = t.MapToTopic(m)
	return
}

func (t *MongoTopicRepository) Create(topic *models.Topic) (string, error) {
	id := bson.NewObjectId().Hex()
	topic.ID = id
	err := t.col().Insert(t.topicToBson(topic))
	if err != nil {
		return "", err
	}
	return id, nil
}

func (t *MongoTopicRepository) GetAll(page paginator.PageRequest) ([]*models.Topic, error) {
	var bsons = make([]bson.M, 0)
	err := ApplyPagination(t.col().Find(nil), page, bsons)
	if err != nil {
		return nil, err
	}

	var result = make([]*models.Topic, len(bsons))
	for i, v := range bsons {
		result[i] = t.MapToTopic(v)
	}
	return result, err
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
