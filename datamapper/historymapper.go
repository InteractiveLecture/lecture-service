package datamapper

import "github.com/richterrettich/lecture-service/paginator"

func (r *DataMapper) GetHintHistory(userId string, pr paginator.PageRequest) ([]byte, error) {
	return r.queryIntoBytes(`SELECT get_hint_purchase_history($1,$2,$3)`, userId, pr.Size, pr.Size*pr.Number)
}

func (r *DataMapper) GetModuleHistory(userId string, pr paginator.PageRequest) ([]byte, error) {
	return r.queryIntoBytes(`SELECT get_module_history($1,$2,$3)`, userId, pr.Size, pr.Size*pr.Number)
}

func (r *DataMapper) GetExerciseHistory(userId string, pr paginator.PageRequest) ([]byte, error) {
	return r.queryIntoBytes(`SELECT get_exercise_history($1,$2,$3)`, userId, pr.Size, pr.Size*pr.Number)
}

func (r *DataMapper) GetNextModulesForUser(id string) ([]byte, error) {
	return r.queryIntoBytes("SELECT get_next_modules_for_user($1)", id)
}

func (r *DataMapper) GetCurrentModulesForUser(id string) ([]byte, error) {
	return r.queryIntoBytes("SELECT get_current_modules_for_user($1)", id)
}
