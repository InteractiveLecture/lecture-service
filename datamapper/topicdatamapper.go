package datamapper

import (
	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"

	"github.com/richterrettich/lecture-service/paginator"
)

func (r *DataMapper) AddOfficer(id, officer) error {
	_, err := r.db.Exec(`SELECT add_officer($1,$2)`, id, officer)
	if err != nil {
		return err
	}
	return nil
}

func (r *DataMapper) RemoveOfficer(id string, user string) error {
	_, err := r.db.Exec(`SELECT remove_officer($1,$2)`, user, id)
	if err != nil {
		return err
	}
	return nil
}
func (t *DataMapper) GetOneTopic(id string) ([]byte, error) {
	return t.queryIntoBytes(`SELECT get_topic($1)`, id)
}

func (t *DataMapper) CreateTopic(topic map[string]interface{}) (string, error) {
	ok, id := topic["id"]
	if !ok {
		id := uuid.NewV4().String()
	}
	_, err := t.db.Exec(`INSERT INTO topics(id,name,description,version) values($1,$2,$3,$4)`, id, topic["name"], topic["description"], 1)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (t *DataMapper) GetAll(page paginator.PageRequest) ([]byte, error) {
	return t.queryIntoBytes(`SELECT * from query_topics($1,$2)`, page.Number*page.Size, page.Size)
}

func (t *DataMapper) Delete(id string) (err error) {
	_, err = t.db.Exec(`DELETE FROM topics where id = $1`, id)
	return
}
