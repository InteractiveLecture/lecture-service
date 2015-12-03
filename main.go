package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/InteractiveLecture/id-extractor"
	"github.com/gorilla/mux"
	"github.com/richterrettich/lecture-service/datamapper"
	"github.com/richterrettich/lecture-service/handler"
)

func main() {
	dbHost := flag.String("dbhost", "localhost", "the database host")
	dbPort := flag.Int("dbport", 5432, "the database port")
	dbUser := flag.String("dbuser", "lectureapp", "the database user")
	dbSsl := flag.Bool("dbssl", false, "database ssl config")
	dbName := flag.String("dbname", "lecture", "the database name")
	dbPassword := flag.String("dbpass", "", "database password")
	flag.Parse()
	config := datamapper.DefaultConfig()
	config.Host = *dbHost
	config.Port = *dbPort
	config.User = *dbUser
	config.Ssl = *dbSsl
	config.Database = *dbName
	config.Password = *dbPassword

	mapper, err := datamapper.New(config)
	if err != nil {
		panic(err)
	}

	if testData {

	}

	extractor := idextractor.MuxIdExtractor("id")
	r := mux.NewRouter()

	//TOPICS
	r.Path("/topics").
		Methods("GET").
		Handler(handler.TopicCollectionHandler(mapper))
	r.Path("/topics").
		Methods("POST").
		Handler(handler.TopicCreateHandler(mapper))
	r.Path("/topics/{id}").
		Methods("GET").
		Handler(handler.TopicFindHandler(mapper, extractor))
	r.Path("/topics/{id}").
		Methods("PATCH").
		Handler(handler.TopicPatchHandler(mapper, extractor))
	r.Path("/topics/{id}/officers").
		Methods("POST").
		Handler(handler.TopicAddOfficerHandler(mapper, extractor))
	r.Path("/topics/{id}/officers").
		Methods("DELETE").
		Handler(handler.TopicRemoveOfficerHandler(mapper, extractor))

	//MODULES
	r.Path("/topics/{id}/modules").
		Methods("GET").
		Handler(handler.
		ModulesTreeHandler(mapper, extractor))
	r.Path("/modules/{id}").
		Methods("GET").
		Handler(handler.ModulesGetHandler(mapper, extractor))
	r.Path("/modules/{id}").
		Methods("PATCH").
		Handler(handler.ModulesPatchHandler(mapper, extractor))

	//EXERCISES
	r.Path("/exercises/{id}").
		Methods("POST").
		Handler(handler.CompleteExerciseHandler(mapper, extractor))
	r.Path("/hints/{id}").
		Methods("GET").
		Handler(handler.GetHintHandler(mapper, extractor))
	r.Path("/hints/{id}").
		Methods("POST").
		Handler(handler.PurchaseHintHandler(mapper, extractor))
	//TODO route for GetOneExercise

	//HISTORIES AND PROGRESS
	r.Path("/users/{id}/hints").
		Methods("GET").
		Handler(handler.HintHistoryHandler(mapper, extractor))
	r.Path("/users/{id}/modules").
		Methods("GET").
		Handler(handler.NextModulesForUserHandler(mapper, extractor))
	r.Path("/users/{id}/modules/start").
		Methods("POST").
		Handler(handler.ModuleStartHandler(mapper, extractor))
	r.Path("/users/{id}/modules/next").
		Methods("GET").
		Handler(handler.ModuleHistoryHandler(mapper, extractor))
	r.Path("/users/{id}/exercises").
		Methods("GET").
		Handler(handler.ExerciseStartHandler(mapper, extractor))
	r.Path("/users/{id}/exercises/start").
		Methods("POST").
		Handler(handler.ExerciseHistoryHandler(mapper, extractor))

	log.Println("listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
