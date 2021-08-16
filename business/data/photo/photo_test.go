package photo_test

import (
	"context"
	"photo-contest/business/data/photo"
	"photo-contest/business/data/schema"
	"photo-contest/business/data/tests"
	"photo-contest/business/data/user"
	"testing"

	"github.com/google/go-cmp/cmp"
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
			// using data in schema/seed.sql
			usr, err := userStore.QueryByID(1)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve user by ID: %s.", tests.Failed, testID, err)
			}

			np := photo.NewPhoto{
				OwnerID:     usr.ID,
				Title:       "Test photo 1",
				Description: "Hopa Hopa Penelopa",
				UpdatedBy:   usr.Name,
			}
			pht, err := photoStore.Create(np)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create photo : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create photo.", tests.Success, testID)

			//-------------------------------------------------------------------------------
			//testID++
			saved, err := photoStore.QueryByID(pht.ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve photo by ID: %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve photo by ID.", tests.Success, testID)
			if diff := cmp.Diff(pht, saved); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get back the same photo. Diff:\n%s", tests.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same photo.", tests.Success, testID)
			//-------------------------------------------------------------------------------
			//}
			testID++
			t.Logf("\tTest %d:\tWhen handling multiple photos.", testID)
			//{
			userPhotos, err := photoStore.QueryByOwnerID(usr.ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve photos by ownerID: %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve photos by ownerID.", tests.Success, testID)

			// should have exactly one pohoto
			if len(userPhotos) != 1 {
				t.Fatalf("\t%s\tTest %d:\tShould have received extactly one photo from this owner: %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tRetrieved exactly one photo from this owner.", tests.Success, testID)

			//testID++
			// should have received the same photo that was added
			if userPhotos[0].ID != pht.ID || userPhotos[0].OwnerID != usr.ID {
				t.Fatalf("\t%s\tTest %d:\tShould have received the same photo from this user: %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tRetrieved exactly the same photo from this owner.", tests.Success, testID)
			//-------------------------------------------------------------------------------
			// expecting errors here..
			userPhotos, err = photoStore.QueryByOwnerID(9213)
			/*fmt.Printf("*****\n%s\n*****\n\n\n", err)
			if errors.Cause(err) != database.ErrNotFound {
				t.Fatalf("\t%s\tTest %d:\tShould NOT retrieve photos by non-existing user: %s.", tests.Failed, testID, err)
			}
			*/
			if len(userPhotos) != 0 {
				t.Fatalf("\t%s\tTest %d:\tShould NOT retrieve photos by non-existing user: %s.", tests.Failed, testID, err)

			}

			t.Logf("\t%s\tTest %d:\tShould NOT retrieve photos by non-existing user.", tests.Success, testID)

			//-------------------------------------------------------------------------------
		}
	}
}
