package userdb

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

type UserInfo struct {
	Username string
	Password string

	// TBC
}

type UserDB struct {
	DBType string
	DBPath string

	UserDB map[string]*UserInfo
	DBLock sync.RWMutex
}

func New(dbpath string) (*UserDB, error) {
	// FIXME: support more db
	if !strings.HasPrefix(dbpath, "file://") {
		return nil, fmt.Errorf("dbpath %s does not support.", dbpath)
	}

	path := dbpath[len("file://"):]
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("Stat dbpath %s error: %s", dbpath, err)
	}

	return &UserDB{
		DBType: "file",
		DBPath: path,
		UserDB: make(map[string]*UserInfo),
	}, nil
}

func (db *UserDB) SyncFromDB() error {
	// FIXME: What should do when removing a user
	db.DBLock.Lock()
	defer db.DBLock.Unlock()

	// FIXME
	f, err := os.Open(db.DBPath)
	if err != nil {
		return err
	}
	defer f.Close()

	r := bufio.NewReader(f)

	next := true
	for next {
		l, err := r.ReadString('\n')
		if err != nil {
			next = false
		}
		l = strings.TrimSpace(l)
		if strings.HasPrefix(l, "#") {
			continue
		}

		field := strings.Split(l, ":")
		// FIXME
		if len(field) != 2 {
			continue
		}
		if ui, ok := db.UserDB[field[0]]; ok {
			ui.Password = field[1]
		} else {
			ui = &UserInfo{
				Username: field[0],
				Password: field[1],
			}
			db.UserDB[field[0]] = ui
		}

		if err != nil {
			break
		}
	}

	return nil
}

func (db *UserDB) GetUserInfo(user string) (*UserInfo, bool) {
	db.DBLock.RLock()
	defer db.DBLock.RUnlock()

	ui, ok := db.UserDB[user]
	return ui, ok
}
