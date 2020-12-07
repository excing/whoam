// Code generated by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"github.com/facebook/ent/dialect/sql"
	"whoam.xyz/ent/oauth"
	"whoam.xyz/ent/service"
	"whoam.xyz/ent/user"
)

// Oauth is the model entity for the Oauth schema.
type Oauth struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// ExpiredAt holds the value of the "expired_at" field.
	ExpiredAt time.Time `json:"expired_at,omitempty"`
	// MainToken holds the value of the "main_token" field.
	MainToken string `json:"main_token,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the OauthQuery when eager-loading is set.
	Edges         OauthEdges `json:"edges"`
	oauth_service *int
	user_oauths   *int
}

// OauthEdges holds the relations/edges for other nodes in the graph.
type OauthEdges struct {
	// Owner holds the value of the owner edge.
	Owner *User
	// Service holds the value of the service edge.
	Service *Service
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// OwnerOrErr returns the Owner value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e OauthEdges) OwnerOrErr() (*User, error) {
	if e.loadedTypes[0] {
		if e.Owner == nil {
			// The edge owner was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: user.Label}
		}
		return e.Owner, nil
	}
	return nil, &NotLoadedError{edge: "owner"}
}

// ServiceOrErr returns the Service value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e OauthEdges) ServiceOrErr() (*Service, error) {
	if e.loadedTypes[1] {
		if e.Service == nil {
			// The edge service was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: service.Label}
		}
		return e.Service, nil
	}
	return nil, &NotLoadedError{edge: "service"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Oauth) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},  // id
		&sql.NullTime{},   // created_at
		&sql.NullTime{},   // expired_at
		&sql.NullString{}, // main_token
	}
}

// fkValues returns the types for scanning foreign-keys values from sql.Rows.
func (*Oauth) fkValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{}, // oauth_service
		&sql.NullInt64{}, // user_oauths
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Oauth fields.
func (o *Oauth) assignValues(values ...interface{}) error {
	if m, n := len(values), len(oauth.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	o.ID = int(value.Int64)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field created_at", values[0])
	} else if value.Valid {
		o.CreatedAt = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field expired_at", values[1])
	} else if value.Valid {
		o.ExpiredAt = value.Time
	}
	if value, ok := values[2].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field main_token", values[2])
	} else if value.Valid {
		o.MainToken = value.String
	}
	values = values[3:]
	if len(values) == len(oauth.ForeignKeys) {
		if value, ok := values[0].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field oauth_service", value)
		} else if value.Valid {
			o.oauth_service = new(int)
			*o.oauth_service = int(value.Int64)
		}
		if value, ok := values[1].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field user_oauths", value)
		} else if value.Valid {
			o.user_oauths = new(int)
			*o.user_oauths = int(value.Int64)
		}
	}
	return nil
}

// QueryOwner queries the owner edge of the Oauth.
func (o *Oauth) QueryOwner() *UserQuery {
	return (&OauthClient{config: o.config}).QueryOwner(o)
}

// QueryService queries the service edge of the Oauth.
func (o *Oauth) QueryService() *ServiceQuery {
	return (&OauthClient{config: o.config}).QueryService(o)
}

// Update returns a builder for updating this Oauth.
// Note that, you need to call Oauth.Unwrap() before calling this method, if this Oauth
// was returned from a transaction, and the transaction was committed or rolled back.
func (o *Oauth) Update() *OauthUpdateOne {
	return (&OauthClient{config: o.config}).UpdateOne(o)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (o *Oauth) Unwrap() *Oauth {
	tx, ok := o.config.driver.(*txDriver)
	if !ok {
		panic("ent: Oauth is not a transactional entity")
	}
	o.config.driver = tx.drv
	return o
}

// String implements the fmt.Stringer.
func (o *Oauth) String() string {
	var builder strings.Builder
	builder.WriteString("Oauth(")
	builder.WriteString(fmt.Sprintf("id=%v", o.ID))
	builder.WriteString(", created_at=")
	builder.WriteString(o.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", expired_at=")
	builder.WriteString(o.ExpiredAt.Format(time.ANSIC))
	builder.WriteString(", main_token=")
	builder.WriteString(o.MainToken)
	builder.WriteByte(')')
	return builder.String()
}

// Oauths is a parsable slice of Oauth.
type Oauths []*Oauth

func (o Oauths) config(cfg config) {
	for _i := range o {
		o[_i].config = cfg
	}
}
