package main

import (
	"os"
	"os/signal"

	"bitbucket.org/vmasych/urllookup/pkg/rest"
	"bitbucket.org/vmasych/urllookup/pkg/store/mockstore"
	"github.com/coreos/pkg/capnslog"
)

var log = capnslog.NewPackageLogger(
	"bitbucket.org/vmasych/urllookup/cmd/restapi/main", "restmain")

func main() {

	db := &mockstore.MockStore{}
	db.Open()

	defer func() {
		log.Infof("closing")
		db.Close()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	rest := rest.New(db)
	go func() {
		log.Infof("starting api server")
		rest.Run()
	}()
	<-quit
}
