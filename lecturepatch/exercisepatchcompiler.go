package lecturepatch

import (
	"fmt"
	"strconv"

	"github.com/ant0ine/go-urlrouter"
	"github.com/richterrettich/jsonpatch"
)

type ExercisePatchCompiler struct{}

func ForExercises() jsonpatch.PatchCompiler {
	return &ExercisePatchCompiler{}
}

var patchRouter urlrouter.Router
var fromRouter urlrouter.Router

func init() {
	patchRouter = urlrouter.Router{
		Routes: []urlrouter.Route{
			urlrouter.Route{
				PathExp: "/tasks/:taskPosition/hints", //ADD
				Dest:    CommandGenerator(generateAddHint),
			},
			urlrouter.Route{
				PathExp: "/tasks/:taskPosition/hints/:hintPosition", //REMOVE, MOVE
				Dest:    CommandGenerator(generateMoveOrRemoveHint),
			},
			urlrouter.Route{
				PathExp: "/tasks/:taskPosition/hints/:hintPosition/content", //REPLACE
				Dest:    CommandGenerator(generateUpdateHintContent),
			},
			urlrouter.Route{
				PathExp: "/tasks/:taskPosition/hints/:hintPosition/cost", //REPLACE
				Dest:    CommandGenerator(generateUpdateHintCost),
			},
			urlrouter.Route{
				PathExp: "/tasks", //ADD
				Dest:    CommandGenerator(generateAddTask),
			},
			urlrouter.Route{
				PathExp: "/tasks/:taskPosition", // REMOVE, MOVE
				Dest:    CommandGenerator(generateMoveOrRemoveTask),
			},
			urlrouter.Route{
				PathExp: "/tasks/:taskPosition/content", //command REPLACE
				Dest:    CommandGenerator(generateUpdateTaskCommand),
			},
		},
	}
	err := patchRouter.Start()
	if err != nil {
		panic(err)
	}
	fromRouter = urlrouter.Router{
		Routes: []urlrouter.Route{
			urlrouter.Route{
				PathExp: "/tasks/:taskPosition/hints/:hintPosition", //REMOVE, MOVE
				Dest:    "HINT_PATH_ROUTER",
			},
			urlrouter.Route{
				PathExp: "/tasks/:taskPosition", // REMOVE, MOVE
				Dest:    "TASK_PATH_ROUTER",
			},
		},
	}
	err = fromRouter.Start()
	if err != nil {
		panic(err)
	}
}

func (c *ExercisePatchCompiler) Compile(patch *jsonpatch.Patch, options map[string]interface{}) (*jsonpatch.CommandList, error) {
	id, userId := options["id"].(string), options["userId"].(string)
	result := NewCommandList()
	AddCommand(result, "SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")
	AddCommand(result, "SELECT check_version($1,$2,$3)", id, "exercises", patch.Version)
	err := translatePatch(result, id, userId, &patchRouter, patch)
	if err != nil {
		return nil, err
	}
	AddCommand(result, "SELECT increment_version($1,$2)", id, "exercises")
	return result, nil
}

//database checked
func generateAddTask(id, userId string, op *jsonpatch.Operation, params map[string]string) (jsonpatch.CommandContainer, error) {
	values := op.Value.(map[string]interface{})
	if op.Type != jsonpatch.ADD {
		return nil, jsonpatch.InvalidPatchError{fmt.Sprintf("Only add allowed for %s", op.Path)}
	}
	stmt, par := prepare("SELECT add_task(%v)", id, values["id"], values["position"], values["content"])
	return createCommand(stmt, par...), nil
}

// database checked
func generateMoveOrRemoveTask(id, userId string, op *jsonpatch.Operation, params map[string]string) (jsonpatch.CommandContainer, error) {
	switch op.Type {
	case jsonpatch.REMOVE:
		newPosition, err := strconv.Atoi(params["taskPosition"])
		if err != nil {
			return nil, jsonpatch.InvalidPatchError{"Error while moving/removing task: not a valid path variable."}
		}
		stmt, par := prepare("SELECT remove_task(%v)", id, newPosition)
		return createCommand(stmt, par...), nil
	case jsonpatch.MOVE:
		fromParams, err := evalFromRoute(op.From, "TASK_PATH_ROUTER", "taskPosition")
		if err != nil {
			return nil, err
		}
		newPosition, err := strconv.Atoi(params["taskPosition"])
		if err != nil {
			return nil, jsonpatch.InvalidPatchError{"Error while moving/removing task: not a valid path variable."}
		}
		return createCommand("SELECT move_task($1,$2,$3)", id, fromParams[0], newPosition), nil
	default:
		return nil, jsonpatch.InvalidPatchError{fmt.Sprintf("Only remove  or move allowed for %s", op.Path)}
	}
}

