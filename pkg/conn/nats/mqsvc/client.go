package mqsvc

import (
	"time"

	"bitbucket.org/vmasych/urllookup/pkg/model"
	"bitbucket.org/vmasych/urllookup/pkg/schema"
	"github.com/coreos/pkg/capnslog"
	"github.com/nats-io/nats"
	pretty "github.com/tonnerre/golang-pretty"
)

var log = capnslog.NewPackageLogger(
	"github.com/vmasich/dynweb/pkg/conn/nats/mqsvc", "mqsvc")

type Nats struct {
	URL     string
	EConn   *nats.EncodedConn
	Backend model.MyModel
}

func (n *Nats) ConnectRest() (err error) {
	conn, err := nats.Connect(n.URL)
	if err != nil {
		return
	}
	n.EConn, err = nats.NewEncodedConn(conn, "json")
	return err
}

func (n *Nats) ConnectStore(store model.MyModel) (err error) {
	n.Backend = store
	conn, err := nats.Connect(n.URL)
	n.EConn, err = nats.NewEncodedConn(conn, "json")
	log.Infof("*** NATS Connected to store %# v", pretty.Formatter(n.Backend))
	return err
}

func (n *Nats) Close() (err error) {
	n.EConn.Close()
	return
}

func (n *Nats) UpdateURLs(urls []schema.UpdateMyUrl) bool {
	log.Infof("*** NATS UpdateUrls %# v", pretty.Formatter(urls))
	return true
}

func (n *Nats) CheckURL(url schema.MyUrl) bool {

	subj := "lookup"
	//	payload := []byte(url.Host)
	var resp bool
	err := n.EConn.Request(subj, url, &resp, 100*time.Millisecond)

	//	msg, err := n.EConn.Request(subj, []byte(payload), 100*time.Millisecond)
	if err != nil {
		if n.EConn.LastError() != nil {
			log.Errorf("Request: %v\n", n.EConn.LastError())
		}
		log.Fatalf("Error in Request: %v\n", err)
	}

	// log.Infof("Published [%s] : '%s'\n", subj, payload)
	// log.Infof("Received [%v] : '%s'\n", msg.Subject, string(msg.Data))
	// log.Infof("*** NATS CheckUrl %v", url)

	return resp
}
