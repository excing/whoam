// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/facebook/ent/dialect/sql"
	"github.com/facebook/ent/dialect/sql/sqlgraph"
	"github.com/facebook/ent/schema/field"
	"whoam.xyz/ent/oauth"
	"whoam.xyz/ent/predicate"
	"whoam.xyz/ent/service"
	"whoam.xyz/ent/user"
)

// OauthUpdate is the builder for updating Oauth entities.
type OauthUpdate struct {
	config
	hooks    []Hook
	mutation *OauthMutation
}

// Where adds a new predicate for the builder.
func (ou *OauthUpdate) Where(ps ...predicate.Oauth) *OauthUpdate {
	ou.mutation.predicates = append(ou.mutation.predicates, ps...)
	return ou
}

// SetExpiredAt sets the expired_at field.
func (ou *OauthUpdate) SetExpiredAt(t time.Time) *OauthUpdate {
	ou.mutation.SetExpiredAt(t)
	return ou
}

// SetOwnerID sets the owner edge to User by id.
func (ou *OauthUpdate) SetOwnerID(id int) *OauthUpdate {
	ou.mutation.SetOwnerID(id)
	return ou
}

// SetOwner sets the owner edge to User.
func (ou *OauthUpdate) SetOwner(u *User) *OauthUpdate {
	return ou.SetOwnerID(u.ID)
}

// SetServiceID sets the service edge to Service by id.
func (ou *OauthUpdate) SetServiceID(id int) *OauthUpdate {
	ou.mutation.SetServiceID(id)
	return ou
}

// SetService sets the service edge to Service.
func (ou *OauthUpdate) SetService(s *Service) *OauthUpdate {
	return ou.SetServiceID(s.ID)
}

// Mutation returns the OauthMutation object of the builder.
func (ou *OauthUpdate) Mutation() *OauthMutation {
	return ou.mutation
}

// ClearOwner clears the "owner" edge to type User.
func (ou *OauthUpdate) ClearOwner() *OauthUpdate {
	ou.mutation.ClearOwner()
	return ou
}

