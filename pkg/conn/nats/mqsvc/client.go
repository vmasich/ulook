package mqsvc

import (
	"fmt"
	"time"

	"bitbucket.org/vmasych/urllookup/pkg/model"
	"bitbucket.org/vmasych/urllookup/pkg/schema"
	"github.com/coreos/pkg/capnslog"
	"github.com/nats-io/nats"
)

var log = capnslog.NewPackageLogger(
	"github.com/vmasich/dynweb/pkg/conn/nats/mqsvc", "mqsvc")

type Nats struct {
	URL     string
	EConn   *nats.EncodedConn
	Backend model.Operations
}

func (n *Nats) ConnectRest() (err error) {
	conn, err := nats.Connect(n.URL)
	if err != nil {
		return
	}
	n.EConn, err = nats.NewEncodedConn(conn, "json")
	log.Infof("*** NATS - API connected, URL: %s", n.URL)

	return err
}

func (n *Nats) ConnectStore(store model.Operations) (err error) {
	n.Backend = store
	log.Infof("%# v", n)

	conn, err := nats.Connect(n.URL)
	if err != nil {
		return err
	}
	n.EConn, err = nats.NewEncodedConn(conn, "json")
	log.Infof("*** NATS - STORE connnected, %# v", n.Backend)
	return err
}

func (n *Nats) Close() (err error) {
	n.EConn.Close()
	return
}

func (n *Nats) UpdateURL(url schema.UpdLookupURL) error {
	return fmt.Errorf("Not implemented")
}

func (n *Nats) UpdateURLs(urls []schema.UpdLookupURL) error {
	subj := "update"
	for _, uu := range urls {
		log.Infof("*** NATS UpdateUrl %v", uu)
		var resp bool
		err := n.EConn.Request(subj, uu, &resp, 100*time.Millisecond)
		if err != nil {
			if n.EConn.LastError() != nil {
				log.Errorf("Request: %v\n", n.EConn.LastError())
			}
			log.Errorf("update request: %v", err)
		}
	}
	return nil
}

func (n *Nats) LookupURL(url schema.LookupURL) (resp bool, err error) {

	subj := "lookup"
	err = n.EConn.Request(subj, url, &resp, 100*time.Millisecond)

	if err != nil {
		if n.EConn.LastError() != nil {
			log.Errorf("Request: %v\n", n.EConn.LastError())
		}
		log.Errorf("lookup request: %v", err)
	}

	// log.Infof("Published [%s] : '%s'\n", subj, payload)
	// log.Infof("Received [%v] : '%s'\n", msg.Subject, string(msg.Data))
	// log.Infof("*** NATS CheckUrl %v", url)

	return
}
