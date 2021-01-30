package schema

import (
	"time"

	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
)

// Oauth holds the schema definition for the Oauth entity.
type Oauth struct {
	ent.Schema
}

// Fields of the Oauth.
func (Oauth) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("expired_at"),
		field.String("main_token").Immutable().Unique().NotEmpty(),
	}
}

// Edges of the Oauth.
func (Oauth) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("oauths").Required().Unique(),
		edge.To("service", Service.Type).Required().Unique(),
	}
}
