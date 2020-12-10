package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
)

// Accord holds the schema definition for the Accord entity.
type Accord struct {
	ent.Schema
}

// Fields of the Accord.
func (Accord) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("about").Optional(),
	}
}

// Edges of the Accord.
func (Accord) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("articles", Article.Type).Required(),
	}
}
