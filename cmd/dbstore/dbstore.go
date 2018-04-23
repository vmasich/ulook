package main

import (
	"os"
	"os/signal"

	"bitbucket.org/vmasych/urllookup/pkg/config"
	"bitbucket.org/vmasych/urllookup/pkg/conn/nats/mqsvc"
	"bitbucket.org/vmasych/urllookup/pkg/schema"
	"bitbucket.org/vmasych/urllookup/pkg/store/boltstore"
	"github.com/coreos/pkg/capnslog"
)

var log = capnslog.NewPackageLogger(
	"bitbucket.org/vmasych/urllookup/cmd/dbstore/main", "dbstore")

func main() {
	//	db := &mockstore.MockStore{}
	db := &boltstore.BoltStore{
		Filename: config.Get().DbFileName,
	}
	db.Open()
	backc := &mqsvc.Nats{
		URL: config.Get().NatsURL,
	}
	if err := backc.ConnectStore(db); err != nil {
		log.Fatalf("cannot connect to NATS, %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		i := 0
		backc.EConn.Subscribe("lookup", func(subj, reply string, url *schema.LookupURL) {
			i++
			found, err := backc.Backend.LookupURL(*url)
			if err != nil {
				log.Errorf("backend: %v", err)
			}
			backc.EConn.Publish(reply, found)
		})
		backc.EConn.Flush()

		if err := backc.EConn.LastError(); err != nil {
			log.Error(err)
		}

	}()

	go func() {
		i := 0
		backc.EConn.Subscribe("update", func(subj, reply string, url *schema.UpdLookupURL) {
			i++
			err := backc.Backend.UpdateURL(*url)
			if err != nil {
				log.Errorf("backend: %v", err)
			}
			backc.EConn.Publish(reply, err)
		})
		backc.EConn.Flush()

		if err := backc.EConn.LastError(); err != nil {
			log.Error(err)
		}

	}()

	<-quit
}
