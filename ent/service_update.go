// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"github.com/facebook/ent/dialect/sql"
	"github.com/facebook/ent/dialect/sql/sqlgraph"
	"github.com/facebook/ent/schema/field"
	"whoam.xyz/ent/method"
	"whoam.xyz/ent/predicate"
	"whoam.xyz/ent/service"
)

// ServiceUpdate is the builder for updating Service entities.
type ServiceUpdate struct {
	config
	hooks    []Hook
	mutation *ServiceMutation
}

// Where adds a new predicate for the builder.
func (su *ServiceUpdate) Where(ps ...predicate.Service) *ServiceUpdate {
	su.mutation.predicates = append(su.mutation.predicates, ps...)
	return su
}

// SetName sets the name field.
func (su *ServiceUpdate) SetName(s string) *ServiceUpdate {
	su.mutation.SetName(s)
	return su
}

// SetSubject sets the subject field.
func (su *ServiceUpdate) SetSubject(s string) *ServiceUpdate {
	su.mutation.SetSubject(s)
	return su
}

// SetDomain sets the domain field.
func (su *ServiceUpdate) SetDomain(s string) *ServiceUpdate {
	su.mutation.SetDomain(s)
	return su
}

// SetCloneURI sets the clone_uri field.
func (su *ServiceUpdate) SetCloneURI(s string) *ServiceUpdate {
	su.mutation.SetCloneURI(s)
	return su
}

// AddMethodIDs adds the methods edge to Method by ids.
func (su *ServiceUpdate) AddMethodIDs(ids ...int) *ServiceUpdate {
	su.mutation.AddMethodIDs(ids...)
	return su
}

// AddMethods adds the methods edges to Method.
func (su *ServiceUpdate) AddMethods(m ...*Method) *ServiceUpdate {
	ids := make([]int, len(m))
	for i := range m {
		ids[i] = m[i].ID
	}
	return su.AddMethodIDs(ids...)
}

// Mutation returns the ServiceMutation object of the builder.
func (su *ServiceUpdate) Mutation() *ServiceMutation {
	return su.mutation
}

// ClearMethods clears all "methods" edges to type Method.
func (su *ServiceUpdate) ClearMethods() *ServiceUpdate {
	su.mutation.ClearMethods()
	return su
}

// RemoveMethodIDs removes the methods edge to Method by ids.
func (su *ServiceUpdate) RemoveMethodIDs(ids ...int) *ServiceUpdate {
	su.mutation.RemoveMethodIDs(ids...)
	return su
}

