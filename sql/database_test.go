package database

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

/*


			FOO
			 |
			BAR
			 |
			BLI
			/ \
  BLUBB  BLA
	  \    /
		 BAZZ


*/

func TestMoveSingle(t *testing.T) {
	db, err := dbConnect()
	assert.Nil(t, err)
	defer db.Close()

	modules := getModules(t, db)
	assert.Equal(t, len(modules), 6)
	_, err = db.Exec(`SELECT move_module($1,$2)`, modules["bazz"].Id, modules["foo"].Id)
	assert.Nil(t, err)
	modules = getModules(t, db)
	assert.Equal(t, len(modules), 6) //lenght shouldnt have changed.
	assert.Equal(t, 1, modules["bazz"].level)
	_, err = db.Exec(`SELECT move_module($1,$2)`, modules["foo"].Id, modules["bla"].Id)
	assert.Nil(t, err)
	modules = getModules(t, db)
	assert.Equal(t, 3, modules["foo"].level)
	assert.Equal(t, 0, modules["bar"].level)
	for _, v := range modules {
		assert.True(t, strings.HasPrefix(v.paths[0], "/"+modules["bar"].Id))
	}
	topicId := getTopicId(t, db, modules["foo"].Id)
	_, err = db.Exec(`SELECT check_version($1,$2)`, topicId, 2)
	assert.Nil(t, err)
}

func getTopicId(t *testing.T, db *sql.DB, moduleId string) string {
	rows, err := db.Query(`SELECT t.id FROM topics t inner join modules m on t.id = m.topic_id where m.id = $1`, moduleId)
	assert.Nil(t, err)
	var id string
	rows.Next()
	rows.Scan(&id)
	return id
}

func TestMoveTree(t *testing.T) {
	db, err := dbConnect()
	assert.Nil(t, err)
	defer db.Close()
	modules := getModules(t, db)
	_, err = db.Exec(`SELECT move_module_tree($1,$2)`, modules["bli"].Id, modules["foo"].Id)
	assert.Nil(t, err)
	modules = getModules(t, db)
	for i := 0; i < 2; i++ {
		assert.False(t, strings.Contains(modules["bazz"].paths[i], modules["bar"].Id))
	}
	assert.Equal(t, 1, modules["bar"].level)
	_, err = db.Exec(`SELECT move_module_tree($1,$2)`, modules["bli"].Id, modules["foo"].Id)
	assert.Nil(t, err)
	assert.Equal(t, modules["bli"].Id, getDirectParents(modules["bla"])[0])
}

func TestDeleteModule(t *testing.T) {
	db, err := dbConnect()
	assert.Nil(t, err)
	defer db.Close()
	modules := getModules(t, db)
	_, err = db.Exec(`SELECT delete_module($1)`, modules["bli"].Id)
	assert.Nil(t, err)
	modules = getModules(t, db)
	parents := getDirectParents(modules["bla"])
	assert.Equal(t, 1, len(parents))
	assert.Equal(t, modules["bar"].Id, parents[0])
	parents = getDirectParents(modules["blubb"])
	assert.Equal(t, 1, len(parents))
	assert.Equal(t, modules["bar"].Id, parents[0])
	_, err = db.Exec(`SELECT delete_module($1)`, modules["foo"].Id)
	assert.Nil(t, err)
	_, err = db.Exec(`SELECT delete_module($1)`, modules["bar"].Id)
	assert.Nil(t, err)
	modules = getModules(t, db)
	assert.Equal(t, 3, len(modules))
	assert.Equal(t, modules["bla"].Id, getDirectParents(modules["blubb"])[0])
	assert.Equal(t, fmt.Sprintf("/%s/%s/%s", modules["bla"].Id, modules["blubb"].Id, modules["bazz"].Id), modules["bazz"].paths[0])
}

func TestInsertModule(t *testing.T) {
	db, err := dbConnect()
	assert.Nil(t, err)
	defer db.Close()
	modules := getModules(t, db)
	_, err = db.Exec(`SELECT insert_module($1,$2,$3,$4)`, uuid.NewV4().String(), modules["foo"].topicId, "hugo", modules["foo"].Id)
	assert.Nil(t, err)
	modules = getModules(t, db)
	val, ok := modules["hugo"]
	assert.True(t, ok)
	assert.Equal(t, modules["foo"].Id, getDirectParents(val)[0])
	assert.Equal(t, fmt.Sprintf("/%s/%s", modules["foo"].Id, val.Id), val.paths[0])
	assert.Equal(t, modules["foo"].Id, getDirectParents(modules["bar"])[0])
}

func TestDeleteModuleTree(t *testing.T) {
	db, err := dbConnect()
	assert.Nil(t, err)
	defer db.Close()
	modules := getModules(t, db)
	_, err = db.Exec(`SELECT insert_module($1,$2,$3,$4)`, uuid.NewV4().String(), modules["foo"].topicId, "hugo", modules["foo"].Id)
	assert.Nil(t, err)
	_, err = db.Exec(`SELECT delete_module_tree($1)`, modules["bar"].Id)
	assert.Nil(t, err)
	modules = getModules(t, db)
	assert.Equal(t, 2, len(modules))
	for _, v := range []string{"foo", "hugo"} {
		_, ok := modules[v]
		assert.True(t, ok)
	}
	assert.Equal(t, fmt.Sprintf("/%s/%s", modules["foo"].Id, modules["hugo"].Id), modules["hugo"].paths[0])
}

func getDirectParents(m module) []string {
	var result = make([]string, 0)
	for _, v := range m.paths {
		parts := strings.Split(v, "/")
		result = append(result, parts[len(parts)-2])
	}
	return result
}

func getModules(t *testing.T, db *sql.DB) map[string]module {
	rows, err := db.Query(`SELECT id,description,level,paths, topic_id FROM module_trees`)
	assert.Nil(t, err)
	defer rows.Close()
	var id, description, paths, topicId string
	var level int
	var modules = make(map[string]module)
	for rows.Next() {
		err = rows.Scan(&id, &description, &level, &paths, &topicId)
		assert.Nil(t, err)
		modules[description] = module{id, description, parseArray(paths), topicId, level}
	}
	return modules
}

type module struct {
	Id          string
	description string
	paths       []string
	topicId     string
	level       int
}

func parseArray(arr string) []string {
	step1 := arr[1 : len(arr)-1]
	return strings.Split(step1, ",")
}

func dbConnect() (*sql.DB, error) {
	host := os.Getenv("PGHOST")
	if host == "" {
		host = "localhost"
	}
	_, err := exec.Command("./prepare_data.sh").Output()
	if err != nil {
		panic(err)
	}
	return sql.Open("postgres", fmt.Sprintf("postgres://lectureapp@%s/lecture?sslmode=disable", host))
}