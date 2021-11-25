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
 //a := "8810630cf1968b901aaf20e896fa00db8cb5c4102917094ab47f2bb96e98af6e0625f274af0889926b8283c6f5ba1bc775200327c714b7ff32f41c7a142fbaf62301a26ab52e22bee7482cada87e282c5b0f7310702bd1d8ac50f451414bd18672b29e67352604850192c0ebc02b796dee91777b4100fdaaa96ff426ce0e83398f6ea560cd96b22d2ed4ada9068965eddffea63145c6543e78c60b21888b133a4816fc16b6e26e23e6f15b4bb1d903277ff295b920334b247b1a2bf66e8d466621b4e72aa25186edcf5ca094f1ab840393f2227cf1cbdb16dfe7bbbc4a20a3c6014ce597d4d58c9a40049ee6f29afdfcfc116e2b2e8a32d696c61789e249568d"
 //b  := "8810630cf1968b901aaf20e896fa00db8cb5c4102917094ab47f2bb96e98af6e0625f274af0889926b8283c6f5ba1bc775200327c714b7ff32f41c7a142fbaf62301a26ab52e22bee7482cada87e282c5b0f7310702bd1d8ac50f451414bd18672b29e67352604850192c0ebc02b796dee91777b4100fdaaa96ff426ce0e83398f6ea560cd96b22d2ed4ada9068965eddffea63145c6543e78c60b21888b133a4816fc16b6e26e23e6f15b4bb1d903277ff295b920334b247b1a2bf66e8d466621b4e72aa25186edcf5ca094f1ab840393f2227cf1cbdb16dfe7bbbc4a20a3c6014ce597d4d58c9a40049ee6f29afdfcfc116e2b2e8a32d696c61789e249568d"
 //log.Println(a == b)
	//return
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
