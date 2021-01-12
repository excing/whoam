package schema

import (
	"regexp"

	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
)

// Service holds the schema definition for the Service entity.
type Service struct {
	ent.Schema
}

// Fields of the Service.
func (Service) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			NotEmpty().
			Unique().
			Immutable(),
		field.String("name"),
		field.String("subject"),
		field.String("domain").Match(httpRegexp()),
		field.String("clone_uri").Match(regexp.MustCompile(`((git|ssh|http(s)?)|(git@[\w\.]+))(:(//)?)([\w\.@\:/\-~]+)(\.git)(/)?`)),
	}
}

// Edges of the Service.
func (Service) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("methods", Method.Type),
		edge.To("permissions", Method.Type),
	}
}
