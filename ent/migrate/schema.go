// Code generated by entc, DO NOT EDIT.

package migrate

import (
	"github.com/facebook/ent/dialect/sql/schema"
	"github.com/facebook/ent/schema/field"
)

var (
	// RaSsColumns holds the columns for the "ra_ss" table.
	RaSsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID},
		{Name: "subject", Type: field.TypeString},
		{Name: "post_uri", Type: field.TypeJSON},
		{Name: "redirect_uri", Type: field.TypeJSON},
		{Name: "state", Type: field.TypeEnum, Enums: []string{"new", "allowed", "rejected", "abstained", "voided"}},
		{Name: "created_at", Type: field.TypeTime},
	}
	// RaSsTable holds the schema information for the "ra_ss" table.
	RaSsTable = &schema.Table{
		Name:        "ra_ss",
		Columns:     RaSsColumns,
		PrimaryKey:  []*schema.Column{RaSsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		RaSsTable,
	}
)

func init() {
}
