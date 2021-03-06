package database

import (
	"bytes"
	"encoding/binary"
	"errors"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"time"

	"github.com/boltdb/bolt"
	"github.com/dfkdream/gallery-plugin/config"
	"github.com/nfnt/resize"
)

var (
	ErrGalleryNotFound = errors.New("gallery not found")
	ErrAlbumNotFound   = errors.New("album not found")
	ErrImageNotFound   = errors.New("image not found")
)

var (
	galleryBucket  = []byte("gallery")
	titleKey       = []byte("title")
	albumsBucket   = []byte("album")
	imagesBucket   = []byte("images")
	imageKey       = []byte("image")
	thumbnailKey   = []byte("thumbnail")
	timestampKey   = []byte("timestamp")
	descriptionKey = []byte("description")
)

type Database struct {
	db  *bolt.DB
	cfg *config.Config
}

func New(db *bolt.DB, cfg *config.Config) (*Database, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(galleryBucket)
		return err
	})

	if err != nil {
		return nil, err
	}

	return &Database{db: db, cfg: cfg}, nil
}

func itob(id uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, id)
	return b
}

func btoi(v []byte) uint64 {
	return binary.BigEndian.Uint64(v)
}

type Gallery struct {
	Id    uint64 `json:"id"`
	Title string `json:"title"`
}

func (d *Database) GetGalleries() ([]Gallery, error) {
	result := make([]Gallery, 0)
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(galleryBucket)
		c := b.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			result = append(result, Gallery{
				Id:    btoi(k),
				Title: string(b.Bucket(k).Get(titleKey)),
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, err
}

func (d *Database) GetGallery(galleryId uint64) (Gallery, error) {
	var result Gallery
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(galleryBucket)
		b = b.Bucket(itob(galleryId))
		if b == nil {
			return ErrGalleryNotFound
		}
		result.Title = string(b.Get(titleKey))
		result.Id = galleryId
		return nil
	})

	return result, err
}

// CreateGallery creates gallery and return uint64 auto-incremental key
func (d *Database) CreateGallery(title string) (uint64, error) {
	var id uint64

	err := d.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(galleryBucket)
		if err != nil {
			return err
		}
		id, err = b.NextSequence()
		if err != nil {
			return err
		}
		bkt, err := b.CreateBucket(itob(id))
		if err != nil {
			return err
		}

		_, err = bkt.CreateBucket(albumsBucket)
		if err != nil {
			return err
		}

		return bkt.Put(titleKey, []byte(title))
	})

	return id, err
}

// DeleteGallery deletes gallery
func (d *Database) DeleteGallery(id uint64) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(galleryBucket).DeleteBucket(itob(id))
	})
}

func (d *Database) SetGalleryTitle(id uint64, title string) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		g := tx.Bucket(galleryBucket)
		b := g.Bucket(itob(id))
		if b == nil {
			return ErrGalleryNotFound
		}

		return b.Put(titleKey, []byte(title))
	})
}

type Album struct {
	Id    uint64 `json:"id"`
	Title string `json:"title"`
	Cover uint64 `json:"cover"`
}

func (d *Database) GetAlbums(galleryId uint64) ([]Album, error) {
	result := make([]Album, 0)

	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(galleryBucket)
		b = b.Bucket(itob(galleryId))
		if b == nil {
			return ErrGalleryNotFound
		}
		b = b.Bucket(albumsBucket)
		c := b.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			title := string(b.Bucket(k).Get(titleKey))

			var cover uint64 = 0
			key, _ := b.Bucket(k).Bucket(imagesBucket).Cursor().First()
			if key != nil {
				cover = btoi(key)
			}
			result = append(result, Album{
				Id:    btoi(k),
				Title: title,
				Cover: cover,
			})
		}
		return nil
	})

	return result, err
}

