package main

import (
	"algorithms/examples/entplay/ent/enttest"
	"algorithms/examples/entplay/ent/proto/entpb"
	"context"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestUserProto(t *testing.T) {
	user := entpb.User{
		Name:         "rotemtam",
		EmailAddress: "rotemtam@example.com",
	}
	if user.GetName() != "rotemtam" {
		t.Fatal("expected user name to be rotemtam")
	}
	if user.GetEmailAddress() != "rotemtam@example.com" {
		t.Fatal("expected email address to be rotemtam@example.com")
	}
}

func TestGet(t *testing.T) {
	// start by initializing an ent client connected to an in memory sqlite instance
	ctx := context.Background()
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	// next, initialize the UserService. Notice we won't be opening an actual port and
	// creating a gRPC server and instead we are just calling the library code directly.
	svc := entpb.NewUserService(client)

	// next, create a user, a category and set that user to be the admin of the category
	user := client.User.Create().
		SetName("rotemtam").
		SetEmailAddress("r@entgo.io").
		SaveX(ctx)

	client.Category.Create().
		SetName("category").
		SetAdmin(user).
		SaveX(ctx)

	// next, retrieve the user without edge information
	get, err := svc.Get(ctx, &entpb.GetUserRequest{
		Id: int64(user.ID),
	})
	if err != nil {
		t.Fatal("failed retrieving the created user", err)
	}
	if len(get.Administered) != 0 {
		t.Fatal("by default edge information is not supposed to be retrieved")
	}

	// next, retrieve the user *WITH* edge information
	get, err = svc.Get(ctx, &entpb.GetUserRequest{
		Id:   int64(user.ID),
		View: entpb.GetUserRequest_WITH_EDGE_IDS,
	})
	if err != nil {
		t.Fatal("failed retrieving the created user", err)
	}
	if len(get.Administered) != 1 {
		t.Fatal("using WITH_EDGE_IDS edges should be returned")
	}
}
