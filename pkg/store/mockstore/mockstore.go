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

func (s *MockStore) CheckURL(url schema.MyUrl) bool {
	log.Infof("CheckUrl %# v", url)
	if url.HostPort == "a" {
		return true
	}
	return false
}

func (s *MockStore) UpdateURLs(urls []schema.UpdateMyUrl) bool {
	log.Infof("*** UpdateUrls %# v", pretty.Formatter(urls))
	return true
}
