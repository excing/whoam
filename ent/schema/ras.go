package schema

import "github.com/facebook/ent"

// RAS holds the schema definition for the RAS entity.
type RAS struct {
	ent.Schema
}

// Fields of the RAS.
func (RAS) Fields() []ent.Field {
	return nil
}

// Edges of the RAS.
func (RAS) Edges() []ent.Edge {
	return nil
}
