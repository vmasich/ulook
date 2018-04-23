package main

import (
	"flag"
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
	pairPtr := flag.String("filter", "\x00\xff", "DB shariding filter range")
	flag.Parse()

	filter := NewFilterPair(*pairPtr)
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
		backc.EConn.Subscribe("lookup", func(subj, reply string, url *schema.LookupURL) {

			if !filter.Match(url.Host) {
				return
			}

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
		backc.EConn.Subscribe("update", func(subj, reply string, url *schema.UpdLookupURL) {
			if !filter.Match(url.Host) {
				return
			}

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

type FilterPair struct {
	EqOrBigger byte
	Less       byte
}

func NewFilterPair(pair string) *FilterPair {

	if len(pair) != 2 {
		log.Fatal("invalid filter pair, %s", pair)
	}
	f := &FilterPair{
		EqOrBigger: pair[0],
		Less:       pair[1],
	}
	log.Infof("accepted host name range: '%c'(%X)* >= HOSTNAME < '%c'(%X)*",
		f.EqOrBigger, f.EqOrBigger, f.Less, f.Less)
	return f
}

func (fp *FilterPair) Match(str string) bool {
	if len(str) == 0 {
		return false
	}
	ch := str[0]
	if ch >= fp.EqOrBigger && ch < fp.Less {
		return true
	}
	return false
}
