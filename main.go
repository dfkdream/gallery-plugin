package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dfkdream/gallery-plugin/api"

	"github.com/boltdb/bolt"
	"github.com/dfkdream/gallery-plugin/config"
	"github.com/dfkdream/gallery-plugin/database"
	"github.com/dfkdream/hugocms/plugin"
)

var (
	version = ""
)

func main() {
	p := plugin.New(plugin.Info{
		Name:        "Gallery",
		Author:      "HugoCMS",
		Description: "Gallery Plugin",
		Version:     version,
		IconClass:   "fas fa-images",
	}, "gallery")

	cfg := config.Get()

	fmt.Println(cfg)

	b, err := bolt.Open(cfg.BoltPath, os.FileMode(0644), nil)
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.New(b, cfg)
	if err != nil {
		log.Fatal(err)
	}

	a := api.New(db)

	a.SetupHandlers(p.APIRouter())

	if err := http.ListenAndServe(":80", p); err != nil {
		log.Fatal(err)
	}
}
