package contest_test

import (
	"context"
	"photo-contest/business/data/contest"
	"photo-contest/business/data/photo"
	"photo-contest/business/data/schema"
	"photo-contest/business/data/tests"
	"testing"
	"time"

	"github.com/avelino/slugify"
	"github.com/google/go-cmp/cmp"
)

func TestPhoto(t *testing.T) {
	log, db, teardown := tests.NewUnit(t)
	t.Cleanup(teardown)

	//userStore := user.NewStore(log, db)
	photoStore := photo.NewStore(log, db)
	contestStore := contest.NewStore(log, db)

	// lets seed some data (we need a user at least)
	ctx := context.Background()
	if err := schema.Seed(ctx, db); err != nil {
		log.Fatal("Can't seed data:", err)
	}

	t.Log("Given the need to work with Contest and ContestEntry records.")
	{
		//var pht photo.Photo
		var photos []photo.Photo
		var c contest.Contest
		//var usr user.AuthUser

		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single contest.", testID)
		{
			var err error

			photos, err = photoStore.QueryByOwnerID(2) // user w/ id=2 added via schema/sql/seed.sql
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve user photos: %s.", tests.Failed, testID, err)
			}
			if len(photos) == 0 {
				t.Fatalf("\t%s\tTest %d:\tShould have retrieve at least one photo.", tests.Failed, testID)
			}

			now := time.Now()
			nc := contest.NewContest{
				Title:       "A Contest",
				Description: "Hopa Hopa Penelopa",
				StartDate:   now.Truncate(time.Hour * 24),
				EndDate:     now.Add(time.Hour * 24 * 10).Truncate(time.Hour * 24),
				UpdatedBy:   "Gopher Tester",
			}
			c, err = contestStore.Create(nc)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create contest : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create contest.", tests.Success, testID)

			//-------------------------------------------------------------------------------
			// this should fail - duplicate slug
			_, err = contestStore.Create(nc)
			if err == nil {
				t.Fatalf("\t%s\tTest %d:\tShould NOT be able to add two contest w/ same slug: %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould NOT be able to add two contests w/ same slug.", tests.Success, testID)
			//-------------------------------------------------------------------------------
			// this should fail as well
			nc.EndDate = now.Add(-24 * time.Hour)
			_, err = contestStore.Create(nc)
			if err == nil {
				t.Fatalf("\t%s\tTest %d:\tShould NOT be able to set StartDate after EndDate: %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould NOT be able to set StartDate after EndDate.", tests.Success, testID)
			//-------------------------------------------------------------------------------
			saved, err := contestStore.QueryByID(c.ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to query contest : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to query contest.", tests.Success, testID)
			//log.Printf("saved = %+v\n", saved)
			if diff := cmp.Diff(c, saved); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get back the same contest. Diff:\n%s", tests.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same contest.", tests.Success, testID)

			//-------------------------------------------------------------------------------
			// query contest by slug
			slug := slugify.Slugify(nc.Title)
			saved, err = contestStore.QueryBySlug(slug)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to query contest by slug: %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to query contest by slug.", tests.Success, testID)
			log.Printf("saved = %+v\n", saved)
			if diff := cmp.Diff(c, saved); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get back the same contest. Diff:\n%s", tests.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same contest.", tests.Success, testID)
		}

		testID++
		t.Logf("\tTest %d:\tWhen handling contest entries.", testID)
		{
			//-------------------------------------------------------------------------------
			ncp := contest.NewContestEntry{
				ContestID: c.ID,
				PhotoID:   photos[0].ID,
				Status:    "active",
				UpdatedBy: "Gopher Tester",
			}
			cp, err := contestStore.CreateContestEntry(ncp)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create contest entry: %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create contest entry.", tests.Success, testID)
			//fmt.Printf("cp = %+v\n", cp)
			//-------------------------------------------------------------------------------
			// TODO check contest photo
			cPhotos, err := contestStore.QueryContestEntrys(c.ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to query contest entries: %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to query contest entries.", tests.Success, testID)
			if diff := cmp.Diff(cp, cPhotos[0]); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get back the same contest entries. Diff:\n%s", tests.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same contest entries.", tests.Success, testID)
			//-------------------------------------------------------------------------------
			// TODO update contest photo
			//-------------------------------------------------------------------------------
			//-------------------------------------------------------------------------------

		}
	}
}
