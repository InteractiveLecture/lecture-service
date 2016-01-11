package lecturepatch

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/InteractiveLecture/jsonpatch"
	"github.com/InteractiveLecture/serviceclient"
	"github.com/ant0ine/go-urlrouter"
)

type TopicPatchCompiler struct{}

func ForTopics() jsonpatch.PatchCompiler {
	return &TopicPatchCompiler{}
}

//database checked
func generateAddModule(id, userId, jwt string, officers, assistants map[string]bool, op *jsonpatch.Operation, params map[string]string) (*jsonpatch.CommandContainer, error) {
	if err := checkAuthorityAndValidatePatch(jsonpatch.ADD, op.Type, userId, officers); err != nil {
		log.Println("Patch not acceptable: ", err)
		return nil, err
	}
	value := make(map[string]interface{})
	err := json.NewDecoder(strings.NewReader(op.Value.(string))).Decode(&value)
	if err != nil {
		log.Println("error while decoding module: ", err)
		return nil, err
	}
	command := buildDefaultCommand("SELECT add_module(%v)", value["id"].(string), id, value["description"].(string), value["video_id"].(string), value["script_id"].(string), value["parents"])
	command.AfterCallback = func(transaction, prev interface{}) (interface{}, error) {
		client := serviceclient.New("acl-service")
		entity := map[string]interface{}{
			"id":     value["id"].(string),
			"parent": id,
		}
		jsonEntity, _ := json.Marshal(entity)
		err := checkStatus(client.Post("/objects", "json", bytes.NewReader(jsonEntity), "Authorization", jwt))
		if err != nil {
			log.Println("error while creating acl entry: ", err)
			return nil, err
		}
		if len(assistants) > 0 {
			assistantsString := ""
			for _, a := range assistants {
				assistantsString = assistantsString + "&sid=" + strconv.FormatBool(a)
			}
			assistantsString = strings.TrimLeft(assistantsString, "&")
			permissions, _ := json.Marshal(map[string]bool{"read": true, "create": true, "update": false, "delete": false})
			return nil, checkStatus(client.Put(fmt.Sprintf("/objects/%s/permissions?%s", value["id"], assistantsString), "application/json", bytes.NewReader(permissions), "Authorization", jwt))
		}
		return nil, nil
	}
	return command, nil
}

//database checked
func generateRemoveModuleTree(id, userId, jwt string, officers, assistants map[string]bool, op *jsonpatch.Operation, params map[string]string) (*jsonpatch.CommandContainer, error) {
	if err := checkAuthorityAndValidatePatch(jsonpatch.REMOVE, op.Type, userId, officers); err != nil {
		return nil, err
	}
	command := new(jsonpatch.CommandContainer)
	command.MainCallback = func(transaction, prev interface{}) (interface{}, error) {
		row := transaction.(*sql.Tx).QueryRow("SELECT string_agg(id::varchar,'&oid=') from remove_module_tree($1,$2) group by id", id, params["moduleId"])
		var val string
		err := row.Scan(&val)
		if err != nil {
			return nil, err
		}
		return val, nil
	}
	command.AfterCallback = func(transaction, prev interface{}) (interface{}, error) {
		client := serviceclient.New("acl-service")
		return nil, checkStatus(client.Delete(fmt.Sprintf("/objects?oid=%s", prev.(string)), "Authorization", jwt))
	}
	return command, nil
}

//database checked
func generateRemoveModule(id, userId, jwt string, officers, assistants map[string]bool, op *jsonpatch.Operation, params map[string]string) (*jsonpatch.CommandContainer, error) {
	if err := checkAuthorityAndValidatePatch(jsonpatch.REMOVE, op.Type, userId, officers); err != nil {
		return nil, err
	}
	command := buildDefaultCommand("SELECT remove_module(%v)", id, params["moduleId"])
	command.AfterCallback = func(transaction, prev interface{}) (interface{}, error) {
		client := serviceclient.New("acl-service")
		return nil, checkStatus(client.Delete(fmt.Sprintf("/objects/%s", params["moduleId"]), "Authorization", jwt))
	}
	return command, nil
}

//database checked
func generateMoveModule(id, userId, jwt string, officers, assistants map[string]bool, op *jsonpatch.Operation, params map[string]string) (*jsonpatch.CommandContainer, error) {
	if err := checkAuthorityAndValidatePatch(jsonpatch.REPLACE, op.Type, userId, officers); err != nil {
		return nil, err
	}
	return buildDefaultCommand("SELECT move_module(%v)", id, params["moduleId"], op.Value), nil
}

