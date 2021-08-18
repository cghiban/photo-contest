package photo_test

import (
	"context"
	"fmt"
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
		var pht photo.Photo
		var usr user.AuthUser
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single photo.", testID)
		{
			var err error
			// using data in schema/seed.sql
			usr, err = userStore.QueryByID(1)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve user by ID: %s.", tests.Failed, testID, err)
			}

			np := photo.NewPhoto{
				OwnerID:     usr.ID,
				Title:       "Test photo 1",
				Description: "Hopa Hopa Penelopa",
				UpdatedBy:   usr.Name,
			}
			pht, err = photoStore.Create(np)
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
				t.Fatalf("\t%s\tTest %d:\tShould have received exactly one photo from this owner: %s.", tests.Failed, testID, err)
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

		}
		{
			testID++
			t.Logf("\tTest %d:\tWhen updating photo.", testID)
			//-------------------------------------------------------------------------------
			upd := photo.UpdatePhoto{
				Title:     tests.StringPointer("Updated photo title"),
				UpdatedBy: tests.StringPointer("System Updater"),
			}
			if err := photoStore.Update(pht.ID, upd); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to update photo : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to update photo.", tests.Success, testID)

			saved, err := photoStore.QueryByID(pht.ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve photo by ID: %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve photo by ID.", tests.Success, testID)

			if saved.Title != *upd.Title {
				t.Errorf("\t%s\tTest %d:\tShould be able to see updates to Title.", tests.Failed, testID)
				t.Logf("\t\tTest %d:\tGot: %v", testID, saved.Title)
				t.Logf("\t\tTest %d:\tExpected: %v", testID, *upd.Title)
			}
			//-------------------------------------------------------------------------------
			if saved.Description != pht.Description {
				t.Errorf("\t%s\tTest %d:\tShould NOT see updates to Description.", tests.Failed, testID)
				t.Logf("\t\tTest %d:\tGot: %v", testID, saved.Description)
				t.Logf("\t\tTest %d:\tExpected: %v", testID, pht.Description)
			}
			//-------------------------------------------------------------------------------
		}
		{
			testID++
			t.Logf("\tTest %d:\tWhen handling photo files.", testID)
			//-------------------------------------------------------------------------------
			photoSizes := []string{"thumb", "small", "medium", "large"}
			for _, size := range photoSizes {
				npf := photo.NewPhotoFile{
					PhotoID:   pht.ID,
					FilePath:  fmt.Sprintf("/tmp/photo1-%s-%s.jpg", pht.ID, size),
					Size:      "small",
					UpdatedBy: usr.Name,
				}
				_, err := photoStore.CreateFile(npf)
				if err != nil {
					t.Fatalf("\t%s\tTest %d:\tShould be able to create %s photo file: %s.", tests.Failed, testID, size, err)
				}
				t.Logf("\t%s\tTest %d:\tShould be able to create %s photo file.", tests.Success, testID, size)
			}
			//-------------------------------------------------------------------------------
			npf := photo.NewPhotoFile{
				PhotoID:   pht.ID,
				FilePath:  fmt.Sprintf("/tmp/photo1-%s-%s.jpg", pht.ID, "ioio90909090"),
				Size:      "ioio90909090",
				UpdatedBy: usr.Name,
			}
			if _, err := photoStore.CreateFile(npf); err == nil {
				t.Fatalf("\t%s\tTest %d:\tShould NOT be able to create photo file w/ invalid size: %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould NOT be able to create photo file w/ invalid size.", tests.Success, testID)
			//-------------------------------------------------------------------------------
			photoFiles, err := photoStore.QueryPhotoFiles(pht.ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve photo files: %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve photo files.", tests.Success, testID)
			//-------------------------------------------------------------------------------
			if len(photoSizes) != len(photoFiles) {
				t.Fatalf("\t%s\tTest %d:\tShould retrieve the same number of photo files. Got %d insted of %d", tests.Failed, testID, len(photoFiles), len(photoSizes))
			}
			t.Logf("\t%s\tTest %d:\tShould retrieve %d photo files.", tests.Success, testID, len(photoSizes))
			//-------------------------------------------------------------------------------
		}
		{
			testID++
			t.Logf("\tTest %d:\tWhen deleting a photo.", testID)
			//-------------------------------------------------------------------------------
			if err := photoStore.Delete(pht.ID); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to delete photo : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to delete photo.", tests.Success, testID)

			dPht, err := photoStore.QueryByID(pht.ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve photo by ID: %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve photo by ID.", tests.Success, testID)

			if dPht.Deleted != true {
				t.Errorf("\t%s\tTest %d:\tShould be able to see updates to Deleted.", tests.Failed, testID)
				t.Logf("\t\tTest %d:\tGot: %v", testID, dPht.Deleted)
				t.Logf("\t\tTest %d:\tExpected: %v", testID, false)
			}
			//-------------------------------------------------------------------------------
			photoFiles, err := photoStore.QueryPhotoFiles(pht.ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve photo files: %s.", tests.Failed, testID, err)
			}

			if len(photoFiles) > 0 {
				t.Fatalf("\t%s\tTest %d:\tShould NOT be able to retrieve any photo file for this photo.", tests.Failed, testID)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to delete photo files.", tests.Success, testID)

		}
	}
}
