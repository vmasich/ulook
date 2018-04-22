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

func (s *MockStore) CheckURL(url schema.LURL) (bool, error) {
	log.Debugf("CheckUrl %# v", url)
	switch url.Host {
	case "a", "a.b.c8":
		return true, nil
	}
	return false, nil
}

func (s *MockStore) UpdateURLs(urls []schema.UpdateURL) error {
	log.Debugf("*** UpdateUrls %# v", pretty.Formatter(urls))
	return nil
}

func (s *MockStore) Open() error {
	return nil
}

func (s *MockStore) Close() error {
	return nil
}
