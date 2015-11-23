package repositories

import (
	"fmt"
	"strings"

	"github.com/richterrettich/lecture-service/models"
	"github.com/richterrettich/lecture-service/modulepatch"
	"github.com/richterrettich/lecture-service/paginator"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/txn"
)

type MModuleRepoFactory struct {
	originalSession *mgo.Session
}

type MModuleRepo struct {
	*mgo.Session
}

func (f *MModuleRepoFactory) CreateRepository() ModuleRepository {
	return &MModuleRepo{f.originalSession.Clone()}
}

func (r *MModuleRepo) col() *mgo.Collection {
	return r.DB("lecture").C("modules")
}

func (r *MModuleRepo) ApplyPatch(topicId string, p *modulepatch.Patch) error {
	return nil
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
		return "", &InvalidPatchError{fmt.Sprintf("Path %s is invalid.", patch.Path)}
	}
	id, parts := parts[1], a[2:]
	return id, parts, nil
}

func translateOperation(patchOperation *modulepatch.Operation) (*txn.Op, error) {
	result := &txn.Op{}
	id, parts, err := extractId(patchOperation)
	if err != nil {
		return nil, err
	}
	result.Id = id
	switch patchOperation.Type {
	case modulepatch.ADD:
		switch parts[0] {
		cae "exercises":
			return nil, nil
		case "recommendations":
			return nil, nil
		case "parents":
			return nil, nil
		default:
			return nil, &InvalidPatchError{fmt.Sprintf("cannot add to %s", parts[0])}
		}
		result.Insert = bson.M{}
	}
	return result, nil
}

func (r *MModuleRepo) GetChildren(id string) ([]*models.Module, error) {
	return nil, nil
}

func (r *MModuleRepo) Create(m *models.Module) error {
	return r.col().Insert(m)
}

func (r *MModuleRepo) GetOne(id string) (*models.Module, error) {
	return nil, nil
}

func (r *MModuleRepo) GetByLectureId(topicId string, dr paginator.DepthRequest) ([]*models.Module, error) {
	var result = make([]*models.Module, 0)
	// return the entire tree.
	if dr.Layer == 0 && dr.Ancestors == 0 && dr.Descendants == -1 {
		return result, r.col().Find(bson.M{"topic_id": topicId}).All(result)
	}

	// start with nodes Layer - Ancestors, return all children.
	if dr.Descendants == -1 {
		return result, r.col().Find(bson.M{
			"$and": bson.M{
				"topic_id": topicId, //TODO inspect if objectid needs to be casted
				"depth": bson.M{
					"$gte": dr.Layer - dr.Ancestors,
				},
			},
		}).All(result)
	}
	// end with layer + descendants. Return all layers above.
	if dr.Ancestors == -1 {
		return result, r.col().Find(bson.M{
			"$and": bson.M{
				"topic_id": topicId, //TODO inspect if objectid needs to be casted
				"depth": bson.M{
					"$lte": dr.Layer + dr.Descendants,
				},
			},
		}).All(result)
	}
	// return a window between layer + descendants and layer - ancestors.
	return result, r.col().Find(bson.M{
		"$and": bson.M{
			"topic_id": topicId, //TODO inspect if objectid needs to be casted
			"depth": bson.M{
				"$gte": dr.Layer - dr.Ancestors,
				"$lte": dr.Layer + dr.Descendants,
			},
		},
	}).All(result)
}
