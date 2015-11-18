package topicshandler

import (
	"net/http"

	"gopkg.in/mgo.v2"
)

func NewCollectionHandler() http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {

		session, err := mgo.Dial("mongo")
		session.DB("lecture-service").C("lectures")
	}
	return createHandler(handlerFunc)
}

func createHandler(handlerFunc http.HandlerFunc) http.Handler {
	return http.Handler(http.HandlerFunc(handlerFunc))
}
