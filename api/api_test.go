package api

import (
	"bytes"
	"context"
	"encoding/json"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/dfkdream/hugocms/plugin"

	"github.com/boltdb/bolt"
	"github.com/dfkdream/gallery-plugin/config"
	"github.com/dfkdream/gallery-plugin/database"
	"github.com/dfkdream/hugocms/user"
	"github.com/gorilla/mux"
	"github.com/nfnt/resize"
)

func createTestImage() *bytes.Buffer {
	i := image.NewRGBA(image.Rect(0, 0, 1280, 1280))
	var b bytes.Buffer
	err := jpeg.Encode(&b, i, &jpeg.Options{Quality: 80})
	if err != nil {
		panic(err)
	}
	return &b
}

func mustMarshalJSON(msg interface{}) []byte {
	res, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	res = append(res, "\n"[0])
	return res
}

func newAuthenticatedRequest(method, target string, body io.Reader) *http.Request {
	req := httptest.NewRequest(method, target, body)
	req = req.WithContext(context.WithValue(req.Context(), plugin.ContextKeyUser, &user.User{Id: "hello", Username: "world"}))
	return req
}

func createTestBolt() *bolt.DB {
	dpath, err := ioutil.TempDir("", "gallery-plugin-test-")
	if err != nil {
		panic(err)
	}
	b, err := bolt.Open(path.Join(dpath, "gallery.db"), os.FileMode(0644), nil)
	if err != nil {
		panic(err)
	}
	return b
}

func createTestDB() *database.Database {
	b := createTestBolt()
	db, err := database.New(b, &config.Config{Interpolation: resize.Lanczos3, Quality: 80})
	if err != nil {
		panic(err)
	}
	return db
}

func createTestAPI() *API {
	return New(createTestDB())
}

