package repositories

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/ant0ine/go-urlrouter"
	"github.com/richterrettich/lecture-service/models"
	"github.com/richterrettich/lecture-service/modulepatch"
	"github.com/richterrettich/lecture-service/paginator"
	"github.com/satori/go.uuid"
)

type PModuleRepoFactory struct {
	DB *sql.DB
}

type PModuleRepo struct {
	session *sql.DB
}

func (f *PModuleRepoFactory) CreateRepository() ModuleRepository {
	return &PModuleRepo{f.DB}
}

type InvalidPatchError struct {
	Message string
}

func (e InvalidPatchError) Error() string {
	return e.Message
}

func extractParts(patch *modulepatch.Operation) (string, []string, error) {
	parts := strings.Split(patch.Path, "/")
	if len(parts) == 0 || parts[0] != "" {
		return "", nil, &InvalidPatchError{fmt.Sprintf("Path %s is invalid.", patch.Path)}
	}
	id, parts := parts[1], parts[2:]
	return id, parts, nil
}

func (r *PModuleRepo) ApplyTreePatch(treePatch *modulepatch.Patch) error {
	return nil
}

type commandList struct {
	commands []command
}

type command struct {
	statement  string
	parameters []interface{}
}

func (c *command) execute(tx *sql.Tx) error {
	_, err := tx.Exec(c.statement, c.parameters)
	return err
}

func (c *commandList) executeCommands(db *sql.DB) error {
	tx, err := db.Begin()
	for _, com := range c.commands {
		err = com.execute(tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func createCommand(command string, parameters ...interface{}) *command {
	return &Command{command, parameters}
}

func buildAddModuleCommand(op *modulepatch.Operation, params map[string]string) *command {
	value := op.Value.(map[string]interface{})
	return createCommand(prepare("SELECT insert_module(%v)", value["id"], value["topic_id"], value["description"], value["video_id"], value["script_id"], value["topics"])), nil
}

func buildDeleteModuleTreeCommand(op *modulepatch.Operation, params map[string]string) *command {
	return createCommand(prepare("SELECT delete_module(%v)", params["moduleId"]))
}

func buildDeleteModuleCommand(op *modulepatch.Operation, params map[string]string) *command {
	return createCommand(prepare("SELECT delete_module_tree(%v)", params["moduleId"]))
}

func buildMoveModuleCommand(op *modulepatch.Operation, params map[string]string) *command {
	return createCommand("SELECT move_module($1,null)", params["moduleId"])
}

func buildMoveModuleTreeCommand(op *modulepatch.Operation, params map[string]string) *command {
	return createCommand(`SELECT move_module_tree($1,null)`, params["moduleId"])
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

func ParsePatch(treePatch *modulepatch.Patch) (*CommandList, error) {
	router := urlrouter.Router{
		Routes: []urlrouter.Route{
			urlrouter.Route{
				PathExp: "/",
				Dest:    buildAddModuleCommand,
			},
			urlrouter.Route{
				PathExp: "/:moduleId/", //watch the slash!
				Dest:    buildDeleteModuleTreeCommand,
			},
			urlrouter.Route{
				PathExp: "/:moduleId",
				Dest:    buildDeleteModuleCommand,
			},
			urlrouter.Route{
				PathExp: "/:moduleId/parents",
				Dest:    buildMoveModuleCommand,
			},
			urlrouter.Route{
				PathExp: "/:moduleId/parents/",
				Dest:    buildMoveModuleTreeCommand,
			},
			urlrouter.Route{
				PathExp: "/:moduleId/recommendations",
				Dest:    buildAddRecommandationCommand,
			},
			urlrouter.Route{
				PathExp: "/:moduleId/recommendations/:recommendationId",
				Dest:    buildRemoveRecommendation,
			},
			urlrouter.Route{
				PathExp: "/:moduleId/video",
				Dest:    buildAddVideo,
			},
			urlrouter.Route{
				PathExp: "/:moduleId/video/:videoId",
				Dest:    buildRemoveVideo,
			},
			urlrouter.Route{
				PathExp: "/:moduleId/script",
				Dest:    buildAddScript,
			},
			urlrouter.Route{
				PathExp: "/:moduleId/script/:scriptId",
				Dest:    buildRemoveScript,
			},
			urlrouter.Route{
				PathExp: "/:moduleId/exercises",
				Dest:    buildAddExercise,
			},
			urlrouter.Route{
				PathExp: "/:moduleId/exercises/:exerciseId",
				Dest:    buildRemoveExercise,
			}, urlrouter.Route{
				PathExp: "/:moduleId/exercises/:exerciseId/hints",
				Dest:    buildAddHint,
			}, urlrouter.Route{
				PathExp: "/:moduleId/exercises/:exerciseId/hints/:hintId",
				Dest:    buildRemoveHint,
			}, urlrouter.Route{
				PathExp: "/:moduleId/exercises/:exerciseId/tasks",
				Dest:    buildAddTask,
			},
			urlrouter.Route{
				PathExp: "/:moduleId/exercises/:exerciseId/tasks/:taskId",
				Dest:    buildRemoveTask,
			},
		},
	}
	result := &CommandList{}
	result.AddCommand(`SET TRANSACTION ISOLATION LEVEL SERIALIZABLE`)
	result.AddCommand(`SELECT check_version($1,$2)`, treePatch.LectureID, treePatch.Version)
	for _, op := range treePatch.Operations {
	}
	return result, nil
}

/*
func translateOperation(patchOperation *modulepatch.Operation) (*txn.Op, error) {
	result := &txn.Op{}
	id, parts, err := extractParts(patchOperation)
	if err != nil {
		return nil, err
	}
	result.Id = id
	switch patchOperation.Type {
	}
	return result, nil
}*/

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

func rowToBytes(row *sql.Row) ([]byte, error) {
	var result = make([]byte, 0)
	err := row.Scan(result)
	return result, err
}

func (r *PModuleRepo) GetOne(id string) ([]byte, error) {
	return rowToBytes(r.session.QueryRow(`SELECT get_module($1)`, id))
}

func (r *PModuleRepo) GetByLectureId(topicId string, dr paginator.DepthRequest) ([]byte, error) {
	return rowToBytes(r.session.QueryRow(`SELECT get_module_tree($1,$2,$3)`, topicId, dr.Descendants, dr.Ancestors))
}
