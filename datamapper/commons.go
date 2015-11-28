package datamapper

import (
	"database/sql"
	"fmt"

	"github.com/richterrettich/lecture-service/modulepatch"
)

type InvalidPatchError struct {
	Message string
}

func (e InvalidPatchError) Error() string {
	return e.Message
}

type commandList struct {
	commands []command
}

type command struct {
	statement  string
	parameters []interface{}
}

func (c *command) execute(tx *sql.Tx) error {
	_, err := tx.Exec(c.statement, c.parameters)
	return err
}

func (c *commandList) executeCommands(db *sql.DB) error {
	tx, err := db.Begin()
	for _, com := range c.commands {
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

type CommandGenerator func(id string, op *modulepatch.Operation, params map[string]string) command

func (c *commandList) translatePatch(id string, router *urlrouter.Router, patch *modulepatch.Patch) error {
	for _, op := range patch.Operations {
		route, params, err := router.FindRoute(op.Path)
		if err != nil {
			return err
		}
		if route != nil {
			return InvalidPatchError{fmt.Sprintf("Invalid Operation. Can't do %s on %s", op.Type, op.Path)}
		}
		builder := route.Dest.(CommandBuilder)
		c.commands = append(c.commands, builder(id, &op, params))
	}
	return nil
}

func rowToBytes(row *sql.Row) ([]byte, error) {
	var result = make([]byte, 0)
	err := row.Scan(result)
	return result, err
}
