package mqsvc

import (
	"testing"
	"time"

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
	ready    = make(chan interface{})
	complete = make(chan interface{})
)

func init() {
	db := &mockstore.MockStore{}
	err := backend.ConnectStore(db)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		i := 0
		backend.EConn.Subscribe("lookup", func(subj, reply string, url *schema.LURL) {
			i++
			found, _ := backend.Backend.CheckURL(*url)

			backend.EConn.Publish(reply, found)
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
		time.Sleep(200 * time.Millisecond)
	}()

	<-ready

	dt := []schema.LURL{
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
		go func(u schema.LURL) {
			found, _ := restc.CheckURL(u)
			t.Log(u, found)
		}(d)
	}
}

func TestClose(t *testing.T) {

	restc.Close()
	backend.Close()
}
