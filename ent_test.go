package main

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"whoam.xyz/ent"
	"whoam.xyz/ent/enttest"
	"whoam.xyz/ent/method"
	"whoam.xyz/ent/migrate"
	"whoam.xyz/ent/oauth"
	"whoam.xyz/ent/permission"
	"whoam.xyz/ent/schema"
	"whoam.xyz/ent/service"
	"whoam.xyz/ent/user"
	"whoam.xyz/ent/vote"
)

func TestVoteRAS(t *testing.T) {
	ctx, client := CreateClient(t)

	ras, err := client.RAS.Create().
		SetSubject("this is a test RAS").
		SetPostURI("https://saynice.whoam.xyz/post/10725").
		SetRedirectURI("https://127.0.0.1:5500").
		Save(ctx)

	if err != nil {
		t.Fatalf("create ras failed %v", err)
	}

	t.Logf("create ras is %v", ras)

	names := make([]*ent.UserCreate, 1000)

	for i := 0; i < len(names); i++ {
		names[i] = client.User.Create().SetEmail(New16BitID() + "@example.com")
	}

	_, err = client.User.CreateBulk(names...).Save(ctx)

	if err != nil {
		t.Fatalf("create users failed %v", err)
	}

	users := make([]int, 10)
	for i := 0; i < len(users); i++ {
		users[i], err = client.User.Query().
			Order(schema.Rand()).
			Limit(1).
			Select(user.FieldID).
			Int(ctx)
		if err != nil {
			t.Logf("query user failed %v", err)
		}
	}

	t.Logf("random users %v", users)

	rasID := ras.ID.String()

	rases := NewBox(4096, 3)
	rases.SetVal(rasID, users)

	var voters []int
	err = rases.Val(rasID, &voters)

	if err != nil {
		t.Fatalf("box hasn't ras %d, failed %v", ras.ID, err)
	}

	voteCreates := make([]*ent.VoteCreate, len(voters))
	rasUUID, _ := uuid.Parse(rasID)
	for i, v := range voters {
		voteCreates[i] = client.Vote.Create().SetState(vote.StateAllowed).SetOwnerID(v).SetDstID(rasUUID)
	}

	votes, err := client.Vote.CreateBulk(voteCreates...).Save(ctx)
	if err != nil {
		t.Fatalf("create votes failed %v", err)
	}

	t.Logf("create votes %v", votes)

	for _, v := range votes {
		u, _ := v.QueryOwner().Only(ctx)
		s, _ := v.QueryDst().Only(ctx)

		t.Logf("user %v, RAS %v", u, s)
	}
}

func TestCreateRAS(t *testing.T) {
	ctx, client := CreateClient(t)

	ras, err := client.RAS.Create().
		SetSubject("this is a test RAS").
		SetPostURI("https://saynice.whoam.xyz/post/10725").
		SetRedirectURI("https://127.0.0.1:5500").
		Save(ctx)

	if err != nil {
		t.Fatalf("create ras failed %v", err)
	}

	t.Logf("create ras is %v", ras)
}

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
		SetID("cloud.com").
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
		SetID("cloud.com").
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
		SetID("cloud.com").
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
		SetID("blog.com").
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
		SetID("cloud.com").
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
		SetID("blog.com").
		SetName("blog service").
		SetDomain("https://xblog.xzy").
		SetCloneURI("https://github.com/excing/BlogZoneServer.git").
		SetSubject("The best ideas can change who we are").
		Save(ctx)

	t.Log(blog, err)

	signingKey := []byte(New16BitID())
	accessToken, err := NewJWTToken(aoli.ID, blog.ID, timeoutAccessToken, signingKey)

	t.Log(accessToken, err)

	mainToken := New64BitID()
	oauth, err := client.Oauth.Create().
		SetMainToken(mainToken).
		SetExpiredAt(time.Now().Add(timeoutAccessToken)).
		SetOwner(aoli).SetService(blog).
		Save(ctx)

	t.Log(oauth, err)

	blogPermissions, err := client.Permission.Create().SetOwner(aoli).SetClient(blog).AddMethods(upload, download).Save(ctx)

	t.Log(blogPermissions, err)

	token, err := FilterJWTToken(accessToken, signingKey)

	t.Log(token, err)

	t.Logf("is blog service? %v", token.Audience == blog.ID)

	queryPermission := client.Permission.Query().
		Where(permission.HasOwnerWith(user.IDEQ(int(token.OtherID)))).
		Where(permission.HasClientWith(service.IDEQ(token.Audience)))

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