// RemoveMethods removes methods edges to Method.
func (su *ServiceUpdate) RemoveMethods(m ...*Method) *ServiceUpdate {
	ids := make([]int, len(m))
	for i := range m {
		ids[i] = m[i].ID
	}
	return su.RemoveMethodIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (su *ServiceUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(su.hooks) == 0 {
		if err = su.check(); err != nil {
			return 0, err
		}
		affected, err = su.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ServiceMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = su.check(); err != nil {
				return 0, err
			}
			su.mutation = mutation
			affected, err = su.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(su.hooks) - 1; i >= 0; i-- {
			mut = su.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, su.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (su *ServiceUpdate) SaveX(ctx context.Context) int {
	affected, err := su.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (su *ServiceUpdate) Exec(ctx context.Context) error {
	_, err := su.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (su *ServiceUpdate) ExecX(ctx context.Context) {
	if err := su.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (su *ServiceUpdate) check() error {
	if v, ok := su.mutation.Domain(); ok {
		if err := service.DomainValidator(v); err != nil {
			return &ValidationError{Name: "domain", err: fmt.Errorf("ent: validator failed for field \"domain\": %w", err)}
		}
	}
	if v, ok := su.mutation.CloneURI(); ok {
		if err := service.CloneURIValidator(v); err != nil {
			return &ValidationError{Name: "clone_uri", err: fmt.Errorf("ent: validator failed for field \"clone_uri\": %w", err)}
		}
	}
	return nil
}

func (su *ServiceUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   service.Table,
			Columns: service.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: service.FieldID,
			},
		},
	}
	if ps := su.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := su.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: service.FieldName,
		})
	}
	if value, ok := su.mutation.Subject(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: service.FieldSubject,
		})
	}
	if value, ok := su.mutation.Domain(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: service.FieldDomain,
		})
	}
	if value, ok := su.mutation.CloneURI(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: service.FieldCloneURI,
		})
	}
	if su.mutation.MethodsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   service.MethodsTable,
			Columns: []string{service.MethodsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: method.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.RemovedMethodsIDs(); len(nodes) > 0 && !su.mutation.MethodsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   service.MethodsTable,
			Columns: []string{service.MethodsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: method.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.MethodsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   service.MethodsTable,
			Columns: []string{service.MethodsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: method.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, su.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{service.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// ServiceUpdateOne is the builder for updating a single Service entity.
type ServiceUpdateOne struct {
	config
	hooks    []Hook
	mutation *ServiceMutation
}

// SetName sets the name field.
func (suo *ServiceUpdateOne) SetName(s string) *ServiceUpdateOne {
	suo.mutation.SetName(s)
	return suo
}

// SetSubject sets the subject field.
func (suo *ServiceUpdateOne) SetSubject(s string) *ServiceUpdateOne {
	suo.mutation.SetSubject(s)
	return suo
}

// SetDomain sets the domain field.
func (suo *ServiceUpdateOne) SetDomain(s string) *ServiceUpdateOne {
	suo.mutation.SetDomain(s)
	return suo
}

// SetCloneURI sets the clone_uri field.
func (suo *ServiceUpdateOne) SetCloneURI(s string) *ServiceUpdateOne {
	suo.mutation.SetCloneURI(s)
	return suo
}

// AddMethodIDs adds the methods edge to Method by ids.
func (suo *ServiceUpdateOne) AddMethodIDs(ids ...int) *ServiceUpdateOne {
	suo.mutation.AddMethodIDs(ids...)
	return suo
}

// AddMethods adds the methods edges to Method.
func (suo *ServiceUpdateOne) AddMethods(m ...*Method) *ServiceUpdateOne {
	ids := make([]int, len(m))
	for i := range m {
		ids[i] = m[i].ID
	}
	return suo.AddMethodIDs(ids...)
}

// Mutation returns the ServiceMutation object of the builder.
func (suo *ServiceUpdateOne) Mutation() *ServiceMutation {
	return suo.mutation
}

// ClearMethods clears all "methods" edges to type Method.
func (suo *ServiceUpdateOne) ClearMethods() *ServiceUpdateOne {
	suo.mutation.ClearMethods()
	return suo
}

// RemoveMethodIDs removes the methods edge to Method by ids.
func (suo *ServiceUpdateOne) RemoveMethodIDs(ids ...int) *ServiceUpdateOne {
	suo.mutation.RemoveMethodIDs(ids...)
	return suo
}

// RemoveMethods removes methods edges to Method.
func (suo *ServiceUpdateOne) RemoveMethods(m ...*Method) *ServiceUpdateOne {
	ids := make([]int, len(m))
	for i := range m {
		ids[i] = m[i].ID
	}
	return suo.RemoveMethodIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (suo *ServiceUpdateOne) Save(ctx context.Context) (*Service, error) {
	var (
		err  error
		node *Service
	)
	if len(suo.hooks) == 0 {
		if err = suo.check(); err != nil {
			return nil, err
		}
		node, err = suo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ServiceMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = suo.check(); err != nil {
				return nil, err
			}
			suo.mutation = mutation
			node, err = suo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(suo.hooks) - 1; i >= 0; i-- {
			mut = suo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, suo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (suo *ServiceUpdateOne) SaveX(ctx context.Context) *Service {
	node, err := suo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (suo *ServiceUpdateOne) Exec(ctx context.Context) error {
	_, err := suo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (suo *ServiceUpdateOne) ExecX(ctx context.Context) {
	if err := suo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (suo *ServiceUpdateOne) check() error {
	if v, ok := suo.mutation.Domain(); ok {
		if err := service.DomainValidator(v); err != nil {
			return &ValidationError{Name: "domain", err: fmt.Errorf("ent: validator failed for field \"domain\": %w", err)}
		}
	}
	if v, ok := suo.mutation.CloneURI(); ok {
		if err := service.CloneURIValidator(v); err != nil {
			return &ValidationError{Name: "clone_uri", err: fmt.Errorf("ent: validator failed for field \"clone_uri\": %w", err)}
		}
	}
	return nil
}

func (suo *ServiceUpdateOne) sqlSave(ctx context.Context) (_node *Service, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   service.Table,
			Columns: service.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: service.FieldID,
			},
		},
	}
	id, ok := suo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "ID", err: fmt.Errorf("missing Service.ID for update")}
	}
	_spec.Node.ID.Value = id
	if value, ok := suo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: service.FieldName,
		})
	}
	if value, ok := suo.mutation.Subject(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: service.FieldSubject,
		})
	}
	if value, ok := suo.mutation.Domain(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: service.FieldDomain,
		})
	}
	if value, ok := suo.mutation.CloneURI(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: service.FieldCloneURI,
		})
	}
	if suo.mutation.MethodsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   service.MethodsTable,
			Columns: []string{service.MethodsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: method.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.RemovedMethodsIDs(); len(nodes) > 0 && !suo.mutation.MethodsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   service.MethodsTable,
			Columns: []string{service.MethodsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: method.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.MethodsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   service.MethodsTable,
			Columns: []string{service.MethodsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: method.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Service{config: suo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues()
	if err = sqlgraph.UpdateNode(ctx, suo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{service.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return _node, nil
}
