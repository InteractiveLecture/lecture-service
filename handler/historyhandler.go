package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/InteractiveLecture/id-extractor"
	"github.com/InteractiveLecture/lecture-service/paginator"
	"github.com/InteractiveLecture/pgmapper"
)

func HintHistoryHandler(mapper *pgmapper.Mapper, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		pr, err := paginator.ParsePages(r.URL)
		if err != nil {
			return http.StatusInternalServerError
		}
		id, err := extractor(r)
		if err != nil {
			return http.StatusInternalServerError
		}
		var result []byte
		ids, ok := r.URL.Query()["exercise_id"]
		if ok {
			result, err = mapper.PreparedQueryIntoBytes("SELECT get_hint_purchase_history(%v)", id, pr.Size, pr.Size*pr.Number, ids[0])
		} else {
			result, err = mapper.PreparedQueryIntoBytes("SELECT get_hint_purchase_history(%v)", id, pr.Size, pr.Size*pr.Number)
		}
		if err != nil {
			log.Println(err)
			return http.StatusNotFound
		}
		reader := bytes.NewReader(result)
		_, err = io.Copy(w, reader)
		if err != nil {
			log.Println(err)
		}
		return -1
	}
	return createHandler(handlerFunc)
}

func ModuleHistoryHandler(mapper *pgmapper.Mapper, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		pr, err := paginator.ParsePages(r.URL)
		if err != nil {
			return http.StatusInternalServerError
		}
		id, err := extractor(r)
		if err != nil {
			return http.StatusInternalServerError
		}
		var result []byte
		limit := pr.Size
		skip := pr.Size * pr.Number
		if pr.Number == -1 || pr.Size == -1 {
			skip = -1
		}
		ids, ok := r.URL.Query()["topic_id"]
		if ok {
			result, err = mapper.PreparedQueryIntoBytes("SELECT get_module_history(%v)", id, limit, skip, ids[0])
		} else {
			result, err = mapper.PreparedQueryIntoBytes("SELECT get_module_history(%v)", id, limit, skip)
		}
		if err != nil {
			return http.StatusNotFound
		}
		reader := bytes.NewReader(result)
		_, err = io.Copy(w, reader)
		if err != nil {
			log.Println(err)
		}
		return -1
	}
	return createHandler(handlerFunc)
}

func ExerciseHistoryHandler(mapper *pgmapper.Mapper, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		pr, err := paginator.ParsePages(r.URL)
		if err != nil {
			return http.StatusInternalServerError
		}
		id, err := extractor(r)
		if err != nil {
			return http.StatusInternalServerError
		}
		limit := pr.Size
		skip := pr.Size * pr.Number
		if pr.Number == -1 || pr.Size == -1 {
			skip = -1
		}
		var result []byte
		if ids, ok := r.URL.Query()["module_id"]; ok {
			result, err = mapper.PreparedQueryIntoBytes("SELECT get_exercise_history(%v)", id, limit, skip, ids[0])
		} else {
			result, err = mapper.PreparedQueryIntoBytes("SELECT get_exercise_history(%v)", id, limit, skip)
		}
		if err != nil {
			return http.StatusNotFound
		}
		reader := bytes.NewReader(result)
		_, err = io.Copy(w, reader)
		if err != nil {
			log.Println(err)
		}
		return -1
	}
	return createHandler(handlerFunc)
}

func NextModulesForUserHandler(mapper *pgmapper.Mapper, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		id, err := extractor(r)
		if err != nil {
			return http.StatusInternalServerError
		}
		result, err := mapper.PreparedQueryIntoBytes("SELECT get_next_modules_for_user(%v)", id)
		if err != nil {
			return http.StatusNotFound
		}
		reader := bytes.NewReader(result)
		_, err = io.Copy(w, reader)
		if err != nil {
			log.Println(err)
		}
		return -1
	}
	return createHandler(handlerFunc)
}

func TopicBalanceHandler(mapper *pgmapper.Mapper, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		id, err := extractor(r)
		if err != nil {
			return http.StatusInternalServerError
		}
		result, err := mapper.PreparedQueryIntoBytes("Select get_balances(%v)", id)
		if err != nil {
			return http.StatusNotFound
		}
		reader := bytes.NewReader(result)
		_, err = io.Copy(w, reader)
		if err != nil {
			log.Println(err)
		}
		return -1
	}
	return createHandler(handlerFunc)
}

func ModuleStartHandler(mapper *pgmapper.Mapper, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		id, err := extractor(r)
		if err != nil {
			return http.StatusInternalServerError
		}
		var moduleId string
		err = json.NewDecoder(r.Body).Decode(moduleId)
		if err != nil {
			return http.StatusBadRequest
		}
		err = mapper.Execute("insert into module_progress_histories(user_id,module_id,amount,time,state) values(%v)", id, moduleId, 0, time.Now(), 1)
		if err != nil {
			return http.StatusInternalServerError
		}
		return -1
	}
	return createHandler(handlerFunc)
}

func ExerciseStartHandler(mapper *pgmapper.Mapper, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		id, err := extractor(r)
		if err != nil {
			return http.StatusInternalServerError
		}
		var exerciseId string
		err = json.NewDecoder(r.Body).Decode(exerciseId)
		if err != nil {
			return http.StatusBadRequest
		}
		err = mapper.Execute("insert into exercise_progress_histories(user_id,exercise_id,amount,time,state) values(%v)", id, exerciseId, 0, time.Now(), 1)
		if err != nil {
			return http.StatusNotFound
		}
		return -1
	}
	return createHandler(handlerFunc)
}
