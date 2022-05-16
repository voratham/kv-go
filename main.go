package main

import (
	"flag"
	"fmt"
	"kv-go/db"
	"kv-go/web"
	"log"
	"net/http"
)

var (
	dbLocation = flag.String("db-location", "", "The path to the bolt db database")
	httpAddr   = flag.String("http-addr", "127.0.0.1:8080", "HTTP host and port")
)

func parseFlag() {
	flag.Parse()

	if *dbLocation == "" {
		log.Fatal("Must provide db-location")
	}
}

func main() {
	fmt.Println("Start kv-go")
	parseFlag()

	dbInstance, close, err := db.NewDatabase(*dbLocation)
	if err != nil {
		log.Fatalf("NewDatabase(%q): %v", *dbLocation, err)
	}
	defer close()

	srv := web.NewServer(dbInstance)
	http.HandleFunc("/get", srv.GetHandler)
	http.HandleFunc("/set", srv.SetHandler)

	log.Fatal(srv.ListenAndServe(httpAddr))

}
