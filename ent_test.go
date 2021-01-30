package main

import (
	"context"
	"testing"

	"whoam.xyz/ent"
	"whoam.xyz/ent/enttest"
	"whoam.xyz/ent/migrate"
)

func CreateClient(t *testing.T) (context.Context, *ent.Client) {
	opts := []enttest.Option{
		enttest.WithOptions(ent.Log(t.Log)),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	}

	// https://godoc.org/github.com/mattn/go-sqlite3#SQLiteDriver.Open
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1", opts...)
	// client := enttest.Open(t, "sqlite3", "file:test.db?_fk=1", opts...)

	if err := client.Schema.Create(context.Background()); err != nil {
		return nil, nil
	}
	return context.Background(), client
}

func TestCreateUser(t *testing.T) {
	ctx, client := CreateClient(t)
	aoli, err := client.User.Create().SetEmail("aoli@example.com").Save(ctx)
	t.Log(aoli, err)
}

func TestCreateService(t *testing.T) {
	ctx, client := CreateClient(t)
	cloud, err := client.Service.Create().
		SetID("cloud.com").
		SetName("cloud service").
		SetSubject("Support file storage, read-write and update services").
		SetDomain("https://xcloud.xzy").
		SetCloneURI("https://github.com/ThreeTenth/Cloud.git").
		Save(ctx)

	t.Log(cloud, err)
}
