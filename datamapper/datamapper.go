package datamapper

import (
	"database/sql"

	"github.com/richterrettich/jsonpatch"
)

type DataMapper struct {
	db *sql.DB
}

func rowToBytes(row *sql.Row) ([]byte, error) {
	var result = make([]byte, 0)
	err := row.Scan(result)
	return result, err
}
func (mapper *DataMapper) ApplyPatch(id string, patch *jsonpatch.Patch, compiler *jsonpatch.PatchCompiler) error {
	commands, err := compiler.Compile(id, patch)
	if err != nil {
		return err
	}

	tx, err := mapper.db.Begin()

	for _, com := range commands.Commands {
		err = com.ExecuteBefore(tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	for _, com := range commands.Commands {
		err = com.ExecuteMain(tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	for _, com := range commands.Commands {
		err = com.ExecuteAfter(tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (t *DataMapper) queryIntoBytes(query string, params ...interface{}) ([]byte, error) {
	row, err := t.db.Query(query, params)
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
