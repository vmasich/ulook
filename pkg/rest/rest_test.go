package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"bitbucket.org/vmasych/urllookup/pkg/schema"
	"bitbucket.org/vmasych/urllookup/pkg/store/mockstore"
	"github.com/stretchr/testify/assert"
)

var (
	rest    *Rest
	baseUrl = "http://localhost:3333"
)

func TestRest(t *testing.T) {
	db := &mockstore.MockStore{}
	rest = New(db)
	go func() {
		rest.Run()
	}()
}

func TestInfo(t *testing.T) {
	td := []struct {
		status int
		host   string
		path   string
	}{
		{200, "a", "b"},
		{404, "c", "d"},
	}
	for i, d := range td {
		resp, err := http.Get(fmt.Sprintf("%s/urlinfo/1/%s/%s", baseUrl, d.host, d.path))
		assert.NoError(t, err)
		_, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("+++ %d, %+v, %d", i, d, resp.StatusCode)
	}
}

func TestUpdate(t *testing.T) {
	//	t.Skip()
	td := []struct {
		status int
		data   []schema.UpdLookupURL
	}{
		{
			200, []schema.UpdLookupURL{
				{"+", schema.LookupURL{"a", "b"}},
				{"+", schema.LookupURL{"c", "d"}},
				{"-", schema.LookupURL{"l", "u"}},
			},
		},
	}
	for i, d := range td {
		body, err := json.Marshal(d.data)
		resp, err := http.Post(
			fmt.Sprintf("%s/urlinfo/bulkupdate", baseUrl),
			"application/json",
			bytes.NewBuffer(body),
		)
		assert.NoError(t, err)
		_, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("+++ %d, %+v, %d", i, string(body), resp.StatusCode)
	}
}
