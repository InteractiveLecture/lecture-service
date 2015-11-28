package datamapper

import (
	"database/sql"

	"github.com/richterrettich/lecture-service/modulepatch"
)

type Datamapper struct {
	db *sql.DB
}

type PatchCompiler func(id string, patch modulepatch.Patch) (*commandList, error)
