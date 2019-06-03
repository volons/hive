package models

type Caps map[string]Def

type Def map[string]string

// Supports checks if a feature is present in the capabilities
func (caps Caps) Supports(feature string) bool {
	_, ok := caps[feature]
	return ok
}
