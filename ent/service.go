// Code generated by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"github.com/facebook/ent/dialect/sql"
	"whoam.xyz/ent/service"
)

// Service is the model entity for the Service schema.
type Service struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Subject holds the value of the "subject" field.
	Subject string `json:"subject,omitempty"`
	// Domain holds the value of the "domain" field.
	Domain string `json:"domain,omitempty"`
	// CloneURI holds the value of the "clone_uri" field.
	CloneURI string `json:"clone_uri,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ServiceQuery when eager-loading is set.
	Edges ServiceEdges `json:"edges"`
}

// ServiceEdges holds the relations/edges for other nodes in the graph.
type ServiceEdges struct {
	// Methods holds the value of the methods edge.
	Methods []*Method
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// MethodsOrErr returns the Methods value or an error if the edge
// was not loaded in eager-loading.
func (e ServiceEdges) MethodsOrErr() ([]*Method, error) {
	if e.loadedTypes[0] {
		return e.Methods, nil
	}
	return nil, &NotLoadedError{edge: "methods"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Service) scanValues() []interface{} {
	return []interface{}{
		&sql.NullString{}, // id
		&sql.NullString{}, // name
		&sql.NullString{}, // subject
		&sql.NullString{}, // domain
		&sql.NullString{}, // clone_uri
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Service fields.
func (s *Service) assignValues(values ...interface{}) error {
	if m, n := len(values), len(service.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	if value, ok := values[0].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field id", values[0])
	} else if value.Valid {
		s.ID = value.String
	}
	values = values[1:]
	if value, ok := values[0].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field name", values[0])
	} else if value.Valid {
		s.Name = value.String
	}
	if value, ok := values[1].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field subject", values[1])
	} else if value.Valid {
		s.Subject = value.String
	}
	if value, ok := values[2].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field domain", values[2])
	} else if value.Valid {
		s.Domain = value.String
	}
	if value, ok := values[3].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field clone_uri", values[3])
	} else if value.Valid {
		s.CloneURI = value.String
	}
	return nil
}

// QueryMethods queries the methods edge of the Service.
func (s *Service) QueryMethods() *MethodQuery {
	return (&ServiceClient{config: s.config}).QueryMethods(s)
}

// Update returns a builder for updating this Service.
// Note that, you need to call Service.Unwrap() before calling this method, if this Service
// was returned from a transaction, and the transaction was committed or rolled back.
func (s *Service) Update() *ServiceUpdateOne {
	return (&ServiceClient{config: s.config}).UpdateOne(s)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (s *Service) Unwrap() *Service {
	tx, ok := s.config.driver.(*txDriver)
	if !ok {
		panic("ent: Service is not a transactional entity")
	}
	s.config.driver = tx.drv
	return s
}

// String implements the fmt.Stringer.
func (s *Service) String() string {
	var builder strings.Builder
	builder.WriteString("Service(")
	builder.WriteString(fmt.Sprintf("id=%v", s.ID))
	builder.WriteString(", name=")
	builder.WriteString(s.Name)
	builder.WriteString(", subject=")
	builder.WriteString(s.Subject)
	builder.WriteString(", domain=")
	builder.WriteString(s.Domain)
	builder.WriteString(", clone_uri=")
	builder.WriteString(s.CloneURI)
	builder.WriteByte(')')
	return builder.String()
}

// Services is a parsable slice of Service.
type Services []*Service

func (s Services) config(cfg config) {
	for _i := range s {
		s[_i].config = cfg
	}
}
