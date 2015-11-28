package datamapper

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"

	"github.com/richterrettich/lecture-service/models"
	"github.com/richterrettich/lecture-service/paginator"
)

type PTRepoFactory struct {
	db *sql.DB
}

func (f *PTRepoFactory) CreateRepository() TopicRepository {
	return &PTRepo{f.db}
}

type PTRepo struct {
	session *sql.DB
}

func (r *PTRepo) AddAuthority(id, officer, kind string) error {
	_, err := r.session.Exec(`INSERT INTO topic_authority (topic_id,user_id,kind), values($1,$2,$3)`, id, officer, strings.ToUpper(kind))
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

func (t *PTRepo) queryIntoBytes(query string, params ...interface{}) ([]byte, error) {
	row, err := t.session.Query(query, params)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	var result = make([]byte, 0)
	for row.Next() {
		err = row.Scan(result)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (t *PTRepo) GetOne(id string) ([]byte, error) {
	return t.queryIntoBytes(`SELECT get_topic($1)`, id)
}

func (t *PTRepo) Create(topic *models.Topic) (string, error) {
	id := uuid.NewV4().String()
	_, err := t.session.Exec(`INSERT INTO topics(id,name,description,version) values($1,$2,$3,$4)`, id, topic.Name, topic.Description, 1)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (t *PTRepo) GetAll(page paginator.PageRequest) ([]byte, error) {
	return t.queryIntoBytes(`SELECT * from query_topics($1,$2)`, page.Number*page.Size, page.Size)
}

func (t *PTRepo) Delete(id string) (err error) {
	_, err = t.session.Exec(`DELETE FROM topics where id = $1`, id)
	return
}

func (t *PTRepo) Update(id string, newValues map[string]interface{}) error {
	stmt := "UPDATE topics set "
	var parameters = make([]interface{}, len(newValues))
	i := 1
	for k, v := range newValues {
		stmt = fmt.Sprintf("%s %s = $%d,", stmt, k, i)
		parameters = append(parameters, v)
		i = i + 1
	}
	stmt = strings.Trim(stmt, ",")
	stmt = fmt.Sprintf("%s where id = $%d", stmt, i)
	_, err := t.session.Exec(stmt, parameters...)
	return err
}
