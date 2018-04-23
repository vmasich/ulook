package schema

type LookupURL struct {
	// Host host or host:port
	Host      string `json:"h"`
	PathQuery string `json:"pq"`
}

// UpdLookupURL used in bulk update
type UpdLookupURL struct {
	Operation string `json:"op"`
	LookupURL
}
