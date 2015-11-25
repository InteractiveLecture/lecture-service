package repositories

import (
	"database/sql"
	"strings"

	_ "github.com/lib/pq"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/richterrettich/lecture-service/models"
	"github.com/richterrettich/lecture-service/paginator"
)

type PTRepoFactory struct {
	db *sql.DB
}

func (f *PTRepoFactory) CreateRepository() TopicRepository {
	return PTRepo{f.db}
}

type PTRepo struct {
	session *sql.DB
}

func (r *PTRepo) AddAuthority(id string, officer, kind string) error {
	_, err = r.session.Exec(`INSERT INTO topic_authority (topic_id,user_id,kind), values($1,$2,$3)`, id, officer, strings.ToUpper(kind))
	if err != nil {
		return err
	}
	return nil
}

func (r *PTRepo) RemoveAuthority(id string, user string) error {
	_, err := r.session.Exec(`DELTE FROM topic_authority where user_id = $1 AND topic_id = $2`, user, id)
	if err != nil {
		return err
	}
	return nil
}

func (t *PTRepo) rowToTopic(row *sql.Row) *models.Topic {
	result := models.Topic{}
	var (
		id, name, description string
	)

}

func (t *PTRepo) GetOne(id string) (topic *models.Topic, err error) {
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
