package handler

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"encoding/binary"

	"github.com/InteractiveLecture/id-extractor"
	"github.com/InteractiveLecture/jsonpatch"
	"github.com/InteractiveLecture/lecture-service/lecturepatch"
	"github.com/InteractiveLecture/middlewares/jwtware"
	"github.com/InteractiveLecture/pgmapper"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
)

func GetHintHandler(mapper *pgmapper.Mapper, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		id, err := extractor(r)
		if err != nil {
			return http.StatusBadRequest
		}
		userId := context.Get(r, "user").(*jwt.Token).Claims["id"].(string)
		result, err := mapper.PreparedQueryIntoBytes("SELECT get_hint(%v)", userId, id)
		switch {
		case err != nil:
			return http.StatusInternalServerError
		case len(result) == 0:
			return http.StatusPaymentRequired
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

func PurchaseHintHandler(mapper *pgmapper.Mapper, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		id, err := extractor(r)
		if err != nil {
			return http.StatusBadRequest
		}
		userId := context.Get(r, "user").(*jwt.Token).Claims["id"].(string)
		result, err := mapper.PreparedQueryIntoBytes("SELECT purchase_hint(%v)", id, userId)
		purchaseResult, _ := binary.Varint(result) //TODO check function...something is fishy
		switch {
		case purchaseResult == 0:
			return -1
		case purchaseResult == 1:
			return 420
		case purchaseResult == 2:
			return http.StatusConflict
		case purchaseResult == 3:
			return http.StatusNotFound
		default:
			return http.StatusInternalServerError
		}
	}
	return jwtware.New(createHandler(handlerFunc))
}

/*
func CompleteExerciseHandler(mapper *datamapper.DataMapper, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		id, err := extractor(r)
		if err != nil {
			return http.StatusBadRequest
		}
		var userId string
		err = json.NewDecoder(r.Body).Decode(userId)
		if err != nil {
			return http.StatusBadRequest
		}
		err = mapper.CompleteExercise(id, userId)
		if err != nil {
			return http.StatusInternalServerError
		}
		return -1
	}
	return createHandler(handlerFunc)
}*/

func ExercisePatchHandler(mapper *pgmapper.Mapper, extractor idextractor.Extractor) http.Handler {
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
		compiler := lecturepatch.ForExercises()
		err = mapper.ApplyPatch(id, userId, patch, compiler)
		if err != nil {
			return http.StatusBadRequest
		}
		return -1
	}
	return jwtware.New(createHandler(handlerFunc))
}
