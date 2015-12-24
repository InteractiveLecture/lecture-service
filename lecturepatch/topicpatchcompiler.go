package lecturepatch

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/InteractiveLecture/serviceclient"
	"github.com/ant0ine/go-urlrouter"
	"github.com/richterrettich/jsonpatch"
)

type TopicPatchCompiler struct{}

func ForTopics() jsonpatch.PatchCompiler {
	return &TopicPatchCompiler{}
}

//database checked
func generateAddModule(id, userId string, officers, assistants map[string]bool, op *jsonpatch.Operation, params map[string]string) (*jsonpatch.CommandContainer, error) {
	if err := checkAuthorityAndValidatePatch(jsonpatch.ADD, op.Type, userId, officers); err != nil {
		return nil, err
	}
	value := op.Value.(map[string]interface{})
	command := buildDefaultCommand("SELECT add_module(%v)", value["id"], id, value["description"], value["video_id"], value["script_id"], value["parents"])
	command.AfterCallback = func(transaction, prev interface{}) (interface{}, error) {
		client := serviceclient.GetInstance("acl-service")
		err := checkStatus(client.Post("/objects", "json", strings.NewReader(value["id"].(string))))
		if err != nil {
			return nil, err
		}
		assistantsString := ""
		for _, a := range assistants {
			assistantsString = assistantsString + "&sid=" + strconv.FormatBool(a)
		}
		assistantsString = strings.TrimLeft(assistantsString, "&")
		permissions, _ := json.Marshal(map[string]bool{"read": true, "create": true, "update": false, "delete": false})
		return nil, checkStatus(client.Put(fmt.Sprintf("/objects/%s/permissions?%s", value["id"], assistantsString), "application/json", bytes.NewReader(permissions)))
	}
	return command, nil
}

//database checked
func generateRemoveModuleTree(id, userId string, officers, assistants map[string]bool, op *jsonpatch.Operation, params map[string]string) (*jsonpatch.CommandContainer, error) {
	if err := checkAuthorityAndValidatePatch(jsonpatch.REMOVE, op.Type, userId, officers); err != nil {
		return nil, err
	}
	command := new(jsonpatch.CommandContainer)
	command.MainCallback = func(transaction, prev interface{}) (interface{}, error) {
		row := transaction.(*sql.Tx).QueryRow("SELECT string_agg(id,'&oid=') from remove_module_tree($1,$2) group by id", id, params["moduleId"])
		var val string
		err := row.Scan(&val)
		if err != nil {
			return nil, err
		}
		return val, nil
	}
	command.AfterCallback = func(transaction, prev interface{}) (interface{}, error) {
		client := serviceclient.GetInstance("acl-service")
		return nil, checkStatus(client.Delete(fmt.Sprintf("/objects?oid=%s", prev.(string))))
	}
	return command, nil
}

//database checked
func generateRemoveModule(id, userId string, officers, assistants map[string]bool, op *jsonpatch.Operation, params map[string]string) (*jsonpatch.CommandContainer, error) {
	if err := checkAuthorityAndValidatePatch(jsonpatch.REMOVE, op.Type, userId, officers); err != nil {
		return nil, err
	}
	command := buildDefaultCommand("SELECT remove_module(%v)", id, params["moduleId"])
	command.AfterCallback = func(transaction, prev interface{}) (interface{}, error) {
		client := serviceclient.GetInstance("acl-service")
		return nil, checkStatus(client.Delete(fmt.Sprintf("/objects/%s", params["moduleId"])))
	}
	return command, nil
}

//database checked
func generateMoveModule(id, userId string, officers, assistants map[string]bool, op *jsonpatch.Operation, params map[string]string) (*jsonpatch.CommandContainer, error) {
	if err := checkAuthorityAndValidatePatch(jsonpatch.REPLACE, op.Type, userId, officers); err != nil {
		return nil, err
	}
	return buildDefaultCommand("SELECT move_module(%v)", id, params["moduleId"], op.Value), nil
}

//database checked
func generateMoveModuleTree(id, userId string, officers, assistants map[string]bool, op *jsonpatch.Operation, params map[string]string) (*jsonpatch.CommandContainer, error) {
	if err := checkAuthorityAndValidatePatch(jsonpatch.REPLACE, op.Type, userId, officers); err != nil {
		return nil, err
	}
	command := buildDefaultCommand("SELECT move_module_tree(%v)", id, params["moduleId"], op.Value)
	return command, nil
}

//database checked
func generateReplaceTopicDescription(id, userId string, officers, assistants map[string]bool, op *jsonpatch.Operation, params map[string]string) (*jsonpatch.CommandContainer, error) {
	if err := checkAuthorityAndValidatePatch(jsonpatch.REPLACE, op.Type, userId, officers); err != nil {
		return nil, err
	}
	return buildDefaultCommand("SELECT replace_topic_description(%v)", id, op.Value), nil
}