//database checked
func generateMoveModuleTree(id, userId, jwt string, officers, assistants map[string]bool, op *jsonpatch.Operation, params map[string]string) (*jsonpatch.CommandContainer, error) {
	if err := checkAuthorityAndValidatePatch(jsonpatch.REPLACE, op.Type, userId, officers); err != nil {
		return nil, err
	}
	command := buildDefaultCommand("SELECT move_module_tree(%v)", id, params["moduleId"], op.Value)
	return command, nil
}

//database checked
func generateReplaceTopicDescription(id, userId, jwt string, officers, assistants map[string]bool, op *jsonpatch.Operation, params map[string]string) (*jsonpatch.CommandContainer, error) {
	if err := checkAuthorityAndValidatePatch(jsonpatch.REPLACE, op.Type, userId, officers); err != nil {
		return nil, err
	}
	return buildDefaultCommand("SELECT replace_topic_description(%v)", id, op.Value), nil
}

//database checked
func generateAddAssistant(id, userId, jwt string, officers, assistants map[string]bool, op *jsonpatch.Operation, params map[string]string) (*jsonpatch.CommandContainer, error) {
	if err := checkAuthorityAndValidatePatch(jsonpatch.ADD, op.Type, userId, officers); err != nil {
		return nil, err
	}
	command := buildDefaultCommand("SELECT add_assistant(%v)", id, op.Value)
	command.AfterCallback = func(transaction, prev interface{}) (interface{}, error) {
		return nil, setPermissions(id, userId, jwt, transaction.(*sql.Tx), map[string]bool{
			"read_permission":   true,
			"create_permission": true,
			"delete_permission": false,
			"update_permission": false,
		})
	}
	return command, nil
}

func setPermissions(topicId, userId, jwt string, txn *sql.Tx, permissions map[string]bool) error {
	stmt := `SELECT string_agg(id::varchar, $1) from modules where topic_id = $2 GROUP BY id`
	row := txn.QueryRow(stmt, "&oid=", topicId)
	oids := ""
	err := row.Scan(&oids)
	if err != nil {
		return nil
	}
	oids = "oid=" + oids
	client := serviceclient.New("acl-service")
	//TODO post multiple sids in acl einfügen endpiont in acl-service einfügen
	newPermissions, _ := json.Marshal(permissions)
	resp, err := client.Put(fmt.Sprintf("/sids/%s/permissions?%s", userId, oids), "application/json", bytes.NewReader(newPermissions), "Authorization", jwt)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		return errors.New("acl-service returned a not successfull statuscode while setting permissions for new assistant.")
	}
	return nil
}

//databsae checked
func generateRemoveAssistant(id, userId, jwt string, officers, assistants map[string]bool, op *jsonpatch.Operation, params map[string]string) (*jsonpatch.CommandContainer, error) {
	if err := checkAuthorityAndValidatePatch(jsonpatch.REMOVE, op.Type, userId, officers); err != nil {
		return nil, err
	}
	command := buildDefaultCommand("SELECT remove_assistant(%v)", id, params["assistantId"])
	command.AfterCallback = func(transaction, prev interface{}) (interface{}, error) {
		return nil, setPermissions(id, userId, jwt, transaction.(*sql.Tx), map[string]bool{
			"read_permission":   false,
			"create_permission": false,
			"update_permission": false,
			"delete_permission": false,
		})
	}
	return command, nil
}

func (c *TopicPatchCompiler) Compile(treePatch *jsonpatch.Patch, options map[string]interface{}) (*jsonpatch.CommandList, error) {
	id, userId := options["id"].(string), options["userId"].(string)
	db := options["db"].(*sql.DB)
	jwt := options["jwt"].(string)
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
		buildTransactionSerializableCommand(),
		buildDefaultCommand("SELECT check_version(%v)", id, "topics", treePatch.Version),
	)
	err = router.Start()
	if err != nil {
		return nil, err
	}
	err = translatePatch(result, id, userId, jwt, officers, assistants, &router, treePatch)
	if err != nil {
		return nil, err
	}
	result.AddCommands(buildDefaultCommand("SELECT increment_version(%v)", id, "topics"))
	return result, nil
}
