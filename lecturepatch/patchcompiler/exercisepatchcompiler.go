package patchcompiler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ant0ine/go-urlrouter"
	"github.com/richterrettich/lecture-service/lecturepatch"
)

type ExercisePatchCompiler struct{}

func ForExercises() PatchCompiler {
	return &ExercisePatchCompiler{}
}

func (c *ExercisePatchCompiler) Compile(id string, patch *lecturepatch.Patch) (*CommandList, error) {
	result := NewCommandList()
	router := urlrouter.Router{
		Routes: []urlrouter.Route{
			urlrouter.Route{
				PathExp: "/hints", //ADD
				Dest:    CommandGenerator(generateAddHint),
			},
			urlrouter.Route{
				PathExp: "/hints/:hintPosition", //REMOVE, MOVE
				Dest:    CommandGenerator(generateMoveOrRemoveHint),
			},
			urlrouter.Route{
				PathExp: "/hints/:hintPosition/content", //REPLACE
				Dest:    CommandGenerator(generateUpdateHintContent),
			},
			urlrouter.Route{
				PathExp: "/hints/:hintPosition/cost", //REPLACE
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
	err := router.Start()
	if err != nil {
		return nil, err
	}
	result.AddCommand("SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")
	result.AddCommand("SELECT check_exercise_version($1,$2)", id, patch.Version)
	err = result.translatePatch(id, &router, patch)
	if err != nil {
		return nil, err
	}
	result.AddCommand("SELECT increment_exercise_version($1)", id)
	return result, nil
}

func generateAddTask(id string, op *lecturepatch.Operation, params map[string]string) (*command, error) {
	values := op.Value.(map[string]interface{})
	if op.Type != lecturepatch.ADD {
		return nil, InvalidPatchError{fmt.Sprintf("Only add allowed for %s", op.Path)}
	}
	stmt, par := prepare("SELECT add_task(%v)", id, values["position"], values["content"])
	return createCommand(stmt, par...), nil
}

func generateMoveOrRemoveTask(id string, op *lecturepatch.Operation, params map[string]string) (*command, error) {
	switch op.Type {
	case lecturepatch.REMOVE:
		newPosition, err := strconv.Atoi(params["taskPosition"])
		if err != nil {
			return nil, InvalidPatchError{"Error while moving/removing task: not a valid path variable."}
		}
		stmt, par := prepare("SELECT remove_task(%v)", id, newPosition)
		return createCommand(stmt, par...), nil
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
		return createCommand("SELECT move_task($1,$2,$3)", id, oldPosition, newPosition), nil
	default:
		return nil, InvalidPatchError{fmt.Sprintf("Only remove  or move allowed for %s", op.Path)}
	}
}

func generateAddHint(id string, op *lecturepatch.Operation, params map[string]string) (*command, error) {
	if op.Type != lecturepatch.ADD {
		return nil, InvalidPatchError{fmt.Sprintf("Only add allowed for %s", op.Path)}
	}
	value := op.Value.(map[string]interface{})
	stmt, par := prepare("SELECT add_hint(%v)", id, value["position"], value["content"], value["cost"])
	return createCommand(stmt, par...), nil
}

func generateMoveOrRemoveHint(id string, op *lecturepatch.Operation, params map[string]string) (*command, error) {
	switch op.Type {
	case lecturepatch.REMOVE:
		position, err := strconv.Atoi(params["hintPosition"])
		if err != nil {
			return nil, InvalidPatchError{"Error while moving/removing hint: not a valid path vairable."}
		}
		stmt, par := prepare("SELECT remove_hint(%v)", id, position)
		return createCommand(stmt, par...), nil
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
		newPosition, err := strconv.Atoi(params["hintPosition"])
		if err != nil {
			return nil, InvalidPatchError{"Position is not valid."}
		}
		stmt, par := prepare("SELECT move_hint(%v)", id, oldPosition, newPosition)
		return createCommand(stmt, par...), nil
	default:
		return nil, InvalidPatchError{fmt.Sprintf("Only remove  or move allowed for %s", op.Path)}
	}
}

func generateUpdateHintContent(id string, op *lecturepatch.Operation, params map[string]string) (*command, error) {
	if op.Type != lecturepatch.REPLACE {
		return nil, InvalidPatchError{fmt.Sprintf("Only add allowed for %s", op.Path)}
	}
	position, err := strconv.Atoi(params["hintPosition"])
	if err != nil {
		return nil, err
	}
	stmt, par := prepare("SELECT replace_hint_content(%v)", id, position, op.Value)
	return createCommand(stmt, par...), nil
}

func generateUpdateHintCost(id string, op *lecturepatch.Operation, params map[string]string) (*command, error) {
	if op.Type != lecturepatch.REPLACE {
		return nil, InvalidPatchError{fmt.Sprintf("Only add allowed for %s", op.Path)}
	}
	position, err := strconv.Atoi(params["hintPosition"])
	if err != nil {
		return nil, err
	}
	stmt, par := prepare("SELECT replace_hint_cost(%v)", id, position, op.Value)
	return createCommand(stmt, par...), nil
}

func generateUpdateTaskCommand(id string, op *lecturepatch.Operation, params map[string]string) (*command, error) {
	if op.Type != lecturepatch.REPLACE {
		return nil, InvalidPatchError{fmt.Sprintf("Only add allowed for %s", op.Path)}
	}
	position, err := strconv.Atoi(params["taskPosition"])
	if err != nil {
		return nil, err
	}
	stmt, par := prepare("SELECT replace_task_content(%v)", op.Value, id, position)
	return createCommand(stmt, par...), nil
}
