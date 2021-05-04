package wit

// Delta represents a page change
type Delta interface {
	Apply(document Document)
	MarshalJSON() ([]byte, error)
}