func evalFromRoute(from, checkString string, params ...string) ([]int, error) {
	route, routeParams, err := fromRouter.FindRoute(from)
	if err != nil {
		return nil, err
	}
	if route.Dest.(string) != checkString {
		return nil, jsonpatch.InvalidPatchError{fmt.Sprintf("invalid 'FROM' argument %s", from)}
	}
	result := make([]int, 0)
	for _, v := range params {
		param, err := strconv.Atoi(routeParams[v])
		if err != nil {
			return nil, jsonpatch.InvalidPatchError{fmt.Sprintf("Error while evaluating 'FROM': not a valid path variable for taskPosition. %v %v", v, routeParams[v])}
		}
		result = append(result, param)
	}
	return result, nil
}

//database checked
func generateAddHint(id, userId string, op *jsonpatch.Operation, params map[string]string) (jsonpatch.CommandContainer, error) {
	if op.Type != jsonpatch.ADD {
		return nil, jsonpatch.InvalidPatchError{fmt.Sprintf("Only add allowed for %s", op.Path)}
	}
	value := op.Value.(map[string]interface{})
	taskPosition, err := strconv.Atoi(params["taskPosition"])
	if err != nil {
		return nil, err
	}
	stmt, par := prepare("SELECT add_hint(%v)", id, taskPosition, value["id"], value["position"], value["content"], value["cost"])
	return createCommand(stmt, par...), nil
}

//database checked
func generateMoveOrRemoveHint(id, userId string, op *jsonpatch.Operation, params map[string]string) (jsonpatch.CommandContainer, error) {
	switch op.Type {
	case jsonpatch.REMOVE:
		taskPosition, err := strconv.Atoi(params["taskPosition"])
		if err != nil {
			return nil, jsonpatch.InvalidPatchError{"Error while moving/removing hint: not a valid path vairable."}
		}
		hintPosition, err := strconv.Atoi(params["hintPosition"])
		if err != nil {
			return nil, jsonpatch.InvalidPatchError{"Error while moving/removing hint: not a valid path vairable."}
		}
		stmt, par := prepare("SELECT remove_hint(%v)", id, taskPosition, hintPosition)
		return createCommand(stmt, par...), nil
	case jsonpatch.MOVE:
		newHintPosition, err := strconv.Atoi(params["hintPosition"])
		if err != nil {
			return nil, jsonpatch.InvalidPatchError{"Position is not valid."}
		}
		fromParams, err := evalFromRoute(op.From, "HINT_PATH_ROUTER", "taskPosition", "hintPosition")
		if err != nil {
			return nil, err
		}
		newTaskPosition, err := strconv.Atoi(params["taskPosition"])
		if err != nil {
			return nil, jsonpatch.InvalidPatchError{"Position is not valid."}
		}
		stmt, par := prepare("SELECT move_hint(%v)", id, fromParams[0], fromParams[1], newTaskPosition, newHintPosition)
		return createCommand(stmt, par...), nil
	default:
		return nil, jsonpatch.InvalidPatchError{fmt.Sprintf("Only remove  or move allowed for %s", op.Path)}
	}
}

//database checked
func generateUpdateHintContent(id, userId string, op *jsonpatch.Operation, params map[string]string) (jsonpatch.CommandContainer, error) {
	if op.Type != jsonpatch.REPLACE {
		return nil, jsonpatch.InvalidPatchError{fmt.Sprintf("Only add allowed for %s", op.Path)}
	}
	hintPosition, err := strconv.Atoi(params["hintPosition"])
	if err != nil {
		return nil, err
	}
	taskPosition, err := strconv.Atoi(params["taskPosition"])
	if err != nil {
		return nil, err
	}
	stmt, par := prepare("SELECT replace_hint_content(%v)", id, taskPosition, hintPosition, op.Value)
	return createCommand(stmt, par...), nil
}

//database checked
func generateUpdateHintCost(id, userId string, op *jsonpatch.Operation, params map[string]string) (jsonpatch.CommandContainer, error) {
	if op.Type != jsonpatch.REPLACE {
		return nil, jsonpatch.InvalidPatchError{fmt.Sprintf("Only add allowed for %s", op.Path)}
	}
	hintPosition, err := strconv.Atoi(params["hintPosition"])
	if err != nil {
		return nil, err
	}
	taskPosition, err := strconv.Atoi(params["taskPosition"])
	if err != nil {
		return nil, err
	}
	stmt, par := prepare("SELECT replace_hint_cost(%v)", id, taskPosition, hintPosition, op.Value)
	return createCommand(stmt, par...), nil
}

//database checked
func generateUpdateTaskCommand(id, userId string, op *jsonpatch.Operation, params map[string]string) (jsonpatch.CommandContainer, error) {
	if op.Type != jsonpatch.REPLACE {
		return nil, jsonpatch.InvalidPatchError{fmt.Sprintf("Only add allowed for %s", op.Path)}
	}
	position, err := strconv.Atoi(params["taskPosition"])
	if err != nil {
		return nil, err
	}
	stmt, par := prepare("SELECT replace_task_content(%v)", id, position, op.Value)
	return createCommand(stmt, par...), nil
}
