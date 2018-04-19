package schema

type MyUrl struct {
	HostPort  string `json:"h"`
	PathQuery string `json:"p"`
}

// UpdateMyUrl used in bulk update
type UpdateMyUrl struct {
	Operation string `json:"op"`
	MyUrl
}
