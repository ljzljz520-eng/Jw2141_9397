package main

import (
	"fmt"
	"os"
	"path/filepath"

	"example.com/xiangzhenfarm/internal/cli"
	"example.com/xiangzhenfarm/internal/service"
	"example.com/xiangzhenfarm/internal/store"
)

func main() {
	path := os.Getenv("FARM_REGISTRY_DB")
	if path == "" {
		path = filepath.Join(".", "farmregistry.db")
	}
	database, err := store.Open(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer database.Close()
	app := cli.New(service.NewRegistry(database), os.Stdin, os.Stdout)
	if err := app.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
