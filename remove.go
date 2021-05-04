package wit

// Remove removes matching elements
type Remove struct{}

// Apply applies the delta to the provided elements
func (r Remove) Apply(d Document) {
	for _, node := range d.nodes {
		parent := node.Parent
		if parent != nil {
			parent.RemoveChild(node)
		}
	}
}

// MarshalJSON marshals the delta to JSON format
func (r Remove) MarshalJSON() ([]byte, error) {
	return []byte("[" + removeLabelJSON + "]"), nil
}
