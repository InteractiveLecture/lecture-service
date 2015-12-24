package lecturepatch

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/InteractiveLecture/pgmapper/pgutil"
	"github.com/ant0ine/go-urlrouter"
	"github.com/richterrettich/jsonpatch"
)

var PermissionDeniedError = errors.New("Permission Denied")

type CommandGenerator func(id, userId string, officers, assistants map[string]bool, op *jsonpatch.Operation, params map[string]string) (*jsonpatch.CommandContainer, error)

func getTopicAuthority(id string, db *sql.DB) (map[string]bool, map[string]bool, error) {
	stmt := `SELECT user_id,kind from topic_authority where topic_id = $1`
	rows, err := db.Query(stmt, id)
	if err != nil {
		return nil, nil, err
	}
	return extractAuthority(rows)
}

func getModuleAuthority(id string, db *sql.DB) (map[string]bool, map[string]bool, error) {
	stmt := `SELECT ta.user_id,ta.kind 
					 FROM topic_authority ta  
					 INNER JOIN topics t on t.id = ta.topic_id 
					 INNER JOIN modules m on m.topic_id = t.id where m.id = $1`
	rows, err := db.Query(stmt, id)
	if err != nil {
		return nil, nil, err
	}
	return extractAuthority(rows)
}

func getExerciseAuthority(id string, db *sql.DB) (map[string]bool, map[string]bool, error) {
	stmt := `SELECT ta.user_id,ta.kind 
					 FROM topic_authority ta  
					 INNER JOIN topics t on t.id = ta.topic_id 
					 INNER JOIN modules m on m.topic_id = t.id 
					 INNER JOIN exercises e on e.module_id = m.id where e.id = $1`
	rows, err := db.Query(stmt, id)
	if err != nil {
		return nil, nil, err
	}
	return extractAuthority(rows)
}

func extractAuthority(rows *sql.Rows) (map[string]bool, map[string]bool, error) {
	record := make(map[string]interface{})
	officers := make(map[string]bool)
	assistants := make(map[string]bool)
	err := rows.Err()
	if err != nil {
		return nil, nil, err
	}
	for rows.Next() {
		err := rows.Scan(&record)
		if err != nil {
			return nil, nil, err
		}
		if record["kind"].(string) == "ASSISTANT" {
			assistants[record["user_id"].(string)] = true
		}
		if record["kind"].(string) == "OFFICER" {
			officers[record["user_id"].(string)] = true
		}
	}
	return officers, assistants, nil
}

func NewCommandList() *jsonpatch.CommandList {
	result := jsonpatch.CommandList{
		Commands: make([]*jsonpatch.CommandContainer, 0),
	}
	return &result
}

/*func AddCommand(c *jsonpatch.CommandList, command string, values ...interface{}) {
	c.Commands = append(c.Commands, createCommand(command, values...))
}*/

func translatePatch(c *jsonpatch.CommandList, id, userId string, officers, assistants map[string]bool, router *urlrouter.Router, patch *jsonpatch.Patch) error {
	for _, op := range patch.Operations {
		route, params, err := router.FindRoute(op.Path)
		if err != nil {
			return err
		}
		if route == nil {
			return jsonpatch.InvalidPatchError{fmt.Sprintf("Invalid Operation. Can't do %s on %s", op.Type, op.Path)}
		}
		builder := route.Dest.(CommandGenerator)

		command, err := builder(id, userId, officers, assistants, &op, params)
		if err != nil {
			return err
		}
		c.Commands = append(c.Commands, command)
	}
	return nil
}

func buildDefaultMainCallback(stmt string, params ...interface{}) jsonpatch.ContainerCallback {
	return func(transaction, prev interface{}) (interface{}, error) {
		_, err := transaction.(*sql.Tx).Exec(pgutil.Prepare(stmt, params...))
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
}

func buildDefaultCommand(stmt string, params ...interface{}) *jsonpatch.CommandContainer {
	command := new(jsonpatch.CommandContainer)
	command.MainCallback = buildDefaultMainCallback(stmt, params...)
	return command
}

func checkAuthority(id string, authorities ...map[string]bool) error {
	found := false
	for _, a := range authorities {
		if a[id] {
			found = true
			break
		}
	}
	if !found {
		return PermissionDeniedError
	}
	return nil
}

func checkAuthorityAndValidatePatch(assumedOperation, realOperation jsonpatch.OperationType, id string, authorities ...map[string]bool) error {
	if assumedOperation != realOperation {
		return jsonpatch.InvalidPatchError{"Operation not allowed here"}
	}
	return checkAuthority(id, authorities...)
}

func checkStatus(resp *http.Response, err error) error {
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		return fmt.Errorf("There was an error while calling a different service. It returened: %d", resp.StatusCode)
	}
	return nil
}
