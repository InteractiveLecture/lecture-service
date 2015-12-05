package datamapper

import (
	"log"

	_ "github.com/lib/pq"

	"github.com/richterrettich/lecture-service/paginator"
)

func (r *DataMapper) AddOfficer(id, officer string) error {
	_, err := r.db.Exec(`SELECT add_officer($1,$2)`, id, officer)
	if err != nil {
		return err
	}
	return nil
}

func (r *DataMapper) RemoveOfficer(id, user string) error {
	_, err := r.db.Exec(`SELECT remove_officer($1,$2)`, id, user)
	if err != nil {
		return err
	}
	return nil
}

func (t *DataMapper) GetOneTopic(id string) ([]byte, error) {
	return t.queryIntoBytes(`SELECT get_topic($1)`, id)
}

func (t *DataMapper) CreateTopic(topic map[string]interface{}) error {
	stmt, parameters := prepare("SELECT add_topic(%v)", topic["id"], topic["name"], topic["description"], topic["officers"])
	_, err := t.db.Exec(stmt, parameters...)
	return err
}

func (t *DataMapper) GetTopicsPage(page paginator.PageRequest) ([]byte, error) {
	log.Println(page)
	result, err := t.queryIntoBytes(`SELECT * from query_topics($1,$2)`, page.Number*page.Size, page.Size)

	if err != nil {
		log.Println(err)
	}
	return result, err
}

func (t *DataMapper) Delete(id string) (err error) {
	_, err = t.db.Exec(`DELETE FROM topics where id = $1`, id)
	return
}