func (d *Database) GetAlbum(galleryId, albumId uint64) (Album, error) {
	var result Album
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(galleryBucket)
		b = b.Bucket(itob(galleryId))
		if b == nil {
			return ErrGalleryNotFound
		}
		b = b.Bucket(albumsBucket)
		b = b.Bucket(itob(albumId))
		if b == nil {
			return ErrAlbumNotFound
		}
		result.Id = albumId
		result.Title = string(b.Get(titleKey))
		var cover uint64 = 0
		if k, _ := b.Bucket(imagesBucket).Cursor().First(); k != nil {
			cover = btoi(k)
		}
		result.Cover = cover
		return nil
	})

	return result, err
}

func (d *Database) CreateAlbum(galleryId uint64, title string) (uint64, error) {
	var albumId uint64

	err := d.db.Update(func(tx *bolt.Tx) error {
		g := tx.Bucket(galleryBucket)
		b := g.Bucket(itob(galleryId))
		if b == nil {
			return ErrGalleryNotFound
		}
		b = b.Bucket(albumsBucket)
		var err error
		albumId, err = b.NextSequence()
		if err != nil {
			return err
		}
		a, err := b.CreateBucket(itob(albumId))
		if err != nil {
			return err
		}
		_, err = a.CreateBucket(imagesBucket)
		if err != nil {
			return err
		}

		return a.Put(titleKey, []byte(title))
	})

	return albumId, err
}

func (d *Database) SetAlbumTitle(galleryId, albumId uint64, title string) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		g := tx.Bucket(galleryBucket)
		b := g.Bucket(itob(galleryId))
		if b == nil {
			return ErrGalleryNotFound
		}
		b = b.Bucket(albumsBucket)
		b = b.Bucket(itob(albumId))
		if b == nil {
			return ErrAlbumNotFound
		}
		return b.Put(titleKey, []byte(title))
	})
}

func (d *Database) DeleteAlbum(galleryId, albumId uint64) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(galleryBucket)
		b = b.Bucket(itob(galleryId))
		if b == nil {
			return ErrGalleryNotFound
		}
		b = b.Bucket(albumsBucket)
		return b.DeleteBucket(itob(albumId))
	})
}

type Image struct {
	Id          uint64 `json:"id"`
	Description string `json:"description"`
}

func (d *Database) GetImages(galleryId, albumId uint64) ([]Image, error) {
	result := make([]Image, 0)

	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(galleryBucket)
		b = b.Bucket(itob(galleryId))
		if b == nil {
			return ErrGalleryNotFound
		}
		b = b.Bucket(albumsBucket)
		b = b.Bucket(itob(albumId))
		if b == nil {
			return ErrAlbumNotFound
		}
		b = b.Bucket(imagesBucket)

		c := b.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			id := btoi(k)
			description := string(b.Bucket(k).Get(descriptionKey))

			result = append(result, Image{
				Id:          id,
				Description: description,
			})
		}

		return nil
	})
	return result, err
}

func (d *Database) AddImage(galleryId, albumId uint64, imageReader io.Reader) (uint64, error) {
	var imgId uint64

	img, _, err := image.Decode(imageReader)
	if err != nil {
		return 0, err
	}

	thumb := resize.Thumbnail(640, 360, img, d.cfg.Interpolation)

	var tBuff, iBuff bytes.Buffer
	err = jpeg.Encode(&tBuff, thumb, &jpeg.Options{Quality: d.cfg.Quality})
	if err != nil {
		return 0, err
	}

	err = jpeg.Encode(&iBuff, img, &jpeg.Options{Quality: d.cfg.Quality})
	if err != nil {
		return 0, err
	}

	err = d.db.Update(func(tx *bolt.Tx) error {
		g := tx.Bucket(galleryBucket)
		b := g.Bucket(itob(galleryId))
		if b == nil {
			return ErrGalleryNotFound
		}
		b = b.Bucket(albumsBucket)

		b = b.Bucket(itob(albumId))
		if b == nil {
			return ErrAlbumNotFound
		}

		imgs := b.Bucket(imagesBucket)

		imgId, err = imgs.NextSequence()
		if err != nil {
			return err
		}

		imgBucket, err := imgs.CreateBucket(itob(imgId))
		if err != nil {
			return err
		}

		err = imgBucket.Put(thumbnailKey, tBuff.Bytes())
		if err != nil {
			return err
		}

		err = imgBucket.Put(imageKey, iBuff.Bytes())
		if err != nil {
			return err
		}

		err = imgBucket.Put(descriptionKey, []byte(""))
		if err != nil {
			return err
		}

		return imgBucket.Put(timestampKey, itob(uint64(time.Now().UnixNano())))
	})

	return imgId, err
}