func TestAPI_SetupHandlers(t *testing.T) {
	for idx, c := range [][]struct {
		req    *http.Request
		code   int
		resp   []byte
		assert func(idx1, idx2 int, a, b []byte, t *testing.T)
	}{
		{
			{
				req:  newAuthenticatedRequest("POST", "/", bytes.NewReader([]byte(`{"title":"hello"}`))),
				code: 200,
				resp: nil,
			}, {
				req:  newAuthenticatedRequest("GET", "/", nil),
				code: 200,
				resp: mustMarshalJSON([]database.Gallery{{Id: 1, Title: "hello"}}),
			}, {
				req:  newAuthenticatedRequest("GET", "/1", nil),
				code: 200,
				resp: mustMarshalJSON(database.Gallery{Id: 1, Title: "hello"}),
			}, {
				req:  newAuthenticatedRequest("POST", "/1", bytes.NewReader([]byte(`{"title":"world"}`))),
				code: 200,
				resp: nil,
			}, {
				req:  newAuthenticatedRequest("GET", "/1", nil),
				code: 200,
				resp: mustMarshalJSON(database.Gallery{Id: 1, Title: "world"}),
			}, {
				req:  newAuthenticatedRequest("DELETE", "/1", nil),
				code: 200,
				resp: nil,
			}, {
				req:  newAuthenticatedRequest("GET", "/", nil),
				code: 200,
				resp: mustMarshalJSON([]database.Gallery{}),
			}, {
				req:  newAuthenticatedRequest("GET", "/1", nil),
				code: 404,
				resp: nil,
			},
		}, {
			{
				req:  newAuthenticatedRequest("POST", "/", bytes.NewReader([]byte(`{"title":"hello"}`))),
				code: 200,
				resp: nil,
			}, {
				req:  newAuthenticatedRequest("GET", "/1/albums", nil),
				code: 200,
				resp: mustMarshalJSON([]database.Album{}),
			}, {
				req:  newAuthenticatedRequest("POST", "/1/albums", bytes.NewReader([]byte(`{"title":"hello"}`))),
				code: 200,
				resp: nil,
			}, {
				req:  newAuthenticatedRequest("GET", "/1/albums", nil),
				code: 200,
				resp: mustMarshalJSON([]database.Album{{Id: 1, Title: "hello", Cover: 0}}),
			}, {
				req:  newAuthenticatedRequest("GET", "/1/album/1", nil),
				code: 200,
				resp: mustMarshalJSON(database.Album{Id: 1, Title: "hello", Cover: 0}),
			}, {
				req:  newAuthenticatedRequest("POST", "/1/album/1", bytes.NewReader([]byte(`{"title":"world"}`))),
				code: 200,
				resp: nil,
			}, {
				req:  newAuthenticatedRequest("GET", "/1/album/1", nil),
				code: 200,
				resp: mustMarshalJSON(database.Album{Id: 1, Title: "world", Cover: 0}),
			}, {
				req:  newAuthenticatedRequest("DELETE", "/1/album/1", nil),
				code: 200,
				resp: nil,
			}, {
				req:  newAuthenticatedRequest("GET", "/1/albums", nil),
				code: 200,
				resp: mustMarshalJSON([]database.Album{}),
			}, {
				req:  newAuthenticatedRequest("GET", "/1/album/1", nil),
				code: 404,
				resp: nil,
			},
		}, {
			{
				req:  newAuthenticatedRequest("POST", "/", bytes.NewReader([]byte(`{"title":"hello"}`))),
				code: 200,
				resp: nil,
			}, {
				req:  newAuthenticatedRequest("POST", "/1/albums", bytes.NewReader([]byte(`{"title":"hello"}`))),
				code: 200,
				resp: nil,
			}, {
				req:  newAuthenticatedRequest("GET", "/1/album/1/images", nil),
				code: 200,
				resp: mustMarshalJSON([]database.Image{}),
			}, {
				req:  newAuthenticatedRequest("POST", "/1/album/1/images", createTestImage()),
				code: 200,
				resp: nil,
			}, {
				req:  newAuthenticatedRequest("GET", "/1/album/1/images", nil),
				code: 200,
				resp: mustMarshalJSON([]database.Image{{Id: 1, Description: ""}}),
			}, {
				req:  newAuthenticatedRequest("GET", "/1/album/1/image/1", nil),
				code: 200,
				assert: func(idx, idx2 int, a, b []byte, t *testing.T) {
					i, _, err := image.Decode(bytes.NewReader(b))
					if err != nil {
						t.Error(idx, idx2, err)
						return
					}
					if i.Bounds() != image.Rect(0, 0, 1280, 1280) {
						t.Error(idx, idx2, "image size not matches:", i.Bounds(), "!=", image.Rect(0, 0, 1280, 1280))
					}
				},
			}, {
				req:  newAuthenticatedRequest("GET", "/1/album/1/image/1?thumb=1", nil),
				code: 200,
				assert: func(idx, idx2 int, a, b []byte, t *testing.T) {
					i, _, err := image.Decode(bytes.NewReader(b))
					if err != nil {
						t.Error(idx, idx2, err)
						return
					}
					if i.Bounds() != image.Rect(0, 0, 720, 720) {
						t.Error(idx, idx2, "image size not matches:", i.Bounds(), "!=", image.Rect(0, 0, 720, 720))
					}
				},
			}, {
				req:  newAuthenticatedRequest("POST", "/1/album/1/image/1", bytes.NewReader([]byte(`{"description":"world"}`))),
				code: 200,
				resp: nil,
			}, {
				req:  newAuthenticatedRequest("GET", "/1/album/1/images", nil),
				code: 200,
				resp: mustMarshalJSON([]database.Image{{Id: 1, Description: "world"}}),
			}, {
				req:  newAuthenticatedRequest("DELETE", "/1/album/1/image/1", nil),
				code: 200,
				resp: nil,
			}, {
				req:  newAuthenticatedRequest("GET", "/1/album/1/images", nil),
				code: 200,
				resp: mustMarshalJSON([]database.Image{}),
			}, {
				req:  newAuthenticatedRequest("GET", "/1/album/1/image/1", nil),
				code: 404,
				resp: nil,
			}, {
				req:  newAuthenticatedRequest("GET", "/1/album/1/image/1?thumb=1", nil),
				code: 404,
				resp: nil,
			},
		},
	} {
		a := createTestAPI()

		m := mux.NewRouter()
		a.SetupHandlers(m)

		for idx2, r := range c {
			res := httptest.NewRecorder()
			m.ServeHTTP(res, r.req)

			if res.Code != r.code {
				t.Error(idx, idx2, "code not matches:", res.Code, "!=", r.code)
				continue
			}

			if r.resp != nil && !bytes.Equal(res.Body.Bytes(), r.resp) {
				t.Error(idx, idx2, "response not matches:", string(res.Body.Bytes()), "!=", string(r.resp))
			}

			if r.assert != nil {
				r.assert(idx, idx2, r.resp, res.Body.Bytes(), t)
			}
		}

	}
}
