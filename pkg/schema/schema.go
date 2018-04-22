package schema

type LURL struct {
	// Host host or host:port
	Host      string `json:"h"`
	PathQuery string `json:"pq"`
}

// UpdateURL used in bulk update
type UpdateURL struct {
	Operation string `json:"op"`
	LURL
}
