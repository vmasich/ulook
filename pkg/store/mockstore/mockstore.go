package mockstore

import (
	"bitbucket.org/vmasych/urllookup/pkg/schema"
	"github.com/coreos/pkg/capnslog"
	pretty "github.com/tonnerre/golang-pretty"
)

var log = capnslog.NewPackageLogger(
	"github.com/vmasich/dynweb/pkg/store/mockstore", "mockstore")

type MockStore struct {
}

func (s *MockStore) LookupURL(url schema.LookupURL) (bool, error) {
	log.Infof("CheckUrl %# v", url)
	switch url.Host {
	case "a:80", "a", "a.b.c8":
		switch url.PathQuery {
		case "u/a/b/c", "b":
			return true, nil
		}
		return false, nil
	}
	return false, nil
}

func (s *MockStore) UpdateURLs(urls []schema.UpdLookupURL) error {
	log.Infof("*** UpdateUrls %# v", pretty.Formatter(urls))
	return nil
}

func (s *MockStore) UpdateURL(url schema.UpdLookupURL) error {
	log.Infof("*** UpdateUrl %# v", pretty.Formatter(url))
	return nil
}

func (s *MockStore) Open() error {
	return nil
}

func (s *MockStore) Close() error {
	return nil
}
