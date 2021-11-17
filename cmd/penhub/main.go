package main

import (
	"github.com/rustamfozilov/penhub/internal/db"
	"github.com/rustamfozilov/penhub/internal/handlers"
	"github.com/rustamfozilov/penhub/internal/services"
	"log"
	"net"

	"net/http"
	"os"
)

func main() {
	host := "0.0.0.0"
	port := "9999"

	newDB, err := db.NewDB()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer newDB.Pool.Close()
	service := services.NewService(newDB)
	handler := handlers.NewHandler(service)
	mux := NewRouter(handler)
	server := http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: mux,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
