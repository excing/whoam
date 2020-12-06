package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
)

// Method holds the schema definition for the Method entity.
type Method struct {
	ent.Schema
}

// Fields of the Method.
func (Method) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("route"),
		field.Enum("scope").Values("public", "private").Default("public"),
	}
}

// Edges of the Method.
func (Method) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", Service.Type).Ref("methods").Required().Unique(),
	}
}
