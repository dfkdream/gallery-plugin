package database

import (
	"bytes"
	"image"
	"image/jpeg"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/dfkdream/gallery-plugin/config"
	"github.com/nfnt/resize"

	"github.com/boltdb/bolt"
)

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

func createTestDB() *Database {
	b := createTestBolt()
	db, err := New(b, &config.Config{Interpolation: resize.Lanczos3, Quality: 80})
	if err != nil {
		panic(err)
	}
	return db
}

func createTestImage() bytes.Buffer {
	i := image.NewRGBA(image.Rect(0, 0, 1280, 1280))
	var b bytes.Buffer
	err := jpeg.Encode(&b, i, &jpeg.Options{Quality: 80})
	if err != nil {
		panic(err)
	}
	return b
}

func TestNew(t *testing.T) {
	b := createTestBolt()
	_, err := New(b, &config.Config{Interpolation: resize.Lanczos3, Quality: 80})
	if err != nil {
		t.Error(err)
	}
}

func TestDatabase_CreateGallery(t *testing.T) {
	db := createTestDB()
	_, err := db.CreateGallery("test-gallery")
	if err != nil {
		t.Error(err)
	}
}

func TestDatabase_GetGalleries(t *testing.T) {
	db := createTestDB()
	_, err := db.CreateGallery("test-gallery")
	if err != nil {
		t.Error(err)
	}
	g, err := db.GetGalleries()
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(g, []Gallery{{Id: 1, Title: "test-gallery"}}) {
		t.Errorf("Assertion Failed: %+v", g)
	}
}

func TestDatabase_GetGallery(t *testing.T) {
	db := createTestDB()
	gid, err := db.CreateGallery("test-gallery")
	if err != nil {
		t.Error(err)
	}
	g, err := db.GetGallery(gid)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(g, Gallery{Id: 1, Title: "test-gallery"}) {
		t.Errorf("Assertion Failed: %+v", g)
	}
}

func TestDatabase_SetGalleryTitle(t *testing.T) {
	db := createTestDB()
	id, err := db.CreateGallery("test-gallery")
	if err != nil {
		t.Error(err)
	}
	err = db.SetGalleryTitle(id, "test-gallery-01")
	if err != nil {
		t.Error(err)
	}
	g, err := db.GetGalleries()
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(g, []Gallery{{Id: 1, Title: "test-gallery-01"}}) {
		t.Errorf("Assertion Failed: %+v", g)
	}
}

func TestDatabase_DeleteGallery(t *testing.T) {
	db := createTestDB()
	gid, err := db.CreateGallery("test-gallery")
	if err != nil {
		t.Error(err)
	}
	err = db.DeleteGallery(gid)
	if err != nil {
		t.Error(err)
	}
	g, err := db.GetGalleries()
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(g, []Gallery{}) {
		t.Errorf("Assertion Failed: %+v", g)
	}
}

func TestDatabase_CreateAlbum(t *testing.T) {
	db := createTestDB()
	gid, err := db.CreateGallery("test-gallery")
	if err != nil {
		t.Error(err)
	}
	_, err = db.CreateAlbum(gid, "test-album")
	if err != nil {
		t.Error(err)
	}
}

func TestDatabase_GetAlbums(t *testing.T) {
	db := createTestDB()
	gid, err := db.CreateGallery("test-gallery")
	if err != nil {
		t.Error(err)
	}
	_, err = db.CreateAlbum(gid, "test-album")
	if err != nil {
		t.Error(err)
	}
	a, err := db.GetAlbums(gid)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(a, []Album{{Id: 1, Title: "test-album", Cover: 0}}) {
		t.Errorf("Assertion Failed: %+v", a)
	}
}

func TestDatabase_GetAlbum(t *testing.T) {
	db := createTestDB()
	gid, err := db.CreateGallery("test-gallery")
	if err != nil {
		t.Error(err)
	}
	aid, err := db.CreateAlbum(gid, "test-album")
	if err != nil {
		t.Error(err)
	}
	a, err := db.GetAlbum(gid, aid)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(a, Album{Id: 1, Title: "test-album", Cover: 0}) {
		t.Errorf("Assertion Failed: %+v", a)
	}
}

func TestDatabase_SetAlbumTitle(t *testing.T) {
	db := createTestDB()
	gid, err := db.CreateGallery("test-gallery")
	if err != nil {
		t.Error(err)
	}
	aid, err := db.CreateAlbum(gid, "test-album")
	if err != nil {
		t.Error(err)
	}
	err = db.SetAlbumTitle(gid, aid, "test-album-1")
	if err != nil {
		t.Error(err)
	}
	a, err := db.GetAlbums(gid)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(a, []Album{{Id: 1, Title: "test-album-1", Cover: 0}}) {
		t.Errorf("Assertion Failed: %+v", a)
	}
}

