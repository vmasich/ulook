package boltstore

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"bitbucket.org/vmasych/urllookup/pkg/schema"
	bolt "github.com/coreos/bbolt"
	"github.com/coreos/pkg/capnslog"
)

var log = capnslog.NewPackageLogger(
	"bitbucket.org/vmasych/urllookup/pkg/boltstore", "boltstore")

type BoltStore struct {
	Db       *bolt.DB
	Filename string
	Path     string
}

func (b *BoltStore) Open() {
	var err error
	b.Path, err = filepath.Abs(b.Filename)
	dir := filepath.Dir(b.Path)
	os.MkdirAll(dir, 0755)
	log.Infof("OPEN db \"%s\", err: %v", b.Path, err)
	b.Db, err = bolt.Open(b.Path, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (b *BoltStore) Close() {
	err := b.Db.Close()
	log.Infof("CLOSE db \"%s\", err: %v", b.Path, err)
}

func (b *BoltStore) AddURL(url *schema.LookupURL) error {
	log.Infof("ADD URL %+v", url)
	if err := b.Db.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists([]byte(url.Host))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		err = bkt.Put([]byte(url.PathQuery), []byte(""))
		if err != nil {
			log.Errorf("ADD URL %+v, err %v", url, err)
		}
		return nil
	}); err != nil {
		log.Errorf("ADD URL %+v, err %v", url, err)
		return err
	}
	return nil
}

func (b *BoltStore) RemoveURL(url *schema.LookupURL) error {
	log.Infof("REMOVE URL %+v", url)

	if err := b.Db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(url.Host))
		if bkt == nil {
			return nil
		}
		err := bkt.Delete([]byte(url.PathQuery))
		keyN := bkt.Stats().KeyN
		if keyN == 1 {
			tx.DeleteBucket([]byte(url.Host))
		}
		log.Infof("KeyN: %d", keyN)

		return err
	}); err != nil {
		log.Errorf("REMOVE URL %+v, err %v", url, err)
		return err
	}
	return nil
}

func (b *BoltStore) UpdateURL(url schema.UpdLookupURL) error {
	log.Infof("Update Url, %v", url)
	switch url.Operation {
	case "+":
		return b.AddURL(&url.LookupURL)
	case "-":
		return b.RemoveURL(&url.LookupURL)
	default:
		return fmt.Errorf("invalid operation %s", url.Operation)
	}
}

func (b *BoltStore) UpdateURLs(url []schema.UpdLookupURL) error {
	return fmt.Errorf("Not implemented")
}

// retrieve the data
func (b *BoltStore) LookupURL(url schema.LookupURL) (bool, error) {

	log.Infof("LOOKUP URL %+v", url)
	var res bool
	err := b.Db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(url.Host))
		if bkt == nil {
			res = false
			return nil
		}

		val := bkt.Get([]byte(url.PathQuery))
		res = val != nil
		return nil
	})
	if err != nil {
		log.Error(err)
	}
	return res, err
}

func recPrint(tx *bolt.Tx, c *bolt.Cursor, indent int) {
	for k, v := c.First(); k != nil; k, v = c.Next() {
		if v == nil {
			fmt.Printf(strings.Repeat("\t", indent)+"[%s]\n", k)
			newBucket := c.Bucket().Bucket(k)
			if newBucket == nil {
				newBucket = tx.Bucket(k)
			}
			newCursor := newBucket.Cursor()
			recPrint(tx, newCursor, indent+1)
		} else {
			fmt.Printf(strings.Repeat("\t", indent)+"%s\n", k)
			fmt.Printf(strings.Repeat("\t", indent+1)+"%s\n", v)
		}
	}
}

func (b *BoltStore) List() {

	err := b.Db.View(func(tx *bolt.Tx) error {
		c := tx.Cursor()
		recPrint(tx, c, 0)
		return nil
	})
	if err != nil {
		log.Error(err)
	}

}
