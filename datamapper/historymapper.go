package datamapper

import (
	"time"

	"github.com/richterrettich/lecture-service/paginator"
)

func (r *DataMapper) GetHintHistory(userId string, pr paginator.PageRequest, exerciseId string) ([]byte, error) {
	if exerciseId != "" {
		return r.queryIntoBytes("SELECT get_hint_purchase_history($1,$2,$3,$4)", userId, pr.Size, pr.Size*pr.Number, exerciseId)
	}
	return r.queryIntoBytes(`SELECT get_hint_purchase_history($1,$2,$3)`, userId, pr.Size, pr.Size*pr.Number)
}

func (r *DataMapper) GetModuleHistory(userId string, pr paginator.PageRequest, topicId string) ([]byte, error) {
	if topicId != "" {
		return r.queryIntoBytes(`SELECT get_module_history($1,$2,$3,$4)`, userId, pr.Size, pr.Size*pr.Number, topicId)
	}
	return r.queryIntoBytes(`SELECT get_module_history($1,$2,$3)`, userId, pr.Size, pr.Size*pr.Number)
}

func (r *DataMapper) GetExerciseHistory(userId string, pr paginator.PageRequest, moduleId string) ([]byte, error) {
	if moduleId != "" {
		return r.queryIntoBytes(`SELECT get_exercise_history($1,$2,$3,$4)`, userId, pr.Size, pr.Size*pr.Number, moduleId)
	}
	return r.queryIntoBytes(`SELECT get_exercise_history($1,$2,$3)`, userId, pr.Size, pr.Size*pr.Number)
}

func (r *DataMapper) GetNextModulesForUser(id string) ([]byte, error) {
	return r.queryIntoBytes("SELECT get_next_modules_for_user($1)", id)
}

func (r *DataMapper) GetTopicBalances(id string) ([]byte, error) {
	return r.queryIntoBytes("Select get_balances($1)", id)
}

func (r *DataMapper) StartExercise(id, exerciseId string) error {
	_, err := r.db.Exec("insert into exercise_progress_histories(user_id,exercise_id,amount,time,state) values($1,$2,$3,$4,$5", id, exerciseId, 0, time.Now(), 1)
	return err
}

func (r *DataMapper) StartModule(id, moduleId string) error {
	_, err := r.db.Exec("insert into module_progress_histories(user_id,module_id,amount,time,state) values($1,$2,$3,$4,$5", id, moduleId, 0, time.Now(), 1)
	return err
}
