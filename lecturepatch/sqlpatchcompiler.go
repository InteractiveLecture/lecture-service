package lecturepatch

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/ant0ine/go-urlrouter"
	"github.com/richterrettich/jsonpatch"
)

type SqlCommandContainer struct {
	statement         string
	parameters        []interface{}
	beforeRunCallback jsonpatch.ContainerCallback
	afterRunCallback  jsonpatch.ContainerCallback
}

func (c *SqlCommandContainer) ExecuteMain(transaction interface{}) error {
	tx := transaction.(*sql.Tx)
	log.Println("executing command: ", c.statement)
	_, err := tx.Exec(c.statement, c.parameters...)
	return err
}

func (c *SqlCommandContainer) ExecuteAfter(transaction interface{}) error {
	if c.afterRunCallback != nil {
		return c.afterRunCallback(transaction)
	}
	return nil
}

func (c *SqlCommandContainer) ExecuteBefore(transaction interface{}) error {
	if c.beforeRunCallback != nil {
		return c.beforeRunCallback(transaction)
	}
	return nil
}

func createCommand(c string, parameters ...interface{}) *SqlCommandContainer {
	return &SqlCommandContainer{c, parameters, nil, nil}
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
				parameters = append(parameters, inval.Interface())
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

type CommandGenerator func(id, userId string, op *jsonpatch.Operation, params map[string]string) (jsonpatch.CommandContainer, error)

func NewCommandList() *jsonpatch.CommandList {
	result := jsonpatch.CommandList{
		Commands: make([]jsonpatch.CommandContainer, 0),
	}
	return &result
}

func AddCommand(c *jsonpatch.CommandList, command string, values ...interface{}) {
	c.Commands = append(c.Commands, createCommand(command, values...))
}

func translatePatch(c *jsonpatch.CommandList, id, userId string, router *urlrouter.Router, patch *jsonpatch.Patch) error {
	for _, op := range patch.Operations {
		route, params, err := router.FindRoute(op.Path)
		if err != nil {
			return err
		}
		if route == nil {
			return jsonpatch.InvalidPatchError{fmt.Sprintf("Invalid Operation. Can't do %s on %s", op.Type, op.Path)}
		}
		builder := route.Dest.(CommandGenerator)

		command, err := builder(id, userId, &op, params)
		if err != nil {
			return err
		}
		c.Commands = append(c.Commands, command)
	}
	return nil
}
