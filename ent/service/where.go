// Code generated by entc, DO NOT EDIT.

package service

import (
	"github.com/facebook/ent/dialect/sql"
	"github.com/facebook/ent/dialect/sql/sqlgraph"
	"whoam.xyz/ent/predicate"
)

// ID filters vertices based on their identifier.
func ID(id int) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.In(s.C(FieldID), v...))
	})
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.NotIn(s.C(FieldID), v...))
	})
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// ServiceID applies equality check predicate on the "service_id" field. It's identical to ServiceIDEQ.
func ServiceID(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldServiceID), v))
	})
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// Subject applies equality check predicate on the "subject" field. It's identical to SubjectEQ.
func Subject(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldSubject), v))
	})
}

// Domain applies equality check predicate on the "domain" field. It's identical to DomainEQ.
func Domain(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldDomain), v))
	})
}

// CloneURI applies equality check predicate on the "clone_uri" field. It's identical to CloneURIEQ.
func CloneURI(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCloneURI), v))
	})
}

// ServiceIDEQ applies the EQ predicate on the "service_id" field.
func ServiceIDEQ(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldServiceID), v))
	})
}

// ServiceIDNEQ applies the NEQ predicate on the "service_id" field.
func ServiceIDNEQ(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldServiceID), v))
	})
}

// ServiceIDIn applies the In predicate on the "service_id" field.
func ServiceIDIn(vs ...string) predicate.Service {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Service(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldServiceID), v...))
	})
}

// ServiceIDNotIn applies the NotIn predicate on the "service_id" field.
func ServiceIDNotIn(vs ...string) predicate.Service {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Service(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldServiceID), v...))
	})
}

// ServiceIDGT applies the GT predicate on the "service_id" field.
func ServiceIDGT(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldServiceID), v))
	})
}

// ServiceIDGTE applies the GTE predicate on the "service_id" field.
func ServiceIDGTE(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldServiceID), v))
	})
}

// ServiceIDLT applies the LT predicate on the "service_id" field.
func ServiceIDLT(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldServiceID), v))
	})
}

// ServiceIDLTE applies the LTE predicate on the "service_id" field.
func ServiceIDLTE(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldServiceID), v))
	})
}

// ServiceIDContains applies the Contains predicate on the "service_id" field.
func ServiceIDContains(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldServiceID), v))
	})
}

// ServiceIDHasPrefix applies the HasPrefix predicate on the "service_id" field.
func ServiceIDHasPrefix(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldServiceID), v))
	})
}

// ServiceIDHasSuffix applies the HasSuffix predicate on the "service_id" field.
func ServiceIDHasSuffix(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldServiceID), v))
	})
}

// ServiceIDEqualFold applies the EqualFold predicate on the "service_id" field.
func ServiceIDEqualFold(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldServiceID), v))
	})
}

// ServiceIDContainsFold applies the ContainsFold predicate on the "service_id" field.
func ServiceIDContainsFold(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldServiceID), v))
	})
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldName), v))
	})
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.Service {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Service(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldName), v...))
	})
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.Service {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Service(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldName), v...))
	})
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldName), v))
	})
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldName), v))
	})
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldName), v))
	})
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldName), v))
	})
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldName), v))
	})
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldName), v))
	})
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldName), v))
	})
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldName), v))
	})
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldName), v))
	})
}

// SubjectEQ applies the EQ predicate on the "subject" field.
func SubjectEQ(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldSubject), v))
	})
}

// SubjectNEQ applies the NEQ predicate on the "subject" field.
func SubjectNEQ(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldSubject), v))
	})
}

// SubjectIn applies the In predicate on the "subject" field.
func SubjectIn(vs ...string) predicate.Service {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Service(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldSubject), v...))
	})
}

// SubjectNotIn applies the NotIn predicate on the "subject" field.
func SubjectNotIn(vs ...string) predicate.Service {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Service(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldSubject), v...))
	})
}

// SubjectGT applies the GT predicate on the "subject" field.
func SubjectGT(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldSubject), v))
	})
}

// SubjectGTE applies the GTE predicate on the "subject" field.
func SubjectGTE(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldSubject), v))
	})
}

// SubjectLT applies the LT predicate on the "subject" field.
func SubjectLT(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldSubject), v))
	})
}

// SubjectLTE applies the LTE predicate on the "subject" field.
func SubjectLTE(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldSubject), v))
	})
}

// SubjectContains applies the Contains predicate on the "subject" field.
func SubjectContains(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldSubject), v))
	})
}

// SubjectHasPrefix applies the HasPrefix predicate on the "subject" field.
func SubjectHasPrefix(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldSubject), v))
	})
}

