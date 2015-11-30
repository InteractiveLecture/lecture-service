package patchcompiler

import (
	"github.com/ant0ine/go-urlrouter"
	"github.com/richterrettich/lecture-service/lecturepatch"
)

type TopicPatchCompiler struct{}

func ForTopics() PatchCompiler {
	return &TopicPatchCompiler{}
}

//database checked
func generateAddModule(id string, op *lecturepatch.Operation, params map[string]string) (*command, error) {
	if op.Type != lecturepatch.ADD {
		return nil, InvalidPatchError{"Operation not allowed here"}
	}
	value := op.Value.(map[string]interface{})
	stmt, parameters := prepare("SELECT add_module(%v)", value["id"], id, value["description"], value["video_id"], value["script_id"], value["parents"])
	return createCommand(stmt, parameters...), nil
}

//database checked
func generateRemoveModuleTree(id string, op *lecturepatch.Operation, params map[string]string) (*command, error) {
	if op.Type != lecturepatch.REMOVE {
		return nil, InvalidPatchError{"Operation not allowed here"}
	}
	return createCommand("SELECT remove_module_tree($1,$2)", id, params["moduleId"]), nil
}

//database checked
func generateRemoveModule(id string, op *lecturepatch.Operation, params map[string]string) (*command, error) {
	if op.Type != lecturepatch.REMOVE {
		return nil, InvalidPatchError{"Operation not allowed here"}
	}
	return createCommand("SELECT remove_module($1,$2)", id, params["moduleId"]), nil
}

//database checked
func generateMoveModule(id string, op *lecturepatch.Operation, params map[string]string) (*command, error) {
	if op.Type != lecturepatch.REPLACE {
		return nil, InvalidPatchError{"Operation not allowed here"}
	}
	stmt, parameters := prepare("SELECT move_module(%v)", id, params["moduleId"], op.Value)
	return createCommand(stmt, parameters...), nil
}

//database checked
func generateMoveModuleTree(id string, op *lecturepatch.Operation, params map[string]string) (*command, error) {
	if op.Type != lecturepatch.REPLACE {
		return nil, InvalidPatchError{"Operation not allowed here"}
	}
	stmt, parameters := prepare("SELECT move_module_tree(%v)", id, params["moduleId"], op.Value)
	return createCommand(stmt, parameters...), nil
}

//database checked
func generateReplaceTopicDescription(id string, op *lecturepatch.Operation, params map[string]string) (*command, error) {
	if op.Type != lecturepatch.REPLACE {
		return nil, InvalidPatchError{"Operation not allowed here"}
	}
	return createCommand("SELECT replace_topic_description($1,$2)", id, op.Value), nil
}

//database checked
func generateAddAssistant(id string, op *lecturepatch.Operation, params map[string]string) (*command, error) {
	if op.Type != lecturepatch.ADD {
		return nil, InvalidPatchError{"Operation not allowed here"}
	}
	return createCommand("SELECT add_assistant($1,$2)", id, op.Value), nil
}

//databsae checked
func generateRemoveAssistant(id string, op *lecturepatch.Operation, params map[string]string) (*command, error) {
	if op.Type != lecturepatch.REMOVE {
		return nil, InvalidPatchError{"Operation not allowed here"}
	}
	return createCommand("SELECT remove_assistant($1,$2)", id, params["assistantId"]), nil
}

func (c *TopicPatchCompiler) Compile(id string, treePatch *lecturepatch.Patch) (*CommandList, error) {
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
	result.AddCommand(`SET TRANSACTION ISOLATION LEVEL SERIALIZABLE`)
	result.AddCommand(`SELECT check_version($1,$2,$3)`, id, "topics", treePatch.Version)
	err := router.Start()
	if err != nil {
		return nil, err
	}
	err = result.translatePatch(id, &router, treePatch)
	if err != nil {
		return nil, err
	}
	result.AddCommand(`SELECT increment_version($1,$2)`, id, "topics")
	result.AddCommand(`REFRESH MATERIALIZED VIEW module_trees`)
	return result, nil
}
