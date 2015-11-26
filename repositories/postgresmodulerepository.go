package repositories

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/richterrettich/lecture-service/models"
	"github.com/richterrettich/lecture-service/modulepatch"
	"github.com/richterrettich/lecture-service/paginator"
	"github.com/satori/go.uuid"
)

type PModuleRepoFactory struct {
	DB *sql.DB
}

type PModuleRepo struct {
	session *sql.DB
}

func (f *PModuleRepoFactory) CreateRepository() ModuleRepository {
	return &PModuleRepo{f.DB}
}

type InvalidPatchError struct {
	Message string
}

func (e *InvalidPatchError) Error() string {
	return e.Message
}

func extractParts(patch *modulepatch.Operation) (string, []string, error) {
	parts := strings.Split(patch.Path, "/")
	if len(parts) == 0 || parts[0] != "" {
		return "", nil, &InvalidPatchError{fmt.Sprintf("Path %s is invalid.", patch.Path)}
	}
	id, parts := parts[1], parts[2:]
	return id, parts, nil
}

/*
func translateOperation(patchOperation *modulepatch.Operation) (*txn.Op, error) {
	result := &txn.Op{}
	id, parts, err := extractParts(patchOperation)
	if err != nil {
		return nil, err
	}
	result.Id = id
	switch patchOperation.Type {
	}
	return result, nil
}*/

func (r *PModuleRepo) GetChildren(id string) ([]byte, error) {
	return nil, nil
}

func prepare(stmt string, values ...interface{}) (string, []interface{}) {
	parametersString := ""
	var parameters = make([]interface{}, 0)
	currentIndex := 1
	for _, v := range values {
		val := reflect.ValueOf(v)
		if val.Kind() == reflect.Slice {
			for i := 0; i < val.Len(); i++ {
				inval := val.Index(i)
				parameters = append(parameters, inval)
				parametersString = fmt.Sprintf("%s,$%d", parametersString, currentIndex)
				currentIndex = currentIndex + 1
			}
		} else {
			parameters = append(parameters, v)
			parametersString = fmt.Sprintf("%s,$%d", parametersString, currentIndex)
			currentIndex = currentIndex + 1
		}
	}
	stmt = fmt.Sprintf(stmt, strings.Trim(parametersString, ","))
	return stmt, parameters
}

func (r *PModuleRepo) Create(m *models.Module) error {
	if m.ID == "" {
		m.ID = uuid.NewV4().String()
	}
	_, err := r.session.Exec(prepare("SELECT insert_module(%s)", m.ID, m.TopicID, m.Description, m.VideoID, m.ScriptID, m.Parents))
	return err
}

func rowToBytes(row *sql.Row) ([]byte, error) {
	var result = make([]byte, 0)
	err := row.Scan(result)
	return result, err
}

func (r *PModuleRepo) GetOne(id string) ([]byte, error) {
	return rowToBytes(r.session.QueryRow(`SELECT get_module($1)`, id))
}

func (r *PModuleRepo) GetByLectureId(topicId string, dr paginator.DepthRequest) ([]byte, error) {
	return rowToBytes(r.session.QueryRow(`SELECT get_module_tree($1,$2,$3)`, topicId, dr.Descendants, dr.Ancestors))
}
