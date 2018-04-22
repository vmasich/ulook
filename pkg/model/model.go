package model

import "bitbucket.org/vmasych/urllookup/pkg/schema"

type URLChecker interface {
	CheckURL(url schema.LURL) (bool, error)
}

type URLUpdater interface {
	UpdateURLs(urls []schema.UpdateURL) error
}

type Operations interface {
	URLChecker
	URLUpdater
}
