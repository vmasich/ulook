package main

import (
	"os"
	"os/signal"

	"bitbucket.org/vmasych/urllookup/pkg/config"
	"bitbucket.org/vmasych/urllookup/pkg/conn/nats/mqsvc"
	"bitbucket.org/vmasych/urllookup/pkg/rest"
	"github.com/coreos/pkg/capnslog"
)

var log = capnslog.NewPackageLogger(
	"bitbucket.org/vmasych/urllookup/cmd/restapi/main", "restapi")

func main() {

	capnslog.SetGlobalLogLevel(capnslog.DEBUG)

	// db := &mockstore.MockStore{}
	// db.Open()
	restc := &mqsvc.Nats{
		URL: config.Get().NatsURL,
	}

	if err := restc.ConnectRest(); err != nil {
		log.Fatalf("cannot connect to NATS, %v", err)
	}

	defer func() {
		log.Infof("closing")
		restc.Close()
		//		db.Close()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	rest := rest.New(restc)
	go func() {
		log.Infof("starting api server")
		rest.Run()
	}()
	<-quit
}
