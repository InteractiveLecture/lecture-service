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

	extractor := idextractor.MuxIdExtractor("id")
	r := mux.NewRouter()
	r.Path("/topics").
		Methods("GET").
		Handler(handler.TopicCollectionHandler(mapper))
	r.Path("/topics/{id}").
		Methods("GET").
		Handler(handler.TopicFindHandler(mapper, extractor))

	r.Path("/topics/{id}/modules").Methods("GET").Handler(handler.ModulesTreeHandler(mapper, extractor))

	log.Fatal(http.ListenAndServe(":8080", r))
}
