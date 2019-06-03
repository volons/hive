package models

type Permissions map[string]bool

func (p Permissions) Allowed(feature string) bool {
	return p[feature]
}

func (p Permissions) Set(feature string, val bool) {
	p[feature] = val
}

func (p Permissions) JSON() interface{} {
	out := make(map[string]bool)

	for feature, val := range p {
		if val {
			out[feature] = true
		}
	}

	return out
}
