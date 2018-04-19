package model

import "bitbucket.org/vmasych/urllookup/pkg/schema"

type URLChecker interface {
	CheckURL(url schema.MyUrl) bool
}

type URLUpdater interface {
	UpdateURLs(urls []schema.UpdateMyUrl) bool
}

type MyModel interface {
	URLChecker
	URLUpdater
}
