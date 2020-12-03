// Code generated by entc, DO NOT EDIT.

package ras

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the ras type in the database.
	Label = "ras"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldSubject holds the string denoting the subject field in the database.
	FieldSubject = "subject"
	// FieldPostURI holds the string denoting the post_uri field in the database.
	FieldPostURI = "post_uri"
	// FieldRedirectURI holds the string denoting the redirect_uri field in the database.
	FieldRedirectURI = "redirect_uri"
	// FieldState holds the string denoting the state field in the database.
	FieldState = "state"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"

	// Table holds the table name of the ras in the database.
	Table = "ra_ss"
)

// Columns holds all SQL columns for ras fields.
var Columns = []string{
	FieldID,
	FieldSubject,
	FieldPostURI,
	FieldRedirectURI,
	FieldState,
	FieldCreatedAt,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultCreatedAt holds the default value on creation for the created_at field.
	DefaultCreatedAt func() time.Time
	// DefaultID holds the default value on creation for the id field.
	DefaultID func() uuid.UUID
)

// State defines the type for the state enum field.
type State string

// State values.
const (
	StateNew       State = "new"
	StateAllowed   State = "allowed"
	StateRejected  State = "rejected"
	StateAbstained State = "abstained"
	StateVoided    State = "voided"
)

func (s State) String() string {
	return string(s)
}

// StateValidator is a validator for the "state" field enum values. It is called by the builders before save.
func StateValidator(s State) error {
	switch s {
	case StateNew, StateAllowed, StateRejected, StateAbstained, StateVoided:
		return nil
	default:
		return fmt.Errorf("ras: invalid enum value for state field: %q", s)
	}
}
