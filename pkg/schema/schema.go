package schema

type MyUrl struct {
	// Host host or host:port
	Host      string `json:"h"`
	PathQuery string `json:"pq"`
}

// UpdateMyUrl used in bulk update
type UpdateMyUrl struct {
	Operation string `json:"op"`
	MyUrl
}