//database checked
func generateAddAssistant(id, userId string, officers, assistants map[string]bool, op *jsonpatch.Operation, params map[string]string) (*jsonpatch.CommandContainer, error) {
	if err := checkAuthorityAndValidatePatch(jsonpatch.ADD, op.Type, userId, officers); err != nil {
		return nil, err
	}
	command := buildDefaultCommand("SELECT add_assistant(%v)", id, op.Value)
	command.AfterCallback = func(transaction, prev interface{}) (interface{}, error) {
		return nil, setPermissions(id, userId, transaction.(*sql.Tx), map[string]bool{
			"read":   true,
			"create": true,
			"delete": false,
			"update": false,
		})
	}
	return command, nil
}

func setPermissions(topicId, userId string, txn *sql.Tx, permissions map[string]bool) error {
	stmt := `SELECT string_agg(id, '&oid=') from modules where topic_id = $1 GROUP BY id`
	row := txn.QueryRow(stmt, topicId)
	oids := ""
	err := row.Scan(&oids)
	if err != nil {
		return nil
	}
	oids = "oid=" + oids
	client := serviceclient.GetInstance("acl-service")
	//TODO post multiple sids in acl einfügen endpiont in acl-service einfügen
	newPermissions, _ := json.Marshal(permissions)
	resp, err := client.Put(fmt.Sprintf("/sids/%s/permissions?%s", userId, oids), "application/json", bytes.NewReader(newPermissions))
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		return errors.New("acl-service returned a not successfull statuscode while setting permissions for new assistant.")
	}
	return nil
}

//databsae checked
func generateRemoveAssistant(id, userId string, officers, assistants map[string]bool, op *jsonpatch.Operation, params map[string]string) (*jsonpatch.CommandContainer, error) {
	if err := checkAuthorityAndValidatePatch(jsonpatch.REMOVE, op.Type, userId, officers); err != nil {
		return nil, err
	}
	command := buildDefaultCommand("SELECT remove_assistant(%v)", id, params["assistantId"])
	command.AfterCallback = func(transaction, prev interface{}) (interface{}, error) {
		return nil, setPermissions(id, userId, transaction.(*sql.Tx), map[string]bool{
			"read":   false,
			"create": false,
			"update": false,
			"delete": false,
		})
	}
	return command, nil
}

func (c *TopicPatchCompiler) Compile(treePatch *jsonpatch.Patch, options map[string]interface{}) (*jsonpatch.CommandList, error) {
	id, userId := options["id"].(string), options["userId"].(string)
	db := options["db"].(*sql.DB)
	officers, assistants, err := getTopicAuthority(id, db)
	if err != nil {
		return nil, err
	}
	router := urlrouter.Router{
		Routes: []urlrouter.Route{
			urlrouter.Route{
				PathExp: "/description",
				Dest:    CommandGenerator(generateReplaceTopicDescription), //REPLACE
			},
			urlrouter.Route{
				PathExp: "/assistants",
				Dest:    CommandGenerator(generateAddAssistant),
			},
			urlrouter.Route{
				PathExp: "/assistants/:assistantId",
				Dest:    CommandGenerator(generateRemoveAssistant),
			},
			urlrouter.Route{
				PathExp: "/modules",
				Dest:    CommandGenerator(generateAddModule),
			},
			urlrouter.Route{
				PathExp: "/modules/:moduleId/tree",
				Dest:    CommandGenerator(generateRemoveModuleTree),
			},
			urlrouter.Route{
				PathExp: "/modules/:moduleId",
				Dest:    CommandGenerator(generateRemoveModule),
			},
			urlrouter.Route{
				PathExp: "/modules/:moduleId/parents",
				Dest:    CommandGenerator(generateMoveModule),
			},
			urlrouter.Route{
				PathExp: "/modules/:moduleId/parents/tree",
				Dest:    CommandGenerator(generateMoveModuleTree),
			},
		},
	}
	result := NewCommandList()
	result.AddCommands(
		buildDefaultCommand(`SET TRANSACTION ISOLATION LEVEL SERIALIZABLE`),
		buildDefaultCommand("SELECT check_version(%v)", id, "topics", treePatch.Version),
	)
	err = router.Start()
	if err != nil {
		return nil, err
	}
	err = translatePatch(result, id, userId, officers, assistants, &router, treePatch)
	if err != nil {
		return nil, err
	}
	result.AddCommands(buildDefaultCommand("SELECT increment_version(%v)", id, "topics"))
	return result, nil
}
