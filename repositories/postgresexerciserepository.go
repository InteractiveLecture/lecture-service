package repositories

import "github.com/richterrettich/lecture-service/modulepatch"

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

func buildAddDescriptionCommand(op *modulepatch.Operation, params map[string]string) *command {

}

func ParseExercisePatch(treePatch *modulepatch.Patch) (*CommandList, error) {
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
				PathExp: "/tasks/:taskId", // REMOVE
				Dest:    buildRemoveTask,
			},
			urlrouter.Route{
				PathExp: "/tasks/:taskId/task", //command REPLACE
				Dest:    buildUpdateTaskCommand,
			},
		},
	}
}
