package main

import (
	"flag"
	"lawindex/internal/api"
	"lawindex/internal/archive"
	"lawindex/internal/catalog"
	"lawindex/internal/review"
	"lawindex/internal/store"
	"log"
	"net/http"
)

func main() {
	path := flag.String("db", "lawindex.db", "database path")
	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()
	storage, err := store.Open(*path)
	if err != nil {
		log.Fatal(err)
	}
	defer storage.Close()
	c := catalog.New(storage)
	r := review.New(storage, nil)
	a := archive.New(storage)
	server := api.New(c)
	_ = r
	_ = a
	log.Printf("law index listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, server.Handler()))
}
