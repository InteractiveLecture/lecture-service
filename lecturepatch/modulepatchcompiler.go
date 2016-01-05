package lecturepatch

import (
	"database/sql"
	"strings"

	"github.com/InteractiveLecture/jsonpatch"
	"github.com/InteractiveLecture/serviceclient"
	"github.com/ant0ine/go-urlrouter"
)

type ModulePatchCompiler struct {
}

func ForModules() jsonpatch.PatchCompiler {
	return ModulePatchCompiler{}
}

func (compiler ModulePatchCompiler) Compile(patch *jsonpatch.Patch, options map[string]interface{}) (*jsonpatch.CommandList, error) {
	id, userId := options["id"].(string), options["userId"].(string)
	db := options["db"].(*sql.DB)
	officers, assistants, err := getModuleAuthority(id, db)
	if err != nil {
		return nil, err
	}
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
	result := NewCommandList()
	result.AddCommands(
		buildDefaultCommand("SET TRANSACTION ISOLATION LEVEL SERIALIZABLE"),
		buildDefaultCommand("SELECT check_version(%v)", id, "modules", patch.Version),
	)
	err = router.Start()
	if err != nil {
		return nil, err
	}
	err = translatePatch(result, id, userId, officers, assistants, router, patch)
	if err != nil {
		return nil, err
	}
	result.AddCommands(buildDefaultCommand("SELECT increment_version(%v)", id, "modules"))
	return result, nil
}

// database checked
func generateReplaceDescription(id, userId string, officers, assistants map[string]bool, op *jsonpatch.Operation, params map[string]string) (*jsonpatch.CommandContainer, error) {
	if err := checkAuthorityAndValidatePatch(jsonpatch.REPLACE, op.Type, userId, officers); err != nil {
		return nil, err
	}
	return buildDefaultCommand("SELECT replace_module_description(%v)", id, op.Value), nil
}

//database checked
func generateAddRecommendation(id, userId string, officers, assistants map[string]bool, op *jsonpatch.Operation, params map[string]string) (*jsonpatch.CommandContainer, error) {
	if err := checkAuthorityAndValidatePatch(jsonpatch.ADD, op.Type, userId, officers); err != nil {
		return nil, err
	}
	return buildDefaultCommand("SELECT add_module_recommendation(%v)", id, op.Value), nil
}

//database checked
func generateRemoveRecommendation(id, userId string, officers, assistants map[string]bool, op *jsonpatch.Operation, params map[string]string) (*jsonpatch.CommandContainer, error) {
	if err := checkAuthorityAndValidatePatch(jsonpatch.REMOVE, op.Type, userId, officers); err != nil {
		return nil, err
	}
	return buildDefaultCommand("SELECT remove_module_recommendation(%v)", id, params["recommendationId"]), nil
}

//database checked
func generateAddVideo(id, userId string, officers, assistants map[string]bool, op *jsonpatch.Operation, params map[string]string) (*jsonpatch.CommandContainer, error) {
	if err := checkAuthorityAndValidatePatch(jsonpatch.ADD, op.Type, userId, officers); err != nil {
		return nil, err
	}
	return buildDefaultCommand("SELECT add_module_video(%v)", id, op.Value), nil
}

// database checked
func generateRemoveVideo(id, userId string, officers, assistants map[string]bool, op *jsonpatch.Operation, params map[string]string) (*jsonpatch.CommandContainer, error) {
	if err := checkAuthorityAndValidatePatch(jsonpatch.REMOVE, op.Type, userId, officers); err != nil {
		return nil, err
	}
	return buildDefaultCommand("SELECT remove_module_video(%v)", id, params["videoId"]), nil
}

//database checked
func generateAddScript(id, userId string, officers, assistants map[string]bool, op *jsonpatch.Operation, params map[string]string) (*jsonpatch.CommandContainer, error) {
	if err := checkAuthorityAndValidatePatch(jsonpatch.ADD, op.Type, userId, officers); err != nil {
		return nil, err
	}
	return buildDefaultCommand("SELECT add_module_script(%v)", id, op.Value), nil
}

//dataase checked
func generateRemoveScript(id, userId string, officers, assistants map[string]bool, op *jsonpatch.Operation, params map[string]string) (*jsonpatch.CommandContainer, error) {
	if err := checkAuthorityAndValidatePatch(jsonpatch.REMOVE, op.Type, userId, officers); err != nil {
		return nil, err
	}
	return buildDefaultCommand("SELECT remove_module_script(%v)", id, params["scriptId"]), nil
}

//database checked
func generateAddExercise(id, userId string, officers, assistants map[string]bool, op *jsonpatch.Operation, params map[string]string) (*jsonpatch.CommandContainer, error) {
	if err := checkAuthorityAndValidatePatch(jsonpatch.ADD, op.Type, userId, officers, assistants); err != nil {
		return nil, err
	}
	value := op.Value.(map[string]interface{})
	command := buildDefaultCommand("SELECT add_exercise(%v)", value["id"], id, value["backend"])
	command.AfterCallback = func(transaction, prev interface{}) (interface{}, error) {
		return nil, checkStatus(serviceclient.New("acl-service").Post("/objects", "json", strings.NewReader(value["id"].(string))))
	}
	return command, nil
}

//database checked
func generateRemoveExercise(id, userId string, officers, assistants map[string]bool, op *jsonpatch.Operation, params map[string]string) (*jsonpatch.CommandContainer, error) {
	if err := checkAuthorityAndValidatePatch(jsonpatch.REMOVE, op.Type, userId, officers, assistants); err != nil {
		return nil, err
	}
	command := buildDefaultCommand("SELECT remove_exercise(%v)", id, params["exerciseId"])
	command.AfterCallback = func(transaction, prev interface{}) (interface{}, error) {
		return nil, checkStatus(serviceclient.New("acl-service").Delete("/objects/" + params["exerciseId"]))
	}
	return command, nil
}