func TestDatabase_DeleteAlbum(t *testing.T) {
	db := createTestDB()
	gid, err := db.CreateGallery("test-gallery")
	if err != nil {
		t.Error(err)
	}
	aid, err := db.CreateAlbum(gid, "test-album")
	if err != nil {
		t.Error(err)
	}
	err = db.DeleteAlbum(gid, aid)
	if err != nil {
		t.Error(err)
	}
	a, err := db.GetAlbums(gid)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(a, []Album{}) {
		t.Errorf("Assertion Failed: %+v", a)
	}
}

func TestDatabase_AddImage(t *testing.T) {
	db := createTestDB()
	gid, err := db.CreateGallery("test-gallery")
	if err != nil {
		t.Error(err)
	}
	aid, err := db.CreateAlbum(gid, "test-album")
	if err != nil {
		t.Error(err)
	}
	img := createTestImage()
	_, err = db.AddImage(gid, aid, &img)
	if err != nil {
		t.Error(err)
	}
}

func TestDatabase_GetImages(t *testing.T) {
	db := createTestDB()
	gid, err := db.CreateGallery("test-gallery")
	if err != nil {
		t.Error(err)
	}
	aid, err := db.CreateAlbum(gid, "test-album")
	if err != nil {
		t.Error(err)
	}
	img := createTestImage()
	_, err = db.AddImage(gid, aid, &img)
	if err != nil {
		t.Error(err)
	}
	i, err := db.GetImages(gid, aid)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(i, []Image{{Id: 1, Description: ""}}) {
		t.Errorf("Assertion Failed: %+v", i)
	}
}

func TestDatabase_SetImageDescription(t *testing.T) {

	db := createTestDB()
	gid, err := db.CreateGallery("test-gallery")
	if err != nil {
		t.Error(err)
	}
	aid, err := db.CreateAlbum(gid, "test-album")
	if err != nil {
		t.Error(err)
	}
	img := createTestImage()
	iid, err := db.AddImage(gid, aid, &img)
	if err != nil {
		t.Error(err)
	}
	err = db.SetImageDescription(gid, aid, iid, "test-image")
	i, err := db.GetImages(gid, aid)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(i, []Image{{Id: 1, Description: "test-image"}}) {
		t.Errorf("Assertion Failed: %+v", i)
	}
}

func TestDatabase_GetImage(t *testing.T) {
	db := createTestDB()
	gid, err := db.CreateGallery("test-gallery")
	if err != nil {
		t.Error(err)
	}
	aid, err := db.CreateAlbum(gid, "test-album")
	if err != nil {
		t.Error(err)
	}
	img := createTestImage()
	iid, err := db.AddImage(gid, aid, &img)
	if err != nil {
		t.Error(err)
	}
	i, err := db.GetImage(gid, aid, iid)
	if err != nil {
		t.Error(err)
	}
	ig, it, err := image.Decode(bytes.NewBuffer(i))
	if err != nil {
		t.Error(err)
	}

	if it != "jpeg" {
		t.Errorf("%s != jpeg", it)
	}

	if ig.Bounds() != image.Rect(0, 0, 1280, 1280) {
		t.Errorf("%+v != %+v", ig.Bounds(), image.Rect(0, 0, 1280, 1280))
	}
}

func TestDatabase_GetThumbnail(t *testing.T) {
	db := createTestDB()
	gid, err := db.CreateGallery("test-gallery")
	if err != nil {
		t.Error(err)
	}
	aid, err := db.CreateAlbum(gid, "test-album")
	if err != nil {
		t.Error(err)
	}
	img := createTestImage()
	iid, err := db.AddImage(gid, aid, &img)
	if err != nil {
		t.Error(err)
	}
	i, err := db.GetThumbnail(gid, aid, iid)
	if err != nil {
		t.Error(err)
	}
	ig, it, err := image.Decode(bytes.NewBuffer(i))
	if err != nil {
		t.Error(err)
	}

	if it != "jpeg" {
		t.Errorf("%s != jpeg", it)
	}

	if ig.Bounds() != image.Rect(0, 0, 360, 360) {
		t.Errorf("%+v != %+v", ig.Bounds(), image.Rect(0, 0, 360, 360))
	}
}

func TestDatabase_DeleteImage(t *testing.T) {
	db := createTestDB()
	gid, err := db.CreateGallery("test-gallery")
	if err != nil {
		t.Error(err)
	}
	aid, err := db.CreateAlbum(gid, "test-album")
	if err != nil {
		t.Error(err)
	}
	img := createTestImage()
	iid, err := db.AddImage(gid, aid, &img)
	if err != nil {
		t.Error(err)
	}

	err = db.DeleteImage(gid, aid, iid)
	if err != nil {
		t.Error(err)
	}

	i, err := db.GetImages(gid, aid)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(i, []Image{}) {
		t.Errorf("Assertion Failed: %+v", i)
	}
}
