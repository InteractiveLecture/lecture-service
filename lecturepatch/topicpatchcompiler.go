package lecturepatch

import (
	"github.com/ant0ine/go-urlrouter"
	"github.com/richterrettich/jsonpatch"
)

type TopicPatchCompiler struct{}

func ForTopics() jsonpatch.PatchCompiler {
	return &TopicPatchCompiler{}
}

//database checked
func generateAddModule(id string, op *jsonpatch.Operation, params map[string]string) (jsonpatch.CommandContainer, error) {
	if op.Type != jsonpatch.ADD {
		return nil, jsonpatch.InvalidPatchError{"Operation not allowed here"}
	}
	value := op.Value.(map[string]interface{})
	stmt, parameters := prepare("SELECT add_module(%v)", value["id"], id, value["description"], value["video_id"], value["script_id"], value["parents"])
	command := createCommand(stmt, parameters...)
	command.afterRunCallback = func() error {
		//	client := serviceclient.GetInstance("acl-service")
		// TODO acl stuff
		return nil
	}
	return command, nil
}

//database checked
func generateRemoveModuleTree(id string, op *jsonpatch.Operation, params map[string]string) (jsonpatch.CommandContainer, error) {
	if op.Type != jsonpatch.REMOVE {
		return nil, jsonpatch.InvalidPatchError{"Operation not allowed here"}
	}
	command := createCommand("SELECT remove_module_tree($1,$2)", id, params["moduleId"])
	command.afterRunCallback = func() error {
		//	client := serviceclient.GetInstance("acl-service")
		// TODO acl stuff
		return nil
	}
	return command, nil
}

//database checked
func generateRemoveModule(id string, op *jsonpatch.Operation, params map[string]string) (jsonpatch.CommandContainer, error) {
	if op.Type != jsonpatch.REMOVE {
		return nil, jsonpatch.InvalidPatchError{"Operation not allowed here"}
	}
	command := createCommand("SELECT remove_module($1,$2)", id, params["moduleId"])
	command.afterRunCallback = func() error {
		//	client := serviceclient.GetInstance("acl-service")
		// TODO acl stuff
		return nil
	}

	return command, nil
}

//database checked
func generateMoveModule(id string, op *jsonpatch.Operation, params map[string]string) (jsonpatch.CommandContainer, error) {
	if op.Type != jsonpatch.REPLACE {
		return nil, jsonpatch.InvalidPatchError{"Operation not allowed here"}
	}
	stmt, parameters := prepare("SELECT move_module(%v)", id, params["moduleId"], op.Value)
	return createCommand(stmt, parameters...), nil
}

//database checked
func generateMoveModuleTree(id string, op *jsonpatch.Operation, params map[string]string) (jsonpatch.CommandContainer, error) {
	if op.Type != jsonpatch.REPLACE {
		return nil, jsonpatch.InvalidPatchError{"Operation not allowed here"}
	}
	stmt, parameters := prepare("SELECT move_module_tree(%v)", id, params["moduleId"], op.Value)
	return createCommand(stmt, parameters...), nil
}

//database checked
func generateReplaceTopicDescription(id string, op *jsonpatch.Operation, params map[string]string) (jsonpatch.CommandContainer, error) {
	if op.Type != jsonpatch.REPLACE {
		return nil, jsonpatch.InvalidPatchError{"Operation not allowed here"}
	}
	return createCommand("SELECT replace_topic_description($1,$2)", id, op.Value), nil
}

//database checked
func generateAddAssistant(id string, op *jsonpatch.Operation, params map[string]string) (jsonpatch.CommandContainer, error) {
	if op.Type != jsonpatch.ADD {
		return nil, jsonpatch.InvalidPatchError{"Operation not allowed here"}
	}
	command := createCommand("SELECT add_assistant($1,$2)", id, op.Value)
	command.afterRunCallback = func() error {
		//	client := serviceclient.GetInstance("acl-service")
		// TODO acl stuff
		return nil
	}

	return command, nil
}

//databsae checked
func generateRemoveAssistant(id string, op *jsonpatch.Operation, params map[string]string) (jsonpatch.CommandContainer, error) {
	if op.Type != jsonpatch.REMOVE {
		return nil, jsonpatch.InvalidPatchError{"Operation not allowed here"}
	}
	command := createCommand("SELECT remove_assistant($1,$2)", id, params["assistantId"])
	command.afterRunCallback = func() error {
		//	client := serviceclient.GetInstance("acl-service")
		// TODO acl stuff
		return nil
	}

	return command, nil

}

func (c *TopicPatchCompiler) Compile(id string, treePatch *jsonpatch.Patch) (*jsonpatch.CommandList, error) {
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
	AddCommand(result, `SET TRANSACTION ISOLATION LEVEL SERIALIZABLE`)
	AddCommand(result, `SELECT check_version($1,$2,$3)`, id, "topics", treePatch.Version)
	err := router.Start()
	if err != nil {
		return nil, err
	}
	err = translatePatch(result, id, &router, treePatch)
	if err != nil {
		return nil, err
	}
	AddCommand(result, `SELECT increment_version($1,$2)`, id, "topics")
	//	AddCommand(result, `REFRESH MATERIALIZED VIEW module_trees`)
	return result, nil
}
