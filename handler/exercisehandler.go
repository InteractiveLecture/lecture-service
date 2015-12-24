package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/InteractiveLecture/id-extractor"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/richterrettich/jsonpatch"
	"github.com/richterrettich/lecture-service/datamapper"
	"github.com/richterrettich/lecture-service/lecturepatch"
)

func GetHintHandler(mapper *datamapper.DataMapper, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		id, err := extractor(r)
		if err != nil {
			return http.StatusBadRequest
		}
		userId := context.Get(r, "user").(*jwt.Token).Claims["id"].(string)
		result, err := mapper.GetHint(id, userId)
		if _, ok := err.(datamapper.PaymentRequiredError); ok {
			return http.StatusPaymentRequired
		}
		if err != nil {
			log.Println("Error while processing hint ", err)
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

func PurchaseHintHandler(mapper *datamapper.DataMapper, extractor idextractor.Extractor) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) int {
		id, err := extractor(r)
		if err != nil {
			return http.StatusBadRequest
		}
		userId := context.Get(r, "user").(*jwt.Token).Claims["id"].(string)
		err = mapper.PurchaseHint(id, userId)
		if _, ok := err.(datamapper.HintNotFoundError); ok {
			return http.StatusNotFound
		}
		if _, ok := err.(datamapper.InsufficientPointsError); ok {
			return 420
		}
		if _, ok := err.(datamapper.AlreadyPurchasedError); ok {
			return http.StatusConflict
		}
		if err != nil {
			return http.StatusInternalServerError
		}
		return -1
	}
	return createHandler(handlerFunc)
}

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
}

func ExercisePatchHandler(mapper *datamapper.DataMapper, extractor idextractor.Extractor) http.Handler {
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
	return createHandler(handlerFunc)
}
