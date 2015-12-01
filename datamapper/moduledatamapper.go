package datamapper

import "github.com/richterrettich/lecture-service/paginator"

func (r *DataMapper) GetOneModule(id string) ([]byte, error) {
	return rowToBytes(r.db.QueryRow(`SELECT get_module($1)`, id))
}

func (r *DataMapper) GetModuleRange(topicId string, dr paginator.DepthRequest) ([]byte, error) {
	return rowToBytes(r.db.QueryRow(`SELECT get_module_tree($1,$2,$3)`, topicId, dr.Descendants, dr.Ancestors))
}
