package patchcompiler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ant0ine/go-urlrouter"
	"github.com/richterrettich/lecture-service/lecturepatch"
)

type ExercisePatchCompiler struct{}

func (c *ExercisePatchCompiler) Compile(id string, patch *lecturepatch.Patch) (*CommandList, error) {
	result = &commandList{}
	router := urlrouter.Router{
		Routes: []urlrouter.Route{
			urlrouter.Route{
				PathExp: "/hints", //ADD
				Dest:    CommandGenerator(generateAddHint),
			},
			urlrouter.Route{
				PathExp: "/hints/:hintId", //REMOVE
				Dest:    CommandGenerator(generateRemoveHint),
			},
			urlrouter.Route{
				PathExp: "/hints/:hintId/content", //REPLACE
				Dest:    CommandGenerator(generateUpdateHint),
			},
			urlrouter.Route{
				PathExp: "/tasks", //ADD
				Dest:    CommandGenerator(generateAddTask),
			},
			urlrouter.Route{
				PathExp: "/tasks/:taskPosition", // REMOVE, MOVE
				Dest:    CommandGenerator(generateMoveOrDeleteTask),
			},
			urlrouter.Route{
				PathExp: "/tasks/:taskPosition/task", //command REPLACE
				Dest:    CommandGenerator(generateUpdateTaskCommand),
			},
		},
	}
	err = router.Start()
	if err != nil {
		return
	}
	err = result.translatePatch(id, &router, treePatch)
	return

}

func generateAddTask(id string, op *lecturepatch.Operation, params map[string]string) (*command, error) {
	values := op.Value.(map[string]interface{})
	if op.Type != lecturepatch.ADD {
		return nil, InvalidPatchError{fmt.Sprintf("Only add allowed for %s", op.Path)}
	}
	return createCommand(prepare("SELECT add_task(%v)", values["id"], id, values["position"], values["task"])), nil
}

func generateMoveOrDeleteTask(id string, op *lecturepatch.Operation, params map[string]string) (*command, error) {
	switch op.Type {
	case lecturepatch.REMOVE:
		newPosition, err := strconv.Atoi(params["taskPosition"])
		if err != nil {
			return nil, InvalidPatchError{"From is not valid."}
		}
		return createCommand(prepare("SELECT remove_task($1,$2)", id, newPosition)), nil
	case lecturepatch.MOVE:
		from := strings.Trim(op.From, "/")
		fromParts := strings.Split(from, "/")
		if len(fromParts) != 2 {
			return nil, InvalidPatchError{"From is not valid."}
		}
		oldPosition, err := strconv.Atoi(fromParts[1])
		if err != nil {
			return nil, InvalidPatchError{"From is not valid."}
		}
		newPosition, err := strconv.Atoi(params["taskPosition"])
		if err != nil {
			return nil, InvalidPatchError{"Position is not valid."}
		}
		return createCommand("SELECT move_task($1,$2)", id, oldPosition, newPosition), nil
	default:
		return nil, InvalidPatchError{fmt.Sprintf("Only remove  or move allowed for %s", op.Path)}
	}
}

func generateAddHint(id string, op *lecturepatch.Operation, params map[string]string) (*command, error) {
	value := op.Value.(map[string]interface{})
	position, err := strconv.Atoi(params["hintPosition"])
	if err != nil {
		return nil, InvalidPatchError{"Position is not valid."}
	}
	return createCommand(prepare("SELECT add_hint(%v)", id, position))
}

func generateRemoveHint(id string, op *lecturepatch.Operation, params map[string]string) (*command, error) {
	return createCommand(prepare("SELECT remove_hint", params[":hintId"]))
}

func generateUpdateHint(id string, op *lecturepatch.Operation, params map[string]string) (*command, error) {
	return createCommand("update hints set content = $1 where id = $2", op.Value, params["hintId"])
}

func generateRemoveTask(id string, op *modulepatch.Operation, params map[string]string) (*command, error) {
	return createCommand("delete from tasks where exercise_id = $1 AND position = $2", id, params["taskPosition"]), nil

}

func generateUpdateTaskCommand(id string, op *modulepatch.Operation, params map[string]string) (*command, error) {
	return createCommand("update tasks set task = $1 where exercise_id = $2 AND position = $3", op.Value, id, params["taskPosition"]), nil
}
