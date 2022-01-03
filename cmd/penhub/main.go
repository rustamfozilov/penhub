package main

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/rustamfozilov/penhub/internal/db"
	"github.com/rustamfozilov/penhub/internal/handlers"
	"github.com/rustamfozilov/penhub/internal/services"
	"github.com/rustamfozilov/penhub/internal/types"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	host := "0.0.0.0"
	port := "9999"
	err := execute(host, port)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func execute(host string, port string) error {
	const pathConfig = `D:\penhub\configuration.json`
	config, err := getConfig(pathConfig)
	if err != nil {
		return err
	}
	newDB, err := db.NewDB(config)
	if err != nil {
		return err
	}
	defer newDB.Pool.Close()
	service := services.NewService(newDB, config.ImagesPath)
	handler := handlers.NewHandler(service)
	mux := NewRouter(handler)
	server := http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: mux,
	}
	err = server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func getConfig(path string) (*types.Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var config types.Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &config, nil
}
