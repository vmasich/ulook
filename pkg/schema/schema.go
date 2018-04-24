package schema

import (
	"strings"
)

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

type UpdatePayload [][2]string

func (p UpdatePayload) Transform() []UpdLookupURL {
	var uu []UpdLookupURL
	for _, p := range p {
		arr := strings.Split(p[1], "/")
		pq := ""
		if len(arr) > 1 {
			pq = strings.Join(arr[1:], "/")
		}
		u := UpdLookupURL{
			Operation: p[0],
			LookupURL: LookupURL{
				Host:      arr[0],
				PathQuery: pq,
			},
		}
		uu = append(uu, u)
	}
	return uu
}
