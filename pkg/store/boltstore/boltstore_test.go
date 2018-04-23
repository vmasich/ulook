package boltstore

import (
	"testing"

	"bitbucket.org/vmasych/urllookup/pkg/schema"
)

var db = &BoltStore{
	Filename: "../../../dbstore.db",
}

func TestOpen(t *testing.T) {
	db.Open()
}

func TestAddUrl(t *testing.T) {
	dt := []schema.LookupURL{
		{"a.b.c", "bum"},
		{"a.b.c", "rum"},
		{"a.b.c", "bum"},
	}
	for _, d := range dt {
		db.AddURL(&d)
	}
}

func TestLookupUrl(t *testing.T) {
	dt := []schema.LookupURL{
		{"a.b.c", "bum"},
		{"a.b.c", "buc"},
	}
	for _, d := range dt {
		found, _ := db.CheckURL(d)
		t.Log(found)
	}
}

func TestList(t *testing.T) {
	db.List()
}

func TestRemoveUrl(t *testing.T) {
	dt := []schema.LookupURL{
		{"a.b.c", "bum"},
		{"a.b.c", "bum"},
		{"a.b.c", "rum"},
		{"a.b", "bum"},
	}
	for _, d := range dt {
		db.RemoveURL(&d)
	}
}

func TestList1(t *testing.T) {
	db.List()
}

func TestClose(t *testing.T) {
	db.Close()
}
