package schema

import "github.com/facebook/ent"

// Accord holds the schema definition for the Accord entity.
type Accord struct {
	ent.Schema
}

// Fields of the Accord.
func (Accord) Fields() []ent.Field {
	return nil
}

// Edges of the Accord.
func (Accord) Edges() []ent.Edge {
	return nil
}