// ClearService clears the "service" edge to type Service.
func (ou *OauthUpdate) ClearService() *OauthUpdate {
	ou.mutation.ClearService()
	return ou
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (ou *OauthUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(ou.hooks) == 0 {
		if err = ou.check(); err != nil {
			return 0, err
		}
		affected, err = ou.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*OauthMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = ou.check(); err != nil {
				return 0, err
			}
			ou.mutation = mutation
			affected, err = ou.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(ou.hooks) - 1; i >= 0; i-- {
			mut = ou.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, ou.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (ou *OauthUpdate) SaveX(ctx context.Context) int {
	affected, err := ou.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (ou *OauthUpdate) Exec(ctx context.Context) error {
	_, err := ou.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ou *OauthUpdate) ExecX(ctx context.Context) {
	if err := ou.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ou *OauthUpdate) check() error {
	if _, ok := ou.mutation.OwnerID(); ou.mutation.OwnerCleared() && !ok {
		return errors.New("ent: clearing a required unique edge \"owner\"")
	}
	if _, ok := ou.mutation.ServiceID(); ou.mutation.ServiceCleared() && !ok {
		return errors.New("ent: clearing a required unique edge \"service\"")
	}
	return nil
}

func (ou *OauthUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   oauth.Table,
			Columns: oauth.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: oauth.FieldID,
			},
		},
	}
	if ps := ou.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ou.mutation.ExpiredAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: oauth.FieldExpiredAt,
		})
	}
	if ou.mutation.OwnerCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   oauth.OwnerTable,
			Columns: []string{oauth.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ou.mutation.OwnerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   oauth.OwnerTable,
			Columns: []string{oauth.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if ou.mutation.ServiceCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   oauth.ServiceTable,
			Columns: []string{oauth.ServiceColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: service.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ou.mutation.ServiceIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   oauth.ServiceTable,
			Columns: []string{oauth.ServiceColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: service.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, ou.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{oauth.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// OauthUpdateOne is the builder for updating a single Oauth entity.
type OauthUpdateOne struct {
	config
	hooks    []Hook
	mutation *OauthMutation
}

// SetExpiredAt sets the expired_at field.
func (ouo *OauthUpdateOne) SetExpiredAt(t time.Time) *OauthUpdateOne {
	ouo.mutation.SetExpiredAt(t)
	return ouo
}

// SetOwnerID sets the owner edge to User by id.
func (ouo *OauthUpdateOne) SetOwnerID(id int) *OauthUpdateOne {
	ouo.mutation.SetOwnerID(id)
	return ouo
}

// SetOwner sets the owner edge to User.
func (ouo *OauthUpdateOne) SetOwner(u *User) *OauthUpdateOne {
	return ouo.SetOwnerID(u.ID)
}

// SetServiceID sets the service edge to Service by id.
func (ouo *OauthUpdateOne) SetServiceID(id int) *OauthUpdateOne {
	ouo.mutation.SetServiceID(id)
	return ouo
}

// SetService sets the service edge to Service.
func (ouo *OauthUpdateOne) SetService(s *Service) *OauthUpdateOne {
	return ouo.SetServiceID(s.ID)
}

// Mutation returns the OauthMutation object of the builder.
func (ouo *OauthUpdateOne) Mutation() *OauthMutation {
	return ouo.mutation
}

// ClearOwner clears the "owner" edge to type User.
func (ouo *OauthUpdateOne) ClearOwner() *OauthUpdateOne {
	ouo.mutation.ClearOwner()
	return ouo
}

// ClearService clears the "service" edge to type Service.
func (ouo *OauthUpdateOne) ClearService() *OauthUpdateOne {
	ouo.mutation.ClearService()
	return ouo
}

// Save executes the query and returns the updated entity.
func (ouo *OauthUpdateOne) Save(ctx context.Context) (*Oauth, error) {
	var (
		err  error
		node *Oauth
	)
	if len(ouo.hooks) == 0 {
		if err = ouo.check(); err != nil {
			return nil, err
		}
		node, err = ouo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*OauthMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = ouo.check(); err != nil {
				return nil, err
			}
			ouo.mutation = mutation
			node, err = ouo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(ouo.hooks) - 1; i >= 0; i-- {
			mut = ouo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, ouo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (ouo *OauthUpdateOne) SaveX(ctx context.Context) *Oauth {
	node, err := ouo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (ouo *OauthUpdateOne) Exec(ctx context.Context) error {
	_, err := ouo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ouo *OauthUpdateOne) ExecX(ctx context.Context) {
	if err := ouo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ouo *OauthUpdateOne) check() error {
	if _, ok := ouo.mutation.OwnerID(); ouo.mutation.OwnerCleared() && !ok {
		return errors.New("ent: clearing a required unique edge \"owner\"")
	}
	if _, ok := ouo.mutation.ServiceID(); ouo.mutation.ServiceCleared() && !ok {
		return errors.New("ent: clearing a required unique edge \"service\"")
	}
	return nil
}

func (ouo *OauthUpdateOne) sqlSave(ctx context.Context) (_node *Oauth, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   oauth.Table,
			Columns: oauth.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: oauth.FieldID,
			},
		},
	}
	id, ok := ouo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "ID", err: fmt.Errorf("missing Oauth.ID for update")}
	}
	_spec.Node.ID.Value = id
	if value, ok := ouo.mutation.ExpiredAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: oauth.FieldExpiredAt,
		})
	}
	if ouo.mutation.OwnerCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   oauth.OwnerTable,
			Columns: []string{oauth.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ouo.mutation.OwnerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   oauth.OwnerTable,
			Columns: []string{oauth.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: user.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if ouo.mutation.ServiceCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   oauth.ServiceTable,
			Columns: []string{oauth.ServiceColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: service.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ouo.mutation.ServiceIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   oauth.ServiceTable,
			Columns: []string{oauth.ServiceColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: service.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Oauth{config: ouo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues()
	if err = sqlgraph.UpdateNode(ctx, ouo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{oauth.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return _node, nil
}