// SubjectHasSuffix applies the HasSuffix predicate on the "subject" field.
func SubjectHasSuffix(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldSubject), v))
	})
}

// SubjectEqualFold applies the EqualFold predicate on the "subject" field.
func SubjectEqualFold(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldSubject), v))
	})
}

// SubjectContainsFold applies the ContainsFold predicate on the "subject" field.
func SubjectContainsFold(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldSubject), v))
	})
}

// DomainEQ applies the EQ predicate on the "domain" field.
func DomainEQ(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldDomain), v))
	})
}

// DomainNEQ applies the NEQ predicate on the "domain" field.
func DomainNEQ(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldDomain), v))
	})
}

// DomainIn applies the In predicate on the "domain" field.
func DomainIn(vs ...string) predicate.Service {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Service(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldDomain), v...))
	})
}

// DomainNotIn applies the NotIn predicate on the "domain" field.
func DomainNotIn(vs ...string) predicate.Service {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Service(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldDomain), v...))
	})
}

// DomainGT applies the GT predicate on the "domain" field.
func DomainGT(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldDomain), v))
	})
}

// DomainGTE applies the GTE predicate on the "domain" field.
func DomainGTE(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldDomain), v))
	})
}

// DomainLT applies the LT predicate on the "domain" field.
func DomainLT(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldDomain), v))
	})
}

// DomainLTE applies the LTE predicate on the "domain" field.
func DomainLTE(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldDomain), v))
	})
}

// DomainContains applies the Contains predicate on the "domain" field.
func DomainContains(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldDomain), v))
	})
}

// DomainHasPrefix applies the HasPrefix predicate on the "domain" field.
func DomainHasPrefix(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldDomain), v))
	})
}

// DomainHasSuffix applies the HasSuffix predicate on the "domain" field.
func DomainHasSuffix(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldDomain), v))
	})
}

// DomainEqualFold applies the EqualFold predicate on the "domain" field.
func DomainEqualFold(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldDomain), v))
	})
}

// DomainContainsFold applies the ContainsFold predicate on the "domain" field.
func DomainContainsFold(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldDomain), v))
	})
}

// CloneURIEQ applies the EQ predicate on the "clone_uri" field.
func CloneURIEQ(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCloneURI), v))
	})
}

// CloneURINEQ applies the NEQ predicate on the "clone_uri" field.
func CloneURINEQ(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCloneURI), v))
	})
}

// CloneURIIn applies the In predicate on the "clone_uri" field.
func CloneURIIn(vs ...string) predicate.Service {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Service(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldCloneURI), v...))
	})
}

// CloneURINotIn applies the NotIn predicate on the "clone_uri" field.
func CloneURINotIn(vs ...string) predicate.Service {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Service(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldCloneURI), v...))
	})
}

// CloneURIGT applies the GT predicate on the "clone_uri" field.
func CloneURIGT(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCloneURI), v))
	})
}

// CloneURIGTE applies the GTE predicate on the "clone_uri" field.
func CloneURIGTE(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCloneURI), v))
	})
}

// CloneURILT applies the LT predicate on the "clone_uri" field.
func CloneURILT(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCloneURI), v))
	})
}

// CloneURILTE applies the LTE predicate on the "clone_uri" field.
func CloneURILTE(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCloneURI), v))
	})
}

// CloneURIContains applies the Contains predicate on the "clone_uri" field.
func CloneURIContains(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldCloneURI), v))
	})
}

// CloneURIHasPrefix applies the HasPrefix predicate on the "clone_uri" field.
func CloneURIHasPrefix(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldCloneURI), v))
	})
}

// CloneURIHasSuffix applies the HasSuffix predicate on the "clone_uri" field.
func CloneURIHasSuffix(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldCloneURI), v))
	})
}

// CloneURIEqualFold applies the EqualFold predicate on the "clone_uri" field.
func CloneURIEqualFold(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldCloneURI), v))
	})
}

// CloneURIContainsFold applies the ContainsFold predicate on the "clone_uri" field.
func CloneURIContainsFold(v string) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldCloneURI), v))
	})
}

// HasMethods applies the HasEdge predicate on the "methods" edge.
func HasMethods() predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(MethodsTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, MethodsTable, MethodsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasMethodsWith applies the HasEdge predicate on the "methods" edge with a given conditions (other predicates).
func HasMethodsWith(preds ...predicate.Method) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(MethodsInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, MethodsTable, MethodsColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups list of predicates with the AND operator between them.
func And(predicates ...predicate.Service) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups list of predicates with the OR operator between them.
func Or(predicates ...predicate.Service) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Service) predicate.Service {
	return predicate.Service(func(s *sql.Selector) {
		p(s.Not())
	})
}
