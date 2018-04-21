package mqsvc

import (
	"testing"

	"bitbucket.org/vmasych/urllookup/pkg/schema"
	"bitbucket.org/vmasych/urllookup/pkg/store/mockstore"
	"github.com/nats-io/nats"
	"github.com/stretchr/testify/assert"
)

var (
	restc = &Nats{
		URL: nats.DefaultURL,
	}
	backend = &Nats{
		URL: nats.DefaultURL,
	}
	ready = make(chan interface{})
)

func printMsg(m *nats.Msg, i int) {
	log.Printf("[#%d] Received on [%s]: '%s'\n", i, m.Subject, string(m.Data))
}

func init() {

	go func() {
		db := &mockstore.MockStore{}
		err := backend.ConnectStore(db)
		if err != nil {
			log.Fatal(err)
		}

		i := 0
		backend.EConn.Subscribe("lookup", func(subj, reply string, url *schema.MyUrl) {
			i++
			//			printMsg(msg, i)
			check := backend.Backend.CheckURL(*url)
			backend.EConn.Publish(reply, check)
		})
		backend.EConn.Flush()

		if err = backend.EConn.LastError(); err != nil {
			log.Fatal(err)
		}
		close(ready)
	}()
}

func TestRestClient(t *testing.T) {
	a := assert.New(t)
	err := restc.ConnectRest()
	a.NoError(err)
}

func TestRestCheck(t *testing.T) {
	defer func() {
		//		time.Sleep(time.Second)
	}()
	<-ready

	dt := []schema.MyUrl{
		{"a.b.c1", "bum"},
		{"a.b.c2", "rum"},
		{"a.b.c3", "bum"},
		{"a.b.c4", "rum"},
		{"a.b.c5", "bum"},
		{"a.b.c6", "rum"},
		{"a.b.c7", "bum"},
		{"a.b.c8", "rum"},
		{"a", "bum"},
	}

	for _, d := range dt {
		go func(u schema.MyUrl) {
			found := restc.CheckURL(u)
			t.Log(u, found)
		}(d)
	}

}
