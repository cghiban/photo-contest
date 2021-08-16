package photo_test

import (
	"context"
	"photo-contest/business/data/photo"
	"photo-contest/business/data/schema"
	"photo-contest/business/data/tests"
	"photo-contest/business/data/user"
	"testing"
)

func TestPhoto(t *testing.T) {
	log, db, teardown := tests.NewUnit(t)
	t.Cleanup(teardown)

	userStore := user.NewStore(log, db)
	photoStore := photo.NewStore(log, db)

	// lets seed some data (we need a user at least)
	ctx := context.Background()
	if err := schema.Seed(ctx, db); err != nil {
		log.Fatal("Can't seed data:", err)
	}

	t.Log("Given the need to work with Photo records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single photo.", testID)
		{
			user, err := userStore.QueryByID(1)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve user by ID: %s.", tests.Failed, testID, err)
			}

			np := photo.NewPhoto{
				OwnerID:     user.ID,
				Title:       "Test photo 1",
				Description: "Hopa Hopa Penelopa",
				UpdatedBy:   user.Name,
			}
			_, err = photoStore.Create(np)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create photo : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create photo.", tests.Success, testID)

		}
	}
}