func (d *Database) SetImageDescription(galleryId, albumId, imageId uint64, description string) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		g := tx.Bucket(galleryBucket)
		b := g.Bucket(itob(galleryId))
		if b == nil {
			return ErrGalleryNotFound
		}
		b = b.Bucket(albumsBucket)
		b = b.Bucket(itob(albumId))
		if b == nil {
			return ErrAlbumNotFound
		}
		b = b.Bucket(imagesBucket)
		i := b.Bucket(itob(imageId))
		if i == nil {
			return ErrImageNotFound
		}
		return i.Put(descriptionKey, []byte(description))
	})
}

func (d *Database) DeleteImage(galleryId, albumId, imageId uint64) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		g := tx.Bucket(galleryBucket)
		b := g.Bucket(itob(galleryId))
		if b == nil {
			return ErrGalleryNotFound
		}
		b = b.Bucket(albumsBucket)
		b = b.Bucket(itob(albumId))
		if b == nil {
			return ErrAlbumNotFound
		}
		b = b.Bucket(imagesBucket)
		return b.DeleteBucket(itob(imageId))
	})
}

func (d *Database) GetImage(galleryId, albumId, imageId uint64) ([]byte, time.Time, error) {
	var img []byte = nil
	timestamp := time.Unix(1, 0)

	err := d.db.View(func(tx *bolt.Tx) error {
		g := tx.Bucket(galleryBucket)
		b := g.Bucket(itob(galleryId))
		if b == nil {
			return ErrGalleryNotFound
		}
		b = b.Bucket(albumsBucket)
		b = b.Bucket(itob(albumId))
		if b == nil {
			return ErrAlbumNotFound
		}
		b = b.Bucket(imagesBucket)
		i := b.Bucket(itob(imageId))
		if i == nil {
			return ErrImageNotFound
		}

		ib := i.Get(imageKey)
		img = make([]byte, len(ib))
		copy(img, ib)

		if mt := i.Get(timestampKey); mt != nil {
			timestamp = time.Unix(0, int64(btoi(mt)))
		}

		return nil
	})
	if err != nil {
		return nil, timestamp, err
	}
	return img, timestamp, nil
}

func (d *Database) GetThumbnail(galleryId, albumId, imageId uint64) ([]byte, time.Time, error) {
	var img []byte = nil
	timestamp := time.Unix(1, 0)
	err := d.db.View(func(tx *bolt.Tx) error {
		g := tx.Bucket(galleryBucket)
		b := g.Bucket(itob(galleryId))
		if b == nil {
			return ErrGalleryNotFound
		}
		b = b.Bucket(albumsBucket)
		b = b.Bucket(itob(albumId))
		if b == nil {
			return ErrAlbumNotFound
		}
		b = b.Bucket(imagesBucket)
		i := b.Bucket(itob(imageId))
		if i == nil {
			return ErrImageNotFound
		}
		ib := i.Get(thumbnailKey)
		img = make([]byte, len(ib))
		copy(img, ib)

		if mt := i.Get(timestampKey); mt != nil {
			timestamp = time.Unix(0, int64(btoi(mt)))
		}

		return nil
	})
	if err != nil {
		return nil, timestamp, err
	}
	return img, timestamp, nil
}
