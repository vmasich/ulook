package rest

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

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
		hpq    string
	}{
		{200, "a/b"},
		{404, "a"},
		{404, ""},
		{200, "a:80/u/a/b/c"},
		{404, "c/d"},
	}
	a := assert.New(t)
	for i, d := range td {
		resp, err := http.Get(fmt.Sprintf("%s/urlinfo/1/%s", baseUrl, d.hpq))
		a.NoError(err)
		_, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			t.Fatal(err)
		}
		a.Equal(d.status, resp.StatusCode, d.hpq)
		t.Logf("+++ %d, %+v, %d", i, d, resp.StatusCode)
	}
}

func TestUpdate(t *testing.T) {
	//	t.Skip()
	a := assert.New(t)
	td := []struct {
		status int
		data   string
	}{
		{
			200, `[["+","a/b"],["+","c/d"],["-","l/u/a/b/c"]]`,
		},
		{
			200, `[["+"],["+","c/d", "u"],["-","l/u/a/b/c"]]`,
		},
		{
			400, `[["+","a":"a"],["+","c/d"],["-","l/u/a/b/c"]]`,
		},
	}
	for i, d := range td {
		//		body, err := json.Marshal(d.data)

		body := []byte(d.data)

		resp, err := http.Post(
			fmt.Sprintf("%s/urlinfo/bulkupdate", baseUrl),
			"application/json",
			bytes.NewBuffer(body),
		)
		a.NoError(err)
		_, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		a.NoError(err)
		a.Equal(d.status, resp.StatusCode)
		t.Logf("+++ %d, %+v, %d", i, string(body), resp.StatusCode)
	}
}
