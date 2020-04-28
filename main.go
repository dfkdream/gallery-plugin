package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dfkdream/gallery-plugin/api"

	"github.com/boltdb/bolt"
	"github.com/dfkdream/gallery-plugin/config"
	"github.com/dfkdream/gallery-plugin/database"
	"github.com/dfkdream/hugocms/plugin"
)

var (
	version = ""
	appHtml = []byte(`<link rel="stylesheet" href="/api/gallery/assets/admin.bundle.css"/>` + "\n" +
		`<div id="app"></div>` + "\n" +
		`<script src="/api/gallery/assets/admin.bundle.js"></script>`)
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

	p.HandleAdminPage("/", "Manage gallery", http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		_, err := res.Write(appHtml)
		if err != nil {
			log.Println(err)
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		}
	}))

	p.AdminPageRouter().PathPrefix("/").HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		_, err := res.Write(appHtml)
		if err != nil {
			log.Println(err)
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		}
	})

	p.APIRouter().PathPrefix("/assets/").HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if strings.HasSuffix(req.URL.Path, "/") {
			http.NotFound(res, req)
			return
		}
		http.StripPrefix("/api/assets", http.FileServer(http.Dir("./assets"))).ServeHTTP(res, req)
	})

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
