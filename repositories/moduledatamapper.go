package datamapper

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ant0ine/go-urlrouter"
	"github.com/richterrettich/lecture-service/models"
	"github.com/richterrettich/lecture-service/modulepatch"
	"github.com/richterrettich/lecture-service/paginator"
	"github.com/satori/go.uuid"
)

func (mapper *Datamapper) ApplyModulePatch(id string, patch *modulepatch.Patch, compiler PatchCompiler) error {
	commands, err := compiler(id, patch)
	if err != nil {
		return err
	}
	return commands.executeCommands(mapper.db)
}

func buildAddRecommendationCommand(op *modulepatch.Operation, params map[string]string) *command {
	return createCommand(prepare("SELECT add_recommendations(%v)", params["moduleId"], op.Value))
}

func buildRemoveRecommendation(op *modulepatch.Operation, params map[string]string) *command {
	return createCommand(prepare("SELECT remove_recommendation(%v)", params["moduleId"], params["recommendationId"]))
}

func buildAddVideo(op *modulepatch.Operation, params map[string]string) *command {
	return createCommand(`SELECT add_video($1,$2)`, params["moduleId"], params["videoId"])
}

func buildRemoveVideo(op *modulepatch.Operation, params map[string]string) *command {
	return createCommand(`SELECT remove_video($1,$2)`, params["moduleId"], params["videoId"])
}

func buildAddScript(op *modulepatch.Operation, params map[string]string) *command {
	return createCommand(`SELECT add_script($1,$2)`, params["moduleId"], params["scriptId"])
}

func buildRemoveScript(op *modulepatch.Operation, params map[string]string) *command {
	return createCommand(`SELECT remove_script($1,$2)`, params["moduleId"], params["scriptId"])
}

func buildAddExercise(op *modulepatch.Operation, params map[string]string) *command {
	value := op.Value.(map[string]interface{})
	return createCommand(prepare("insert into exercises values (%v)", value["id"], params["moduleId"], value["backend"]))
}

func buildRemoveExercise(op *modulepatch.Operation, params map[string]string) *command {
	return createCommand("delete from exercises where id = $1", params["exerciseId"])
}

func parseModulePatch(id string, treepatch *modulepatch.Patch) (*commandList, error) {
	result := commandlist{}
	router := urlrouter.Router{
		Routes: []urlrouter.Route{
			urlrouter.Route{
				pathexp: "/recommendations",
				dest:    buildaddrecommandationcommand,
			},
			urlrouter.Route{
				pathexp: "/recommendations/:recommendationid",
				dest:    buildremoverecommendation,
			},
			urlrouter.Route{
				pathexp: "/video",
				dest:    buildaddvideo,
			},
			urlrouter.Route{
				pathexp: "/video/:videoid",
				dest:    buildremovevideo,
			},
			urlrouter.Route{
				pathexp: "/script",
				dest:    buildaddscript,
			},
			urlrouter.Route{
				pathexp: "/script/:scriptid",
				dest:    buildremovescript,
			},
			urlrouter.Route{
				pathexp: "/exercises",
				dest:    buildaddexercise,
			},
			urlrouter.Route{
				pathexp: "/exercises/:exerciseid",
				dest:    buildaddexercise,
			},
		},
	}
	for _, op := range treepatch.Operations {
		route, params, err := router.FindRoute(op.Path)
		if err != nil {
			return err
		}
		if route == nil {
			return InvalidPatchError{"Operation not supported"}
		}
		commandBuilder := route.Dest.(CommandBuilder)
		result = append(result, commandBuilder(id, op, params))
	}
}

func prepare(stmt string, values ...interface{}) (string, []interface{}) {
	parametersString := ""
	var parameters = make([]interface{}, 0)
	currentIndex := 1
	for _, v := range values {
		val := reflect.ValueOf(v)
		if val.Kind() == reflect.Slice {
			for i := 0; i < val.Len(); i++ {
				inval := val.Index(i)
				parameters = append(parameters, inval)
				parametersString = fmt.Sprintf("%s,$%d", parametersString, currentIndex)
				currentIndex = currentIndex + 1
			}
		} else {
			parameters = append(parameters, v)
			parametersString = fmt.Sprintf("%s,$%d", parametersString, currentIndex)
			currentIndex = currentIndex + 1
		}
	}
	stmt = fmt.Sprintf(stmt, strings.Trim(parametersString, ","))
	return stmt, parameters
}

func (r *PModuleRepo) Create(m *models.Module) error {
	if m.ID == "" {
		m.ID = uuid.NewV4().String()
	}
	_, err := r.session.Exec(prepare("SELECT insert_module(%s)", m.ID, m.TopicID, m.Description, m.VideoID, m.ScriptID, m.Parents))
	return err
}

func (r *PModuleRepo) GetOne(id string) ([]byte, error) {
	return rowToBytes(r.session.QueryRow(`SELECT get_module($1)`, id))
}

func (r *PModuleRepo) GetByLectureId(topicId string, dr paginator.DepthRequest) ([]byte, error) {
	return rowToBytes(r.session.QueryRow(`SELECT get_module_tree($1,$2,$3)`, topicId, dr.Descendants, dr.Ancestors))
}
