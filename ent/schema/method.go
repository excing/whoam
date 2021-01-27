package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
	"github.com/facebook/ent/schema/index"
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
		field.String("subject").Optional(),
		field.Enum("scope").Values("public", "private").Default("public"),
	}
}

// Edges of the Method.
func (Method) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", Service.Type).Ref("methods").Required().Unique(),
	}
}

// Indexs of the Method.
func (Method) Indexs() []ent.Index {
	return []ent.Index{
		index.Fields("name").Edges("owner").Unique(),
	}
}
