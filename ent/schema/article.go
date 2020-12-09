package schema

import "github.com/facebook/ent"

// Article holds the schema definition for the Article entity.
type Article struct {
	ent.Schema
}

// Fields of the Article.
func (Article) Fields() []ent.Field {
	return nil
}

// Edges of the Article.
func (Article) Edges() []ent.Edge {
	return nil
}
