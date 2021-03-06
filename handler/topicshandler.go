package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/InteractiveLecture/id-extractor"
	"github.com/InteractiveLecture/jsonpatch"
	"github.com/InteractiveLecture/lecture-service/lecturepatch"
	"github.com/InteractiveLecture/lecture-service/paginator"
	"github.com/InteractiveLecture/middlewares/jwtware"
	"github.com/InteractiveLecture/pgmapper"
	"github.com/InteractiveLecture/serviceclient"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
)

func TopicCollectionHandler(mapper *pgmapper.Mapper) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		pageRequest, err := paginator.ParsePages(r.URL)
		if err != nil {
			return http.StatusInternalServerError
		}
		result, err := mapper.PreparedQueryIntoBytes("SELECT * from query_topics(%v)", pageRequest.Number*pageRequest.Size, pageRequest.Size)
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
	return jwtware.New(createHandler(handlerFunc))
}

func TopicFindHandler(mapper *pgmapper.Mapper, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		id, err := extractor(r)
		if err != nil {
			return http.StatusInternalServerError
		}
		result, err := mapper.PreparedQueryIntoBytes("SELECT get_topic(%v)", id)
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
	return jwtware.New(createHandler(handlerFunc))
}

func TopicCreateHandler(mapper *pgmapper.Mapper) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		var topic = make(map[string]interface{})
		err := json.NewDecoder(r.Body).Decode(&topic)
		if err != nil {
			log.Println("error while decoding new topic json: ", err)
			return http.StatusBadRequest
		}
		err = mapper.Execute("SELECT add_topic(%v)", topic["id"], topic["name"], topic["description"], topic["officers"])
		if err != nil {
			log.Println("error while inserting new topic into database")
			return http.StatusBadRequest // TODO it could be an internal server error as well. need distinction
		}
		client := serviceclient.New("acl-service")
		aclEntity, _ := json.Marshal(topic)
		resp, err := client.Post("/objects", "application/json", bytes.NewReader(aclEntity), "Authorization", r.Header.Get("Authorization"))
		if err != nil {
			log.Println("error while creating acl-object: ", err)
			return http.StatusInternalServerError
		}
		if resp.StatusCode >= 300 {
			log.Println("got unexpected statuscode from acl-service while creating object: ", resp.StatusCode)
			return http.StatusInternalServerError
		}
		return http.StatusCreated
	}
	return jwtware.New(createHandler(handlerFunc))
}

func TopicPatchHandler(mapper *pgmapper.Mapper, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		id, err := extractor(r)
		if err != nil {
			return http.StatusInternalServerError
		}
		patch, err := jsonpatch.Decode(r.Body)
		if err != nil {
			log.Println("error while decoding patch: ", err)
			return http.StatusBadRequest
		}
		userId := context.Get(r, "user").(*jwt.Token).Claims["id"].(string)
		options := map[string]interface{}{
			"userId": userId,
			"id":     id,
			"jwt":    r.Header.Get("Authorization"),
		}
		compiler := lecturepatch.ForTopics()
		err = mapper.ApplyPatch(patch, compiler, options)
		if err != nil {
			if err == lecturepatch.PermissionDeniedError {
				log.Println(err)
				return http.StatusUnauthorized
			}
			log.Println("error while applying patch: ", err)
			return http.StatusBadRequest
		}
		return -1
	}
	return jwtware.New(createHandler(handlerFunc))
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
		err = mapper.Execute("SELECT add_officer(%v)", id, officer)
		if err != nil {
			return http.StatusInternalServerError
		}
		return -1
	}
	return jwtware.New(createHandler(handlerFunc))
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
		err = mapper.Execute("SELECT remove_officer(%v)", id, officer)
		if err != nil {
			return http.StatusInternalServerError
		}
		return -1
	}
	return jwtware.New(createHandler(handlerFunc))
}
