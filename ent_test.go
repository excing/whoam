package main

import (
	"context"
	"testing"

	"whoam.xyz/ent"
	"whoam.xyz/ent/enttest"
	"whoam.xyz/ent/method"
	"whoam.xyz/ent/migrate"
	"whoam.xyz/ent/oauth"
	"whoam.xyz/ent/permission"
	"whoam.xyz/ent/service"
	"whoam.xyz/ent/user"
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
		SetServiceID("cloud.com").
		SetName("cloud service").
		SetSubject("Support file storage, read-write and update services").
		SetCloneURI("https://github.com/ThreeTenth/Cloud.git").
		Save(ctx)

	t.Log(cloud, err)
}

func TestCreateMethods(t *testing.T) {
	ctx, client := CreateClient(t)
	upload, err := client.Method.Create().
		SetName("upload").
		SetRoute("/upload").
		Save(ctx)

	t.Log(upload, err)

	download, err := client.Method.Create().
		SetName("download").
		SetRoute("/download").
		SetScope("private").
		Save(ctx)

	t.Log(download, err)

	cloud, err := client.Service.Create().
		SetServiceID("cloud.com").
		SetName("cloud service").
		SetSubject("Support file storage, read-write and update services").
		SetDomain("https://xcloud.xzy").
		SetCloneURI("https://github.com/ThreeTenth/Cloud.git").
		AddMethods(upload, download).
		Save(ctx)

	t.Log(cloud, err)

	methods, err := cloud.QueryMethods().All(ctx)

	t.Log(methods, err)

	service, err := upload.QueryOwner().Only(ctx)

	t.Log(service, err)
}

func TestCreateOAuth(t *testing.T) {
	ctx, client := CreateClient(t)
	aoli, err := client.User.Create().
		SetEmail("aoli@example.com").
		Save(ctx)

	t.Log(aoli, err)

	cloud, err := client.Service.Create().
		SetServiceID("cloud.com").
		SetName("cloud service").
		SetSubject("Support file storage, read-write and update services").
		SetDomain("https://xcloud.xzy").
		SetCloneURI("https://github.com/ThreeTenth/Cloud.git").
		Save(ctx)

	t.Log(cloud, err)

	mainToken := New64BitID()

	oAuth, err := client.Oauth.Create().
		SetMainToken(mainToken).
		SetOwner(aoli).
		SetService(cloud).
		Save(ctx)

	t.Log(oAuth, err)

	owner := client.Oauth.Query().
		Where(oauth.MainTokenEQ(mainToken)).
		QueryOwner().
		FirstX(ctx)

	t.Log(owner, err)

	service := client.Oauth.Query().
		Where(oauth.MainTokenEQ(mainToken)).
		QueryService().
		FirstX(ctx)

	t.Log(service, err)
}

func TestCreatePermission(t *testing.T) {
	ctx, client := CreateClient(t)

	cloud, err := client.Service.Create().
		SetServiceID("cloud.com").
		SetName("cloud service").
		SetSubject("Support file storage, read-write and update services").
		SetDomain("https://xcloud.xzy").
		SetCloneURI("https://github.com/ThreeTenth/Cloud.git").
		Save(ctx)

	t.Log(cloud, err)

	upload, err := client.Method.Create().
		SetName("upload").
		SetRoute("/upload").
		SetOwner(cloud).
		Save(ctx)

	t.Log(upload, err)

	download, err := client.Method.Create().
		SetName("download").
		SetRoute("/download").
		SetScope("private").
		SetOwner(cloud).
		Save(ctx)

	t.Log(download, err)

	methods, err := cloud.QueryMethods().All(ctx)

	t.Log(methods, err)

	blog, err := client.Service.Create().
		SetServiceID("blog.com").
		SetName("blog service").
		SetDomain("https://xblog.xzy").
		SetCloneURI("https://github.com/excing/BlogZoneServer.git").
		SetSubject("The best ideas can change who we are").
		Save(ctx)

	t.Log(blog, err)

	aoli, err := client.User.Create().
		SetEmail("aoli@example.com").
		Save(ctx)

	t.Log(aoli, err)

	blogPremissions, err := client.Permission.Create().
		AddMethods(upload, download).
		SetOwner(aoli).
		SetClient(blog).
		Save(ctx)

	t.Log(blogPremissions, err)

	permissions, err := aoli.QueryPermissions().QueryMethods().All(ctx)

	t.Log(permissions, err)
}

func TestUserLogin(t *testing.T) {
	ctx, client := CreateClient(t)

	cloud, err := client.Service.Create().
		SetServiceID("cloud.com").
		SetName("cloud service").
		SetSubject("Support file storage, read-write and update services").
		SetDomain("https://xcloud.xzy").
		SetCloneURI("https://github.com/ThreeTenth/Cloud.git").
		Save(ctx)

	t.Log(cloud, err)

	upload, err := client.Method.Create().
		SetName("upload").
		SetRoute("/upload").
		SetOwner(cloud).
		Save(ctx)

	t.Log(upload, err)

	download, err := client.Method.Create().
		SetName("download").
		SetRoute("/download").
		SetScope("private").
		SetOwner(cloud).
		Save(ctx)

	t.Log(download, err)

	aoli, err := client.User.Create().
		SetEmail("aoli@example.com").
		Save(ctx)

	t.Log(aoli, err)

	kra, err := client.User.Create().
		SetEmail("kra@example.com").
		Save(ctx)

	t.Log(kra, err)

	blog, err := client.Service.Create().
		SetServiceID("blog.com").
		SetName("blog service").
		SetDomain("https://xblog.xzy").
		SetCloneURI("https://github.com/excing/BlogZoneServer.git").
		SetSubject("The best ideas can change who we are").
		Save(ctx)

	t.Log(blog, err)

	signingKey := []byte(New16BitID())
	accessToken, err := NewJWTToken(aoli.ID, blog.ServiceID, timeoutAccessToken, signingKey)

	t.Log(accessToken, err)

	mainToken := New64BitID()
	oauth, err := client.Oauth.Create().SetMainToken(mainToken).SetOwner(aoli).SetService(blog).Save(ctx)

	t.Log(oauth, err)

	blogPermissions, err := client.Permission.Create().SetOwner(aoli).SetClient(blog).AddMethods(upload, download).Save(ctx)

	t.Log(blogPermissions, err)

	token, err := FilterJWTToken(accessToken, signingKey)

	t.Log(token, err)

	t.Logf("is blog service? %v", token.Audience == blog.ServiceID)

	queryPermission := client.Permission.Query().
		Where(permission.HasOwnerWith(user.IDEQ(int(token.OtherID)))).
		Where(permission.HasClientWith(service.ServiceIDEQ(token.Audience)))

	t.Log(queryPermission.All(ctx))

	hasUpload := queryPermission.
		QueryMethods().
		Where(method.NameEQ(upload.Name)).
		ExistX(ctx)

	t.Logf("has upload permission? %v", hasUpload)

	hasShared := queryPermission.
		QueryMethods().
		Where(method.NameEQ("share")).
		ExistX(ctx)

	t.Logf("has share permission? %v", hasShared)
}
