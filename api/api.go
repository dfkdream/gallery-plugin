package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/dfkdream/gallery-plugin/database"

	"github.com/dfkdream/hugocms/plugin"
	"github.com/gorilla/mux"
)

func atou(s string) (uint64, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return uint64(i), nil
}

type API struct {
	db *database.Database
}

func New(db *database.Database) *API {
	return &API{db: db}
}

// GET: get galleries
// POST: create gallery
func (a *API) galleriesHandler(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		r, err := a.db.GetGalleries()
		if err != nil {
			log.Println(err)
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(res).Encode(r)
		if err != nil {
			log.Println(err)
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	case "POST":
		var values struct {
			Title string `json:"title"`
		}

		err := json.NewDecoder(req.Body).Decode(&values)
		if err != nil {
			http.Error(res, "Bad Request", http.StatusBadRequest)
		}

		gid, err := a.db.CreateGallery(values.Title)
		if err != nil {
			log.Println(err)
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		_, _ = res.Write([]byte(string(gid)))
	default:
		http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// GET: get gallery
// POST: set gallery title
// DELETE: delete gallery
func (a *API) galleryHandler(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	gid, err := atou(vars["gid"])
	if err != nil {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	g, err := a.db.GetGallery(gid)
	if err != nil {
		if err == database.ErrGalleryNotFound {
			http.Error(res, "Not Found", http.StatusNotFound)
			return
		} else {
			log.Println(err)
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	switch req.Method {
	case "GET":
		err := json.NewEncoder(res).Encode(g)
		if err != nil {
			log.Println(err)
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	case "POST":
		var values struct {
			Title string `json:"title"`
		}

		err := json.NewDecoder(req.Body).Decode(&values)
		if err != nil {
			http.Error(res, "Bad Request", http.StatusBadRequest)
			return
		}

		err = a.db.SetGalleryTitle(gid, values.Title)
		if err != nil {
			log.Println(err)
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		_, _ = res.Write([]byte(string(gid)))
	case "DELETE":
		err := a.db.DeleteGallery(gid)
		if err != nil {
			log.Println(err)
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		_, _ = res.Write([]byte(string(gid)))
	default:
		http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// GET: get albums
// POST: create album
func (a *API) albumsHandler(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	gid, err := atou(vars["gid"])
	if err != nil {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	switch req.Method {
	case "GET":
		a, err := a.db.GetAlbums(gid)
		if err != nil {
			log.Println(err)
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(res).Encode(a)
		if err != nil {
			log.Println(err)
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	case "POST":
		var values struct {
			Title string `json:"title"`
		}

		err := json.NewDecoder(req.Body).Decode(&values)
		if err != nil {
			http.Error(res, "Bad Request", http.StatusBadRequest)
			return
		}

		aid, err := a.db.CreateAlbum(gid, values.Title)
		if err != nil {
			log.Println(err)
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		_, _ = res.Write([]byte(string(aid)))
	default:
		http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// GET: get album
// POST: set album title
// DELETE: delete album
func (a *API) albumHandler(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	gid, err := atou(vars["gid"])
	if err != nil {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}
	aid, err := atou(vars["aid"])
	if err != nil {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	album, err := a.db.GetAlbum(gid, aid)
	if err != nil {
		if err == database.ErrGalleryNotFound || err == database.ErrAlbumNotFound {
			http.Error(res, "Not Found", http.StatusNotFound)
			return
		} else {
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	switch req.Method {
	case "GET":
		err := json.NewEncoder(res).Encode(album)
		if err != nil {
			log.Println(err)
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	case "POST":
		var values struct {
			Title string `json:"title"`
		}

		err := json.NewDecoder(req.Body).Decode(&values)
		if err != nil {
			http.Error(res, "Bad Request", http.StatusBadRequest)
			return
		}

		err = a.db.SetAlbumTitle(gid, aid, values.Title)
		if err != nil {
			log.Println(err)
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		_, _ = res.Write([]byte(string(gid)))
	case "DELETE":
		err := a.db.DeleteAlbum(gid, aid)
		if err != nil {
			log.Println(err)
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		_, _ = res.Write([]byte(string(aid)))
	default:
		http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// GET: get images
// POST: add image
func (a *API) imagesHandler(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	gid, err := atou(vars["gid"])
	if err != nil {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}
	aid, err := atou(vars["aid"])
	if err != nil {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	i, err := a.db.GetImages(gid, aid)
	if err != nil {
		if err == database.ErrGalleryNotFound || err == database.ErrAlbumNotFound {
			http.Error(res, "Not Found", http.StatusNotFound)
			return
		} else {
			log.Println(err)
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	switch req.Method {
	case "GET":
		err := json.NewEncoder(res).Encode(i)
		if err != nil {
			log.Println(err)
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	case "POST":
		iid, err := a.db.AddImage(gid, aid, req.Body)
		if err != nil {
			log.Println(err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		_, _ = res.Write([]byte(string(iid)))
	default:
		http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// GET: get image
// POST: set image description
// DELETE: delete image
func (a *API) imageHandler(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	gid, err := atou(vars["gid"])
	if err != nil {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}
	aid, err := atou(vars["aid"])
	if err != nil {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}
	iid, err := atou(vars["iid"])
	if err != nil {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	var img []byte = nil

	if req.URL.Query().Get("thumb") != "" {
		img, err = a.db.GetThumbnail(gid, aid, iid)
	} else {
		img, err = a.db.GetImage(gid, aid, iid)
	}

	if err != nil {
		if err == database.ErrAlbumNotFound || err == database.ErrGalleryNotFound || err == database.ErrImageNotFound {
			http.Error(res, "Not Found", http.StatusNotFound)
			return
		} else {
			log.Println(err)
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	switch req.Method {
	case "GET":
		res.Header().Set("Content-Type", "image/jpeg")
		_, _ = res.Write(img)
	case "POST":
		var values struct {
			Description string `json:"description"`
		}

		err := json.NewDecoder(req.Body).Decode(&values)
		if err != nil {
			http.Error(res, "Bad Request", http.StatusBadRequest)
			return
		}

		err = a.db.SetImageDescription(gid, aid, iid, values.Description)
		if err != nil {
			log.Println(err)
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		_, _ = res.Write([]byte(string(iid)))
	case "DELETE":
		err := a.db.DeleteImage(gid, aid, iid)
		if err != nil {
			log.Println(err)
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		_, _ = res.Write([]byte(string(iid)))
	default:
		http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (a *API) SetupHandlers(r *mux.Router) {
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			if req.Method != "GET" {
				if plugin.GetUser(req) == nil {
					http.Error(res, "Forbidden", http.StatusForbidden)
					return
				}
			}
			next.ServeHTTP(res, req)
		})
	})

	r.HandleFunc("/", a.galleriesHandler)
	r.HandleFunc("/{gid}", a.galleryHandler)
	r.HandleFunc("/{gid}/albums", a.albumsHandler)
	r.HandleFunc("/{gid}/album/{aid}", a.albumHandler)
	r.HandleFunc("/{gid}/album/{aid}/images", a.imagesHandler)
	r.HandleFunc("/{gid}/album/{aid}/image/{iid}", a.imageHandler)
}
