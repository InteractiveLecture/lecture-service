package lecturepatch

import (
	"github.com/ant0ine/go-urlrouter"
	"github.com/richterrettich/jsonpatch"
)

type ModulePatchCompiler struct {
}

func ForModules() jsonpatch.PatchCompiler {
	return ModulePatchCompiler{}
}

func (compiler ModulePatchCompiler) Compile(id string, patch *jsonpatch.Patch) (*jsonpatch.CommandList, error) {
	result := &jsonpatch.CommandList{}
	router := &urlrouter.Router{
		Routes: []urlrouter.Route{
			urlrouter.Route{
				PathExp: "/recommendations",
				Dest:    CommandGenerator(generateAddRecommendation),
			},
			urlrouter.Route{
				PathExp: "/description",
				Dest:    CommandGenerator(generateReplaceDescription),
			},
			urlrouter.Route{
				PathExp: "/recommendations/:recommendationId",
				Dest:    CommandGenerator(generateRemoveRecommendation),
			},
			urlrouter.Route{
				PathExp: "/video",
				Dest:    CommandGenerator(generateAddVideo),
			},
			urlrouter.Route{
				PathExp: "/video/:videoId",
				Dest:    CommandGenerator(generateRemoveVideo),
			},
			urlrouter.Route{
				PathExp: "/script",
				Dest:    CommandGenerator(generateAddScript),
			},
			urlrouter.Route{
				PathExp: "/script/:scriptId",
				Dest:    CommandGenerator(generateRemoveScript),
			},
			urlrouter.Route{
				PathExp: "/exercises",
				Dest:    CommandGenerator(generateAddExercise),
			},
			urlrouter.Route{
				PathExp: "/exercises/:exerciseId",
				Dest:    CommandGenerator(generateRemoveExercise),
			},
		},
	}
	AddCommand(result, `SET TRANSACTION ISOLATION LEVEL SERIALIZABLE`)
	AddCommand(result, `SELECT check_version($1,$2,$3)`, id, "modules", patch.Version)
	err := router.Start()
	if err != nil {
		return nil, err
	}
	err = translatePatch(result, id, router, patch)
	if err != nil {
		return nil, err
	}
	AddCommand(result, "SELECT increment_version($1,$2)", id, "modules")
	AddCommand(result, `REFRESH MATERIALIZED VIEW module_trees`)
	return result, nil
}

// database checked
func generateReplaceDescription(id string, op *jsonpatch.Operation, params map[string]string) (jsonpatch.CommandContainer, error) {
	if op.Type != jsonpatch.REPLACE {
		return nil, jsonpatch.InvalidPatchError{"Operation Not allowed here."}
	}
	return createCommand("SELECT replace_module_description($1,$2)", id, op.Value), nil
}

//database checked
func generateAddRecommendation(id string, op *jsonpatch.Operation, params map[string]string) (jsonpatch.CommandContainer, error) {
	if op.Type != jsonpatch.ADD {
		return nil, jsonpatch.InvalidPatchError{"Operation Not allowed here."}
	}
	return createCommand("SELECT add_module_recommendation($1,$2)", id, op.Value), nil
}

//database checked
func generateRemoveRecommendation(id string, op *jsonpatch.Operation, params map[string]string) (jsonpatch.CommandContainer, error) {
	if op.Type != jsonpatch.REMOVE {
		return nil, jsonpatch.InvalidPatchError{"Operation Not allowed here."}
	}
	return createCommand("SELECT remove_module_recommendation($1,$2)", id, params["recommendationId"]), nil
}

//database checked
func generateAddVideo(id string, op *jsonpatch.Operation, params map[string]string) (jsonpatch.CommandContainer, error) {
	if op.Type != jsonpatch.ADD {
		return nil, jsonpatch.InvalidPatchError{"Operation Not allowed here."}
	}
	return createCommand(`SELECT add_module_video($1,$2)`, id, op.Value), nil
}

// database checked
func generateRemoveVideo(id string, op *jsonpatch.Operation, params map[string]string) (jsonpatch.CommandContainer, error) {
	if op.Type != jsonpatch.REMOVE {
		return nil, jsonpatch.InvalidPatchError{"Operation Not allowed here."}
	}
	return createCommand(`SELECT remove_module_video($1,$2)`, id, params["videoId"]), nil
}

//database checked
func generateAddScript(id string, op *jsonpatch.Operation, params map[string]string) (jsonpatch.CommandContainer, error) {
	if op.Type != jsonpatch.ADD {
		return nil, jsonpatch.InvalidPatchError{"Operation Not allowed here."}
	}
	return createCommand(`SELECT add_module_script($1,$2)`, id, op.Value), nil
}

//dataase checked
func generateRemoveScript(id string, op *jsonpatch.Operation, params map[string]string) (jsonpatch.CommandContainer, error) {
	if op.Type != jsonpatch.REMOVE {
		return nil, jsonpatch.InvalidPatchError{"Operation Not allowed here."}
	}
	return createCommand(`SELECT remove_module_script($1,$2)`, id, params["scriptId"]), nil
}

//database checked
func generateAddExercise(id string, op *jsonpatch.Operation, params map[string]string) (jsonpatch.CommandContainer, error) {
	if op.Type != jsonpatch.ADD {
		return nil, jsonpatch.InvalidPatchError{"Operation Not allowed here."}
	}
	value := op.Value.(map[string]interface{})
	stmt, par := prepare("SELECT add_exercise(%v)", value["id"], id, value["backend"])
	return createCommand(stmt, par...), nil
}

//database checked
func generateRemoveExercise(id string, op *jsonpatch.Operation, params map[string]string) (jsonpatch.CommandContainer, error) {
	if op.Type != jsonpatch.REMOVE {
		return nil, jsonpatch.InvalidPatchError{"Operation Not allowed here."}
	}
	return createCommand("SELECT remove_exercise($1,$2)", id, params["exerciseId"]), nil
}
