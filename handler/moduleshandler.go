package handler

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/InteractiveLecture/id-extractor"
	"github.com/richterrettich/jsonpatch"
	"github.com/richterrettich/lecture-service/datamapper"
	"github.com/richterrettich/lecture-service/lecturepatch"
	"github.com/richterrettich/lecture-service/paginator"
)

func ModulesTreeHandler(mapper *datamapper.DataMapper, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		log.Println("Within ModulesTreeHandler")
		id, err := extractor(r)
		if err != nil {
			log.Println("error with extractor")
			return http.StatusBadRequest
		}
		dr, err := paginator.ParseDepth(r.URL)
		if err != nil {
			log.Println("error parsing depth request")
			return http.StatusInternalServerError
		}
		log.Println(dr)
		result, err := mapper.GetModuleRange(id, dr)
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

func ModulesGetHandler(mapper *datamapper.DataMapper, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		id, err := extractor(r)
		if err != nil {
			return http.StatusInternalServerError
		}
		result, err := mapper.GetOneModule(id)
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

func ModulesPatchHandler(mapper *datamapper.DataMapper, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		id, err := extractor(r)
		if err != nil {
			return http.StatusInternalServerError
		}
		patch, err := jsonpatch.Decode(r.Body)
		if err != nil {
			return http.StatusBadRequest
		}
		compiler := lecturepatch.ForModules()
		err = mapper.ApplyPatch(id, patch, compiler)
		if err != nil {
			return http.StatusBadRequest
		}
		return -1
	}
	return createHandler(handlerFunc)
}
