package main

import (
	"flag"
	"fmt"
	"kv-go/config"
	"kv-go/db"
	"kv-go/web"
	"log"
	"net/http"

	"github.com/BurntSushi/toml"
)

var (
	dbLocation = flag.String("db-location", "", "The path to the bolt db database")
	httpAddr   = flag.String("http-addr", "127.0.0.1:8080", "HTTP host and port")
	configFile = flag.String("config-file", "sharding.toml", "Config file for static sharding")
	shard      = flag.String("shard", "", "The name of the shard for the data")
)

func parseFlag() {
	flag.Parse()

	if *dbLocation == "" {
		log.Fatal("Must provide db-location")
	}

	if *shard == "" {
		log.Fatal("Must provide shard")
	}
}

func main() {
	fmt.Println("Start kv-go")
	parseFlag()
	var c config.Config
	if _, err := toml.DecodeFile(*configFile, &c); err != nil {
		log.Fatalf("toml.DecodeFile(%q): %v", *dbLocation, err)
	}

	shardCount := len(c.Shard)
	shardIndex := -1
	addressesMap := map[int]string{}

	for _, s := range c.Shard {
		addressesMap[s.Idx] = s.Address
		if s.Name == *shard {
			shardIndex = s.Idx
		}
	}

	if shardIndex < 0 {
		log.Fatalf("Shard %q was not found", *shard)
	}

	log.Printf("Shard count is %d, current shard: %d", shardCount, shardIndex)

	dbInstance, close, err := db.NewDatabase(*dbLocation)
	if err != nil {
		log.Fatalf("NewDatabase(%q): %v", *dbLocation, err)
	}
	defer close()

	srv := web.NewServer(dbInstance, shardIndex, shardCount, addressesMap)
	http.HandleFunc("/get", srv.GetHandler)
	http.HandleFunc("/set", srv.SetHandler)

	log.Fatal(srv.ListenAndServe(httpAddr))

}
