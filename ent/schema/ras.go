package schema

import (
	"net/url"
	"time"

	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
	"github.com/google/uuid"
)

// RAS holds the schema definition for the RAS entity.
type RAS struct {
	ent.Schema
}

// Fields of the RAS.
func (RAS) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.String("subject"),
		field.JSON("post_uri", &url.URL{}),
		field.JSON("redirect_uri", &url.URL{}),
		field.Enum("state").Values("new", "allowed", "rejected", "abstained", "voided"),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

// Edges of the RAS.
func (RAS) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("organizer", User.Type),
	}
}
