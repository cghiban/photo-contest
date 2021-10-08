package user_test

import (
	"photo-contest/business/data/tests"
	"photo-contest/business/data/user"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUser(t *testing.T) {
	log, db, teardown := tests.NewUnit(t)
	t.Cleanup(teardown)

	store := user.NewStore(log, db)

	t.Log("Given the need to work with User records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single User.", testID)
		{
			nu := user.NewAuthUser{
				Name:        "Cornel Ghiban",
				Email:       "ghiban@cshl.edu",
				Pass:        "HopaHopaPenelopa",
				PassConfirm: "HopaHopaPenelopa",
				Street:      "334 Main St",
				City:        "Cold Spring Harbor",
				State:       "NY",
				Zip:         "11724",
				Phone:       "(516)367-5170",
				Age:         30,
				Gender:      "M",
				Ethnicity:   "wh",
			}
			usr, err := store.Create(nu)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create user : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create user.", tests.Success, testID)

			saved, err := store.QueryByID(usr.ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve user by ID: %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve user by ID.", tests.Success, testID)
			if diff := cmp.Diff(usr, saved); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get back the same user. Diff:\n%s", tests.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same user.", tests.Success, testID)

			uu := user.UpdateAuthUser{
				Name:      "Cornel",
				Email:     "ghiban@cshl.edu",
				Street:    "334 Main St",
				City:      "Cold Spring Harbor",
				State:     "NY",
				Zip:       "11724",
				Phone:     "(516)367-5170",
				Age:       30,
				Gender:    "M",
				Ethnicity: "wh",
			}
			updated_user, err := store.Update(usr.ID, uu)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to update user with ID: %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to update user with ID.", tests.Success, testID)
			if usr.ID != updated_user.ID {
				t.Fatalf("\t%s\tTest %d:\tShould get back the same user id. %d != %d", tests.Failed, testID, usr.ID, updated_user.ID)
			}

			up := user.UpdateAuthUserPass{
				Pass:        "DecentTulipCat",
				PassConfirm: "DecentTulipCat",
			}
			updated_pass_user, err := store.UpdatePass(usr.ID, up)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to update password of user with ID: %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to update password of user with ID.", tests.Success, testID)
			if usr.ID != updated_pass_user.ID {
				t.Fatalf("\t%s\tTest %d:\tShould get back the same user id. %d != %d", tests.Failed, testID, usr.ID, updated_user.ID)
			}
		}
	}
}
