package schema

import (
	"regexp"
	"time"

	"github.com/facebook/ent"
	"github.com/facebook/ent/dialect/sql"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
	"github.com/google/uuid"
)

// RAS holds the schema definition for the RAS(random anonymous space) entity.
type RAS struct {
	ent.Schema
}

// Fields of the RAS.
func (RAS) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("udpated_at").Default(time.Now).UpdateDefault(time.Now),
		field.String("subject"),
		field.String("post_uri").Match(httpRegexp()).Immutable(),
		field.String("redirect_uri").Match(httpRegexp()).Immutable(),
		field.Enum("state").Values("new", "allowed", "rejected", "abstained", "voided").Nillable().Optional(),
	}
}

// Edges of the RAS.
func (RAS) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("organizer", User.Type).Ref("sessions").Unique(),
		edge.To("votes", Vote.Type),
		edge.To("accord", Accord.Type).Unique().Required(),
	}
}

func httpRegexp() *regexp.Regexp {
	return regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)
}

// Rand order by `rand()`
func Rand() func(*sql.Selector, func(string) bool) {
	return func(s *sql.Selector, check func(string) bool) {
		s.OrderBy("random()")
	}
}
