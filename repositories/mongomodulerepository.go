package repositories

import (
	"github.com/richterrettich/lecture-service/models"
	"github.com/richterrettich/lecture-service/paginator"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

func (r *MModuleRepo) Create(m *models.Module) error {
	lastModule := &models.Module{}
	err := r.col().Find(bson.M{"topic_id": m.TopicID}).Sort("-depth").Limit(1).One(lastModule)
	if err != nil {
		// if this error occurs, this must be the root module.
		if err == mgo.ErrNotFound {
			m.Depth = 0
			insertError := r.col().Insert(m)
			if insertError != nil {
				return insertError
			}
		} else {
			//something went wrong with the database
			return err
		}
	}
	//new last module
	if lastModule.Depth+1 <= m.Depth {
		m.Depth = lastModule.Depth + 1
	}
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
