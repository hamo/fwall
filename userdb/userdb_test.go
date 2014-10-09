package userdb

import (
	"testing"
)

var db DB

func Test_NewFileDB(t *testing.T) {
	var err error
	db, err = NewDB("file://./test.db")
	if err != nil {
		t.Errorf("new db failed: %s", err)
	}
}

func Test_SyncFromDB(t *testing.T) {
	if err := db.SyncFromDB(); err != nil {
		t.Errorf("SyncFromDB failed: %s", err)
	}
}

func Test_GetUserInfo(t *testing.T) {
	//a
	if ui, ok := db.GetUserInfo("a"); ok {
		t.Errorf("commented out key a has result: %v", ui)
	}

	//b
	if ui, ok := db.GetUserInfo("b"); !ok {
		t.Errorf("Can not get b")
	} else {
		if ui.Password != "b" {
			t.Errorf("Key b get wrong result: %v", ui)
		}
	}

	//c
	if ui, ok := db.GetUserInfo("c"); !ok {
		t.Errorf("Can not get c")
	} else {
		if ui.Password != "c" {
			t.Errorf("Key c get wrong result: %v", ui)
		}
	}

	//d
	if ui, ok := db.GetUserInfo("d"); ok {
		t.Errorf("commented out key d has result: %v", ui)
	}
}
