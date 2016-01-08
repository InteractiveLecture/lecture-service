package handler

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/InteractiveLecture/id-extractor"
	"github.com/InteractiveLecture/jsonpatch"
	"github.com/InteractiveLecture/lecture-service/lecturepatch"
	"github.com/InteractiveLecture/lecture-service/paginator"
	"github.com/InteractiveLecture/middlewares/jwtware"
	"github.com/InteractiveLecture/pgmapper"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
)

func ModulesTreeHandler(mapper *pgmapper.Mapper, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		id, err := extractor(r)
		if err != nil {
			log.Println("error with extractor in ModulesTreeHandler")
			return http.StatusBadRequest
		}
		dr, err := paginator.ParseDepth(r.URL)
		if err != nil {
			log.Println("error parsing depth request in ModulesTreeHandler")
			return http.StatusInternalServerError
		}
		result, err := mapper.PreparedQueryIntoBytes("SELECT get_module_tree(%v)", id, dr.Descendants, dr.Ancestors)
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

func ModulesGetHandler(mapper *pgmapper.Mapper, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		id, err := extractor(r)
		if err != nil {
			log.Println("error with extractor in ModulesGetHandler")
			return http.StatusInternalServerError
		}
		result, err := mapper.QueryIntoBytes(`SELECT details from module_details where id = $1`, id)
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

func ModulesPatchHandler(mapper *pgmapper.Mapper, extractor idextractor.Extractor) http.Handler {
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
		compiler := lecturepatch.ForModules()
		err = mapper.ApplyPatch(id, userId, patch, compiler)
		if err != nil {
			return http.StatusBadRequest
		}
		return -1
	}
	return jwtware.New(createHandler(handlerFunc))
}
