package model

import "bitbucket.org/vmasych/urllookup/pkg/schema"

type URLChecker interface {
	LookupURL(url schema.LookupURL) (bool, error)
}

type URLBulkUpdater interface {
	UpdateURLs(urls []schema.UpdLookupURL) error
}

type URLUpdater interface {
	UpdateURL(url schema.UpdLookupURL) error
}

type Operations interface {
	URLBulkUpdater
	URLChecker
	URLUpdater
}
