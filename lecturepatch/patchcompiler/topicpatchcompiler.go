package patchcompiler

import (
	"github.com/ant0ine/go-urlrouter"
	"github.com/richterrettich/lecture-service/lecturepatch"
)

type TopicPatchCompiler struct{}

func ForTopics() PatchCompiler {
	return &TopicPatchCompiler{}
}

func generateAddModule(id string, op *lecturepatch.Operation, params map[string]string) (*command, error) {
	if op.Type != lecturepatch.ADD {
		return nil, InvalidPatchError{"Operation not allowed here"}
	}
	value := op.Value.(map[string]interface{})
	return createCommand(prepare("SELECT insert_module(%v)", value["id"], id, value["description"], value["video_id"], value["script_id"], value["parents"])), nil
}

func generateDeleteModuleTree(id string, op *lecturepatch.Operation, params map[string]string) (*command, error) {
	if op.Type != lecturepatch.REMOVE {
		return nil, InvalidPatchError{"Operation not allowed here"}
	}
	return createCommand(prepare("SELECT delete_module(%v)", id, params["moduleId"])), nil
}

func generateDeleteModule(id string, op *lecturepatch.Operation, params map[string]string) (*command, error) {
	if op.Type != lecturepatch.REMOVE {
		return nil, InvalidPatchError{"Operation not allowed here"}
	}
	return createCommand(prepare("SELECT delete_module_tree(%v)", id, params["moduleId"])), nil
}

func generateMoveModule(id string, op *lecturepatch.Operation, params map[string]string) (*command, error) {
	if op.Type != lecturepatch.REPLACE {
		return nil, InvalidPatchError{"Operation not allowed here"}
	}
	return createCommand(prepare("SELECT move_module(%v)", id, params["moduleId"])), nil
}

func generateMoveModuleTree(id string, op *lecturepatch.Operation, params map[string]string) (*command, error) {
	if op.Type != lecturepatch.REPLACE {
		return nil, InvalidPatchError{"Operation not allowed here"}
	}
	return createCommand(prepare("SELECT move_module_tree(%v)", id, params["moduleId"])), nil
}

func generateReplaceTopicDescription(id string, op *lecturepatch.Operation, params map[string]string) (*command, error) {
	if op.Type != lecturepatch.REPLACE {
		return nil, InvalidPatchError{"Operation not allowed here"}
	}
	return createCommand("SELECT update_topic_description($1,$2)", id, op.Value), nil
}

func generateAddAssistant(id string, op *lecturepatch.Operation, params map[string]string) (*command, error) {
	if op.Type != lecturepatch.ADD {
		return nil, InvalidPatchError{"Operation not allowed here"}
	}
	return createCommand("SELECT add_assistant($1,$2)", id, op.Value, "ASSISTANT"), nil
}

func generateRemoveAssistant(id string, op *lecturepatch.Operation, params map[string]string) (*command, error) {
	if op.Type != lecturepatch.REMOVE {
		return nil, InvalidPatchError{"Operation not allowed here"}
	}
	return createCommand("SELECT remove_assistant($1,$2)", id, op.Value), nil
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
				Dest:    CommandGenerator(generateDeleteModuleTree),
			},
			urlrouter.Route{
				PathExp: "/modules/:moduleId",
				Dest:    CommandGenerator(generateDeleteModule),
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
	result.Commands = append(result.Commands, createCommand(`SET TRANSACTION ISOLATION LEVEL SERIALIZABLE`))
	result.Commands = append(result.Commands, createCommand(`SELECT check_topic_version($1,$2)`, id, treePatch.Version))
	err := router.Start()
	if err != nil {
		return nil, err
	}
	err = result.translatePatch(id, &router, treePatch)
	if err != nil {
		return nil, err
	}
	return result, nil
}
