package patchcompiler

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/ant0ine/go-urlrouter"
	"github.com/richterrettich/lecture-service/lecturepatch"
)

type InvalidPatchError struct {
	Message string
}

type PatchCompiler interface {
	Compile(id string, patch *lecturepatch.Patch) (*CommandList, error)
}

func (e InvalidPatchError) Error() string {
	return e.Message
}

type CommandList struct {
	Commands []*command
}

type command struct {
	statement  string
	parameters []interface{}
}

func (c *command) execute(tx *sql.Tx) error {
	_, err := tx.Exec(c.statement, c.parameters)
	return err
}

func (c *CommandList) executeCommands(db *sql.DB) error {
	tx, err := db.Begin()
	for _, com := range c.Commands {
		err = com.execute(tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func createCommand(c string, parameters ...interface{}) *command {
	return &command{c, parameters}
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

type CommandGenerator func(id string, op *lecturepatch.Operation, params map[string]string) (*command, error)

func NewCommandList() *CommandList {
	result := CommandList{
		Commands: make([]*command, 0),
	}

	return &result

}

func (c *CommandList) translatePatch(id string, router *urlrouter.Router, patch *lecturepatch.Patch) error {
	for _, op := range patch.Operations {
		route, params, err := router.FindRoute(op.Path)
		if err != nil {
			return err
		}
		if route == nil {
			return InvalidPatchError{fmt.Sprintf("Invalid Operation. Can't do %s on %s", op.Type, op.Path)}
		}
		builder := route.Dest.(CommandGenerator)

		command, err := builder(id, &op, params)
		if err != nil {
			return err
		}
		c.Commands = append(c.Commands, command)
	}
	return nil
}
