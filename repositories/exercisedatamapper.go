package datamapper

import (
	"github.com/ant0ine/go-urlrouter"
	"github.com/richterrettich/lecture-service/modulepatch"
)

func buildAddTask(op *modulepatch.Operation, params map[string]string) *command {
	values := op.Value.(map[string]interface{})
	return createCommand(prepare("insert into tasks values(%v)", values["id"], params["exerciseId"], values["task"]))
}

func buildDeleteTask(op *modulepatch.Operation, params map[string]string) *command {
	return createCommand(prepare("delete from tasks where id =%v", params["taskId"]))
}

func buildAddHint(op *modulepatch.Operation, params map[string]string) *command {
	return createCommand(prepare("delete from tasks where id =%v", params["taskId"]))
}

func buildRemoveHint(op *modulepatch.Operation, params map[string]string) *command {
	return createCommand(prepare("delete from hints where id = %v", params[":hintId"]))
}

func buildUpdateHint(op *modulepatch.Operation, params map[string]string) *command {
	return createCommand("update hints set content = $1 where id = $2", op.Value, params["hintId"])
}

func buildRemoveTask(id string, op *modulepatch.Operation, params map[string]string) *command {
	return createCommand("delete from tasks where exercise_id = $1 AND position = $2", id, params["taskPosition"])

}

func buildUpdateTaskCommand(id string, op *modulepatch.Operation, params map[string]string) *command {
	return createCommand("update tasks set task = $1 where exercise_id = $2 AND position = $3", op.Value, id, params["taskPosition"])
}

func ParseExercisePatch(id string, treePatch *modulepatch.Patch) (result *commandList, err error) {
	result = &commandList{}
	router := urlrouter.Router{
		Routes: []urlrouter.Route{
			urlrouter.Route{
				PathExp: "/hints", //ADD
				Dest:    buildAddHint,
			},
			urlrouter.Route{
				PathExp: "/hints/:hintId", //REMOVE
				Dest:    buildRemoveHint,
			},
			urlrouter.Route{
				PathExp: "/hints/:hintId/content", //REPLACE
				Dest:    buildUpdateHint,
			},
			urlrouter.Route{
				PathExp: "/tasks", //ADD
				Dest:    buildAddTask,
			},
			urlrouter.Route{
				PathExp: "/tasks/:taskPosition", // REMOVE
				Dest:    buildRemoveTask,
			},
			urlrouter.Route{
				PathExp: "/tasks/:taskPosition/task", //command REPLACE
				Dest:    buildUpdateTaskCommand,
			},
		},
	}
	err = result.translatePatch(id, &router, treePatch)
	return
}
