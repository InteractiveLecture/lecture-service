package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/InteractiveLecture/id-extractor"
	"github.com/InteractiveLecture/pgmapper"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/richterrettich/jsonpatch"
	"github.com/richterrettich/lecture-service/lecturepatch"
	"github.com/richterrettich/lecture-service/paginator"
)

func TopicCollectionHandler(mapper *pgmapper.Mapper) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		pageRequest, err := paginator.ParsePages(r.URL)
		if err != nil {
			return http.StatusInternalServerError
		}
		result, err := mapper.PreparedQueryIntoBytes(`SELECT * from query_topics($1,$2)`, pageRequest.Number*pageRequest.Size, pageRequest.Size)
		if err != nil {
			log.Println(err)
			return http.StatusInternalServerError
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

func TopicFindHandler(mapper *pgmapper.Mapper, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		id, err := extractor(r)
		if err != nil {
			return http.StatusInternalServerError
		}
		result, err := mapper.PreparedQueryIntoBytes(`SELECT get_topic($1)`, id)
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

func TopicCreateHandler(mapper *pgmapper.Mapper) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		var topic = make(map[string]interface{})
		err := json.NewDecoder(r.Body).Decode(topic)
		if err != nil {
			return http.StatusBadRequest
		}
		err = mapper.Execute("SELECT add_topic(%v)", topic["id"], topic["name"], topic["description"], topic["officers"])
		if err != nil {
			return http.StatusBadRequest // TODO it could be an internal server error as well. need distinction
		}
		w.WriteHeader(http.StatusCreated)
		return -1
	}
	return createHandler(handlerFunc)
}

func TopicPatchHandler(mapper *pgmapper.Mapper, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		id, err := extractor(r)
		if err != nil {
			return http.StatusInternalServerError
		}
		patch, err := jsonpatch.Decode(r.Body)
		if err != nil {
			return http.StatusBadRequest
		}
		userId := context.Get(r, "user").(*jwt.Token).Claims["id"].(string)
		compiler := lecturepatch.ForTopics()
		err = mapper.ApplyPatch(id, userId, patch, compiler)
		if err != nil {
			return http.StatusBadRequest
		}
		return -1
	}
	return createHandler(handlerFunc)
}

func TopicAddOfficerHandler(mapper *pgmapper.Mapper, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		id, err := extractor(r)
		if err != nil {
			return http.StatusInternalServerError
		}
		var officer string
		err = json.NewDecoder(r.Body).Decode(officer)
		if err != nil {
			return http.StatusBadRequest
		}
		err = mapper.Execute(`SELECT add_officer($1,$2)`, id, officer)
		if err != nil {
			return http.StatusInternalServerError
		}
		return -1
	}
	return createHandler(handlerFunc)
}

func TopicRemoveOfficerHandler(mapper *pgmapper.Mapper, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		id, err := extractor(r)
		if err != nil {
			return http.StatusInternalServerError
		}
		var officer string
		err = json.NewDecoder(r.Body).Decode(officer)
		if err != nil {
			return http.StatusBadRequest
		}
		err = mapper.Execute(`SELECT remove_officer($1,$2)`, id, officer)
		if err != nil {
			return http.StatusInternalServerError
		}
		return -1
	}
	return createHandler(handlerFunc)
}
